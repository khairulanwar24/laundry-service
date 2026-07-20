// registry.go adalah pusat pendaftaran (registry) seluruh controller per-domain.
package controllers

import (
	accountController "laundry-service/controllers/account"
	catalogController "laundry-service/controllers/catalog"
	customerController "laundry-service/controllers/customer"
	dashboardController "laundry-service/controllers/dashboard"
	expenseController "laundry-service/controllers/expense"
	orderController "laundry-service/controllers/order"
	outletController "laundry-service/controllers/outlet"
	paymentMethodController "laundry-service/controllers/paymentmethod"
	reportController "laundry-service/controllers/report"
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
	GetOrder() orderController.IOrderController
	GetDashboard() dashboardController.IDashboardController
	GetExpense() expenseController.IExpenseController
	GetReport() reportController.IReportController
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

// GetOrder mengembalikan controller pesanan.
func (r *Registry) GetOrder() orderController.IOrderController {
	return orderController.NewOrderController(r.service)
}

// GetDashboard mengembalikan controller dashboard.
func (r *Registry) GetDashboard() dashboardController.IDashboardController {
	return dashboardController.NewDashboardController(r.service)
}

// GetExpense mengembalikan controller pengeluaran.
func (r *Registry) GetExpense() expenseController.IExpenseController {
	return expenseController.NewExpenseController(r.service)
}

// GetReport mengembalikan controller laporan.
func (r *Registry) GetReport() reportController.IReportController {
	return reportController.NewReportController(r.service)
}
