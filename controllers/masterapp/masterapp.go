// Package controllers (masterapp) adalah lapisan HTTP handler (Fiber) untuk domain master aplikasi.
package controllers

import (
	"fmt"
	"time"

	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// MasterAppController membungkus service registry.
type MasterAppController struct {
	service services.IServiceRegistry
}

// IMasterAppController adalah kontrak handler HTTP domain master aplikasi.
type IMasterAppController interface {
	GetMasterApps(*fiber.Ctx) error
	GetMasterAppById(*fiber.Ctx) error
	CreateMasterApp(*fiber.Ctx) error
	UpdateMasterApps(*fiber.Ctx) error
	DeleteMasterApp(*fiber.Ctx) error
}

// NewMasterAppController membuat instance MasterAppController baru.
func NewMasterAppController(service services.IServiceRegistry) IMasterAppController {
	return &MasterAppController{service: service}
}

func (ctrl *MasterAppController) GetMasterApps(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.GetData)
	data := ctrl.service.GetMasterApp().GetMasterApps(form.Order, form.Filter, form.Limit, form.Offset)
	return c.JSON(data)
}

func (ctrl *MasterAppController) GetMasterAppById(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetMasterAppByIdParams)
	data := ctrl.service.GetMasterApp().GetMasterAppById(c.UserContext(), params.Id_master_aplikasi)
	return c.JSON(data)
}

func (ctrl *MasterAppController) CreateMasterApp(c *fiber.Ctx) error {
	// Ambil file image (jika ada)
	_, errs := c.FormFile("image")
	image := ""
	if errs == nil {
		maxFileSize := int64(2 * 1024 * 1024) // 2 MB
		allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

		err := middleware.FileUploadToS3Middleware(c, "image", "apps", allowedTypes, maxFileSize)
		if err != nil {
			fmt.Println("File upload error:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "File upload failed: " + err.Error(),
				"data":    nil,
			})
		}
		image = c.Locals("fileName").(string)
	}

	form := c.Locals("validatedForm").(*dto.CreateMasterAppForm)

	// Konversi tgl_version dari string ke time.Time
	layout := "2006-01-02"
	tglVersion, err := time.Parse(layout, form.Tgl_Version)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid date format",
		})
	}

	data := ctrl.service.GetMasterApp().CreateMasterApp(c.UserContext(), form.Nama_Aplikasi, form.Deskripsi, tglVersion, form.Url, form.Versi_Aplikasi, image)
	return c.JSON(data)
}

func (ctrl *MasterAppController) UpdateMasterApps(c *fiber.Ctx) error {
	// Handle file upload
	_, errs := c.FormFile("image")
	image := ""
	if errs == nil {
		maxFileSize := int64(2 * 1024 * 1024) // 2 MB
		allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

		err := middleware.FileUploadToS3Middleware(c, "image", "apps", allowedTypes, maxFileSize)
		if err != nil {
			fmt.Println("File upload error:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "File upload failed: " + err.Error(),
				"data":    nil,
			})
		}
		image = c.Locals("fileName").(string)
	}

	params := c.Locals("validatedParams").(*dto.UpdateMasterAppsParams)
	form := c.Locals("validatedForm").(*dto.UpdateMasterAppsForm)

	// Convert tgl_version from string to time.Time
	tglVersion, err := time.Parse("2006-01-02", form.Tgl_Version)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid date format for tgl_version",
			"data":    nil,
		})
	}

	data := ctrl.service.GetMasterApp().UpdateMasterApp(c.UserContext(), params.Id_master_aplikasi, image, form.Nama_Aplikasi, form.Deskripsi, form.Versi_Aplikasi, tglVersion, form.URL)
	return c.JSON(data)
}

func (ctrl *MasterAppController) DeleteMasterApp(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.DeleteMasterAppParams)
	data := ctrl.service.GetMasterApp().DeleteMasterApp(c.UserContext(), params.Id_Master_Aplikasi)
	return c.JSON(data)
}
