package util

import (
	"context"
	"net"
	"regexp"
	"time"
)

func ResolveDomain(dnsServer, domain string) []string {
	ctx := context.TODO()
	var resolver = net.DefaultResolver
	if dnsServer != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Second,
				}
				return d.DialContext(ctx, "UDP", dnsServer)
			},
		}
	}
	hosts, err := resolver.LookupHost(ctx, domain)
	if err != nil {
		return []string{}
	}
	return hosts
}

func IsIpv4(ip string) bool {
	var rule = regexp.MustCompile(`^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$`)
	return rule.MatchString(ip)
}
