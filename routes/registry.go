// registry.go adalah pusat perakitan rute berbasis pola berlapis (registry).
package routes

import (
	"laundry-service/controllers"
	accountRoute "laundry-service/routes/account"
	catalogRoute "laundry-service/routes/catalog"
	outletRoute "laundry-service/routes/outlet"
	paymentMethodRoute "laundry-service/routes/paymentmethod"

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
	r.outletRoute().Run()
	r.paymentMethodRoute().Run()
	r.catalogRoute().Run()
}

func (r *Registry) accountRoute() accountRoute.IAccountRoute {
	return accountRoute.NewAccountRoute(r.controller, r.router)
}

func (r *Registry) outletRoute() outletRoute.IOutletRoute {
	return outletRoute.NewOutletRoute(r.controller, r.router)
}

func (r *Registry) paymentMethodRoute() paymentMethodRoute.IPaymentMethodRoute {
	return paymentMethodRoute.NewPaymentMethodRoute(r.controller, r.router)
}

func (r *Registry) catalogRoute() catalogRoute.ICatalogRoute {
	return catalogRoute.NewCatalogRoute(r.controller, r.router)
}
