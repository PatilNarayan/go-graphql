package dto

import "fmt"

// CustomError is a custom error type that includes an errorCode and message
type CustomError struct {
	ErrorCode    string
	ErrorMessage string
	ErrorDetails string
}

// Implement the Error method to satisfy the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %s, message: %s, details: %s", e.ErrorCode, e.ErrorMessage, e.ErrorDetails)
}

// NewCustomError creates a new instance of CustomError
func NewCustomError(errorCode, errorMessage, errorDetails string) *CustomError {
	return &CustomError{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		ErrorDetails: errorDetails,
	}
}
