// Package routes (ref) mendaftarkan endpoint HTTP untuk domain referensi.
package routes

import (
	"sso-service/controllers"
	"sso-service/domain/dto"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// RefRoute mewadahi controller registry & router Fiber.
type RefRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IRefRoute adalah kontrak pendaftaran rute referensi.
type IRefRoute interface {
	Run()
}

// NewRefRoute membuat instance RefRoute baru.
func NewRefRoute(controller controllers.IControllerRegistry, router fiber.Router) IRefRoute {
	return &RefRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain referensi.
func (r *RefRoute) Run() {
	ref := r.router.Group("/ref", middleware.JWTMiddleware)
	ref.Get("/get_prodi", r.controller.GetRef().GetMasterProdi)
	ref.Get("/get_angkatan/:id_prodi", middleware.ValidatedParams(&dto.GetAngkatanParams{}), r.controller.GetRef().GetAngkatan)
}
