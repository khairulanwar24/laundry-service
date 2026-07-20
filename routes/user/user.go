// Package routes (user) mendaftarkan endpoint HTTP untuk domain user.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// UserRoute mewadahi controller registry & router Fiber.
type UserRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IUserRoute adalah kontrak pendaftaran rute user.
type IUserRoute interface {
	Run()
}

// NewUserRoute membuat instance UserRoute baru.
func NewUserRoute(controller controllers.IControllerRegistry, router fiber.Router) IUserRoute {
	return &UserRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain user.
func (r *UserRoute) Run() {
	users := r.router.Group("/users", middleware.JWTMiddleware)
	users.Get("/", middleware.ValidatedParams2(&dto.GetUsersForm{}), r.controller.GetUser().GetUsers)
	users.Post("/", middleware.ValidateForm(&dto.CreateUsersForm{}), r.controller.GetUser().CreateUsers)
	users.Get("/get_mahasiswa/:id_prodi", middleware.ValidatedParams(&dto.GetProdiParams{}), r.controller.GetUser().GetUsersMahasiswa)
	users.Get("/get_mahasiswa_data/:id_prodi", middleware.ValidatedParams(&dto.GetProdiParams{}), r.controller.GetUser().GetUsersMahasiswaData)
	users.Get("/get_dosen", r.controller.GetUser().GetUsersDosen)
	users.Get("/get_detail_dosen/:person_id", middleware.ValidatedParams(&dto.GetIDUserParams{}), r.controller.GetUser().GetDetailDosen)
	users.Get("/get_detail_mahasiswa/:id_registrasi_mahasiswa", middleware.ValidatedParams(&dto.GetMahasiswaParams{}), r.controller.GetUser().GetDetailMahasiswa)
	users.Get("/:id_user/", middleware.ValidatedParams(&dto.GetUserParams{}), r.controller.GetUser().GetUser)
	users.Put("/:id_user", middleware.ValidatedParams(&dto.UpdateUsersParams{}), middleware.ValidateForm(&dto.UpdateUsersForm{}), r.controller.GetUser().UpdateUsers)
	users.Put("/updatepassword/:id_user", middleware.ValidatedParams(&dto.UpdatePasswordParams{}), middleware.ValidateForm(&dto.UpdatePasswordForm{}), r.controller.GetUser().UpdatePassword)
	users.Delete("/:id_user", middleware.ValidatedParams(&dto.DeleteUserParams{}), r.controller.GetUser().DeleteUser)
	users.Post("/bulk_mahasiswa", r.controller.GetUser().BulkCreateUsersMahasiswa)
	users.Post("/generate", r.controller.GetUser().GenerateUserMahasiswa)
	users.Post("/generate_dosen", r.controller.GetUser().GenerateUserDosen)
	users.Post("/generate_tendik", r.controller.GetUser().GenerateUserTendik)
}
