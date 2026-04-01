package domain_test

import (
	"testing"

	"assets-api/internal/domain"
)

func TestAssetValidation(t *testing.T) {
	tests := []struct {
		name    string
		asset   domain.Asset
		wantErr bool
	}{
		{
			name:    "valid domain asset",
			asset:   domain.Asset{Name: "example.com", Type: domain.AssetTypeDomain},
			wantErr: false,
		},
		{
			name:    "valid ip asset",
			asset:   domain.Asset{Name: "127.0.0.1", Type: domain.AssetTypeIP},
			wantErr: false,
		},
		{
			name:    "invalid asset type",
			asset:   domain.Asset{Name: "test", Type: "invalid"},
			wantErr: true,
		},
		{
			name:    "invalid ip value",
			asset:   domain.Asset{Name: "999.999.1.1", Type: domain.AssetTypeIP},
			wantErr: true,
		},
		{
			name:    "empty asset name",
			asset:   domain.Asset{Name: "", Type: domain.AssetTypeDomain},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.asset.Validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
		})
	}
}
