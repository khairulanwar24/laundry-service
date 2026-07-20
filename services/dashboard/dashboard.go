// Package services (dashboard) berisi logika bisnis domain dashboard: ringkasan
// operasional & keuangan outlet.
package services

import (
	"context"
	"time"

	"laundry-service/common/response"
	"laundry-service/repositories"
)

var fixedCostCategories = []string{"Sewa", "Gaji", "Listrik", "Air", "Internet"}

// DashboardService membungkus akses ke repository registry.
type DashboardService struct {
	repository repositories.IRepositoryRegistry
}

// IDashboardService adalah kontrak logika bisnis domain dashboard.
type IDashboardService interface {
	Summary(ctx context.Context, idOutlet, dateStr string) response.Response
	DailyReport(ctx context.Context, idOutlet, dateStr string) response.Response
	FinancialSummary(ctx context.Context, idOutlet, dateStr string) response.Response
	RevenueSummary(ctx context.Context, idOutlet string) response.Response
	ProfitLossSummary(ctx context.Context, idOutlet string) response.Response
	CustomerReport(ctx context.Context, idOutlet, search, status, dateFrom, dateTo string) response.Response
	CustomerSearch(ctx context.Context, idOutlet, search, status, dateFrom, dateTo string, page, perPage int) response.Response
}

// NewDashboardService membuat instance DashboardService baru.
func NewDashboardService(repository repositories.IRepositoryRegistry) IDashboardService {
	return &DashboardService{repository: repository}
}

func parseDateOrToday(dateStr string) time.Time {
	if dateStr == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	t, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	return t
}

func calculateGrowth(current, previous float64) map[string]interface{} {
	nominal := current - previous
	var persentase float64
	if previous == 0 {
		if current > 0 {
			persentase = 100
		}
	} else {
		persentase = (current - previous) / previous * 100
	}
	return map[string]interface{}{"nominal": nominal, "persentase": persentase}
}

func (s *DashboardService) Summary(ctx context.Context, idOutlet, dateStr string) response.Response {
	date := parseDateOrToday(dateStr)
	start := date
	end := start.AddDate(0, 0, 1)
	repo := s.repository.GetDashboard()

	summary, err := repo.OrdersSummaryOnDate(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil ringkasan pesanan: " + err.Error()}
	}
	pendapatan, err := repo.PendapatanRange(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pendapatan: " + err.Error()}
	}
	pengeluaran, err := repo.PengeluaranRange(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pengeluaran: " + err.Error()}
	}

	summary["date"] = date.Format("2006-01-02")
	summary["pendapatan"] = pendapatan
	summary["pengeluaran"] = pengeluaran

	return response.Response{Success: true, Message: "sukses", Data: summary}
}

func (s *DashboardService) DailyReport(ctx context.Context, idOutlet, dateStr string) response.Response {
	date := parseDateOrToday(dateStr)
	start := date
	end := start.AddDate(0, 0, 1)
	repo := s.repository.GetDashboard()

	summary, err := repo.OrdersSummaryOnDate(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil ringkasan pesanan: " + err.Error()}
	}
	transitions, err := repo.StatusTransitionCounts(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil transisi status: " + err.Error()}
	}
	pendapatan, err := repo.PendapatanRange(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pendapatan: " + err.Error()}
	}
	pengeluaran, err := repo.PengeluaranRange(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pengeluaran: " + err.Error()}
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"date":                    date.Format("2006-01-02"),
			"total_transaksi_masuk":   summary["masuk"],
			"total_transaksi_selesai": transitions["SELESAI"],
			"total_transaksi_batal":   transitions["BATAL"],
			"total_pendapatan":        pendapatan,
			"total_omset":             summary["omset"],
			"total_pengeluaran":       pengeluaran,
		},
	}
}

func (s *DashboardService) FinancialSummary(ctx context.Context, idOutlet, dateStr string) response.Response {
	date := parseDateOrToday(dateStr)
	start := date
	end := start.AddDate(0, 0, 1)
	repo := s.repository.GetDashboard()

	pendapatan, err := repo.PendapatanRange(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pendapatan: " + err.Error()}
	}
	pengeluaran, err := repo.PengeluaranRange(ctx, idOutlet, start, end)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pengeluaran: " + err.Error()}
	}
	transaksiTerakhir, err := repo.LatestTransactions(ctx, idOutlet, start, end, 5)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil transaksi terakhir: " + err.Error()}
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"date":               date.Format("2006-01-02"),
			"total_keuangan":     pendapatan - pengeluaran,
			"pendapatan":         pendapatan,
			"pengeluaran":        pengeluaran,
			"transaksi_terakhir": transaksiTerakhir,
		},
	}
}

func (s *DashboardService) RevenueSummary(ctx context.Context, idOutlet string) response.Response {
	repo := s.repository.GetDashboard()
	now := time.Now()

	hariIniStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	hariIniEnd := hariIniStart.AddDate(0, 0, 1)
	kemarinStart := hariIniStart.AddDate(0, 0, -1)

	bulanIniStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	bulanIniEnd := bulanIniStart.AddDate(0, 1, 0)
	bulanLaluStart := bulanIniStart.AddDate(0, -1, 0)

	tahunIniStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	tahunIniEnd := tahunIniStart.AddDate(1, 0, 0)
	tahunLaluStart := tahunIniStart.AddDate(-1, 0, 0)

	omsetHariIni, transHariIni, err := repo.RevenuePeriod(ctx, idOutlet, hariIniStart, hariIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil omset hari ini: " + err.Error()}
	}
	omsetKemarin, _, err := repo.RevenuePeriod(ctx, idOutlet, kemarinStart, hariIniStart)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil omset kemarin: " + err.Error()}
	}
	omsetBulanIni, transBulanIni, err := repo.RevenuePeriod(ctx, idOutlet, bulanIniStart, bulanIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil omset bulan ini: " + err.Error()}
	}
	omsetBulanLalu, _, err := repo.RevenuePeriod(ctx, idOutlet, bulanLaluStart, bulanIniStart)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil omset bulan lalu: " + err.Error()}
	}
	omsetTahunIni, transTahunIni, err := repo.RevenuePeriod(ctx, idOutlet, tahunIniStart, tahunIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil omset tahun ini: " + err.Error()}
	}
	omsetTahunLalu, _, err := repo.RevenuePeriod(ctx, idOutlet, tahunLaluStart, tahunIniStart)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil omset tahun lalu: " + err.Error()}
	}

	topLayanan, err := repo.TopServices(ctx, idOutlet, bulanIniStart, bulanIniEnd, 5)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil top layanan: " + err.Error()}
	}
	topPelanggan, err := repo.TopCustomers(ctx, idOutlet, bulanIniStart, bulanIniEnd, 5)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil top pelanggan: " + err.Error()}
	}

	rataRata := func(omset float64, total int64) float64 {
		if total == 0 {
			return 0
		}
		return omset / float64(total)
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"periode":         map[string]interface{}{"hari_ini": hariIniStart.Format("2006-01-02"), "bulan_ini": bulanIniStart.Format("2006-01"), "tahun_ini": tahunIniStart.Format("2006")},
			"omset":           map[string]interface{}{"hari_ini": omsetHariIni, "bulan_ini": omsetBulanIni, "tahun_ini": omsetTahunIni},
			"total_transaksi": map[string]interface{}{"hari_ini": transHariIni, "bulan_ini": transBulanIni, "tahun_ini": transTahunIni},
			"rata_rata_transaksi": map[string]interface{}{
				"hari_ini":  rataRata(omsetHariIni, transHariIni),
				"bulan_ini": rataRata(omsetBulanIni, transBulanIni),
				"tahun_ini": rataRata(omsetTahunIni, transTahunIni),
			},
			"pertumbuhan": map[string]interface{}{
				"hari_ini":  calculateGrowth(omsetHariIni, omsetKemarin),
				"bulan_ini": calculateGrowth(omsetBulanIni, omsetBulanLalu),
				"tahun_ini": calculateGrowth(omsetTahunIni, omsetTahunLalu),
			},
			"top_layanan":   topLayanan,
			"top_pelanggan": topPelanggan,
		},
	}
}

func (s *DashboardService) ProfitLossSummary(ctx context.Context, idOutlet string) response.Response {
	repo := s.repository.GetDashboard()
	now := time.Now()

	hariIniStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	hariIniEnd := hariIniStart.AddDate(0, 0, 1)
	bulanIniStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	bulanIniEnd := bulanIniStart.AddDate(0, 1, 0)
	tahunIniStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	tahunIniEnd := tahunIniStart.AddDate(1, 0, 0)
	tigaBulanLaluStart := bulanIniStart.AddDate(0, -3, 0)
	tigaBulanLaluEnd := tigaBulanLaluStart.AddDate(0, 1, 0)

	pendapatanHariIni, err := repo.PendapatanRange(ctx, idOutlet, hariIniStart, hariIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pendapatan hari ini: " + err.Error()}
	}
	pengeluaranHariIni, err := repo.PengeluaranRange(ctx, idOutlet, hariIniStart, hariIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pengeluaran hari ini: " + err.Error()}
	}
	pendapatanBulanIni, err := repo.PendapatanRange(ctx, idOutlet, bulanIniStart, bulanIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pendapatan bulan ini: " + err.Error()}
	}
	pengeluaranBulanIni, err := repo.PengeluaranRange(ctx, idOutlet, bulanIniStart, bulanIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pengeluaran bulan ini: " + err.Error()}
	}
	pendapatanTahunIni, err := repo.PendapatanRange(ctx, idOutlet, tahunIniStart, tahunIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pendapatan tahun ini: " + err.Error()}
	}
	pengeluaranTahunIni, err := repo.PengeluaranRange(ctx, idOutlet, tahunIniStart, tahunIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pengeluaran tahun ini: " + err.Error()}
	}
	pengeluaranTigaBulanLalu, err := repo.PengeluaranRange(ctx, idOutlet, tigaBulanLaluStart, tigaBulanLaluEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil pengeluaran 3 bulan lalu: " + err.Error()}
	}

	labaHariIni := pendapatanHariIni - pengeluaranHariIni
	labaBulanIni := pendapatanBulanIni - pengeluaranBulanIni
	labaTahunIni := pendapatanTahunIni - pengeluaranTahunIni

	marginPersen := func(laba, omset float64) float64 {
		if omset == 0 {
			return 0
		}
		return laba / omset * 100
	}
	modalAwal := pengeluaranTigaBulanLalu
	roiPersen := func(laba float64) float64 {
		if modalAwal == 0 {
			return 0
		}
		return laba / modalAwal * 100
	}

	rincianPengeluaranBulanIni, err := repo.ExpensesByCategory(ctx, idOutlet, bulanIniStart, bulanIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil rincian pengeluaran: " + err.Error()}
	}
	bep, err := repo.ExpensesSumByCategories(ctx, idOutlet, bulanIniStart, bulanIniEnd, fixedCostCategories)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menghitung break event point: " + err.Error()}
	}
	_, totalTransaksiBulanIni, err := repo.RevenuePeriod(ctx, idOutlet, bulanIniStart, bulanIniEnd)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil total transaksi bulan ini: " + err.Error()}
	}
	costPerTransaksi := 0.0
	if totalTransaksiBulanIni > 0 {
		costPerTransaksi = pengeluaranBulanIni / float64(totalTransaksiBulanIni)
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"laba": map[string]interface{}{"hari_ini": labaHariIni, "bulan_ini": labaBulanIni, "tahun_ini": labaTahunIni},
			"margin_persen": map[string]interface{}{
				"hari_ini": marginPersen(labaHariIni, pendapatanHariIni), "bulan_ini": marginPersen(labaBulanIni, pendapatanBulanIni), "tahun_ini": marginPersen(labaTahunIni, pendapatanTahunIni),
			},
			"roi_persen": map[string]interface{}{
				"hari_ini": roiPersen(labaHariIni), "bulan_ini": roiPersen(labaBulanIni), "tahun_ini": roiPersen(labaTahunIni),
			},
			"ringkasan_keuangan": map[string]interface{}{
				"hari_ini":  map[string]interface{}{"pendapatan": pendapatanHariIni, "pengeluaran": pengeluaranHariIni, "laba_bersih": labaHariIni},
				"bulan_ini": map[string]interface{}{"pendapatan": pendapatanBulanIni, "pengeluaran": pengeluaranBulanIni, "laba_bersih": labaBulanIni},
				"tahun_ini": map[string]interface{}{"pendapatan": pendapatanTahunIni, "pengeluaran": pengeluaranTahunIni, "laba_bersih": labaTahunIni},
			},
			"rincian_pengeluaran": map[string]interface{}{"bulan_ini": rincianPengeluaranBulanIni},
			"perbandingan_bulanan": map[string]interface{}{
				"bulan_ini":       pengeluaranBulanIni,
				"tiga_bulan_lalu": pengeluaranTigaBulanLalu,
			},
			"metrik_utama": map[string]interface{}{
				"break_event_point_nominal":    bep,
				"cost_per_transaction_nominal": costPerTransaksi,
			},
		},
	}
}

func (s *DashboardService) CustomerReport(ctx context.Context, idOutlet, search, status, dateFrom, dateTo string) response.Response {
	repo := s.repository.GetDashboard()

	total, aktif, nonAktif, err := repo.CustomerCounts(ctx, idOutlet)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menghitung total pelanggan: " + err.Error()}
	}
	list, err := repo.CustomerList(ctx, idOutlet, search, status, dateFrom, dateTo, 0, 0)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar pelanggan: " + err.Error()}
	}
	for _, row := range list {
		if isActive, _ := row["is_active"].(bool); isActive {
			row["status"] = "aktif"
		} else {
			row["status"] = "non-aktif"
		}
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"total_pelanggan":  map[string]interface{}{"semua": total, "aktif": aktif, "non_aktif": nonAktif},
			"filter":           map[string]interface{}{"search": search, "status": status, "date_from": dateFrom, "date_to": dateTo},
			"daftar_pelanggan": list,
		},
	}
}

func (s *DashboardService) CustomerSearch(ctx context.Context, idOutlet, search, status, dateFrom, dateTo string, page, perPage int) response.Response {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	repo := s.repository.GetDashboard()

	total, aktif, nonAktif, err := repo.CustomerCounts(ctx, idOutlet)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menghitung total pelanggan: " + err.Error()}
	}
	data, err := repo.CustomerList(ctx, idOutlet, search, status, dateFrom, dateTo, perPage, (page-1)*perPage)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar pelanggan: " + err.Error()}
	}
	filteredTotal, err := repo.CustomerCount(ctx, idOutlet, search, status, dateFrom, dateTo)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menghitung total pelanggan terfilter: " + err.Error()}
	}

	return response.Response{
		Success: true,
		Message: "sukses",
		Data: map[string]interface{}{
			"total_pelanggan": map[string]interface{}{"semua": total, "aktif": aktif, "non_aktif": nonAktif},
			"data":            data,
			"meta":            map[string]interface{}{"current_page": page, "per_page": perPage, "total": filteredTotal},
			"filter":          map[string]interface{}{"search": search, "status": status, "date_from": dateFrom, "date_to": dateTo},
		},
	}
}
