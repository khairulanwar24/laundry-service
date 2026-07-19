package routes

import (
	"sso-service/controllers"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	users := app.Group("/users", middleware.JWTMiddleware)
	users.Get("/", middleware.ValidatedParams2(&controllers.GetUsersForm{}), controllers.GetUsers)
	users.Post("/", middleware.ValidateForm(&controllers.CreateUsersForm{}), controllers.CreateUsers)
	users.Get("/get_mahasiswa/:id_prodi", middleware.ValidatedParams(&controllers.GetProdiParams{}), controllers.GetUsersMahasiswa)
	users.Get("/get_mahasiswa_data/:id_prodi", middleware.ValidatedParams(&controllers.GetProdiParams{}), controllers.GetUsersMahasiswaData)
	users.Get("/get_dosen", controllers.GetUsersDosen)
	users.Get("/get_detail_dosen/:person_id", middleware.ValidatedParams(&controllers.GetIDUserParams{}), controllers.GetDetailDosen)
	users.Get("/get_detail_mahasiswa/:id_registrasi_mahasiswa", middleware.ValidatedParams(&controllers.GetMahasiswaParams{}), controllers.GetDetailMahasiswa)
	users.Get("/:id_user/", middleware.ValidatedParams(&controllers.GetUserParams{}), controllers.GetUser)
	users.Put("/:id_user", middleware.ValidatedParams(&controllers.UpdateUsersParams{}), middleware.ValidateForm(&controllers.UpdateUsersForm{}), controllers.UpdateUsers)
	users.Put("/updatepassword/:id_user", middleware.ValidatedParams(&controllers.UpdatePasswordParams{}), middleware.ValidateForm(&controllers.UpdatePasswordForm{}), controllers.UpdatePassword)
	users.Delete("/:id_user", middleware.ValidatedParams(&controllers.DeleteUserParams{}), controllers.DeleteUser)
	users.Post("/bulk_mahasiswa", controllers.BulkCreateUsersMahasiswa)
	users.Post("/generate", controllers.GenerateUserMahasiswa)
	users.Post("/generate_dosen", controllers.GenerateUserDosen)
	users.Post("/generate_tendik", controllers.GenerateUserTendik)

	// users.Get("/getTendik", controllers.GetUsersTendik)
}
