package repositories

import (
	"context"

	"gorm.io/gorm"
)

// ServiceVariantRepository menangani akses data ke laundry.service_variants.
type ServiceVariantRepository struct {
	db *gorm.DB
}

// IServiceVariantRepository adalah kontrak akses data domain varian layanan.
type IServiceVariantRepository interface {
	Create(ctx context.Context, idService, nama, satuan string, hargaPerSatuan float64, durasiJam int, gambarPath, catatan string, isActive bool) (string, error)
	FindByID(ctx context.Context, idServiceVariant string) (map[string]interface{}, error)
	Update(ctx context.Context, idServiceVariant, nama, satuan string, hargaPerSatuan float64, durasiJam int, gambarPath, catatan string, isActive bool) error
}

// NewServiceVariantRepository membuat instance ServiceVariantRepository baru.
func NewServiceVariantRepository(db *gorm.DB) IServiceVariantRepository {
	return &ServiceVariantRepository{db: db}
}

func (r *ServiceVariantRepository) Create(ctx context.Context, idService, nama, satuan string, hargaPerSatuan float64, durasiJam int, gambarPath, catatan string, isActive bool) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.service_variants (id_service, nama, satuan, harga_per_satuan, durasi_pengerjaan_jam, gambar_path, catatan, is_active)
		VALUES (?, ?, ?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), ?)
		RETURNING id_service_variant
	`, idService, nama, satuan, hargaPerSatuan, durasiJam, gambarPath, catatan, isActive).Scan(&id)
	return id, result.Error
}

// FindByID mengembalikan varian beserta id_outlet pemiliknya (lewat join service)
// supaya service layer bisa memverifikasi kepemilikan outlet.
func (r *ServiceVariantRepository) FindByID(ctx context.Context, idServiceVariant string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT v.id_service_variant, v.id_service, s.id_outlet, v.nama, v.satuan,
		       v.harga_per_satuan, v.durasi_pengerjaan_jam, v.gambar_path, v.catatan, v.is_active
		FROM laundry.service_variants v
		JOIN laundry.services s ON s.id_service = v.id_service
		WHERE v.id_service_variant = ? AND v.status_data = true
	`, idServiceVariant).First(&data)
	return data, result.Error
}

func (r *ServiceVariantRepository) Update(ctx context.Context, idServiceVariant, nama, satuan string, hargaPerSatuan float64, durasiJam int, gambarPath, catatan string, isActive bool) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.service_variants
		SET nama = ?, satuan = ?, harga_per_satuan = ?, durasi_pengerjaan_jam = ?,
		    gambar_path = NULLIF(?, ''), catatan = NULLIF(?, ''), is_active = ?, tgl_update = NOW()
		WHERE id_service_variant = ?
	`, nama, satuan, hargaPerSatuan, durasiJam, gambarPath, catatan, isActive, idServiceVariant).Error
}
