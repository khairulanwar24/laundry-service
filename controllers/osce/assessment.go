// Package osce adalah lapisan HTTP handler (Fiber) untuk domain OSCE penilaian.
package osce

import (
	"sso-service/domain/dto/osce"
	"sso-service/services"

	"github.com/gofiber/fiber/v2"
)

// OsceAssessmentController membungkus service registry.
type OsceAssessmentController struct {
	service services.IServiceRegistry
}

// IOsceAssessmentController adalah kontrak handler HTTP domain OSCE penilaian.
type IOsceAssessmentController interface {
	GetFormPenilaian(*fiber.Ctx) error
	CreatePenilaian(*fiber.Ctx) error
	UpdatePenilaian(*fiber.Ctx) error
	KirimPenilaian(*fiber.Ctx) error
	GetPenilaianByID(*fiber.Ctx) error
	GetHasilUjian(*fiber.Ctx) error
	GetHasilSaya(*fiber.Ctx) error
	GetHasilByID(*fiber.Ctx) error
	KalkulasiHasil(*fiber.Ctx) error
	KalkulasiBorderline(*fiber.Ctx) error
	GetStatistikUjian(*fiber.Ctx) error
	GetJadwalPenguji(*fiber.Ctx) error
}

// NewOsceAssessmentController membuat instance baru.
func NewOsceAssessmentController(service services.IServiceRegistry) IOsceAssessmentController {
	return &OsceAssessmentController{service: service}
}

// ===== PENILAIAN =====

func (ctrl *OsceAssessmentController) GetFormPenilaian(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceAssessment().GetFormPenilaian(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceAssessmentController) CreatePenilaian(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.PenilaianForm)
	return c.JSON(ctrl.service.GetOsceAssessment().CreatePenilaian(c.UserContext(), *form))
}

func (ctrl *OsceAssessmentController) UpdatePenilaian(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.PenilaianUpdateForm)
	return c.JSON(ctrl.service.GetOsceAssessment().UpdatePenilaian(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceAssessmentController) KirimPenilaian(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceAssessment().KirimPenilaian(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceAssessmentController) GetPenilaianByID(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceAssessment().GetPenilaianByID(c.UserContext(), c.Params("id")))
}

// ===== HASIL UJIAN =====

func (ctrl *OsceAssessmentController) GetHasilUjian(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceAssessment().GetHasilUjian(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceAssessmentController) GetHasilSaya(c *fiber.Ctx) error {
	idUser := c.Locals("id_user").(string)
	return c.JSON(ctrl.service.GetOsceAssessment().GetHasilSaya(c.UserContext(), idUser))
}

func (ctrl *OsceAssessmentController) GetHasilByID(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceAssessment().GetHasilByID(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceAssessmentController) KalkulasiHasil(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.KalkulasiHasilRequest)
	return c.JSON(ctrl.service.GetOsceAssessment().KalkulasiHasil(c.UserContext(), *form))
}

func (ctrl *OsceAssessmentController) KalkulasiBorderline(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.BorderlineRequest)
	return c.JSON(ctrl.service.GetOsceAssessment().KalkulasiBorderline(c.UserContext(), *form))
}

func (ctrl *OsceAssessmentController) GetStatistikUjian(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceAssessment().GetStatistikUjian(c.UserContext(), c.Params("id")))
}

// ===== JADWAL PENGUJI =====

func (ctrl *OsceAssessmentController) GetJadwalPenguji(c *fiber.Ctx) error {
	idPenguji := c.Query("id_penguji", "")
	return c.JSON(ctrl.service.GetOsceAssessment().GetJadwalPenguji(c.UserContext(), idPenguji))
}
