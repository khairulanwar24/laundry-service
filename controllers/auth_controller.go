package controllers

import (
	"os"
	middleware "sso-service/middlewares"
	"sso-service/models"
	"sso-service/types"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MyClaims struct {
	jwt.RegisteredClaims
	// Add custom claims here if necessary
}

type LoginForm struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password"  form:"password" validate:"required"`
}

func LoginController(c *fiber.Ctx) error {

	form := c.Locals("validatedForm").(*LoginForm)

	data := models.Login(form.Username, form.Password)

	if !data.Success {
		return c.JSON(data)
	}

	user := data.Data.(map[string]interface{})

	resp := map[string]interface{}{}

	resp["user"] = user

	RefreshToken, _ := middleware.RefreshToken(user["username"].(string), user["id_user"].(string))
	setHttpCookie(c, RefreshToken)
	// setCookieHandler(RefreshToken)

	// fmt.Println(tes)

	AccessToken, _ := middleware.AccessToken(user["username"].(string), user["id_user"].(string))

	resp["access_token"] = AccessToken

	respJson := types.Response{
		Success: data.Success,
		Message: data.Message,
		Data:    resp,
	}

	// cookieRefreshToken := fiber.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    RefreshToken,
	// 	Expires:  time.Now().Add(time.Hour * 8),
	// 	HTTPOnly: true,
	// 	// SameSite: "lax",
	// 	Secure: true,
	// 	Attributes: map[string]string{
	// 		"Partitioned": "true", // This is for demonstration; check compatibility in your environment
	// 	},
	// 	// Domain:   ".dnglab.id",
	// 	// Domain:   "localhost",
	// }
	// c.Cookie(&cookieRefreshToken)

	// app := fiber.New()

	// app.Get("/set-partitioned-cookie", func(c *fiber.Ctx) error {
	// 	cookie := new(fiber.Cookie)
	// 	cookie.Name = "refresh_token"
	// 	cookie.Value = "your-token-value-here"
	// 	cookie.Expires = time.Now().Add(24 * time.Hour) // Expiry time for the cookie
	// 	cookie.HTTPOnly = true                          // Prevent JavaScript access
	// 	cookie.Secure = true                            // Send cookie only over HTTPS
	// 	cookie.SameSite = "None"                        // Allow cross-site usage
	// 	cookie.Domain = "api.dnglab.id"                 // Set the correct domain
	// 	cookie.Path = "/"                               // Path for the cookie

	// 	// Add the Partitioned attribute
	// 	cookie.Attributes = map[string]string{
	// 		"Partitioned": "true", // This is for demonstration; check compatibility in your environment
	// 	}

	// 	// Set the cookie
	// 	c.Cookie(cookie)
	// 	return c.SendString("Partitioned cookie set")
	// })

	return c.JSON(respJson)
}

func LogoutController(c *fiber.Ctx) error {

	cookieRefreshToken := fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   true,
		// Domain:   "localhost",
	}
	c.Cookie(&cookieRefreshToken)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully, session cookie deleted!",
	})
}

type ResetForm struct {
	Username string `json:"username" form:"username" validate:"required"`
}

func ResetPassword(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*ResetForm)

	uuid := uuid.New().String()

	data := models.ResetPassword(form.Username, uuid)

	if data.Success {
		OTP, _ := middleware.OTPResetpassword(uuid, 0)
		setOTPCookie(c, OTP, "OTP")
	}

	return c.JSON(data)
}

type ChangePasswordForm struct {
	Username string `json:"username" form:"username" validate:""`
	Password string `json:"password"  form:"password" validate:"required"`
}

func ChangePassword(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*ChangePasswordForm)

	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing or invalid Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// fmt.p

	data := models.ChangePassword(form.Username, form.Password, tokenString)

	return c.JSON(data)
}

type OtpForm struct {
	// Username string `json:"username" form:"username" validate:"required"`
	Otp string `json:"otp" form:"otp" validate:"required,numeric"`
}

func CekOtp(c *fiber.Ctx) error {
	var resp types.Response

	form := c.Locals("validatedForm").(*OtpForm)

	token := middleware.CheckjwtOTP(c)
	// fmt.Println(token)

	// return c.JSON(resp)
	if !token {
		resp.Success = false
		resp.Message = "OTP Salah, Silahkan Coba Lagi"
		resp.Data = nil
		return c.JSON(resp)
	}

	jwtotp, err := middleware.ParseTokenOTP(c)

	if err != nil {
		resp.Success = false
		resp.Message = "OTP Salah, Silahkan Coba Lagi"
		resp.Data = nil
		return c.JSON(resp)
	}

	IDResetPassword := jwtotp.IDResetPassword
	otp := form.Otp

	// fmt.Println(IDResetPassword, otp)

	data := models.CekOtp(IDResetPassword, otp)

	// fmt.Println(data)

	if data.Success {

		token := models.GenerateTokenChangePassword(IDResetPassword)

		setOTPCookie(c, token, "access_token")
	} else {

		cekpercobaan := models.CekPercobaan(IDResetPassword)

		OTP, _ := middleware.OTPResetpassword(IDResetPassword, cekpercobaan)

		setOTPCookie(c, OTP, "OTP")
	}

	return c.JSON(data)
}

// GenerateUserMahasiswa/Dosen/Tendik dipindah ke domain user (controllers/user, services/user).

func RefreshToken(c *fiber.Ctx) error {
	var resp types.Response
	// c.Cookies("refresh_token")

	// refreshtoken := c.Cookies("refresh_token")
	// fmt.Println(refreshtoken)
	// panic(refreshtoken)

	// return c.JSON(refreshtoken)

	token := middleware.CheckjwtRefresh(c)
	// return c.JSON(resp)
	if !token {
		resp.Success = false
		resp.Message = "unauthenticated"
		resp.Data = nil
	} else {
		resp = middleware.JWTParse(c)
	}
	return c.JSON(resp)
}

// func setCookieHandler(tokenValue string) http.HandlerFunc {
// 	fmt.Println(tokenValue)
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("ini")
// 		// Create a new cookie
// 		cookie := &http.Cookie{
// 			Name:        "refresh_token",
// 			Value:       tokenValue,
// 			Expires:     time.Now().Add(24 * time.Hour), // Expiry time for the cookie
// 			HttpOnly:    true,                           // Prevent JavaScript access
// 			Secure:      true,                           // Send cookie only over HTTPS
// 			SameSite:    http.SameSiteNoneMode,          // Allow cross-site usage
// 			Domain:      ".farmasiunissula.com",         // Set the correct domain
// 			Path:        "/",                            // Path for the cookie
// 			Partitioned: true,                           // This is for demonstration; check compatibility in your environment
// 		}

// 		// Set the cookie
// 		http.SetCookie(w, cookie)

// 		// Optionally, set the cookie with the Partitioned attribute in the response header
// 		w.Header().Set("Set-Cookie", "refresh_token="+tokenValue+"; SameSite=None; Secure; Partitioned")

// 		w.Write([]byte("Refresh token cookie set"))
// 	}
// }

// func setHttpCookie(c *fiber.Ctx, tokenValue string) {
// 	// Create a new cookie using net/http
// 	cookieRefreshToken := fiber.Cookie{
// 		Name:     "refresh_token",
// 		Value:    tokenValue,
// 		Expires:  time.Now().Add(time.Hour * 8),
// 		HTTPOnly: true,
// 		SameSite: "None",
// 		Secure:   true,
// 		// Attributes: map[string]string{
// 		// 	"Partitioned": "true", // This is for demonstration; check compatibility in your environment
// 		// },
// 		Domain: ".farmasiunissula.com",
// 		// Domain:   "localhost",
// 	}
// 	c.Cookie(&cookieRefreshToken)

// }

func setHttpCookie(c *fiber.Ctx, tokenValue string) {
	appEnv := os.Getenv("APP_ENV") // pastikan sudah dipanggil godotenv.Load()
	isProd := appEnv == "production"

	var domain string
	secure := false
	sameSite := "Lax"

	if isProd {
		// ⛳ Jika production, pakai domain dan secure tetap
		domain = ".farmasiunissula.com"
		secure = true
		sameSite = "None"
	} else {
		// 🧪 Jika development, fleksibel
		isHTTPS := c.Protocol() == "https"
		host := c.Hostname()
		isLocal := strings.Contains(host, "localhost") || strings.HasPrefix(host, "127.") || strings.HasSuffix(host, ".test")

		if isLocal {
			domain = "localhost"
		} else if strings.HasSuffix(host, ".farmasiunissula.com") {
			domain = ".farmasiunissula.com"
		} else {
			domain = host // fallback
		}

		secure = isHTTPS
		if secure {
			sameSite = "None"
		}
	}

	// Buat cookie
	cookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokenValue,
		Expires:  time.Now().Add(8 * time.Hour),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
		Domain:   domain,
	}

	c.Cookie(&cookie)
}

func setOTPCookie(c *fiber.Ctx, tokenValue string, nama string) {
	appEnv := os.Getenv("APP_ENV") // pastikan sudah dipanggil godotenv.Load()
	isProd := appEnv == "production"

	var domain string
	secure := false
	sameSite := "Lax"

	if isProd {
		// ⛳ Jika production, pakai domain dan secure tetap
		domain = ".farmasiunissula.com"
		secure = true
		sameSite = "None"
	} else {
		// 🧪 Jika development, fleksibel
		isHTTPS := c.Protocol() == "https"
		host := c.Hostname()
		isLocal := strings.Contains(host, "localhost") || strings.HasPrefix(host, "127.") || strings.HasSuffix(host, ".test")

		if isLocal {
			domain = "localhost"
		} else if strings.HasSuffix(host, ".farmasiunissula.com") {
			domain = ".farmasiunissula.com"
		} else {
			domain = host // fallback
		}

		secure = isHTTPS
		if secure {
			sameSite = "None"
		}
	}

	// Buat cookie
	cookie := fiber.Cookie{
		Name:     nama,
		Value:    tokenValue,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		Path:     "/",
		Domain:   domain,
	}

	c.Cookie(&cookie)
}

// digunakanm jika dibutuhkan logout session dari database
// func CreateSession(c *fiber.Ctx) error {
// 	var resp types.Response

// 	return c.JSON(resp)
// }

// func DeleteSession(c *fiber.Ctx) error {
// 	var resp types.Response
// 	return c.JSON(resp)
// }

// func GetSession(c *fiber.Ctx) error {
// 	var resp types.Response
// 	return c.JSON(resp)
// }

// func RevokeSession(c *fiber.Ctx) error {
// 	var resp types.Response
// 	return c.JSON(resp)
// }
