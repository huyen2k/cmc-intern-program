package service_test

import (
	"errors"
	"testing"
	"time"

	"assets-api/internal/domain"
	"assets-api/internal/service"
)

type mockAssetRepository struct {
	createFn      func(asset domain.Asset) (*domain.Asset, error)
	listFn        func(page, limit int, t, status string) ([]domain.Asset, int, error)
	batchCreateFn func(assets []domain.Asset) ([]string, error)
}

func (m *mockAssetRepository) Create(asset domain.Asset) (*domain.Asset, error) {
	if m.createFn != nil {
		return m.createFn(asset)
	}
	return &asset, nil
}

func (m *mockAssetRepository) GetByID(id string) (*domain.Asset, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAssetRepository) DeleteByID(id string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *mockAssetRepository) GetStats() (*domain.Stats, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAssetRepository) Count(t, status string) (int, error) {
	return 0, errors.New("not implemented")
}

func (m *mockAssetRepository) BatchCreate(assets []domain.Asset) ([]string, error) {
	if m.batchCreateFn != nil {
		return m.batchCreateFn(assets)
	}
	return []string{}, nil
}

func (m *mockAssetRepository) BatchDelete(ids []string) (int, int, error) {
	return 0, 0, errors.New("not implemented")
}

func (m *mockAssetRepository) List(page, limit int, t, status string) ([]domain.Asset, int, error) {
	if m.listFn != nil {
		return m.listFn(page, limit, t, status)
	}
	return []domain.Asset{}, 0, nil
}

func (m *mockAssetRepository) Search(q string) ([]domain.Asset, error) {
	return []domain.Asset{}, nil
}

func TestAssetService_Create(t *testing.T) {
	tests := []struct {
		name    string
		asset   domain.Asset
		wantErr bool
	}{
		{
			name:    "valid asset with default status",
			asset:   domain.Asset{Name: "example.com", Type: domain.AssetTypeDomain},
			wantErr: false,
		},
		{
			name:    "invalid type",
			asset:   domain.Asset{Name: "example.com", Type: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAssetRepository{}
			svc := service.NewAssetService(repo)

			created, err := svc.Create(tt.asset)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if created.Status != "active" {
				t.Fatalf("expected default status active, got: %s", created.Status)
			}
		})
	}
}

func TestAssetService_BatchCreate(t *testing.T) {
	t.Run("reject over limit", func(t *testing.T) {
		repo := &mockAssetRepository{}
		svc := service.NewAssetService(repo)

		assets := make([]domain.Asset, 101)
		for i := range assets {
			assets[i] = domain.Asset{Name: "example.com", Type: domain.AssetTypeDomain}
		}

		_, err := svc.BatchCreate(assets)
		if err == nil {
			t.Fatalf("expected error for >100 assets")
		}
	})

	t.Run("validate each asset", func(t *testing.T) {
		repo := &mockAssetRepository{}
		svc := service.NewAssetService(repo)

		assets := []domain.Asset{
			{Name: "ok.com", Type: domain.AssetTypeDomain},
			{Name: "bad", Type: "invalid"},
		}

		_, err := svc.BatchCreate(assets)
		if err == nil {
			t.Fatalf("expected validation error")
		}
	})

	t.Run("success delegates repository", func(t *testing.T) {
		repo := &mockAssetRepository{
			batchCreateFn: func(assets []domain.Asset) ([]string, error) {
				return []string{"id-1", "id-2"}, nil
			},
		}
		svc := service.NewAssetService(repo)

		ids, err := svc.BatchCreate([]domain.Asset{
			{Name: "a.com", Type: domain.AssetTypeDomain},
			{Name: "127.0.0.1", Type: domain.AssetTypeIP},
		})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(ids) != 2 {
			t.Fatalf("expected 2 ids, got: %d", len(ids))
		}
	})
}

func TestAssetService_List(t *testing.T) {
	now := time.Now().UTC()
	repo := &mockAssetRepository{
		listFn: func(page, limit int, t, status string) ([]domain.Asset, int, error) {
			if page != 2 || limit != 20 || t != domain.AssetTypeDomain || status != "active" {
				return nil, 0, errors.New("unexpected arguments")
			}
			return []domain.Asset{{ID: "1", Name: "a.com", Type: domain.AssetTypeDomain, Status: "active", CreatedAt: now}}, 1, nil
		},
	}

	svc := service.NewAssetService(repo)
	items, total, err := svc.List(2, 20, domain.AssetTypeDomain, "active")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("unexpected list result total=%d len=%d", total, len(items))
	}
}
