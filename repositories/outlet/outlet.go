// Package repositories (outlet) adalah lapisan akses data untuk domain outlet.
package repositories

import (
	"context"

	"gorm.io/gorm"
)

// OutletRepository menangani akses data ke laundry.outlets & laundry.user_outlets.
type OutletRepository struct {
	db *gorm.DB
}

// IOutletRepository adalah kontrak akses data domain outlet.
type IOutletRepository interface {
	ListOutletsForUser(ctx context.Context, idUser string) ([]map[string]interface{}, error)
	CreateOutlet(ctx context.Context, idOwnerUser, nama, alamat, telepon, logoPath string) (string, error)
	FindOutletByID(ctx context.Context, idOutlet string) (map[string]interface{}, error)
	UpdateOutlet(ctx context.Context, idOutlet, nama, alamat, telepon, logoPath string) error
	AssignUser(ctx context.Context, idOutlet, idUser, role, permissionsJSON string) (string, error)
	FindMemberByUser(ctx context.Context, idOutlet, idUser string) (map[string]interface{}, error)
	FindMemberByID(ctx context.Context, idOutlet, idUserOutlet string) (map[string]interface{}, error)
	ListStaff(ctx context.Context, idOutlet string) ([]map[string]interface{}, error)
	UpdateMember(ctx context.Context, idUserOutlet, role, permissionsJSON string) error
	DeactivateMember(ctx context.Context, idUserOutlet string) error
	CountActiveOwners(ctx context.Context, idOutlet string) (int64, error)
}

// NewOutletRepository membuat instance OutletRepository baru.
func NewOutletRepository(db *gorm.DB) IOutletRepository {
	return &OutletRepository{db: db}
}

func (r *OutletRepository) ListOutletsForUser(ctx context.Context, idUser string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT o.id_outlet, o.nama, o.logo_path, o.alamat, o.telepon, o.is_active,
		       uo.role AS user_role, uo.permissions_json AS user_permissions, uo.tgl_insert AS joined_at
		FROM laundry.outlets o
		JOIN laundry.user_outlets uo ON uo.id_outlet = o.id_outlet
		WHERE uo.id_user = ? AND uo.is_active = true AND o.status_data = true
		ORDER BY o.nama
	`, idUser).Scan(&data)
	return data, result.Error
}

func (r *OutletRepository) CreateOutlet(ctx context.Context, idOwnerUser, nama, alamat, telepon, logoPath string) (string, error) {
	var idOutlet string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.outlets (id_owner_user, nama, alamat, telepon, logo_path)
		VALUES (?, ?, ?, ?, NULLIF(?, ''))
		RETURNING id_outlet
	`, idOwnerUser, nama, alamat, telepon, logoPath).Scan(&idOutlet)
	return idOutlet, result.Error
}

func (r *OutletRepository) FindOutletByID(ctx context.Context, idOutlet string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_outlet, id_owner_user, nama, logo_path, alamat, telepon, is_active, tgl_insert, tgl_update
		FROM laundry.outlets
		WHERE id_outlet = ? AND status_data = true
	`, idOutlet).First(&data)
	return data, result.Error
}

func (r *OutletRepository) UpdateOutlet(ctx context.Context, idOutlet, nama, alamat, telepon, logoPath string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.outlets
		SET nama = ?, alamat = ?, telepon = ?, logo_path = NULLIF(?, ''), tgl_update = NOW()
		WHERE id_outlet = ?
	`, nama, alamat, telepon, logoPath, idOutlet).Error
}

// AssignUser meng-upsert baris keanggotaan outlet (unique di id_user+id_outlet).
func (r *OutletRepository) AssignUser(ctx context.Context, idOutlet, idUser, role, permissionsJSON string) (string, error) {
	var idUserOutlet string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.user_outlets (id_user, id_outlet, role, permissions_json, is_active)
		VALUES (?, ?, ?, ?::jsonb, true)
		ON CONFLICT (id_user, id_outlet) DO UPDATE
		SET role = EXCLUDED.role, permissions_json = EXCLUDED.permissions_json, is_active = true, tgl_update = NOW()
		RETURNING id_user_outlet
	`, idUser, idOutlet, role, permissionsJSON).Scan(&idUserOutlet)
	return idUserOutlet, result.Error
}

func (r *OutletRepository) FindMemberByUser(ctx context.Context, idOutlet, idUser string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_user_outlet, id_user, id_outlet, role, permissions_json, is_active
		FROM laundry.user_outlets
		WHERE id_outlet = ? AND id_user = ? AND status_data = true
	`, idOutlet, idUser).First(&data)
	return data, result.Error
}

func (r *OutletRepository) FindMemberByID(ctx context.Context, idOutlet, idUserOutlet string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_user_outlet, id_user, id_outlet, role, permissions_json, is_active
		FROM laundry.user_outlets
		WHERE id_outlet = ? AND id_user_outlet = ? AND status_data = true
	`, idOutlet, idUserOutlet).First(&data)
	return data, result.Error
}

func (r *OutletRepository) ListStaff(ctx context.Context, idOutlet string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT uo.id_user_outlet, uo.id_user, u.nama, u.email, u.telepon,
		       uo.role, uo.permissions_json, uo.is_active, uo.tgl_insert AS joined_at
		FROM laundry.user_outlets uo
		JOIN laundry.users u ON u.id_user = uo.id_user
		WHERE uo.id_outlet = ? AND uo.is_active = true
		ORDER BY CASE WHEN uo.role = 'owner' THEN 0 ELSE 1 END, u.nama
	`, idOutlet).Scan(&data)
	return data, result.Error
}

func (r *OutletRepository) UpdateMember(ctx context.Context, idUserOutlet, role, permissionsJSON string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.user_outlets
		SET role = ?, permissions_json = ?::jsonb, tgl_update = NOW()
		WHERE id_user_outlet = ?
	`, role, permissionsJSON, idUserOutlet).Error
}

func (r *OutletRepository) DeactivateMember(ctx context.Context, idUserOutlet string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.user_outlets SET is_active = false, tgl_update = NOW() WHERE id_user_outlet = ?
	`, idUserOutlet).Error
}

func (r *OutletRepository) CountActiveOwners(ctx context.Context, idOutlet string) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Table("laundry.user_outlets").
		Where("id_outlet = ? AND role = 'owner' AND is_active = true", idOutlet).Count(&count)
	return count, result.Error
}
