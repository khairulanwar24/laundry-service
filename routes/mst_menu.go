package routes

import (
	"sso-service/controllers"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupMstMenuRoutes(app *fiber.App) {
	mstmenu := app.Group("/mstmenu", middleware.JWTMiddleware)

	mstmenu.Get("/:id_master_aplikasi",
		middleware.ValidatedParams(&controllers.GetMstMenuParams{}),
		middleware.ValidatedParams2(&controllers.GetMstMenuForm{}),
		controllers.GetMstMenu)
	mstmenu.Get("/detail/:id_master_menu",
		middleware.ValidatedParams(&controllers.GetDetailMstMenuParams{}),
		controllers.GetDetailMstMenu)
	mstmenu.Get("/modul/:id_master_menu",
		middleware.ValidatedParams(&controllers.GetMstMenuModulParams{}),
		middleware.ValidatedParams2(&controllers.GetMstMenuModulForm{}),
		controllers.GetMstMenuModul)
	mstmenu.Get("/modul/detail/:id_master_modul",
		middleware.ValidatedParams(&controllers.GetDetailMstMenuModulParams{}),
		controllers.GetDetailMstMenuModul)
	mstmenu.Post("/",
		middleware.ValidateForm(&controllers.CreateMstMenuForm{}),
		controllers.CreateMstMenu)
	mstmenu.Post("/modul",
		middleware.ValidateForm(&controllers.CreateMstMenuModulForm{}),
		controllers.CreateMstMenuModul)
	mstmenu.Put("/:id_master_menu",
		middleware.ValidatedParams(&controllers.UpdateMstMenuParams{}),
		middleware.ValidateForm(&controllers.UpdateMstMenuForm{}),
		controllers.UpdateMstMenu)
	mstmenu.Delete("/:id_master_menu",
		middleware.ValidatedParams(&controllers.DeleteMstMenuParams{}),
		controllers.DeleteMstMenu)
	mstmenu.Put("/modul/:id_master_modul",
		middleware.ValidatedParams(&controllers.UpdateMstMenuModulParams{}),
		middleware.ValidateForm(&controllers.UpdateMstMenuModulForm{}),
		controllers.UpdateMstMenuModul)
	mstmenu.Delete("/modul/:id_master_modul",
		middleware.ValidatedParams(&controllers.DeleteMstMenuModulParams{}),
		controllers.DeleteMstMenuModul)

}
