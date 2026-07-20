// registry.go adalah pusat pendaftaran (registry) seluruh controller per-domain.
package controllers

import (
	accountController "laundry-service/controllers/account"
	catalogController "laundry-service/controllers/catalog"
	customerController "laundry-service/controllers/customer"
	outletController "laundry-service/controllers/outlet"
	paymentMethodController "laundry-service/controllers/paymentmethod"
	"laundry-service/services"
)

// Registry menyimpan service registry sebagai dependency.
type Registry struct {
	service services.IServiceRegistry
}

// IControllerRegistry adalah kontrak untuk mengambil controller per-domain.
type IControllerRegistry interface {
	GetAccount() accountController.IAccountController
	GetOutlet() outletController.IOutletController
	GetPaymentMethod() paymentMethodController.IPaymentMethodController
	GetService() catalogController.IServiceController
	GetServiceVariant() catalogController.IServiceVariantController
	GetPerfume() catalogController.IPerfumeController
	GetDiscount() catalogController.IDiscountController
	GetCustomer() customerController.ICustomerController
}

// NewControllerRegistry membuat controller registry baru.
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

// GetAccount mengembalikan controller akun.
func (r *Registry) GetAccount() accountController.IAccountController {
	return accountController.NewAccountController(r.service)
}

// GetOutlet mengembalikan controller outlet.
func (r *Registry) GetOutlet() outletController.IOutletController {
	return outletController.NewOutletController(r.service)
}

// GetPaymentMethod mengembalikan controller metode pembayaran.
func (r *Registry) GetPaymentMethod() paymentMethodController.IPaymentMethodController {
	return paymentMethodController.NewPaymentMethodController(r.service)
}

// GetService mengembalikan controller layanan.
func (r *Registry) GetService() catalogController.IServiceController {
	return catalogController.NewServiceController(r.service)
}

// GetServiceVariant mengembalikan controller varian layanan.
func (r *Registry) GetServiceVariant() catalogController.IServiceVariantController {
	return catalogController.NewServiceVariantController(r.service)
}

// GetPerfume mengembalikan controller parfum.
func (r *Registry) GetPerfume() catalogController.IPerfumeController {
	return catalogController.NewPerfumeController(r.service)
}

// GetDiscount mengembalikan controller diskon.
func (r *Registry) GetDiscount() catalogController.IDiscountController {
	return catalogController.NewDiscountController(r.service)
}

// GetCustomer mengembalikan controller pelanggan.
func (r *Registry) GetCustomer() customerController.ICustomerController {
	return customerController.NewCustomerController(r.service)
}
