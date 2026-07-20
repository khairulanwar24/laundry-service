// registry.go adalah pusat pendaftaran (registry) seluruh controller per-domain.
package controllers

import (
	accountController "laundry-service/controllers/account"
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
