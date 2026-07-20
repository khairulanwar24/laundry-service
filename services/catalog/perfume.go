package services

import (
	"context"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	"laundry-service/repositories"
)

// PerfumeService membungkus akses ke repository registry.
type PerfumeService struct {
	repository repositories.IRepositoryRegistry
}

// IPerfumeService adalah kontrak logika bisnis domain parfum.
type IPerfumeService interface {
	ListActive(ctx context.Context, idOutlet string) response.Response
	Create(ctx context.Context, idOutlet string, form dto.PerfumeForm) response.Response
	Update(ctx context.Context, idOutlet, idPerfume string, form dto.PerfumeForm) response.Response
}

// NewPerfumeService membuat instance PerfumeService baru.
func NewPerfumeService(repository repositories.IRepositoryRegistry) IPerfumeService {
	return &PerfumeService{repository: repository}
}

func (s *PerfumeService) ListActive(ctx context.Context, idOutlet string) response.Response {
	data, err := s.repository.GetPerfume().ListActive(ctx, idOutlet)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar parfum: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *PerfumeService) Create(ctx context.Context, idOutlet string, form dto.PerfumeForm) response.Response {
	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}
	id, err := s.repository.GetPerfume().Create(ctx, idOutlet, form.Nama, form.Catatan, isActive)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat parfum: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Parfum berhasil dibuat", Data: map[string]interface{}{"id_perfume": id}}
}

func (s *PerfumeService) Update(ctx context.Context, idOutlet, idPerfume string, form dto.PerfumeForm) response.Response {
	repo := s.repository.GetPerfume()

	if _, err := repo.FindByID(ctx, idOutlet, idPerfume); err != nil {
		return response.Response{Success: false, Message: "Parfum tidak ditemukan"}
	}

	isActive := true
	if form.IsActive != nil {
		isActive = *form.IsActive
	}
	if err := repo.Update(ctx, idPerfume, form.Nama, form.Catatan, isActive); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui parfum: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Parfum berhasil diperbarui"}
}
