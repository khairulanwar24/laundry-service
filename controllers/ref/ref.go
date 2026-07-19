// Package controllers (ref) adalah lapisan HTTP handler (Fiber) untuk domain referensi.
// Handler dibuat tipis: baca input -> panggil service -> balas response.
package controllers

import (
	"sso-service/common/response"
	"sso-service/domain/dto"
	"sso-service/services"

	"github.com/gofiber/fiber/v2"
)

// RefController membungkus service registry.
type RefController struct {
	service services.IServiceRegistry
}

// IRefController adalah kontrak handler HTTP domain referensi.
type IRefController interface {
	GetMasterProdi(*fiber.Ctx) error
	GetAngkatan(*fiber.Ctx) error
}

// NewRefController membuat instance RefController baru.
func NewRefController(service services.IServiceRegistry) IRefController {
	return &RefController{service: service}
}

// GetMasterProdi menangani GET /ref/get_prodi.
func (ctrl *RefController) GetMasterProdi(c *fiber.Ctx) error {
	data, err := ctrl.service.GetRef().GetMasterProdi(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, data)
}

// GetAngkatan menangani GET /ref/get_angkatan/:id_prodi.
func (ctrl *RefController) GetAngkatan(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetAngkatanParams)
	data, err := ctrl.service.GetRef().GetAngkatan(c.UserContext(), params.IDProdi)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, data)
}
