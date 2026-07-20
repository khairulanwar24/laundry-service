// Package services adalah pusat pendaftaran (registry) seluruh service (logika bisnis).
package services

import (
	"laundry-service/repositories"
	authService "laundry-service/services/auth"
	groupAksesService "laundry-service/services/groupakses"
	masterAppService "laundry-service/services/masterapp"
	mstMenuService "laundry-service/services/mstmenu"
	osceService "laundry-service/services/osce"
	refService "laundry-service/services/ref"
	userService "laundry-service/services/user"
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
	GetAuth() authService.IAuthService
	GetOsceStation() osceService.IOsceStationService
	GetOsceExam() osceService.IOsceExamService
	GetOsceAssessment() osceService.IOsceAssessmentService
	GetOsceUser() osceService.IOsceUserService
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

// GetAuth mengembalikan service autentikasi.
func (r *Registry) GetAuth() authService.IAuthService {
	return authService.NewAuthService(r.repository)
}

// GetOsceStation mengembalikan service OSCE station.
func (r *Registry) GetOsceStation() osceService.IOsceStationService {
	return osceService.NewOsceStationService(r.repository)
}

// GetOsceExam mengembalikan service OSCE exam.
func (r *Registry) GetOsceExam() osceService.IOsceExamService {
	return osceService.NewOsceExamService(r.repository)
}

// GetOsceAssessment mengembalikan service OSCE assessment.
func (r *Registry) GetOsceAssessment() osceService.IOsceAssessmentService {
	return osceService.NewOsceAssessmentService(r.repository)
}

// GetOsceUser mengembalikan service OSCE user.
func (r *Registry) GetOsceUser() osceService.IOsceUserService {
	return osceService.NewOsceUserService(r.repository)
}
