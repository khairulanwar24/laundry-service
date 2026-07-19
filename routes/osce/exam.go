// Package routes (osce) mendaftarkan endpoint HTTP untuk domain OSCE ujian.
package routes

import (
	"sso-service/controllers"
	"sso-service/domain/dto/osce"
	middleware "sso-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

// OsceExamRoute mewadahi controller registry & router Fiber.
type OsceExamRoute struct {
	controller controllers.IControllerRegistry
	router     fiber.Router
}

// IOsceExamRoute adalah kontrak pendaftaran rute ujian OSCE.
type IOsceExamRoute interface {
	Run()
}

// NewOsceExamRoute membuat instance baru.
func NewOsceExamRoute(controller controllers.IControllerRegistry, router fiber.Router) IOsceExamRoute {
	return &OsceExamRoute{controller: controller, router: router}
}

// Run mendaftarkan seluruh endpoint domain ujian OSCE.
func (r *OsceExamRoute) Run() {
	api := r.router.Group("/osce")

	// Tahun Akademik
	api.Get("/tahun-akademik", r.controller.GetOsceExam().GetTahunAkademik)
	api.Get("/tahun-akademik/:id", r.controller.GetOsceExam().GetTahunAkademikByID)
	api.Post("/tahun-akademik", middleware.ValidateForm(&osce.TahunAkademikForm{}), r.controller.GetOsceExam().CreateTahunAkademik)
	api.Put("/tahun-akademik/:id", middleware.ValidateForm(&osce.TahunAkademikForm{}), r.controller.GetOsceExam().UpdateTahunAkademik)
	api.Delete("/tahun-akademik/:id", r.controller.GetOsceExam().DeleteTahunAkademik)

	// Semester
	api.Get("/semester", r.controller.GetOsceExam().GetSemester)
	api.Get("/semester/:id", r.controller.GetOsceExam().GetSemesterByID)
	api.Post("/semester", middleware.ValidateForm(&osce.SemesterForm{}), r.controller.GetOsceExam().CreateSemester)
	api.Put("/semester/:id", middleware.ValidateForm(&osce.SemesterForm{}), r.controller.GetOsceExam().UpdateSemester)
	api.Delete("/semester/:id", r.controller.GetOsceExam().DeleteSemester)

	// Program Studi
	api.Get("/program-studi", r.controller.GetOsceExam().GetProgramStudi)
	api.Get("/program-studi/:id", r.controller.GetOsceExam().GetProgramStudiByID)
	api.Post("/program-studi", middleware.ValidateForm(&osce.ProgramStudiForm{}), r.controller.GetOsceExam().CreateProgramStudi)
	api.Put("/program-studi/:id", middleware.ValidateForm(&osce.ProgramStudiForm{}), r.controller.GetOsceExam().UpdateProgramStudi)
	api.Delete("/program-studi/:id", r.controller.GetOsceExam().DeleteProgramStudi)

	// Ujian
	api.Get("/ujian", middleware.ValidateForm(&osce.UjianQuery{}), r.controller.GetOsceExam().GetUjian)
	api.Get("/ujian/:id", r.controller.GetOsceExam().GetUjianByID)
	api.Post("/ujian", middleware.ValidateForm(&osce.UjianForm{}), r.controller.GetOsceExam().CreateUjian)
	api.Put("/ujian/:id", middleware.ValidateForm(&osce.UjianForm{}), r.controller.GetOsceExam().UpdateUjian)
	api.Delete("/ujian/:id", r.controller.GetOsceExam().DeleteUjian)

	// Ujian Station
	api.Get("/ujian/:id/stations", r.controller.GetOsceExam().GetUjianStationByUjian)
	api.Post("/ujian/:id/stations", middleware.ValidateForm(&osce.UjianStationForm{}), r.controller.GetOsceExam().CreateUjianStation)
	api.Put("/ujian-station/:id", middleware.ValidateForm(&osce.UjianStationForm{}), r.controller.GetOsceExam().UpdateUjianStation)
	api.Delete("/ujian-station/:id", r.controller.GetOsceExam().DeleteUjianStation)

	// Sesi Ujian
	api.Get("/ujian/:id/sesi", r.controller.GetOsceExam().GetSesiByUjian)
	api.Post("/ujian/:id/sesi", middleware.ValidateForm(&osce.SesiUjianForm{}), r.controller.GetOsceExam().CreateSesi)
	api.Put("/sesi/:id", middleware.ValidateForm(&osce.SesiUjianForm{}), r.controller.GetOsceExam().UpdateSesi)
	api.Delete("/sesi/:id", r.controller.GetOsceExam().DeleteSesi)

	// Rotasi
	api.Get("/sesi/:id/rotasi", r.controller.GetOsceExam().GetRotasiBySesi)
	api.Post("/sesi/:id/rotasi", middleware.ValidateForm(&osce.RotasiUjianForm{}), r.controller.GetOsceExam().CreateRotasi)
	api.Put("/rotasi/:id", middleware.ValidateForm(&osce.RotasiUjianForm{}), r.controller.GetOsceExam().UpdateRotasi)
	api.Delete("/rotasi/:id", r.controller.GetOsceExam().DeleteRotasi)

	// Jadwal
	api.Get("/sesi/:id/jadwal", r.controller.GetOsceExam().GetJadwalBySesi)
	api.Post("/sesi/:id/jadwal/generate", middleware.ValidateForm(&osce.GenerateJadwalRequest{}), r.controller.GetOsceExam().GenerateJadwal)
	api.Post("/sesi/:id/jadwal", middleware.ValidateForm(&osce.JadwalUjianForm{}), r.controller.GetOsceExam().CreateJadwal)
	api.Put("/jadwal/:id", middleware.ValidateForm(&osce.JadwalUjianForm{}), r.controller.GetOsceExam().UpdateJadwal)
	api.Delete("/jadwal/:id", r.controller.GetOsceExam().DeleteJadwal)
}
