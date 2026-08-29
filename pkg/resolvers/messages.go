package resolvers

import (
	"context"
	"crypto/sha256"
	"ech-injector/pkg/workers"
	"encoding/base64"
	"encoding/hex"

	"github.com/miekg/dns"
)

func PackMessage(dnsMessage *dns.Msg) ([]byte, error) {
	content, err := dnsMessage.Pack()
	if err != nil {
		return nil, err
	}

	return content, nil
}

func UnpackMessage(content []byte) (*dns.Msg, error) {
	dnsMessage := new(dns.Msg)
	err := dnsMessage.Unpack(content)
	if err != nil {
		return nil, err
	}

	return dnsMessage, nil
}

func ExchangeMessage(context context.Context, dnsQuestion *dns.Msg) (*dns.Msg, error) {
	id := dnsQuestion.Id
	dnsQuestion.Id = 0

	content, err := PackMessage(dnsQuestion)
	if err != nil {
		return nil, err
	}

	checksum := sha256.Sum256(content)
	key := hex.EncodeToString(checksum[:])
	cache, err := workers.CacheFunc(key, func() (map[string]string, int, error) {
		httpRequest := workers.HTTPRequest{
			Method:  "POST",
			Scheme:  "https",
			Host:    host,
			Path:    path,
			Queries: map[string]string{},
			Headers: map[string]string{
				"Content-Type": "application/dns-message",
				"Accept":       "application/dns-message",
			},
			Content: content,
		}
		content, err = workers.Fetch(context, httpRequest)
		if err != nil {
			return nil, 0, err
		}

		cache := map[string]string{
			"answer": base64.StdEncoding.EncodeToString(content),
		}

		dnsAnwser, err := UnpackMessage(content)
		if err != nil {
			return nil, 0, err
		}

		var minTTL int
		for _, resourceRecord := range dnsAnwser.Answer {
			ttl := int(resourceRecord.Header().Ttl)
			if minTTL == 0 || ttl < minTTL {
				minTTL = ttl
			}
		}

		return cache, minTTL, nil
	})
	if err != nil {
		return nil, err
	}

	content, err = base64.StdEncoding.DecodeString(cache["answer"])
	if err != nil {
		return nil, err
	}

	dnsAnswer, err := UnpackMessage(content)
	if err != nil {
		return nil, err
	}

	dnsAnswer.Id = id

	return dnsAnswer, nil
}
