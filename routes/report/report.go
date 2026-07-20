// Package routes (report) mendaftarkan endpoint HTTP untuk domain laporan.
package routes

import (
	"laundry-service/controllers"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// ReportRoute mewadahi controller registry & router Fiber.
type ReportRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IReportRoute adalah kontrak pendaftaran rute domain laporan.
type IReportRoute interface {
	Run()
}

// NewReportRoute membuat instance baru.
func NewReportRoute(controller controllers.IControllerRegistry, router fiber.Router) IReportRoute {
	return &ReportRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain laporan.
func (r *ReportRoute) Run() {
	api := r.router.Group("/outlets/:id_outlet/reports", middleware.LaundryJWTMiddleware, middleware.RequireOutletMember)

	api.Get("/transactions", r.controller.GetReport().Transactions)
}
