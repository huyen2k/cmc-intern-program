package service

import (
	"database/sql"
	"errors"

	"assets-api/internal/domain"
)

type AssetRepository interface {
	Create(asset domain.Asset) (*domain.Asset, error)
	GetByID(id string) (*domain.Asset, error)
	DeleteByID(id string) (bool, error)
	GetStats() (*domain.Stats, error)
	Count(t, status string) (int, error)
	BatchCreate(assets []domain.Asset) ([]string, error)
	BatchDelete(ids []string) (int, int, error)
	List(page, limit int, t, status string) ([]domain.Asset, int, error)
	Search(q string) ([]domain.Asset, error)
}

type AssetService struct {
	repo AssetRepository
}

func NewAssetService(r AssetRepository) *AssetService {
	return &AssetService{r}
}

func (s *AssetService) Create(asset domain.Asset) (*domain.Asset, error) {
	if asset.Status == "" {
		asset.Status = "active"
	}

	if err := asset.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Create(asset)
}

func (s *AssetService) GetByID(id string) (*domain.Asset, error) {
	return s.repo.GetByID(id)
}

func (s *AssetService) DeleteByID(id string) (bool, error) {
	return s.repo.DeleteByID(id)
}

func (s *AssetService) GetStats() (*domain.Stats, error) {
	return s.repo.GetStats()
}

func (s *AssetService) Count(t, status string) (int, error) {
	return s.repo.Count(t, status)
}

func (s *AssetService) BatchCreate(assets []domain.Asset) ([]string, error) {

	if len(assets) > 100 {
		return nil, errors.New("max 100 assets per request")
	}

	for _, a := range assets {
		if a.Status == "" {
			a.Status = "active"
		}

		if err := a.Validate(); err != nil {
			return nil, err
		}
	}

	return s.repo.BatchCreate(assets)
}

func (s *AssetService) BatchDelete(ids []string) (int, int, error) {
	return s.repo.BatchDelete(ids)
}

func (s *AssetService) List(page, limit int, t, status string) ([]domain.Asset, int, error) {
	return s.repo.List(page, limit, t, status)
}

func (s *AssetService) Search(q string) ([]domain.Asset, error) {
	return s.repo.Search(q)
}

func IsNotFoundErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
