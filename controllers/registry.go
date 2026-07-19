// registry.go adalah pusat pendaftaran (registry) seluruh controller.
// Berada di package `controllers` (root) berdampingan dengan handler flat lama
// selama masa migrasi bertahap.
package controllers

import (
	masterAppController "sso-service/controllers/masterapp"
	mstMenuController "sso-service/controllers/mstmenu"
	refController "sso-service/controllers/ref"
	userController "sso-service/controllers/user"
	"sso-service/services"
)

// Registry menyimpan service registry sebagai dependency.
type Registry struct {
	service services.IServiceRegistry
}

// IControllerRegistry adalah kontrak untuk mengambil controller per-domain.
type IControllerRegistry interface {
	GetRef() refController.IRefController
	GetUser() userController.IUserController
	GetMasterApp() masterAppController.IMasterAppController
	GetMstMenu() mstMenuController.IMstMenuController
}

// NewControllerRegistry membuat controller registry baru.
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

// GetRef mengembalikan controller referensi.
func (r *Registry) GetRef() refController.IRefController {
	return refController.NewRefController(r.service)
}

// GetUser mengembalikan controller user.
func (r *Registry) GetUser() userController.IUserController {
	return userController.NewUserController(r.service)
}

// GetMasterApp mengembalikan controller master aplikasi.
func (r *Registry) GetMasterApp() masterAppController.IMasterAppController {
	return masterAppController.NewMasterAppController(r.service)
}

// GetMstMenu mengembalikan controller master menu & modul.
func (r *Registry) GetMstMenu() mstMenuController.IMstMenuController {
	return mstMenuController.NewMstMenuController(r.service)
}
