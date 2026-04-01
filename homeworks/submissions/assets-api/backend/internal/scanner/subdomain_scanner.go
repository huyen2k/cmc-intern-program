package scanner

import (
	"fmt"
	"net"
	"time"
)

type SubdomainScanner struct{}

func NewSubdomainScanner() *SubdomainScanner {
	return &SubdomainScanner{}
}

func (s *SubdomainScanner) Scan(target string) ([]any, error) {
	common := []string{"www", "api", "dev", "staging", "mail"}
	found := []map[string]any{}

	for _, prefix := range common {
		sub := fmt.Sprintf("%s.%s", prefix, target)
		ips, err := net.LookupIP(sub)
		if err != nil || len(ips) == 0 {
			continue
		}

		resolved := make([]string, 0, len(ips))
		for _, ip := range ips {
			resolved = append(resolved, ip.String())
		}

		found = append(found, map[string]any{
			"name":       sub,
			"ip_addresses": resolved,
		})
	}

	result := map[string]any{
		"domain":      target,
		"subdomains":  found,
		"total_found": len(found),
		"created_at":  time.Now().UTC(),
	}

	return []any{result}, nil
}
