// Package repositories (customer) adalah lapisan akses data untuk domain pelanggan.
package repositories

import (
	"context"

	"gorm.io/gorm"
)

// CustomerRepository menangani akses data ke laundry.customers.
type CustomerRepository struct {
	db *gorm.DB
}

// ICustomerRepository adalah kontrak akses data domain pelanggan.
type ICustomerRepository interface {
	List(ctx context.Context, idOutlet, search string, isActive *bool, limit, offset int) ([]map[string]interface{}, error)
	Count(ctx context.Context, idOutlet, search string, isActive *bool) (int64, error)
	Create(ctx context.Context, idOutlet, nama, telepon, email, alamat string, isActive bool) (string, error)
	FindByID(ctx context.Context, idOutlet, idCustomer string) (map[string]interface{}, error)
	Update(ctx context.Context, idCustomer, nama, telepon, email, alamat string, isActive bool) error
	Delete(ctx context.Context, idCustomer string) error
	CountOrders(ctx context.Context, idCustomer string) (int64, error)
	ExistsPhoneInOutlet(ctx context.Context, idOutlet, telepon, excludeID string) (bool, error)
	ExistsEmailInOutlet(ctx context.Context, idOutlet, email, excludeID string) (bool, error)
	ListOrders(ctx context.Context, idCustomer, status, fromDate, toDate string, limit, offset int) ([]map[string]interface{}, error)
}

// NewCustomerRepository membuat instance CustomerRepository baru.
func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) List(ctx context.Context, idOutlet, search string, isActive *bool, limit, offset int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Table("laundry.customers").
		Select("id_customer, nama, telepon, email, alamat, is_active, tgl_insert").
		Where("id_outlet = ? AND status_data = true", idOutlet)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("nama ILIKE ? OR telepon ILIKE ? OR email ILIKE ?", like, like, like)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	result := query.Order("nama").Limit(limit).Offset(offset).Scan(&data)
	return data, result.Error
}

func (r *CustomerRepository) Count(ctx context.Context, idOutlet, search string, isActive *bool) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Table("laundry.customers").
		Where("id_outlet = ? AND status_data = true", idOutlet)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("nama ILIKE ? OR telepon ILIKE ? OR email ILIKE ?", like, like, like)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *CustomerRepository) Create(ctx context.Context, idOutlet, nama, telepon, email, alamat string, isActive bool) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.customers (id_outlet, nama, telepon, email, alamat, is_active)
		VALUES (?, ?, NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), ?)
		RETURNING id_customer
	`, idOutlet, nama, telepon, email, alamat, isActive).Scan(&id)
	return id, result.Error
}

func (r *CustomerRepository) FindByID(ctx context.Context, idOutlet, idCustomer string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_customer, id_outlet, nama, telepon, email, alamat, is_active, tgl_insert
		FROM laundry.customers
		WHERE id_outlet = ? AND id_customer = ? AND status_data = true
	`, idOutlet, idCustomer).First(&data)
	return data, result.Error
}

func (r *CustomerRepository) Update(ctx context.Context, idCustomer, nama, telepon, email, alamat string, isActive bool) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.customers
		SET nama = ?, telepon = NULLIF(?, ''), email = NULLIF(?, ''), alamat = NULLIF(?, ''), is_active = ?, tgl_update = NOW()
		WHERE id_customer = ?
	`, nama, telepon, email, alamat, isActive, idCustomer).Error
}

func (r *CustomerRepository) Delete(ctx context.Context, idCustomer string) error {
	return r.db.WithContext(ctx).Exec(`DELETE FROM laundry.customers WHERE id_customer = ?`, idCustomer).Error
}

func (r *CustomerRepository) CountOrders(ctx context.Context, idCustomer string) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Table("laundry.orders").Where("id_customer = ?", idCustomer).Count(&count)
	return count, result.Error
}

func (r *CustomerRepository) ExistsPhoneInOutlet(ctx context.Context, idOutlet, telepon, excludeID string) (bool, error) {
	if telepon == "" {
		return false, nil
	}
	var count int64
	query := r.db.WithContext(ctx).Table("laundry.customers").
		Where("id_outlet = ? AND telepon = ? AND status_data = true", idOutlet, telepon)
	if excludeID != "" {
		query = query.Where("id_customer <> ?", excludeID)
	}
	result := query.Count(&count)
	return count > 0, result.Error
}

func (r *CustomerRepository) ExistsEmailInOutlet(ctx context.Context, idOutlet, email, excludeID string) (bool, error) {
	if email == "" {
		return false, nil
	}
	var count int64
	query := r.db.WithContext(ctx).Table("laundry.customers").
		Where("id_outlet = ? AND email = ? AND status_data = true", idOutlet, email)
	if excludeID != "" {
		query = query.Where("id_customer <> ?", excludeID)
	}
	result := query.Count(&count)
	return count > 0, result.Error
}

func (r *CustomerRepository) ListOrders(ctx context.Context, idCustomer, status, fromDate, toDate string, limit, offset int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Raw(`
		SELECT o.id_order, o.invoice_no, o.status, o.status_pembayaran, o.subtotal, o.total,
		       o.checkin_at, o.eta_at, o.diambil_at,
		       COALESCE(json_agg(json_build_object(
		           'id_order_item', oi.id_order_item,
		           'nama_layanan', sv.nama,
		           'satuan', oi.satuan,
		           'qty', oi.qty,
		           'harga_per_satuan_snapshot', oi.harga_per_satuan_snapshot,
		           'total_harga', oi.total_harga
		       )) FILTER (WHERE oi.id_order_item IS NOT NULL), '[]') AS items
		FROM laundry.orders o
		LEFT JOIN laundry.order_items oi ON oi.id_order = o.id_order
		LEFT JOIN laundry.service_variants sv ON sv.id_service_variant = oi.id_service_variant
		WHERE o.id_customer = ?
		  AND (? = '' OR o.status = ?)
		  AND (? = '' OR o.tgl_insert >= ?::date)
		  AND (? = '' OR o.tgl_insert < (?::date + INTERVAL '1 day'))
		GROUP BY o.id_order
		ORDER BY o.checkin_at DESC
		LIMIT ? OFFSET ?
	`, idCustomer, status, status, fromDate, fromDate, toDate, toDate, limit, offset).Scan(&data)
	return data, query.Error
}
