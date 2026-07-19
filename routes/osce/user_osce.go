// Package routes (osce) mendaftarkan endpoint HTTP untuk domain pengguna OSCE.
package routes

import (
	"sso-service/controllers"
	"sso-service/domain/dto/osce"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// OsceUserRoute mewadahi controller registry & router Fiber.
type OsceUserRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IOsceUserRoute adalah kontrak pendaftaran rute pengguna OSCE.
type IOsceUserRoute interface {
	Run()
}

// NewOsceUserRoute membuat instance baru.
func NewOsceUserRoute(controller controllers.IControllerRegistry, router fiber.Router) IOsceUserRoute {
	return &OsceUserRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain pengguna OSCE.
func (r *OsceUserRoute) Run() {
	api := r.router.Group("/osce")

	// Penguji
	api.Get("/penguji", middleware.ValidateForm(&osce.PengujiQuery{}), r.controller.GetOsceUser().GetPenguji)
	api.Get("/penguji/:id", r.controller.GetOsceUser().GetPengujiByID)
	api.Post("/penguji", middleware.ValidateForm(&osce.PengujiForm{}), r.controller.GetOsceUser().CreatePenguji)
	api.Put("/penguji/:id", middleware.ValidateForm(&osce.PengujiForm{}), r.controller.GetOsceUser().UpdatePenguji)
	api.Delete("/penguji/:id", r.controller.GetOsceUser().DeletePenguji)

	// Kelompok Mahasiswa
	api.Get("/kelompok-mahasiswa", middleware.ValidateForm(&osce.KelompokMahasiswaQuery{}), r.controller.GetOsceUser().GetKelompokMahasiswa)
	api.Get("/kelompok-mahasiswa/:id", r.controller.GetOsceUser().GetKelompokMahasiswaByID)
	api.Post("/kelompok-mahasiswa", middleware.ValidateForm(&osce.KelompokMahasiswaForm{}), r.controller.GetOsceUser().CreateKelompokMahasiswa)
	api.Put("/kelompok-mahasiswa/:id", middleware.ValidateForm(&osce.KelompokMahasiswaForm{}), r.controller.GetOsceUser().UpdateKelompokMahasiswa)
	api.Delete("/kelompok-mahasiswa/:id", r.controller.GetOsceUser().DeleteKelompokMahasiswa)

	// Anggota Kelompok
	api.Get("/kelompok-mahasiswa/:id/anggota", r.controller.GetOsceUser().GetAnggotaByKelompok)
	api.Post("/kelompok-mahasiswa/:id/anggota", middleware.ValidateForm(&osce.AnggotaKelompokForm{}), r.controller.GetOsceUser().TambahAnggota)
	api.Post("/kelompok-mahasiswa/:id/anggota/bulk", middleware.ValidateForm(&osce.BulkAnggotaKelompokForm{}), r.controller.GetOsceUser().TambahAnggotaBulk)
	api.Delete("/anggota-kelompok/:id", r.controller.GetOsceUser().HapusAnggota)

	// Penugasan Penguji
	api.Get("/sesi/:id/penugasan-penguji", r.controller.GetOsceUser().GetPenugasanBySesi)
	api.Post("/sesi/:id/penugasan-penguji", middleware.ValidateForm(&osce.PenugasanPengujiForm{}), r.controller.GetOsceUser().AssignPenguji)
	api.Delete("/penugasan-penguji/:id", r.controller.GetOsceUser().HapusPenugasan)
}
