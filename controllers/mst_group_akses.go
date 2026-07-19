package controllers

import (
	"fmt"
	"sso-service/models"

	"github.com/gofiber/fiber/v2"
)

type CreateMstGroupAksesForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Nama_Group         string `json:"nama_group" form:"nama_group" validate:"required"`
	Deskripsi          string `json:"deskripsi" form:"deskripsi" validate:"required"`
}

func CreateMstGroupAkses(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*CreateMstGroupAksesForm)

	id_master_aplikasi := form.Id_Master_Aplikasi
	nama_group := form.Nama_Group
	deskripsi := form.Deskripsi

	data := models.CreateMstGroupAkses(id_master_aplikasi, nama_group, deskripsi)

	return c.JSON(data)
}

type GetMstGroupAksesParams struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type GetMstGroupAksesForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

func GetMstGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetMstGroupAksesParams)
	id_master_aplikasi := params.Id_Master_Aplikasi

	form := c.Locals("validatedForm").(*GetMstGroupAksesForm)

	limit := form.Limit
	offset := form.Offset
	order := form.Order
	filter := form.Filter

	data := models.GetMstGroupAkses(id_master_aplikasi, limit, offset, order, filter)

	return c.JSON(data)
}

type GetMstGroupAksesModulParams struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Group    string `json:"id_master_group" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type GetMstGroupAksesModulForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

func GetMstGroupAksesModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetMstGroupAksesModulParams)
	id_master_aplikasi := params.Id_Master_Aplikasi
	id_master_group := params.Id_Master_Group

	form := c.Locals("validatedForm").(*GetMstGroupAksesModulForm)

	limit := form.Limit
	offset := form.Offset
	order := form.Order
	filter := form.Filter

	data := models.GetMstGroupAksesModul(id_master_aplikasi, id_master_group, limit, offset, order, filter)

	return c.JSON(data)
}

type GetGroupAksesParams struct {
	Id_Master_Group string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
}

type GetGroupAksesForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

func GetGroupAkses(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*GetGroupAksesParams)

	id_master_group := params.Id_Master_Group

	form := c.Locals("validatedForm").(*GetGroupAksesForm)

	limit := form.Limit
	offset := form.Offset
	order := form.Order
	filter := form.Filter

	data := models.GetGroupAkses(id_master_group, limit, offset, order, filter)
	return c.JSON(data)
}

type UpdateMstGroupAksesParams struct {
	Id_Master_Group string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
}

type UpdateMstGroupAksesForm struct {
	Nama_Group string `json:"nama_group" form:"nama_group" validate:"required"`
	Deskripsi  string `json:"deskripsi" form:"deskripsi" validate:"required"`
}

func UpdateMstGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*UpdateMstGroupAksesParams)
	id_master_group := params.Id_Master_Group

	form := c.Locals("validatedForm").(*UpdateMstGroupAksesForm)

	nama_group := form.Nama_Group
	deskripsi := form.Deskripsi

	data := models.UpdateMstGroupAkses(id_master_group, nama_group, deskripsi)

	return c.JSON(data)
}

type DeleteMstGroupAksesParams struct {
	Id_Master_Group string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
}

func DeleteMstGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*DeleteMstGroupAksesParams)

	id_master_group := params.Id_Master_Group

	data := models.DeleteMstGroupAkses(id_master_group)

	return c.JSON(data)
}

type GetDetailMstGroupAksesParams struct {
	Id_Master_Group string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
}

func GetDetailMstGroupAkses(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetDetailMstGroupAksesParams)

	id_master_group := params.Id_Master_Group

	data := models.GetDetailMstGroupAkses(id_master_group)

	return c.JSON(data)
}

type CreateGroupAksesForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Group    string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
	Id_Master_Modul    string `json:"id_master_modul" form:"id_master_modul" validate:"required,uuid4"`
}

func CreateGroupAkses(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*CreateGroupAksesForm)

	id_master_aplikasi := form.Id_Master_Aplikasi
	id_master_group := form.Id_Master_Group
	id_master_modul := form.Id_Master_Modul
	akses := "1,1,1,1"

	data := models.CreateGroupAkses(id_master_aplikasi, id_master_group, id_master_modul, akses)

	return c.JSON(data)
}

type DeleteGroupAksesParams struct {
	Id_Group_Akses string `json:"id_group_akses" form:"id_group_akses" validate:"required,uuid4"`
}

func DeleteGroupAkses(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*DeleteGroupAksesParams)
	id_group_akses := params.Id_Group_Akses

	data := models.DeleteGroupAkses(id_group_akses)

	return c.JSON(data)
}

type GetGroupAksesUserMenuParams struct {
	Id_User            string `json:"id_user" form:"id_user" validate:"required,uuid4"`
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_group_akses" validate:"required,uuid4"`
}

func GetGroupAksesUserMenu(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*GetGroupAksesUserMenuParams)
	id_user := params.Id_User
	id_master_aplikasi := params.Id_Master_Aplikasi

	data := models.GetGroupAksesUserMenu(id_user, id_master_aplikasi)
	return c.JSON(data)
}

type GetGroupAksesUserAppsParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

func GetGroupAksesUserApps(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*GetGroupAksesUserAppsParams)
	id_user := params.Id_User

	data := models.GetGroupAksesUserApps(id_user)
	return c.JSON(data)
}

type CreateGroupAksesUserAppsForm struct {
	Id_User            string `json:"id_user" form:"id_user" validate:"required,uuid4"`
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Group    string `json:"id_master_group" form:"id_master_group" validate:"required,uuid4"`
	Status_Data        string `json:"status_data" form:"status_data" validate:"required,oneof='true' 'false' "`
}

func CreateGroupAksesUserApps(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*CreateGroupAksesUserAppsForm)

	id_user := form.Id_User
	id_master_aplikasi := form.Id_Master_Aplikasi
	id_master_group := form.Id_Master_Group
	status_data := form.Status_Data

	data := models.CreateGroupAksesUserApps(id_user, id_master_aplikasi, id_master_group, status_data)

	return c.JSON(data)
}

type BulkGroupAksesUserAppsForm struct {
	IDUsers          []string `json:"id_users" validate:"required,min=1,dive,uuid"`
	IDMasterAplikasi string   `json:"id_master_aplikasi" validate:"required,uuid"`
	IDMasterGroup    string   `json:"id_master_group" validate:"required,uuid"`
	StatusData       bool     `json:"status_data" validate:"required"`
}

func BulkCreateGroupAksesUserApps(c *fiber.Ctx) error {

	form := c.Locals("validatedForm").(*BulkGroupAksesUserAppsForm)

	resp := models.BulkCreateGroupAksesUserApps(
		form.IDUsers,
		form.IDMasterAplikasi,
		form.IDMasterGroup,
		form.StatusData,
	)

	return c.JSON(resp)
}

type UpdateGroupAksesUserAppsParams struct {
	Id_Trans_User_Group string `json:"id_trans_user_group" form:"id_trans_user_group" validate:"required,uuid4"`
}

type UpdateGroupAksesUserAppsForm struct {
	Status_Data string `json:"status_data" form:"status_data" validate:"required,oneof='true' 'false' "`
}

func UpdateGroupAksesUserApps(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*UpdateGroupAksesUserAppsParams)
	id_trans_user_group := params.Id_Trans_User_Group

	form := c.Locals("validatedForm").(*UpdateGroupAksesUserAppsForm)

	status_data := form.Status_Data

	data := models.UpdateGroupAksesUserApps(id_trans_user_group, status_data)

	return c.JSON(data)
}

type DeleteGroupAksesUserAppsParams struct {
	Id_Trans_User_Group string `json:"id_trans_user_group" form:"id_trans_user_group" validate:"required,uuid4"`
}

func DeleteGroupAksesUserApps(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*DeleteGroupAksesUserAppsParams)

	id_trans_user_group := params.Id_Trans_User_Group

	fmt.Println(id_trans_user_group)

	data := models.DeleteGroupAksesUserApps(id_trans_user_group)

	return c.JSON(data)
}
