// Package services adalah pusat pendaftaran (registry) seluruh service (logika bisnis).
package services

import (
	"sso-service/repositories"
	refService "sso-service/services/ref"
)

// Registry menyimpan repository registry sebagai dependency.
type Registry struct {
	repository repositories.IRepositoryRegistry
}

// IServiceRegistry adalah kontrak untuk mengambil service per-domain.
type IServiceRegistry interface {
	GetRef() refService.IRefService
}

// NewServiceRegistry membuat service registry baru.
func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

// GetRef mengembalikan service referensi.
func (r *Registry) GetRef() refService.IRefService {
	return refService.NewRefService(r.repository)
}
