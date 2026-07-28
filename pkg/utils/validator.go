package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validate is the singleton validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using go-playground/validator tags.
//
// Returns a formatted error message with all validation errors.
//
// Example usage:
//
//	type CreateUserRequest struct {
//	    Name  string `validate:"required,min=3,max=255"`
//	    Email string `validate:"required,email"`
//	}
//
//	req := CreateUserRequest{Name: "Jo", Email: "invalid"}
//	if err := ValidateStruct(req); err != nil {
//	    // Handle validation error
//	}
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Format validation errors
	validationErrors := err.(validator.ValidationErrors)
	errorMessages := make([]string, 0, len(validationErrors))

	for _, e := range validationErrors {
		errorMessages = append(errorMessages, formatValidationError(e))
	}

	return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
}

// formatValidationError converts a validation error to a human-readable message.
func formatValidationError(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, e.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	default:
		return fmt.Sprintf("%s failed validation for '%s'", field, tag)
	}
}

// GetValidator returns the singleton validator instance.
//
// Use this to register custom validation functions.
//
// Example:
//
//	v := GetValidator()
//	v.RegisterValidation("custom_tag", func(fl validator.FieldLevel) bool {
//	    // Custom validation logic
//	    return true
//	})
func GetValidator() *validator.Validate {
	return validate
}
