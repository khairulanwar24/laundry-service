// Package services (masterapp) berisi logika bisnis domain master aplikasi.
package services

import (
	"context"
	"time"

	"sso-service/common/response"
	"sso-service/repositories"
)

// MasterAppService membungkus akses ke repository registry.
type MasterAppService struct {
	repository repositories.IRepositoryRegistry
}

// IMasterAppService adalah kontrak logika bisnis domain master aplikasi.
type IMasterAppService interface {
	GetMasterApps(order, filter string, limit, offset int) response.Response
	GetMasterAppById(ctx context.Context, id string) response.Response
	CreateMasterApp(
		ctx context.Context,
		namaAplikasi,
		deskripsi string,
		tglVersion time.Time,
		url,
		versiAplikasi,
		image string,
	) response.Response
	UpdateMasterApp(
		ctx context.Context,
		id,
		image,
		namaAplikasi,
		deskripsi,
		versiAplikasi string,
		tglVersion time.Time,
		url string,
	) response.Response
	DeleteMasterApp(ctx context.Context, id string) response.Response
}

// NewMasterAppService membuat instance MasterAppService baru.
func NewMasterAppService(repository repositories.IRepositoryRegistry) IMasterAppService {
	return &MasterAppService{repository: repository}
}

// GetMasterApps mengembalikan daftar master aplikasi (datatable).
func (s *MasterAppService) GetMasterApps(order, filter string, limit, offset int) response.Response {
	data := s.repository.GetMasterApp().GetMasterApps(order, filter, limit, offset)
	return response.Response{Success: true, Message: "Success", Data: data}
}

// GetMasterAppById mengambil detail satu master aplikasi.
func (s *MasterAppService) GetMasterAppById(ctx context.Context, id string) response.Response {
	master, rows, err := s.repository.GetMasterApp().GetMasterAppById(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mendapatkan aplikasi"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "aplikasi tidak tersedia"}
	}
	return response.Response{Success: true, Message: "Success", Data: master}
}

// CreateMasterApp membuat master aplikasi baru.
func (s *MasterAppService) CreateMasterApp(
	ctx context.Context,
	namaAplikasi,
	deskripsi string,
	tglVersion time.Time,
	url,
	versiAplikasi,
	image string,
) response.Response {
	if err := s.repository.GetMasterApp().Create(ctx, namaAplikasi, deskripsi, tglVersion, url, versiAplikasi, image); err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}
	return response.Response{Success: true, Message: "Success"}
}

// UpdateMasterApp memperbarui master aplikasi (dengan / tanpa image).
func (s *MasterAppService) UpdateMasterApp(
	ctx context.Context,
	id,
	image,
	namaAplikasi,
	deskripsi,
	versiAplikasi string,
	tglVersion time.Time,
	url string,
) response.Response {
	var rows int64
	var err error
	if image == "" {
		rows, err = s.repository.GetMasterApp().UpdateWithoutImage(ctx, namaAplikasi, deskripsi, versiAplikasi, tglVersion, url, id)
	} else {
		rows, err = s.repository.GetMasterApp().UpdateWithImage(ctx, image, namaAplikasi, deskripsi, versiAplikasi, tglVersion, url, id)
	}
	if err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui aplikasi"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "Aplikasi tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "Aplikasi berhasil diperbarui"}
}

// DeleteMasterApp menghapus master aplikasi.
func (s *MasterAppService) DeleteMasterApp(ctx context.Context, id string) response.Response {
	rows, err := s.repository.GetMasterApp().Delete(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Delete master aplikasi"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "Master aplikasi tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "Success"}
}
