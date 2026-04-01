package domain

import "time"

const (
	ScanTypeDNS       = "dns"
	ScanTypeWhois     = "whois"
	ScanTypeSubdomain = "subdomain"
	ScanTypeCertTrans = "cert_trans"
	ScanTypeASN       = "asn"
	ScanTypeAll       = "all"
	ScanTypeIP        = "ip"
	ScanTypePort      = "port"
	ScanTypeSSL       = "ssl"
	ScanTypeTech      = "tech"
)

const (
	ScanStatusPending   = "pending"
	ScanStatusRunning   = "running"
	ScanStatusCompleted = "completed"
	ScanStatusFailed    = "failed"
	ScanStatusPartial   = "partial"
)

type ScanJob struct {
	ID        string     `json:"id"`
	AssetID   string     `json:"asset_id"`
	ScanType  string     `json:"scan_type"`
	Status    string     `json:"status"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	Error     string     `json:"error"`
	Results   int        `json:"results"`
	CreatedAt time.Time  `json:"created_at"`
}

type ScanResult struct {
	ID        string    `json:"id"`
	JobID     string    `json:"job_id"`
	AssetID   string    `json:"asset_id"`
	ScanType  string    `json:"scan_type"`
	Data      any       `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type StartScanRequest struct {
	ScanType string `json:"scan_type"`
}

type ScanResultsResponse struct {
	JobID    string `json:"job_id"`
	ScanType string `json:"scan_type"`
	Results  []any  `json:"results"`
}

func IsValidScanType(scanType string) bool {
	switch scanType {
	case ScanTypeDNS, ScanTypeWhois, ScanTypeSubdomain, ScanTypeCertTrans, ScanTypeASN, ScanTypeAll,
		ScanTypeIP, ScanTypePort, ScanTypeSSL, ScanTypeTech:
		return true
	default:
		return false
	}
}
