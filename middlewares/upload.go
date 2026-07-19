package middleware

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	// "github.com/gofiber/fiber"
)

func FileUploadMiddleware(c *fiber.Ctx, uploadPath string, nameFormat string, allowedTypes []string, maxSize int64) error {
	// Get the file from the form field "file"
	file, err := c.FormFile(uploadPath)
	if err != nil {
		return fmt.Errorf("file upload error: %v", err)

	}
	// var maxSize int64 = 2 * 1024 * 1024
	if maxSize == 0 {
		maxSize = 2 * 1024 * 1024
	}
	if file.Size > maxSize {
		return fmt.Errorf(fmt.Sprintf("File too large. Maximum allowed size is %d bytes", maxSize))
	}

	conf := "asset/public"

	// Validate if the file type is allowed
	if !isValidFileType(file, allowedTypes) {
		allowedTypesStr := strings.Join(allowedTypes, ",")
		return fmt.Errorf("File " + file.Header["Content-Type"][0] + " type not allowed, allowed types: " + allowedTypesStr)
	}

	// Generate custom file name based on the format
	fileName := generateFileName(file.Filename, nameFormat)

	// Ensure the upload directory exists
	if _, err := os.Stat(conf + "/" + uploadPath); os.IsNotExist(err) {
		os.Mkdir(conf+"/"+uploadPath, os.ModePerm)
	}

	// Save the file to the desired location
	filePath := fmt.Sprintf("%s/%s", conf+"/"+uploadPath, fileName)
	err = c.SaveFile(file, filePath)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	// Store the file name in locals for use in the controller
	c.Locals("fileName", fileName)

	// Return nil to indicate successful processing
	return nil
}

// Function to generate a custom file name
func generateFileName(originalName string, format string) string {
	extension := originalName[strings.LastIndex(originalName, "."):]
	switch format {
	case "timestamp":
		return fmt.Sprintf("%d%s", time.Now().Unix(), extension)
	case "uuid":
		return fmt.Sprintf("%s%s", uuid.New().String(), extension)
	// case "custom":
	// 	return fmt.Sprintf("custom_name_%d%s", time.Now().Unix(), extension)
	default:
		// fmt.Println(format)
		return fmt.Sprintf(format+"_%d%s", time.Now().Unix(), extension)
	}
}

// Function to validate if the file type is allowed
func isValidFileType(file *multipart.FileHeader, allowedTypes []string) bool {
	openedFile, err := file.Open()
	if err != nil {
		return false
	}
	defer openedFile.Close()

	buf := make([]byte, 512)
	_, err = openedFile.Read(buf)
	if err != nil {
		return false
	}
	fileType := http.DetectContentType(buf)

	for _, t := range allowedTypes {
		if fileType == t {
			return true
		}
	}

	return false
}
