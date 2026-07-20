// Package repositories (expense) adalah lapisan akses data untuk domain pengeluaran outlet.
package repositories

import (
	"context"

	"gorm.io/gorm"
)

// ExpenseRepository menangani akses data ke laundry.expenses.
type ExpenseRepository struct {
	db *gorm.DB
}

// IExpenseRepository adalah kontrak akses data domain pengeluaran.
type IExpenseRepository interface {
	List(ctx context.Context, idOutlet, kategori, startDate, endDate string, limit, offset int) ([]map[string]interface{}, error)
	Count(ctx context.Context, idOutlet, kategori, startDate, endDate string) (int64, error)
	Create(ctx context.Context, idOutlet, kategori, deskripsi, tanggal string, jumlah float64) (string, error)
	FindByID(ctx context.Context, idOutlet, idExpense string) (map[string]interface{}, error)
	Update(ctx context.Context, idExpense, kategori, deskripsi, tanggal string, jumlah float64) error
	Delete(ctx context.Context, idExpense string) error
}

// NewExpenseRepository membuat instance ExpenseRepository baru.
func NewExpenseRepository(db *gorm.DB) IExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) filterQuery(idOutlet, kategori, startDate, endDate string) *gorm.DB {
	query := r.db.Table("laundry.expenses").Where("id_outlet = ? AND status_data = true", idOutlet)
	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	if startDate != "" {
		query = query.Where("tanggal >= ?::date", startDate)
	}
	if endDate != "" {
		query = query.Where("tanggal <= ?::date", endDate)
	}
	return query
}

func (r *ExpenseRepository) List(ctx context.Context, idOutlet, kategori, startDate, endDate string, limit, offset int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	query := r.filterQuery(idOutlet, kategori, startDate, endDate).WithContext(ctx).
		Select("id_expense, kategori, jumlah, deskripsi, tanggal")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	result := query.Order("tanggal DESC").Scan(&data)
	return data, result.Error
}

func (r *ExpenseRepository) Count(ctx context.Context, idOutlet, kategori, startDate, endDate string) (int64, error) {
	var count int64
	result := r.filterQuery(idOutlet, kategori, startDate, endDate).WithContext(ctx).Count(&count)
	return count, result.Error
}

func (r *ExpenseRepository) Create(ctx context.Context, idOutlet, kategori, deskripsi, tanggal string, jumlah float64) (string, error) {
	var id string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.expenses (id_outlet, kategori, jumlah, deskripsi, tanggal)
		VALUES (?, ?, ?, ?, ?::date)
		RETURNING id_expense
	`, idOutlet, kategori, jumlah, deskripsi, tanggal).Scan(&id)
	return id, result.Error
}

func (r *ExpenseRepository) FindByID(ctx context.Context, idOutlet, idExpense string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_expense, id_outlet, kategori, jumlah, deskripsi, tanggal
		FROM laundry.expenses
		WHERE id_outlet = ? AND id_expense = ? AND status_data = true
	`, idOutlet, idExpense).First(&data)
	return data, result.Error
}

func (r *ExpenseRepository) Update(ctx context.Context, idExpense, kategori, deskripsi, tanggal string, jumlah float64) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.expenses
		SET kategori = ?, jumlah = ?, deskripsi = ?, tanggal = ?::date, tgl_update = NOW()
		WHERE id_expense = ?
	`, kategori, jumlah, deskripsi, tanggal, idExpense).Error
}

func (r *ExpenseRepository) Delete(ctx context.Context, idExpense string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.expenses SET status_data = false, tgl_update = NOW() WHERE id_expense = ?
	`, idExpense).Error
}
