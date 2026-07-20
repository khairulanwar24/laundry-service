// Package repositories (account) adalah lapisan akses data untuk domain akun laundry.
package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// AccountRepository menangani akses data ke laundry.users & laundry.password_resets.
type AccountRepository struct {
	db *gorm.DB
}

// IAccountRepository adalah kontrak akses data domain akun.
type IAccountRepository interface {
	FindUserByEmail(ctx context.Context, email string) (map[string]interface{}, error)
	FindUserByID(ctx context.Context, idUser string) (map[string]interface{}, error)
	ExistsEmail(ctx context.Context, email string) (bool, error)
	ExistsPhone(ctx context.Context, telepon string) (bool, error)
	CreateUser(ctx context.Context, nama, email, telepon, passwordHash, alamat string) (string, error)
	UpsertPasswordReset(ctx context.Context, email, otp, tokenHash string, sentAt time.Time) error
	FindPasswordResetByEmail(ctx context.Context, email string) (map[string]interface{}, error)
	DeletePasswordResetByEmail(ctx context.Context, email string) error
	UpdateUserPassword(ctx context.Context, idUser, passwordHash string) error
}

// NewAccountRepository membuat instance AccountRepository baru.
func NewAccountRepository(db *gorm.DB) IAccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) FindUserByEmail(ctx context.Context, email string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_user, nama, email, telepon, alamat, password, is_active
		FROM laundry.users
		WHERE email = ? AND status_data = true
	`, email).First(&data)
	return data, result.Error
}

func (r *AccountRepository) FindUserByID(ctx context.Context, idUser string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_user, nama, email, telepon, alamat, is_active
		FROM laundry.users
		WHERE id_user = ? AND status_data = true
	`, idUser).First(&data)
	return data, result.Error
}

func (r *AccountRepository) ExistsEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Table("laundry.users").
		Where("email = ? AND status_data = true", email).Count(&count)
	return count > 0, result.Error
}

func (r *AccountRepository) ExistsPhone(ctx context.Context, telepon string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Table("laundry.users").
		Where("telepon = ? AND status_data = true", telepon).Count(&count)
	return count > 0, result.Error
}

func (r *AccountRepository) CreateUser(ctx context.Context, nama, email, telepon, passwordHash, alamat string) (string, error) {
	var idUser string
	result := r.db.WithContext(ctx).Raw(`
		INSERT INTO laundry.users (nama, email, telepon, password, alamat)
		VALUES (?, ?, NULLIF(?, ''), ?, ?)
		RETURNING id_user
	`, nama, email, telepon, passwordHash, alamat).Scan(&idUser)
	return idUser, result.Error
}

func (r *AccountRepository) UpsertPasswordReset(ctx context.Context, email, otp, tokenHash string, sentAt time.Time) error {
	if err := r.db.WithContext(ctx).Exec(`DELETE FROM laundry.password_resets WHERE email = ?`, email).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO laundry.password_resets (email, otp, token_hash, otp_sent_at)
		VALUES (?, ?, ?, ?)
	`, email, otp, tokenHash, sentAt).Error
}

func (r *AccountRepository) FindPasswordResetByEmail(ctx context.Context, email string) (map[string]interface{}, error) {
	var data map[string]interface{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT id_password_reset, email, otp, token_hash, otp_sent_at, otp_verified
		FROM laundry.password_resets
		WHERE email = ?
	`, email).First(&data)
	return data, result.Error
}

func (r *AccountRepository) DeletePasswordResetByEmail(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).Exec(`DELETE FROM laundry.password_resets WHERE email = ?`, email).Error
}

func (r *AccountRepository) UpdateUserPassword(ctx context.Context, idUser, passwordHash string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE laundry.users SET password = ?, tgl_update = NOW() WHERE id_user = ?
	`, passwordHash, idUser).Error
}
