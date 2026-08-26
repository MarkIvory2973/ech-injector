package cidrs

import (
	"context"
	"ech-injector/pkg/resolvers"
	"net/netip"
	"slices"
)

func contains(context context.Context, prefixes []netip.Prefix, name string) (bool, error) {
	ips, err := resolvers.ResolveA(context, name)
	if err != nil || len(ips) == 0 {
		ips, err = resolvers.ResolveAAAA(context, name)
		if err != nil {
			return false, err
		}

		if len(ips) == 0 {
			return false, nil
		}
	}

	address, err := netip.ParseAddr(ips[0])
	if err != nil {
		return false, err
	}

	yes := slices.ContainsFunc(prefixes, func(prefix netip.Prefix) bool {
		return prefix.Contains(address)
	})

	return yes, nil
}
