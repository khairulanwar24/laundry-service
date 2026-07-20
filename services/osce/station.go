// Package osce adalah lapisan business logic untuk domain OSCE station.
package osce

import (
	"context"
	"laundry-service/common/response"
	dtos "laundry-service/domain/dto/osce"
	"laundry-service/repositories"
)

// OsceStationService membungkus repository registry.
type OsceStationService struct {
	repo repositories.IRepositoryRegistry
}

// IOsceStationService adalah kontrak business logic domain station.
type IOsceStationService interface {
	// Tipe Station
	GetTipeStation(ctx context.Context, limit, offset int, order, filter string) response.Response
	GetTipeStationByID(ctx context.Context, id string) response.Response
	CreateTipeStation(ctx context.Context, form dtos.TipeStationForm) response.Response
	UpdateTipeStation(ctx context.Context, id string, form dtos.TipeStationForm) response.Response
	DeleteTipeStation(ctx context.Context, id string) response.Response

	// Master Station
	GetStations(ctx context.Context, form dtos.StationQuery) response.Response
	GetStationByID(ctx context.Context, id string) response.Response
	CreateStation(ctx context.Context, form dtos.StationForm) response.Response
	UpdateStation(ctx context.Context, id string, form dtos.StationForm) response.Response
	DeleteStation(ctx context.Context, id string) response.Response

	// Kompetensi
	GetKompetensiByStation(ctx context.Context, idStation string) response.Response
	CreateKompetensi(ctx context.Context, form dtos.KompetensiStationForm) response.Response
	UpdateKompetensi(ctx context.Context, id string, form dtos.KompetensiStationForm) response.Response
	DeleteKompetensi(ctx context.Context, id string) response.Response

	// Checklist
	GetChecklistByStation(ctx context.Context, idStation string) response.Response
	CreateChecklist(ctx context.Context, form dtos.ChecklistStationForm) response.Response
	UpdateChecklist(ctx context.Context, id string, form dtos.ChecklistStationForm) response.Response
	DeleteChecklist(ctx context.Context, id string) response.Response
	ReorderChecklist(ctx context.Context, idStation string, form dtos.ChecklistReorderRequest) response.Response
}

// NewOsceStationService membuat instance baru.
func NewOsceStationService(repo repositories.IRepositoryRegistry) IOsceStationService {
	return &OsceStationService{repo: repo}
}

// ===== TIPE STATION =====

func (s *OsceStationService) GetTipeStation(ctx context.Context, limit, offset int, order, filter string) response.Response {
	data, err := s.repo.GetOsceStation().FindTipeStation(ctx, limit, offset, order, filter)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	count, _ := s.repo.GetOsceStation().CountTipeStation(ctx, filter)
	return response.Response{
		Success: true, Message: "sukses",
		Data: map[string]interface{}{"data": data, "total_data": count},
	}
}

func (s *OsceStationService) GetTipeStationByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceStation().FindTipeStationByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceStationService) CreateTipeStation(ctx context.Context, form dtos.TipeStationForm) response.Response {
	if err := s.repo.GetOsceStation().InsertTipeStation(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat tipe station: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Tipe station berhasil dibuat"}
}

func (s *OsceStationService) UpdateTipeStation(ctx context.Context, id string, form dtos.TipeStationForm) response.Response {
	if err := s.repo.GetOsceStation().UpdateTipeStation(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Tipe station berhasil diperbarui"}
}

func (s *OsceStationService) DeleteTipeStation(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceStation().DeleteTipeStation(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Tipe station berhasil dihapus"}
}

// ===== MASTER STATION =====

func (s *OsceStationService) GetStations(ctx context.Context, form dtos.StationQuery) response.Response {
	data, err := s.repo.GetOsceStation().FindStations(ctx, form.Limit, form.Offset, form.Order, form.Filter)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	count, _ := s.repo.GetOsceStation().CountStations(ctx, form.Filter)
	return response.Response{
		Success: true, Message: "sukses",
		Data:    map[string]interface{}{"data": data, "total_data": count},
	}
}

func (s *OsceStationService) GetStationByID(ctx context.Context, id string) response.Response {
	data, err := s.repo.GetOsceStation().FindStationByID(ctx, id)
	if err != nil {
		return response.Response{Success: false, Message: "Data tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceStationService) CreateStation(ctx context.Context, form dtos.StationForm) response.Response {
	id, err := s.repo.GetOsceStation().InsertStation(ctx, form)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat station: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Station berhasil dibuat", Data: map[string]string{"id": id}}
}

func (s *OsceStationService) UpdateStation(ctx context.Context, id string, form dtos.StationForm) response.Response {
	if err := s.repo.GetOsceStation().UpdateStation(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Station berhasil diperbarui"}
}

func (s *OsceStationService) DeleteStation(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceStation().DeleteStation(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Station berhasil dihapus"}
}

// ===== KOMPETENSI =====

func (s *OsceStationService) GetKompetensiByStation(ctx context.Context, idStation string) response.Response {
	data, err := s.repo.GetOsceStation().FindKompetensiByStation(ctx, idStation)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceStationService) CreateKompetensi(ctx context.Context, form dtos.KompetensiStationForm) response.Response {
	if err := s.repo.GetOsceStation().InsertKompetensi(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat kompetensi: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Kompetensi berhasil dibuat"}
}

func (s *OsceStationService) UpdateKompetensi(ctx context.Context, id string, form dtos.KompetensiStationForm) response.Response {
	if err := s.repo.GetOsceStation().UpdateKompetensi(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Kompetensi berhasil diperbarui"}
}

func (s *OsceStationService) DeleteKompetensi(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceStation().DeleteKompetensi(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Kompetensi berhasil dihapus"}
}

// ===== CHECKLIST =====

func (s *OsceStationService) GetChecklistByStation(ctx context.Context, idStation string) response.Response {
	data, err := s.repo.GetOsceStation().FindChecklistByStation(ctx, idStation)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mengambil data: " + err.Error()}
	}
	return response.Response{Success: true, Message: "sukses", Data: data}
}

func (s *OsceStationService) CreateChecklist(ctx context.Context, form dtos.ChecklistStationForm) response.Response {
	if err := s.repo.GetOsceStation().InsertChecklist(ctx, form); err != nil {
		return response.Response{Success: false, Message: "Gagal membuat checklist: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Checklist berhasil dibuat"}
}

func (s *OsceStationService) UpdateChecklist(ctx context.Context, id string, form dtos.ChecklistStationForm) response.Response {
	if err := s.repo.GetOsceStation().UpdateChecklist(ctx, id, form); err != nil {
		return response.Response{Success: false, Message: "Gagal memperbarui: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Checklist berhasil diperbarui"}
}

func (s *OsceStationService) DeleteChecklist(ctx context.Context, id string) response.Response {
	if err := s.repo.GetOsceStation().DeleteChecklist(ctx, id); err != nil {
		return response.Response{Success: false, Message: "Gagal menghapus: " + err.Error()}
	}
	return response.Response{Success: true, Message: "Checklist berhasil dihapus"}
}

func (s *OsceStationService) ReorderChecklist(ctx context.Context, idStation string, form dtos.ChecklistReorderRequest) response.Response {
	for _, item := range form.Items {
		if err := s.repo.GetOsceStation().ReorderChecklist(ctx, item.IDChecklist, item.Urutan); err != nil {
			return response.Response{Success: false, Message: "Gagal reorder: " + err.Error()}
		}
	}
	return response.Response{Success: true, Message: "Urutan checklist berhasil diperbarui"}
}
