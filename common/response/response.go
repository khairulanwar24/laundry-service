// Package response berisi struktur & helper standar untuk membungkus balasan HTTP.
// Dipindahkan dari package `types` agar mengikuti pola berlapis (clean architecture).
package response

import "github.com/gofiber/fiber/v2"

// Response adalah bentuk baku balasan JSON API.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success mengirim balasan 200 OK dengan data.
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Message: "Success",
		Data:    data,
	})
}

// Error mengirim balasan gagal dengan kode status & pesan tertentu.
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Response{
		Success: false,
		Message: message,
	})
}
