// Package repositories (report) adalah lapisan akses data untuk domain laporan transaksi.
package repositories

import (
	"context"

	"gorm.io/gorm"
)

// ReportRepository menangani query laporan transaksi outlet.
type ReportRepository struct {
	db *gorm.DB
}

// IReportRepository adalah kontrak akses data domain laporan.
type IReportRepository interface {
	Transactions(ctx context.Context, idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir string, limit, offset int) ([]map[string]interface{}, error)
	CountTransactions(ctx context.Context, idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir string) (int64, error)
}

// NewReportRepository membuat instance ReportRepository baru.
func NewReportRepository(db *gorm.DB) IReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) filterQuery(idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir string) *gorm.DB {
	query := r.db.Table("laundry.orders o").
		Joins("JOIN laundry.customers c ON c.id_customer = o.id_customer").
		Where("o.id_outlet = ? AND o.status_data = true", idOutlet)
	if statusPesanan != "" {
		query = query.Where("o.status = ?", statusPesanan)
	}
	if statusPembayaran != "" {
		query = query.Where("o.status_pembayaran = ?", statusPembayaran)
	}
	if tanggalMulai != "" {
		query = query.Where("o.checkin_at >= ?::date", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("o.checkin_at < (?::date + INTERVAL '1 day')", tanggalAkhir)
	}
	return query
}

func (r *ReportRepository) Transactions(ctx context.Context, idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir string, limit, offset int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.filterQuery(idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir).WithContext(ctx).
		Select("o.id_order, o.invoice_no, o.status, o.status_pembayaran, o.subtotal, o.nilai_diskon, o.total, o.checkin_at, c.nama AS nama_pelanggan, c.telepon AS telepon_pelanggan")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	result := query.Order("o.checkin_at DESC").Scan(&data)
	return data, result.Error
}

func (r *ReportRepository) CountTransactions(ctx context.Context, idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir string) (int64, error) {
	var count int64
	result := r.filterQuery(idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir).WithContext(ctx).Count(&count)
	return count, result.Error
}
