package cidrs

import (
	"context"
	"net/netip"
	"ech-injector/pkg/resolvers"
)

func contains(context context.Context, prefixes []netip.Prefix, name string) (bool, error) {
	ips, err := resolvers.ResolveA(context, name)
	if err != nil {
		ips, err = resolvers.ResolveAAAA(context, name)
		if err != nil {
			return false, err
		}
	}

	if len(ips[0]) == 0 {
		return false, err
	}

	address, err := netip.ParseAddr(ips[0])
	if err != nil {
		return false, err
	}

	for _, prefix := range prefixes {
		if !prefix.Contains(address) {
			continue
		}

		return true, nil
	}

	return false, nil
}
