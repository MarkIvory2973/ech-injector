package injectors

import (
	"context"
	"ech-injector/pkg/resolvers"
	"fmt"

	"github.com/miekg/dns"
)

func getECHConfig(context context.Context, name string) ([]byte, error) {
	svcbs, err := resolvers.ResolveHTTPS(context, name)
	if err != nil {
		return nil, err
	}

	for _, svcb := range svcbs {
		for _, value := range svcb {
			svcbECHConfig, ok := value.(*dns.SVCBECHConfig)
			if !ok {
				continue
			}

			return svcbECHConfig.ECH, nil
		}
	}

	return nil, fmt.Errorf("ECH config not found")
}
