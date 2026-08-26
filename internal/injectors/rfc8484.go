package injectors

import (
	"context"
	"ech-injector/pkg/cidrs"
	"ech-injector/pkg/resolvers"
	"encoding/base64"
	"fmt"

	"github.com/miekg/dns"
)

func getECHConfig(context context.Context, name string) ([]byte, error) {
	svcbs, err := resolvers.ResolveHTTPS(context, name)
	if err != nil {
		return nil, err
	}

	for _, value := range svcbs[0] {
		svcbECHConfig, ok := value.(*dns.SVCBECHConfig)
		if !ok {
			continue
		}

		echConfig := svcbECHConfig.ECH

		return echConfig, err
	}

	return nil, fmt.Errorf("ECH not found")
}

func InjectRFC8484(context context.Context, content []byte) ([]byte, error) {
	dnsQuestion, err := resolvers.UnpackMessage(content)
	if err != nil {
		return nil, err
	}

	dnsAnwser, err := resolvers.ExchangeMessage(context, dnsQuestion)
	if err != nil {
		return nil, err
	}

	exists := false
	for _, resourceRecord := range dnsAnwser.Answer {
		https, ok := resourceRecord.(*dns.HTTPS)
		if !ok {
			continue
		}

		exists = true

		if !exists {
			yes, err := cidrs.IsCloudflare(context, https.Hdr.Name)
			if err != nil {
				continue
			}

			if yes {
				echConfig, err := getECHConfig(context, "cloudflare-ech.com.")
				if err != nil {
					return nil, err
				}

				resourceRecord := &dns.SVCB{
					Hdr: dns.RR_Header{
						Name:   https.Hdr.Name,
						Rrtype: dns.TypeHTTPS,
						Class:  dns.ClassINET,
						Ttl:    300,
					},
					Priority: 1,
					Target:   ".",
					Value: []dns.SVCBKeyValue{
						&dns.SVCBECHConfig{
							ECH: echConfig,
						},
					},
				}
				dnsAnwser.Answer = append(dnsAnwser.Answer, resourceRecord)
			}

			yes, err = cidrs.IsMeta(context, https.Hdr.Name)
			if err != nil {
				continue
			}

			if yes {
				echConfig, err := base64.StdEncoding.DecodeString("AEj+DQBEAQAgACAdd+scUi0IYFsXnUIU7ko2Nd9+F8M26pAGZVpz/KrWPgAEAAEAAWQVZWNoLXB1YmxpYy5hdG1ldGEuY29tAAA=")
				if err != nil {
					return nil, err
				}

				resourceRecord := &dns.SVCB{
					Hdr: dns.RR_Header{
						Name:   https.Hdr.Name,
						Rrtype: dns.TypeHTTPS,
						Class:  dns.ClassINET,
						Ttl:    300,
					},
					Priority: 1,
					Target:   ".",
					Value: []dns.SVCBKeyValue{
						&dns.SVCBECHConfig{
							ECH: echConfig,
						},
					},
				}
				dnsAnwser.Answer = append(dnsAnwser.Answer, resourceRecord)
			}
		}

		hasValidECHConfig := false
		for _, value := range https.Value {
			svcbECHConfig, ok := value.(*dns.SVCBECHConfig)
			if !ok {
				continue
			}

			if len(svcbECHConfig.ECH) == 0 {
				continue
			}

			hasValidECHConfig = true
		}

		if !hasValidECHConfig {
			yes, err := cidrs.IsCloudflare(context, https.Hdr.Name)
			if err != nil {
				continue
			}

			if yes {
				echConfig, err := getECHConfig(context, "cloudflare-ech.com.")
				if err != nil {
					return nil, err
				}

				value := &dns.SVCBECHConfig{
					ECH: echConfig,
				}
				https.Value = append(https.Value, value)
			}

			yes, err = cidrs.IsMeta(context, https.Hdr.Name)
			if err != nil {
				continue
			}

			if yes {
				echConfig, err := base64.StdEncoding.DecodeString("AEj+DQBEAQAgACAdd+scUi0IYFsXnUIU7ko2Nd9+F8M26pAGZVpz/KrWPgAEAAEAAWQVZWNoLXB1YmxpYy5hdG1ldGEuY29tAAA=")
				if err != nil {
					return nil, err
				}

				value := &dns.SVCBECHConfig{
					ECH: echConfig,
				}
				https.Value = append(https.Value, value)
			}
		}
	}

	content, err = resolvers.PackMessage(dnsAnwser)
	if err != nil {
		return nil, err
	}

	return content, nil
}
