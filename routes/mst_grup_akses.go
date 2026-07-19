package routes

import (
	"sso-service/controllers"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupMstGrupAksesRoutes(app *fiber.App) {

	// middleware.ValidatedParams(&controllers.UpdateUsersParams{}), middleware.ValidateForm(&controllers.UpdateUsersForm{}),

	mstgroupakses := app.Group("/mstgroupakses", middleware.JWTMiddleware)
	mstgroupakses.Get("/:id_master_aplikasi", middleware.ValidatedParams(&controllers.GetMstGroupAksesParams{}), middleware.ValidatedParams2(&controllers.GetMstGroupAksesForm{}), controllers.GetMstGroupAkses)
	mstgroupakses.Get("/modul/:id_master_aplikasi/:id_master_group", middleware.ValidatedParams(&controllers.GetMstGroupAksesModulParams{}), middleware.ValidatedParams2(&controllers.GetMstGroupAksesModulForm{}), controllers.GetMstGroupAksesModul)
	mstgroupakses.Post("/", middleware.ValidateForm(&controllers.CreateMstGroupAksesForm{}), controllers.CreateMstGroupAkses)
	mstgroupakses.Put("/:id_master_group", middleware.ValidatedParams(&controllers.UpdateMstGroupAksesParams{}), middleware.ValidateForm(&controllers.UpdateMstGroupAksesForm{}), controllers.UpdateMstGroupAkses)
	mstgroupakses.Get("/detail/:id_master_group", middleware.ValidatedParams(&controllers.GetDetailMstGroupAksesParams{}), controllers.GetDetailMstGroupAkses)
	mstgroupakses.Delete("/:id_master_group", middleware.ValidatedParams(&controllers.DeleteMstGroupAksesParams{}), controllers.DeleteMstGroupAkses)

	groupakses := app.Group("/groupakses", middleware.JWTMiddleware)

	groupakses.Get("/:id_master_group", middleware.ValidatedParams(&controllers.GetGroupAksesParams{}), middleware.ValidatedParams2(&controllers.GetGroupAksesForm{}), controllers.GetGroupAkses)
	groupakses.Post("/", middleware.ValidateForm(&controllers.CreateGroupAksesForm{}), controllers.CreateGroupAkses)
	groupakses.Delete("/:id_group_akses", middleware.ValidatedParams(&controllers.DeleteGroupAksesParams{}), controllers.DeleteGroupAkses)
	groupakses.Get("/usermenu/:id_user/:id_master_aplikasi", middleware.ValidatedParams(&controllers.GetGroupAksesUserMenuParams{}), controllers.GetGroupAksesUserMenu)
	groupakses.Get("/userapps/:id_user", middleware.ValidatedParams(&controllers.GetGroupAksesUserAppsParams{}), controllers.GetGroupAksesUserApps)
	groupakses.Post("/userapps/", middleware.ValidateForm(&controllers.CreateGroupAksesUserAppsForm{}), controllers.CreateGroupAksesUserApps)
	groupakses.Post("/userapps/bulk", middleware.ValidateForm(&controllers.BulkGroupAksesUserAppsForm{}), controllers.BulkCreateGroupAksesUserApps)

	groupakses.Put("/userapps/:id_trans_user_group", middleware.ValidatedParams(&controllers.UpdateGroupAksesUserAppsParams{}), middleware.ValidateForm(&controllers.UpdateGroupAksesUserAppsForm{}), controllers.UpdateGroupAksesUserApps)
	groupakses.Delete("/userapps/:id_trans_user_group", middleware.ValidatedParams(&controllers.DeleteGroupAksesUserAppsParams{}), controllers.DeleteGroupAksesUserApps)
}
