// Package controllers (catalog) adalah lapisan HTTP handler (Fiber) untuk domain katalog.
package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// ServiceController membungkus service registry.
type ServiceController struct {
	service services.IServiceRegistry
}

// IServiceController adalah kontrak handler HTTP domain layanan.
type IServiceController interface {
	ListActive(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Deactivate(*fiber.Ctx) error
}

// NewServiceController membuat instance baru.
func NewServiceController(service services.IServiceRegistry) IServiceController {
	return &ServiceController{service: service}
}

func (ctrl *ServiceController) ListActive(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	q := c.Query("q", "")
	data := ctrl.service.GetService().ListActive(c.UserContext(), idOutlet, q)
	return c.JSON(data)
}

func (ctrl *ServiceController) Create(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	form := c.Locals("validatedForm").(*dto.ServiceForm)
	data := ctrl.service.GetService().Create(c.UserContext(), idOutlet, *form)
	return c.JSON(data)
}

func (ctrl *ServiceController) Update(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idService := c.Params("id_service")
	form := c.Locals("validatedForm").(*dto.ServiceForm)
	data := ctrl.service.GetService().Update(c.UserContext(), idOutlet, idService, *form)
	return c.JSON(data)
}

func (ctrl *ServiceController) Deactivate(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idService := c.Params("id_service")
	data := ctrl.service.GetService().Deactivate(c.UserContext(), idOutlet, idService)
	return c.JSON(data)
}
