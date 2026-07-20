package repositories

import (
	"context"

	"gorm.io/gorm"
)

// DiscountRepository menangani akses data ke laundry.discounts.
type DiscountRepository struct {
	db *gorm.DB
}

// IDiscountRepository adalah kontrak akses data domain diskon.
type IDiscountRepository interface {
	ListActive(ctx context.Context, idOutlet string) ([]map[string]interface{}, error)
	Create(ctx context.Context, idOutlet, nama, jenis string, nilai float64, catatan string, isActive bool) (string, error)
	FindByID(ctx context.Context, idOutlet, idDiscount string) (map[string]interface{}, error)
	Update(ctx context.Context, idDiscount, nama, jenis string, nilai float64, catatan string, isActive bool) error
}

// NewDiscountRepository membuat instance DiscountRepository baru.
func NewDiscountRepository(db *gorm.DB) IDiscountRepository {
	return &DiscountRepository{db: db}
}

func (r *DiscountRepository) ListActive(ctx context.Context, idOutlet string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_discount, nama, jenis, nilai, catatan, is_active
		FROM laundry.discounts
		WHERE id_outlet = ? AND is_active = true AND status_data = true
		ORDER BY nama
	`, idOutlet).Scan(&data)
	return data, result.Error
}

func (r *DiscountRepository) Create(ctx context.Context, idOutlet, nama, jenis string, nilai float64, catatan string, isActive bool) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.discounts (id_outlet, nama, jenis, nilai, catatan, is_active)
		VALUES (?, ?, ?, ?, NULLIF(?, ''), ?)
		RETURNING id_discount
	`, idOutlet, nama, jenis, nilai, catatan, isActive).Scan(&id)
	return id, result.Error
}

func (r *DiscountRepository) FindByID(ctx context.Context, idOutlet, idDiscount string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_discount, id_outlet, nama, jenis, nilai, catatan, is_active
		FROM laundry.discounts
		WHERE id_outlet = ? AND id_discount = ? AND status_data = true
	`, idOutlet, idDiscount).First(&data)
	return data, result.Error
}

func (r *DiscountRepository) Update(ctx context.Context, idDiscount, nama, jenis string, nilai float64, catatan string, isActive bool) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.discounts SET nama = ?, jenis = ?, nilai = ?, catatan = NULLIF(?, ''), is_active = ?, tgl_update = NOW()
		WHERE id_discount = ?
	`, nama, jenis, nilai, catatan, isActive, idDiscount).Error
}
