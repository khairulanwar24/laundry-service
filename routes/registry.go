// registry.go adalah pusat perakitan rute berbasis pola berlapis (registry).
package routes

import (
	"laundry-service/controllers"
	accountRoute "laundry-service/routes/account"
	catalogRoute "laundry-service/routes/catalog"
	customerRoute "laundry-service/routes/customer"
	dashboardRoute "laundry-service/routes/dashboard"
	expenseRoute "laundry-service/routes/expense"
	orderRoute "laundry-service/routes/order"
	outletRoute "laundry-service/routes/outlet"
	paymentMethodRoute "laundry-service/routes/paymentmethod"
	reportRoute "laundry-service/routes/report"

	"github.com/gofiber/fiber/v2"
)

// Registry menyimpan controller registry & router Fiber.
type Registry struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IRouteRegistry adalah kontrak untuk menyajikan seluruh rute domain laundry.
type IRouteRegistry interface {
	Serve()
}

// NewRouteRegistry membuat route registry baru.
func NewRouteRegistry(controller controllers.IControllerRegistry, router fiber.Router) IRouteRegistry {
	return &Registry{controller: controller, router: router}
}

// Serve mendaftarkan seluruh rute domain laundry.
func (r *Registry) Serve() {
	r.accountRoute().Run()
	r.outletRoute().Run()
	r.paymentMethodRoute().Run()
	r.catalogRoute().Run()
	r.customerRoute().Run()
	r.orderRoute().Run()
	r.dashboardRoute().Run()
	r.expenseRoute().Run()
	r.reportRoute().Run()
}

func (r *Registry) accountRoute() accountRoute.IAccountRoute {
	return accountRoute.NewAccountRoute(r.controller, r.router)
}

func (r *Registry) outletRoute() outletRoute.IOutletRoute {
	return outletRoute.NewOutletRoute(r.controller, r.router)
}

func (r *Registry) paymentMethodRoute() paymentMethodRoute.IPaymentMethodRoute {
	return paymentMethodRoute.NewPaymentMethodRoute(r.controller, r.router)
}

func (r *Registry) catalogRoute() catalogRoute.ICatalogRoute {
	return catalogRoute.NewCatalogRoute(r.controller, r.router)
}

func (r *Registry) customerRoute() customerRoute.ICustomerRoute {
	return customerRoute.NewCustomerRoute(r.controller, r.router)
}

func (r *Registry) orderRoute() orderRoute.IOrderRoute {
	return orderRoute.NewOrderRoute(r.controller, r.router)
}

func (r *Registry) dashboardRoute() dashboardRoute.IDashboardRoute {
	return dashboardRoute.NewDashboardRoute(r.controller, r.router)
}

func (r *Registry) expenseRoute() expenseRoute.IExpenseRoute {
	return expenseRoute.NewExpenseRoute(r.controller, r.router)
}

func (r *Registry) reportRoute() reportRoute.IReportRoute {
	return reportRoute.NewReportRoute(r.controller, r.router)
}
