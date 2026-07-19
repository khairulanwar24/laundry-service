// Package services (groupakses) berisi logika bisnis domain master group akses & group akses.
package services

import (
	"context"

	"sso-service/common/response"
	"sso-service/repositories"
)

// GroupAksesService membungkus akses ke repository registry.
type GroupAksesService struct {
	repository repositories.IRepositoryRegistry
}

// IGroupAksesService adalah kontrak logika bisnis domain group akses.
type IGroupAksesService interface {
	CreateMstGroupAkses(ctx context.Context, idMasterAplikasi, namaGroup, deskripsi string) response.Response
	GetMstGroupAkses(idMasterAplikasi string, limit, offset int, order, filter string) response.Response
	GetMstGroupAksesModul(idMasterAplikasi, idMasterGroup string, limit, offset int, order, filter string) response.Response
	GetGroupAkses(idMasterGroup string, limit, offset int, order, filter string) response.Response
	UpdateMstGroupAkses(ctx context.Context, idMasterGroup, namaGroup, deskripsi string) response.Response
	GetDetailMstGroupAkses(ctx context.Context, idMasterGroup string) response.Response
	DeleteMstGroupAkses(ctx context.Context, idMasterGroup string) response.Response
	CreateGroupAkses(ctx context.Context, idMasterAplikasi, idMasterGroup, idMasterModul, akses string) response.Response
	DeleteGroupAkses(ctx context.Context, idGroupAkses string) response.Response
	GetGroupAksesUserMenu(ctx context.Context, idUser, idMasterAplikasi string) response.Response
	GetGroupAksesUserApps(ctx context.Context, idUser string) response.Response
	CreateGroupAksesUserApps(ctx context.Context, idUser, idMasterAplikasi, idMasterGroup, statusData string) response.Response
	BulkCreateGroupAksesUserApps(ctx context.Context, idUsers []string, idMasterAplikasi, idMasterGroup string, statusData bool) response.Response
	UpdateGroupAksesUserApps(ctx context.Context, idTransUserGroup, statusData string) response.Response
	DeleteGroupAksesUserApps(ctx context.Context, idTransUserGroup string) response.Response
}

// NewGroupAksesService membuat instance GroupAksesService baru.
func NewGroupAksesService(repository repositories.IRepositoryRegistry) IGroupAksesService {
	return &GroupAksesService{repository: repository}
}

func (s *GroupAksesService) CreateMstGroupAkses(ctx context.Context, idMasterAplikasi, namaGroup, deskripsi string) response.Response {
	if err := s.repository.GetGroupAkses().CreateMstGroupAkses(ctx, idMasterAplikasi, namaGroup, deskripsi); err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *GroupAksesService) GetMstGroupAkses(idMasterAplikasi string, limit, offset int, order, filter string) response.Response {
	data := s.repository.GetGroupAkses().GetMstGroupAkses(idMasterAplikasi, limit, offset, order, filter)
	return response.Response{Success: true, Message: "Success", Data: data}
}

func (s *GroupAksesService) GetMstGroupAksesModul(idMasterAplikasi, idMasterGroup string, limit, offset int, order, filter string) response.Response {
	data := s.repository.GetGroupAkses().GetMstGroupAksesModul(idMasterAplikasi, idMasterGroup, limit, offset, order, filter)
	return response.Response{Success: true, Message: "Success", Data: data}
}

func (s *GroupAksesService) GetGroupAkses(idMasterGroup string, limit, offset int, order, filter string) response.Response {
	data := s.repository.GetGroupAkses().GetGroupAkses(idMasterGroup, limit, offset, order, filter)
	return response.Response{Success: true, Message: "Success", Data: data}
}

func (s *GroupAksesService) UpdateMstGroupAkses(ctx context.Context, idMasterGroup, namaGroup, deskripsi string) response.Response {
	rows, err := s.repository.GetGroupAkses().UpdateMstGroupAkses(ctx, idMasterGroup, namaGroup, deskripsi)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Update master_group"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "master_group tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *GroupAksesService) GetDetailMstGroupAkses(ctx context.Context, idMasterGroup string) response.Response {
	user, rows, err := s.repository.GetGroupAkses().GetDetailMstGroupAkses(ctx, idMasterGroup)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan Master Group"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "Master Group tidak ada"}
	}
	return response.Response{Success: true, Message: "Success", Data: user}
}

func (s *GroupAksesService) DeleteMstGroupAkses(ctx context.Context, idMasterGroup string) response.Response {
	rows, err := s.repository.GetGroupAkses().DeleteMstGroupAkses(ctx, idMasterGroup)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Delete Master Group"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "Master Group tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *GroupAksesService) CreateGroupAkses(ctx context.Context, idMasterAplikasi, idMasterGroup, idMasterModul, akses string) response.Response {
	if err := s.repository.GetGroupAkses().CreateGroupAkses(ctx, idMasterAplikasi, idMasterGroup, idMasterModul, akses); err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *GroupAksesService) DeleteGroupAkses(ctx context.Context, idGroupAkses string) response.Response {
	rows, err := s.repository.GetGroupAkses().DeleteGroupAkses(ctx, idGroupAkses)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Delete Group Akses"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "Group Akses tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *GroupAksesService) GetGroupAksesUserMenu(ctx context.Context, idUser, idMasterAplikasi string) response.Response {
	respmenu, err := s.repository.GetGroupAkses().GetGroupAksesUserMenu(ctx, idUser, idMasterAplikasi)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan Group Akses"}
	}
	return response.Response{Success: true, Message: "Success", Data: respmenu}
}

func (s *GroupAksesService) GetGroupAksesUserApps(ctx context.Context, idUser string) response.Response {
	respapps, err := s.repository.GetGroupAkses().GetGroupAksesUserApps(ctx, idUser)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan Akses Aplikasi"}
	}
	return response.Response{Success: true, Message: "Success", Data: respapps}
}

func (s *GroupAksesService) CreateGroupAksesUserApps(ctx context.Context, idUser, idMasterAplikasi, idMasterGroup, statusData string) response.Response {
	if err := s.repository.GetGroupAkses().CreateGroupAksesUserApps(ctx, idUser, idMasterAplikasi, idMasterGroup, statusData); err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *GroupAksesService) BulkCreateGroupAksesUserApps(ctx context.Context, idUsers []string, idMasterAplikasi, idMasterGroup string, statusData bool) response.Response {
	if err := s.repository.GetGroupAkses().BulkCreateGroupAksesUserApps(ctx, idUsers, idMasterAplikasi, idMasterGroup, statusData); err != nil {
		return response.Response{Success: false, Message: err.Error()}
	}
	return response.Response{Success: true, Message: "Bulk group akses berhasil disimpan"}
}

func (s *GroupAksesService) UpdateGroupAksesUserApps(ctx context.Context, idTransUserGroup, statusData string) response.Response {
	rows, err := s.repository.GetGroupAkses().UpdateGroupAksesUserApps(ctx, idTransUserGroup, statusData)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Update trans_user_group"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "trans_user_group tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}

func (s *GroupAksesService) DeleteGroupAksesUserApps(ctx context.Context, idTransUserGroup string) response.Response {
	rows, err := s.repository.GetGroupAkses().DeleteGroupAksesUserApps(ctx, idTransUserGroup)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Delete Master Group"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "Master Group tidak ada"}
	}
	return response.Response{Success: true, Message: "Success"}
}
