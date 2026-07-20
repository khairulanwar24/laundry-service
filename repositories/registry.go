// registry.go adalah pusat pendaftaran (registry) seluruh repository per-domain.
package repositories

import (
	accountRepo "laundry-service/repositories/account"
	outletRepo "laundry-service/repositories/outlet"
	paymentMethodRepo "laundry-service/repositories/paymentmethod"

	"gorm.io/gorm"
)

// Registry menyimpan koneksi database sebagai dependency.
type Registry struct {
	db *gorm.DB
}

// IRepositoryRegistry adalah kontrak untuk mengambil repository per-domain.
type IRepositoryRegistry interface {
	GetAccount() accountRepo.IAccountRepository
	GetOutlet() outletRepo.IOutletRepository
	GetPaymentMethod() paymentMethodRepo.IPaymentMethodRepository
}

// NewRepositoryRegistry membuat repository registry baru.
func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

// GetAccount mengembalikan repository akun.
func (r *Registry) GetAccount() accountRepo.IAccountRepository {
	return accountRepo.NewAccountRepository(r.db)
}

// GetOutlet mengembalikan repository outlet.
func (r *Registry) GetOutlet() outletRepo.IOutletRepository {
	return outletRepo.NewOutletRepository(r.db)
}

// GetPaymentMethod mengembalikan repository metode pembayaran.
func (r *Registry) GetPaymentMethod() paymentMethodRepo.IPaymentMethodRepository {
	return paymentMethodRepo.NewPaymentMethodRepository(r.db)
}
