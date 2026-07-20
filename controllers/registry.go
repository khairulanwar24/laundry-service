// registry.go adalah pusat pendaftaran (registry) seluruh controller per-domain.
package controllers

import (
	"laundry-service/services"
)

// Registry menyimpan service registry sebagai dependency.
type Registry struct {
	service services.IServiceRegistry
}

// IControllerRegistry adalah kontrak untuk mengambil controller per-domain.
type IControllerRegistry interface {
}

// NewControllerRegistry membuat controller registry baru.
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}
