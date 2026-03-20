package scanner

import (
	"net"
	"time"
)

type DNSScanner struct{}

func NewDNSScanner() *DNSScanner {
	return &DNSScanner{}
}

func (s *DNSScanner) Scan(target string) ([]any, error) {
	ips, _ := net.LookupIP(target)
	mxRecords, _ := net.LookupMX(target)
	nsRecords, _ := net.LookupNS(target)

	aRecords := make([]string, 0, len(ips))
	for _, ip := range ips {
		aRecords = append(aRecords, ip.String())
	}

	mx := make([]map[string]any, 0, len(mxRecords))
	for _, record := range mxRecords {
		mx = append(mx, map[string]any{
			"host": record.Host,
			"pref": record.Pref,
		})
	}

	ns := make([]string, 0, len(nsRecords))
	for _, record := range nsRecords {
		ns = append(ns, record.Host)
	}

	result := map[string]any{
		"domain": target,
		"records": map[string]any{
			"a":  aRecords,
			"mx": mx,
			"ns": ns,
		},
		"created_at": time.Now().UTC(),
	}

	return []any{result}, nil
}
