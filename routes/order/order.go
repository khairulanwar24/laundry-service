// Package routes (order) mendaftarkan endpoint HTTP untuk domain pesanan.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// OrderRoute mewadahi controller registry & router Fiber.
type OrderRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IOrderRoute adalah kontrak pendaftaran rute domain pesanan.
type IOrderRoute interface {
	Run()
}

// NewOrderRoute membuat instance baru.
func NewOrderRoute(controller controllers.IControllerRegistry, router fiber.Router) IOrderRoute {
	return &OrderRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain pesanan.
func (r *OrderRoute) Run() {
	api := r.router.Group("/outlets/:id_outlet/orders", middleware.LaundryJWTMiddleware, middleware.RequireOutletMember)

	api.Get("/", r.controller.GetOrder().List)
	api.Get("/search", r.controller.GetOrder().Search)
	api.Post("/", middleware.ValidateForm(&dto.CreateOrderForm{}), r.controller.GetOrder().Create)
	api.Get("/:id_order", r.controller.GetOrder().Get)
	api.Get("/:id_order/history", r.controller.GetOrder().GetHistory)
	api.Post("/:id_order/status", middleware.ValidateForm(&dto.ChangeStatusForm{}), r.controller.GetOrder().ChangeStatus)
	api.Post("/:id_order/pay", middleware.ValidateForm(&dto.PayOrderForm{}), r.controller.GetOrder().Pay)
	api.Post("/:id_order/pickup", middleware.ValidateForm(&dto.PickupForm{}), r.controller.GetOrder().Pickup)
}
