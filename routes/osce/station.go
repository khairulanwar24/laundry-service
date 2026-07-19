// Package routes (osce) mendaftarkan endpoint HTTP untuk domain OSCE station.
package routes

import (
	"sso-service/controllers"
	"sso-service/domain/dto/osce"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// OsceStationRoute mewadahi controller registry & router Fiber.
type OsceStationRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IOsceStationRoute adalah kontrak pendaftaran rute station OSCE.
type IOsceStationRoute interface {
	Run()
}

// NewOsceStationRoute membuat instance baru.
func NewOsceStationRoute(controller controllers.IControllerRegistry, router fiber.Router) IOsceStationRoute {
	return &OsceStationRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain station OSCE.
func (r *OsceStationRoute) Run() {
	api := r.router.Group("/osce")

	// Tipe Station
	api.Get("/tipe-station", r.controller.GetOsceStation().GetTipeStation)
	api.Get("/tipe-station/:id", r.controller.GetOsceStation().GetTipeStationByID)
	api.Post("/tipe-station", middleware.ValidateForm(&osce.TipeStationForm{}), r.controller.GetOsceStation().CreateTipeStation)
	api.Put("/tipe-station/:id", middleware.ValidateForm(&osce.TipeStationForm{}), r.controller.GetOsceStation().UpdateTipeStation)
	api.Delete("/tipe-station/:id", r.controller.GetOsceStation().DeleteTipeStation)

	// Master Station
	api.Get("/stations", middleware.ValidateForm(&osce.StationQuery{}), r.controller.GetOsceStation().GetStations)
	api.Get("/stations/:id", r.controller.GetOsceStation().GetStationByID)
	api.Post("/stations", middleware.ValidateForm(&osce.StationForm{}), r.controller.GetOsceStation().CreateStation)
	api.Put("/stations/:id", middleware.ValidateForm(&osce.StationForm{}), r.controller.GetOsceStation().UpdateStation)
	api.Delete("/stations/:id", r.controller.GetOsceStation().DeleteStation)

	// Kompetensi
	api.Get("/stations/:id/kompetensi", r.controller.GetOsceStation().GetKompetensiByStation)
	api.Post("/stations/:id/kompetensi", middleware.ValidateForm(&osce.KompetensiStationForm{}), r.controller.GetOsceStation().CreateKompetensi)
	api.Put("/kompetensi/:id", middleware.ValidateForm(&osce.KompetensiStationForm{}), r.controller.GetOsceStation().UpdateKompetensi)
	api.Delete("/kompetensi/:id", r.controller.GetOsceStation().DeleteKompetensi)

	// Checklist
	api.Get("/stations/:id/checklists", r.controller.GetOsceStation().GetChecklistByStation)
	api.Post("/stations/:id/checklists", middleware.ValidateForm(&osce.ChecklistStationForm{}), r.controller.GetOsceStation().CreateChecklist)
	api.Put("/checklists/:id", middleware.ValidateForm(&osce.ChecklistStationForm{}), r.controller.GetOsceStation().UpdateChecklist)
	api.Delete("/checklists/:id", r.controller.GetOsceStation().DeleteChecklist)
	api.Put("/stations/:id/checklists/reorder", middleware.ValidateForm(&osce.ChecklistReorderRequest{}), r.controller.GetOsceStation().ReorderChecklist)
}
