package validation

import (
	"fmt"
	"strings"

	"github.com/KILLERGTG01/smart-task-planner-be/internal/db"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateUser(user *db.User) error {
	if err := validate.Struct(user); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func ValidatePlan(plan *db.Plan) error {
	if err := validate.Struct(plan); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// formatValidationError formats validation errors into a more readable format
func formatValidationError(err error) error {
	var errorMessages []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errorMessages = append(errorMessages, formatFieldError(fieldError))
		}
	}

	return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, ", "))
}

// formatFieldError formats individual field validation errors
func formatFieldError(fieldError validator.FieldError) string {
	field := fieldError.Field()
	tag := fieldError.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, fieldError.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, fieldError.Param())
	default:
		return fmt.Sprintf("%s failed validation for tag '%s'", field, tag)
	}
}
