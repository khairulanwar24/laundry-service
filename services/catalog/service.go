// Package services (catalog) berisi logika bisnis domain katalog: layanan,
// varian layanan, parfum, dan diskon.
package services

import (
	"context"
	"encoding/json"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	"laundry-service/repositories"
)

var defaultLangkahProses = []string{"cuci", "kering", "setrika"}

// ServiceService membungkus akses ke repository registry.
type ServiceService struct {
	repository repositories.IRepositoryRegistry
}

// IServiceService adalah kontrak logika bisnis domain layanan.
type IServiceService interface {
	ListActive(ctx context.Context, idOutlet, q string) response.Response
	Create(ctx context.Context, idOutlet string, form dto.ServiceForm) response.Response
	Update(ctx context.Context, idOutlet, idService string, form dto.ServiceForm) response.Response
	Deactivate(ctx context.Context, idOutlet, idService string) response.Response
}

// NewServiceService membuat instance ServiceService baru.
func NewServiceService(repository repositories.IRepositoryRegistry) IServiceService {
	return &ServiceService{repository: repository}
}

func (s *ServiceService) ListActive(ctx context.Context, idOutlet, q string) response.Response {
	data, err := s.repository.GetService().ListActive(ctx, idOutlet, q)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar layanan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *ServiceService) Create(ctx context.Context, idOutlet string, form dto.ServiceForm) response.Response {
	prioritas := 50
	if form.Prioritas != nil {
		prioritas = *form.Prioritas
	}
	langkah := form.LangkahProses
	if len(langkah) == 0 {
		langkah = defaultLangkahProses
	}
	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}
	langkahJSON, _ := json.Marshal(langkah)

	id, err := s.repository.GetService().Create(ctx, idOutlet, form.Nama, prioritas, string(langkahJSON), isActive)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat layanan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Layanan berhasil dibuat", Data: map[string]interface{}{"id_service": id}}
}

func (s *ServiceService) Update(ctx context.Context, idOutlet, idService string, form dto.ServiceForm) response.Response {
	repo := s.repository.GetService()

	if _, err := repo.FindByID(ctx, idOutlet, idService); err != nil {
		return response.Response{Success: false, Message: "Layanan tidak ditemukan"}
	}

	prioritas := 50
	if form.Prioritas != nil {
		prioritas = *form.Prioritas
	}
	langkah := form.LangkahProses
	if len(langkah) == 0 {
		langkah = defaultLangkahProses
	}
	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}
	langkahJSON, _ := json.Marshal(langkah)

	if err := repo.Update(ctx, idService, form.Nama, prioritas, string(langkahJSON), isActive); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui layanan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Layanan berhasil diperbarui"}
}

func (s *ServiceService) Deactivate(ctx context.Context, idOutlet, idService string) response.Response {
	repo := s.repository.GetService()

	if _, err := repo.FindByID(ctx, idOutlet, idService); err != nil {
		return response.Response{Success: false, Message: "Layanan tidak ditemukan"}
	}
	if err := repo.DeactivateVariantsByService(ctx, idService); err != nil {
		return response.Response{Success: false, Message: "Gagal menonaktifkan varian layanan: " + err.Error()}
	}
	if err := repo.Deactivate(ctx, idService); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus layanan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Layanan berhasil dihapus"}
}
