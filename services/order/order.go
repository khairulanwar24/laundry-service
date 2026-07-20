// Package services (order) berisi logika bisnis domain pesanan: pembuatan
// pesanan, perhitungan harga/diskon, mesin status, pembayaran, dan pengambilan.
package services

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	"laundry-service/repositories"
)

var validNextStatus = map[string][]string{
	"ANTRIAN":      {"PROSES", "BATAL"},
	"PROSES":       {"SIAP_DIAMBIL", "BATAL"},
	"SIAP_DIAMBIL": {"SELESAI", "BATAL"},
	"SELESAI":      {},
	"BATAL":        {},
}

func isValidTransition(from, to string) bool {
	for _, s := range validNextStatus[from] {
		if s == to {
			return true
		}
	}
	return false
}

// toFloat64 menormalkan nilai kolom NUMERIC/DECIMAL yang dibaca lewat
// map[string]interface{} — driver Postgres bisa mengembalikannya sebagai
// float64, string, maupun []byte tergantung tipe kolom.
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int64:
		return float64(val)
	case int:
		return float64(val)
	case []byte:
		f, _ := strconv.ParseFloat(string(val), 64)
		return f
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}

func ptr(s string) *string { return &s }

func mapStatusesForTab(tab string) []string {
	switch strings.ToUpper(strings.ReplaceAll(tab, "-", "_")) {
	case "ANTRIAN":
		return []string{"ANTRIAN"}
	case "PROSES":
		return []string{"PROSES"}
	case "SIAP_AMBIL", "SIAP_DIAMBIL":
		return []string{"SIAP_DIAMBIL"}
	case "SELESAI":
		return []string{"SELESAI"}
	case "BATAL":
		return []string{"BATAL"}
	default:
		return nil
	}
}

// OrderService membungkus akses ke repository registry.
type OrderService struct {
	repository repositories.IRepositoryRegistry
}

// IOrderService adalah kontrak logika bisnis domain pesanan.
type IOrderService interface {
	Create(ctx context.Context, idOutlet, idUser string, form dto.CreateOrderForm) response.Response
	List(ctx context.Context, idOutlet, tab, q string, page, perPage int) response.Response
	Search(ctx context.Context, idOutlet, q string, page, perPage int) response.Response
	Get(ctx context.Context, idOutlet, idOrder string) response.Response
	GetHistory(ctx context.Context, idOutlet, idOrder string) response.Response
	ChangeStatus(ctx context.Context, idOutlet, idOrder, idUser string, form dto.ChangeStatusForm) response.Response
	Pay(ctx context.Context, idOutlet, idOrder, idUser string, form dto.PayOrderForm) response.Response
	Pickup(ctx context.Context, idOutlet, idOrder, idUser string, form dto.PickupForm) response.Response
}

// NewOrderService membuat instance OrderService baru.
func NewOrderService(repository repositories.IRepositoryRegistry) IOrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) Create(ctx context.Context, idOutlet, idUser string, form dto.CreateOrderForm) response.Response {
	if _, err := s.repository.GetCustomer().FindByID(ctx, idOutlet, form.IDCustomer); err != nil {
		return response.Response{Success: false, Message: "Pelanggan tidak ditemukan"}
	}

	var idPerfume *string
	if form.IDPerfume != "" {
		if _, err := s.repository.GetPerfume().FindByID(ctx, idOutlet, form.IDPerfume); err != nil {
			return response.Response{Success: false, Message: "Parfum tidak ditemukan"}
		}
		idPerfume = &form.IDPerfume
	}

	var idDiscount *string
	var jenisDiskon string
	var nilaiDiskonRaw float64
	if form.IDDiscount != "" {
		discount, err := s.repository.GetDiscount().FindByID(ctx, idOutlet, form.IDDiscount)
		if err != nil {
			return response.Response{Success: false, Message: "Diskon tidak ditemukan"}
		}
		jenisDiskon, _ = discount["jenis"].(string)
		nilaiDiskonRaw = toFloat64(discount["nilai"])
		idDiscount = &form.IDDiscount
	}

	items := make([]dto.OrderItemInput, 0, len(form.Items))
	var subtotal float64
	var maxDurasiJam int

	for _, item := range form.Items {
		variant, err := s.repository.GetServiceVariant().FindByID(ctx, item.IDServiceVariant)
		if err != nil {
			return response.Response{Success: false, Message: "Varian layanan tidak ditemukan"}
		}
		if outletID, _ := variant["id_outlet"].(string); outletID != idOutlet {
			return response.Response{Success: false, Message: "Varian layanan tidak ditemukan"}
		}
		if isActive, _ := variant["is_active"].(bool); !isActive {
			return response.Response{Success: false, Message: "Varian layanan sudah tidak aktif"}
		}

		satuan, _ := variant["satuan"].(string)
		if satuan == "pcs" && item.Qty != math.Trunc(item.Qty) {
			return response.Response{Success: false, Message: "Qty harus bilangan bulat untuk satuan pcs"}
		}

		hargaPerSatuan := toFloat64(variant["harga_per_satuan"])
		totalHarga := item.Qty * hargaPerSatuan
		subtotal += totalHarga

		durasiJam := int(toFloat64(variant["durasi_pengerjaan_jam"]))
		if durasiJam > maxDurasiJam {
			maxDurasiJam = durasiJam
		}

		items = append(items, dto.OrderItemInput{
			IDServiceVariant:       item.IDServiceVariant,
			Satuan:                 satuan,
			Qty:                    item.Qty,
			HargaPerSatuanSnapshot: hargaPerSatuan,
			TotalHarga:             totalHarga,
			Catatan:                item.Catatan,
		})
	}

	var nilaiDiskon float64
	if idDiscount != nil {
		if jenisDiskon == "percent" {
			nilaiDiskon = subtotal * nilaiDiskonRaw / 100
		} else {
			nilaiDiskon = nilaiDiskonRaw
		}
		if nilaiDiskon > subtotal {
			nilaiDiskon = subtotal
		}
	}
	total := subtotal - nilaiDiskon

	checkinAt := time.Now()
	etaAt := checkinAt.Add(time.Duration(maxDurasiJam) * time.Hour)

	idOrder, invoiceNo, err := s.repository.GetOrder().CreateOrder(
		ctx, idOutlet, form.IDCustomer, idPerfume, idDiscount,
		nilaiDiskon, subtotal, total, form.Catatan, checkinAt, etaAt, idUser, items,
	)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat pesanan: " + err.Error()}
	}

	return response.Response{
		Success: true,
		Message: "Pesanan berhasil dibuat",
		Data: map[string]interface{}{
			"id_order":     idOrder,
			"invoice_no":   invoiceNo,
			"subtotal":     subtotal,
			"nilai_diskon": nilaiDiskon,
			"total":        total,
			"eta_at":       etaAt,
		},
	}
}

func (s *OrderService) list(ctx context.Context, idOutlet string, statuses []string, q string, page, perPage int) response.Response {
	page, perPage = normalizePage(page, perPage)
	repo := s.repository.GetOrder()

	data, err := repo.List(ctx, idOutlet, statuses, q, perPage, (page-1)*perPage)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar pesanan: " + err.Error()}
	}
	total, err := repo.Count(ctx, idOutlet, statuses, q)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menghitung total pesanan: " + err.Error()}
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

func normalizePage(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	return page, perPage
}

func (s *OrderService) List(ctx context.Context, idOutlet, tab, q string, page, perPage int) response.Response {
	return s.list(ctx, idOutlet, mapStatusesForTab(tab), q, page, perPage)
}

func (s *OrderService) Search(ctx context.Context, idOutlet, q string, page, perPage int) response.Response {
	if q == "" {
		return response.Response{Success: false, Message: "Kata kunci pencarian wajib diisi"}
	}
	return s.list(ctx, idOutlet, []string{"ANTRIAN", "PROSES", "SIAP_DIAMBIL"}, q, page, perPage)
}

func (s *OrderService) Get(ctx context.Context, idOutlet, idOrder string) response.Response {
	repo := s.repository.GetOrder()

	order, err := repo.FindByID(ctx, idOutlet, idOrder)
	if err != nil {
		return response.Response{Success: false, Message: "Pesanan tidak ditemukan"}
	}
	items, err := repo.ListItems(ctx, idOrder)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil item pesanan: " + err.Error()}
	}
	order["items"] = items

	return response.Response{Success: true, Message: "sukses", Data: order}
}

func (s *OrderService) GetHistory(ctx context.Context, idOutlet, idOrder string) response.Response {
	repo := s.repository.GetOrder()

	if _, err := repo.FindByID(ctx, idOutlet, idOrder); err != nil {
		return response.Response{Success: false, Message: "Pesanan tidak ditemukan"}
	}
	data, err := repo.ListHistory(ctx, idOrder)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil riwayat status: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OrderService) ChangeStatus(ctx context.Context, idOutlet, idOrder, idUser string, form dto.ChangeStatusForm) response.Response {
	repo := s.repository.GetOrder()

	order, err := repo.FindByID(ctx, idOutlet, idOrder)
	if err != nil {
		return response.Response{Success: false, Message: "Pesanan tidak ditemukan"}
	}
	currentStatus, _ := order["status"].(string)

	if !isValidTransition(currentStatus, form.To) {
		return response.Response{Success: false, Message: "Transisi status tidak valid dari " + currentStatus + " ke " + form.To}
	}

	if err := repo.UpdateStatus(ctx, idOrder, form.To); err != nil {
		return response.Response{Success: false, Message: "Gagal mengubah status pesanan: " + err.Error()}
	}
	if err := repo.InsertStatusHistory(ctx, idOrder, ptr(currentStatus), ptr(form.To), &idUser, form.Catatan); err != nil {
		return response.Response{Success: false, Message: "Gagal mencatat riwayat status: " + err.Error()}
	}

	if form.To == "BATAL" {
		if payment, err := repo.FindPaymentByOrder(ctx, idOrder); err == nil {
			if status, _ := payment["status"].(string); status == "SUCCESS" {
				if err := repo.VoidPayment(ctx, idOrder); err != nil {
					return response.Response{Success: false, Message: "Status dibatalkan tapi gagal membatalkan pembayaran: " + err.Error()}
				}
			}
		}
	}

	return response.Response{Success: true, Message: "Status pesanan berhasil diubah"}
}

func (s *OrderService) Pay(ctx context.Context, idOutlet, idOrder, idUser string, form dto.PayOrderForm) response.Response {
	repo := s.repository.GetOrder()

	order, err := repo.FindByID(ctx, idOutlet, idOrder)
	if err != nil {
		return response.Response{Success: false, Message: "Pesanan tidak ditemukan"}
	}

	if _, err := repo.FindPaymentByOrder(ctx, idOrder); err == nil {
		return response.Response{Success: false, Message: "Pesanan sudah memiliki pembayaran"}
	}

	if _, err := s.repository.GetPaymentMethod().FindByID(ctx, idOutlet, form.IDPaymentMethod); err != nil {
		return response.Response{Success: false, Message: "Metode pembayaran tidak ditemukan"}
	}

	total := toFloat64(order["total"])
	if form.Jumlah < total {
		return response.Response{Success: false, Message: "Jumlah pembayaran kurang dari total pesanan"}
	}

	if _, err := repo.CreatePayment(ctx, idOrder, form.IDPaymentMethod, form.Jumlah, form.NoReferensi); err != nil {
		return response.Response{Success: false, Message: "Gagal mencatat pembayaran: " + err.Error()}
	}
	if err := repo.SetPaymentStatus(ctx, idOrder, "PAID"); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui status pembayaran: " + err.Error()}
	}

	currentStatus, _ := order["status"].(string)
	if currentStatus == "ANTRIAN" {
		if err := repo.UpdateStatus(ctx, idOrder, "PROSES"); err != nil {
			return response.Response{Success: false, Message: "Gagal memperbarui status pesanan: " + err.Error()}
		}
		_ = repo.InsertStatusHistory(ctx, idOrder, ptr(currentStatus), ptr("PROSES"), &idUser, "Pembayaran berhasil diproses")
	}

	return response.Response{Success: true, Message: "Pembayaran berhasil dicatat"}
}

func (s *OrderService) Pickup(ctx context.Context, idOutlet, idOrder, idUser string, form dto.PickupForm) response.Response {
	repo := s.repository.GetOrder()

	order, err := repo.FindByID(ctx, idOutlet, idOrder)
	if err != nil {
		return response.Response{Success: false, Message: "Pesanan tidak ditemukan"}
	}
	currentStatus, _ := order["status"].(string)
	if currentStatus != "SIAP_DIAMBIL" {
		return response.Response{Success: false, Message: "Pesanan belum siap diambil"}
	}

	if err := repo.UpdateStatus(ctx, idOrder, "SELESAI"); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui status pesanan: " + err.Error()}
	}
	if err := repo.InsertStatusHistory(ctx, idOrder, ptr(currentStatus), ptr("SELESAI"), &idUser, form.Catatan); err != nil {
		return response.Response{Success: false, Message: "Gagal mencatat riwayat status: " + err.Error()}
	}
	if err := repo.SetPickup(ctx, idOrder, idUser); err != nil {
		return response.Response{Success: false, Message: "Gagal mencatat pengambilan pesanan: " + err.Error()}
	}

	return response.Response{Success: true, Message: "Pesanan berhasil ditandai sudah diambil"}
}
