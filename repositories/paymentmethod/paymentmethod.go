// Package repositories (paymentmethod) adalah lapisan akses data untuk domain metode pembayaran.
package repositories

import (
	"context"

	"gorm.io/gorm"
)

// PaymentMethodRepository menangani akses data ke laundry.payment_methods.
type PaymentMethodRepository struct {
	db *gorm.DB
}

// IPaymentMethodRepository adalah kontrak akses data domain metode pembayaran.
type IPaymentMethodRepository interface {
	ListActive(ctx context.Context, idOutlet string) ([]map[string]interface{}, error)
	Create(ctx context.Context, idOutlet, kategori, nama, logo, namaPemilik, tagsJSON string, isActive bool) (string, error)
	FindByID(ctx context.Context, idOutlet, idPaymentMethod string) (map[string]interface{}, error)
	Update(ctx context.Context, idPaymentMethod, kategori, nama, logo, namaPemilik, tagsJSON string, isActive bool) error
	Deactivate(ctx context.Context, idPaymentMethod string) error
}

// NewPaymentMethodRepository membuat instance PaymentMethodRepository baru.
func NewPaymentMethodRepository(db *gorm.DB) IPaymentMethodRepository {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) ListActive(ctx context.Context, idOutlet string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_payment_method, id_outlet, kategori, nama, logo, nama_pemilik, tags, is_active
		FROM laundry.payment_methods
		WHERE id_outlet = ? AND is_active = true AND status_data = true
		ORDER BY nama
	`, idOutlet).Scan(&data)
	return data, result.Error
}

func (r *PaymentMethodRepository) Create(ctx context.Context, idOutlet, kategori, nama, logo, namaPemilik, tagsJSON string, isActive bool) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.payment_methods (id_outlet, kategori, nama, logo, nama_pemilik, tags, is_active)
		VALUES (?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), ?::jsonb, ?)
		RETURNING id_payment_method
	`, idOutlet, kategori, nama, logo, namaPemilik, tagsJSON, isActive).Scan(&id)
	return id, result.Error
}

func (r *PaymentMethodRepository) FindByID(ctx context.Context, idOutlet, idPaymentMethod string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_payment_method, id_outlet, kategori, nama, logo, nama_pemilik, tags, is_active
		FROM laundry.payment_methods
		WHERE id_outlet = ? AND id_payment_method = ? AND status_data = true
	`, idOutlet, idPaymentMethod).First(&data)
	return data, result.Error
}

func (r *PaymentMethodRepository) Update(ctx context.Context, idPaymentMethod, kategori, nama, logo, namaPemilik, tagsJSON string, isActive bool) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.payment_methods
		SET kategori = ?, nama = ?, logo = NULLIF(?, ''), nama_pemilik = NULLIF(?, ''), tags = ?::jsonb, is_active = ?, tgl_update = NOW()
		WHERE id_payment_method = ?
	`, kategori, nama, logo, namaPemilik, tagsJSON, isActive, idPaymentMethod).Error
}

func (r *PaymentMethodRepository) Deactivate(ctx context.Context, idPaymentMethod string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.payment_methods SET is_active = false, tgl_update = NOW() WHERE id_payment_method = ?
	`, idPaymentMethod).Error
}
