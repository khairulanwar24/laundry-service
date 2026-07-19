package models

import (
	"fmt"
	"math/rand"
	"regexp"
	"sso-service/database"
	middleware "sso-service/middlewares"
	"sso-service/types"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u *User) TableName() string {
	return "users"
}

func Login(username, password string) types.Response {
	var resp types.Response
	var user map[string]interface{}

	result := database.DB.Raw("select username,nama_lengkap,email,no_hp,avatar,id_user,id_person,first_login,password,jenis_user from users where username = ? ", username).First(&user)

	if result.Error != nil {
		resp.Success = false
		resp.Message = "username or password is incorrect"
		resp.Data = nil
		return resp
	}

	hash := user["password"].(string)

	match := middleware.CheckPasswordHash(password, hash)

	if !match {
		resp.Success = false
		resp.Message = "username or password is incorrect"
		resp.Data = nil
	} else {

		delete(user, "password")

		resp.Success = true
		resp.Message = "Login Success"
		resp.Data = user
	}

	return resp
}

func isValidEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

type UserBlock struct {
	StillBlocked     bool    `json:"still_blocked"`
	SecondsRemaining float64 `json:"seconds_remaining"`
	IdUser           string  `json:"id_user"`
	Email            string  `json:"email"`
}

func ResetPassword(username, uuid string) types.Response {
	var resp types.Response
	var user UserBlock
	// var insertUser []map[string]interface{}

	result := database.DB.Raw("select id_user,email,tgl_lock > NOW() AS still_blocked, EXTRACT(EPOCH FROM (tgl_lock - NOW())) AS seconds_remaining from users where username = ? ", username).Scan(&user)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Seed the random number generator
	randomNumber := r.Intn(900000) + 100000       // Generate a random 6-digit number
	randomNumberStr := strconv.Itoa(randomNumber) // Convert the number to a string

	// untuk menipu hacker jangan tampilkan info kalau username valid , semua respon kembalikan true
	if result.Error != nil || user.IdUser == "" {
		// buat menangani maksimal input otp 5 kali
		database.DB.Exec(`INSERT INTO password_reset (id_password_reset,username, otp, tgl_insert) VALUES (?,?, ?, ?)`, uuid, username, randomNumberStr, time.Now())

		resp.Success = true
		// resp.Message = result.Error.Error()
		// resp.Message = "username is incorrect"
		resp.Message = "Silahkan Cek Email Anda, Token akan dikirim melalui Email"
		resp.Data = uuid
		return resp
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

		resp.Success = false
		resp.Message = fmt.Sprintf("Akun Anda terkunci selama %s  lagi. Silahkan coba lagi dalam %s kembali", waktu, waktu)
		resp.Data = uuid
		return resp

	}

	id_user := user.IdUser
	email := user.Email

	// fmt.Println("Random 6-digit number:", randomNumberStr)

	result = database.DB.Exec(`Update password_reset set status_data = false where id_user= ?`, id_user)
	if result.Error != nil {
		database.DB.Exec(`INSERT INTO password_reset (id_password_reset,username, otp, tgl_insert) VALUES (?,?, ?, ?)`, uuid, username, randomNumberStr, time.Now())
		resp.Success = true
		// resp.Message = "Failed to reset password"
		resp.Message = "Silahkan Cek Email Anda, Token akan dikirim melalui Email"
		resp.Data = uuid
		return resp
	}

	result = database.DB.Exec(`INSERT INTO password_reset (id_password_reset,username,  id_user, otp, tgl_insert) VALUES (?,?, ?, ?, ?)`, uuid, username, id_user, randomNumberStr, time.Now())
	if result.Error != nil {
		database.DB.Exec(`INSERT INTO password_reset (id_password_reset,username, otp, tgl_insert) VALUES (?,?, ?, ?)`, uuid, username, randomNumberStr, time.Now())
		resp.Success = true
		resp.Message = "Silahkan Cek Email Anda, Token akan dikirim melalui Email"
		// resp.Message = "Failed to reset password"
		resp.Data = uuid
		return resp
	}

	// result = database.DB.Create(&reset)

	if result.RowsAffected == 0 {
		database.DB.Exec(`INSERT INTO password_reset (id_password_reset,username, otp, tgl_insert) VALUES (?,?, ?, ?)`, uuid, username, randomNumberStr, time.Now())
		resp.Success = true
		resp.Message = "Silahkan Cek Email Anda, Token akan dikirim melalui Email"
		// resp.Message = "Failed to reset password"
		resp.Data = uuid
		return resp
	} else {

		// randomNumberStr := "319935" // Example OTP
		to := email
		subject := "Reset Password"
		body := "Kode OTP Anda = " + randomNumberStr + " Tolong Jangan Bagi OTP Ini Dengan Orang Lain. OTP ini berlaku selama 10 menit. Thank you."

		validemail := isValidEmailFormat(to)

		if !validemail {
			resp.Success = false
			resp.Message = "Silahkan Cek Email Anda, Token akan dikirim melalui Email"
			// resp.Message = "Email tidak valid"
			resp.Data = uuid
			return resp
		}

		_, err := middleware.Mail(to, subject, body)

		if err != nil {
			resp.Success = false
			resp.Message = "Email tidak valid"
			resp.Data = uuid
			return resp
		}

		resp.Success = true
		// resp.Message = "Rows affected:" + strconv.FormatInt(result.RowsAffected, 10)
		// resp.Message = "Email Berhasil dikirim ke " + maskEmail(email)
		resp.Message = "Silahkan Cek Email Anda, Token akan dikirim melalui Email"
		resp.Data = uuid

		// _, err := middleware.WhatsApp("628985576111", "Your OTP is: "+randomNumberStr+" \n\nPlease do not share your OTP with anyone. \n\nThank you.")
		// if err != nil {
		// 	t.Fatalf("Failed to publish message: %v", err)
		// }
		// test := middleware.PublishMessageToQueue(
		// 	"dhanang.hadiyanto@gmail.com",
		// 	"Reset Password",
		// 	"Your OTP is: "+randomNumberStr+" \n\nPlease do not share your OTP with anyone. \n\nThank you.")
		// fmt.Println(message, err)
	}

	return resp
}

func maskEmail(email string) string {
	// Split the email into username and domain
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // return original if not a valid email format
	}
	username, domain := parts[0], parts[1]

	// Mask the username and return the masked email
	if len(username) > 1 {
		return fmt.Sprintf("%s****@%s", string(username[0]), domain)
	}
	return fmt.Sprintf("****@%s", domain)
}

func CekOtp(id_password_reset, otp string) types.Response {
	var resp types.Response
	var user []map[string]interface{}
	// var insertUser []map[string]interface{}

	result := database.DB.Raw("select username,id_user,tgl_insert,percobaan , otp from password_reset where id_password_reset = ? and status_data = true ", id_password_reset).First(&user)

	if result.Error != nil {
		resp.Success = false
		// resp.Message = result.Error.Error()
		resp.Message = "OTP Salah, Silahkan Coba Lagi"
		resp.Data = nil
		return resp
	}

	percobaan := user[0]["percobaan"].(int32)
	percobaan = percobaan + 1
	result = database.DB.Exec(`Update password_reset set percobaan = ? where id_password_reset = ? and status_data = true`, percobaan, id_password_reset)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Failed Update DB password reset"
		resp.Data = nil
		return resp
	}
	// fmt.Println("ini user", user[0]["id_user"])

	if user[0]["id_user"] == nil {
		resp.Success = false
		resp.Message = "OTP Salah, Silahkan Coba Lagi"
		resp.Data = nil
		return resp
	}

	tgl_insert := user[0]["tgl_insert"].(time.Time)
	id_user := user[0]["id_user"].(string)
	// fmt.Println("mencari percobaan")
	if percobaan > 5 {
		if id_user != "" {
			database.DB.Exec(`Update users set tgl_lock =  NOW() + INTERVAL '30 minutes' where id_user = ? and status_data = true`, id_user)
		}

		resp.Success = false
		resp.Message = "Failed: Limit percobaan OTP telah terpenuhi, Silahkan Coba Lagi Dalam 30 Menit"
		resp.Data = nil
		return resp
	}

	otpdb := user[0]["otp"].(string)

	if otp != otpdb {
		resp.Success = false
		resp.Message = "OTP Failed"
		resp.Data = nil
		return resp
	}

	now := time.Now()
	diff := now.Sub(tgl_insert)

	if diff > 10*time.Minute {
		resp.Success = false
		resp.Message = "OTP Expired"
		resp.Data = nil
		return resp
	}

	username := user[0]["username"].(string)
	token, code := middleware.AccessToken(username, id_user)
	if code != 200 {
		resp.Success = false
		resp.Message = "error generate token"
		resp.Data = nil
		return resp
	}

	respon := make(map[string]interface{})

	respon["token"] = token

	resp.Success = true
	resp.Message = "OTP Valid"
	resp.Data = nil
	return resp
}

func GenerateTokenChangePassword(id_password_reset string) string {
	var user []map[string]interface{}

	result := database.DB.Raw("select username,id_user,tgl_insert,percobaan , otp from password_reset where id_password_reset = ? and status_data = true ", id_password_reset).First(&user)

	username := user[0]["username"].(string)
	id_user := user[0]["id_user"].(string)

	token, _ := middleware.AccessTokenOTP(username, id_user)
	if result.Error != nil {
		return ""
	}

	return token

}

func CekPercobaan(id_password_reset string) int32 {
	var user []map[string]interface{}

	result := database.DB.Raw("select username,id_user,tgl_insert,percobaan , otp from password_reset where id_password_reset = ? and status_data = true ", id_password_reset).First(&user)

	if result.Error != nil {
		return 0
	}

	// username := user[0]["username"].(string)
	percobaan := user[0]["percobaan"].(int32)

	return percobaan

}

func ChangePassword(username, password, token string) types.Response {
	var resp types.Response

	parsetoken, err := middleware.ParseTokenChangePassword(token)

	// fmt.Println(parsetoken)
	if err != nil {
		resp.Success = false
		resp.Message = "Invalid Token" + err.Error()
		resp.Data = nil
		return resp
	}

	id_user := parsetoken.Id_user

	// fmt.Println(id_user)
	var user []map[string]interface{}

	result := database.DB.Raw(`SELECT
								id_user
								, first_login
								FROM  users where id_user = ? and status_data = true`, id_user).First(&user)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Mendapatkan users"
		resp.Data = nil
		return resp
	}
	password_hash, _ := middleware.HashPassword(password)
	first_login := user[0]["first_login"].(bool)
	if first_login {
		result = database.DB.Exec(`UPDATE users
								SET password = ?,tgl_update = NOW(),first_login = false,tgl_first_login = NOW()
								WHERE id_user = ?`, password_hash, id_user)

	} else {
		result = database.DB.Exec(`UPDATE users
								SET password = ?,tgl_update = NOW()
								WHERE id_user = ?`, password_hash, id_user)
	}

	// panic(result.RowsAffected)
	if result.Error != nil {
		resp.Success = false
		resp.Message = "Gagal Update users"
		resp.Data = nil
		return resp
	} else if result.RowsAffected == 0 {
		resp.Success = false
		resp.Message = "users tidak ada"
		resp.Data = nil
		return resp
	}

	resp.Success = true
	resp.Message = "Success"
	resp.Data = nil

	return resp
}

// func CreateSession(username string) types.Response {
// 	var resp types.Response
// 	var user []map[string]interface{}
// 	// var insertUser []map[string]interface{}
