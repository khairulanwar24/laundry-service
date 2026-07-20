// Package controllers (groupakses) adalah lapisan HTTP handler (Fiber) untuk domain group akses.
package controllers

import (
	"fmt"

	"laundry-service/domain/dto"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// GroupAksesController membungkus service registry.
type GroupAksesController struct {
	service services.IServiceRegistry
}

// IGroupAksesController adalah kontrak handler HTTP domain group akses.
type IGroupAksesController interface {
	CreateMstGroupAkses(*fiber.Ctx) error
	GetMstGroupAkses(*fiber.Ctx) error
	GetMstGroupAksesModul(*fiber.Ctx) error
	GetGroupAkses(*fiber.Ctx) error
	UpdateMstGroupAkses(*fiber.Ctx) error
	GetDetailMstGroupAkses(*fiber.Ctx) error
	DeleteMstGroupAkses(*fiber.Ctx) error
	CreateGroupAkses(*fiber.Ctx) error
	DeleteGroupAkses(*fiber.Ctx) error
	GetGroupAksesUserMenu(*fiber.Ctx) error
	GetGroupAksesUserApps(*fiber.Ctx) error
	CreateGroupAksesUserApps(*fiber.Ctx) error
	BulkCreateGroupAksesUserApps(*fiber.Ctx) error
	UpdateGroupAksesUserApps(*fiber.Ctx) error
	DeleteGroupAksesUserApps(*fiber.Ctx) error
}

// NewGroupAksesController membuat instance GroupAksesController baru.
func NewGroupAksesController(service services.IServiceRegistry) IGroupAksesController {
	return &GroupAksesController{service: service}
}

func (ctrl *GroupAksesController) CreateMstGroupAkses(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.CreateMstGroupAksesForm)
	data := ctrl.service.GetGroupAkses().CreateMstGroupAkses(c.UserContext(), form.Id_Master_Aplikasi, form.Nama_Group, form.Deskripsi)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) GetMstGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetMstGroupAksesParams)
	form := c.Locals("validatedForm").(*dto.GetMstGroupAksesForm)
	data := ctrl.service.GetGroupAkses().GetMstGroupAkses(params.Id_Master_Aplikasi, form.Limit, form.Offset, form.Order, form.Filter)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) GetMstGroupAksesModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetMstGroupAksesModulParams)
	form := c.Locals("validatedForm").(*dto.GetMstGroupAksesModulForm)
	data := ctrl.service.GetGroupAkses().GetMstGroupAksesModul(params.Id_Master_Aplikasi, params.Id_Master_Group, form.Limit, form.Offset, form.Order, form.Filter)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) GetGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetGroupAksesParams)
	form := c.Locals("validatedForm").(*dto.GetGroupAksesForm)
	data := ctrl.service.GetGroupAkses().GetGroupAkses(params.Id_Master_Group, form.Limit, form.Offset, form.Order, form.Filter)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) UpdateMstGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.UpdateMstGroupAksesParams)
	form := c.Locals("validatedForm").(*dto.UpdateMstGroupAksesForm)
	data := ctrl.service.GetGroupAkses().UpdateMstGroupAkses(c.UserContext(), params.Id_Master_Group, form.Nama_Group, form.Deskripsi)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) GetDetailMstGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetDetailMstGroupAksesParams)
	data := ctrl.service.GetGroupAkses().GetDetailMstGroupAkses(c.UserContext(), params.Id_Master_Group)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) DeleteMstGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.DeleteMstGroupAksesParams)
	data := ctrl.service.GetGroupAkses().DeleteMstGroupAkses(c.UserContext(), params.Id_Master_Group)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) CreateGroupAkses(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.CreateGroupAksesForm)
	akses := "1,1,1,1"
	data := ctrl.service.GetGroupAkses().CreateGroupAkses(c.UserContext(), form.Id_Master_Aplikasi, form.Id_Master_Group, form.Id_Master_Modul, akses)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) DeleteGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.DeleteGroupAksesParams)
	data := ctrl.service.GetGroupAkses().DeleteGroupAkses(c.UserContext(), params.Id_Group_Akses)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) GetGroupAksesUserMenu(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetGroupAksesUserMenuParams)
	data := ctrl.service.GetGroupAkses().GetGroupAksesUserMenu(c.UserContext(), params.Id_User, params.Id_Master_Aplikasi)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) GetGroupAksesUserApps(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetGroupAksesUserAppsParams)
	data := ctrl.service.GetGroupAkses().GetGroupAksesUserApps(c.UserContext(), params.Id_User)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) CreateGroupAksesUserApps(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.CreateGroupAksesUserAppsForm)
	data := ctrl.service.GetGroupAkses().CreateGroupAksesUserApps(c.UserContext(), form.Id_User, form.Id_Master_Aplikasi, form.Id_Master_Group, form.Status_Data)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) BulkCreateGroupAksesUserApps(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.BulkGroupAksesUserAppsForm)
	resp := ctrl.service.GetGroupAkses().BulkCreateGroupAksesUserApps(c.UserContext(), form.IDUsers, form.IDMasterAplikasi, form.IDMasterGroup, form.StatusData)
	return c.JSON(resp)
}

func (ctrl *GroupAksesController) UpdateGroupAksesUserApps(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.UpdateGroupAksesUserAppsParams)
	form := c.Locals("validatedForm").(*dto.UpdateGroupAksesUserAppsForm)
	data := ctrl.service.GetGroupAkses().UpdateGroupAksesUserApps(c.UserContext(), params.Id_Trans_User_Group, form.Status_Data)
	return c.JSON(data)
}

func (ctrl *GroupAksesController) DeleteGroupAksesUserApps(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.DeleteGroupAksesUserAppsParams)
	fmt.Println(params.Id_Trans_User_Group)
	data := ctrl.service.GetGroupAkses().DeleteGroupAksesUserApps(c.UserContext(), params.Id_Trans_User_Group)
	return c.JSON(data)
}
