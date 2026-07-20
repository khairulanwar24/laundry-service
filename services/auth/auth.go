// Package services (auth) berisi logika bisnis domain autentikasi:
// login, reset password, verifikasi OTP, dan ganti password.
package services

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"laundry-service/common/response"
	middleware "laundry-service/middlewares"
	"laundry-service/repositories"
)

// AuthService membungkus akses ke repository registry.
type AuthService struct {
	repository repositories.IRepositoryRegistry
}

// IAuthService adalah kontrak logika bisnis domain autentikasi.
type IAuthService interface {
	Login(ctx context.Context, username, password string) response.Response
	ResetPassword(ctx context.Context, username, uuid string) response.Response
	CekOtp(ctx context.Context, idPasswordReset, otp string) response.Response
	GenerateTokenChangePassword(ctx context.Context, idPasswordReset string) string
	CekPercobaan(ctx context.Context, idPasswordReset string) int32
	ChangePassword(ctx context.Context, username, password, token string) response.Response
}

// NewAuthService membuat instance AuthService baru.
func NewAuthService(repository repositories.IRepositoryRegistry) IAuthService {
	return &AuthService{repository: repository}
}

func isValidEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Login memverifikasi username & password.
func (s *AuthService) Login(ctx context.Context, username, password string) response.Response {
	user, err := s.repository.GetAuth().FindUserByUsername(ctx, username)
	if err != nil {
		return response.Response{Success: false, Message: "username or password is incorrect"}
	}

	hash := user["password"].(string)
	match := middleware.CheckPasswordHash(password, hash)
	if !match {
		return response.Response{Success: false, Message: "username or password is incorrect"}
	}

	delete(user, "password")
	return response.Response{Success: true, Message: "Login Success", Data: user}
}

// ResetPassword membuat OTP reset password & mengirim email.
func (s *AuthService) ResetPassword(ctx context.Context, username, uuid string) response.Response {
	repo := s.repository.GetAuth()
	user, err := repo.FindUserBlockByUsername(ctx, username)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := r.Intn(900000) + 100000
	randomNumberStr := strconv.Itoa(randomNumber)

	// untuk menipu hacker jangan tampilkan info kalau username valid, semua respon kembalikan true
	if err != nil || user.IdUser == "" {
		_ = repo.InsertPasswordResetNoUser(ctx, uuid, username, randomNumberStr, time.Now())
		return response.Response{Success: true, Message: "Silahkan Cek Email Anda, Token akan dikirim melalui Email", Data: uuid}
	}

	if user.StillBlocked {
		totalSeconds := int(user.SecondsRemaining)
		hours := totalSeconds / 3600
		minutes := (totalSeconds % 3600) / 60

		var waktu string
		if hours > 0 {
			waktu = fmt.Sprintf("%d jam %d menit", hours, minutes)
		} else {
			waktu = fmt.Sprintf("%d menit", minutes)
		}

		return response.Response{Success: false, Message: fmt.Sprintf("Akun Anda terkunci selama %s  lagi. Silahkan coba lagi dalam %s kembali", waktu, waktu), Data: uuid}
	}

	idUser := user.IdUser
	email := user.Email

	if err := repo.DeactivatePasswordResetByUser(ctx, idUser); err != nil {
		_ = repo.InsertPasswordResetNoUser(ctx, uuid, username, randomNumberStr, time.Now())
		return response.Response{Success: true, Message: "Silahkan Cek Email Anda, Token akan dikirim melalui Email", Data: uuid}
	}

	rows, err := repo.InsertPasswordResetWithUser(ctx, uuid, username, idUser, randomNumberStr, time.Now())
	if err != nil {
		_ = repo.InsertPasswordResetNoUser(ctx, uuid, username, randomNumberStr, time.Now())
		return response.Response{Success: true, Message: "Silahkan Cek Email Anda, Token akan dikirim melalui Email", Data: uuid}
	}

	if rows == 0 {
		_ = repo.InsertPasswordResetNoUser(ctx, uuid, username, randomNumberStr, time.Now())
		return response.Response{Success: true, Message: "Silahkan Cek Email Anda, Token akan dikirim melalui Email", Data: uuid}
	}

	to := email
	subject := "Reset Password"
	body := "Kode OTP Anda = " + randomNumberStr + " Tolong Jangan Bagi OTP Ini Dengan Orang Lain. OTP ini berlaku selama 10 menit. Thank you."

	if !isValidEmailFormat(to) {
		return response.Response{Success: false, Message: "Silahkan Cek Email Anda, Token akan dikirim melalui Email", Data: uuid}
	}

	if _, err := middleware.Mail(to, subject, body); err != nil {
		return response.Response{Success: false, Message: "Email tidak valid", Data: uuid}
	}

	return response.Response{Success: true, Message: "Silahkan Cek Email Anda, Token akan dikirim melalui Email", Data: uuid}
}

// CekOtp memverifikasi OTP reset password.
func (s *AuthService) CekOtp(ctx context.Context, idPasswordReset, otp string) response.Response {
	repo := s.repository.GetAuth()
	user, err := repo.FindPasswordReset(ctx, idPasswordReset)
	if err != nil {
		return response.Response{Success: false, Message: "OTP Salah, Silahkan Coba Lagi"}
	}

	percobaan := user[0]["percobaan"].(int32)
	percobaan = percobaan + 1
	if err := repo.UpdatePasswordResetPercobaan(ctx, percobaan, idPasswordReset); err != nil {
		return response.Response{Success: false, Message: "Failed Update DB password reset"}
	}

	if user[0]["id_user"] == nil {
		return response.Response{Success: false, Message: "OTP Salah, Silahkan Coba Lagi"}
	}

	tglInsert := user[0]["tgl_insert"].(time.Time)
	idUser := user[0]["id_user"].(string)

	if percobaan > 5 {
		if idUser != "" {
			_ = repo.LockUser(ctx, idUser)
		}
		return response.Response{Success: false, Message: "Failed: Limit percobaan OTP telah terpenuhi, Silahkan Coba Lagi Dalam 30 Menit"}
	}

	otpdb := user[0]["otp"].(string)
	if otp != otpdb {
		return response.Response{Success: false, Message: "OTP Failed"}
	}

	now := time.Now()
	if now.Sub(tglInsert) > 10*time.Minute {
		return response.Response{Success: false, Message: "OTP Expired"}
	}

	username := user[0]["username"].(string)
	_, code := middleware.AccessToken(username, idUser)
	if code != 200 {
		return response.Response{Success: false, Message: "error generate token"}
	}

	return response.Response{Success: true, Message: "OTP Valid"}
}

// GenerateTokenChangePassword membuat token khusus untuk ganti password.
func (s *AuthService) GenerateTokenChangePassword(ctx context.Context, idPasswordReset string) string {
	user, err := s.repository.GetAuth().FindPasswordReset(ctx, idPasswordReset)
	if err != nil {
		return ""
	}

	username := user[0]["username"].(string)
	idUser := user[0]["id_user"].(string)

	token, _ := middleware.AccessTokenOTP(username, idUser)
	return token
}

// CekPercobaan mengambil jumlah percobaan OTP.
func (s *AuthService) CekPercobaan(ctx context.Context, idPasswordReset string) int32 {
	user, err := s.repository.GetAuth().FindPasswordReset(ctx, idPasswordReset)
	if err != nil {
		return 0
	}

	return user[0]["percobaan"].(int32)
}

// ChangePassword mengganti password user berdasarkan token yang valid.
func (s *AuthService) ChangePassword(ctx context.Context, username, password, token string) response.Response {
	repo := s.repository.GetAuth()

	parsetoken, err := middleware.ParseTokenChangePassword(token)
	if err != nil {
		return response.Response{Success: false, Message: "Invalid Token" + err.Error()}
	}

	idUser := parsetoken.Id_user

	user, err := repo.FindUserFirstLogin(ctx, idUser)
	if err != nil {
		return response.Response{Success: false, Message: "Gagal Mendapatkan users"}
	}

	passwordHash, _ := middleware.HashPassword(password)
	firstLogin := user[0]["first_login"].(bool)

	var rows int64
	if firstLogin {
		rows, err = repo.UpdateUserPasswordFirstLogin(ctx, passwordHash, idUser)
	} else {
		rows, err = repo.UpdateUserPasswordNormal(ctx, passwordHash, idUser)
	}

	if err != nil {
		return response.Response{Success: false, Message: "Gagal Update users"}
	} else if rows == 0 {
		return response.Response{Success: false, Message: "users tidak ada"}
	}

	return response.Response{Success: true, Message: "Success"}
}
