package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"assets-api/internal/domain"
	"assets-api/internal/handler"
	"assets-api/internal/service"
)

type mockAssetRepository struct {
	createFn func(asset domain.Asset) (*domain.Asset, error)
	listFn   func(page, limit int, t, status string) ([]domain.Asset, int, error)
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
	return nil, errors.New("not implemented")
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

func TestCreateAssetHandler(t *testing.T) {
	repo := &mockAssetRepository{
		createFn: func(asset domain.Asset) (*domain.Asset, error) {
			asset.ID = "id-1"
			asset.CreatedAt = time.Now().UTC()
			if asset.Status == "" {
				asset.Status = "active"
			}
			return &asset, nil
		},
	}
	h := handler.NewAssetHandler(service.NewAssetService(repo), nil)

	body := `{"name":"test.com","type":"domain"}`
	req := httptest.NewRequest(http.MethodPost, "/assets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var got domain.Asset
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.ID == "" || got.Name != "test.com" || got.Type != "domain" {
		t.Fatalf("unexpected created asset: %+v", got)
	}
}

func TestListAssetHandler_DefaultAndCapLimit(t *testing.T) {
	repo := &mockAssetRepository{
		listFn: func(page, limit int, t, status string) ([]domain.Asset, int, error) {
			if page != 1 {
				return nil, 0, errors.New("expected default page=1")
			}
			if limit != 100 {
				return nil, 0, errors.New("expected capped limit=100")
			}
			return []domain.Asset{}, 0, nil
		},
	}
	h := handler.NewAssetHandler(service.NewAssetService(repo), nil)

	req := httptest.NewRequest(http.MethodGet, "/assets?page=0&limit=1000", nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}
}
