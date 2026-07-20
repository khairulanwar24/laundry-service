// Package repositories (catalog) adalah lapisan akses data untuk domain katalog:
// layanan (service), varian layanan, parfum, dan diskon.
package repositories

import (
	"context"

	"gorm.io/gorm"
)

// ServiceRepository menangani akses data ke laundry.services.
type ServiceRepository struct {
	db *gorm.DB
}

// IServiceRepository adalah kontrak akses data domain layanan.
type IServiceRepository interface {
	ListActive(ctx context.Context, idOutlet, q string) ([]map[string]interface{}, error)
	Create(ctx context.Context, idOutlet, nama string, prioritas int, langkahProsesJSON string, isActive bool) (string, error)
	FindByID(ctx context.Context, idOutlet, idService string) (map[string]interface{}, error)
	Update(ctx context.Context, idService, nama string, prioritas int, langkahProsesJSON string, isActive bool) error
	DeactivateVariantsByService(ctx context.Context, idService string) error
	Deactivate(ctx context.Context, idService string) error
}

// NewServiceRepository membuat instance ServiceRepository baru.
func NewServiceRepository(db *gorm.DB) IServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) ListActive(ctx context.Context, idOutlet, q string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Raw(`
		SELECT s.id_service, s.nama, s.prioritas, s.langkah_proses, s.is_active,
		       COALESCE(json_agg(json_build_object(
		           'id_service_variant', v.id_service_variant,
		           'nama', v.nama,
		           'satuan', v.satuan,
		           'harga_per_satuan', v.harga_per_satuan,
		           'durasi_pengerjaan_jam', v.durasi_pengerjaan_jam,
		           'gambar_path', v.gambar_path,
		           'catatan', v.catatan,
		           'is_active', v.is_active
		       )) FILTER (WHERE v.id_service_variant IS NOT NULL), '[]') AS variants
		FROM laundry.services s
		LEFT JOIN laundry.service_variants v ON v.id_service = s.id_service AND v.is_active = true AND v.status_data = true
		WHERE s.id_outlet = ? AND s.is_active = true AND s.status_data = true
		  AND (? = '' OR s.nama ILIKE '%' || ? || '%' OR v.nama ILIKE '%' || ? || '%')
		GROUP BY s.id_service
		ORDER BY s.prioritas DESC, s.nama ASC
	`, idOutlet, q, q, q)
	result := query.Scan(&data)
	return data, result.Error
}

func (r *ServiceRepository) Create(ctx context.Context, idOutlet, nama string, prioritas int, langkahProsesJSON string, isActive bool) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.services (id_outlet, nama, prioritas, langkah_proses, is_active)
		VALUES (?, ?, ?, ?::jsonb, ?)
		RETURNING id_service
	`, idOutlet, nama, prioritas, langkahProsesJSON, isActive).Scan(&id)
	return id, result.Error
}

func (r *ServiceRepository) FindByID(ctx context.Context, idOutlet, idService string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_service, id_outlet, nama, prioritas, langkah_proses, is_active
		FROM laundry.services
		WHERE id_outlet = ? AND id_service = ? AND status_data = true
	`, idOutlet, idService).First(&data)
	return data, result.Error
}

func (r *ServiceRepository) Update(ctx context.Context, idService, nama string, prioritas int, langkahProsesJSON string, isActive bool) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.services
		SET nama = ?, prioritas = ?, langkah_proses = ?::jsonb, is_active = ?, tgl_update = NOW()
		WHERE id_service = ?
	`, nama, prioritas, langkahProsesJSON, isActive, idService).Error
}

func (r *ServiceRepository) DeactivateVariantsByService(ctx context.Context, idService string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.service_variants SET is_active = false, tgl_update = NOW() WHERE id_service = ?
	`, idService).Error
}

func (r *ServiceRepository) Deactivate(ctx context.Context, idService string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.services SET is_active = false, tgl_update = NOW() WHERE id_service = ?
	`, idService).Error
}
