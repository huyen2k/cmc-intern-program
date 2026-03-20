package scanner_test

import (
	"testing"

	"assets-api/internal/scanner"
)

func TestDNSScanner_Scan(t *testing.T) {
	s := scanner.NewDNSScanner()
	results, err := s.Scan("localhost")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}

func TestWhoisScanner_Scan(t *testing.T) {
	s := scanner.NewWhoisScanner()
	results, err := s.Scan("example.com")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}

func TestSubdomainScanner_Scan(t *testing.T) {
	s := scanner.NewSubdomainScanner()
	results, err := s.Scan("localhost")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}

func TestCertTransScanner_Scan(t *testing.T) {
	s := scanner.NewCertTransScanner()
	results, err := s.Scan("example.com")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}

func TestASNScanner_Scan(t *testing.T) {
	s := scanner.NewASNScanner()
	results, err := s.Scan("127.0.0.1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}
