// Package controllers (mstmenu) adalah lapisan HTTP handler (Fiber) untuk domain master menu & modul.
package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// MstMenuController membungkus service registry.
type MstMenuController struct {
	service services.IServiceRegistry
}

// IMstMenuController adalah kontrak handler HTTP domain master menu & modul.
type IMstMenuController interface {
	GetMstMenu(*fiber.Ctx) error
	GetDetailMstMenu(*fiber.Ctx) error
	GetMstMenuModul(*fiber.Ctx) error
	GetDetailMstMenuModul(*fiber.Ctx) error
	CreateMstMenu(*fiber.Ctx) error
	CreateMstMenuModul(*fiber.Ctx) error
	UpdateMstMenu(*fiber.Ctx) error
	UpdateMstMenuModul(*fiber.Ctx) error
	DeleteMstMenu(*fiber.Ctx) error
	DeleteMstMenuModul(*fiber.Ctx) error
}

// NewMstMenuController membuat instance MstMenuController baru.
func NewMstMenuController(service services.IServiceRegistry) IMstMenuController {
	return &MstMenuController{service: service}
}

func (ctrl *MstMenuController) GetMstMenu(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetMstMenuParams)
	form := c.Locals("validatedForm").(*dto.GetMstMenuForm)
	data := ctrl.service.GetMstMenu().GetMstMenu(params.Id_Master_Aplikasi, form.Limit, form.Offset, form.Order, form.Filter)
	return c.JSON(data)
}

func (ctrl *MstMenuController) GetDetailMstMenu(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetDetailMstMenuParams)
	data := ctrl.service.GetMstMenu().GetDetailMstMenu(c.UserContext(), params.Id_Master_Menu)
	return c.JSON(data)
}

func (ctrl *MstMenuController) GetMstMenuModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetMstMenuModulParams)
	form := c.Locals("validatedForm").(*dto.GetMstMenuModulForm)
	data := ctrl.service.GetMstMenu().GetMstMenuModul(params.Id_Master_Menu, form.Limit, form.Offset, form.Order, form.Filter)
	return c.JSON(data)
}

func (ctrl *MstMenuController) GetDetailMstMenuModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetDetailMstMenuModulParams)
	data := ctrl.service.GetMstMenu().GetDetailMstMenuModul(c.UserContext(), params.Id_Master_Modul)
	return c.JSON(data)
}

func (ctrl *MstMenuController) CreateMstMenu(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.CreateMstMenuForm)
	data := ctrl.service.GetMstMenu().CreateMstMenu(c.UserContext(), form.Id_Master_Aplikasi, form.Nama_Menu, form.Deskripsi, form.Order, form.Icon)
	return c.JSON(data)
}

func (ctrl *MstMenuController) CreateMstMenuModul(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.CreateMstMenuModulForm)
	data := ctrl.service.GetMstMenu().CreateMstMenuModul(c.UserContext(), form.Id_Master_Aplikasi, form.Id_Master_Menu, form.Nama_Modul, form.Path, form.Deskripsi, form.Order, form.Icon)
	return c.JSON(data)
}

func (ctrl *MstMenuController) UpdateMstMenu(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.UpdateMstMenuParams)
	form := c.Locals("validatedForm").(*dto.UpdateMstMenuForm)
	data := ctrl.service.GetMstMenu().UpdateMstMenu(c.UserContext(), params.Id_Master_Menu, form.Nama_Menu, form.Deskripsi, form.Order, form.Icon)
	return c.JSON(data)
}

func (ctrl *MstMenuController) UpdateMstMenuModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.UpdateMstMenuModulParams)
	form := c.Locals("validatedForm").(*dto.UpdateMstMenuModulForm)
	data := ctrl.service.GetMstMenu().UpdateMstMenuModul(c.UserContext(), params.Id_Master_Modul, form.Id_Master_Aplikasi, form.Id_Master_Menu, form.Nama_Modul, form.Path, form.Deskripsi, form.Order, form.Icon)
	return c.JSON(data)
}

func (ctrl *MstMenuController) DeleteMstMenu(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.DeleteMstMenuParams)
	data := ctrl.service.GetMstMenu().DeleteMstMenu(c.UserContext(), params.Id_Master_Menu)
	return c.JSON(data)
}

func (ctrl *MstMenuController) DeleteMstMenuModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.DeleteMstMenuModulParams)
	data := ctrl.service.GetMstMenu().DeleteMstMenuModul(c.UserContext(), params.Id_Master_Modul)
	return c.JSON(data)
}
