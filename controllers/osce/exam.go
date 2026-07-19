// Package osce adalah lapisan HTTP handler (Fiber) untuk domain OSCE ujian.
package osce

import (
	"sso-service/domain/dto/osce"
	"sso-service/services"

	"github.com/gofiber/fiber/v2"
)

// OsceExamController membungkus service registry.
type OsceExamController struct {
	service services.IServiceRegistry
}

// IOsceExamController adalah kontrak handler HTTP domain OSCE ujian.
type IOsceExamController interface {
	// Tahun Akademik
	GetTahunAkademik(*fiber.Ctx) error
	GetTahunAkademikByID(*fiber.Ctx) error
	CreateTahunAkademik(*fiber.Ctx) error
	UpdateTahunAkademik(*fiber.Ctx) error
	DeleteTahunAkademik(*fiber.Ctx) error

	// Semester
	GetSemester(*fiber.Ctx) error
	GetSemesterByID(*fiber.Ctx) error
	CreateSemester(*fiber.Ctx) error
	UpdateSemester(*fiber.Ctx) error
	DeleteSemester(*fiber.Ctx) error

	// Program Studi
	GetProgramStudi(*fiber.Ctx) error
	GetProgramStudiByID(*fiber.Ctx) error
	CreateProgramStudi(*fiber.Ctx) error
	UpdateProgramStudi(*fiber.Ctx) error
	DeleteProgramStudi(*fiber.Ctx) error

	// Ujian
	GetUjian(*fiber.Ctx) error
	GetUjianByID(*fiber.Ctx) error
	CreateUjian(*fiber.Ctx) error
	UpdateUjian(*fiber.Ctx) error
	DeleteUjian(*fiber.Ctx) error

	// Ujian Station
	GetUjianStationByUjian(*fiber.Ctx) error
	CreateUjianStation(*fiber.Ctx) error
	UpdateUjianStation(*fiber.Ctx) error
	DeleteUjianStation(*fiber.Ctx) error

	// Sesi
	GetSesiByUjian(*fiber.Ctx) error
	CreateSesi(*fiber.Ctx) error
	UpdateSesi(*fiber.Ctx) error
	DeleteSesi(*fiber.Ctx) error

	// Rotasi
	GetRotasiBySesi(*fiber.Ctx) error
	CreateRotasi(*fiber.Ctx) error
	UpdateRotasi(*fiber.Ctx) error
	DeleteRotasi(*fiber.Ctx) error

	// Jadwal
	GetJadwalBySesi(*fiber.Ctx) error
	GenerateJadwal(*fiber.Ctx) error
	CreateJadwal(*fiber.Ctx) error
	UpdateJadwal(*fiber.Ctx) error
	DeleteJadwal(*fiber.Ctx) error
}

// NewOsceExamController membuat instance baru.
func NewOsceExamController(service services.IServiceRegistry) IOsceExamController {
	return &OsceExamController{service: service}
}

// ===== TAHUN AKADEMIK =====

func (ctrl *OsceExamController) GetTahunAkademik(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetTahunAkademik(c.UserContext()))
}

func (ctrl *OsceExamController) GetTahunAkademikByID(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetTahunAkademikByID(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceExamController) CreateTahunAkademik(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.TahunAkademikForm)
	return c.JSON(ctrl.service.GetOsceExam().CreateTahunAkademik(c.UserContext(), *form))
}

func (ctrl *OsceExamController) UpdateTahunAkademik(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.TahunAkademikForm)
	return c.JSON(ctrl.service.GetOsceExam().UpdateTahunAkademik(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceExamController) DeleteTahunAkademik(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().DeleteTahunAkademik(c.UserContext(), c.Params("id")))
}

// ===== SEMESTER =====

func (ctrl *OsceExamController) GetSemester(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetSemester(c.UserContext()))
}

func (ctrl *OsceExamController) GetSemesterByID(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetSemesterByID(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceExamController) CreateSemester(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.SemesterForm)
	return c.JSON(ctrl.service.GetOsceExam().CreateSemester(c.UserContext(), *form))
}

func (ctrl *OsceExamController) UpdateSemester(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.SemesterForm)
	return c.JSON(ctrl.service.GetOsceExam().UpdateSemester(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceExamController) DeleteSemester(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().DeleteSemester(c.UserContext(), c.Params("id")))
}

// ===== PROGRAM STUDI =====

func (ctrl *OsceExamController) GetProgramStudi(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetProgramStudi(c.UserContext()))
}

func (ctrl *OsceExamController) GetProgramStudiByID(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetProgramStudiByID(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceExamController) CreateProgramStudi(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.ProgramStudiForm)
	return c.JSON(ctrl.service.GetOsceExam().CreateProgramStudi(c.UserContext(), *form))
}

func (ctrl *OsceExamController) UpdateProgramStudi(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.ProgramStudiForm)
	return c.JSON(ctrl.service.GetOsceExam().UpdateProgramStudi(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceExamController) DeleteProgramStudi(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().DeleteProgramStudi(c.UserContext(), c.Params("id")))
}

// ===== UJIAN =====

func (ctrl *OsceExamController) GetUjian(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.UjianQuery)
	return c.JSON(ctrl.service.GetOsceExam().GetUjian(c.UserContext(), *form))
}

func (ctrl *OsceExamController) GetUjianByID(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetUjianByID(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceExamController) CreateUjian(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.UjianForm)
	return c.JSON(ctrl.service.GetOsceExam().CreateUjian(c.UserContext(), *form))
}

func (ctrl *OsceExamController) UpdateUjian(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.UjianForm)
	return c.JSON(ctrl.service.GetOsceExam().UpdateUjian(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceExamController) DeleteUjian(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().DeleteUjian(c.UserContext(), c.Params("id")))
}

// ===== UJIAN STATION =====

func (ctrl *OsceExamController) GetUjianStationByUjian(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetUjianStationByUjian(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceExamController) CreateUjianStation(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.UjianStationForm)
	return c.JSON(ctrl.service.GetOsceExam().CreateUjianStation(c.UserContext(), *form))
}

func (ctrl *OsceExamController) UpdateUjianStation(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.UjianStationForm)
	return c.JSON(ctrl.service.GetOsceExam().UpdateUjianStation(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceExamController) DeleteUjianStation(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().DeleteUjianStation(c.UserContext(), c.Params("id")))
}

// ===== SESI =====

func (ctrl *OsceExamController) GetSesiByUjian(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetSesiByUjian(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceExamController) CreateSesi(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.SesiUjianForm)
	return c.JSON(ctrl.service.GetOsceExam().CreateSesi(c.UserContext(), *form))
}

func (ctrl *OsceExamController) UpdateSesi(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.SesiUjianForm)
	return c.JSON(ctrl.service.GetOsceExam().UpdateSesi(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceExamController) DeleteSesi(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().DeleteSesi(c.UserContext(), c.Params("id")))
}

// ===== ROTASI =====

func (ctrl *OsceExamController) GetRotasiBySesi(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetRotasiBySesi(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceExamController) CreateRotasi(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.RotasiUjianForm)
	return c.JSON(ctrl.service.GetOsceExam().CreateRotasi(c.UserContext(), *form))
}

func (ctrl *OsceExamController) UpdateRotasi(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.RotasiUjianForm)
	return c.JSON(ctrl.service.GetOsceExam().UpdateRotasi(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceExamController) DeleteRotasi(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().DeleteRotasi(c.UserContext(), c.Params("id")))
}

// ===== JADWAL =====

func (ctrl *OsceExamController) GetJadwalBySesi(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().GetJadwalBySesi(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceExamController) GenerateJadwal(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.GenerateJadwalRequest)
	return c.JSON(ctrl.service.GetOsceExam().GenerateJadwal(c.UserContext(), *form))
}

func (ctrl *OsceExamController) CreateJadwal(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.JadwalUjianForm)
	return c.JSON(ctrl.service.GetOsceExam().CreateJadwal(c.UserContext(), *form))
}

func (ctrl *OsceExamController) UpdateJadwal(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.JadwalUjianForm)
	return c.JSON(ctrl.service.GetOsceExam().UpdateJadwal(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceExamController) DeleteJadwal(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceExam().DeleteJadwal(c.UserContext(), c.Params("id")))
}
