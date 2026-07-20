// registry.go adalah pusat pendaftaran (registry) seluruh repository per-domain.
package repositories

import (
	accountRepo "laundry-service/repositories/account"
	catalogRepo "laundry-service/repositories/catalog"
	customerRepo "laundry-service/repositories/customer"
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
	GetService() catalogRepo.IServiceRepository
	GetServiceVariant() catalogRepo.IServiceVariantRepository
	GetPerfume() catalogRepo.IPerfumeRepository
	GetDiscount() catalogRepo.IDiscountRepository
	GetCustomer() customerRepo.ICustomerRepository
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

// GetService mengembalikan repository layanan.
func (r *Registry) GetService() catalogRepo.IServiceRepository {
	return catalogRepo.NewServiceRepository(r.db)
}

// GetServiceVariant mengembalikan repository varian layanan.
func (r *Registry) GetServiceVariant() catalogRepo.IServiceVariantRepository {
	return catalogRepo.NewServiceVariantRepository(r.db)
}

// GetPerfume mengembalikan repository parfum.
func (r *Registry) GetPerfume() catalogRepo.IPerfumeRepository {
	return catalogRepo.NewPerfumeRepository(r.db)
}

// GetDiscount mengembalikan repository diskon.
func (r *Registry) GetDiscount() catalogRepo.IDiscountRepository {
	return catalogRepo.NewDiscountRepository(r.db)
}

// GetCustomer mengembalikan repository pelanggan.
func (r *Registry) GetCustomer() customerRepo.ICustomerRepository {
	return customerRepo.NewCustomerRepository(r.db)
}
