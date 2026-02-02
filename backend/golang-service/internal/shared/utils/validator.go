package utils

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator returns a singleton validator instance
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		// Register custom tag name func to use json tags
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register custom validations here if needed
		// Example: validate.RegisterValidation("custom", customValidationFunc)
	})
	return validate
}

// ValidateStruct validates a struct and returns formatted error message
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}

// FormatValidationErrors formats validation errors into a readable message
func FormatValidationErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			messages = append(messages, formatFieldError(e))
		}
		return strings.Join(messages, "; ")
	}
	return err.Error()
}

func formatFieldError(e validator.FieldError) string {
	field := e.Field()
	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + e.Param() + " characters"
	case "max":
		return field + " must be at most " + e.Param() + " characters"
	case "oneof":
		return field + " must be one of: " + e.Param()
	case "uuid":
		return field + " must be a valid UUID"
	case "gte":
		return field + " must be greater than or equal to " + e.Param()
	case "lte":
		return field + " must be less than or equal to " + e.Param()
	default:
		return field + " is invalid"
	}
}

// SetupGinValidator sets up the validator for Gin binding
func SetupGinValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}
