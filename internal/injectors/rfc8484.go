package injectors

import (
	"context"
	"slices"

	"github.com/miekg/dns"
)

func InjectRFC8484(context context.Context, dnsQuestion *dns.Msg, dnsAnswer *dns.Msg) error {
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
				return err
			}
		}
	}

	return nil
}
