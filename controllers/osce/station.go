// Package osce adalah lapisan HTTP handler (Fiber) untuk domain OSCE station.
package osce

import (
	"sso-service/domain/dto/osce"
	"sso-service/services"

	"github.com/gofiber/fiber/v2"
)

// OsceStationController membungkus service registry.
type OsceStationController struct {
	service services.IServiceRegistry
}

// IOsceStationController adalah kontrak handler HTTP domain OSCE station.
type IOsceStationController interface {
	// Tipe Station
	GetTipeStation(*fiber.Ctx) error
	GetTipeStationByID(*fiber.Ctx) error
	CreateTipeStation(*fiber.Ctx) error
	UpdateTipeStation(*fiber.Ctx) error
	DeleteTipeStation(*fiber.Ctx) error

	// Master Station
	GetStations(*fiber.Ctx) error
	GetStationByID(*fiber.Ctx) error
	CreateStation(*fiber.Ctx) error
	UpdateStation(*fiber.Ctx) error
	DeleteStation(*fiber.Ctx) error

	// Kompetensi
	GetKompetensiByStation(*fiber.Ctx) error
	CreateKompetensi(*fiber.Ctx) error
	UpdateKompetensi(*fiber.Ctx) error
	DeleteKompetensi(*fiber.Ctx) error

	// Checklist
	GetChecklistByStation(*fiber.Ctx) error
	CreateChecklist(*fiber.Ctx) error
	UpdateChecklist(*fiber.Ctx) error
	DeleteChecklist(*fiber.Ctx) error
	ReorderChecklist(*fiber.Ctx) error
}

// NewOsceStationController membuat instance baru.
func NewOsceStationController(service services.IServiceRegistry) IOsceStationController {
	return &OsceStationController{service: service}
}

// ===== TIPE STATION =====

func (ctrl *OsceStationController) GetTipeStation(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)
	order := c.Query("order", "")
	filter := c.Query("filter", "")
	data := ctrl.service.GetOsceStation().GetTipeStation(c.UserContext(), limit, offset, order, filter)
	return c.JSON(data)
}

func (ctrl *OsceStationController) GetTipeStationByID(c *fiber.Ctx) error {
	id := c.Params("id")
	data := ctrl.service.GetOsceStation().GetTipeStationByID(c.UserContext(), id)
	return c.JSON(data)
}

func (ctrl *OsceStationController) CreateTipeStation(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.TipeStationForm)
	data := ctrl.service.GetOsceStation().CreateTipeStation(c.UserContext(), *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) UpdateTipeStation(c *fiber.Ctx) error {
	id := c.Params("id")
	form := c.Locals("validatedForm").(*osce.TipeStationForm)
	data := ctrl.service.GetOsceStation().UpdateTipeStation(c.UserContext(), id, *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) DeleteTipeStation(c *fiber.Ctx) error {
	id := c.Params("id")
	data := ctrl.service.GetOsceStation().DeleteTipeStation(c.UserContext(), id)
	return c.JSON(data)
}

// ===== MASTER STATION =====

func (ctrl *OsceStationController) GetStations(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.StationQuery)
	data := ctrl.service.GetOsceStation().GetStations(c.UserContext(), *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) GetStationByID(c *fiber.Ctx) error {
	id := c.Params("id")
	data := ctrl.service.GetOsceStation().GetStationByID(c.UserContext(), id)
	return c.JSON(data)
}

func (ctrl *OsceStationController) CreateStation(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.StationForm)
	data := ctrl.service.GetOsceStation().CreateStation(c.UserContext(), *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) UpdateStation(c *fiber.Ctx) error {
	id := c.Params("id")
	form := c.Locals("validatedForm").(*osce.StationForm)
	data := ctrl.service.GetOsceStation().UpdateStation(c.UserContext(), id, *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) DeleteStation(c *fiber.Ctx) error {
	id := c.Params("id")
	data := ctrl.service.GetOsceStation().DeleteStation(c.UserContext(), id)
	return c.JSON(data)
}

// ===== KOMPETENSI =====

func (ctrl *OsceStationController) GetKompetensiByStation(c *fiber.Ctx) error {
	idStation := c.Params("id")
	data := ctrl.service.GetOsceStation().GetKompetensiByStation(c.UserContext(), idStation)
	return c.JSON(data)
}

func (ctrl *OsceStationController) CreateKompetensi(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.KompetensiStationForm)
	data := ctrl.service.GetOsceStation().CreateKompetensi(c.UserContext(), *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) UpdateKompetensi(c *fiber.Ctx) error {
	id := c.Params("id")
	form := c.Locals("validatedForm").(*osce.KompetensiStationForm)
	data := ctrl.service.GetOsceStation().UpdateKompetensi(c.UserContext(), id, *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) DeleteKompetensi(c *fiber.Ctx) error {
	id := c.Params("id")
	data := ctrl.service.GetOsceStation().DeleteKompetensi(c.UserContext(), id)
	return c.JSON(data)
}

// ===== CHECKLIST =====

func (ctrl *OsceStationController) GetChecklistByStation(c *fiber.Ctx) error {
	idStation := c.Params("id")
	data := ctrl.service.GetOsceStation().GetChecklistByStation(c.UserContext(), idStation)
	return c.JSON(data)
}

func (ctrl *OsceStationController) CreateChecklist(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*osce.ChecklistStationForm)
	data := ctrl.service.GetOsceStation().CreateChecklist(c.UserContext(), *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) UpdateChecklist(c *fiber.Ctx) error {
	id := c.Params("id")
	form := c.Locals("validatedForm").(*osce.ChecklistStationForm)
	data := ctrl.service.GetOsceStation().UpdateChecklist(c.UserContext(), id, *form)
	return c.JSON(data)
}

func (ctrl *OsceStationController) DeleteChecklist(c *fiber.Ctx) error {
	id := c.Params("id")
	data := ctrl.service.GetOsceStation().DeleteChecklist(c.UserContext(), id)
	return c.JSON(data)
}

func (ctrl *OsceStationController) ReorderChecklist(c *fiber.Ctx) error {
	idStation := c.Params("id")
	form := c.Locals("validatedForm").(*osce.ChecklistReorderRequest)
	data := ctrl.service.GetOsceStation().ReorderChecklist(c.UserContext(), idStation, *form)
	return c.JSON(data)
}
