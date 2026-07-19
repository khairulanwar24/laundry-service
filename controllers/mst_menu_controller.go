package controllers

import (
	"sso-service/models"

	"github.com/gofiber/fiber/v2"
)

type GetMstMenuParams struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type GetMstMenuForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

func GetMstMenu(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetMstMenuParams)
	id_master_aplikasi := params.Id_Master_Aplikasi

	form := c.Locals("validatedForm").(*GetMstMenuForm)

	limit := form.Limit
	offset := form.Offset
	order := form.Order
	filter := form.Filter

	data := models.GetMstMenu(id_master_aplikasi, limit, offset, order, filter)

	return c.JSON(data)
}

type GetDetailMstMenuParams struct {
	Id_Master_Menu string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
}

func GetDetailMstMenu(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*GetDetailMstMenuParams)

	id_master_menu := params.Id_Master_Menu

	data := models.GetDetailMstMenu(id_master_menu)

	return c.JSON(data)
}

type GetMstMenuModulParams struct {
	Id_Master_Menu string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
}

type GetMstMenuModulForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

type GetDetailMstMenuModulParams struct {
	Id_Master_Modul string `json:"id_master_modul" form:"id_master_modul" validate:"required,uuid4"`
}

func GetDetailMstMenuModul(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*GetDetailMstMenuModulParams)

	id_master_modul := params.Id_Master_Modul

	data := models.GetDetailMstMenuModul(id_master_modul)

	return c.JSON(data)
}

func GetMstMenuModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetMstMenuModulParams)
	id_master_menu := params.Id_Master_Menu

	form := c.Locals("validatedForm").(*GetMstMenuModulForm)

	limit := form.Limit
	offset := form.Offset
	order := form.Order
	filter := form.Filter

	data := models.GetMstMenuModul(id_master_menu, limit, offset, order, filter)

	return c.JSON(data)
}

type CreateMstMenuForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Nama_Menu          string `json:"nama_menu" form:"nama_menu" validate:"required"`
	Deskripsi          string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Order              string `json:"order" form:"order" validate:"required"`
	Icon               string `json:"icon" form:"icon" validate:"required"`
}

func CreateMstMenu(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*CreateMstMenuForm)

	id_master_aplikasi := form.Id_Master_Aplikasi
	nama_menu := form.Nama_Menu
	deskripsi := form.Deskripsi
	order := form.Order
	icon := form.Icon

	data := models.CreateMstMenu(id_master_aplikasi, nama_menu, deskripsi, order, icon)

	return c.JSON(data)
}

type CreateMstMenuModulForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Menu     string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
	Nama_Modul         string `json:"nama_modul" form:"nama_modul" validate:"required"`
	Path               string `json:"path" form:"path" validate:"required"`
	Deskripsi          string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Order              string `json:"order" form:"order" validate:"required"`
	Icon               string `json:"icon" form:"icon" validate:"required"`
}

func CreateMstMenuModul(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*CreateMstMenuModulForm)

	id_master_aplikasi := form.Id_Master_Aplikasi
	id_master_menu := form.Id_Master_Menu
	nama_modul := form.Nama_Modul
	path := form.Path
	deskripsi := form.Deskripsi
	order := form.Order
	icon := form.Icon

	data := models.CreateMstMenuModul(id_master_aplikasi, id_master_menu, nama_modul, path, deskripsi, order, icon)

	return c.JSON(data)
}

type UpdateMstMenuParams struct {
	Id_Master_Menu string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
}

type UpdateMstMenuForm struct {
	Nama_Menu string `json:"nama_menu" form:"nama_menu" validate:"required"`
	Deskripsi string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Order     string `json:"order" form:"order" validate:"required,numeric"`
	Icon      string `json:"icon" form:"icon" validate:"required"`
}

func UpdateMstMenu(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*UpdateMstMenuParams)
	id_master_menu := params.Id_Master_Menu

	form := c.Locals("validatedForm").(*UpdateMstMenuForm)

	nama_menu := form.Nama_Menu
	deskripsi := form.Deskripsi
	order := form.Order
	icon := form.Icon

	data := models.UpdateMstMenu(id_master_menu, nama_menu, deskripsi, order, icon)

	return c.JSON(data)
}

type UpdateMstMenuModulParams struct {
	Id_Master_Modul string `json:"id_master_modul" form:"id_master_modul" validate:"required,uuid4"`
}

type UpdateMstMenuModulForm struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
	Id_Master_Menu     string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
	Nama_Modul         string `json:"nama_modul" form:"nama_modul" validate:"required"`
	Path               string `json:"path" form:"path" validate:"required"`
	Deskripsi          string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Order              string `json:"order" form:"order" validate:"required"`
	Icon               string `json:"icon" form:"icon" validate:"required"`
}

func UpdateMstMenuModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*UpdateMstMenuModulParams)
	id_master_modul := params.Id_Master_Modul

	form := c.Locals("validatedForm").(*UpdateMstMenuModulForm)

	id_master_aplikasi := form.Id_Master_Aplikasi
	id_master_menu := form.Id_Master_Menu
	nama_modul := form.Nama_Modul
	path := form.Path
	deskripsi := form.Deskripsi
	order := form.Order
	icon := form.Icon

	data := models.UpdateMstMenuModul(id_master_modul, id_master_aplikasi, id_master_menu, nama_modul, path, deskripsi, order, icon)

	return c.JSON(data)
}

type DeleteMstMenuParams struct {
	Id_Master_Menu string `json:"id_master_menu" form:"id_master_menu" validate:"required,uuid4"`
}

func DeleteMstMenu(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*DeleteMstMenuParams)

	id_master_menu := params.Id_Master_Menu

	data := models.DeleteMstMenu(id_master_menu)

	return c.JSON(data)
}

type DeleteMstMenuModulParams struct {
	Id_Master_Modul string `json:"id_master_modul" form:"id_master_modul" validate:"required,uuid4"`
}

func DeleteMstMenuModul(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*DeleteMstMenuModulParams)

	id_master_modul := params.Id_Master_Modul

	data := models.DeleteMstMenuModul(id_master_modul)

	return c.JSON(data)
}
