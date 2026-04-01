package scanner

import (
	"errors"
	"net"
)

type Scanner interface {
	Scan(target string) ([]any, error)
}

func IsPrivateOrLocalIP(raw string) bool {
	if raw == "localhost" {
		return true
	}

	ip := net.ParseIP(raw)
	if ip == nil {
		return false
	}

	if ip.IsLoopback() || ip.IsPrivate() {
		return true
	}

	return false
}

func EnsureSafeActiveTarget(target string) error {
	if !IsPrivateOrLocalIP(target) {
		return errors.New("active scan is only allowed for localhost and private IP ranges")
	}

	return nil
}
