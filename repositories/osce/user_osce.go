// Package osce adalah lapisan akses data untuk domain pengguna OSCE.
package osce

import (
	"context"
	"laundry-service/domain/dto/osce"

	"gorm.io/gorm"
)

// OsceUserRepository memegang koneksi database utama & akademik.
type OsceUserRepository struct {
	db         *gorm.DB
	dbAkademik *gorm.DB
}

// IOsceUserRepository adalah kontrak akses data domain pengguna OSCE.
type IOsceUserRepository interface {
	// Penguji
	FindPenguji(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error)
	CountPenguji(ctx context.Context, filter string) (int64, error)
	FindPengujiByID(ctx context.Context, id string) (map[string]interface{}, error)
	FindPengujiByUserID(ctx context.Context, idUser string) (map[string]interface{}, error)
	InsertPenguji(ctx context.Context, form osce.PengujiForm) error
	UpdatePenguji(ctx context.Context, id string, form osce.PengujiForm) error
	DeletePenguji(ctx context.Context, id string) error

	// Kelompok Mahasiswa
	FindKelompokMahasiswa(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error)
	CountKelompokMahasiswa(ctx context.Context, filter string) (int64, error)
	FindKelompokMahasiswaByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertKelompokMahasiswa(ctx context.Context, form osce.KelompokMahasiswaForm) (string, error)
	UpdateKelompokMahasiswa(ctx context.Context, id string, form osce.KelompokMahasiswaForm) error
	DeleteKelompokMahasiswa(ctx context.Context, id string) error

	// Anggota Kelompok
	FindAnggotaByKelompok(ctx context.Context, idKelompok string) ([]map[string]interface{}, error)
	InsertAnggota(ctx context.Context, idKelompok, idUser string) error
	InsertAnggotaBulk(ctx context.Context, idKelompok string, idUsers []string) (int, int, []string, error)
	DeleteAnggota(ctx context.Context, id string) error
	CountAnggotaByKelompok(ctx context.Context, idKelompok string) (int64, error)

	// Penugasan Penguji
	FindPenugasanBySesi(ctx context.Context, idSesi string) ([]map[string]interface{}, error)
	InsertPenugasan(ctx context.Context, form osce.PenugasanPengujiForm) error
	DeletePenugasan(ctx context.Context, id string) error
}

// NewOsceUserRepository membuat instance baru.
func NewOsceUserRepository(db, dbAkademik *gorm.DB) IOsceUserRepository {
	return &OsceUserRepository{db: db, dbAkademik: dbAkademik}
}

// ===== PENGUJI =====

func (r *OsceUserRepository) FindPenguji(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Table("osce.penguji p").
		Select("p.*, u.username, u.nama_lengkap, u.email").
		Joins("JOIN public.users u ON p.id_user = u.id_user").
		Where("p.status_data = true")
	if filter != "" {
		query = query.Where("u.nama_lengkap ILIKE ? OR u.username ILIKE ?", "%"+filter+"%", "%"+filter+"%")
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

func (r *OsceUserRepository) CountPenguji(ctx context.Context, filter string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Table("osce.penguji p").
		Joins("JOIN public.users u ON p.id_user = u.id_user").
		Where("p.status_data = true")
	if filter != "" {
		query = query.Where("u.nama_lengkap ILIKE ?", "%"+filter+"%")
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *OsceUserRepository) FindPengujiByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.penguji p").
		Select("p.*, u.username, u.nama_lengkap, u.email").
		Joins("JOIN public.users u ON p.id_user = u.id_user").
		Where("p.id_penguji = ? AND p.status_data = true", id).
		First(&data)
	return data, result.Error
}

func (r *OsceUserRepository) FindPengujiByUserID(ctx context.Context, idUser string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.penguji").
		Where("id_user = ? AND status_data = true", idUser).
		First(&data)
	return data, result.Error
}

func (r *OsceUserRepository) InsertPenguji(ctx context.Context, form osce.PengujiForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.penguji (id_user, gelar_depan, gelar_belakang, keahlian, no_str, status_aktif)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		form.IDUser, form.GelarDepan, form.GelarBelakang, form.Keahlian, form.NoSTR, form.StatusAktif,
	).Error
}

func (r *OsceUserRepository) UpdatePenguji(ctx context.Context, id string, form osce.PengujiForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.penguji
		 SET gelar_depan = ?, gelar_belakang = ?, keahlian = ?, no_str = ?, status_aktif = ?, tgl_update = NOW()
		 WHERE id_penguji = ?`,
		form.GelarDepan, form.GelarBelakang, form.Keahlian, form.NoSTR, form.StatusAktif, id,
	).Error
}

func (r *OsceUserRepository) DeletePenguji(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.penguji SET status_data = false, tgl_update = NOW() WHERE id_penguji = ?`, id,
	).Error
}

// ===== KELOMPOK MAHASISWA =====

func (r *OsceUserRepository) FindKelompokMahasiswa(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Table("osce.kelompok_mahasiswa km").
		Select("km.*, u.nama_ujian, (SELECT COUNT(*) FROM osce.anggota_kelompok WHERE id_kelompok_mahasiswa = km.id_kelompok_mahasiswa AND status_data = true) as jumlah_anggota").
		Joins("LEFT JOIN osce.ujian u ON km.id_ujian = u.id_ujian").
		Where("km.status_data = true")
	if filter != "" {
		query = query.Where("km.nama_kelompok ILIKE ?", "%"+filter+"%")
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

func (r *OsceUserRepository) CountKelompokMahasiswa(ctx context.Context, filter string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Table("osce.kelompok_mahasiswa").Where("status_data = true")
	if filter != "" {
		query = query.Where("nama_kelompok ILIKE ?", "%"+filter+"%")
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *OsceUserRepository) FindKelompokMahasiswaByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.kelompok_mahasiswa km").
		Select("km.*, u.nama_ujian, (SELECT COUNT(*) FROM osce.anggota_kelompok WHERE id_kelompok_mahasiswa = km.id_kelompok_mahasiswa AND status_data = true) as jumlah_anggota").
		Joins("LEFT JOIN osce.ujian u ON km.id_ujian = u.id_ujian").
		Where("km.id_kelompok_mahasiswa = ? AND km.status_data = true", id).
		First(&data)
	return data, result.Error
}

func (r *OsceUserRepository) InsertKelompokMahasiswa(ctx context.Context, form osce.KelompokMahasiswaForm) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(
		`INSERT INTO osce.kelompok_mahasiswa (id_ujian, nama_kelompok, deskripsi)
		 VALUES (?, ?, ?) RETURNING id_kelompok_mahasiswa`,
		form.IDUjian, form.NamaKelompok, form.Deskripsi,
	).Scan(&id)
	return id, result.Error
}

func (r *OsceUserRepository) UpdateKelompokMahasiswa(ctx context.Context, id string, form osce.KelompokMahasiswaForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.kelompok_mahasiswa
		 SET id_ujian = ?, nama_kelompok = ?, deskripsi = ?, tgl_update = NOW()
		 WHERE id_kelompok_mahasiswa = ?`,
		form.IDUjian, form.NamaKelompok, form.Deskripsi, id,
	).Error
}

func (r *OsceUserRepository) DeleteKelompokMahasiswa(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.kelompok_mahasiswa SET status_data = false, tgl_update = NOW() WHERE id_kelompok_mahasiswa = ?`, id,
	).Error
}

// ===== ANGGOTA KELOMPOK =====

func (r *OsceUserRepository) FindAnggotaByKelompok(ctx context.Context, idKelompok string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.anggota_kelompok ak").
		Select("ak.*, u.username, u.nama_lengkap, u.email").
		Joins("JOIN public.users u ON ak.id_user = u.id_user").
		Where("ak.id_kelompok_mahasiswa = ? AND ak.status_data = true", idKelompok).
		Order("u.nama_lengkap ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceUserRepository) InsertAnggota(ctx context.Context, idKelompok, idUser string) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.anggota_kelompok (id_kelompok_mahasiswa, id_user)
		 VALUES (?, ?)
		 ON CONFLICT (id_kelompok_mahasiswa, id_user) DO NOTHING`,
		idKelompok, idUser,
	).Error
}

func (r *OsceUserRepository) InsertAnggotaBulk(ctx context.Context, idKelompok string, idUsers []string) (int, int, []string, error) {
	berhasil := 0
	gagal := 0
	var detailGagal []string

	for _, idUser := range idUsers {
		err := r.InsertAnggota(ctx, idKelompok, idUser)
		if err != nil {
			gagal++
			detailGagal = append(detailGagal, idUser+": "+err.Error())
		} else {
			berhasil++
		}
	}

	return berhasil, gagal, detailGagal, nil
}

func (r *OsceUserRepository) DeleteAnggota(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.anggota_kelompok SET status_data = false, tgl_update = NOW() WHERE id_anggota_kelompok = ?`, id,
	).Error
}

func (r *OsceUserRepository) CountAnggotaByKelompok(ctx context.Context, idKelompok string) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Table("osce.anggota_kelompok").
		Where("id_kelompok_mahasiswa = ? AND status_data = true", idKelompok).
		Count(&count)
	return count, result.Error
}

// ===== PENUGASAN PENGUJI =====

func (r *OsceUserRepository) FindPenugasanBySesi(ctx context.Context, idSesi string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.penugasan_penguji pp").
		Select("pp.*, u.nama_lengkap as nama_penguji, s.kode_station, s.nama_station, su.nama_sesi, su.tgl_sesi").
		Joins("JOIN osce.penguji p ON pp.id_penguji = p.id_penguji").
		Joins("JOIN public.users u ON p.id_user = u.id_user").
		Joins("JOIN osce.ujian_station us ON pp.id_ujian_station = us.id_ujian_station").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Joins("JOIN osce.sesi_ujian su ON pp.id_sesi_ujian = su.id_sesi_ujian").
		Where("pp.id_sesi_ujian = ? AND pp.status_data = true", idSesi).
		Order("u.nama_lengkap ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceUserRepository) InsertPenugasan(ctx context.Context, form osce.PenugasanPengujiForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.penugasan_penguji (id_penguji, id_sesi_ujian, id_ujian_station)
		 VALUES (?, ?, ?)
		 ON CONFLICT (id_penguji, id_sesi_ujian, id_ujian_station) DO NOTHING`,
		form.IDPenguji, form.IDSesiUjian, form.IDUjianStation,
	).Error
}

func (r *OsceUserRepository) DeletePenugasan(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.penugasan_penguji SET status_data = false, tgl_update = NOW() WHERE id_penugasan_penguji = ?`, id,
	).Error
}
