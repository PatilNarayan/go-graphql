package validations

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Define a struct to represent each validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func init() {
	validate = validator.New()
}

func ValidateStruct(u interface{}) error {
	// Validate the struct
	err := validate.Struct(u)
	if err != nil {
		// If validation fails, the error might contain multiple fields
		// and we want to iterate over each error to provide better feedback

		// Define a slice to hold the validation errors
		var validationErrors []ValidationError

		// Loop over the validation errors and format them
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: fmt.Sprintf("Field '%s' failed on '%s' tag", err.Field(), err.Tag()),
			})
		}

		// Marshal the validation errors into JSON
		validationErrorsJSON, _ := json.Marshal(validationErrors)

		// Return a formatted error message with all validation issues
		return fmt.Errorf("Validation failed: %v", string(validationErrorsJSON))
	}

	// Return nil if validation is successful
	return nil
}

func UpdateDeletedMap() map[string]interface{} {
	return map[string]interface{}{
		"row_status": 0,
	}
}

// ValidateName validates that the input string matches the regex "^[A-Za-z0-9\\-_]+$".
func ValidateName(name string) error {
	// Define the regex pattern
	pattern := `^[A-Za-z0-9\-_]+$`
	// Compile the regex
	re := regexp.MustCompile(pattern)
	// Check if the name matches the regex
	if !re.MatchString(name) {
		return errors.New("invalid name: must contain only alphanumeric characters, hyphens, or underscores")
	}
	return nil
}
