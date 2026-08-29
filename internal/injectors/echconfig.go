package injectors

import (
	"context"
	"ech-injector/pkg/resolvers"
	"ech-injector/pkg/workers"
	"encoding/base64"
	"fmt"

	"github.com/miekg/dns"
)

func getECHConfig(context context.Context, name string) ([]byte, error) {
	content, err := workers.CacheFunc(name, func() (map[string]string, int, error) {
		svcbs, err := resolvers.ResolveHTTPS(context, name)
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
