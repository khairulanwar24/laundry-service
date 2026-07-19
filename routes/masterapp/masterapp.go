// Package routes (masterapp) mendaftarkan endpoint HTTP untuk domain master aplikasi.
package routes

import (
	"sso-service/controllers"
	"sso-service/domain/dto"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// MasterAppRoute mewadahi controller registry & router Fiber.
type MasterAppRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IMasterAppRoute adalah kontrak pendaftaran rute master aplikasi.
type IMasterAppRoute interface {
	Run()
}

// NewMasterAppRoute membuat instance MasterAppRoute baru.
func NewMasterAppRoute(controller controllers.IControllerRegistry, router fiber.Router) IMasterAppRoute {
	return &MasterAppRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain master aplikasi.
func (r *MasterAppRoute) Run() {
	masterapps := r.router.Group("/masterapps", middleware.JWTMiddleware)
	masterapps.Get("/", middleware.ValidatedParams2(&dto.GetData{}), r.controller.GetMasterApp().GetMasterApps)
	masterapps.Get("/:id_master_aplikasi", middleware.ValidatedParams(&dto.GetMasterAppByIdParams{}), r.controller.GetMasterApp().GetMasterAppById)
	masterapps.Post("/", middleware.ValidateForm(&dto.CreateMasterAppForm{}), r.controller.GetMasterApp().CreateMasterApp)
	masterapps.Put("/:id_master_aplikasi", middleware.ValidatedParams(&dto.UpdateMasterAppsParams{}), middleware.ValidateForm(&dto.UpdateMasterAppsForm{}), r.controller.GetMasterApp().UpdateMasterApps)
	masterapps.Delete("/:id_master_aplikasi", middleware.ValidatedParams(&dto.DeleteMasterAppParams{}), r.controller.GetMasterApp().DeleteMasterApp)
}
