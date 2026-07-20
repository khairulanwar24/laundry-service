// Package routes (osce) mendaftarkan endpoint HTTP untuk domain OSCE penilaian.
package routes

import (
	"laundry-service/controllers"
	"laundry-service/domain/dto/osce"
	middleware "laundry-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// OsceAssessmentRoute mewadahi controller registry & router Fiber.
type OsceAssessmentRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IOsceAssessmentRoute adalah kontrak pendaftaran rute penilaian OSCE.
type IOsceAssessmentRoute interface {
	Run()
}

// NewOsceAssessmentRoute membuat instance baru.
func NewOsceAssessmentRoute(controller controllers.IControllerRegistry, router fiber.Router) IOsceAssessmentRoute {
	return &OsceAssessmentRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain penilaian OSCE.
func (r *OsceAssessmentRoute) Run() {
	api := r.router.Group("/osce")

	// Penilaian (Penguji)
	api.Get("/penguji/jadwal", r.controller.GetOsceAssessment().GetJadwalPenguji)
	api.Get("/penguji/jadwal/:id/penilaian", r.controller.GetOsceAssessment().GetFormPenilaian)
	api.Post("/penilaian", middleware.ValidateForm(&osce.PenilaianForm{}), r.controller.GetOsceAssessment().CreatePenilaian)
	api.Put("/penilaian/:id", middleware.ValidateForm(&osce.PenilaianUpdateForm{}), r.controller.GetOsceAssessment().UpdatePenilaian)
	api.Put("/penilaian/:id/kirim", r.controller.GetOsceAssessment().KirimPenilaian)
	api.Get("/penilaian/:id", r.controller.GetOsceAssessment().GetPenilaianByID)

	// Hasil Ujian
	api.Get("/ujian/:id/hasil", r.controller.GetOsceAssessment().GetHasilUjian)
	api.Post("/ujian/:id/hasil/kalkulasi", middleware.ValidateForm(&osce.KalkulasiHasilRequest{}), r.controller.GetOsceAssessment().KalkulasiHasil)
	api.Post("/ujian/:id/borderline", middleware.ValidateForm(&osce.BorderlineRequest{}), r.controller.GetOsceAssessment().KalkulasiBorderline)
	api.Get("/ujian/:id/statistik", r.controller.GetOsceAssessment().GetStatistikUjian)
	api.Get("/hasil/saya", r.controller.GetOsceAssessment().GetHasilSaya)
	api.Get("/hasil/:id", r.controller.GetOsceAssessment().GetHasilByID)
}
