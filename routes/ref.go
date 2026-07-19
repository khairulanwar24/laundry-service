package routes

import (
	"sso-service/controllers"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRefRoutes(app *fiber.App) {
	ref := app.Group("/ref", middleware.JWTMiddleware)
	ref.Get("/get_prodi", controllers.GetMasterProdi)
	ref.Get("/get_angkatan/:id_prodi", middleware.ValidatedParams(&controllers.GetAngkatanParams{}), controllers.GetAngkatan)
}
