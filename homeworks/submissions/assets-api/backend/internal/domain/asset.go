package domain

import (
	"errors"
	"net"
	"strings"
	"time"
)

type Asset struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Stats struct {
	Total    int            `json:"total"`
	ByType   map[string]int `json:"by_type"`
	ByStatus map[string]int `json:"by_status"`
}

func (a Asset) Validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("asset name is required")
	}

	switch a.Type {
	case AssetTypeDomain:
		if strings.Contains(a.Name, " ") {
			return errors.New("invalid domain")
		}
	case AssetTypeIP:
		if net.ParseIP(a.Name) == nil {
			return errors.New("invalid ip address")
		}
	case AssetTypeService:
		// service type remains permissive for backward compatibility.
	default:
		return errors.New("invalid asset type")
	}

	return nil
}

const (
	AssetTypeDomain  = "domain"
	AssetTypeIP      = "ip"
	AssetTypeService = "service"
)
