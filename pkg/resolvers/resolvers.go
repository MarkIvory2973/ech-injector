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

	dnsAnwser, err := ExchangeMessage(context, dnsQuestion)
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, resourceRecord := range dnsAnwser.Answer {
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

	dnsAnwser, err := ExchangeMessage(context, dnsQuestion)
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, resourceRecord := range dnsAnwser.Answer {
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

	dnsAnwser, err := ExchangeMessage(context, dnsQuestion)
	if err != nil {
		return nil, err
	}

	var svcbs [][]dns.SVCBKeyValue
	for _, resourceRecord := range dnsAnwser.Answer {
		https, ok := resourceRecord.(*dns.HTTPS)
		if !ok {
			continue
		}

		svcbs = append(svcbs, https.Value)
	}

	return svcbs, nil
}
