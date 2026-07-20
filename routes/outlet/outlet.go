// Package routes (outlet) mendaftarkan endpoint HTTP untuk domain outlet.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// OutletRoute mewadahi controller registry & router Fiber.
type OutletRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IOutletRoute adalah kontrak pendaftaran rute domain outlet.
type IOutletRoute interface {
	Run()
}

// NewOutletRoute membuat instance baru.
func NewOutletRoute(controller controllers.IControllerRegistry, router fiber.Router) IOutletRoute {
	return &OutletRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain outlet.
func (r *OutletRoute) Run() {
	api := r.router.Group("/outlets", middleware.LaundryJWTMiddleware)

	api.Get("/", r.controller.GetOutlet().ListMyOutlets)
	api.Post("/", middleware.ValidateForm(&dto.CreateOutletForm{}), r.controller.GetOutlet().CreateOutlet)

	api.Get("/:id_outlet", middleware.RequireOutletMember, r.controller.GetOutlet().GetOutlet)
	api.Put("/:id_outlet", middleware.RequireOutletOwner, middleware.ValidateForm(&dto.UpdateOutletForm{}), r.controller.GetOutlet().UpdateOutlet)
	api.Get("/:id_outlet/staff", middleware.RequireOutletMember, r.controller.GetOutlet().ListStaff)
	api.Post("/:id_outlet/invite", middleware.RequireOutletOwner, middleware.ValidateForm(&dto.InviteEmployeeForm{}), r.controller.GetOutlet().InviteEmployee)
	api.Put("/:id_outlet/members/:id_member", middleware.RequireOutletOwner, middleware.ValidateForm(&dto.UpdateMemberForm{}), r.controller.GetOutlet().UpdateMember)
	api.Delete("/:id_outlet/members/:id_member", middleware.RequireOutletOwner, r.controller.GetOutlet().RemoveMember)
}
