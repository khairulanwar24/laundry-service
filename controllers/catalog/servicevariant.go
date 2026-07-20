package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// ServiceVariantController membungkus service registry.
type ServiceVariantController struct {
	service services.IServiceRegistry
}

// IServiceVariantController adalah kontrak handler HTTP domain varian layanan.
type IServiceVariantController interface {
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
}

// NewServiceVariantController membuat instance baru.
func NewServiceVariantController(service services.IServiceRegistry) IServiceVariantController {
	return &ServiceVariantController{service: service}
}

func (ctrl *ServiceVariantController) Create(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idService := c.Params("id_service")
	form := c.Locals("validatedForm").(*dto.ServiceVariantForm)
	data := ctrl.service.GetServiceVariant().Create(c.UserContext(), idOutlet, idService, *form)
	return c.JSON(data)
}

func (ctrl *ServiceVariantController) Update(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idServiceVariant := c.Params("id_service_variant")
	form := c.Locals("validatedForm").(*dto.ServiceVariantForm)
	data := ctrl.service.GetServiceVariant().Update(c.UserContext(), idOutlet, idServiceVariant, *form)
	return c.JSON(data)
}
