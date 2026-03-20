package scanner

import (
	"errors"
	"net"
	"time"
)

type IPScanner struct{}

func NewIPScanner() *IPScanner {
	return &IPScanner{}
}

func (s *IPScanner) Scan(target string) ([]any, error) {
	ip := net.ParseIP(target)
	if ip == nil {
		return nil, errors.New("ip scan requires a valid IP asset")
	}

	geo := map[string]any{
		"country":      "Unknown",
		"country_code": "UN",
		"city":         "Unknown",
		"region":       "Unknown",
		"latitude":     0.0,
		"longitude":    0.0,
		"isp":          "Unknown",
		"org":          "Unknown",
	}

	asn := map[string]any{
		"number":      0,
		"name":        "UNKNOWN",
		"description": "Unknown ASN",
	}

	if ip.IsLoopback() {
		geo["country"] = "Local"
		geo["country_code"] = "LO"
		geo["city"] = "localhost"
		asn["name"] = "LOOPBACK"
		asn["description"] = "Loopback network"
	}

	reverse := ""
	names, err := net.LookupAddr(target)
	if err == nil && len(names) > 0 {
		reverse = names[0]
	}

	result := map[string]any{
		"ip_address":  target,
		"geolocation": geo,
		"asn":         asn,
		"reverse_dns": reverse,
		"created_at":  time.Now().UTC(),
	}

	return []any{result}, nil
}
