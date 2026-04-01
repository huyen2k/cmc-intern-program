package scanner_test

import (
	"net"
	"testing"

	"assets-api/internal/scanner"
)

func TestPortScanner_Scan(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:80")
	if err != nil {
		t.Skipf("cannot bind port 80 in test environment: %v", err)
	}
	defer ln.Close()

	s := scanner.NewPortScanner()
	results, err := s.Scan("127.0.0.1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}

func TestPortScanner_Scan_BlockPublicIP(t *testing.T) {
	s := scanner.NewPortScanner()
	_, err := s.Scan("8.8.8.8")
	if err == nil {
		t.Fatalf("expected safety check error")
	}
}
