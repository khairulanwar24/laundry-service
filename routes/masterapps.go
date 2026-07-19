package routes

import (
	"sso-service/controllers"
	middleware "sso-service/middlewares"
	"sso-service/types"

	"github.com/gofiber/fiber/v2"
)

func SetupMasterRoutes(app *fiber.App) {
	masterapps := app.Group("/masterapps", middleware.JWTMiddleware)
	masterapps.Get("/", middleware.ValidatedParams2(&types.GetData{}), controllers.GetMasterApps)
	masterapps.Get("/:id_master_aplikasi", middleware.ValidatedParams(&controllers.GetMasterAppByIdParams{}), controllers.GetMasterAppById)
	masterapps.Post("/", middleware.ValidateForm(&controllers.CreateMasterAppForm{}), controllers.CreateMasterApp)
	masterapps.Put("/:id_master_aplikasi", middleware.ValidatedParams(&controllers.UpdateMasterAppsParams{}), middleware.ValidateForm(&controllers.UpdateMasterAppsForm{}), controllers.UpdateMasterApps)
	masterapps.Delete("/:id_master_aplikasi", middleware.ValidatedParams(&controllers.DeleteMasterAppParams{}), controllers.DeleteMasterApp)
}
