// Package routes (groupakses) mendaftarkan endpoint HTTP untuk domain master group akses & group akses.
package routes

import (
	"sso-service/controllers"
	"sso-service/domain/dto"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// GroupAksesRoute mewadahi controller registry & router Fiber.
type GroupAksesRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IGroupAksesRoute adalah kontrak pendaftaran rute group akses.
type IGroupAksesRoute interface {
	Run()
}

// NewGroupAksesRoute membuat instance GroupAksesRoute baru.
func NewGroupAksesRoute(controller controllers.IControllerRegistry, router fiber.Router) IGroupAksesRoute {
	return &GroupAksesRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain master group akses & group akses.
func (r *GroupAksesRoute) Run() {
	mstgroupakses := r.router.Group("/mstgroupakses", middleware.JWTMiddleware)
	mstgroupakses.Get("/:id_master_aplikasi", middleware.ValidatedParams(&dto.GetMstGroupAksesParams{}), middleware.ValidatedParams2(&dto.GetMstGroupAksesForm{}), r.controller.GetGroupAkses().GetMstGroupAkses)
	mstgroupakses.Get("/modul/:id_master_aplikasi/:id_master_group", middleware.ValidatedParams(&dto.GetMstGroupAksesModulParams{}), middleware.ValidatedParams2(&dto.GetMstGroupAksesModulForm{}), r.controller.GetGroupAkses().GetMstGroupAksesModul)
	mstgroupakses.Post("/", middleware.ValidateForm(&dto.CreateMstGroupAksesForm{}), r.controller.GetGroupAkses().CreateMstGroupAkses)
	mstgroupakses.Put("/:id_master_group", middleware.ValidatedParams(&dto.UpdateMstGroupAksesParams{}), middleware.ValidateForm(&dto.UpdateMstGroupAksesForm{}), r.controller.GetGroupAkses().UpdateMstGroupAkses)
	mstgroupakses.Get("/detail/:id_master_group", middleware.ValidatedParams(&dto.GetDetailMstGroupAksesParams{}), r.controller.GetGroupAkses().GetDetailMstGroupAkses)
	mstgroupakses.Delete("/:id_master_group", middleware.ValidatedParams(&dto.DeleteMstGroupAksesParams{}), r.controller.GetGroupAkses().DeleteMstGroupAkses)

	groupakses := r.router.Group("/groupakses", middleware.JWTMiddleware)
	groupakses.Get("/:id_master_group", middleware.ValidatedParams(&dto.GetGroupAksesParams{}), middleware.ValidatedParams2(&dto.GetGroupAksesForm{}), r.controller.GetGroupAkses().GetGroupAkses)
	groupakses.Post("/", middleware.ValidateForm(&dto.CreateGroupAksesForm{}), r.controller.GetGroupAkses().CreateGroupAkses)
	groupakses.Delete("/:id_group_akses", middleware.ValidatedParams(&dto.DeleteGroupAksesParams{}), r.controller.GetGroupAkses().DeleteGroupAkses)
	groupakses.Get("/usermenu/:id_user/:id_master_aplikasi", middleware.ValidatedParams(&dto.GetGroupAksesUserMenuParams{}), r.controller.GetGroupAkses().GetGroupAksesUserMenu)
	groupakses.Get("/userapps/:id_user", middleware.ValidatedParams(&dto.GetGroupAksesUserAppsParams{}), r.controller.GetGroupAkses().GetGroupAksesUserApps)
	groupakses.Post("/userapps/", middleware.ValidateForm(&dto.CreateGroupAksesUserAppsForm{}), r.controller.GetGroupAkses().CreateGroupAksesUserApps)
	groupakses.Post("/userapps/bulk", middleware.ValidateForm(&dto.BulkGroupAksesUserAppsForm{}), r.controller.GetGroupAkses().BulkCreateGroupAksesUserApps)
	groupakses.Put("/userapps/:id_trans_user_group", middleware.ValidatedParams(&dto.UpdateGroupAksesUserAppsParams{}), middleware.ValidateForm(&dto.UpdateGroupAksesUserAppsForm{}), r.controller.GetGroupAkses().UpdateGroupAksesUserApps)
	groupakses.Delete("/userapps/:id_trans_user_group", middleware.ValidatedParams(&dto.DeleteGroupAksesUserAppsParams{}), r.controller.GetGroupAkses().DeleteGroupAksesUserApps)
}
