package repositories

import (
	"context"

	"gorm.io/gorm"
)

// PerfumeRepository menangani akses data ke laundry.perfumes.
type PerfumeRepository struct {
	db *gorm.DB
}

// IPerfumeRepository adalah kontrak akses data domain parfum.
type IPerfumeRepository interface {
	ListActive(ctx context.Context, idOutlet string) ([]map[string]interface{}, error)
	Create(ctx context.Context, idOutlet, nama, catatan string, isActive bool) (string, error)
	FindByID(ctx context.Context, idOutlet, idPerfume string) (map[string]interface{}, error)
	Update(ctx context.Context, idPerfume, nama, catatan string, isActive bool) error
}

// NewPerfumeRepository membuat instance PerfumeRepository baru.
func NewPerfumeRepository(db *gorm.DB) IPerfumeRepository {
	return &PerfumeRepository{db: db}
}

func (r *PerfumeRepository) ListActive(ctx context.Context, idOutlet string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_perfume, nama, catatan, is_active
		FROM laundry.perfumes
		WHERE id_outlet = ? AND is_active = true AND status_data = true
		ORDER BY nama
	`, idOutlet).Scan(&data)
	return data, result.Error
}

func (r *PerfumeRepository) Create(ctx context.Context, idOutlet, nama, catatan string, isActive bool) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.perfumes (id_outlet, nama, catatan, is_active)
		VALUES (?, ?, NULLIF(?, ''), ?)
		RETURNING id_perfume
	`, idOutlet, nama, catatan, isActive).Scan(&id)
	return id, result.Error
}

func (r *PerfumeRepository) FindByID(ctx context.Context, idOutlet, idPerfume string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_perfume, id_outlet, nama, catatan, is_active
		FROM laundry.perfumes
		WHERE id_outlet = ? AND id_perfume = ? AND status_data = true
	`, idOutlet, idPerfume).First(&data)
	return data, result.Error
}

func (r *PerfumeRepository) Update(ctx context.Context, idPerfume, nama, catatan string, isActive bool) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.perfumes SET nama = ?, catatan = NULLIF(?, ''), is_active = ?, tgl_update = NOW()
		WHERE id_perfume = ?
	`, nama, catatan, isActive, idPerfume).Error
}
