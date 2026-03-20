package repository

import (
	"assets-api/internal/domain"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ScanRepository struct {
	db *sql.DB
}

func NewScanRepository(db *sql.DB) *ScanRepository {
	return &ScanRepository{db: db}
}

func (r *ScanRepository) EnsureSchema() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS scan_jobs (
			id UUID PRIMARY KEY,
			asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
			scan_type TEXT NOT NULL CHECK (scan_type IN ('dns','whois','subdomain','cert_trans','asn','all','ip','port','ssl','tech')),
			status TEXT NOT NULL,
			started_at TIMESTAMP NULL,
			ended_at TIMESTAMP NULL,
			error TEXT NOT NULL DEFAULT '',
			results INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_scan_jobs_asset_id ON scan_jobs(asset_id);
		CREATE INDEX IF NOT EXISTS idx_scan_jobs_status ON scan_jobs(status);
	`)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS scan_results (
			id UUID PRIMARY KEY,
			job_id UUID NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
			asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
			scan_type TEXT NOT NULL,
			data JSONB NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_scan_results_job_id ON scan_results(job_id);
		CREATE INDEX IF NOT EXISTS idx_scan_results_asset_id ON scan_results(asset_id);
	`)

	return err
}

func (r *ScanRepository) CreateJob(assetID, scanType string) (*domain.ScanJob, error) {
	now := time.Now().UTC()
	job := &domain.ScanJob{
		ID:        uuid.New().String(),
		AssetID:   assetID,
		ScanType:  scanType,
		Status:    domain.ScanStatusPending,
		StartedAt: &now,
		Error:     "",
		Results:   0,
		CreatedAt: now,
	}

	err := r.db.QueryRow(`
		INSERT INTO scan_jobs(id,asset_id,scan_type,status,started_at,error,results,created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING created_at
	`, job.ID, job.AssetID, job.ScanType, job.Status, job.StartedAt, job.Error, job.Results, job.CreatedAt).Scan(&job.CreatedAt)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (r *ScanRepository) SetJobStatus(jobID, status, errText string, results int, endedAt *time.Time) error {
	_, err := r.db.Exec(`
		UPDATE scan_jobs
		SET status=$2, error=$3, results=$4, ended_at=$5
		WHERE id=$1
	`, jobID, status, errText, results, endedAt)
	return err
}

func (r *ScanRepository) GetJobByID(jobID string) (*domain.ScanJob, error) {
	var j domain.ScanJob
	err := r.db.QueryRow(`
		SELECT id,asset_id,scan_type,status,started_at,ended_at,error,results,created_at
		FROM scan_jobs WHERE id=$1
	`, jobID).Scan(&j.ID, &j.AssetID, &j.ScanType, &j.Status, &j.StartedAt, &j.EndedAt, &j.Error, &j.Results, &j.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (r *ScanRepository) ListJobsByAssetID(assetID string) ([]domain.ScanJob, error) {
	rows, err := r.db.Query(`
		SELECT id,asset_id,scan_type,status,started_at,ended_at,error,results,created_at
		FROM scan_jobs WHERE asset_id=$1 ORDER BY created_at DESC
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.ScanJob{}
	for rows.Next() {
		var j domain.ScanJob
		if err := rows.Scan(&j.ID, &j.AssetID, &j.ScanType, &j.Status, &j.StartedAt, &j.EndedAt, &j.Error, &j.Results, &j.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, j)
	}

	return out, rows.Err()
}

func (r *ScanRepository) SaveResults(jobID, assetID, scanType string, rows []any) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range rows {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
			INSERT INTO scan_results(id,job_id,asset_id,scan_type,data)
			VALUES($1,$2,$3,$4,$5)
		`, uuid.New().String(), jobID, assetID, scanType, data)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ScanRepository) GetResultsByJobID(jobID string) ([]any, error) {
	rows, err := r.db.Query(`SELECT data FROM scan_results WHERE job_id=$1 ORDER BY created_at ASC`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []any{}
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}

		var item any
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, rows.Err()
}

func (r *ScanRepository) ListResultsByAssetID(assetID string) ([]domain.ScanResult, error) {
	rows, err := r.db.Query(`
		SELECT id,job_id,asset_id,scan_type,data,created_at
		FROM scan_results WHERE asset_id=$1 ORDER BY created_at DESC
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []domain.ScanResult{}
	for rows.Next() {
		var s domain.ScanResult
		var raw []byte
		if err := rows.Scan(&s.ID, &s.JobID, &s.AssetID, &s.ScanType, &raw, &s.CreatedAt); err != nil {
			return nil, err
		}

		var data any
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		s.Data = data
		results = append(results, s)
	}

	return results, rows.Err()
}
