// Package controllers (user) adalah lapisan HTTP handler (Fiber) untuk domain user.
package controllers

import (
	"fmt"
	"strconv"

	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"
	"laundry-service/services"

	"github.com/gofiber/fiber/v2"
)

// UserController membungkus service registry.
type UserController struct {
	service services.IServiceRegistry
}

// IUserController adalah kontrak handler HTTP domain user.
type IUserController interface {
	GetUsers(*fiber.Ctx) error
	CreateUsers(*fiber.Ctx) error
	GetUser(*fiber.Ctx) error
	UpdateUsers(*fiber.Ctx) error
	UpdatePassword(*fiber.Ctx) error
	DeleteUser(*fiber.Ctx) error
	GetUsersDosen(*fiber.Ctx) error
	GetDetailDosen(*fiber.Ctx) error
	BulkCreateUsersMahasiswa(*fiber.Ctx) error
	GetDetailMahasiswa(*fiber.Ctx) error
	GetUsersMahasiswa(*fiber.Ctx) error
	GetUsersMahasiswaData(*fiber.Ctx) error
	GenerateUserMahasiswa(*fiber.Ctx) error
	GenerateUserDosen(*fiber.Ctx) error
	GenerateUserTendik(*fiber.Ctx) error
}

// NewUserController membuat instance UserController baru.
func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{service: service}
}

func (ctrl *UserController) GetUsers(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.GetUsersForm)
	data := ctrl.service.GetUser().GetUsers(form.Order, form.Filter, form.Limit, form.Offset)
	return c.JSON(data)
}

func (ctrl *UserController) CreateUsers(c *fiber.Ctx) error {
	_, errs := c.FormFile("avatar")
	avatar := ""
	if errs == nil {
		maxFileSize := int64(2 * 1024 * 1024) // 2 MB
		allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

		err := middleware.FileUploadToS3Middleware(c, "avatar", "avatars", allowedTypes, maxFileSize)
		if err != nil {
			fmt.Println("err")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "File upload failed: " + err.Error(),
				"data":    nil,
			})
		}
		avatar = c.Locals("fileName").(string)
	}

	form := c.Locals("validatedForm").(*dto.CreateUsersForm)
	data := ctrl.service.GetUser().CreateUser(c.UserContext(), form.Email, form.Id_Person, form.Jenis_User, form.Nama_Lengkap, form.No_Hp, form.Username, form.Password, avatar)
	return c.JSON(data)
}

func (ctrl *UserController) GetUser(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetUserParams)
	data := ctrl.service.GetUser().GetUser(c.UserContext(), params.Id_User)
	return c.JSON(data)
}

func (ctrl *UserController) UpdateUsers(c *fiber.Ctx) error {
	_, errs := c.FormFile("avatar")
	avatar := ""
	if errs == nil {
		maxFileSize := int64(2 * 1024 * 1024) // 2 MB
		allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

		err := middleware.FileUploadToS3Middleware(c, "avatar", "avatars", allowedTypes, maxFileSize)
		if err != nil {
			fmt.Println("File upload error:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "File upload failed: " + err.Error(),
				"data":    nil,
			})
		}
		avatar = c.Locals("fileName").(string)
	}

	params := c.Locals("validatedParams").(*dto.UpdateUsersParams)
	form := c.Locals("validatedForm").(*dto.UpdateUsersForm)
	data := ctrl.service.GetUser().UpdateUser(c.UserContext(), params.Id_User, avatar, form.Email, form.Id_Person, form.Jenis_User, form.Nama_Lengkap, form.No_Hp, form.Username)
	return c.JSON(data)
}

func (ctrl *UserController) UpdatePassword(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.UpdatePasswordParams)
	form := c.Locals("validatedForm").(*dto.UpdatePasswordForm)
	data := ctrl.service.GetUser().UpdatePassword(c.UserContext(), params.Id_User, form.Password)
	return c.JSON(data)
}

func (ctrl *UserController) DeleteUser(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.DeleteUserParams)
	data := ctrl.service.GetUser().DeleteUser(c.UserContext(), params.Id_User)
	return c.JSON(data)
}

func (ctrl *UserController) GetUsersDosen(c *fiber.Ctx) error {
	data := ctrl.service.GetUser().GetUsersDosen(c.UserContext())
	return c.JSON(data)
}

func (ctrl *UserController) GetDetailDosen(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetIDUserParams)
	data := ctrl.service.GetUser().GetDetailDosen(c.UserContext(), params.Person_ID)
	return c.JSON(data)
}

func (ctrl *UserController) BulkCreateUsersMahasiswa(c *fiber.Ctx) error {
	var items []dto.BulkMahasiswaItem
	if err := c.BodyParser(&items); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid: " + err.Error(),
			"data":    nil,
		})
	}
	if len(items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data mahasiswa tidak boleh kosong",
			"data":    nil,
		})
	}
	data := ctrl.service.GetUser().BulkCreateUsersMahasiswa(c.UserContext(), items)
	return c.JSON(data)
}

func (ctrl *UserController) GetDetailMahasiswa(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetMahasiswaParams)
	data := ctrl.service.GetUser().GetDetailMahasiswa(c.UserContext(), params.ID_Registrasi_Mahasiswa)
	return c.JSON(data)
}

func (ctrl *UserController) GetUsersMahasiswa(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetProdiParams)
	data := ctrl.service.GetUser().GetUsersMahasiswa(c.UserContext(), params.ID_Prodi)
	return c.JSON(data)
}

func (ctrl *UserController) GetUsersMahasiswaData(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*dto.GetProdiParams)
	idProdi := params.ID_Prodi

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	order := c.Query("order", "nama_mahasiswa asc")
	filter := c.Query("filter", "")
	idAngkatan := c.Query("id_angkatan", "")

	data := ctrl.service.GetUser().GetUsersMahasiswaData(c.UserContext(), idProdi, idAngkatan, order, filter, limit, offset)
	return c.JSON(data)
}

func (ctrl *UserController) GenerateUserMahasiswa(c *fiber.Ctx) error {
	data := ctrl.service.GetUser().GenerateUserMahasiswa(c.UserContext())
	return c.JSON(data)
}

func (ctrl *UserController) GenerateUserDosen(c *fiber.Ctx) error {
	data := ctrl.service.GetUser().GenerateUserDosen(c.UserContext())
	return c.JSON(data)
}

func (ctrl *UserController) GenerateUserTendik(c *fiber.Ctx) error {
	data := ctrl.service.GetUser().GenerateUserTendik(c.UserContext())
	return c.JSON(data)
}
