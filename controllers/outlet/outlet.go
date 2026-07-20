// Package controllers (outlet) adalah lapisan HTTP handler (Fiber) untuk domain outlet.
package controllers

import (
	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// OutletController membungkus service registry.
type OutletController struct {
	service services.IServiceRegistry
}

// IOutletController adalah kontrak handler HTTP domain outlet.
type IOutletController interface {
	ListMyOutlets(*fiber.Ctx) error
	CreateOutlet(*fiber.Ctx) error
	GetOutlet(*fiber.Ctx) error
	UpdateOutlet(*fiber.Ctx) error
	ListStaff(*fiber.Ctx) error
	InviteEmployee(*fiber.Ctx) error
	UpdateMember(*fiber.Ctx) error
	RemoveMember(*fiber.Ctx) error
}

// NewOutletController membuat instance baru.
func NewOutletController(service services.IServiceRegistry) IOutletController {
	return &OutletController{service: service}
}

func (ctrl *OutletController) ListMyOutlets(c *fiber.Ctx) error {
	idUser, _ := c.Locals("id_user").(string)
	data := ctrl.service.GetOutlet().ListMyOutlets(c.UserContext(), idUser)
	return c.JSON(data)
}

func (ctrl *OutletController) CreateOutlet(c *fiber.Ctx) error {
	idUser, _ := c.Locals("id_user").(string)
	form := c.Locals("validatedForm").(*dto.CreateOutletForm)
	data := ctrl.service.GetOutlet().CreateOutlet(c.UserContext(), idUser, *form)
	return c.JSON(data)
}

func (ctrl *OutletController) GetOutlet(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetOutlet().GetOutlet(c.UserContext(), idOutlet)
	return c.JSON(data)
}

func (ctrl *OutletController) UpdateOutlet(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	form := c.Locals("validatedForm").(*dto.UpdateOutletForm)
	data := ctrl.service.GetOutlet().UpdateOutlet(c.UserContext(), idOutlet, *form)
	return c.JSON(data)
}

func (ctrl *OutletController) ListStaff(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	data := ctrl.service.GetOutlet().ListStaff(c.UserContext(), idOutlet)
	return c.JSON(data)
}

func (ctrl *OutletController) InviteEmployee(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	form := c.Locals("validatedForm").(*dto.InviteEmployeeForm)
	data := ctrl.service.GetOutlet().InviteEmployee(c.UserContext(), idOutlet, *form)
	return c.JSON(data)
}

func (ctrl *OutletController) UpdateMember(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idMember := c.Params("id_member")
	form := c.Locals("validatedForm").(*dto.UpdateMemberForm)
	data := ctrl.service.GetOutlet().UpdateMember(c.UserContext(), idOutlet, idMember, *form)
	return c.JSON(data)
}

func (ctrl *OutletController) RemoveMember(c *fiber.Ctx) error {
	idOutlet := c.Params("id_outlet")
	idMember := c.Params("id_member")
	data := ctrl.service.GetOutlet().RemoveMember(c.UserContext(), idOutlet, idMember)
	return c.JSON(data)
}
