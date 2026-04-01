package scanner

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type TechScanner struct {
	client *http.Client
}

func NewTechScanner() *TechScanner {
	return &TechScanner{client: &http.Client{Timeout: 8 * time.Second}}
}

func (s *TechScanner) Scan(target string) ([]any, error) {
	url := target
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + target
	}

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*512))
	if err != nil {
		return nil, err
	}

	text := strings.ToLower(string(body))
	techs := detectTechnologies(resp, text)
	if len(techs) == 0 {
		return nil, errors.New("no technology signals detected")
	}

	headers := map[string]string{}
	for k, values := range resp.Header {
		if len(values) > 0 {
			headers[strings.ToLower(k)] = values[0]
		}
	}

	meta := map[string]string{}
	metaRe := regexp.MustCompile(`(?i)<meta\s+name=["']([^"']+)["']\s+content=["']([^"']+)["']`)
	for _, m := range metaRe.FindAllStringSubmatch(string(body), -1) {
		meta[strings.ToLower(m[1])] = m[2]
	}

	result := map[string]any{
		"domain":       target,
		"technologies": techs,
		"headers":      headers,
		"meta_tags":    meta,
		"created_at":   time.Now().UTC(),
	}

	return []any{result}, nil
}

func detectTechnologies(resp *http.Response, body string) []map[string]any {
	out := []map[string]any{}
	seen := map[string]bool{}

	add := func(name, category, version string, confidence int) {
		if seen[name] {
			return
		}
		seen[name] = true
		item := map[string]any{
			"name":       name,
			"category":   category,
			"version":    nil,
			"confidence": confidence,
		}
		if version != "" {
			item["version"] = version
		}
		out = append(out, item)
	}

	server := strings.ToLower(resp.Header.Get("Server"))
	xpb := strings.ToLower(resp.Header.Get("X-Powered-By"))
	if strings.Contains(server, "nginx") {
		add("nginx", "Web Server", "", 100)
	}
	if strings.Contains(server, "apache") {
		add("Apache", "Web Server", "", 100)
	}
	if strings.Contains(server, "cloudflare") || resp.Header.Get("CF-Ray") != "" {
		add("Cloudflare", "CDN", "", 100)
	}
	if strings.Contains(xpb, "express") {
		add("Express", "Backend Framework", "", 95)
	}
	if strings.Contains(xpb, "php") {
		add("PHP", "Language Runtime", "", 95)
	}

	if strings.Contains(body, "react") || strings.Contains(body, "_next") {
		add("React", "JavaScript Framework", "", 85)
	}
	if strings.Contains(body, "vue") || strings.Contains(body, "nuxt") {
		add("Vue.js", "JavaScript Framework", "", 85)
	}
	if strings.Contains(body, "wp-content") {
		add("WordPress", "CMS", "", 90)
	}

	return out
}
