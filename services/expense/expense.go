// Package services (expense) berisi logika bisnis domain pengeluaran outlet.
package services

import (
	"context"
	"time"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	"laundry-service/repositories"
)

// ExpenseService membungkus akses ke repository registry.
type ExpenseService struct {
	repository repositories.IRepositoryRegistry
}

// IExpenseService adalah kontrak logika bisnis domain pengeluaran.
type IExpenseService interface {
	List(ctx context.Context, idOutlet, kategori, startDate, endDate string, page, perPage int) response.Response
	Create(ctx context.Context, idOutlet string, form dto.ExpenseForm) response.Response
	Get(ctx context.Context, idOutlet, idExpense string) response.Response
	Update(ctx context.Context, idOutlet, idExpense string, form dto.ExpenseForm) response.Response
	Delete(ctx context.Context, idOutlet, idExpense string) response.Response
}

// NewExpenseService membuat instance ExpenseService baru.
func NewExpenseService(repository repositories.IRepositoryRegistry) IExpenseService {
	return &ExpenseService{repository: repository}
}

func validateTanggalNotFuture(tanggal string) bool {
	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return false
	}
	return !t.After(time.Now())
}

func (s *ExpenseService) List(ctx context.Context, idOutlet, kategori, startDate, endDate string, page, perPage int) response.Response {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	repo := s.repository.GetExpense()

	data, err := repo.List(ctx, idOutlet, kategori, startDate, endDate, perPage, (page-1)*perPage)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil daftar pengeluaran: " + err.Error()}
	}
	total, err := repo.Count(ctx, idOutlet, kategori, startDate, endDate)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal menghitung total pengeluaran: " + err.Error()}
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

func (s *ExpenseService) Create(ctx context.Context, idOutlet string, form dto.ExpenseForm) response.Response {
	if !validateTanggalNotFuture(form.Tanggal) {
		return response.Response{Success: false, Message: "Tanggal tidak boleh di masa depan"}
	}

	id, err := s.repository.GetExpense().Create(ctx, idOutlet, form.Kategori, form.Deskripsi, form.Tanggal, form.Jumlah)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat pengeluaran: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Pengeluaran berhasil dicatat", Data: map[string]interface{}{"id_expense": id}}
}

func (s *ExpenseService) Get(ctx context.Context, idOutlet, idExpense string) response.Response {
	data, err := s.repository.GetExpense().FindByID(ctx, idOutlet, idExpense)
	if err != nil {
		return response.Response{Success: false, Message: "Pengeluaran tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *ExpenseService) Update(ctx context.Context, idOutlet, idExpense string, form dto.ExpenseForm) response.Response {
	repo := s.repository.GetExpense()

	if _, err := repo.FindByID(ctx, idOutlet, idExpense); err != nil {
		return response.Response{Success: false, Message: "Pengeluaran tidak ditemukan"}
	}
	if !validateTanggalNotFuture(form.Tanggal) {
		return response.Response{Success: false, Message: "Tanggal tidak boleh di masa depan"}
	}

	if err := repo.Update(ctx, idExpense, form.Kategori, form.Deskripsi, form.Tanggal, form.Jumlah); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui pengeluaran: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Pengeluaran berhasil diperbarui"}
}

func (s *ExpenseService) Delete(ctx context.Context, idOutlet, idExpense string) response.Response {
	repo := s.repository.GetExpense()

	if _, err := repo.FindByID(ctx, idOutlet, idExpense); err != nil {
		return response.Response{Success: false, Message: "Pengeluaran tidak ditemukan"}
	}
	if err := repo.Delete(ctx, idExpense); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus pengeluaran: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Pengeluaran berhasil dihapus"}
}
