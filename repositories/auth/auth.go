// Package repositories (auth) adalah lapisan akses data untuk domain autentikasi.
package repositories

import (
	"context"
	"time"

	"sso-service/domain/dto"

	"gorm.io/gorm"
)

// AuthRepository memegang koneksi database utama.
type AuthRepository struct {
	db *gorm.DB
}

// IAuthRepository adalah kontrak akses data domain autentikasi.
type IAuthRepository interface {
	FindUserByUsername(ctx context.Context, username string) (map[string]interface{}, error)
	FindUserBlockByUsername(ctx context.Context, username string) (dto.UserBlock, error)
	InsertPasswordResetNoUser(
		ctx context.Context,
		uuid,
		username,
		otp string,
		tglInsert time.Time,
	) error
	DeactivatePasswordResetByUser(ctx context.Context, idUser string) error
	InsertPasswordResetWithUser(
		ctx context.Context,
		uuid,
		username,
		idUser,
		otp string,
		tglInsert time.Time,
	) (int64, error)
	FindPasswordReset(ctx context.Context, idPasswordReset string) ([]map[string]interface{}, error)
	UpdatePasswordResetPercobaan(ctx context.Context, percobaan int32, idPasswordReset string) error
	LockUser(ctx context.Context, idUser string) error
	FindUserFirstLogin(ctx context.Context, idUser string) ([]map[string]interface{}, error)
	UpdateUserPasswordFirstLogin(ctx context.Context, passwordHash, idUser string) (int64, error)
	UpdateUserPasswordNormal(ctx context.Context, passwordHash, idUser string) (int64, error)
}

// NewAuthRepository membuat instance AuthRepository baru.
func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &AuthRepository{db: db}
}

// FindUserByUsername mengambil data user (termasuk hash password) berdasarkan username.
func (r *AuthRepository) FindUserByUsername(ctx context.Context, username string) (map[string]interface{}, error) {
	var user map[string]interface{}
	result := r.db.WithContext(ctx).Raw("select username,nama_lengkap,email,no_hp,avatar,id_user,id_person,first_login,password,jenis_user from users where username = ? ", username).First(&user)
	return user, result.Error
}

// FindUserBlockByUsername mengambil status penguncian akun berdasarkan username.
func (r *AuthRepository) FindUserBlockByUsername(ctx context.Context, username string) (dto.UserBlock, error) {
	var user dto.UserBlock
	result := r.db.WithContext(ctx).Raw("select id_user,email,tgl_lock > NOW() AS still_blocked, EXTRACT(EPOCH FROM (tgl_lock - NOW())) AS seconds_remaining from users where username = ? ", username).Scan(&user)
	return user, result.Error
}

// InsertPasswordResetNoUser menyimpan record reset password tanpa id_user (fallback).
func (r *AuthRepository) InsertPasswordResetNoUser(
	ctx context.Context,
	uuid,
	username,
	otp string,
	tglInsert time.Time,
) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO password_reset (id_password_reset,username, otp, tgl_insert) VALUES (?,?, ?, ?)`, uuid, username, otp, tglInsert).Error
}

// DeactivatePasswordResetByUser menonaktifkan record reset password milik user.
func (r *AuthRepository) DeactivatePasswordResetByUser(ctx context.Context, idUser string) error {
	return r.db.WithContext(ctx).Exec(`Update password_reset set status_data = false where id_user= ?`, idUser).Error
}

// InsertPasswordResetWithUser menyimpan record reset password dengan id_user.
func (r *AuthRepository) InsertPasswordResetWithUser(
	ctx context.Context,
	uuid,
	username,
	idUser,
	otp string,
	tglInsert time.Time,
) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`INSERT INTO password_reset (id_password_reset,username,  id_user, otp, tgl_insert) VALUES (?,?, ?, ?, ?)`, uuid, username, idUser, otp, tglInsert)
	return result.RowsAffected, result.Error
}

// FindPasswordReset mengambil record reset password aktif berdasarkan id.
func (r *AuthRepository) FindPasswordReset(ctx context.Context, idPasswordReset string) ([]map[string]interface{}, error) {
	var user []map[string]interface{}
	result := r.db.WithContext(ctx).Raw("select username,id_user,tgl_insert,percobaan , otp from password_reset where id_password_reset = ? and status_data = true ", idPasswordReset).First(&user)
	return user, result.Error
}

// UpdatePasswordResetPercobaan memperbarui jumlah percobaan OTP.
func (r *AuthRepository) UpdatePasswordResetPercobaan(ctx context.Context, percobaan int32, idPasswordReset string) error {
	return r.db.WithContext(ctx).Exec(`Update password_reset set percobaan = ? where id_password_reset = ? and status_data = true`, percobaan, idPasswordReset).Error
}

// LockUser mengunci akun user selama 30 menit.
func (r *AuthRepository) LockUser(ctx context.Context, idUser string) error {
	return r.db.WithContext(ctx).Exec(`Update users set tgl_lock =  NOW() + INTERVAL '30 minutes' where id_user = ? and status_data = true`, idUser).Error
}

// FindUserFirstLogin mengambil status first_login user.
func (r *AuthRepository) FindUserFirstLogin(ctx context.Context, idUser string) ([]map[string]interface{}, error) {
	var user []map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`SELECT
								id_user
								, first_login
								FROM  users where id_user = ? and status_data = true`, idUser).First(&user)
	return user, result.Error
}

// UpdateUserPasswordFirstLogin memperbarui password saat login pertama.
func (r *AuthRepository) UpdateUserPasswordFirstLogin(ctx context.Context, passwordHash, idUser string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE users
								SET password = ?,tgl_update = NOW(),first_login = false,tgl_first_login = NOW()
								WHERE id_user = ?`, passwordHash, idUser)
	return result.RowsAffected, result.Error
}

// UpdateUserPasswordNormal memperbarui password biasa.
func (r *AuthRepository) UpdateUserPasswordNormal(ctx context.Context, passwordHash, idUser string) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`UPDATE users
								SET password = ?,tgl_update = NOW()
								WHERE id_user = ?`, passwordHash, idUser)
	return result.RowsAffected, result.Error
}
