package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// PerfumeController membungkus service registry.
type PerfumeController struct {
	service services.IServiceRegistry
}

// IPerfumeController adalah kontrak handler HTTP domain parfum.
type IPerfumeController interface {
	ListActive(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
}

// NewPerfumeController membuat instance baru.
func NewPerfumeController(service services.IServiceRegistry) IPerfumeController {
	return &PerfumeController{service: service}
}

func (ctrl *PerfumeController) ListActive(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetPerfume().ListActive(c.UserContext(), idOutlet)
	return c.JSON(data)
}

func (ctrl *PerfumeController) Create(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	form := c.Locals("validatedForm").(*dto.PerfumeForm)
	data := ctrl.service.GetPerfume().Create(c.UserContext(), idOutlet, *form)
	return c.JSON(data)
}

func (ctrl *PerfumeController) Update(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idPerfume := c.Params("id_perfume")
	form := c.Locals("validatedForm").(*dto.PerfumeForm)
	data := ctrl.service.GetPerfume().Update(c.UserContext(), idOutlet, idPerfume, *form)
	return c.JSON(data)
}
