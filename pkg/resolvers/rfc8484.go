package resolvers

import (
	"context"
	"ech-injector/pkg/workers"
	"net/http"

	"github.com/miekg/dns"
)

const hostRFC8484 = "dns.google"
const pathRFC8484 = "/dns-query"

var client *dns.Client

func init() {
	client = new(dns.Client)
}

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
		Method:  http.MethodPost,
		Scheme:  "https",
		Host:    hostRFC8484,
		Path:    pathRFC8484,
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

func ExchangeA(context context.Context, name string) ([]string, error) {
	dnsQuestion := new(dns.Msg)
	dnsQuestion.SetQuestion(name, dns.TypeA)

	dnsAnswer, err := ExchangeMessage(context, dnsQuestion)
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, resourceRecord := range dnsAnswer.Answer {
		a, ok := resourceRecord.(*dns.A)
		if !ok {
			continue
		}

		ip := a.A.String()
		ips = append(ips, ip)
	}

	return ips, nil
}

func ExchangeAAAA(context context.Context, name string) ([]string, error) {
	dnsQuestion := new(dns.Msg)
	dnsQuestion.SetQuestion(name, dns.TypeAAAA)

	dnsAnswer, err := ExchangeMessage(context, dnsQuestion)
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, resourceRecord := range dnsAnswer.Answer {
		aaaa, ok := resourceRecord.(*dns.AAAA)
		if !ok {
			continue
		}

		ip := aaaa.AAAA.String()
		ips = append(ips, ip)
	}

	return ips, nil
}

func ExchangeHTTPS(context context.Context, name string) ([][]dns.SVCBKeyValue, error) {
	dnsQuestion := new(dns.Msg)
	dnsQuestion.SetQuestion(name, dns.TypeHTTPS)

	dnsAnswer, err := ExchangeMessage(context, dnsQuestion)
	if err != nil {
		return nil, err
	}

	var svcbs [][]dns.SVCBKeyValue
	for _, resourceRecord := range dnsAnswer.Answer {
		https, ok := resourceRecord.(*dns.HTTPS)
		if !ok {
			continue
		}

		svcbs = append(svcbs, https.Value)
	}

	return svcbs, nil
}
