// Package services (ref) berisi logika bisnis untuk domain referensi.
package services

import (
	"context"

	"laundry-service/repositories"
)

// RefService membungkus akses ke repository registry.
type RefService struct {
	repository repositories.IRepositoryRegistry
}

// IRefService adalah kontrak logika bisnis domain referensi.
type IRefService interface {
	GetMasterProdi(ctx context.Context) ([]map[string]any, error)
	GetAngkatan(ctx context.Context, idProdi string) ([]map[string]any, error)
}

// NewRefService membuat instance RefService baru.
func NewRefService(repository repositories.IRepositoryRegistry) IRefService {
	return &RefService{repository: repository}
}

// GetMasterProdi meneruskan permintaan daftar prodi ke repository.
func (s *RefService) GetMasterProdi(ctx context.Context) ([]map[string]any, error) {
	return s.repository.GetRef().GetMasterProdi(ctx)
}

// GetAngkatan meneruskan permintaan daftar angkatan ke repository.
func (s *RefService) GetAngkatan(ctx context.Context, idProdi string) ([]map[string]any, error) {
	return s.repository.GetRef().GetAngkatan(ctx, idProdi)
}
