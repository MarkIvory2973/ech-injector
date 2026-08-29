package resolvers

import (
	"context"
	"ech-injector/pkg/workers"

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
		return nil, err
	}

	dnsAnswer, err := UnpackMessage(content)
	if err != nil {
		return nil, err
	}

	return dnsAnswer, nil
}
