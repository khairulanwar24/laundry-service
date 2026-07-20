// Package routes (catalog) mendaftarkan endpoint HTTP untuk domain katalog:
// layanan, varian layanan, parfum, dan diskon.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// CatalogRoute mewadahi controller registry & router Fiber.
type CatalogRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// ICatalogRoute adalah kontrak pendaftaran rute domain katalog.
type ICatalogRoute interface {
	Run()
}

// NewCatalogRoute membuat instance baru.
func NewCatalogRoute(controller controllers.IControllerRegistry, router fiber.Router) ICatalogRoute {
	return &CatalogRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain katalog.
func (r *CatalogRoute) Run() {
	api := r.router.Group("/outlets/:id_outlet", middleware.LaundryJWTMiddleware, middleware.RequireOutletMember)

	// Layanan & varian
	api.Get("/services", r.controller.GetService().ListActive)
	api.Post("/services", middleware.ValidateForm(&dto.ServiceForm{}), r.controller.GetService().Create)
	api.Put("/services/:id_service", middleware.ValidateForm(&dto.ServiceForm{}), r.controller.GetService().Update)
	api.Delete("/services/:id_service", r.controller.GetService().Deactivate)
	api.Post("/services/:id_service/variants", middleware.ValidateForm(&dto.ServiceVariantForm{}), r.controller.GetServiceVariant().Create)
	api.Put("/service-variants/:id_service_variant", middleware.ValidateForm(&dto.ServiceVariantForm{}), r.controller.GetServiceVariant().Update)

	// Parfum
	api.Get("/perfumes", r.controller.GetPerfume().ListActive)
	api.Post("/perfumes", middleware.ValidateForm(&dto.PerfumeForm{}), r.controller.GetPerfume().Create)
	api.Put("/perfumes/:id_perfume", middleware.ValidateForm(&dto.PerfumeForm{}), r.controller.GetPerfume().Update)

	// Diskon
	api.Get("/discounts", r.controller.GetDiscount().ListActive)
	api.Post("/discounts", middleware.ValidateForm(&dto.DiscountForm{}), r.controller.GetDiscount().Create)
	api.Put("/discounts/:id_discount", middleware.ValidateForm(&dto.DiscountForm{}), r.controller.GetDiscount().Update)
}
