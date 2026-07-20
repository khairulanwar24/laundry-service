package middleware

import (
	"fmt"
	"laundry-service/types"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Username string `json:"username"`
	Id_user  string `json:"id_user"`
	jwt.RegisteredClaims
}

type JwtOTPCustomClaims struct {
	IDResetPassword string `json:"id_reset_password"`
	Percobaan       int32  `json:"percobaan"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("kaucantikhariini")
var jwtSecret2 = []byte("lama2akukangenkamu")
var jwtSecret3 = []byte("kulaluibersamamu")

// generate token  and refresh token jwt
func AccessToken(username, id_user string) (string, int) {
	// resp := make(map[string]interface{})

	claims := &JwtCustomClaims{
		username,
		id_user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		resp := err.Error()
		return resp, 400
	}

	return t, 200
}

func RefreshToken(username, id_user string) (string, int) {

	claims := &JwtCustomClaims{
		username,
		id_user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	r, err := token.SignedString(jwtSecret2)

	if err != nil {
		resp := err.Error()
		return resp, 400
	}

	return r, 200
}

func Checkjwt(c *fiber.Ctx) bool {

	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	t, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return false
	}
	if !t.Valid {
		return false
	}
	return true
}

func CheckjwtRefresh(c *fiber.Ctx) bool {

	authHeader := c.Cookies("refresh_token")
	if authHeader == "" {
		return false
	}

	tokenString := authHeader

	t, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret2, nil
	})
	if err != nil {
		return false
	}
	if !t.Valid {
		return false
	}
	return true
}

// JWTMiddleware checks for a valid JWT token in the request
func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing or invalid Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// Extract custom claims
	// if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
	// 	fmt.Println("Username:", claims.Userid)
	// } else {
	// 	fmt.Println("Invalid token")
	// }

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	return c.Next()
}

// token gRPC
func ValidateTokenMiddleware(token string) (bool, string) {
	// fmt.Println("token", token)
	if token == "" {
		return false, "Missing or invalid Authorization header"
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")

	parsedToken, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// fmt.Println("sastro")
	// fmt.Println(parsedToken)

	if err != nil || !parsedToken.Valid {
		return false, "Invalid token"
	}

	return true, ""
}

func JWTParse(c *fiber.Ctx) types.Response {
	// authHeader := c.Get("Authorization")
	var resp types.Response
	authHeader := c.Cookies("refresh_token")

	if authHeader == "" {
		resp.Success = false
		resp.Message = "Missing or invalid Authorization Cookie"
		resp.Data = nil

		return resp
	}

	tokenString := authHeader

	// if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
	// 	resp.Success = false
	// 	resp.Message = "Missing or invalid Authorization header"
	// 	resp.Data = nil

	// 	return resp
	// }

	// tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// fmt.Println(tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret2, nil
	})

	if err != nil || !token.Valid {
		resp.Success = false
		resp.Message = "Invalid token"
		resp.Data = nil

		return resp
		// return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}
	// userid := ""

	// Extract custom claims
	claims, ok := token.Claims.(*JwtCustomClaims)
	data := make(map[string]interface{})

	// fmt.Println(claims)
	if ok && token.Valid {
		token, _ := AccessToken(claims.Username, claims.Id_user)
		data["access_token"] = token
		resp.Success = true
		resp.Message = "success"
		resp.Data = data
	} else {
		resp.Success = false
		resp.Message = "Invalid token"
		resp.Data = nil
		return resp
	}

	return resp
}

func ParseToken(tokenString string) (*JwtCustomClaims, error) {
	// Parse the token

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func OTPResetpassword(id_reset_password string, percobaan int32) (string, int) {

	claims := &JwtOTPCustomClaims{
		id_reset_password,
		percobaan,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	r, err := token.SignedString(jwtSecret3)

	if err != nil {
		resp := err.Error()
		return resp, 400
	}

	return r, 200
}

func CheckjwtOTP(c *fiber.Ctx) bool {

	authHeader := c.Cookies("OTP")
	if authHeader == "" {
		return false
	}

	tokenString := authHeader

	t, err := jwt.ParseWithClaims(tokenString, &JwtOTPCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret3, nil
	})
	if err != nil {
		return false
	}
	if !t.Valid {
		return false
	}

	return true

}
func ParseTokenOTP(c *fiber.Ctx) (*JwtOTPCustomClaims, error) {

	authHeader := c.Cookies("OTP")
	if authHeader == "" {
		return nil, fmt.Errorf("missing or invalid authorization cookie")
	}

	tokenString := authHeader
	// Parse the token

	token, err := jwt.ParseWithClaims(tokenString, &JwtOTPCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret3, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*JwtOTPCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func AccessTokenOTP(username, id_user string) (string, int) {
	// resp := make(map[string]interface{})

	claims := &JwtCustomClaims{
		username,
		id_user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(jwtSecret3)
	if err != nil {
		resp := err.Error()
		return resp, 400
	}

	return t, 200
}

func ParseTokenChangePassword(tokenString string) (*JwtCustomClaims, error) {
	// Parse the token

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret3, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
