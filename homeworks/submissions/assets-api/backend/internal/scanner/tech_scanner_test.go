package scanner_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"assets-api/internal/scanner"
)

func TestTechScanner_Scan(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "nginx")
		w.Header().Set("X-Powered-By", "Express")
		_, _ = w.Write([]byte("<html><head><meta name=\"generator\" content=\"React\"></head><body>react app</body></html>"))
	}))
	defer ts.Close()

	s := scanner.NewTechScanner()
	results, err := s.Scan(ts.URL)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}
