// Package routes (mstmenu) mendaftarkan endpoint HTTP untuk domain master menu & modul.
package routes

import (
	"sso-service/controllers"
	"sso-service/domain/dto"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// MstMenuRoute mewadahi controller registry & router Fiber.
type MstMenuRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IMstMenuRoute adalah kontrak pendaftaran rute master menu & modul.
type IMstMenuRoute interface {
	Run()
}

// NewMstMenuRoute membuat instance MstMenuRoute baru.
func NewMstMenuRoute(controller controllers.IControllerRegistry, router fiber.Router) IMstMenuRoute {
	return &MstMenuRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain master menu & modul.
func (r *MstMenuRoute) Run() {
	mstmenu := r.router.Group("/mstmenu", middleware.JWTMiddleware)

	mstmenu.Get("/:id_master_aplikasi",
		middleware.ValidatedParams(&dto.GetMstMenuParams{}),
		middleware.ValidatedParams2(&dto.GetMstMenuForm{}),
		r.controller.GetMstMenu().GetMstMenu)
	mstmenu.Get("/detail/:id_master_menu",
		middleware.ValidatedParams(&dto.GetDetailMstMenuParams{}),
		r.controller.GetMstMenu().GetDetailMstMenu)
	mstmenu.Get("/modul/:id_master_menu",
		middleware.ValidatedParams(&dto.GetMstMenuModulParams{}),
		middleware.ValidatedParams2(&dto.GetMstMenuModulForm{}),
		r.controller.GetMstMenu().GetMstMenuModul)
	mstmenu.Get("/modul/detail/:id_master_modul",
		middleware.ValidatedParams(&dto.GetDetailMstMenuModulParams{}),
		r.controller.GetMstMenu().GetDetailMstMenuModul)
	mstmenu.Post("/",
		middleware.ValidateForm(&dto.CreateMstMenuForm{}),
		r.controller.GetMstMenu().CreateMstMenu)
	mstmenu.Post("/modul",
		middleware.ValidateForm(&dto.CreateMstMenuModulForm{}),
		r.controller.GetMstMenu().CreateMstMenuModul)
	mstmenu.Put("/:id_master_menu",
		middleware.ValidatedParams(&dto.UpdateMstMenuParams{}),
		middleware.ValidateForm(&dto.UpdateMstMenuForm{}),
		r.controller.GetMstMenu().UpdateMstMenu)
	mstmenu.Delete("/:id_master_menu",
		middleware.ValidatedParams(&dto.DeleteMstMenuParams{}),
		r.controller.GetMstMenu().DeleteMstMenu)
	mstmenu.Put("/modul/:id_master_modul",
		middleware.ValidatedParams(&dto.UpdateMstMenuModulParams{}),
		middleware.ValidateForm(&dto.UpdateMstMenuModulForm{}),
		r.controller.GetMstMenu().UpdateMstMenuModul)
	mstmenu.Delete("/modul/:id_master_modul",
		middleware.ValidatedParams(&dto.DeleteMstMenuModulParams{}),
		r.controller.GetMstMenu().DeleteMstMenuModul)
}
