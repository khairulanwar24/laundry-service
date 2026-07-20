// registry.go adalah pusat pendaftaran (registry) seluruh service (logika bisnis) per-domain.
package services

import (
	"laundry-service/repositories"
	accountService "laundry-service/services/account"
	catalogService "laundry-service/services/catalog"
	customerService "laundry-service/services/customer"
	dashboardService "laundry-service/services/dashboard"
	expenseService "laundry-service/services/expense"
	orderService "laundry-service/services/order"
	outletService "laundry-service/services/outlet"
	paymentMethodService "laundry-service/services/paymentmethod"
	reportService "laundry-service/services/report"
)

// Registry menyimpan repository registry sebagai dependency.
type Registry struct {
	repository repositories.IRepositoryRegistry
}

// IServiceRegistry adalah kontrak untuk mengambil service per-domain.
type IServiceRegistry interface {
	GetAccount() accountService.IAccountService
	GetOutlet() outletService.IOutletService
	GetPaymentMethod() paymentMethodService.IPaymentMethodService
	GetService() catalogService.IServiceService
	GetServiceVariant() catalogService.IServiceVariantService
	GetPerfume() catalogService.IPerfumeService
	GetDiscount() catalogService.IDiscountService
	GetCustomer() customerService.ICustomerService
	GetOrder() orderService.IOrderService
	GetDashboard() dashboardService.IDashboardService
	GetExpense() expenseService.IExpenseService
	GetReport() reportService.IReportService
}

// NewServiceRegistry membuat service registry baru.
func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

// GetAccount mengembalikan service akun.
func (r *Registry) GetAccount() accountService.IAccountService {
	return accountService.NewAccountService(r.repository)
}

// GetOutlet mengembalikan service outlet.
func (r *Registry) GetOutlet() outletService.IOutletService {
	return outletService.NewOutletService(r.repository)
}

// GetPaymentMethod mengembalikan service metode pembayaran.
func (r *Registry) GetPaymentMethod() paymentMethodService.IPaymentMethodService {
	return paymentMethodService.NewPaymentMethodService(r.repository)
}

// GetService mengembalikan service layanan.
func (r *Registry) GetService() catalogService.IServiceService {
	return catalogService.NewServiceService(r.repository)
}

// GetServiceVariant mengembalikan service varian layanan.
func (r *Registry) GetServiceVariant() catalogService.IServiceVariantService {
	return catalogService.NewServiceVariantService(r.repository)
}

// GetPerfume mengembalikan service parfum.
func (r *Registry) GetPerfume() catalogService.IPerfumeService {
	return catalogService.NewPerfumeService(r.repository)
}

// GetDiscount mengembalikan service diskon.
func (r *Registry) GetDiscount() catalogService.IDiscountService {
	return catalogService.NewDiscountService(r.repository)
}

// GetCustomer mengembalikan service pelanggan.
func (r *Registry) GetCustomer() customerService.ICustomerService {
	return customerService.NewCustomerService(r.repository)
}

// GetOrder mengembalikan service pesanan.
func (r *Registry) GetOrder() orderService.IOrderService {
	return orderService.NewOrderService(r.repository)
}

// GetDashboard mengembalikan service dashboard.
func (r *Registry) GetDashboard() dashboardService.IDashboardService {
	return dashboardService.NewDashboardService(r.repository)
}

// GetExpense mengembalikan service pengeluaran.
func (r *Registry) GetExpense() expenseService.IExpenseService {
	return expenseService.NewExpenseService(r.repository)
}

// GetReport mengembalikan service laporan.
func (r *Registry) GetReport() reportService.IReportService {
	return reportService.NewReportService(r.repository)
}
