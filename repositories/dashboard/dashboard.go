// Package repositories (dashboard) adalah lapisan akses data untuk domain dashboard
// (ringkasan operasional & keuangan outlet).
package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// DashboardRepository menangani query agregat untuk kebutuhan dashboard.
type DashboardRepository struct {
	db *gorm.DB
}

// IDashboardRepository adalah kontrak akses data domain dashboard.
type IDashboardRepository interface {
	OrdersSummaryOnDate(ctx context.Context, idOutlet string, start, end time.Time) (map[string]interface{}, error)
	PendapatanRange(ctx context.Context, idOutlet string, start, end time.Time) (float64, error)
	PengeluaranRange(ctx context.Context, idOutlet string, start, end time.Time) (float64, error)
	StatusTransitionCounts(ctx context.Context, idOutlet string, start, end time.Time) (map[string]int64, error)
	RevenuePeriod(ctx context.Context, idOutlet string, start, end time.Time) (float64, int64, error)
	TopServices(ctx context.Context, idOutlet string, start, end time.Time, limit int) ([]map[string]interface{}, error)
	TopCustomers(ctx context.Context, idOutlet string, start, end time.Time, limit int) ([]map[string]interface{}, error)
	ExpensesByCategory(ctx context.Context, idOutlet string, start, end time.Time) ([]map[string]interface{}, error)
	ExpensesSumByCategories(ctx context.Context, idOutlet string, start, end time.Time, categories []string) (float64, error)
	LatestTransactions(ctx context.Context, idOutlet string, start, end time.Time, limit int) ([]map[string]interface{}, error)
	CustomerCounts(ctx context.Context, idOutlet string) (total, aktif, nonAktif int64, err error)
	CustomerList(ctx context.Context, idOutlet, search, status, dateFrom, dateTo string, limit, offset int) ([]map[string]interface{}, error)
	CustomerCount(ctx context.Context, idOutlet, search, status, dateFrom, dateTo string) (int64, error)
}

// NewDashboardRepository membuat instance DashboardRepository baru.
func NewDashboardRepository(db *gorm.DB) IDashboardRepository {
	return &DashboardRepository{db: db}
}

// OrdersSummaryOnDate mengembalikan masuk, harus_selesai, terlambat, item_diambil, omset untuk satu hari.
func (r *DashboardRepository) OrdersSummaryOnDate(ctx context.Context, idOutlet string, start, end time.Time) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) FILTER (WHERE tgl_insert >= ? AND tgl_insert < ? AND status <> 'BATAL') AS masuk,
			COUNT(*) FILTER (WHERE eta_at >= ? AND eta_at < ? AND status NOT IN ('SELESAI','BATAL')) AS harus_selesai,
			COUNT(*) FILTER (WHERE eta_at < NOW() AND status NOT IN ('SELESAI','BATAL')) AS terlambat,
			COUNT(*) FILTER (WHERE diambil_at >= ? AND diambil_at < ?) AS item_diambil,
			COALESCE(SUM(total) FILTER (WHERE tgl_insert >= ? AND tgl_insert < ? AND status <> 'BATAL'), 0) AS omset
		FROM laundry.orders
		WHERE id_outlet = ? AND status_data = true
	`, start, end, start, end, start, end, start, end, idOutlet).First(&data)
	return data, result.Error
}

func (r *DashboardRepository) PendapatanRange(ctx context.Context, idOutlet string, start, end time.Time) (float64, error) {
	var total float64
	result := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(p.jumlah), 0)
		FROM laundry.payments p
		JOIN laundry.orders o ON o.id_order = p.id_order
		WHERE o.id_outlet = ? AND p.status = 'SUCCESS' AND p.dibayar_at >= ? AND p.dibayar_at < ?
	`, idOutlet, start, end).Scan(&total)
	return total, result.Error
}

func (r *DashboardRepository) PengeluaranRange(ctx context.Context, idOutlet string, start, end time.Time) (float64, error) {
	var total float64
	result := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(jumlah), 0)
		FROM laundry.expenses
		WHERE id_outlet = ? AND status_data = true AND tanggal >= ?::date AND tanggal < ?::date
	`, idOutlet, start, end).Scan(&total)
	return total, result.Error
}

func (r *DashboardRepository) StatusTransitionCounts(ctx context.Context, idOutlet string, start, end time.Time) (map[string]int64, error) {
	var rows []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT h.to_status, COUNT(*) AS jumlah
		FROM laundry.order_status_histories h
		JOIN laundry.orders o ON o.id_order = h.id_order
		WHERE o.id_outlet = ? AND h.waktu_perubahan >= ? AND h.waktu_perubahan < ? AND h.to_status IN ('SELESAI','BATAL')
		GROUP BY h.to_status
	`, idOutlet, start, end).Scan(&rows)

	counts := map[string]int64{"SELESAI": 0, "BATAL": 0}
	for _, row := range rows {
		status, _ := row["to_status"].(string)
		switch v := row["jumlah"].(type) {
		case int64:
			counts[status] = v
		case int32:
			counts[status] = int64(v)
		}
	}
	return counts, result.Error
}

func (r *DashboardRepository) RevenuePeriod(ctx context.Context, idOutlet string, start, end time.Time) (float64, int64, error) {
	var row struct {
		Omset          float64
		TotalTransaksi int64
	}
	result := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(total), 0) AS omset, COUNT(*) AS total_transaksi
		FROM laundry.orders
		WHERE id_outlet = ? AND status_data = true AND status <> 'BATAL' AND tgl_insert >= ? AND tgl_insert < ?
	`, idOutlet, start, end).Scan(&row)
	return row.Omset, row.TotalTransaksi, result.Error
}

func (r *DashboardRepository) TopServices(ctx context.Context, idOutlet string, start, end time.Time, limit int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT s.nama AS nama_layanan, SUM(oi.total_harga) AS total_omset, COUNT(DISTINCT o.id_order) AS total_transaksi
		FROM laundry.order_items oi
		JOIN laundry.orders o ON o.id_order = oi.id_order
		JOIN laundry.service_variants sv ON sv.id_service_variant = oi.id_service_variant
		JOIN laundry.services s ON s.id_service = sv.id_service
		WHERE o.id_outlet = ? AND o.status <> 'BATAL' AND o.tgl_insert >= ? AND o.tgl_insert < ?
		GROUP BY s.nama
		ORDER BY total_omset DESC
		LIMIT ?
	`, idOutlet, start, end, limit).Scan(&data)
	return data, result.Error
}

func (r *DashboardRepository) TopCustomers(ctx context.Context, idOutlet string, start, end time.Time, limit int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT c.nama AS nama_pelanggan, SUM(o.total) AS total_omset, COUNT(o.id_order) AS total_transaksi
		FROM laundry.orders o
		JOIN laundry.customers c ON c.id_customer = o.id_customer
		WHERE o.id_outlet = ? AND o.status <> 'BATAL' AND o.tgl_insert >= ? AND o.tgl_insert < ?
		GROUP BY c.nama
		ORDER BY total_omset DESC
		LIMIT ?
	`, idOutlet, start, end, limit).Scan(&data)
	return data, result.Error
}

func (r *DashboardRepository) ExpensesByCategory(ctx context.Context, idOutlet string, start, end time.Time) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT kategori, COALESCE(SUM(jumlah), 0) AS total
		FROM laundry.expenses
		WHERE id_outlet = ? AND status_data = true AND tanggal >= ?::date AND tanggal < ?::date
		GROUP BY kategori
		ORDER BY total DESC
	`, idOutlet, start, end).Scan(&data)
	return data, result.Error
}

func (r *DashboardRepository) ExpensesSumByCategories(ctx context.Context, idOutlet string, start, end time.Time, categories []string) (float64, error) {
	var total float64
	result := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(jumlah), 0)
		FROM laundry.expenses
		WHERE id_outlet = ? AND status_data = true AND tanggal >= ?::date AND tanggal < ?::date AND kategori IN (?)
	`, idOutlet, start, end, categories).Scan(&total)
	return total, result.Error
}

// LatestTransactions menggabungkan pembayaran & pengeluaran terbaru dalam satu rentang waktu.
func (r *DashboardRepository) LatestTransactions(ctx context.Context, idOutlet string, start, end time.Time, limit int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		(
			SELECT s.nama AS nama_service, c.nama AS nama_pelanggan, p.jumlah AS total_harga,
			       o.status_pembayaran, 'pendapatan' AS tipe, p.dibayar_at AS waktu, NULL::text AS description
			FROM laundry.payments p
			JOIN laundry.orders o ON o.id_order = p.id_order
			JOIN laundry.customers c ON c.id_customer = o.id_customer
			LEFT JOIN laundry.order_items oi ON oi.id_order = o.id_order
			LEFT JOIN laundry.service_variants sv ON sv.id_service_variant = oi.id_service_variant
			LEFT JOIN laundry.services s ON s.id_service = sv.id_service
			WHERE o.id_outlet = ? AND p.status = 'SUCCESS' AND p.dibayar_at >= ? AND p.dibayar_at < ?
			GROUP BY s.nama, c.nama, p.jumlah, o.status_pembayaran, p.dibayar_at
		)
		UNION ALL
		(
			SELECT NULL, NULL, jumlah, NULL, 'pengeluaran', tgl_insert, deskripsi
			FROM laundry.expenses
			WHERE id_outlet = ? AND status_data = true AND tanggal >= ?::date AND tanggal < ?::date
		)
		ORDER BY waktu DESC
		LIMIT ?
	`, idOutlet, start, end, idOutlet, start, end, limit).Scan(&data)
	return data, result.Error
}

func (r *DashboardRepository) CustomerCounts(ctx context.Context, idOutlet string) (int64, int64, int64, error) {
	var row struct {
		Total    int64
		Aktif    int64
		NonAktif int64
	}
	result := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) AS total,
		       COUNT(*) FILTER (WHERE is_active = true) AS aktif,
		       COUNT(*) FILTER (WHERE is_active = false) AS non_aktif
		FROM laundry.customers
		WHERE id_outlet = ? AND status_data = true
	`, idOutlet).Scan(&row)
	return row.Total, row.Aktif, row.NonAktif, result.Error
}

func (r *DashboardRepository) customerFilterQuery(idOutlet, search, status, dateFrom, dateTo string) *gorm.DB {
	query := r.db.Table("laundry.customers").Where("id_outlet = ? AND status_data = true", idOutlet)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("nama ILIKE ? OR telepon ILIKE ?", like, like)
	}
	switch status {
	case "active":
		query = query.Where("is_active = true")
	case "inactive":
		query = query.Where("is_active = false")
	}
	if dateFrom != "" {
		query = query.Where("tgl_insert >= ?::date", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("tgl_insert < (?::date + INTERVAL '1 day')", dateTo)
	}
	return query
}

func (r *DashboardRepository) CustomerList(ctx context.Context, idOutlet, search, status, dateFrom, dateTo string, limit, offset int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.customerFilterQuery(idOutlet, search, status, dateFrom, dateTo).WithContext(ctx).
		Select("id_customer AS id, nama, telepon AS nomor_telepon, email, alamat, is_active, tgl_insert AS tanggal_bergabung")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	result := query.Order("nama").Scan(&data)
	return data, result.Error
}

func (r *DashboardRepository) CustomerCount(ctx context.Context, idOutlet, search, status, dateFrom, dateTo string) (int64, error) {
	var count int64
	result := r.customerFilterQuery(idOutlet, search, status, dateFrom, dateTo).WithContext(ctx).Count(&count)
	return count, result.Error
}
