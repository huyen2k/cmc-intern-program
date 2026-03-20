package scanner

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type SSLScanner struct{}

func NewSSLScanner() *SSLScanner {
	return &SSLScanner{}
}

func (s *SSLScanner) Scan(target string) ([]any, error) {
	if ip := net.ParseIP(target); ip != nil {
		if err := EnsureSafeActiveTarget(target); err != nil {
			return nil, err
		}
	}

	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:443", target), &tls.Config{MinVersion: tls.VersionTLS12})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return nil, fmt.Errorf("no peer certificate")
	}

	cert := state.PeerCertificates[0]
	now := time.Now()
	days := int(time.Until(cert.NotAfter).Hours() / 24)

	san := make([]string, 0, len(cert.DNSNames))
	san = append(san, cert.DNSNames...)

	issues := []string{}
	if cert.NotAfter.Before(now) {
		issues = append(issues, "certificate expired")
	}
	if cert.Subject.String() == cert.Issuer.String() {
		issues = append(issues, "self signed certificate")
	}

	grade := "A"
	if len(issues) > 0 {
		grade = "C"
	}

	result := map[string]any{
		"domain": target,
		"certificate": map[string]any{
			"subject":           cert.Subject.String(),
			"issuer":            cert.Issuer.String(),
			"serial_number":     cert.SerialNumber.String(),
			"valid_from":        cert.NotBefore.UTC(),
			"valid_until":       cert.NotAfter.UTC(),
			"days_until_expiry": days,
			"is_expired":        cert.NotAfter.Before(now),
			"is_self_signed":    cert.Subject.String() == cert.Issuer.String(),
			"san":               san,
		},
		"connection": map[string]any{
			"tls_version":  tlsVersionString(state.Version),
			"cipher_suite": tls.CipherSuiteName(state.CipherSuite),
			"key_exchange": "",
		},
		"grade":      grade,
		"issues":     issues,
		"created_at": time.Now().UTC(),
	}

	return []any{result}, nil
}

func tlsVersionString(v uint16) string {
	switch v {
	case tls.VersionTLS13:
		return "TLS 1.3"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS10:
		return "TLS 1.0"
	default:
		return "unknown"
	}
}
