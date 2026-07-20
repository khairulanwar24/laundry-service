// Package routes (customer) mendaftarkan endpoint HTTP untuk domain pelanggan.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// CustomerRoute mewadahi controller registry & router Fiber.
type CustomerRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// ICustomerRoute adalah kontrak pendaftaran rute domain pelanggan.
type ICustomerRoute interface {
	Run()
}

// NewCustomerRoute membuat instance baru.
func NewCustomerRoute(controller controllers.IControllerRegistry, router fiber.Router) ICustomerRoute {
	return &CustomerRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain pelanggan.
func (r *CustomerRoute) Run() {
	api := r.router.Group("/outlets/:id_outlet/customers", middleware.LaundryJWTMiddleware, middleware.RequireOutletMember)

	api.Get("/", r.controller.GetCustomer().List)
	api.Post("/", middleware.ValidateForm(&dto.CustomerForm{}), r.controller.GetCustomer().Create)
	api.Get("/:id_customer", r.controller.GetCustomer().Get)
	api.Put("/:id_customer", middleware.ValidateForm(&dto.CustomerForm{}), r.controller.GetCustomer().Update)
	api.Delete("/:id_customer", r.controller.GetCustomer().Delete)
	api.Get("/:id_customer/orders", r.controller.GetCustomer().ListOrders)
}
