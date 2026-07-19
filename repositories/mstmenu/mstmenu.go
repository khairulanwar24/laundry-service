// Package repositories (mstmenu) adalah lapisan akses data untuk domain master menu & modul.
package repositories

import (
	"context"
	"strings"

	middleware "sso-service/middlewares"

	"gorm.io/gorm"
)

// MstMenuRepository memegang koneksi database utama.
type MstMenuRepository struct {
	db *gorm.DB
}

// IMstMenuRepository adalah kontrak akses data domain master menu & modul.
type IMstMenuRepository interface {
	GetMstMenu(idMasterAplikasi string, limit, offset int, order, filter string) map[string]interface{}
	GetDetailMstMenu(ctx context.Context, idMasterMenu string) ([]map[string]interface{}, int64, error)
	GetMstMenuModul(idMasterMenu string, limit, offset int, order, filter string) map[string]interface{}
	GetDetailMstMenuModul(ctx context.Context, idMasterModul string) ([]map[string]interface{}, int64, error)
	CreateMstMenu(ctx context.Context, idMasterAplikasi, namaMenu, deskripsi, order, icon string) error
	CreateMstMenuModul(ctx context.Context, idMasterAplikasi, idMasterMenu, namaModul, path, deskripsi, order, icon string) error
	UpdateMstMenu(ctx context.Context, idMasterMenu, namaMenu, deskripsi, order, icon string) (int64, error)
	UpdateMstMenuModul(ctx context.Context, idMasterModul, idMasterAplikasi, idMasterMenu, namaModul, path, deskripsi, order, icon string) (int64, error)
	DeleteMstMenu(ctx context.Context, idMasterMenu string) (int64, error)
	DeleteMstMenuModul(ctx context.Context, idMasterModul string) (int64, error)
}

// NewMstMenuRepository membuat instance MstMenuRepository baru.
func NewMstMenuRepository(db *gorm.DB) IMstMenuRepository {
	return &MstMenuRepository{db: db}
}

// GetMstMenu membangun query datatable master menu & mengembalikan hasil dari helper Datatables.
func (r *MstMenuRepository) GetMstMenu(idMasterAplikasi string, limit, offset int, order, filter string) map[string]interface{} {
	sRecursive := ``
	sTable := `SELECT id_master_menu, id_master_aplikasi, nama_menu, deskripsi, "order", icon FROM master_menu WHERE id_master_aplikasi = '` + idMasterAplikasi + `' AND status_data = true`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `AND (LOWER(nama_menu) LIKE ` + "'" + filter + "'" + ` OR LOWER(deskripsi) LIKE ` + "'" + filter + "')"
	}

	return middleware.Datatables(sRecursive, sTable, order, sFilter, limit, offset)
}

// GetDetailMstMenu mengambil detail satu master menu (data, rowsAffected, error).
func (r *MstMenuRepository) GetDetailMstMenu(ctx context.Context, idMasterMenu string) ([]map[string]interface{}, int64, error) {
	var user []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`Select id_master_menu, id_master_aplikasi, nama_menu, deskripsi, "order", icon FROM master_menu where id_master_menu = ?`, idMasterMenu).First(&user)
	return user, result.RowsAffected, result.Error
}

// GetMstMenuModul membangun query datatable master modul & mengembalikan hasil dari helper Datatables.
func (r *MstMenuRepository) GetMstMenuModul(idMasterMenu string, limit, offset int, order, filter string) map[string]interface{} {
	sRecursive := ``
	sTable := `SELECT id_master_modul, id_master_menu, "order" as order, nama_modul, path, deskripsi, id_master_aplikasi, icon FROM master_modul WHERE id_master_menu = '` + idMasterMenu + `' AND status_data = true`

	sFilter := ``
	if filter != "" {
		filter = "%" + strings.ToLower(filter) + "%"
		sFilter = `AND (LOWER(nama_modul) LIKE ` + "'" + filter + "'" + ` OR LOWER(path) LIKE ` + "'" + filter + "'" + ` OR LOWER(deskripsi) LIKE ` + "'" + filter + "')"
	}

	return middleware.Datatables(sRecursive, sTable, order, sFilter, limit, offset)
}

// GetDetailMstMenuModul mengambil detail satu master modul (data, rowsAffected, error).
func (r *MstMenuRepository) GetDetailMstMenuModul(ctx context.Context, idMasterModul string) ([]map[string]interface{}, int64, error) {
	var user []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(` SELECT id_master_modul, id_master_menu, "order" as order, nama_modul, path, deskripsi, id_master_aplikasi, icon FROM master_modul where id_master_modul = ?`, idMasterModul).First(&user)
	return user, result.RowsAffected, result.Error
}

// CreateMstMenu menyimpan master menu baru.
func (r *MstMenuRepository) CreateMstMenu(ctx context.Context, idMasterAplikasi, namaMenu, deskripsi, order, icon string) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO master_menu (id_master_aplikasi, nama_menu, deskripsi, "order", icon) VALUES (?, ?, ?, ?, ?)`, idMasterAplikasi, namaMenu, deskripsi, order, icon).Error
}

// CreateMstMenuModul menyimpan master modul baru.
func (r *MstMenuRepository) CreateMstMenuModul(ctx context.Context, idMasterAplikasi, idMasterMenu, namaModul, path, deskripsi, order, icon string) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO master_modul ( id_master_menu, "order", nama_modul, path, deskripsi, id_master_aplikasi, icon) VALUES (?, ?, ?, ?, ?, ?, ?)`, idMasterMenu, order, namaModul, path, deskripsi, idMasterAplikasi, icon).Error
}

// UpdateMstMenu memperbarui master menu.
func (r *MstMenuRepository) UpdateMstMenu(ctx context.Context, idMasterMenu, namaMenu, deskripsi, order, icon string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE master_menu
								SET nama_menu = ?, deskripsi = ?, "order" = ?, icon = ?
								WHERE id_master_menu = ?`, namaMenu, deskripsi, order, icon, idMasterMenu)
	return result.RowsAffected, result.Error
}

// UpdateMstMenuModul memperbarui master modul.
func (r *MstMenuRepository) UpdateMstMenuModul(ctx context.Context, idMasterModul, idMasterAplikasi, idMasterMenu, namaModul, path, deskripsi, order, icon string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE master_modul
								SET id_master_menu = ?, "order" = ?, nama_modul = ?, path = ?, deskripsi = ?, id_master_aplikasi = ?, icon = ?
								WHERE id_master_modul = ?`, idMasterMenu, order, namaModul, path, deskripsi, idMasterAplikasi, icon, idMasterModul)
	return result.RowsAffected, result.Error
}

// DeleteMstMenu menonaktifkan (soft delete) master menu.
func (r *MstMenuRepository) DeleteMstMenu(ctx context.Context, idMasterMenu string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE master_menu
								SET status_data = false
								WHERE id_master_menu = ?`, idMasterMenu)
	return result.RowsAffected, result.Error
}

// DeleteMstMenuModul menonaktifkan (soft delete) master modul.
func (r *MstMenuRepository) DeleteMstMenuModul(ctx context.Context, idMasterModul string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE master_modul
								SET status_data = false
								WHERE id_master_modul = ?`, idMasterModul)
	return result.RowsAffected, result.Error
}
