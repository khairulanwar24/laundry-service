// Package osce adalah lapisan business logic untuk domain pengguna OSCE.
package osce

import (
	"context"
	"sso-service/common/response"
	dtos "sso-service/domain/dto/osce"
	"sso-service/repositories"
)

// OsceUserService membungkus repository registry.
type OsceUserService struct {
	repo repositories.IRepositoryRegistry
}

// IOsceUserService adalah kontrak business logic domain pengguna OSCE.
type IOsceUserService interface {
	// Penguji
	GetPenguji(ctx context.Context, form dtos.PengujiQuery) response.Response
	GetPengujiByID(ctx context.Context, id string) response.Response
	CreatePenguji(ctx context.Context, form dtos.PengujiForm) response.Response
	UpdatePenguji(ctx context.Context, id string, form dtos.PengujiForm) response.Response
	DeletePenguji(ctx context.Context, id string) response.Response

	// Kelompok Mahasiswa
	GetKelompokMahasiswa(ctx context.Context, form dtos.KelompokMahasiswaQuery) response.Response
	GetKelompokMahasiswaByID(ctx context.Context, id string) response.Response
	CreateKelompokMahasiswa(ctx context.Context, form dtos.KelompokMahasiswaForm) response.Response
	UpdateKelompokMahasiswa(ctx context.Context, id string, form dtos.KelompokMahasiswaForm) response.Response
	DeleteKelompokMahasiswa(ctx context.Context, id string) response.Response

	// Anggota Kelompok
	GetAnggotaByKelompok(ctx context.Context, idKelompok string) response.Response
	TambahAnggota(ctx context.Context, idKelompok string, form dtos.AnggotaKelompokForm) response.Response
	TambahAnggotaBulk(ctx context.Context, idKelompok string, form dtos.BulkAnggotaKelompokForm) response.Response
	HapusAnggota(ctx context.Context, id string) response.Response

	// Penugasan Penguji
	GetPenugasanBySesi(ctx context.Context, idSesi string) response.Response
	AssignPenguji(ctx context.Context, form dtos.PenugasanPengujiForm) response.Response
	HapusPenugasan(ctx context.Context, id string) response.Response
}

// NewOsceUserService membuat instance baru.
func NewOsceUserService(repo repositories.IRepositoryRegistry) IOsceUserService {
	return &OsceUserService{repo: repo}
}

// ===== PENGUJI =====

func (s *OsceUserService) GetPenguji(ctx context.Context, form dtos.PengujiQuery) response.Response {
	data, err := s.repo.GetOsceUser().FindPenguji(ctx, form.Limit, form.Offset, form.Order, form.Filter)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	count, _ := s.repo.GetOsceUser().CountPenguji(ctx, form.Filter)
	return response.Response{
		Success: true, Message: "sukses",
		Data:    map[string]interface{}{"data": data, "total_data": count},
	}
}

func (s *OsceUserService) GetPengujiByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceUser().FindPengujiByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceUserService) CreatePenguji(ctx context.Context, form dtos.PengujiForm) response.Response {
	// Cek duplikat
	existing, _ := s.repo.GetOsceUser().FindPengujiByUserID(ctx, form.IDUser)
	if existing != nil {
		return response.Response{Success: false, Message: "User sudah terdaftar sebagai penguji"}
	}
	if err := s.repo.GetOsceUser().InsertPenguji(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal mendaftarkan penguji: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Penguji berhasil didaftarkan"}
}

func (s *OsceUserService) UpdatePenguji(ctx context.Context, id string, form dtos.PengujiForm) response.Response {
	if err := s.repo.GetOsceUser().UpdatePenguji(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Data penguji berhasil diperbarui"}
}

func (s *OsceUserService) DeletePenguji(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceUser().DeletePenguji(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Penguji berhasil dinonaktifkan"}
}

// ===== KELOMPOK MAHASISWA =====

func (s *OsceUserService) GetKelompokMahasiswa(ctx context.Context, form dtos.KelompokMahasiswaQuery) response.Response {
	data, err := s.repo.GetOsceUser().FindKelompokMahasiswa(ctx, form.Limit, form.Offset, form.Order, form.Filter)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	count, _ := s.repo.GetOsceUser().CountKelompokMahasiswa(ctx, form.Filter)
	return response.Response{
		Success: true, Message: "sukses",
		Data:    map[string]interface{}{"data": data, "total_data": count},
	}
}

func (s *OsceUserService) GetKelompokMahasiswaByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceUser().FindKelompokMahasiswaByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceUserService) CreateKelompokMahasiswa(ctx context.Context, form dtos.KelompokMahasiswaForm) response.Response {
	id, err := s.repo.GetOsceUser().InsertKelompokMahasiswa(ctx, form)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat kelompok: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Kelompok berhasil dibuat", Data: map[string]string{"id": id}}
}

func (s *OsceUserService) UpdateKelompokMahasiswa(ctx context.Context, id string, form dtos.KelompokMahasiswaForm) response.Response {
	if err := s.repo.GetOsceUser().UpdateKelompokMahasiswa(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Kelompok berhasil diperbarui"}
}

func (s *OsceUserService) DeleteKelompokMahasiswa(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceUser().DeleteKelompokMahasiswa(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Kelompok berhasil dihapus"}
}

// ===== ANGGOTA KELOMPOK =====

func (s *OsceUserService) GetAnggotaByKelompok(ctx context.Context, idKelompok string) response.Response {
	data, err := s.repo.GetOsceUser().FindAnggotaByKelompok(ctx, idKelompok)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceUserService) TambahAnggota(ctx context.Context, idKelompok string, form dtos.AnggotaKelompokForm) response.Response {
	if err := s.repo.GetOsceUser().InsertAnggota(ctx, idKelompok, form.IDUser); err != nil {
		return response.Response{Success: false, Message: "Gagal menambah anggota: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Anggota berhasil ditambahkan"}
}

func (s *OsceUserService) TambahAnggotaBulk(ctx context.Context, idKelompok string, form dtos.BulkAnggotaKelompokForm) response.Response {
	berhasil, gagal, detailGagal, _ := s.repo.GetOsceUser().InsertAnggotaBulk(ctx, idKelompok, form.IDUsers)
	return response.Response{
		Success: true,
		Message: "Bulk insert selesai",
		Data: map[string]interface{}{
			"total":        len(form.IDUsers),
			"berhasil":     berhasil,
			"gagal":        gagal,
			"detail_gagal": detailGagal,
		},
	}
}

func (s *OsceUserService) HapusAnggota(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceUser().DeleteAnggota(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus anggota: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Anggota berhasil dihapus dari kelompok"}
}

// ===== PENUGASAN PENGUJI =====

func (s *OsceUserService) GetPenugasanBySesi(ctx context.Context, idSesi string) response.Response {
	data, err := s.repo.GetOsceUser().FindPenugasanBySesi(ctx, idSesi)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceUserService) AssignPenguji(ctx context.Context, form dtos.PenugasanPengujiForm) response.Response {
	if err := s.repo.GetOsceUser().InsertPenugasan(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal assign penguji: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Penguji berhasil di-assign ke station"}
}

func (s *OsceUserService) HapusPenugasan(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceUser().DeletePenugasan(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus penugasan: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Penugasan penguji berhasil dihapus"}
}
