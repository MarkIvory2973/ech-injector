package resolvers

import (
	"context"

	"github.com/miekg/dns"
)

const host = "dns.google"
const path = "/dns-query"

var client *dns.Client

func init() {
	client = new(dns.Client)
}

func ResolveA(context context.Context, name string) ([]string, error) {
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

func ResolveAAAA(context context.Context, name string) ([]string, error) {
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

func ResolveHTTPS(context context.Context, name string) ([][]dns.SVCBKeyValue, error) {
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
