// Package services adalah pusat pendaftaran (registry) seluruh service (logika bisnis).
package services

import (
	"sso-service/repositories"
	groupAksesService "sso-service/services/groupakses"
	masterAppService "sso-service/services/masterapp"
	mstMenuService "sso-service/services/mstmenu"
	refService "sso-service/services/ref"
	userService "sso-service/services/user"
)

// Registry menyimpan repository registry sebagai dependency.
type Registry struct {
	repository repositories.IRepositoryRegistry
}

// IServiceRegistry adalah kontrak untuk mengambil service per-domain.
type IServiceRegistry interface {
	GetRef() refService.IRefService
	GetUser() userService.IUserService
	GetMasterApp() masterAppService.IMasterAppService
	GetMstMenu() mstMenuService.IMstMenuService
	GetGroupAkses() groupAksesService.IGroupAksesService
}

// NewServiceRegistry membuat service registry baru.
func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

// GetRef mengembalikan service referensi.
func (r *Registry) GetRef() refService.IRefService {
	return refService.NewRefService(r.repository)
}

// GetUser mengembalikan service user.
func (r *Registry) GetUser() userService.IUserService {
	return userService.NewUserService(r.repository)
}

// GetMasterApp mengembalikan service master aplikasi.
func (r *Registry) GetMasterApp() masterAppService.IMasterAppService {
	return masterAppService.NewMasterAppService(r.repository)
}

// GetMstMenu mengembalikan service master menu & modul.
func (r *Registry) GetMstMenu() mstMenuService.IMstMenuService {
	return mstMenuService.NewMstMenuService(r.repository)
}

// GetGroupAkses mengembalikan service group akses.
func (r *Registry) GetGroupAkses() groupAksesService.IGroupAksesService {
	return groupAksesService.NewGroupAksesService(r.repository)
}
