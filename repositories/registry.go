// Package repositories adalah pusat pendaftaran (registry) seluruh repository.
// Registry menyimpan koneksi database dan menyediakan repository per-domain.
package repositories

import (
	masterAppRepo "sso-service/repositories/masterapp"
	refRepo "sso-service/repositories/ref"
	userRepo "sso-service/repositories/user"

	"gorm.io/gorm"
)

// Registry adalah wadah koneksi database. ios-service memakai 3 koneksi:
//   - db          : database utama (sso / public)
//   - dbAkademik  : database akademik
//   - dbDigiclass : database digiclass
type Registry struct {
	db          *gorm.DB
	dbAkademik  *gorm.DB
	dbDigiclass *gorm.DB
}

// IRepositoryRegistry adalah kontrak untuk mengambil repository per-domain.
type IRepositoryRegistry interface {
	GetRef() refRepo.IRefRepository
	GetUser() userRepo.IUserRepository
	GetMasterApp() masterAppRepo.IMasterAppRepository
}

// NewRepositoryRegistry membuat registry baru dengan ketiga koneksi database.
func NewRepositoryRegistry(db, dbAkademik, dbDigiclass *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db, dbAkademik: dbAkademik, dbDigiclass: dbDigiclass}
}

// GetRef mengembalikan repository referensi (memakai koneksi akademik).
func (r *Registry) GetRef() refRepo.IRefRepository {
	return refRepo.NewRefRepository(r.dbAkademik)
}

// GetUser mengembalikan repository user (memakai ketiga koneksi database).
func (r *Registry) GetUser() userRepo.IUserRepository {
	return userRepo.NewUserRepository(r.db, r.dbAkademik, r.dbDigiclass)
}

// GetMasterApp mengembalikan repository master aplikasi (koneksi utama).
func (r *Registry) GetMasterApp() masterAppRepo.IMasterAppRepository {
	return masterAppRepo.NewMasterAppRepository(r.db)
}
