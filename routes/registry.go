// registry.go adalah pusat perakitan rute berbasis pola berlapis (registry).
package routes

import (
	"laundry-service/controllers"
	accountRoute "laundry-service/routes/account"

	"github.com/gofiber/fiber/v2"
)

// Registry menyimpan controller registry & router Fiber.
type Registry struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IRouteRegistry adalah kontrak untuk menyajikan seluruh rute domain laundry.
type IRouteRegistry interface {
	Serve()
}

// NewRouteRegistry membuat route registry baru.
func NewRouteRegistry(controller controllers.IControllerRegistry, router fiber.Router) IRouteRegistry {
	return &Registry{controller: controller, router: router}
}

// Serve mendaftarkan seluruh rute domain laundry.
func (r *Registry) Serve() {
	r.accountRoute().Run()
}

func (r *Registry) accountRoute() accountRoute.IAccountRoute {
	return accountRoute.NewAccountRoute(r.controller, r.router)
}
