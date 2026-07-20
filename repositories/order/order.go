// Package repositories (order) adalah lapisan akses data untuk domain pesanan:
// orders, order_items, order_status_histories, payments.
package repositories

import (
	"context"
	"fmt"
	"time"

	"laundry-service/domain/dto"

	"gorm.io/gorm"
)

// OrderRepository menangani akses data ke laundry.orders dan tabel terkait.
type OrderRepository struct {
	db *gorm.DB
}

// IOrderRepository adalah kontrak akses data domain pesanan.
type IOrderRepository interface {
	CreateOrder(ctx context.Context, idOutlet, idCustomer string, idPerfume, idDiscount *string,
		nilaiDiskon, subtotal, total float64, catatan string, checkinAt, etaAt time.Time,
		idUser string, items []dto.OrderItemInput) (string, string, error)
	FindByID(ctx context.Context, idOutlet, idOrder string) (map[string]interface{}, error)
	ListItems(ctx context.Context, idOrder string) ([]map[string]interface{}, error)
	List(ctx context.Context, idOutlet string, statuses []string, q string, limit, offset int) ([]map[string]interface{}, error)
	Count(ctx context.Context, idOutlet string, statuses []string, q string) (int64, error)
	ListHistory(ctx context.Context, idOrder string) ([]map[string]interface{}, error)
	InsertStatusHistory(ctx context.Context, idOrder string, fromStatus, toStatus *string, idUser *string, catatan string) error
	UpdateStatus(ctx context.Context, idOrder, newStatus string) error
	FindPaymentByOrder(ctx context.Context, idOrder string) (map[string]interface{}, error)
	CreatePayment(ctx context.Context, idOrder, idPaymentMethod string, jumlah float64, refNo string) (string, error)
	VoidPayment(ctx context.Context, idOrder string) error
	SetPaymentStatus(ctx context.Context, idOrder, status string) error
	SetPickup(ctx context.Context, idOrder, idUser string) error
}

// NewOrderRepository membuat instance OrderRepository baru.
func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder membuat invoice_no yang aman dari race condition (advisory lock
// per outlet+tanggal dalam transaksi), lalu menyimpan order + item + histori
// status awal secara atomik.
func (r *OrderRepository) CreateOrder(ctx context.Context, idOutlet, idCustomer string, idPerfume, idDiscount *string,
	nilaiDiskon, subtotal, total float64, catatan string, checkinAt, etaAt time.Time,
	idUser string, items []dto.OrderItemInput) (string, string, error) {

	var idOrder, invoiceNo string
	dateSegment := checkinAt.Format("060102")
	lockKey := fmt.Sprintf("laundry_invoice:%s:%s", idOutlet, dateSegment)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`SELECT pg_advisory_xact_lock(hashtext(?))`, lockKey).Error; err != nil {
			return err
		}

		var seq int
		if err := tx.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(invoice_no FROM 10 FOR 4) AS INT)), 0)
			FROM laundry.orders
			WHERE id_outlet = ? AND invoice_no LIKE ?
		`, idOutlet, "JL-"+dateSegment+"%").Scan(&seq).Error; err != nil {
			return err
		}
		invoiceNo = fmt.Sprintf("JL-%s%04d", dateSegment, seq+1)

		result := tx.Raw(`
			INSERT INTO laundry.orders (
				id_outlet, id_customer, invoice_no, status, status_pembayaran,
				id_perfume, id_discount, nilai_diskon, subtotal, total, catatan,
				checkin_at, eta_at, dibuat_oleh_id_user
			) VALUES (
				?, ?, ?, 'ANTRIAN', 'UNPAID',
				?, ?, ?, ?, ?, NULLIF(?, ''),
				?, ?, ?
			) RETURNING id_order
		`, idOutlet, idCustomer, invoiceNo,
			idPerfume, idDiscount, nilaiDiskon, subtotal, total, catatan,
			checkinAt, etaAt, idUser,
		).Scan(&idOrder)
		if result.Error != nil {
			return result.Error
		}

		for _, item := range items {
			if err := tx.Exec(`
				INSERT INTO laundry.order_items (id_order, id_service_variant, satuan, qty, harga_per_satuan_snapshot, total_harga, catatan)
				VALUES (?, ?, ?, ?, ?, ?, NULLIF(?, ''))
			`, idOrder, item.IDServiceVariant, item.Satuan, item.Qty, item.HargaPerSatuanSnapshot, item.TotalHarga, item.Catatan).Error; err != nil {
				return err
			}
		}

		return tx.Exec(`
			INSERT INTO laundry.order_status_histories (id_order, from_status, to_status, id_user)
			VALUES (?, NULL, 'ANTRIAN', ?)
		`, idOrder, idUser).Error
	})

	return idOrder, invoiceNo, err
}

func (r *OrderRepository) FindByID(ctx context.Context, idOutlet, idOrder string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT o.id_order, o.id_outlet, o.id_customer, c.nama AS nama_pelanggan, c.telepon AS telepon_pelanggan,
		       o.invoice_no, o.status, o.status_pembayaran, o.id_payment_method, o.id_perfume, o.id_discount,
		       o.nilai_diskon, o.subtotal, o.total, o.catatan,
		       o.checkin_at, o.eta_at, o.selesai_at, o.batal_at, o.diambil_at, o.diambil_oleh_id_user,
		       o.dibuat_oleh_id_user
		FROM laundry.orders o
		JOIN laundry.customers c ON c.id_customer = o.id_customer
		WHERE o.id_outlet = ? AND o.id_order = ? AND o.status_data = true
	`, idOutlet, idOrder).First(&data)
	return data, result.Error
}

func (r *OrderRepository) ListItems(ctx context.Context, idOrder string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT oi.id_order_item, oi.id_service_variant, sv.nama AS nama_layanan, oi.satuan,
		       oi.qty, oi.harga_per_satuan_snapshot, oi.total_harga, oi.catatan
		FROM laundry.order_items oi
		JOIN laundry.service_variants sv ON sv.id_service_variant = oi.id_service_variant
		WHERE oi.id_order = ?
		ORDER BY oi.tgl_insert
	`, idOrder).Scan(&data)
	return data, result.Error
}

func (r *OrderRepository) List(ctx context.Context, idOutlet string, statuses []string, q string, limit, offset int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.db.WithContext(ctx).Table("laundry.orders o").
		Select("o.id_order, o.invoice_no, o.status, o.status_pembayaran, o.subtotal, o.total, o.checkin_at, o.eta_at, c.nama AS nama_pelanggan, c.telepon AS telepon_pelanggan").
		Joins("JOIN laundry.customers c ON c.id_customer = o.id_customer").
		Where("o.id_outlet = ? AND o.status_data = true", idOutlet)
	if len(statuses) > 0 {
		query = query.Where("o.status IN (?)", statuses)
	}
	if q != "" {
		like := "%" + q + "%"
		query = query.Where("o.invoice_no ILIKE ? OR c.nama ILIKE ? OR c.telepon ILIKE ?", like, like, like)
	}
	result := query.Order("o.checkin_at DESC").Limit(limit).Offset(offset).Scan(&data)
	return data, result.Error
}

func (r *OrderRepository) Count(ctx context.Context, idOutlet string, statuses []string, q string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Table("laundry.orders o").
		Joins("JOIN laundry.customers c ON c.id_customer = o.id_customer").
		Where("o.id_outlet = ? AND o.status_data = true", idOutlet)
	if len(statuses) > 0 {
		query = query.Where("o.status IN (?)", statuses)
	}
	if q != "" {
		like := "%" + q + "%"
		query = query.Where("o.invoice_no ILIKE ? OR c.nama ILIKE ? OR c.telepon ILIKE ?", like, like, like)
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *OrderRepository) ListHistory(ctx context.Context, idOrder string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT h.id_history, h.from_status, h.to_status, h.id_user, u.nama AS nama_user, h.catatan, h.waktu_perubahan
		FROM laundry.order_status_histories h
		LEFT JOIN laundry.users u ON u.id_user = h.id_user
		WHERE h.id_order = ?
		ORDER BY h.waktu_perubahan DESC
	`, idOrder).Scan(&data)
	return data, result.Error
}

func (r *OrderRepository) InsertStatusHistory(ctx context.Context, idOrder string, fromStatus, toStatus *string, idUser *string, catatan string) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO laundry.order_status_histories (id_order, from_status, to_status, id_user, catatan)
		VALUES (?, ?, ?, ?, NULLIF(?, ''))
	`, idOrder, fromStatus, toStatus, idUser, catatan).Error
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, idOrder, newStatus string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.orders
		SET status = ?,
		    selesai_at = CASE WHEN ? = 'SELESAI' THEN NOW() ELSE selesai_at END,
		    batal_at = CASE WHEN ? = 'BATAL' THEN NOW() ELSE batal_at END,
		    tgl_update = NOW()
		WHERE id_order = ?
	`, newStatus, newStatus, newStatus, idOrder).Error
}

func (r *OrderRepository) FindPaymentByOrder(ctx context.Context, idOrder string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_payment, id_order, id_payment_method, jumlah, dibayar_at, no_referensi, catatan, status
		FROM laundry.payments
		WHERE id_order = ?
	`, idOrder).First(&data)
	return data, result.Error
}

func (r *OrderRepository) CreatePayment(ctx context.Context, idOrder, idPaymentMethod string, jumlah float64, refNo string) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.payments (id_order, id_payment_method, jumlah, dibayar_at, no_referensi, status)
		VALUES (?, ?, ?, NOW(), NULLIF(?, ''), 'SUCCESS')
		RETURNING id_payment
	`, idOrder, idPaymentMethod, jumlah, refNo).Scan(&id)
	return id, result.Error
}

func (r *OrderRepository) VoidPayment(ctx context.Context, idOrder string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.payments SET status = 'VOID', tgl_update = NOW() WHERE id_order = ? AND status = 'SUCCESS'
	`, idOrder).Error
}

func (r *OrderRepository) SetPaymentStatus(ctx context.Context, idOrder, status string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.orders SET status_pembayaran = ?, tgl_update = NOW() WHERE id_order = ?
	`, status, idOrder).Error
}

func (r *OrderRepository) SetPickup(ctx context.Context, idOrder, idUser string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.orders SET diambil_at = NOW(), diambil_oleh_id_user = ?, tgl_update = NOW() WHERE id_order = ?
	`, idUser, idOrder).Error
}
