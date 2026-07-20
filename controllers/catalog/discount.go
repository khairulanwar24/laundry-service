package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// DiscountController membungkus service registry.
type DiscountController struct {
	service services.IServiceRegistry
}

// IDiscountController adalah kontrak handler HTTP domain diskon.
type IDiscountController interface {
	ListActive(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
}

// NewDiscountController membuat instance baru.
func NewDiscountController(service services.IServiceRegistry) IDiscountController {
	return &DiscountController{service: service}
}

func (ctrl *DiscountController) ListActive(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetDiscount().ListActive(c.UserContext(), idOutlet)
	return c.JSON(data)
}

func (ctrl *DiscountController) Create(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	form := c.Locals("validatedForm").(*dto.DiscountForm)
	data := ctrl.service.GetDiscount().Create(c.UserContext(), idOutlet, *form)
	return c.JSON(data)
}

func (ctrl *DiscountController) Update(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idDiscount := c.Params("id_discount")
	form := c.Locals("validatedForm").(*dto.DiscountForm)
	data := ctrl.service.GetDiscount().Update(c.UserContext(), idOutlet, idDiscount, *form)
	return c.JSON(data)
}
