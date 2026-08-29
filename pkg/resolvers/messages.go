package resolvers

import (
	"context"
	"crypto/sha256"
	"ech-injector/pkg/workers"
	"encoding/base64"
	"encoding/hex"
	"time"

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
	content, err := PackMessage(dnsQuestion)
	if err != nil {
		return nil, err
	}

	checksum := sha256.Sum256(content)
	key := hex.EncodeToString(checksum[:])
	cached, err := workers.CacheFunc(key, func() (map[string]string, int, error) {
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

		cached := map[string]string{
			"response": base64.StdEncoding.EncodeToString(content),
		}

		dnsAnwser, err := UnpackMessage(content)
		if err != nil {
			return nil, 0, err
		}

		var minTTL time.Duration
		for _, resourceRecord := range dnsAnwser.Answer {
			ttl := time.Duration(resourceRecord.Header().Ttl) * time.Second
			if minTTL == 0 || ttl < minTTL {
				minTTL = ttl
			}
		}

		return cached, int(minTTL), nil
	})

	content, err = base64.StdEncoding.DecodeString(cached["response"])
	if err != nil {
		return nil, err
	}

	dnsAnswer, err := UnpackMessage(content)
	if err != nil {
		return nil, err
	}

	return dnsAnswer, nil
}
