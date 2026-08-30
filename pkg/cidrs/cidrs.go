package cidrs

import (
	"context"
	"ech-injector/pkg/rfc8484"
	"net/netip"
)

func contains(context context.Context, prefixes []netip.Prefix, name string) (bool, error) {
	ips, err := rfc8484.ExchangeA(context, name)
	if err != nil || len(ips) == 0 {
		ips, err = rfc8484.ExchangeAAAA(context, name)
		if err != nil {
			return false, err
		} else if len(ips) == 0 {
			return false, nil
		}
	}

	for _, ip := range ips {
		address, err := netip.ParseAddr(ip)
		if err != nil {
			return false, err
		}

		for _, prefix := range prefixes {
			if !prefix.Contains(address) {
				continue
			}

			return true, nil
		}
	}

	return false, nil
}
