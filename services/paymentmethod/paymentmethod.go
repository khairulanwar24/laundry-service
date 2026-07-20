// Package services (paymentmethod) berisi logika bisnis domain metode pembayaran.
package services

import (
	"context"
	"encoding/json"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	"laundry-service/repositories"
)

// PaymentMethodService membungkus akses ke repository registry.
type PaymentMethodService struct {
	repository repositories.IRepositoryRegistry
}

// IPaymentMethodService adalah kontrak logika bisnis domain metode pembayaran.
type IPaymentMethodService interface {
	ListActive(ctx context.Context, idOutlet string) response.Response
	Create(ctx context.Context, idOutlet string, form dto.PaymentMethodForm) response.Response
	Update(ctx context.Context, idOutlet, idPaymentMethod string, form dto.PaymentMethodForm) response.Response
	Deactivate(ctx context.Context, idOutlet, idPaymentMethod string) response.Response
}

// NewPaymentMethodService membuat instance PaymentMethodService baru.
func NewPaymentMethodService(repository repositories.IRepositoryRegistry) IPaymentMethodService {
	return &PaymentMethodService{repository: repository}
}

func (s *PaymentMethodService) ListActive(ctx context.Context, idOutlet string) response.Response {
	data, err := s.repository.GetPaymentMethod().ListActive(ctx, idOutlet)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar metode pembayaran: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *PaymentMethodService) Create(ctx context.Context, idOutlet string, form dto.PaymentMethodForm) response.Response {
	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}
	tags := form.Tags
	if tags == nil {
		tags = []string{}
	}
	tagsJSON, _ := json.Marshal(tags)

	id, err := s.repository.GetPaymentMethod().Create(ctx, idOutlet, form.Kategori, form.Nama, form.Logo, form.NamaPemilik, string(tagsJSON), isActive)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat metode pembayaran: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Metode pembayaran berhasil dibuat", Data: map[string]interface{}{"id_payment_method": id}}
}

func (s *PaymentMethodService) Update(ctx context.Context, idOutlet, idPaymentMethod string, form dto.PaymentMethodForm) response.Response {
	repo := s.repository.GetPaymentMethod()

	if _, err := repo.FindByID(ctx, idOutlet, idPaymentMethod); err != nil {
		return response.Response{Success: false, Message: "Metode pembayaran tidak ditemukan"}
	}

	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}
	tags := form.Tags
	if tags == nil {
		tags = []string{}
	}
	tagsJSON, _ := json.Marshal(tags)

	if err := repo.Update(ctx, idPaymentMethod, form.Kategori, form.Nama, form.Logo, form.NamaPemilik, string(tagsJSON), isActive); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui metode pembayaran: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Metode pembayaran berhasil diperbarui"}
}

func (s *PaymentMethodService) Deactivate(ctx context.Context, idOutlet, idPaymentMethod string) response.Response {
	repo := s.repository.GetPaymentMethod()

	if _, err := repo.FindByID(ctx, idOutlet, idPaymentMethod); err != nil {
		return response.Response{Success: false, Message: "Metode pembayaran tidak ditemukan"}
	}
	if err := repo.Deactivate(ctx, idPaymentMethod); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus metode pembayaran: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Metode pembayaran berhasil dihapus"}
}
