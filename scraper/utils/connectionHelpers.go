package utils

import (
	"context"
	"net"
	"net/http"
	"time"
)

func GetTransportForRequest(dnsAddress string) *http.Transport {
	// set the default DNS for the http request. by default this is will be the cloudflare DNS Address 1.1.1.1:53.
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", dnsAddress) // using DNS address
		},
	}

	dialer := &net.Dialer{
		Timeout:  5 * time.Second,
		Resolver: resolver, // attach custom resolver
	}
	return &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}
}
