package scanner

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

type PortScanner struct {
	ports []int
}

func NewPortScanner() *PortScanner {
	ports := make([]int, 0, 1000)
	for i := 1; i <= 1000; i++ {
		ports = append(ports, i)
	}

	return &PortScanner{ports: ports}
}

func (s *PortScanner) Scan(target string) ([]any, error) {
	if err := EnsureSafeActiveTarget(target); err != nil {
		return nil, err
	}

	start := time.Now()
	type openPort struct {
		port     int
		protocol string
		state    string
		service  string
		version  string
	}

	openPorts := make([]openPort, 0)
	closed := 0

	var mu sync.Mutex
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 100)

	for _, p := range s.ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			addr := fmt.Sprintf("%s:%d", target, port)
			conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
			if err != nil {
				mu.Lock()
				closed++
				mu.Unlock()
				return
			}
			_ = conn.Close()

			mu.Lock()
			openPorts = append(openPorts, openPort{
				port:     port,
				protocol: "tcp",
				state:    "open",
				service:  guessService(port),
				version:  "",
			})
			mu.Unlock()
		}(p)
	}

	wg.Wait()

	sort.Slice(openPorts, func(i, j int) bool { return openPorts[i].port < openPorts[j].port })

	results := make([]map[string]any, 0, len(openPorts))
	for _, p := range openPorts {
		results = append(results, map[string]any{
			"port":     p.port,
			"protocol": p.protocol,
			"state":    p.state,
			"service":  p.service,
			"version":  p.version,
		})
	}

	if len(results) == 0 && closed == 0 {
		return nil, errors.New("no ports were scanned")
	}

	payload := map[string]any{
		"ip_address":       strings.TrimSpace(target),
		"open_ports":       results,
		"closed_ports":     closed,
		"total_scanned":    len(s.ports),
		"scan_duration_ms": time.Since(start).Milliseconds(),
		"created_at":       time.Now().UTC(),
	}

	return []any{payload}, nil
}

func guessService(port int) string {
	switch port {
	case 22:
		return "ssh"
	case 53:
		return "dns"
	case 80:
		return "http"
	case 443:
		return "https"
	case 5432:
		return "postgresql"
	case 6379:
		return "redis"
	default:
		return "unknown"
	}
}
