package controllers

import (
	"sso-service/models"

	"github.com/gofiber/fiber/v2"
)

func GetMasterProdi(c *fiber.Ctx) error {
	data := models.GetMasterProdi()
	return c.JSON(data)
}

type GetAngkatanParams struct {
	IDProdi string `json:"id_prodi" form:"id_prodi" validate:"required,uuid4"`
}

func GetAngkatan(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetAngkatanParams)
	data := models.GetAngkatan(params.IDProdi)
	return c.JSON(data)
}
