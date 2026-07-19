// Package osce adalah lapisan HTTP handler (Fiber) untuk domain pengguna OSCE.
package osce

import (
	"sso-service/domain/dto/osce"
	"sso-service/services"

	"github.com/gofiber/fiber/v2"
)

// OsceUserController membungkus service registry.
type OsceUserController struct {
	service services.IServiceRegistry
}

// IOsceUserController adalah kontrak handler HTTP domain pengguna OSCE.
type IOsceUserController interface {
	// Penguji
	GetPenguji(*fiber.Ctx) error
	GetPengujiByID(*fiber.Ctx) error
	CreatePenguji(*fiber.Ctx) error
	UpdatePenguji(*fiber.Ctx) error
	DeletePenguji(*fiber.Ctx) error

	// Kelompok Mahasiswa
	GetKelompokMahasiswa(*fiber.Ctx) error
	GetKelompokMahasiswaByID(*fiber.Ctx) error
	CreateKelompokMahasiswa(*fiber.Ctx) error
	UpdateKelompokMahasiswa(*fiber.Ctx) error
	DeleteKelompokMahasiswa(*fiber.Ctx) error

	// Anggota Kelompok
	GetAnggotaByKelompok(*fiber.Ctx) error
	TambahAnggota(*fiber.Ctx) error
	TambahAnggotaBulk(*fiber.Ctx) error
	HapusAnggota(*fiber.Ctx) error

	// Penugasan Penguji
	GetPenugasanBySesi(*fiber.Ctx) error
	AssignPenguji(*fiber.Ctx) error
	HapusPenugasan(*fiber.Ctx) error
}

// NewOsceUserController membuat instance baru.
func NewOsceUserController(service services.IServiceRegistry) IOsceUserController {
	return &OsceUserController{service: service}
}

// ===== PENGUJI =====

func (ctrl *OsceUserController) GetPenguji(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.PengujiQuery)
	return c.JSON(ctrl.service.GetOsceUser().GetPenguji(c.UserContext(), *form))
}

func (ctrl *OsceUserController) GetPengujiByID(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceUser().GetPengujiByID(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceUserController) CreatePenguji(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.PengujiForm)
	return c.JSON(ctrl.service.GetOsceUser().CreatePenguji(c.UserContext(), *form))
}

func (ctrl *OsceUserController) UpdatePenguji(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.PengujiForm)
	return c.JSON(ctrl.service.GetOsceUser().UpdatePenguji(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceUserController) DeletePenguji(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceUser().DeletePenguji(c.UserContext(), c.Params("id")))
}

// ===== KELOMPOK MAHASISWA =====

func (ctrl *OsceUserController) GetKelompokMahasiswa(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.KelompokMahasiswaQuery)
	return c.JSON(ctrl.service.GetOsceUser().GetKelompokMahasiswa(c.UserContext(), *form))
}

func (ctrl *OsceUserController) GetKelompokMahasiswaByID(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceUser().GetKelompokMahasiswaByID(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceUserController) CreateKelompokMahasiswa(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.KelompokMahasiswaForm)
	return c.JSON(ctrl.service.GetOsceUser().CreateKelompokMahasiswa(c.UserContext(), *form))
}

func (ctrl *OsceUserController) UpdateKelompokMahasiswa(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.KelompokMahasiswaForm)
	return c.JSON(ctrl.service.GetOsceUser().UpdateKelompokMahasiswa(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceUserController) DeleteKelompokMahasiswa(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceUser().DeleteKelompokMahasiswa(c.UserContext(), c.Params("id")))
}

// ===== ANGGOTA KELOMPOK =====

func (ctrl *OsceUserController) GetAnggotaByKelompok(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceUser().GetAnggotaByKelompok(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceUserController) TambahAnggota(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.AnggotaKelompokForm)
	return c.JSON(ctrl.service.GetOsceUser().TambahAnggota(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceUserController) TambahAnggotaBulk(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.BulkAnggotaKelompokForm)
	return c.JSON(ctrl.service.GetOsceUser().TambahAnggotaBulk(c.UserContext(), c.Params("id"), *form))
}

func (ctrl *OsceUserController) HapusAnggota(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceUser().HapusAnggota(c.UserContext(), c.Params("id")))
}

// ===== PENUGASAN PENGUJI =====

func (ctrl *OsceUserController) GetPenugasanBySesi(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceUser().GetPenugasanBySesi(c.UserContext(), c.Params("id")))
}

func (ctrl *OsceUserController) AssignPenguji(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.PenugasanPengujiForm)
	return c.JSON(ctrl.service.GetOsceUser().AssignPenguji(c.UserContext(), *form))
}

func (ctrl *OsceUserController) HapusPenugasan(c *fiber.Ctx) error {
	return c.JSON(ctrl.service.GetOsceUser().HapusPenugasan(c.UserContext(), c.Params("id")))
}
