package middlewares

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Global validator instance
var validate = validator.New()

// Function to get JSON field name from struct
func getJSONFieldName[T any](fieldName string) string {
	var t T
	typ := reflect.TypeOf(t)
	field, found := typ.FieldByName(fieldName)
	if !found {
		return fieldName
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}
	return strings.Split(jsonTag, ",")[0]
}

// Parse form-data manually
func parseFormData[T any](c *fiber.Ctx) (T, error) {
	var body T
	typ := reflect.TypeOf(body)
	val := reflect.New(typ).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		formKey := strings.Split(jsonTag, ",")[0] // Get JSON/form field name
		formValue := c.FormValue(formKey)         // Get form-data value

		if formValue != "" {
			switch field.Type.Kind() {
			case reflect.String:
				val.Field(i).SetString(formValue)

			case reflect.Int, reflect.Int64:
				if intValue, err := strconv.ParseInt(formValue, 10, 64); err == nil {
					val.Field(i).SetInt(intValue)
				}

			case reflect.Float64:
				if floatValue, err := strconv.ParseFloat(formValue, 64); err == nil {
					val.Field(i).SetFloat(floatValue)
				}

			case reflect.Struct:
				// Handle MongoDB ObjectID conversion
				if field.Type == reflect.TypeOf(primitive.ObjectID{}) {
					objID, err := primitive.ObjectIDFromHex(formValue)
					if err == nil {
						val.Field(i).Set(reflect.ValueOf(objID))
					}
				}
			}
		}
	}

	return val.Interface().(T), nil
}

// Generic validation middleware
func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		var err error

		contentType := c.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			err = c.BodyParser(&body)
		} else if strings.Contains(contentType, "multipart/form-data") {
			body, err = parseFormData[T](c)
		}

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := validate.Struct(body); err != nil {
			errors := make(map[string]string)
			for _, e := range err.(validator.ValidationErrors) {
				jsonField := getJSONFieldName[T](e.StructField())
				errors[jsonField] = fmt.Sprintf("%s is required", jsonField)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
		}

		c.Locals("body", body)
		return c.Next()
	}
}
