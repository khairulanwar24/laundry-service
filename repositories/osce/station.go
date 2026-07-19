// Package osce adalah lapisan akses data untuk domain OSCE station.
package osce

import (
	"context"
	"sso-service/domain/dto/osce"

	"gorm.io/gorm"
)

// OsceStationRepository memegang koneksi database utama.
type OsceStationRepository struct {
	db *gorm.DB
}

// IOsceStationRepository adalah kontrak akses data domain station OSCE.
type IOsceStationRepository interface {
	// Tipe Station
	FindTipeStation(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error)
	CountTipeStation(ctx context.Context, filter string) (int64, error)
	FindTipeStationByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertTipeStation(ctx context.Context, form osce.TipeStationForm) error
	UpdateTipeStation(ctx context.Context, id string, form osce.TipeStationForm) error
	DeleteTipeStation(ctx context.Context, id string) error

	// Master Station
	FindStations(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error)
	CountStations(ctx context.Context, filter string) (int64, error)
	FindStationByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertStation(ctx context.Context, form osce.StationForm) (string, error)
	UpdateStation(ctx context.Context, id string, form osce.StationForm) error
	DeleteStation(ctx context.Context, id string) error

	// Kompetensi Station
	FindKompetensiByStation(ctx context.Context, idStation string) ([]map[string]interface{}, error)
	FindKompetensiByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertKompetensi(ctx context.Context, form osce.KompetensiStationForm) error
	UpdateKompetensi(ctx context.Context, id string, form osce.KompetensiStationForm) error
	DeleteKompetensi(ctx context.Context, id string) error

	// Checklist Station
	FindChecklistByStation(ctx context.Context, idStation string) ([]map[string]interface{}, error)
	FindChecklistByID(ctx context.Context, id string) (map[string]interface{}, error)
	InsertChecklist(ctx context.Context, form osce.ChecklistStationForm) error
	UpdateChecklist(ctx context.Context, id string, form osce.ChecklistStationForm) error
	DeleteChecklist(ctx context.Context, id string) error
	ReorderChecklist(ctx context.Context, id string, urutan int) error
}

// NewOsceStationRepository membuat instance baru.
func NewOsceStationRepository(db *gorm.DB) IOsceStationRepository {
	return &OsceStationRepository{db: db}
}

// ===== TIPE STATION =====

func (r *OsceStationRepository) FindTipeStation(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Table("osce.tipe_station").Where("status_data = true")
	if filter != "" {
		query = query.Where("nama_tipe ILIKE ?", "%"+filter+"%")
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

func (r *OsceStationRepository) CountTipeStation(ctx context.Context, filter string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Table("osce.tipe_station").Where("status_data = true")
	if filter != "" {
		query = query.Where("nama_tipe ILIKE ?", "%"+filter+"%")
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *OsceStationRepository) FindTipeStationByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.tipe_station").
		Where("id_tipe_station = ? AND status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceStationRepository) InsertTipeStation(ctx context.Context, form osce.TipeStationForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.tipe_station (nama_tipe, deskripsi) VALUES (?, ?)`,
		form.NamaTipe, form.Deskripsi,
	).Error
}

func (r *OsceStationRepository) UpdateTipeStation(ctx context.Context, id string, form osce.TipeStationForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.tipe_station SET nama_tipe = ?, deskripsi = ?, tgl_update = NOW() WHERE id_tipe_station = ?`,
		form.NamaTipe, form.Deskripsi, id,
	).Error
}

func (r *OsceStationRepository) DeleteTipeStation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.tipe_station SET status_data = false, tgl_update = NOW() WHERE id_tipe_station = ?`, id,
	).Error
}

// ===== MASTER STATION =====

func (r *OsceStationRepository) FindStations(ctx context.Context, limit, offset int, order, filter string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Table("osce.mst_station s").
		Select("s.*, t.nama_tipe").
		Joins("JOIN osce.tipe_station t ON s.id_tipe_station = t.id_tipe_station").
		Where("s.status_data = true")
	if filter != "" {
		query = query.Where("s.nama_station ILIKE ? OR s.kode_station ILIKE ?", "%"+filter+"%", "%"+filter+"%")
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

func (r *OsceStationRepository) CountStations(ctx context.Context, filter string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Table("osce.mst_station").Where("status_data = true")
	if filter != "" {
		query = query.Where("nama_station ILIKE ? OR kode_station ILIKE ?", "%"+filter+"%", "%"+filter+"%")
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *OsceStationRepository) FindStationByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.mst_station s").
		Select("s.*, t.nama_tipe").
		Joins("JOIN osce.tipe_station t ON s.id_tipe_station = t.id_tipe_station").
		Where("s.id_station = ? AND s.status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceStationRepository) InsertStation(ctx context.Context, form osce.StationForm) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(
		`INSERT INTO osce.mst_station (id_tipe_station, kode_station, nama_station, deskripsi, durasi_default, bobot)
		 VALUES (?, ?, ?, ?, ?, ?) RETURNING id_station`,
		form.IDTipeStation, form.KodeStation, form.NamaStation, form.Deskripsi, form.DurasiDefault, form.Bobot,
	).Scan(&id)
	return id, result.Error
}

func (r *OsceStationRepository) UpdateStation(ctx context.Context, id string, form osce.StationForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.mst_station
		 SET id_tipe_station = ?, kode_station = ?, nama_station = ?, deskripsi = ?, durasi_default = ?, bobot = ?, tgl_update = NOW()
		 WHERE id_station = ?`,
		form.IDTipeStation, form.KodeStation, form.NamaStation, form.Deskripsi, form.DurasiDefault, form.Bobot, id,
	).Error
}

func (r *OsceStationRepository) DeleteStation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.mst_station SET status_data = false, tgl_update = NOW() WHERE id_station = ?`, id,
	).Error
}

// ===== KOMPETENSI STATION =====

func (r *OsceStationRepository) FindKompetensiByStation(ctx context.Context, idStation string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.kompetensi_station").
		Where("id_station = ? AND status_data = true", idStation).
		Order("urutan ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceStationRepository) FindKompetensiByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.kompetensi_station").
		Where("id_kompetensi_station = ? AND status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceStationRepository) InsertKompetensi(ctx context.Context, form osce.KompetensiStationForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.kompetensi_station (id_station, nama_kompetensi, deskripsi, bobot, urutan)
		 VALUES (?, ?, ?, ?, ?)`,
		form.IDStation, form.NamaKompetensi, form.Deskripsi, form.Bobot, form.Urutan,
	).Error
}

func (r *OsceStationRepository) UpdateKompetensi(ctx context.Context, id string, form osce.KompetensiStationForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.kompetensi_station
		 SET nama_kompetensi = ?, deskripsi = ?, bobot = ?, urutan = ?, tgl_update = NOW()
		 WHERE id_kompetensi_station = ?`,
		form.NamaKompetensi, form.Deskripsi, form.Bobot, form.Urutan, id,
	).Error
}

func (r *OsceStationRepository) DeleteKompetensi(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.kompetensi_station SET status_data = false, tgl_update = NOW() WHERE id_kompetensi_station = ?`, id,
	).Error
}

// ===== CHECKLIST STATION =====

func (r *OsceStationRepository) FindChecklistByStation(ctx context.Context, idStation string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.checklist_station c").
		Select("c.*, comp.nama_kompetensi").
		Joins("LEFT JOIN osce.kompetensi_station comp ON c.id_kompetensi_station = comp.id_kompetensi_station").
		Where("c.id_station = ? AND c.status_data = true", idStation).
		Order("c.urutan ASC").
		Find(&data)
	return data, result.Error
}

func (r *OsceStationRepository) FindChecklistByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Table("osce.checklist_station c").
		Select("c.*, comp.nama_kompetensi").
		Joins("LEFT JOIN osce.kompetensi_station comp ON c.id_kompetensi_station = comp.id_kompetensi_station").
		Where("c.id_checklist_station = ? AND c.status_data = true", id).First(&data)
	return data, result.Error
}

func (r *OsceStationRepository) InsertChecklist(ctx context.Context, form osce.ChecklistStationForm) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO osce.checklist_station
		 (id_station, id_kompetensi_station, urutan, deskripsi_item, skor_maksimal, bobot_item, tipe_penilaian, petunjuk_penilaian)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		form.IDStation, form.IDKompetensiStation, form.Urutan, form.DeskripsiItem,
		form.SkorMaksimal, form.BobotItem, form.TipePenilaian, form.PetunjukPenilaian,
	).Error
}

func (r *OsceStationRepository) UpdateChecklist(ctx context.Context, id string, form osce.ChecklistStationForm) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.checklist_station
		 SET id_kompetensi_station = ?, urutan = ?, deskripsi_item = ?, skor_maksimal = ?,
		     bobot_item = ?, tipe_penilaian = ?, petunjuk_penilaian = ?, tgl_update = NOW()
		 WHERE id_checklist_station = ?`,
		form.IDKompetensiStation, form.Urutan, form.DeskripsiItem, form.SkorMaksimal,
		form.BobotItem, form.TipePenilaian, form.PetunjukPenilaian, id,
	).Error
}

func (r *OsceStationRepository) DeleteChecklist(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.checklist_station SET status_data = false, tgl_update = NOW() WHERE id_checklist_station = ?`, id,
	).Error
}

func (r *OsceStationRepository) ReorderChecklist(ctx context.Context, id string, urutan int) error {
	return r.db.WithContext(ctx).Exec(
		`UPDATE osce.checklist_station SET urutan = ?, tgl_update = NOW() WHERE id_checklist_station = ?`, urutan, id,
	).Error
}
