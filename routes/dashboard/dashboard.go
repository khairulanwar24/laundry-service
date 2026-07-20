// Package routes (dashboard) mendaftarkan endpoint HTTP untuk domain dashboard.
package routes

import (
	"laundry-service/controllers"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// DashboardRoute mewadahi controller registry & router Fiber.
type DashboardRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IDashboardRoute adalah kontrak pendaftaran rute domain dashboard.
type IDashboardRoute interface {
	Run()
}

// NewDashboardRoute membuat instance baru.
func NewDashboardRoute(controller controllers.IControllerRegistry, router fiber.Router) IDashboardRoute {
	return &DashboardRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain dashboard.
func (r *DashboardRoute) Run() {
	api := r.router.Group("/outlets/:id_outlet/dashboard", middleware.LaundryJWTMiddleware, middleware.RequireOutletMember)

	api.Get("/summary", r.controller.GetDashboard().Summary)
	api.Get("/daily-report", r.controller.GetDashboard().DailyReport)
	api.Get("/financial-summary", r.controller.GetDashboard().FinancialSummary)
	api.Get("/revenue-summary", r.controller.GetDashboard().RevenueSummary)
	api.Get("/profit-loss-summary", r.controller.GetDashboard().ProfitLossSummary)
	api.Get("/customer-report", r.controller.GetDashboard().CustomerReport)
	api.Get("/customer-search", r.controller.GetDashboard().CustomerSearch)
}
