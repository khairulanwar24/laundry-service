package middleware

import (
	"fmt"
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

var jwtSecret = []byte("kaucantikhariini")
var jwtSecret2 = []byte("lama2akukangenkamu")

// AccessToken membuat access token (berlaku 5 menit).
func AccessToken(username, id_user string) (string, int) {
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

// RefreshToken membuat refresh token (berlaku 8 jam).
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

// JWTMiddleware memvalidasi access token pada header Authorization.
func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing or invalid Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	return c.Next()
}

// ParseToken memvalidasi & mengekstrak klaim dari access token.
func ParseToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ParseRefreshToken memvalidasi & mengekstrak klaim dari refresh token.
func ParseRefreshToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret2, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
