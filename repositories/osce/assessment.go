// Package osce adalah lapisan akses data untuk domain OSCE penilaian.
package osce

import (
	"context"
	"sso-service/domain/dto/osce"

	"gorm.io/gorm"
)

// OsceAssessmentRepository memegang koneksi database utama.
type OsceAssessmentRepository struct {
	db *gorm.DB
}

// IOsceAssessmentRepository adalah kontrak akses data domain penilaian OSCE.
type IOsceAssessmentRepository interface {
	// Penilaian
	FindPenilaianByJadwal(ctx context.Context, idJadwal string) (map[string]interface{}, error)
	FindPenilaianByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertPenilaian(ctx context.Context, form osce.PenilaianForm, totalSkor, skorTerbobot float64) (string, error)
	UpdatePenilaian(ctx context.Context, id string, form osce.PenilaianUpdateForm, totalSkor, skorTerbobot float64) error
	KirimPenilaian(ctx context.Context, id string) error

	// Detail Penilaian
	InsertDetailPenilaian(ctx context.Context, idPenilaian, idChecklist string, skor int, catatan string) error
	DeleteDetailPenilaian(ctx context.Context, idPenilaian string) error
	FindDetailPenilaian(ctx context.Context, idPenilaian string) ([]map[string]interface{}, error)

	// Hasil Ujian
	FindHasilByUjian(ctx context.Context, idUjian string) ([]map[string]interface{}, error)
	FindHasilByID(ctx context.Context, id string) (map[string]interface{}, error)
	FindHasilByUjianAndUser(ctx context.Context, idUjian, idUser string) (map[string]interface{}, error)
	UpsertHasilUjian(ctx context.Context, idUjian, idUser string, totalSkor, persentase, batasLulus float64, status string) error
	UpdateStatusHasil(ctx context.Context, idUjian string, batasLulus float64) error

	// Hasil Ujian Station
	UpsertHasilStation(ctx context.Context, idHasil, idUjianStation string, skor, skorTerbobot float64, penilaianGlobal string) error
	FindHasilStation(ctx context.Context, idHasil string) ([]map[string]interface{}, error)

	// Statistik
	FindStatistikUjian(ctx context.Context, idUjian string) (map[string]interface{}, error)

	// Borderline Regression data
	FindPenilaianUntukBorderline(ctx context.Context, idUjian string) ([]map[string]interface{}, error)
}

// NewOsceAssessmentRepository membuat instance baru.
func NewOsceAssessmentRepository(db *gorm.DB) IOsceAssessmentRepository {
	return &OsceAssessmentRepository{db: db}
}

// ===== PENILAIAN =====

func (r *OsceAssessmentRepository) FindPenilaianByJadwal(ctx context.Context, idJadwal string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.penilaian").
		Where("id_jadwal_ujian = ? AND status_data = true", idJadwal).
		Order("tgl_insert DESC").First(&data)
	return data, result.Error
}

func (r *OsceAssessmentRepository) FindPenilaianByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.penilaian p").
		Select("p.*, s.kode_station, s.nama_station, u.nama_lengkap as nama_mahasiswa, u2.nama_lengkap as nama_penguji").
		Joins("JOIN osce.jadwal_ujian ju ON p.id_jadwal_ujian = ju.id_jadwal_ujian").
		Joins("JOIN osce.ujian_station us ON ju.id_ujian_station = us.id_ujian_station").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Joins("JOIN public.users u ON ju.id_user = u.id_user").
		Joins("JOIN osce.penguji pg ON p.id_penguji = pg.id_penguji").
		Joins("JOIN public.users u2 ON pg.id_user = u2.id_user").
		Where("p.id_penilaian = ? AND p.status_data = true", id).
		First(&data)
	return data, result.Error
}

func (r *OsceAssessmentRepository) InsertPenilaian(ctx context.Context, form osce.PenilaianForm, totalSkor, skorTerbobot float64) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(
		`INSERT INTO osce.penilaian
		 (id_jadwal_ujian, id_penguji, total_skor, skor_terbobot, penilaian_global, catatan_penguji, tgl_penilaian)
		 VALUES (?, ?, ?, ?, ?, ?, NOW()) RETURNING id_penilaian`,
		form.IDJadwalUjian, form.IDPenguji, totalSkor, skorTerbobot, form.PenilaianGlobal, form.CatatanPenguji,
	).Scan(&id)
	return id, result.Error
}

func (r *OsceAssessmentRepository) UpdatePenilaian(ctx context.Context, id string, form osce.PenilaianUpdateForm, totalSkor, skorTerbobot float64) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.penilaian
		 SET total_skor = ?, skor_terbobot = ?, penilaian_global = ?, catatan_penguji = ?,
		     tgl_penilaian = NOW(), tgl_update = NOW()
		 WHERE id_penilaian = ? AND status_penilaian = 'draft'`,
		totalSkor, skorTerbobot, form.PenilaianGlobal, form.CatatanPenguji, id,
	).Error
}

func (r *OsceAssessmentRepository) KirimPenilaian(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.penilaian
		 SET status_penilaian = 'terkirim', tgl_penilaian = NOW(), tgl_update = NOW()
		 WHERE id_penilaian = ?`,
		id,
	).Error
}

// ===== DETAIL PENILAIAN =====

func (r *OsceAssessmentRepository) InsertDetailPenilaian(ctx context.Context, idPenilaian, idChecklist string, skor int, catatan string) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.detail_penilaian (id_penilaian, id_checklist_station, skor, catatan_item)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT (id_penilaian, id_checklist_station) DO UPDATE SET skor = ?, catatan_item = ?, tgl_update = NOW()`,
		idPenilaian, idChecklist, skor, catatan, skor, catatan,
	).Error
}

func (r *OsceAssessmentRepository) DeleteDetailPenilaian(ctx context.Context, idPenilaian string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.detail_penilaian SET status_data = false, tgl_update = NOW() WHERE id_penilaian = ?`, idPenilaian,
	).Error
}

func (r *OsceAssessmentRepository) FindDetailPenilaian(ctx context.Context, idPenilaian string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.detail_penilaian dp").
		Select("dp.*, cs.deskripsi_item, cs.skor_maksimal, cs.bobot_item").
		Joins("JOIN osce.checklist_station cs ON dp.id_checklist_station = cs.id_checklist_station").
		Where("dp.id_penilaian = ? AND dp.status_data = true", idPenilaian).
		Order("cs.urutan ASC").
		Find(&data)
	return data, result.Error
}

// ===== HASIL UJIAN =====

func (r *OsceAssessmentRepository) FindHasilByUjian(ctx context.Context, idUjian string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.hasil_ujian hu").
		Select("hu.*, u.nama_lengkap, u.username").
		Joins("JOIN public.users u ON hu.id_user = u.id_user").
		Where("hu.id_ujian = ? AND hu.status_data = true", idUjian).
		Order("u.nama_lengkap ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceAssessmentRepository) FindHasilByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.hasil_ujian hu").
		Select("hu.*, u.nama_ujian, u_lengkap.nama_lengkap, u_lengkap.username").
		Joins("JOIN osce.ujian u ON hu.id_ujian = u.id_ujian").
		Joins("JOIN public.users u_lengkap ON hu.id_user = u_lengkap.id_user").
		Where("hu.id_hasil_ujian = ? AND hu.status_data = true", id).
		First(&data)
	return data, result.Error
}

func (r *OsceAssessmentRepository) FindHasilByUjianAndUser(ctx context.Context, idUjian, idUser string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.hasil_ujian").
		Where("id_ujian = ? AND id_user = ? AND status_data = true", idUjian, idUser).
		First(&data)
	return data, result.Error
}

func (r *OsceAssessmentRepository) UpsertHasilUjian(ctx context.Context, idUjian, idUser string, totalSkor, persentase, batasLulus float64, status string) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.hasil_ujian (id_ujian, id_user, total_skor, persentase_skor, batas_lulus, status_kelulusan)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id_ujian, id_user) DO UPDATE
		 SET total_skor = ?, persentase_skor = ?, batas_lulus = ?, status_kelulusan = ?, tgl_update = NOW()`,
		idUjian, idUser, totalSkor, persentase, batasLulus, status,
		totalSkor, persentase, batasLulus, status,
	).Error
}

func (r *OsceAssessmentRepository) UpdateStatusHasil(ctx context.Context, idUjian string, batasLulus float64) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.hasil_ujian
		 SET batas_lulus = ?,
		     status_kelulusan = CASE
		       WHEN persentase_skor >= ? THEN 'lulus'
		       ELSE 'tidak_lulus'
		     END,
		     tgl_update = NOW()
		 WHERE id_ujian = ? AND status_data = true`,
		batasLulus, batasLulus, idUjian,
	).Error
}

// ===== HASIL UJIAN STATION =====

func (r *OsceAssessmentRepository) UpsertHasilStation(ctx context.Context, idHasil, idUjianStation string, skor, skorTerbobot float64, penilaianGlobal string) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.hasil_ujian_station (id_hasil_ujian, id_ujian_station, skor_station, skor_terbobot, penilaian_global)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT (id_hasil_ujian, id_ujian_station) DO UPDATE
		 SET skor_station = ?, skor_terbobot = ?, penilaian_global = ?, tgl_update = NOW()`,
		idHasil, idUjianStation, skor, skorTerbobot, penilaianGlobal, skor, skorTerbobot, penilaianGlobal,
	).Error
}

func (r *OsceAssessmentRepository) FindHasilStation(ctx context.Context, idHasil string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.hasil_ujian_station hus").
		Select("hus.*, us.urutan, s.kode_station, s.nama_station").
		Joins("JOIN osce.ujian_station us ON hus.id_ujian_station = us.id_ujian_station").
		Joins("JOIN osce.mst_station s ON us.id_station = s.id_station").
		Where("hus.id_hasil_ujian = ? AND hus.status_data = true", idHasil).
		Order("us.urutan ASC").
		Find(&data)
	return data, result.Error
}

// ===== STATISTIK =====

func (r *OsceAssessmentRepository) FindStatistikUjian(ctx context.Context, idUjian string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(
		`SELECT
			COUNT(*) as total_peserta,
			COUNT(CASE WHEN status_kelulusan = 'lulus' THEN 1 END) as jumlah_lulus,
			COUNT(CASE WHEN status_kelulusan = 'tidak_lulus' THEN 1 END) as jumlah_tidak_lulus,
			COUNT(CASE WHEN status_kelulusan = 'remedi' THEN 1 END) as jumlah_remedi,
			COALESCE(AVG(total_skor), 0) as rata_rata_skor,
			COALESCE(MAX(total_skor), 0) as skor_tertinggi,
			COALESCE(MIN(total_skor), 0) as skor_terendah
		FROM osce.hasil_ujian
		WHERE id_ujian = ? AND status_data = true`, idUjian,
	).Scan(&data)
	return data, result.Error
}

// ===== BORDERLINE REGRESSION DATA =====

func (r *OsceAssessmentRepository) FindPenilaianUntukBorderline(ctx context.Context, idUjian string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(
		`SELECT
			p.id_penilaian,
			p.id_penguji,
			p.total_skor,
			p.skor_terbobot,
			p.penilaian_global,
			us.id_ujian_station,
			ju.id_user,
			s.kode_station,
			s.nama_station
		FROM osce.penilaian p
		JOIN osce.jadwal_ujian ju ON p.id_jadwal_ujian = ju.id_jadwal_ujian
		JOIN osce.ujian_station us ON ju.id_ujian_station = us.id_ujian_station
		JOIN osce.mst_station s ON us.id_station = s.id_station
		WHERE ju.id_ujian = ? AND p.status_data = true AND p.status_penilaian IN ('terkirim','terverifikasi')
		ORDER BY s.kode_station`, idUjian,
	).Scan(&data)
	return data, result.Error
}
