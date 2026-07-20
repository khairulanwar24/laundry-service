// Package services (report) berisi logika bisnis domain laporan transaksi outlet.
package services

import (
	"context"

	"laundry-service/common/response"
	"laundry-service/repositories"
)

// ReportService membungkus akses ke repository registry.
type ReportService struct {
	repository repositories.IRepositoryRegistry
}

// IReportService adalah kontrak logika bisnis domain laporan.
type IReportService interface {
	Transactions(ctx context.Context, idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir string, page, perPage int) response.Response
}

// NewReportService membuat instance ReportService baru.
func NewReportService(repository repositories.IRepositoryRegistry) IReportService {
	return &ReportService{repository: repository}
}

func (s *ReportService) Transactions(ctx context.Context, idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir string, page, perPage int) response.Response {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	repo := s.repository.GetReport()

	data, err := repo.Transactions(ctx, idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir, perPage, (page-1)*perPage)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil laporan transaksi: " + err.Error()}
	}
	total, err := repo.CountTransactions(ctx, idOutlet, statusPesanan, statusPembayaran, tanggalMulai, tanggalAkhir)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menghitung total transaksi: " + err.Error()}
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{"current_page": page, "per_page": perPage, "total": total},
		},
	}
}
