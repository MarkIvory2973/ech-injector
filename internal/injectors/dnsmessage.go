package injectors

import (
	"context"
	"ech-injector/pkg/resolvers"
	"slices"

	"github.com/miekg/dns"
)

func InjectDNSMessage(context context.Context, content []byte) ([]byte, error) {
	dnsQuestion, err := resolvers.UnpackMessage(content)
	if err != nil {
		return nil, err
	}

	dnsAnswer, err := resolvers.ExchangeMessage(context, dnsQuestion)
	if err != nil {
		return nil, err
	}

	for _, question := range dnsQuestion.Question {
		if question.Qtype != dns.TypeHTTPS {
			continue
		}

		exists := false
		for _, resourceRecord := range dnsAnswer.Answer {
			https, ok := resourceRecord.(*dns.HTTPS)
			if !ok {
				continue
			}

			if question.Name != https.Hdr.Name ||
				question.Qclass != https.Hdr.Class ||
				question.Qtype != https.Hdr.Rrtype {
				continue
			}

			exists = true
		}

		if !exists {
			resourceRecord := &dns.SVCB{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeHTTPS,
					Class:  dns.ClassINET,
					Ttl:    300,
				},
				Priority: 1,
				Target:   ".",
			}
			dnsAnswer.Answer = append(dnsAnswer.Answer, resourceRecord)
		}
	}

	for _, resourceRecord := range dnsAnswer.Answer {
		https, ok := resourceRecord.(*dns.HTTPS)
		if !ok {
			continue
		}

		hasValidECHConfig := slices.ContainsFunc(https.Value, func(value dns.SVCBKeyValue) bool {
			svcbECHConfig, ok := value.(*dns.SVCBECHConfig)
			return ok && len(svcbECHConfig.ECH) != 0
		})
		if !hasValidECHConfig {
			err := setECHConfig(context, https)
			if err != nil {
				return nil, err
			}
		}
	}

	content, err = resolvers.PackMessage(dnsAnswer)
	if err != nil {
		return nil, err
	}

	return content, nil
}
