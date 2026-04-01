package scanner

import (
	"net"
	"time"
)

type ASNScanner struct{}

func NewASNScanner() *ASNScanner {
	return &ASNScanner{}
}

func (s *ASNScanner) Scan(target string) ([]any, error) {
	reverse := ""
	names, err := net.LookupAddr(target)
	if err == nil && len(names) > 0 {
		reverse = names[0]
	}

	result := map[string]any{
		"ip_address": target,
		"asn": map[string]any{
			"number":      0,
			"name":        "UNKNOWN",
			"description": "ASN data unavailable",
		},
		"reverse_dns": reverse,
		"created_at":  time.Now().UTC(),
	}

	return []any{result}, nil
}
