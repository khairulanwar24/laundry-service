// Package routes (paymentmethod) mendaftarkan endpoint HTTP untuk domain metode pembayaran.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// PaymentMethodRoute mewadahi controller registry & router Fiber.
type PaymentMethodRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IPaymentMethodRoute adalah kontrak pendaftaran rute domain metode pembayaran.
type IPaymentMethodRoute interface {
	Run()
}

// NewPaymentMethodRoute membuat instance baru.
func NewPaymentMethodRoute(controller controllers.IControllerRegistry, router fiber.Router) IPaymentMethodRoute {
	return &PaymentMethodRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain metode pembayaran.
func (r *PaymentMethodRoute) Run() {
	api := r.router.Group("/outlets/:id_outlet/payment-methods", middleware.LaundryJWTMiddleware, middleware.RequireOutletMember)

	api.Get("/", r.controller.GetPaymentMethod().ListActive)
	api.Post("/", middleware.RequireOutletOwner, middleware.ValidateForm(&dto.PaymentMethodForm{}), r.controller.GetPaymentMethod().Create)
	api.Put("/:id_payment_method", middleware.RequireOutletOwner, middleware.ValidateForm(&dto.PaymentMethodForm{}), r.controller.GetPaymentMethod().Update)
	api.Delete("/:id_payment_method", middleware.RequireOutletOwner, r.controller.GetPaymentMethod().Deactivate)
}
