package scanner_test

import (
	"testing"

	"assets-api/internal/scanner"
)

func TestIPScanner_Scan(t *testing.T) {
	s := scanner.NewIPScanner()

	results, err := s.Scan("127.0.0.1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}

func TestIPScanner_Scan_InvalidIP(t *testing.T) {
	s := scanner.NewIPScanner()

	_, err := s.Scan("not-an-ip")
	if err == nil {
		t.Fatalf("expected error for invalid ip")
	}
}
