// Package services adalah pusat pendaftaran (registry) seluruh service (logika bisnis).
package services

import (
	"sso-service/repositories"
	masterAppService "sso-service/services/masterapp"
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
