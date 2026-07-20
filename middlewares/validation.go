package middleware

import (
	"errors"
	"reflect"
	"laundry-service/types"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// Fungsi middleware untuk validasi form
func ValidateForm(form interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parsing body berdasarkan form yang diberikan
		if err := c.BodyParser(form); err != nil {
			respJson := types.Response{
				Success: false,
				Message: err.Error(), // "Invalid input",
				Data:    nil,
			}
			return c.Status(fiber.StatusBadRequest).JSON(respJson)

		}

		// Validasi form
		err := validate.Struct(form)
		if err != nil {
			errors := make(map[string]string)

			// Iterasi melalui error validasi dan menambahkan ke map errors
			for _, err := range err.(validator.ValidationErrors) {
				errors[err.Field()] = getErrorMessage(err)

			}
			// Jika ada error validasi, kirimkan response dengan status 400

			respJson := types.Response{
				Success: false,
				Message: "Invalid input",
				Data:    errors,
			}
			return c.Status(fiber.StatusBadRequest).JSON(respJson)
		}

		// Simpan data yang telah divalidasi ke dalam konteks Fiber untuk diakses di handler
		c.Locals("validatedForm", form)

		// Lanjutkan ke handler berikutnya jika validasi berhasil
		return c.Next()
	}
}

func getErrorMessage(err validator.FieldError) string {
	fieldName := err.Field()
	errors := make(map[string]string)
	switch err.Tag() {
	case "required":
		errors[fieldName] = fieldName + " Wajib Diisi"
	case "min":
		errors[fieldName] = fieldName + " Minimal " + err.Param() + " Karakter"
	case "max":
		errors[fieldName] = fieldName + " Maksimal " + err.Param() + " Karakter"
	case "email":
		errors[fieldName] = "Kesalahan format email"
	case "gte":
		errors[fieldName] = fieldName + " Minimal " + err.Param() + " Karakter"
	case "lte":
		errors[fieldName] = fieldName + " Maksimal " + err.Param() + " Karakter"
	case "numeric":
		errors[fieldName] = fieldName + " Hanya Boleh Angka"
	case "oneof":
		errors[fieldName] = fieldName + " Harus Pilih salah satu dari " + err.Param()
	case "url":
		errors[fieldName] = "Kesalahan format url"
	case "uuid":
		errors[fieldName] = "Kesalahan format uuid"
	case "uuid4":
		errors[fieldName] = "Kesalahan format uuid4"
	case "uuid3":
		errors[fieldName] = "Kesalahan format uuid3"
	case "uuid5":
		errors[fieldName] = "Kesalahan format uuid5"
	case "uuid7":
		errors[fieldName] = "Kesalahan format uuid7"
	case "uuid8":
		errors[fieldName] = "Kesalahan format uuid8"
	case "uuid9":
		errors[fieldName] = "Kesalahan format uuid9"
	case "uuid10":
		errors[fieldName] = "Kesalahan format uuid10"
	case "contains":
		errors[fieldName] = fieldName + " Harus Memuat " + err.Param()
	case "containsany":
		errors[fieldName] = fieldName + " Harus Memuat 1 karakter dari " + err.Param()
	case "containsrune":
		errors[fieldName] = fieldName + " Harus Memuat 1 karakter dari " + err.Param()
	case "excludes":
		errors[fieldName] = fieldName + " Tidak Boleh Memuat " + err.Param()
	case "excludesall":
		errors[fieldName] = fieldName + " Tidak Boleh Memuat " + err.Param()
	case "excludesrune":
		errors[fieldName] = fieldName + " Tidak Boleh Memuat 1 karakter dari " + err.Param()
	case "isdefault":
		errors[fieldName] = fieldName + " Harus Memuat 1 karakter dari " + err.Param()
	case "isemail":
		errors[fieldName] = "Kesalahan format email"
	case "length":
		errors[fieldName] = fieldName + " Minimal " + err.Param() + " Karakter"
	case "greaterthan":
		errors[fieldName] = fieldName + " Minimal " + err.Param() + " Karakter"
	case "lessthan":
		errors[fieldName] = fieldName + " Maksimal " + err.Param() + " Karakter"
	case "datetime":
		errors[fieldName] = "Kesalahan format datetime"
	case "date":
		errors[fieldName] = "Kesalahan format date"
	case "time":
		errors[fieldName] = "Kesalahan format time"
	case "json":
		errors[fieldName] = "Kesalahan format json"
	case "base64":
		errors[fieldName] = "Kesalahan format base64"
	case "base64url":
		errors[fieldName] = "Kesalahan format base64url"
	case "file":
		errors[fieldName] = "Kesalahan format file"
	case "fileext":
		errors[fieldName] = "Kesalahan format fileext"
	case "filesize":
		errors[fieldName] = "Kesalahan format filesize"
	case "filemime":
		errors[fieldName] = "Kesalahan format filemime"
	default:
		return "Kesalahan pada " + err.Field()
	}

	return errors[fieldName]
}

func ValidatedParams(form interface{}) fiber.Handler {

	return func(c *fiber.Ctx) error {
		// Create a new instance of the form (pointer to struct)
		formPtr := reflect.New(reflect.TypeOf(form).Elem()).Interface()

		// Extract route params and bind them to the form fields
		if err := c.ParamsParser(formPtr); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Failed to parse parameters: " + err.Error(),
				"Data":    nil,
			})
		}

		// Validate the form
		if err := validate.Struct(formPtr); err != nil {
			var validationErrors validator.ValidationErrors

			// Use errors.As to check if the error is of type validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				// Iterate over validation errors and add them to the map
				errorsMap := make(map[string]string)
				for _, err := range validationErrors {
					errorsMap[err.Field()] = getErrorMessage(err)
				}

				// If validation fails, send a response with status 400
				respJson := types.Response{
					Success: false,
					Message: "Invalid input",
					Data:    errorsMap,
				}
				return c.Status(fiber.StatusBadRequest).JSON(respJson)
			} else {
				// Handle invalid validation error
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Invalid validation error: " + err.Error(),
					"Data":    nil,
				})
			}
		}

		// Store the validated form in Fiber's context
		c.Locals("validatedParams", formPtr)
		// fmt.Println(c.Locals("validatedParams"))
		// Continue to the next handler
		return c.Next()
	}
}

func ValidatedParams2(form interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new instance of the form type dynamically
		formType := reflect.TypeOf(form).Elem()
		formInstance := reflect.New(formType).Interface()

		// Parse query parameters into the form struct
		if err := c.QueryParser(formInstance); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Failed to parse parameters: " + err.Error(),
				"Data":    nil,
			})
		}

		// Perform validation
		if err := validate.Struct(formInstance); err != nil {
			validationErrors := err.(validator.ValidationErrors)

			if errors.As(err, &validationErrors) {
				// Create a map to store error messages
				errors := make(map[string]string)

				// Iterate over validation errors and get custom messages
				for _, fieldError := range validationErrors {
					errors[fieldError.Field()] = getErrorMessage(fieldError)
				}

				respJson := types.Response{
					Success: false,
					Message: "Invalid input",
					Data:    errors,
				}
				return c.Status(fiber.StatusBadRequest).JSON(respJson)

			} else {
				// Handle invalid validation error
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Invalid validation error: " + err.Error(),
					"Data":    nil,
				})
			}

		}

		// Store the validated form in Locals
		c.Locals("validatedForm", formInstance)

		// Proceed to the next handler
		return c.Next()
	}
}
