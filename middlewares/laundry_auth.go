// Package middleware (laundry_auth) berisi middleware autentikasi & otorisasi
// khusus domain laundry: JWT context (id_user) dan keanggotaan outlet.
package middleware

import (
	"strings"

	"laundry-service/database"

	"github.com/gofiber/fiber/v2"
)

// LaundryJWTMiddleware memvalidasi access token & menyimpan id_user ke c.Locals
// supaya handler di domain laundry bisa tahu siapa user yang sedang login.
func LaundryJWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing or invalid Authorization header"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := ParseToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	c.Locals("id_user", claims.Id_user)
	return c.Next()
}

// resolveOutletRole mengambil role user pada outlet tanpa memanggil c.Next(),
// supaya bisa dipakai ulang oleh RequireOutletMember maupun RequireOutletOwner.
func resolveOutletRole(c *fiber.Ctx) (string, error) {
	idOutlet := c.Params("id_outlet")
	idUser, _ := c.Locals("id_user").(string)

	if idOutlet == "" {
		return "", c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "id_outlet wajib diisi"})
	}
	if idUser == "" {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	var rows []map[string]interface{}
	result := database.DB.Raw(`
		SELECT role FROM laundry.user_outlets
		WHERE id_user = ? AND id_outlet = ? AND is_active = true AND status_data = true
	`, idUser, idOutlet).Scan(&rows)

	if result.Error != nil || len(rows) == 0 {
		return "", c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Anda tidak terdaftar di outlet ini"})
	}

	role, _ := rows[0]["role"].(string)
	return role, nil
}

// RequireOutletMember memastikan user yang login adalah anggota aktif outlet
// pada parameter :id_outlet, lalu menyimpan role-nya ke c.Locals("outlet_role").
func RequireOutletMember(c *fiber.Ctx) error {
	role, err := resolveOutletRole(c)
	if err != nil {
		return err
	}

	c.Locals("outlet_role", role)
	return c.Next()
}

// RequireOutletOwner memastikan user adalah anggota outlet DENGAN role owner.
func RequireOutletOwner(c *fiber.Ctx) error {
	role, err := resolveOutletRole(c)
	if err != nil {
		return err
	}

	if role != "owner" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Hanya owner outlet yang dapat melakukan aksi ini"})
	}

	c.Locals("outlet_role", role)
	return c.Next()
}
