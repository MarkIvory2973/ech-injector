package injectors

import (
	"context"
	"ech-injector/pkg/cidrs"
	"ech-injector/pkg/resolvers"
	"ech-injector/pkg/workers"
	"encoding/base64"
	"fmt"

	"github.com/miekg/dns"
)

func getECHConfig(context context.Context, name string) ([]byte, error) {
	content, err := workers.CacheFunc(name, func() (map[string]string, int, error) {
		svcbs, err := resolvers.ExchangeHTTPS(context, name)
		if err != nil {
			return nil, 0, err
		}

		for _, svcb := range svcbs {
			for _, value := range svcb {
				svcbECHConfig, ok := value.(*dns.SVCBECHConfig)
				if !ok {
					continue
				}

				echConfig := base64.StdEncoding.EncodeToString(svcbECHConfig.ECH)
				content := map[string]string{
					"echConfig": echConfig,
				}

				return content, 0, err
			}
		}

		return nil, 3600, fmt.Errorf("ECH config not found")
	})
	if err != nil {
		return nil, err
	}

	echConfig, err := base64.StdEncoding.DecodeString(content["echConfig"])
	if err != nil {
		return nil, err
	}

	return echConfig, nil
}

func setECHConfig(context context.Context, https *dns.HTTPS) error {
	yes, err := cidrs.IsCloudflare(context, https.Hdr.Name)
	if err != nil {
		return err
	}

	if yes {
		echConfig, err := getECHConfig(context, "cloudflare-ech.com.")
		if err != nil {
			return err
		}

		value := &dns.SVCBECHConfig{
			ECH: echConfig,
		}
		https.Value = append(https.Value, value)
		return nil
	}

	yes, err = cidrs.IsMeta(context, https.Hdr.Name)
	if err != nil {
		return err
	}

	if yes {
		echConfig, err := base64.StdEncoding.DecodeString("AEj+DQBEAQAgACAdd+scUi0IYFsXnUIU7ko2Nd9+F8M26pAGZVpz/KrWPgAEAAEAAWQVZWNoLXB1YmxpYy5hdG1ldGEuY29tAAA=")
		if err != nil {
			return err
		}

		value := &dns.SVCBECHConfig{
			ECH: echConfig,
		}
		https.Value = append(https.Value, value)
		return nil
	}

	return nil
}
