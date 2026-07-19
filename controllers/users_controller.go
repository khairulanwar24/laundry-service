package controllers

import (
	"fmt"
	middleware "sso-service/middlewares"
	"sso-service/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type GetUsersForm struct {
	Limit  int    `json:"limit" form:"limit" validate:"required,numeric,oneof=10 25 50 100"  default:"10"`
	Offset int    `json:"offset"  form:"offset" validate:"numeric"`
	Order  string `json:"order"  form:"order" validate:""`
	Filter string `json:"filter"  form:"filter" validate:""`
}

func GetUsers(c *fiber.Ctx) error {

	form := c.Locals("validatedForm").(*GetUsersForm)

	limit := form.Limit

	offset := form.Offset

	order := form.Order
	filter := form.Filter

	data := models.GetUsers(limit, offset, order, filter)
	// fmt.Println(data)
	return c.JSON(data)
	// return c.JSON("oke")
}

type CreateUsersForm struct {
	Email        string `json:"email" form:"email" validate:"required,email"`
	Id_Person    string `json:"id_person" form:"id_person" validate:"required,uuid4"`
	Jenis_User   string `json:"jenis_user" form:"jenis_user" validate:"required,oneof='dosen' 'tenaga pendidik' 'mahasiswa' 'orang tua' 'perseptor'"`
	Nama_Lengkap string `json:"nama_lengkap" form:"nama_lengkap" validate:"required"`
	No_Hp        string `json:"no_hp" form:"no_hp" validate:"required,numeric"`
	Username     string `json:"username" form:"username" validate:"required"`
	Password     string `json:"password"  form:"password" validate:"required"`
}

func CreateUsers(c *fiber.Ctx) error {

	_, errs := c.FormFile("avatar")
	avatar := ""
	if errs == nil {
		maxFileSize := int64(2 * 1024 * 1024) // 2 MB
		allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

		err := middleware.FileUploadToS3Middleware(c, "avatar", "avatars", allowedTypes, maxFileSize)
		// err := middleware.FileUploadMiddleware(c, "avatar", "avatar", allowedTypes, maxFileSize)
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

	form := c.Locals("validatedForm").(*CreateUsersForm)

	email := form.Email
	id_person := form.Id_Person
	jenis_user := form.Jenis_User
	nama_lengkap := form.Nama_Lengkap
	no_hp := form.No_Hp
	username := form.Username
	password := form.Password

	data := models.CreateUsers(email, id_person, jenis_user, nama_lengkap, no_hp, username, password, avatar)

	return c.JSON(data)
}

type GetUserParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

func GetUser(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*GetUserParams)

	id_user := params.Id_User

	data := models.GetUser(id_user)

	return c.JSON(data)
}

type UpdateUsersParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

type UpdateUsersForm struct {
	Email        string `json:"email" form:"email" validate:"required,email"`
	Id_Person    string `json:"id_person" form:"id_person" validate:"required,uuid4"`
	Jenis_User   string `json:"jenis_user" form:"jenis_user" validate:"required,oneof='dosen' 'tenaga tendidik' 'mahasiswa' 'orang tua' 'perseptor'"`
	Nama_Lengkap string `json:"nama_lengkap" form:"nama_lengkap" validate:"required"`
	No_Hp        string `json:"no_hp" form:"no_hp" validate:"required,numeric"`
	Username     string `json:"username" form:"username" validate:"required"`
}

func UpdateUsers(c *fiber.Ctx) error {

	_, errs := c.FormFile("avatar")
	avatar := ""
	if errs == nil {
		maxFileSize := int64(2 * 1024 * 1024) // 2 MB
		allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

		// err := middleware.FileUploadMiddleware(c, "avatar", "avatar", allowedTypes, maxFileSize)
		// // return err
		// // panic(err)
		// if err != nil {
		// 	fmt.Println("err")
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"success": false,
		// 		"message": "File upload failed: " + err.Error(),
		// 		"data":    nil,
		// 	})
		// }
		// avatar = c.Locals("fileName").(string)
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

	params := c.Locals("validatedParams").(*UpdateUsersParams)

	id_user := params.Id_User

	form := c.Locals("validatedForm").(*UpdateUsersForm)

	// avatar := c.FormValue("avatar")
	email := form.Email
	id_person := form.Id_Person
	jenis_user := form.Jenis_User
	nama_lengkap := form.Nama_Lengkap
	no_hp := form.No_Hp
	username := form.Username

	// middleware.FileUploadMiddleware("avatar", "avatar")

	// avatar := ""

	// fmt.Println(avatar)
	data := models.UpdateUsers(id_user, avatar, email, id_person, jenis_user, nama_lengkap, no_hp, username)

	return c.JSON(data)
}

type UpdatePasswordParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

type UpdatePasswordForm struct {
	Password string `json:"password" form:"password" validate:"required"`
}

func UpdatePassword(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*UpdatePasswordParams)

	id_user := params.Id_User

	form := c.Locals("validatedForm").(*UpdatePasswordForm)

	password := form.Password

	data := models.UpdatePassword(id_user, password)

	return c.JSON(data)
}

type DeleteUserParams struct {
	Id_User string `json:"id_user" form:"id_user" validate:"required,uuid4"`
}

func DeleteUser(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*DeleteUserParams)
	id_user := params.Id_User

	data := models.DeleteUser(id_user)

	return c.JSON(data)
}

func GetUsersDosen(c *fiber.Ctx) error {
	// fmt.Println("masuk")
	data := models.GetUsersDosen()
	return c.JSON(data)
}

type GetIDUserParams struct {
	Person_ID string `json:"person_id" form:"person_id" validate:"required,uuid4"`
}

func GetDetailDosen(c *fiber.Ctx) error {

	params := c.Locals("validatedParams").(*GetIDUserParams)

	person_id := params.Person_ID
	// fmt.Println("masuk")
	data := models.GetDetailDosen(person_id)
	return c.JSON(data)
}

// func getTenagaPendidik(c *fiber.Ctx) error {
// 	data := models.GetUsersTenagaPendidik()
// 	return c.JSON(data)
// }

type GetProdiParams struct {
	ID_Prodi string `json:"id_prodi" form:"id_prodi" validate:"required,uuid4"`
}

type GetMahasiswaParams struct {
	ID_Registrasi_Mahasiswa string `json:"id_registrasi_mahasiswa" form:"id_registrasi_mahasiswa" validate:"required,uuid4"`
}

func BulkCreateUsersMahasiswa(c *fiber.Ctx) error {
	var items []models.BulkMahasiswaItem
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
	data := models.BulkCreateUsersMahasiswa(items)
	return c.JSON(data)
}

func GetDetailMahasiswa(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetMahasiswaParams)
	data := models.GetDetailMahasiswa(params.ID_Registrasi_Mahasiswa)
	return c.JSON(data)
}

func GetUsersMahasiswa(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetProdiParams)

	id_prodi := params.ID_Prodi

	data := models.GetUsersMahasiswa(id_prodi)
	return c.JSON(data)
}

func GetUsersMahasiswaData(c *fiber.Ctx) error {
	params := c.Locals("validatedParams").(*GetProdiParams)
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

	data := models.GetUsersMahasiswaData(idProdi, idAngkatan, order, filter, limit, offset)
	return c.JSON(data)
}
