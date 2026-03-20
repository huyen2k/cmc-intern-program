package scanner

import "time"

type CertTransScanner struct{}

func NewCertTransScanner() *CertTransScanner {
	return &CertTransScanner{}
}

func (s *CertTransScanner) Scan(target string) ([]any, error) {
	result := map[string]any{
		"domain": target,
		"entries": []map[string]any{
			{
				"source":      "ct_logs",
				"common_name": target,
				"issuer":      "unknown",
				"not_before":  nil,
				"not_after":   nil,
			},
		},
		"created_at": time.Now().UTC(),
	}

	return []any{result}, nil
}
