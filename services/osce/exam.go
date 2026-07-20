// Package osce adalah lapisan business logic untuk domain OSCE ujian.
package osce

import (
	"context"
	"laundry-service/common/response"
	dtos "laundry-service/domain/dto/osce"
	"laundry-service/repositories"
)

// OsceExamService membungkus repository registry.
type OsceExamService struct {
	repo repositories.IRepositoryRegistry
}

// IOsceExamService adalah kontrak business logic domain ujian.
type IOsceExamService interface {
	// Tahun Akademik
	GetTahunAkademik(ctx context.Context) response.Response
	GetTahunAkademikByID(ctx context.Context, id string) response.Response
	CreateTahunAkademik(ctx context.Context, form dtos.TahunAkademikForm) response.Response
	UpdateTahunAkademik(ctx context.Context, id string, form dtos.TahunAkademikForm) response.Response
	DeleteTahunAkademik(ctx context.Context, id string) response.Response

	// Semester
	GetSemester(ctx context.Context) response.Response
	GetSemesterByID(ctx context.Context, id string) response.Response
	CreateSemester(ctx context.Context, form dtos.SemesterForm) response.Response
	UpdateSemester(ctx context.Context, id string, form dtos.SemesterForm) response.Response
	DeleteSemester(ctx context.Context, id string) response.Response

	// Program Studi
	GetProgramStudi(ctx context.Context) response.Response
	GetProgramStudiByID(ctx context.Context, id string) response.Response
	CreateProgramStudi(ctx context.Context, form dtos.ProgramStudiForm) response.Response
	UpdateProgramStudi(ctx context.Context, id string, form dtos.ProgramStudiForm) response.Response
	DeleteProgramStudi(ctx context.Context, id string) response.Response

	// Ujian
	GetUjian(ctx context.Context, form dtos.UjianQuery) response.Response
	GetUjianByID(ctx context.Context, id string) response.Response
	CreateUjian(ctx context.Context, form dtos.UjianForm) response.Response
	UpdateUjian(ctx context.Context, id string, form dtos.UjianForm) response.Response
	DeleteUjian(ctx context.Context, id string) response.Response
	UpdateStatusUjian(ctx context.Context, id, status string) response.Response

	// Ujian Station
	GetUjianStationByUjian(ctx context.Context, idUjian string) response.Response
	CreateUjianStation(ctx context.Context, form dtos.UjianStationForm) response.Response
	UpdateUjianStation(ctx context.Context, id string, form dtos.UjianStationForm) response.Response
	DeleteUjianStation(ctx context.Context, id string) response.Response

	// Sesi Ujian
	GetSesiByUjian(ctx context.Context, idUjian string) response.Response
	CreateSesi(ctx context.Context, form dtos.SesiUjianForm) response.Response
	UpdateSesi(ctx context.Context, id string, form dtos.SesiUjianForm) response.Response
	DeleteSesi(ctx context.Context, id string) response.Response

	// Rotasi Ujian
	GetRotasiBySesi(ctx context.Context, idSesi string) response.Response
	CreateRotasi(ctx context.Context, form dtos.RotasiUjianForm) response.Response
	UpdateRotasi(ctx context.Context, id string, form dtos.RotasiUjianForm) response.Response
	DeleteRotasi(ctx context.Context, id string) response.Response

	// Jadwal Ujian
	GetJadwalBySesi(ctx context.Context, idSesi string) response.Response
	GenerateJadwal(ctx context.Context, form dtos.GenerateJadwalRequest) response.Response
	CreateJadwal(ctx context.Context, form dtos.JadwalUjianForm) response.Response
	UpdateJadwal(ctx context.Context, id string, form dtos.JadwalUjianForm) response.Response
	DeleteJadwal(ctx context.Context, id string) response.Response
}

// NewOsceExamService membuat instance baru.
func NewOsceExamService(repo repositories.IRepositoryRegistry) IOsceExamService {
	return &OsceExamService{repo: repo}
}

// ===== TAHUN AKADEMIK =====

func (s *OsceExamService) GetTahunAkademik(ctx context.Context) response.Response {
	data, err := s.repo.GetOsceExam().FindTahunAkademik(ctx)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) GetTahunAkademikByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceExam().FindTahunAkademikByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) CreateTahunAkademik(ctx context.Context, form dtos.TahunAkademikForm) response.Response {
	if err := s.repo.GetOsceExam().InsertTahunAkademik(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat tahun akademik: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Tahun akademik berhasil dibuat"}
}

func (s *OsceExamService) UpdateTahunAkademik(ctx context.Context, id string, form dtos.TahunAkademikForm) response.Response {
	if err := s.repo.GetOsceExam().UpdateTahunAkademik(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Tahun akademik berhasil diperbarui"}
}

func (s *OsceExamService) DeleteTahunAkademik(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceExam().DeleteTahunAkademik(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Tahun akademik berhasil dihapus"}
}

// ===== SEMESTER =====

func (s *OsceExamService) GetSemester(ctx context.Context) response.Response {
	data, err := s.repo.GetOsceExam().FindSemester(ctx)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) GetSemesterByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceExam().FindSemesterByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) CreateSemester(ctx context.Context, form dtos.SemesterForm) response.Response {
	if err := s.repo.GetOsceExam().InsertSemester(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat semester: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Semester berhasil dibuat"}
}

func (s *OsceExamService) UpdateSemester(ctx context.Context, id string, form dtos.SemesterForm) response.Response {
	if err := s.repo.GetOsceExam().UpdateSemester(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Semester berhasil diperbarui"}
}

func (s *OsceExamService) DeleteSemester(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceExam().DeleteSemester(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Semester berhasil dihapus"}
}

// ===== PROGRAM STUDI =====

func (s *OsceExamService) GetProgramStudi(ctx context.Context) response.Response {
	data, err := s.repo.GetOsceExam().FindProgramStudi(ctx)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) GetProgramStudiByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceExam().FindProgramStudiByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) CreateProgramStudi(ctx context.Context, form dtos.ProgramStudiForm) response.Response {
	if err := s.repo.GetOsceExam().InsertProgramStudi(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat program studi: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Program studi berhasil dibuat"}
}

func (s *OsceExamService) UpdateProgramStudi(ctx context.Context, id string, form dtos.ProgramStudiForm) response.Response {
	if err := s.repo.GetOsceExam().UpdateProgramStudi(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Program studi berhasil diperbarui"}
}

func (s *OsceExamService) DeleteProgramStudi(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceExam().DeleteProgramStudi(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Program studi berhasil dihapus"}
}

// ===== UJIAN =====

func (s *OsceExamService) GetUjian(ctx context.Context, form dtos.UjianQuery) response.Response {
	data, err := s.repo.GetOsceExam().FindUjian(ctx, form.Limit, form.Offset, form.Order, form.Filter)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	count, _ := s.repo.GetOsceExam().CountUjian(ctx, form.Filter)
	return response.Response{
		Success: true, Message: "sukses",
		Data:    map[string]interface{}{"data": data, "total_data": count},
	}
}

func (s *OsceExamService) GetUjianByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceExam().FindUjianByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) CreateUjian(ctx context.Context, form dtos.UjianForm) response.Response {
	id, err := s.repo.GetOsceExam().InsertUjian(ctx, form)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat ujian: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Ujian berhasil dibuat", Data: map[string]string{"id": id}}
}

func (s *OsceExamService) UpdateUjian(ctx context.Context, id string, form dtos.UjianForm) response.Response {
	if err := s.repo.GetOsceExam().UpdateUjian(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Ujian berhasil diperbarui"}
}

func (s *OsceExamService) DeleteUjian(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceExam().DeleteUjian(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Ujian berhasil dihapus"}
}

func (s *OsceExamService) UpdateStatusUjian(ctx context.Context, id, status string) response.Response {
	if err := s.repo.GetOsceExam().UpdateStatusUjian(ctx, id, status); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui status: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Status ujian berhasil diperbarui"}
}

// ===== UJIAN STATION =====

func (s *OsceExamService) GetUjianStationByUjian(ctx context.Context, idUjian string) response.Response {
	data, err := s.repo.GetOsceExam().FindUjianStationByUjian(ctx, idUjian)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) CreateUjianStation(ctx context.Context, form dtos.UjianStationForm) response.Response {
	if err := s.repo.GetOsceExam().InsertUjianStation(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal menambah station: " + err.Error()}
	}
	_ = s.repo.GetOsceExam().UpdateJumlahStationUjian(ctx, form.IDUjian)
	return response.Response{Success: true, Message: "Station berhasil ditambahkan ke ujian"}
}

func (s *OsceExamService) UpdateUjianStation(ctx context.Context, id string, form dtos.UjianStationForm) response.Response {
	if err := s.repo.GetOsceExam().UpdateUjianStation(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Station ujian berhasil diperbarui"}
}

func (s *OsceExamService) DeleteUjianStation(ctx context.Context, id string) response.Response {
	data, _ := s.repo.GetOsceExam().FindUjianStationByID(ctx, id)
	if err := s.repo.GetOsceExam().DeleteUjianStation(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	if data != nil {
		if idUjian, ok := data["id_ujian"].(string); ok {
			_ = s.repo.GetOsceExam().UpdateJumlahStationUjian(ctx, idUjian)
		}
	}
	return response.Response{Success: true, Message: "Station berhasil dihapus dari ujian"}
}

// ===== SESI UJIAN =====

func (s *OsceExamService) GetSesiByUjian(ctx context.Context, idUjian string) response.Response {
	data, err := s.repo.GetOsceExam().FindSesiByUjian(ctx, idUjian)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) CreateSesi(ctx context.Context, form dtos.SesiUjianForm) response.Response {
	if err := s.repo.GetOsceExam().InsertSesi(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat sesi: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Sesi ujian berhasil dibuat"}
}

func (s *OsceExamService) UpdateSesi(ctx context.Context, id string, form dtos.SesiUjianForm) response.Response {
	if err := s.repo.GetOsceExam().UpdateSesi(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Sesi ujian berhasil diperbarui"}
}

func (s *OsceExamService) DeleteSesi(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceExam().DeleteSesi(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Sesi ujian berhasil dihapus"}
}

// ===== ROTASI UJIAN =====

func (s *OsceExamService) GetRotasiBySesi(ctx context.Context, idSesi string) response.Response {
	data, err := s.repo.GetOsceExam().FindRotasiBySesi(ctx, idSesi)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) CreateRotasi(ctx context.Context, form dtos.RotasiUjianForm) response.Response {
	if err := s.repo.GetOsceExam().InsertRotasi(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat rotasi: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Rotasi berhasil dibuat"}
}

func (s *OsceExamService) UpdateRotasi(ctx context.Context, id string, form dtos.RotasiUjianForm) response.Response {
	if err := s.repo.GetOsceExam().UpdateRotasi(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Rotasi berhasil diperbarui"}
}

func (s *OsceExamService) DeleteRotasi(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceExam().DeleteRotasi(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Rotasi berhasil dihapus"}
}

// ===== JADWAL UJIAN =====

func (s *OsceExamService) GetJadwalBySesi(ctx context.Context, idSesi string) response.Response {
	data, err := s.repo.GetOsceExam().FindJadwalBySesi(ctx, idSesi)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceExamService) GenerateJadwal(ctx context.Context, form dtos.GenerateJadwalRequest) response.Response {
	// Ambil daftar ujian station untuk ujian ini
	ujianStations, err := s.repo.GetOsceExam().FindUjianStationByUjian(ctx, form.IDUjian)
	if err != nil || len(ujianStations) == 0 {
		return response.Response{Success: false, Message: "Ujian tidak memiliki station. Tambah station terlebih dahulu."}
	}

	// Ambil anggota dari kelompok mahasiswa
	anggota, err := s.repo.GetOsceUser().FindAnggotaByKelompok(ctx, form.IDKelompokMahasiswa)
	if err != nil || len(anggota) == 0 {
		return response.Response{Success: false, Message: "Kelompok mahasiswa kosong atau tidak ditemukan."}
	}

	// Validasi sesi ada
	_, err = s.repo.GetOsceExam().FindSesiByID(ctx, form.IDSesiUjian)
	if err != nil {
		return response.Response{Success: false, Message: "Sesi tidak ditemukan."}
	}

	// Generate jadwal: setiap mahasiswa masuk ke tiap station
	generated := 0
	for _, m := range anggota {
		idUser, _ := m["id_user"].(string)
		for i, es := range ujianStations {
			idUjianStation, _ := es["id_ujian_station"].(string)
			urutan := i + 1

			jadwalForm := dtos.JadwalUjianForm{
				IDUjian:        form.IDUjian,
				IDSesiUjian:    form.IDSesiUjian,
				IDUjianStation: idUjianStation,
				IDUser:         idUser,
				UrutanMasuk:    urutan,
			}
			if err := s.repo.GetOsceExam().InsertJadwal(ctx, jadwalForm); err != nil {
				continue // skip jika sudah ada (unique constraint)
			}
			generated++
		}
	}

	return response.Response{
		Success: true,
		Message: "Penjadwalan selesai",
		Data:    map[string]interface{}{"total_mahasiswa": len(anggota), "total_jadwal_dibuat": generated},
	}
}

func (s *OsceExamService) CreateJadwal(ctx context.Context, form dtos.JadwalUjianForm) response.Response {
	if err := s.repo.GetOsceExam().InsertJadwal(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat jadwal: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Jadwal berhasil dibuat"}
}

func (s *OsceExamService) UpdateJadwal(ctx context.Context, id string, form dtos.JadwalUjianForm) response.Response {
	if err := s.repo.GetOsceExam().UpdateJadwal(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Jadwal berhasil diperbarui"}
}

func (s *OsceExamService) DeleteJadwal(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceExam().DeleteJadwal(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Jadwal berhasil dihapus"}
}
