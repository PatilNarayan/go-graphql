package utils

import (
	"errors"
	"fmt"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"
	"regexp"
)

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

func CreateActionMap(store map[string]interface{}, actions []string) map[string]interface{} {
	for _, action := range actions {
		store[action] = map[string]interface{}{
			"name": action,
		}
	}
	return store
}

func GetActionMap(data map[string]interface{}, key string) map[string]interface{} {
	actionMap := data["actions"].(map[string]interface{})
	for _, value := range actionMap {
		value := value.(map[string]interface{})
		for key1 := range value {
			if key1 != "name" {
				delete(value, key1)
			}
		}
	}
	return actionMap
}

// FormatError formats errors into a custom ErrorResponse
func FormatError(err error) models.OperationResult {
	var errorCode, errorMessage, errorDetails string

	// Check if the error is of type *CustomError
	if customErr, ok := err.(*dto.CustomError); ok {
		// Extract errorCode, message, and details from the custom error
		errorCode = customErr.ErrorCode
		errorMessage = customErr.ErrorMessage
		errorDetails = customErr.ErrorDetails
	} else {
		// Fallback to default error code and message if it's not a CustomError
		errorCode = "GENERIC_ERROR"
		errorMessage = "An unknown error occurred"
		errorDetails = err.Error()
	}

	// Return the custom ErrorResponse
	errorResponse := &models.ResponseError{
		IsSuccess:    false,
		Message:      errorMessage,
		ErrorCode:    errorCode,
		ErrorDetails: &errorDetails, // Could be nil if you don't want to provide extra details
	}
	var opResult models.OperationResult = errorResponse
	return opResult
}

func FormatErrorStruct(errorCode, errorMessage, errorDetails string) error {
	// Return the custom ErrorResponse
	errorResponse := &dto.CustomError{
		ErrorMessage: errorMessage,
		ErrorCode:    errorCode,
		ErrorDetails: errorDetails,
	}
	return errorResponse
}

// formatSuccess formats a successful response in the `OperationResult` union
func FormatSuccess(data interface{}) (models.OperationResult, error) {
	// Type assertion: Convert 'data' from 'interface{}' to '[]models.Data'
	if typedData, ok := data.([]models.Data); ok {
		// If the assertion succeeds, return the result in OperationResult
		successResponse := &models.SuccessResponse{
			IsSuccess: true,
			Message:   "Operation successful",
			Data:      typedData, // Now, typedData is of type []models.Data
		}
		var opResult models.OperationResult = successResponse
		return opResult, nil
	}
	// If the type assertion fails, return an error
	return nil, fmt.Errorf("expected data to be of type []models.Data, but got %T", data)
}
