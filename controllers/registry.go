// registry.go adalah pusat pendaftaran (registry) seluruh controller.
// Berada di package `controllers` (root) berdampingan dengan handler flat lama
// selama masa migrasi bertahap.
package controllers

import (
	refController "sso-service/controllers/ref"
	"sso-service/services"
)

// Registry menyimpan service registry sebagai dependency.
type Registry struct {
	service services.IServiceRegistry
}

// IControllerRegistry adalah kontrak untuk mengambil controller per-domain.
type IControllerRegistry interface {
	GetRef() refController.IRefController
}

// NewControllerRegistry membuat controller registry baru.
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

// GetRef mengembalikan controller referensi.
func (r *Registry) GetRef() refController.IRefController {
	return refController.NewRefController(r.service)
}
