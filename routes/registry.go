// registry.go adalah pusat perakitan rute berbasis pola berlapis (registry).
// Berdampingan dengan SetupRoutes lama selama migrasi bertahap; domain yang
// sudah dimigrasi didaftarkan di sini, sisanya masih lewat SetupRoutes.
package routes

import (
	"sso-service/controllers"
	groupAksesRoute "sso-service/routes/groupakses"
	masterAppRoute "sso-service/routes/masterapp"
	mstMenuRoute "sso-service/routes/mstmenu"
	refRoute "sso-service/routes/ref"
	userRoute "sso-service/routes/user"

	"github.com/gofiber/fiber/v2"
)

// Registry menyimpan controller registry & router Fiber.
type Registry struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IRouteRegistry adalah kontrak untuk menyajikan seluruh rute yang telah dimigrasi.
type IRouteRegistry interface {
	Serve()
}

// NewRouteRegistry membuat route registry baru.
func NewRouteRegistry(controller controllers.IControllerRegistry, router fiber.Router) IRouteRegistry {
	return &Registry{controller: controller, router: router}
}

// Serve mendaftarkan seluruh rute domain yang telah dimigrasi ke pola berlapis.
func (r *Registry) Serve() {
	r.refRoute().Run()
	r.userRoute().Run()
	r.masterAppRoute().Run()
	r.mstMenuRoute().Run()
	r.groupAksesRoute().Run()
}

func (r *Registry) refRoute() refRoute.IRefRoute {
	return refRoute.NewRefRoute(r.controller, r.router)
}

func (r *Registry) userRoute() userRoute.IUserRoute {
	return userRoute.NewUserRoute(r.controller, r.router)
}

func (r *Registry) masterAppRoute() masterAppRoute.IMasterAppRoute {
	return masterAppRoute.NewMasterAppRoute(r.controller, r.router)
}

func (r *Registry) mstMenuRoute() mstMenuRoute.IMstMenuRoute {
	return mstMenuRoute.NewMstMenuRoute(r.controller, r.router)
}

func (r *Registry) groupAksesRoute() groupAksesRoute.IGroupAksesRoute {
	return groupAksesRoute.NewGroupAksesRoute(r.controller, r.router)
}
