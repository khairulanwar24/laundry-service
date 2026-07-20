// Package osce adalah lapisan akses data untuk domain OSCE ujian.
package osce

import (
	"context"
	"laundry-service/domain/dto/osce"

	"gorm.io/gorm"
)

// OsceExamRepository memegang koneksi database utama.
type OsceExamRepository struct {
	db *gorm.DB
}

// IOsceExamRepository adalah kontrak akses data domain ujian OSCE.
type IOsceExamRepository interface {
	// Tahun Akademik
	FindTahunAkademik(ctx context.Context) ([]map[string]interface{}, error)
	FindTahunAkademikByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertTahunAkademik(ctx context.Context, form osce.TahunAkademikForm) error
	UpdateTahunAkademik(ctx context.Context, id string, form osce.TahunAkademikForm) error
	DeleteTahunAkademik(ctx context.Context, id string) error

	// Semester
	FindSemester(ctx context.Context) ([]map[string]interface{}, error)
	FindSemesterByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertSemester(ctx context.Context, form osce.SemesterForm) error
	UpdateSemester(ctx context.Context, id string, form osce.SemesterForm) error
	DeleteSemester(ctx context.Context, id string) error

	// Program Studi
	FindProgramStudi(ctx context.Context) ([]map[string]interface{}, error)
	FindProgramStudiByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertProgramStudi(ctx context.Context, form osce.ProgramStudiForm) error
	UpdateProgramStudi(ctx context.Context, id string, form osce.ProgramStudiForm) error
	DeleteProgramStudi(ctx context.Context, id string) error

	// Ujian
	FindUjian(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error)
	CountUjian(ctx context.Context, filter string) (int64, error)
	FindUjianByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertUjian(ctx context.Context, form osce.UjianForm) (string, error)
	UpdateUjian(ctx context.Context, id string, form osce.UjianForm) error
	DeleteUjian(ctx context.Context, id string) error
	UpdateStatusUjian(ctx context.Context, id, status string) error
	UpdateJumlahStationUjian(ctx context.Context, idUjian string) error

	// Ujian Station
	FindUjianStationByUjian(ctx context.Context, idUjian string) ([]map[string]interface{}, error)
	FindUjianStationByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertUjianStation(ctx context.Context, form osce.UjianStationForm) error
	UpdateUjianStation(ctx context.Context, id string, form osce.UjianStationForm) error
	DeleteUjianStation(ctx context.Context, id string) error

	// Sesi Ujian
	FindSesiByUjian(ctx context.Context, idUjian string) ([]map[string]interface{}, error)
	FindSesiByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertSesi(ctx context.Context, form osce.SesiUjianForm) error
	UpdateSesi(ctx context.Context, id string, form osce.SesiUjianForm) error
	DeleteSesi(ctx context.Context, id string) error
	UpdateStatusSesi(ctx context.Context, id, status string) error

	// Rotasi Ujian
	FindRotasiBySesi(ctx context.Context, idSesi string) ([]map[string]interface{}, error)
	InsertRotasi(ctx context.Context, form osce.RotasiUjianForm) error
	UpdateRotasi(ctx context.Context, id string, form osce.RotasiUjianForm) error
	DeleteRotasi(ctx context.Context, id string) error

	// Jadwal Ujian
	FindJadwalBySesi(ctx context.Context, idSesi string) ([]map[string]interface{}, error)
	FindJadwalByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertJadwal(ctx context.Context, form osce.JadwalUjianForm) error
	UpdateJadwal(ctx context.Context, id string, form osce.JadwalUjianForm) error
	DeleteJadwal(ctx context.Context, id string) error
	FindJadwalByPenguji(ctx context.Context, idPenguji string) ([]map[string]interface{}, error)
}

// NewOsceExamRepository membuat instance baru.
func NewOsceExamRepository(db *gorm.DB) IOsceExamRepository {
	return &OsceExamRepository{db: db}
}

// ===== TAHUN AKADEMIK =====

func (r *OsceExamRepository) FindTahunAkademik(ctx context.Context) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.tahun_akademik").
		Where("status_data = true").Order("nama_tahun_akademik DESC").Find(&data)
	return data, result.Error
}

func (r *OsceExamRepository) FindTahunAkademikByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.tahun_akademik").
		Where("id_tahun_akademik = ? AND status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceExamRepository) InsertTahunAkademik(ctx context.Context, form osce.TahunAkademikForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.tahun_akademik (nama_tahun_akademik, tgl_mulai, tgl_selesai, status_aktif) VALUES (?, ?, ?, ?)`,
		form.Nama, form.TglMulai, form.TglSelesai, form.StatusAktif,
	).Error
}

func (r *OsceExamRepository) UpdateTahunAkademik(ctx context.Context, id string, form osce.TahunAkademikForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.tahun_akademik
		 SET nama_tahun_akademik = ?, tgl_mulai = ?, tgl_selesai = ?, status_aktif = ?, tgl_update = NOW()
		 WHERE id_tahun_akademik = ?`,
		form.Nama, form.TglMulai, form.TglSelesai, form.StatusAktif, id,
	).Error
}

func (r *OsceExamRepository) DeleteTahunAkademik(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.tahun_akademik SET status_data = false, tgl_update = NOW() WHERE id_tahun_akademik = ?`, id,
	).Error
}

// ===== SEMESTER =====

func (r *OsceExamRepository) FindSemester(ctx context.Context) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.semester s").
		Select("s.*, ta.nama_tahun_akademik").
		Joins("JOIN osce.tahun_akademik ta ON s.id_tahun_akademik = ta.id_tahun_akademik").
		Where("s.status_data = true").Order("ta.nama_tahun_akademik DESC, s.nama_semester ASC").Find(&data)
	return data, result.Error
}

func (r *OsceExamRepository) FindSemesterByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.semester s").
		Select("s.*, ta.nama_tahun_akademik").
		Joins("JOIN osce.tahun_akademik ta ON s.id_tahun_akademik = ta.id_tahun_akademik").
		Where("s.id_semester = ? AND s.status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceExamRepository) InsertSemester(ctx context.Context, form osce.SemesterForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.semester (id_tahun_akademik, nama_semester, tgl_mulai, tgl_selesai, status_aktif)
		 VALUES (?, ?, ?, ?, ?)`,
		form.IDTahunAkademik, form.NamaSemester, form.TglMulai, form.TglSelesai, form.StatusAktif,
	).Error
}

func (r *OsceExamRepository) UpdateSemester(ctx context.Context, id string, form osce.SemesterForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.semester
		 SET nama_semester = ?, tgl_mulai = ?, tgl_selesai = ?, status_aktif = ?, tgl_update = NOW()
		 WHERE id_semester = ?`,
		form.NamaSemester, form.TglMulai, form.TglSelesai, form.StatusAktif, id,
	).Error
}

func (r *OsceExamRepository) DeleteSemester(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.semester SET status_data = false, tgl_update = NOW() WHERE id_semester = ?`, id,
	).Error
}

// ===== PROGRAM STUDI =====

func (r *OsceExamRepository) FindProgramStudi(ctx context.Context) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.program_studi").
		Where("status_data = true").Order("nama_prodi ASC").Find(&data)
	return data, result.Error
}

func (r *OsceExamRepository) FindProgramStudiByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.program_studi").
		Where("id_program_studi = ? AND status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceExamRepository) InsertProgramStudi(ctx context.Context, form osce.ProgramStudiForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.program_studi (id_prodi, kode_prodi, nama_prodi, jenjang) VALUES (?, ?, ?, ?)`,
		form.IDProdi, form.KodeProdi, form.NamaProdi, form.Jenjang,
	).Error
}

func (r *OsceExamRepository) UpdateProgramStudi(ctx context.Context, id string, form osce.ProgramStudiForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.program_studi SET kode_prodi = ?, nama_prodi = ?, jenjang = ?, tgl_update = NOW()
		 WHERE id_program_studi = ?`,
		form.KodeProdi, form.NamaProdi, form.Jenjang, id,
	).Error
}

func (r *OsceExamRepository) DeleteProgramStudi(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.program_studi SET status_data = false, tgl_update = NOW() WHERE id_program_studi = ?`, id,
	).Error
}

// ===== UJIAN =====

func (r *OsceExamRepository) FindUjian(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Table("osce.ujian e").
		Select("e.*, ta.nama_tahun_akademik, s.nama_semester, ps.nama_prodi").
		Joins("JOIN osce.tahun_akademik ta ON e.id_tahun_akademik = ta.id_tahun_akademik").
		Joins("JOIN osce.semester s ON e.id_semester = s.id_semester").
		Joins("JOIN osce.program_studi ps ON e.id_program_studi = ps.id_program_studi").
		Where("e.status_data = true")
	if filter != "" {
		query = query.Where("e.nama_ujian ILIKE ?", "%"+filter+"%")
	}
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if order != "" {
		query = query.Order(order)
	}
	result := query.Find(&data)
	return data, result.Error
}

func (r *OsceExamRepository) CountUjian(ctx context.Context, filter string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Table("osce.ujian").Where("status_data = true")
	if filter != "" {
		query = query.Where("nama_ujian ILIKE ?", "%"+filter+"%")
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *OsceExamRepository) FindUjianByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.ujian e").
		Select("e.*, ta.nama_tahun_akademik, s.nama_semester, ps.nama_prodi").
		Joins("JOIN osce.tahun_akademik ta ON e.id_tahun_akademik = ta.id_tahun_akademik").
		Joins("JOIN osce.semester s ON e.id_semester = s.id_semester").
		Joins("JOIN osce.program_studi ps ON e.id_program_studi = ps.id_program_studi").
		Where("e.id_ujian = ? AND e.status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceExamRepository) InsertUjian(ctx context.Context, form osce.UjianForm) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(
		`INSERT INTO osce.ujian
		 (id_tahun_akademik, id_semester, id_program_studi, nama_ujian, deskripsi, tgl_mulai, tgl_selesai, durasi_per_station, batas_lulus)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id_ujian`,
		form.IDTahunAkademik, form.IDSemester, form.IDProgramStudi, form.NamaUjian, form.Deskripsi,
		form.TglMulai, form.TglSelesai, form.DurasiPerStation, form.BatasLulus,
	).Scan(&id)
	return id, result.Error
}

func (r *OsceExamRepository) UpdateUjian(ctx context.Context, id string, form osce.UjianForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.ujian
		 SET nama_ujian = ?, deskripsi = ?, tgl_mulai = ?, tgl_selesai = ?, durasi_per_station = ?, batas_lulus = ?, tgl_update = NOW()
		 WHERE id_ujian = ?`,
		form.NamaUjian, form.Deskripsi, form.TglMulai, form.TglSelesai, form.DurasiPerStation, form.BatasLulus, id,
	).Error
}

func (r *OsceExamRepository) DeleteUjian(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.ujian SET status_data = false, tgl_update = NOW() WHERE id_ujian = ?`, id,
	).Error
}

func (r *OsceExamRepository) UpdateStatusUjian(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.ujian SET status_ujian = ?, tgl_update = NOW() WHERE id_ujian = ?`, status, id,
	).Error
}

func (r *OsceExamRepository) UpdateJumlahStationUjian(ctx context.Context, idUjian string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.ujian
		 SET jumlah_station = (SELECT COUNT(*) FROM osce.ujian_station WHERE id_ujian = ? AND status_data = true), tgl_update = NOW()
		 WHERE id_ujian = ?`,
		idUjian, idUjian,
	).Error
}

// ===== UJIAN STATION =====

func (r *OsceExamRepository) FindUjianStationByUjian(ctx context.Context, idUjian string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.ujian_station us").
		Select("us.*, s.kode_station, s.nama_station, st.nama_tipe").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Joins("JOIN osce.tipe_station st ON s.id_tipe_station = st.id_tipe_station").
		Where("us.id_ujian = ? AND us.status_data = true", idUjian).
		Order("us.urutan ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceExamRepository) FindUjianStationByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.ujian_station us").
		Select("us.*, s.kode_station, s.nama_station, st.nama_tipe").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Joins("JOIN osce.tipe_station st ON s.id_tipe_station = st.id_tipe_station").
		Where("us.id_ujian_station = ? AND us.status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceExamRepository) InsertUjianStation(ctx context.Context, form osce.UjianStationForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.ujian_station (id_ujian, id_station, urutan, durasi, bobot) VALUES (?, ?, ?, ?, ?)`,
		form.IDUjian, form.IDStation, form.Urutan, form.Durasi, form.Bobot,
	).Error
}

func (r *OsceExamRepository) UpdateUjianStation(ctx context.Context, id string, form osce.UjianStationForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.ujian_station SET urutan = ?, durasi = ?, bobot = ?, tgl_update = NOW() WHERE id_ujian_station = ?`,
		form.Urutan, form.Durasi, form.Bobot, id,
	).Error
}

func (r *OsceExamRepository) DeleteUjianStation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.ujian_station SET status_data = false, tgl_update = NOW() WHERE id_ujian_station = ?`, id,
	).Error
}

// ===== SESI UJIAN =====

func (r *OsceExamRepository) FindSesiByUjian(ctx context.Context, idUjian string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.sesi_ujian").
		Where("id_ujian = ? AND status_data = true", idUjian).
		Order("tgl_sesi ASC, jam_mulai ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceExamRepository) FindSesiByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.sesi_ujian").
		Where("id_sesi_ujian = ? AND status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceExamRepository) InsertSesi(ctx context.Context, form osce.SesiUjianForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.sesi_ujian (id_ujian, nama_sesi, tgl_sesi, jam_mulai, jam_selesai, lokasi, kuota_peserta)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		form.IDUjian, form.NamaSesi, form.TglSesi, form.JamMulai, form.JamSelesai, form.Lokasi, form.KuotaPeserta,
	).Error
}

func (r *OsceExamRepository) UpdateSesi(ctx context.Context, id string, form osce.SesiUjianForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.sesi_ujian
		 SET nama_sesi = ?, tgl_sesi = ?, jam_mulai = ?, jam_selesai = ?, lokasi = ?, kuota_peserta = ?, tgl_update = NOW()
		 WHERE id_sesi_ujian = ?`,
		form.NamaSesi, form.TglSesi, form.JamMulai, form.JamSelesai, form.Lokasi, form.KuotaPeserta, id,
	).Error
}

func (r *OsceExamRepository) DeleteSesi(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.sesi_ujian SET status_data = false, tgl_update = NOW() WHERE id_sesi_ujian = ?`, id,
	).Error
}

func (r *OsceExamRepository) UpdateStatusSesi(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.sesi_ujian SET status_sesi = ?, tgl_update = NOW() WHERE id_sesi_ujian = ?`, status, id,
	).Error
}

// ===== ROTASI UJIAN =====

func (r *OsceExamRepository) FindRotasiBySesi(ctx context.Context, idSesi string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.rotasi_ujian ru").
		Select("ru.*, s.kode_station, s.nama_station").
		Joins("JOIN osce.ujian_station us ON ru.id_ujian_station = us.id_ujian_station").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Where("ru.id_sesi_ujian = ? AND ru.status_data = true", idSesi).
		Order("ru.urutan_rotasi ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceExamRepository) InsertRotasi(ctx context.Context, form osce.RotasiUjianForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.rotasi_ujian
		 (id_ujian, id_sesi_ujian, id_ujian_station, urutan_rotasi, durasi, jeda_antar_station)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		form.IDUjian, form.IDSesiUjian, form.IDUjianStation, form.UrutanRotasi, form.Durasi, form.JedaAntarStation,
	).Error
}

func (r *OsceExamRepository) UpdateRotasi(ctx context.Context, id string, form osce.RotasiUjianForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.rotasi_ujian
		 SET id_ujian_station = ?, urutan_rotasi = ?, durasi = ?, jeda_antar_station = ?, tgl_update = NOW()
		 WHERE id_rotasi_ujian = ?`,
		form.IDUjianStation, form.UrutanRotasi, form.Durasi, form.JedaAntarStation, id,
	).Error
}

func (r *OsceExamRepository) DeleteRotasi(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.rotasi_ujian SET status_data = false, tgl_update = NOW() WHERE id_rotasi_ujian = ?`, id,
	).Error
}

// ===== JADWAL UJIAN =====

func (r *OsceExamRepository) FindJadwalBySesi(ctx context.Context, idSesi string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.jadwal_ujian ju").
		Select("ju.*, s.kode_station, s.nama_station, u.nama_lengkap, u.username, p2.nama_lengkap as nama_penguji").
		Joins("JOIN osce.ujian_station us ON ju.id_ujian_station = us.id_ujian_station").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Joins("JOIN public.users u ON ju.id_user = u.id_user").
		Joins("LEFT JOIN osce.penguji p ON ju.id_penguji = p.id_penguji").
		Joins("LEFT JOIN public.users p2 ON p.id_user = p2.id_user").
		Where("ju.id_sesi_ujian = ? AND ju.status_data = true", idSesi).
		Order("ju.urutan_masuk ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceExamRepository) FindJadwalByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.jadwal_ujian ju").
		Select("ju.*, s.kode_station, s.nama_station, u.nama_lengkap, u.username").
		Joins("JOIN osce.ujian_station us ON ju.id_ujian_station = us.id_ujian_station").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Joins("JOIN public.users u ON ju.id_user = u.id_user").
		Where("ju.id_jadwal_ujian = ? AND ju.status_data = true", id).
		First(&data)
	return data, result.Error
}

func (r *OsceExamRepository) InsertJadwal(ctx context.Context, form osce.JadwalUjianForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.jadwal_ujian
		 (id_ujian, id_sesi_ujian, id_ujian_station, id_user, id_penguji, urutan_masuk, jam_mulai, jam_selesai)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		form.IDUjian, form.IDSesiUjian, form.IDUjianStation, form.IDUser,
		form.IDPenguji, form.UrutanMasuk, form.JamMulai, form.JamSelesai,
	).Error
}

func (r *OsceExamRepository) UpdateJadwal(ctx context.Context, id string, form osce.JadwalUjianForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.jadwal_ujian
		 SET id_penguji = ?, urutan_masuk = ?, jam_mulai = ?, jam_selesai = ?, tgl_update = NOW()
		 WHERE id_jadwal_ujian = ?`,
		form.IDPenguji, form.UrutanMasuk, form.JamMulai, form.JamSelesai, id,
	).Error
}

func (r *OsceExamRepository) DeleteJadwal(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.jadwal_ujian SET status_data = false, tgl_update = NOW() WHERE id_jadwal_ujian = ?`, id,
	).Error
}

func (r *OsceExamRepository) FindJadwalByPenguji(ctx context.Context, idPenguji string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.jadwal_ujian ju").
		Select("ju.*, s.kode_station, s.nama_station, u.nama_lengkap, u.username, su.nama_sesi, su.tgl_sesi").
		Joins("JOIN osce.ujian_station us ON ju.id_ujian_station = us.id_ujian_station").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Joins("JOIN public.users u ON ju.id_user = u.id_user").
		Joins("JOIN osce.sesi_ujian su ON ju.id_sesi_ujian = su.id_sesi_ujian").
		Where("ju.id_penguji = ? AND ju.status_data = true", idPenguji).
		Order("su.tgl_sesi ASC, ju.urutan_masuk ASC").
		Find(&data)
	return data, result.Error
}
