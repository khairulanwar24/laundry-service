// registry.go adalah pusat pendaftaran (registry) seluruh service (logika bisnis) per-domain.
package services

import (
	"laundry-service/repositories"
	accountService "laundry-service/services/account"
)

// Registry menyimpan repository registry sebagai dependency.
type Registry struct {
	repository repositories.IRepositoryRegistry
}

// IServiceRegistry adalah kontrak untuk mengambil service per-domain.
type IServiceRegistry interface {
	GetAccount() accountService.IAccountService
}

// NewServiceRegistry membuat service registry baru.
func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

// GetAccount mengembalikan service akun.
func (r *Registry) GetAccount() accountService.IAccountService {
	return accountService.NewAccountService(r.repository)
}
