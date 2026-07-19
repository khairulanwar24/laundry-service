// Package services (mstmenu) berisi logika bisnis domain master menu & modul.
package services

import (
	"context"

	"sso-service/common/response"
	"sso-service/repositories"
)

// MstMenuService membungkus akses ke repository registry.
type MstMenuService struct {
	repository repositories.IRepositoryRegistry
}

// IMstMenuService adalah kontrak logika bisnis domain master menu & modul.
type IMstMenuService interface {
	GetMstMenu(idMasterAplikasi string, limit, offset int, order, filter string) response.Response
	GetDetailMstMenu(ctx context.Context, idMasterMenu string) response.Response
	GetMstMenuModul(idMasterMenu string, limit, offset int, order, filter string) response.Response
	GetDetailMstMenuModul(ctx context.Context, idMasterModul string) response.Response
	CreateMstMenu(
		ctx context.Context,
		idMasterAplikasi,
		namaMenu,
		deskripsi,
		order,
		icon string,
	) response.Response
	CreateMstMenuModul(
		ctx context.Context,
		idMasterAplikasi,
		idMasterMenu,
		namaModul,
		path,
		deskripsi,
		order,
		icon string,
	) response.Response
	UpdateMstMenu(
		ctx context.Context,
		idMasterMenu,
		namaMenu,
		deskripsi,
		order,
		icon string,
	) response.Response
	UpdateMstMenuModul(
		ctx context.Context,
		idMasterModul,
		idMasterAplikasi,
		idMasterMenu,
		namaModul,
		path,
		deskripsi,
		order,
		icon string,
	) response.Response
	DeleteMstMenu(ctx context.Context, idMasterMenu string) response.Response
	DeleteMstMenuModul(ctx context.Context, idMasterModul string) response.Response
}

// NewMstMenuService membuat instance MstMenuService baru.
func NewMstMenuService(repository repositories.IRepositoryRegistry) IMstMenuService {
	return &MstMenuService{repository: repository}
}

func (s *MstMenuService) GetMstMenu(idMasterAplikasi string, limit, offset int, order, filter string) response.Response {
	data := s.repository.GetMstMenu().GetMstMenu(idMasterAplikasi, limit, offset, order, filter)
	return response.Response{Success: true, Message: "Success", Data: data}
}

func (s *MstMenuService) GetDetailMstMenu(ctx context.Context, idMasterMenu string) response.Response {
	user, rows, err := s.repository.GetMstMenu().GetDetailMstMenu(ctx, idMasterMenu)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan users"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "users tidak ada"}
	}
	return response.Response{Success: true, Message: "Success", Data: user}
}

func (s *MstMenuService) GetMstMenuModul(idMasterMenu string, limit, offset int, order, filter string) response.Response {
	data := s.repository.GetMstMenu().GetMstMenuModul(idMasterMenu, limit, offset, order, filter)
	return response.Response{Success: true, Message: "Success", Data: data}
}

func (s *MstMenuService) GetDetailMstMenuModul(ctx context.Context, idMasterModul string) response.Response {
	user, rows, err := s.repository.GetMstMenu().GetDetailMstMenuModul(ctx, idMasterModul)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan users"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "users tidak ada"}
	}
	return response.Response{Success: true, Message: "Success", Data: user}
}

func (s *MstMenuService) CreateMstMenu(
	ctx context.Context,
	idMasterAplikasi,
	namaMenu,
	deskripsi,
	order,
	icon string,
) response.Response {
	if err := s.repository.GetMstMenu().CreateMstMenu(ctx, idMasterAplikasi, namaMenu, deskripsi, order, icon); err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *MstMenuService) CreateMstMenuModul(
	ctx context.Context,
	idMasterAplikasi,
	idMasterMenu,
	namaModul,
	path,
	deskripsi,
	order,
	icon string,
) response.Response {
	if err := s.repository.GetMstMenu().CreateMstMenuModul(ctx, idMasterAplikasi, idMasterMenu, namaModul, path, deskripsi, order, icon); err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *MstMenuService) UpdateMstMenu(
	ctx context.Context,
	idMasterMenu,
	namaMenu,
	deskripsi,
	order,
	icon string,
) response.Response {
	rows, err := s.repository.GetMstMenu().UpdateMstMenu(ctx, idMasterMenu, namaMenu, deskripsi, order, icon)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Update master_menu"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "master_menu tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *MstMenuService) UpdateMstMenuModul(
	ctx context.Context,
	idMasterModul,
	idMasterAplikasi,
	idMasterMenu,
	namaModul,
	path,
	deskripsi,
	order,
	icon string,
) response.Response {
	rows, err := s.repository.GetMstMenu().UpdateMstMenuModul(ctx, idMasterModul, idMasterAplikasi, idMasterMenu, namaModul, path, deskripsi, order, icon)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Update master modul"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "master modul tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *MstMenuService) DeleteMstMenu(ctx context.Context, idMasterMenu string) response.Response {
	rows, err := s.repository.GetMstMenu().DeleteMstMenu(ctx, idMasterMenu)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Delete Master Menu"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "Master Menu tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *MstMenuService) DeleteMstMenuModul(ctx context.Context, idMasterModul string) response.Response {
	rows, err := s.repository.GetMstMenu().DeleteMstMenuModul(ctx, idMasterModul)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Delete Master Modul"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "Master Modul tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}
