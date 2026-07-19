// Package controllers (auth) adalah lapisan HTTP handler (Fiber) untuk domain autentikasi.
package controllers

import (
	"os"
	"strings"
	"time"

	"sso-service/common/response"
	"sso-service/domain/dto"
	middleware "sso-service/middlewares"
	"sso-service/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthController membungkus service registry.
type AuthController struct {
	service services.IServiceRegistry
}

// IAuthController adalah kontrak handler HTTP domain autentikasi.
type IAuthController interface {
	LoginController(*fiber.Ctx) error
	LogoutController(*fiber.Ctx) error
	ResetPassword(*fiber.Ctx) error
	ChangePassword(*fiber.Ctx) error
	CekOtp(*fiber.Ctx) error
	RefreshToken(*fiber.Ctx) error
}

// NewAuthController membuat instance AuthController baru.
func NewAuthController(service services.IServiceRegistry) IAuthController {
	return &AuthController{service: service}
}

func (ctrl *AuthController) LoginController(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.LoginForm)

	data := ctrl.service.GetAuth().Login(c.UserContext(), form.Username, form.Password)
	if !data.Success {
		return c.JSON(data)
	}

	user := data.Data.(map[string]interface{})

	resp := map[string]interface{}{}
	resp["user"] = user

	refreshToken, _ := middleware.RefreshToken(user["username"].(string), user["id_user"].(string))
	setHttpCookie(c, refreshToken)

	accessToken, _ := middleware.AccessToken(user["username"].(string), user["id_user"].(string))
	resp["access_token"] = accessToken

	respJson := response.Response{
		Success: data.Success,
		Message: data.Message,
		Data:    resp,
	}

	return c.JSON(respJson)
}

func (ctrl *AuthController) LogoutController(c *fiber.Ctx) error {
	cookieRefreshToken := fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   true,
	}
	c.Cookie(&cookieRefreshToken)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully, session cookie deleted!",
	})
}

func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.ResetForm)

	uuidStr := uuid.New().String()

	data := ctrl.service.GetAuth().ResetPassword(c.UserContext(), form.Username, uuidStr)

	if data.Success {
		OTP, _ := middleware.OTPResetpassword(uuidStr, 0)
		setOTPCookie(c, OTP, "OTP")
	}

	return c.JSON(data)
}

func (ctrl *AuthController) ChangePassword(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.ChangePasswordForm)

	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing or invalid Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	data := ctrl.service.GetAuth().ChangePassword(c.UserContext(), form.Username, form.Password, tokenString)

	return c.JSON(data)
}

func (ctrl *AuthController) CekOtp(c *fiber.Ctx) error {
	form := c.Locals("validatedForm").(*dto.OtpForm)

	token := middleware.CheckjwtOTP(c)
	if !token {
		return c.JSON(response.Response{Success: false, Message: "OTP Salah, Silahkan Coba Lagi"})
	}

	jwtotp, err := middleware.ParseTokenOTP(c)
	if err != nil {
		return c.JSON(response.Response{Success: false, Message: "OTP Salah, Silahkan Coba Lagi"})
	}

	IDResetPassword := jwtotp.IDResetPassword
	otp := form.Otp

	data := ctrl.service.GetAuth().CekOtp(c.UserContext(), IDResetPassword, otp)

	if data.Success {
		token := ctrl.service.GetAuth().GenerateTokenChangePassword(c.UserContext(), IDResetPassword)
		setOTPCookie(c, token, "access_token")
	} else {
		cekpercobaan := ctrl.service.GetAuth().CekPercobaan(c.UserContext(), IDResetPassword)
		OTP, _ := middleware.OTPResetpassword(IDResetPassword, cekpercobaan)
		setOTPCookie(c, OTP, "OTP")
	}

	return c.JSON(data)
}

func (ctrl *AuthController) RefreshToken(c *fiber.Ctx) error {
	resp := response.Response{Success: false, Message: "unauthenticated"}

	token := middleware.CheckjwtRefresh(c)
	if token {
		parsed := middleware.JWTParse(c)
		resp = response.Response{Success: parsed.Success, Message: parsed.Message, Data: parsed.Data}
	}

	return c.JSON(resp)
}

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
