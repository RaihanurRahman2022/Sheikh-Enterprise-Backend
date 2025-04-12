package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Custom validation tags
const (
	passwordMinLen = 8
	passwordMaxLen = 72
)

// Initialize sets up the validator with custom validation rules
func Initialize() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validation tags
		_ = v.RegisterValidation("password", validatePassword)
		_ = v.RegisterValidation("phone", validatePhone)
		_ = v.RegisterValidation("username", validateUsername)

		// Register custom type for handling UUID strings
		v.RegisterCustomTypeFunc(validateUUID, string(""))

		// Register function to get json tag names
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// validatePassword checks if the password meets security requirements:
// - At least 8 characters
// - Maximum 72 characters (bcrypt limitation)
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one number
// - At least one special character
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < passwordMinLen || len(password) > passwordMaxLen {
		return false
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validatePhone checks if the phone number is valid
// Supports BD phone numbers (example: +8801XXXXXXXXX)
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	match, _ := regexp.MatchString(`^\+8801[3-9]\d{8}$`, phone)
	return match
}

// validateUsername checks if the username is valid:
// - 3-30 characters long
// - Can contain letters, numbers, dots, and underscores
// - Must start with a letter
// - Cannot have consecutive dots
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	match, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9._]*[a-zA-Z0-9]$`, username)
	if !match {
		return false
	}

	// Check for consecutive dots
	if strings.Contains(username, "..") {
		return false
	}

	return true
}

// validateUUID checks if the string is a valid UUID
func validateUUID(field reflect.Value) interface{} {
	if field.String() == "" {
		return nil
	}
	match, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, field.String())
	if !match {
		return nil
	}
	return field.String()
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a slice of ValidationError
type ValidationErrors []ValidationError

// FormatError formats validator.ValidationErrors into ValidationErrors
func FormatError(err error) ValidationErrors {
	var errors ValidationErrors
	validatorErrs := err.(validator.ValidationErrors)

	for _, e := range validatorErrs {
		var message string

		switch e.Tag() {
		case "required":
			message = "This field is required"
		case "email":
			message = "Invalid email address"
		case "password":
			message = "Password must be 8-72 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"
		case "phone":
			message = "Invalid phone number format. Must be a valid BD number (e.g., +8801XXXXXXXXX)"
		case "username":
			message = "Username must be 3-30 characters long, start with a letter, and can contain only letters, numbers, dots, and underscores"
		case "min":
			message = "Value must be greater than or equal to " + e.Param()
		case "max":
			message = "Value must be less than or equal to " + e.Param()
		case "uuid":
			message = "Invalid UUID format"
		default:
			message = "Invalid value"
		}

		errors = append(errors, ValidationError{
			Field:   e.Field(),
			Message: message,
		})
	}

	return errors
}
