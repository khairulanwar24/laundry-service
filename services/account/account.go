// Package services (account) berisi logika bisnis domain akun laundry:
// registrasi, login, profil, dan reset password via OTP email.
package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"laundry-service/common/response"
	"laundry-service/domain/dto"
	middleware "laundry-service/middlewares"
	"laundry-service/repositories"
)

const otpValidDuration = 10 * time.Minute

// AccountService membungkus akses ke repository registry.
type AccountService struct {
	repository repositories.IRepositoryRegistry
}

// IAccountService adalah kontrak logika bisnis domain akun.
type IAccountService interface {
	Register(ctx context.Context, form dto.RegisterForm) response.Response
	Login(ctx context.Context, form dto.LoginForm) response.Response
	Me(ctx context.Context, idUser string) response.Response
	Logout(ctx context.Context) response.Response
	ForgotPassword(ctx context.Context, form dto.ForgotPasswordForm) response.Response
	ResetPassword(ctx context.Context, form dto.ResetPasswordForm) response.Response
}

// NewAccountService membuat instance AccountService baru.
func NewAccountService(repository repositories.IRepositoryRegistry) IAccountService {
	return &AccountService{repository: repository}
}

func generateOTP() (string, error) {
	max := 1000000
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n := (int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])) % max
	if n < 0 {
		n = -n
	}
	return fmt.Sprintf("%06d", n), nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Register mendaftarkan akun laundry baru (pemilik/karyawan outlet).
func (s *AccountService) Register(ctx context.Context, form dto.RegisterForm) response.Response {
	repo := s.repository.GetAccount()

	existsEmail, err := repo.ExistsEmail(ctx, form.Email)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal memeriksa email: " + err.Error()}
	}
	if existsEmail {
		return response.Response{Success: false, Message: "Email sudah terdaftar"}
	}

	if form.Telepon != "" {
		existsPhone, err := repo.ExistsPhone(ctx, form.Telepon)
		if err != nil {
			return response.Response{Success: false, Message: "Gagal memeriksa telepon: " + err.Error()}
		}
		if existsPhone {
			return response.Response{Success: false, Message: "Telepon sudah terdaftar"}
		}
	}

	passwordHash, err := middleware.HashPassword(form.Password)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal memproses password"}
	}

	idUser, err := repo.CreateUser(ctx, form.Nama, form.Email, form.Telepon, passwordHash, form.Alamat)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal mendaftarkan akun: " + err.Error()}
	}

	accessToken, _ := middleware.AccessToken(form.Email, idUser)
	refreshToken, _ := middleware.RefreshToken(form.Email, idUser)

	return response.Response{
		Success: true,
		Message: "Registrasi berhasil",
		Data: map[string]interface{}{
			"id_user":       idUser,
			"nama":          form.Nama,
			"email":         form.Email,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	}
}

// Login memverifikasi email & password akun laundry.
func (s *AccountService) Login(ctx context.Context, form dto.LoginForm) response.Response {
	user, err := s.repository.GetAccount().FindUserByEmail(ctx, form.Email)
	if err != nil {
		return response.Response{Success: false, Message: "Email atau password salah"}
	}

	hash, _ := user["password"].(string)
	if !middleware.CheckPasswordHash(form.Password, hash) {
		return response.Response{Success: false, Message: "Email atau password salah"}
	}

	isActive, _ := user["is_active"].(bool)
	if !isActive {
		return response.Response{Success: false, Message: "Akun tidak aktif"}
	}

	idUser, _ := user["id_user"].(string)
	accessToken, _ := middleware.AccessToken(form.Email, idUser)
	refreshToken, _ := middleware.RefreshToken(form.Email, idUser)

	delete(user, "password")
	return response.Response{
		Success: true,
		Message: "Login berhasil",
		Data: map[string]interface{}{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	}
}

// Me mengembalikan profil akun yang sedang login.
func (s *AccountService) Me(ctx context.Context, idUser string) response.Response {
	user, err := s.repository.GetAccount().FindUserByID(ctx, idUser)
	if err != nil {
		return response.Response{Success: false, Message: "Akun tidak ditemukan"}
	}
	return response.Response{Success: true, Message: "sukses", Data: user}
}

// Logout tidak melakukan revoke di server (JWT stateless) — client cukup
// membuang token yang tersimpan.
func (s *AccountService) Logout(ctx context.Context) response.Response {
	return response.Response{Success: true, Message: "Logout berhasil"}
}

// ForgotPassword membuat OTP + token reset password dan mengirimkannya ke email.
func (s *AccountService) ForgotPassword(ctx context.Context, form dto.ForgotPasswordForm) response.Response {
	repo := s.repository.GetAccount()

	user, err := repo.FindUserByEmail(ctx, form.Email)
	if err != nil {
		// Jangan bocorkan apakah email terdaftar atau tidak.
		return response.Response{Success: true, Message: "Jika email terdaftar, OTP telah dikirim"}
	}
	if isActive, _ := user["is_active"].(bool); !isActive {
		return response.Response{Success: true, Message: "Jika email terdaftar, OTP telah dikirim"}
	}

	otp, err := generateOTP()
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat OTP"}
	}
	token, err := generateToken()
	if err != nil {
		return response.Response{Success: false, Message: "Gagal membuat token"}
	}
	tokenHash, err := middleware.HashPassword(token)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal memproses token"}
	}

	if err := repo.UpsertPasswordReset(ctx, form.Email, otp, tokenHash, time.Now()); err != nil {
		return response.Response{Success: false, Message: "Gagal menyimpan permintaan reset password"}
	}

	subject := "Reset Password Laundry Service"
	body := fmt.Sprintf("Kode OTP Anda: %s. Berlaku %d menit. Jangan bagikan kode ini ke siapa pun.", otp, int(otpValidDuration.Minutes()))
	if _, err := middleware.Mail(form.Email, subject, body); err != nil {
		return response.Response{Success: false, Message: "Gagal mengirim email OTP"}
	}

	return response.Response{
		Success: true,
		Message: "Jika email terdaftar, OTP telah dikirim",
		Data: map[string]interface{}{
			"token":          token,
			"otp_expires_in": int(otpValidDuration.Seconds()),
		},
	}
}

// ResetPassword mengganti password memakai token + OTP dari ForgotPassword.
func (s *AccountService) ResetPassword(ctx context.Context, form dto.ResetPasswordForm) response.Response {
	repo := s.repository.GetAccount()

	reset, err := repo.FindPasswordResetByEmail(ctx, form.Email)
	if err != nil {
		return response.Response{Success: false, Message: "Permintaan reset password tidak ditemukan"}
	}

	tokenHash, _ := reset["token_hash"].(string)
	if !middleware.CheckPasswordHash(form.Token, tokenHash) {
		return response.Response{Success: false, Message: "Token tidak valid"}
	}

	otpDB, _ := reset["otp"].(string)
	if form.Otp != otpDB {
		return response.Response{Success: false, Message: "OTP tidak valid"}
	}

	otpSentAt, _ := reset["otp_sent_at"].(time.Time)
	if time.Since(otpSentAt) > otpValidDuration {
		_ = repo.DeletePasswordResetByEmail(ctx, form.Email)
		return response.Response{Success: false, Message: "OTP sudah kedaluwarsa"}
	}

	user, err := repo.FindUserByEmail(ctx, form.Email)
	if err != nil {
		return response.Response{Success: false, Message: "Akun tidak ditemukan"}
	}
	idUser, _ := user["id_user"].(string)

	passwordHash, err := middleware.HashPassword(form.Password)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal memproses password"}
	}

	if err := repo.UpdateUserPassword(ctx, idUser, passwordHash); err != nil {
		return response.Response{Success: false, Message: "Gagal mengubah password"}
	}
	_ = repo.DeletePasswordResetByEmail(ctx, form.Email)

	return response.Response{Success: true, Message: "Password berhasil diubah"}
}
