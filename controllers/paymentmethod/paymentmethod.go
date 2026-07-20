// Package controllers (paymentmethod) adalah lapisan HTTP handler (Fiber) untuk domain metode pembayaran.
package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// PaymentMethodController membungkus service registry.
type PaymentMethodController struct {
	service services.IServiceRegistry
}

// IPaymentMethodController adalah kontrak handler HTTP domain metode pembayaran.
type IPaymentMethodController interface {
	ListActive(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Deactivate(*fiber.Ctx) error
}

// NewPaymentMethodController membuat instance baru.
func NewPaymentMethodController(service services.IServiceRegistry) IPaymentMethodController {
	return &PaymentMethodController{service: service}
}

func (ctrl *PaymentMethodController) ListActive(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetPaymentMethod().ListActive(c.UserContext(), idOutlet)
	return c.JSON(data)
}

func (ctrl *PaymentMethodController) Create(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	form := c.Locals("validatedForm").(*dto.PaymentMethodForm)
	data := ctrl.service.GetPaymentMethod().Create(c.UserContext(), idOutlet, *form)
	return c.JSON(data)
}

func (ctrl *PaymentMethodController) Update(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idPaymentMethod := c.Params("id_payment_method")
	form := c.Locals("validatedForm").(*dto.PaymentMethodForm)
	data := ctrl.service.GetPaymentMethod().Update(c.UserContext(), idOutlet, idPaymentMethod, *form)
	return c.JSON(data)
}

func (ctrl *PaymentMethodController) Deactivate(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idPaymentMethod := c.Params("id_payment_method")
	data := ctrl.service.GetPaymentMethod().Deactivate(c.UserContext(), idOutlet, idPaymentMethod)
	return c.JSON(data)
}
