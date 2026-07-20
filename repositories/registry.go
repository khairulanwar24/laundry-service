// registry.go adalah pusat pendaftaran (registry) seluruh repository per-domain.
package repositories

import (
	"gorm.io/gorm"
)

// Registry menyimpan koneksi database sebagai dependency.
type Registry struct {
	db *gorm.DB
}

// IRepositoryRegistry adalah kontrak untuk mengambil repository per-domain.
type IRepositoryRegistry interface {
}

// NewRepositoryRegistry membuat repository registry baru.
func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}
