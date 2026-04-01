package scanner

import (
	"strings"
	"time"
)

type WhoisScanner struct{}

func NewWhoisScanner() *WhoisScanner {
	return &WhoisScanner{}
}

func (s *WhoisScanner) Scan(target string) ([]any, error) {
	tld := ""
	parts := strings.Split(strings.TrimSpace(target), ".")
	if len(parts) > 1 {
		tld = parts[len(parts)-1]
	}

	result := map[string]any{
		"domain": target,
		"whois": map[string]any{
			"registrar":       "Unknown",
			"whois_server":    "Unknown",
			"status":          "available",
			"name_servers":    []string{},
			"registrant":      "redacted",
			"tld":             strings.ToLower(tld),
			"last_checked_at": time.Now().UTC(),
		},
		"created_at": time.Now().UTC(),
	}

	return []any{result}, nil
}
