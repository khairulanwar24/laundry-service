package controllers

import (
	"fmt"
	middleware "sso-service/middlewares"
	"sso-service/models"
	"sso-service/types"
	"time"

	"github.com/gofiber/fiber/v2"
)

// type GetMasterAppsParams struct {
// 	Id_Master_Aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
// }

type GetMasterAppsForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100" default:"10"`
	Offset int    `json:"offset" form:"offset" validate:"numeric"`
	Order  string `json:"order" form:"order" validate:""`
	Filter string `json:"filter" form:"filter" validate:""`
}

func GetMasterApps(c *fiber.Ctx) error {
	// Retrieve the validated form
	form := c.Locals("validatedForm").(*types.GetData)

	// Retrieve values ​​from a validated form

	// Call the model to get data from the database
	data := models.GetMasterApps(form)

	return c.JSON(data)
}

type GetMasterAppByIdParams struct {
	Id_master_aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
}

func GetMasterAppById(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetMasterAppByIdParams)

	id_master_aplikasi := params.Id_master_aplikasi

	data := models.GetMasterAppById(id_master_aplikasi)

	return c.JSON(data)
}

type CreateMasterAppForm struct {
	Nama_Aplikasi  string `json:"nama_aplikasi" form:"nama_aplikasi" validate:"required"`
	Deskripsi      string `json:"deskripsi" form:"deskripsi" validate:""`
	Versi_Aplikasi string `json:"versi_aplikasi" form:"versi_aplikasi" validate:"required"`
	Tgl_Version    string `json:"tgl_version" form:"tgl_version" validate:"required,datetime=2006-01-02"`
	Url            string `json:"url" form:"url" validate:"required,url"`
}

func CreateMasterApp(c *fiber.Ctx) error {

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

	// Ambil form yang telah divalidasi dari context
	form := c.Locals("validatedForm").(*CreateMasterAppForm)
	nama_aplikasi := form.Nama_Aplikasi
	deskripsi := form.Deskripsi
	versi_aplikasi := form.Versi_Aplikasi
	tgl_version := form.Tgl_Version
	url := form.Url

	// Konversi tgl_version dari string ke time.Time
	layout := "2006-01-02" // Format tanggal yang sesuai dengan input
	tglVersion, err := time.Parse(layout, tgl_version)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid date format",
		})
	}

	// Panggil model untuk membuat master aplikasi
	data := models.CreateMasterApps(
		nama_aplikasi, deskripsi, tglVersion, url, versi_aplikasi, image,
	)

	// Kembalikan response
	return c.JSON(data)
}

type UpdateMasterAppsParams struct {
	Id_master_aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
}

type UpdateMasterAppsForm struct {
	Nama_Aplikasi  string `json:"nama_aplikasi" form:"nama_aplikasi" validate:"required"`
	Deskripsi      string `json:"deskripsi" form:"deskripsi" validate:"required"`
	Versi_Aplikasi string `json:"versi_aplikasi" form:"versi_aplikasi" validate:"required"`
	Tgl_Version    string `json:"tgl_version" form:"tgl_version" validate:"required"` // Ubah menjadi string
	URL            string `json:"url" form:"url" validate:"required"`
}

func UpdateMasterApps(c *fiber.Ctx) error {
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

	// Validate parameters
	params := c.Locals("validatedParams").(*UpdateMasterAppsParams)
	id_master_aplikasi := params.Id_master_aplikasi

	// Validate form data
	form := c.Locals("validatedForm").(*UpdateMasterAppsForm)

	nama_aplikasi := form.Nama_Aplikasi
	deskripsi := form.Deskripsi
	versi_aplikasi := form.Versi_Aplikasi

	// Convert tgl_version from string to time.Time
	tgl_version, err := time.Parse("2006-01-02", form.Tgl_Version)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid date format for tgl_version",
			"data":    nil,
		})
	}

	url := form.URL

	// Call model to update master app
	data := models.UpdateMasterApp(id_master_aplikasi, image, nama_aplikasi, deskripsi, versi_aplikasi, tgl_version, url)

	return c.JSON(data)
}

type DeleteMasterAppParams struct {
	Id_Master_Aplikasi string `json:"id_master_aplikasi" validate:"required,uuid4"`
}

func DeleteMasterApp(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*DeleteMasterAppParams)
	id_master_aplikasi := params.Id_Master_Aplikasi

	data := models.DeleteMasterApp(id_master_aplikasi)

	return c.JSON(data)
}

// type UpdateMasterAppsParams struct {
// 	Id_master_aplikasi string `json:"id_master_aplikasi" form:"id_master_aplikasi" validate:"required,uuid4"`
// }

// type UpdateMasterAppsForm struct {
// 	Nama_Aplikasi  string    `json:"nama_aplikasi" form:"nama_aplikasi" validate:"required"`
// 	Deskripsi      string    `json:"deskripsi" form:"deskripsi" validate:"required"`
// 	Versi_Aplikasi string    `json:"versi_aplikasi" form:"versi_aplikasi" validate:"required"`
// 	Tgl_Version    time.Time `json:"tgl_version" form:"tgl_version"`
// 	URL            string    `json:"url" form:"url" validate:"required"`
// }

// func UpdateMasterApps(c *fiber.Ctx) error {
// 	_, errs := c.FormFile("image")
// 	image := ""
// 	if errs == nil {
// 		maxFileSize := int64(2 * 1024 * 1024)
// 		allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

// 		err := middleware.FileUploadToS3Middleware(c, "image", "images", allowedTypes, maxFileSize)
// 		if err != nil {
// 			fmt.Println("File upload error:", err)
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"success": false,
// 				"message": "File upload failed: " + err.Error(),
// 				"data":    nil,
// 			})
// 		}
// 		image = c.Locals("filename").(string)
// 	}
// 	params := c.Locals("validatedParams").(*UpdateMasterAppsParams)
// 	id_master_aplikasi := params.Id_master_aplikasi
// 	form := c.Locals("validatedForm").(*UpdateMasterAppsForm)

// 	nama_aplikasi := form.Nama_Aplikasi
// 	deskripsi := form.Deskripsi
// 	versi_aplikasi := form.Versi_Aplikasi
// 	tgl_version := form.Tgl_Version
// 	url := form.URL

// 	data := models.UpdateMasterApp(id_master_aplikasi, image, nama_aplikasi, deskripsi, versi_aplikasi, tgl_version, url)

// 	return c.JSON(data)
// }

// func GetMasterApps(c *fiber.Ctx) error {
// 	limit := c.FormValue("limit")
// 	limitParam, err := strconv.Atoi(limit)
// 	if err != nil {
// 		return c.JSON("paramater limit harus angka")
// 	}

// 	offset := c.FormValue("offset")
// 	offsetParam, err := strconv.Atoi(offset)
// 	if err != nil {
// 		return c.JSON("parameter offset harus angka")
// 	}

// 	order := c.FormValue("order")
// 	filter := c.FormValue("filter")

// 	// Log parameter yang diterima dari request untuk debugging
// 	fmt.Println("Limit:", limitParam)
// 	fmt.Println("Offset:", offsetParam)
// 	fmt.Println("Order:", order)
// 	fmt.Println("Filter:", filter)

// 	data := models.GetMasterApps(limitParam, offsetParam, order, filter)
// 	return c.JSON(data)

// }

// func GetMasterAppById(c *fiber.Ctx) error {
// 	id_master_aplikasi := c.Params("id_master_aplikasi")
// 	data := models.GetMasterAppById(id_master_aplikasi)

// 	return c.JSON(data)
// }

// func CreateMasterApps(c *fiber.Ctx) error {
// 	nama := c.FormValue("nama_aplikasi")
// 	deskripsi := c.FormValue("deskripsi")
// 	tgl_version := c.FormValue("tgl_version")
// 	url := c.FormValue("url")
// 	versi_aplikasi := c.FormValue("versi_aplikasi")
// 	image := c.FormValue("image")

// 	// Konversi tgl_version dari string ke time.Time
// 	layout := "2006-01-02 15:04:05" // Sesuaikan dengan format datetime yang Anda terima
// 	tglVersion, err := time.Parse(layout, tgl_version)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"success": false,
// 			"message": "Invalid date format",
// 		})
// 	}

// 	data := models.CreateMasterApps(nama, deskripsi, tglVersion, url, versi_aplikasi, image)

// 	return c.JSON(data)

// }
