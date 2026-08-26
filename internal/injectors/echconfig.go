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
	content, err := workers.CacheFunc(name, func() (map[string]string, error) {
		svcbs, err := resolvers.ResolveHTTPS(context, name)
		if err != nil {
			return nil, err
		}

		for _, value := range svcbs[0] {
			svcbECHConfig, ok := value.(*dns.SVCBECHConfig)
			if !ok {
				continue
			}

			echConfig := base64.StdEncoding.EncodeToString(svcbECHConfig.ECH)
			content := map[string]string{
				"echConfig": echConfig,
			}

			return content, err
		}

		return nil, fmt.Errorf("ECH config not found")
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
