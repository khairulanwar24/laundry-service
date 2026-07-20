package services

import (
	"context"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	"laundry-service/repositories"
)

// DiscountService membungkus akses ke repository registry.
type DiscountService struct {
	repository repositories.IRepositoryRegistry
}

// IDiscountService adalah kontrak logika bisnis domain diskon.
type IDiscountService interface {
	ListActive(ctx context.Context, idOutlet string) response.Response
	Create(ctx context.Context, idOutlet string, form dto.DiscountForm) response.Response
	Update(ctx context.Context, idOutlet, idDiscount string, form dto.DiscountForm) response.Response
}

// NewDiscountService membuat instance DiscountService baru.
func NewDiscountService(repository repositories.IRepositoryRegistry) IDiscountService {
	return &DiscountService{repository: repository}
}

func (s *DiscountService) ListActive(ctx context.Context, idOutlet string) response.Response {
	data, err := s.repository.GetDiscount().ListActive(ctx, idOutlet)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar diskon: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *DiscountService) Create(ctx context.Context, idOutlet string, form dto.DiscountForm) response.Response {
	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}
	id, err := s.repository.GetDiscount().Create(ctx, idOutlet, form.Nama, form.Jenis, form.Nilai, form.Catatan, isActive)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat diskon: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Diskon berhasil dibuat", Data: map[string]interface{}{"id_discount": id}}
}

func (s *DiscountService) Update(ctx context.Context, idOutlet, idDiscount string, form dto.DiscountForm) response.Response {
	repo := s.repository.GetDiscount()

	if _, err := repo.FindByID(ctx, idOutlet, idDiscount); err != nil {
		return response.Response{Success: false, Message: "Diskon tidak ditemukan"}
	}

	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}
	if err := repo.Update(ctx, idDiscount, form.Nama, form.Jenis, form.Nilai, form.Catatan, isActive); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui diskon: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Diskon berhasil diperbarui"}
}
