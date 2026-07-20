package services

import (
	"context"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	"laundry-service/repositories"
)

// ServiceVariantService membungkus akses ke repository registry.
type ServiceVariantService struct {
	repository repositories.IRepositoryRegistry
}

// IServiceVariantService adalah kontrak logika bisnis domain varian layanan.
type IServiceVariantService interface {
	Create(ctx context.Context, idOutlet, idService string, form dto.ServiceVariantForm) response.Response
	Update(ctx context.Context, idOutlet, idServiceVariant string, form dto.ServiceVariantForm) response.Response
}

// NewServiceVariantService membuat instance ServiceVariantService baru.
func NewServiceVariantService(repository repositories.IRepositoryRegistry) IServiceVariantService {
	return &ServiceVariantService{repository: repository}
}

func (s *ServiceVariantService) Create(ctx context.Context, idOutlet, idService string, form dto.ServiceVariantForm) response.Response {
	if _, err := s.repository.GetService().FindByID(ctx, idOutlet, idService); err != nil {
		return response.Response{Success: false, Message: "Layanan tidak ditemukan"}
	}

	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}

	id, err := s.repository.GetServiceVariant().Create(ctx, idService, form.Nama, form.Satuan, form.HargaPerSatuan, form.DurasiPengerjaanJam, form.GambarPath, form.Catatan, isActive)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat varian layanan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Varian layanan berhasil dibuat", Data: map[string]interface{}{"id_service_variant": id}}
}

func (s *ServiceVariantService) Update(ctx context.Context, idOutlet, idServiceVariant string, form dto.ServiceVariantForm) response.Response {
	repo := s.repository.GetServiceVariant()

	variant, err := repo.FindByID(ctx, idServiceVariant)
	if err != nil {
		return response.Response{Success: false, Message: "Varian layanan tidak ditemukan"}
	}
	if outletID, _ := variant["id_outlet"].(string); outletID != idOutlet {
		return response.Response{Success: false, Message: "Varian layanan tidak ditemukan"}
	}

	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}

	if err := repo.Update(ctx, idServiceVariant, form.Nama, form.Satuan, form.HargaPerSatuan, form.DurasiPengerjaanJam, form.GambarPath, form.Catatan, isActive); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui varian layanan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Varian layanan berhasil diperbarui"}
}
