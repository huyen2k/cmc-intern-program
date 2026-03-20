package service

import (
	"assets-api/internal/domain"
	"assets-api/internal/repository"
	"assets-api/internal/scanner"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type ScanService struct {
	assets   *AssetService
	scanRepo *repository.ScanRepository
	scanners map[string]scanner.Scanner
}

func NewScanService(assetService *AssetService, scanRepo *repository.ScanRepository) *ScanService {
	return &ScanService{
		assets:   assetService,
		scanRepo: scanRepo,
		scanners: map[string]scanner.Scanner{
			domain.ScanTypeDNS:       scanner.NewDNSScanner(),
			domain.ScanTypeWhois:     scanner.NewWhoisScanner(),
			domain.ScanTypeSubdomain: scanner.NewSubdomainScanner(),
			domain.ScanTypeCertTrans: scanner.NewCertTransScanner(),
			domain.ScanTypeASN:       scanner.NewASNScanner(),
			domain.ScanTypeIP:   scanner.NewIPScanner(),
			domain.ScanTypePort: scanner.NewPortScanner(),
			domain.ScanTypeSSL:  scanner.NewSSLScanner(),
			domain.ScanTypeTech: scanner.NewTechScanner(),
		},
	}
}

func (s *ScanService) StartScan(assetID, scanType string) (*domain.ScanJob, error) {
	if !domain.IsValidScanType(scanType) {
		return nil, errors.New("invalid scan_type")
	}

	asset, err := s.assets.GetByID(assetID)
	if err != nil {
		return nil, err
	}

	if err := validateScanTypeForAsset(asset, scanType); err != nil {
		return nil, err
	}

	job, err := s.scanRepo.CreateJob(assetID, scanType)
	if err != nil {
		return nil, err
	}

	go s.runScan(job, asset)
	return job, nil
}

func (s *ScanService) runScan(job *domain.ScanJob, asset *domain.Asset) {
	log.Printf("scan job %s started: asset=%s type=%s", job.ID, job.AssetID, job.ScanType)

	now := time.Now().UTC()
	_ = s.scanRepo.SetJobStatus(job.ID, domain.ScanStatusRunning, "", 0, nil)

	results := []any{}
	var err error

	target := strings.TrimSpace(asset.Name)

	if job.ScanType == domain.ScanTypeAll {
		results, err = s.runAllPassiveScans(target)
	} else {
		sc, ok := s.scanners[job.ScanType]
		if !ok {
			err = fmt.Errorf("scan type %s is not implemented", job.ScanType)
		} else {
			results, err = sc.Scan(target)
		}
	}

	ended := time.Now().UTC()
	if err != nil {
		_ = s.scanRepo.SetJobStatus(job.ID, domain.ScanStatusFailed, err.Error(), 0, &ended)
		log.Printf("scan job %s failed: %v", job.ID, err)
		return
	}

	if err := s.scanRepo.SaveResults(job.ID, job.AssetID, job.ScanType, results); err != nil {
		_ = s.scanRepo.SetJobStatus(job.ID, domain.ScanStatusFailed, err.Error(), 0, &ended)
		log.Printf("scan job %s failed while saving results: %v", job.ID, err)
		return
	}

	status := domain.ScanStatusCompleted
	errText := ""
	if job.ScanType == domain.ScanTypeAll && err != nil {
		status = domain.ScanStatusPartial
		errText = err.Error()
	}

	_ = now
	_ = s.scanRepo.SetJobStatus(job.ID, status, errText, len(results), &ended)
	log.Printf("scan job %s completed with %d results", job.ID, len(results))
}

func (s *ScanService) GetJob(jobID string) (*domain.ScanJob, error) {
	return s.scanRepo.GetJobByID(jobID)
}

func (s *ScanService) GetJobResults(jobID string) (*domain.ScanResultsResponse, error) {
	job, err := s.scanRepo.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}

	results, err := s.scanRepo.GetResultsByJobID(jobID)
	if err != nil {
		return nil, err
	}

	return &domain.ScanResultsResponse{
		JobID:    job.ID,
		ScanType: job.ScanType,
		Results:  results,
	}, nil
}

func (s *ScanService) ListAssetScans(assetID string) ([]domain.ScanJob, error) {
	return s.scanRepo.ListJobsByAssetID(assetID)
}

func (s *ScanService) ListAssetResults(assetID string) ([]domain.ScanResult, error) {
	return s.scanRepo.ListResultsByAssetID(assetID)
}

func (s *ScanService) GetLatestAssetResultByType(assetID, scanType string) (*domain.ScanResult, error) {
	results, err := s.scanRepo.ListResultsByAssetID(assetID)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if result.ScanType == scanType {
			copied := result
			return &copied, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (s *ScanService) runAllPassiveScans(target string) ([]any, error) {
	passiveTypes := []string{
		domain.ScanTypeDNS,
		domain.ScanTypeWhois,
		domain.ScanTypeSubdomain,
		domain.ScanTypeCertTrans,
	}

	results := make([]any, 0, len(passiveTypes))
	errorsFound := []string{}

	for _, scanType := range passiveTypes {
		sc, ok := s.scanners[scanType]
		if !ok {
			errorsFound = append(errorsFound, fmt.Sprintf("%s: not implemented", scanType))
			continue
		}

		items, err := sc.Scan(target)
		if err != nil {
			errorsFound = append(errorsFound, fmt.Sprintf("%s: %v", scanType, err))
			continue
		}

		results = append(results, map[string]any{
			"scan_type": scanType,
			"results":   items,
		})
	}

	if len(results) == 0 && len(errorsFound) > 0 {
		return nil, errors.New(strings.Join(errorsFound, "; "))
	}

	if len(errorsFound) > 0 {
		return results, errors.New(strings.Join(errorsFound, "; "))
	}

	return results, nil
}

func validateScanTypeForAsset(asset *domain.Asset, scanType string) error {
	isDomain := asset.Type == domain.AssetTypeDomain
	isIP := asset.Type == domain.AssetTypeIP

	switch scanType {
	case domain.ScanTypeIP, domain.ScanTypePort:
		if !isIP {
			return errors.New("scan_type requires ip asset")
		}
	case domain.ScanTypeSSL, domain.ScanTypeTech:
		if !isDomain && !isIP {
			return errors.New("scan_type requires domain or ip asset")
		}
		if isIP {
			if ip := net.ParseIP(asset.Name); ip != nil {
				if !ip.IsPrivate() && !ip.IsLoopback() {
					return errors.New("active scan for public IP is not allowed")
				}
			}
		}
	default:
		if !isDomain && scanType != domain.ScanTypeASN {
			return errors.New("scan_type requires domain asset")
		}
	}

	return nil
}
