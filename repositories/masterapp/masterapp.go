// Package repositories (masterapp) adalah lapisan akses data untuk domain master aplikasi.
package repositories

import (
	"context"
	"strings"
	"time"

	middleware "sso-service/middlewares"

	"gorm.io/gorm"
)

// MasterAppRepository memegang koneksi database utama.
type MasterAppRepository struct {
	db *gorm.DB
}

// IMasterAppRepository adalah kontrak akses data domain master aplikasi.
type IMasterAppRepository interface {
	GetMasterApps(order, filter string, limit, offset int) map[string]interface{}
	GetMasterAppById(ctx context.Context, id string) ([]map[string]interface{}, int64, error)
	Create(ctx context.Context, namaAplikasi, deskripsi string, tglVersion time.Time, url, versiAplikasi, image string) error
	UpdateWithoutImage(ctx context.Context, namaAplikasi, deskripsi, versiAplikasi string, tglVersion time.Time, url, id string) (int64, error)
	UpdateWithImage(ctx context.Context, image, namaAplikasi, deskripsi, versiAplikasi string, tglVersion time.Time, url, id string) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

// NewMasterAppRepository membuat instance MasterAppRepository baru.
func NewMasterAppRepository(db *gorm.DB) IMasterAppRepository {
	return &MasterAppRepository{db: db}
}

// GetMasterApps membangun query datatable & mengembalikan hasil dari helper Datatables.
func (r *MasterAppRepository) GetMasterApps(order, filter string, limit, offset int) map[string]interface{} {
	sRecursive := ``
	sTable := ` SELECT id_master_aplikasi
											, nama_aplikasi
											, deskripsi
											, versi_aplikasi
											, tgl_version
											, url
											, image from
											master_aplikasi where status_data = true`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `and LOWER(nama_aplikasi) LIKE ` + "'" + filter + "'" + ` OR LOWER(deskripsi) LIKE ` + "'" + filter + "'" + ` OR LOWER(versi_aplikasi) LIKE ` + "'" + filter + "'" + ` OR LOWER(url) LIKE ` + "'" + filter + "'"
	}

	return middleware.Datatables(sRecursive, sTable, order, sFilter, limit, offset)
}

// GetMasterAppById mengambil satu master aplikasi (data, rowsAffected, error).
func (r *MasterAppRepository) GetMasterAppById(ctx context.Context, id string) ([]map[string]interface{}, int64, error) {
	var master []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`SELECT id_master_aplikasi
								, nama_aplikasi
								, deskripsi
								, versi_aplikasi
								, tgl_version
								, url
								, image FROM master_aplikasi WHERE id_master_aplikasi = ? AND status_data = true`, id).First(&master)
	return master, result.RowsAffected, result.Error
}

// Create menyimpan master aplikasi baru.
func (r *MasterAppRepository) Create(ctx context.Context, namaAplikasi, deskripsi string, tglVersion time.Time, url, versiAplikasi, image string) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO master_aplikasi
								(nama_aplikasi
								, deskripsi
								, tgl_version
								, url
								, versi_aplikasi
								, image
								, status_data
								)
								VALUES
								(?, ?, ?, ?, ?, ?, ?)`, namaAplikasi, deskripsi, tglVersion, url, versiAplikasi, image, true).Error
}

// UpdateWithoutImage memperbarui master aplikasi tanpa mengubah kolom image.
func (r *MasterAppRepository) UpdateWithoutImage(ctx context.Context, namaAplikasi, deskripsi, versiAplikasi string, tglVersion time.Time, url, id string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE master_aplikasi
								SET nama_aplikasi = ?, deskripsi = ?, versi_aplikasi = ?, tgl_version = ?, url = ?
								WHERE id_master_aplikasi = ?`,
		namaAplikasi, deskripsi, versiAplikasi, tglVersion, url, id)
	return result.RowsAffected, result.Error
}

// UpdateWithImage memperbarui master aplikasi termasuk kolom image.
func (r *MasterAppRepository) UpdateWithImage(ctx context.Context, image, namaAplikasi, deskripsi, versiAplikasi string, tglVersion time.Time, url, id string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE master_aplikasi
								SET image = ?, nama_aplikasi = ?, deskripsi = ?, versi_aplikasi = ?, tgl_version = ?, url = ?
								WHERE id_master_aplikasi = ?`,
		image, namaAplikasi, deskripsi, versiAplikasi, tglVersion, url, id)
	return result.RowsAffected, result.Error
}

// Delete menghapus master aplikasi.
func (r *MasterAppRepository) Delete(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`DELETE FROM master_aplikasi
								WHERE id_master_aplikasi = ?`, id)
	return result.RowsAffected, result.Error
}
