// Package repositories (ref) adalah lapisan akses data untuk domain referensi (prodi, angkatan).
// Query database mentah dipindahkan ke sini dari package `models` yang lama.
package repositories

import (
	"context"

	"gorm.io/gorm"
)

// RefRepository menyimpan koneksi database akademik.
type RefRepository struct {
	db *gorm.DB
}

// IRefRepository adalah kontrak aksi data untuk domain referensi.
type IRefRepository interface {
	GetMasterProdi(ctx context.Context) ([]map[string]any, error)
	GetAngkatan(ctx context.Context, idProdi string) ([]map[string]any, error)
}

// NewRefRepository membuat instance RefRepository baru.
func NewRefRepository(db *gorm.DB) IRefRepository {
	return &RefRepository{db: db}
}

// GetMasterProdi mengambil seluruh data master prodi.
func (r *RefRepository) GetMasterProdi(ctx context.Context) ([]map[string]any, error) {
	var prodi []map[string]any
	err := r.db.WithContext(ctx).Raw(`select * from master_prodi`).Scan(&prodi).Error
	if err != nil {
		return nil, err
	}
	return prodi, nil
}

// GetAngkatan mengambil daftar angkatan mahasiswa aktif berdasarkan id prodi.
func (r *RefRepository) GetAngkatan(ctx context.Context, idProdi string) ([]map[string]any, error) {
	var angkatan []map[string]any
	err := r.db.WithContext(ctx).Raw(
		`SELECT id_angkatan_mahasiswa, nama_angkatan_mahasiswa
		 FROM master_angkatan_mahasiswa
		 WHERE id_prodi = ? AND status_data = true
		 ORDER BY nama_angkatan_mahasiswa DESC`,
		idProdi,
	).Scan(&angkatan).Error
	if err != nil {
		return nil, err
	}
	return angkatan, nil
}
