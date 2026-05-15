package common

// Code represents an application-specific error code.
type Code int

const (
	CodeSuccess Code = 0

	CodeInvalidParam  Code = 40001
	CodePasswordError Code = 40002
	CodeUserNotFound  Code = 40003
	CodePostNotFound  Code = 40004

	CodeUnauthorized Code = 40100
	CodeForbidden    Code = 40300

	CodeNotFound    Code = 40400
	CodeEmailExists Code = 40901

	CodeInternalError Code = 50000
)

// AppError is the standard application error type.
// Message is sent to the frontend as a user-friendly description.
// Err holds the original Go error for backend logging only.
type AppError struct {
	Code    Code
	Message string
	Err     error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates an AppError with the given code and message.
func NewAppError(code Code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// WithErr wraps an AppError with the original error for backend logging.
func (e *AppError) WithErr(err error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

// Predefined application errors.
var (
	ErrEmailExists = &AppError{
		Code:    CodeEmailExists,
		Message: "This email is already registered. Please use a different email or try logging in.",
	}

	ErrInvalidParam = &AppError{
		Code:    CodeInvalidParam,
		Message: "One or more parameters provided are invalid or missing.",
	}

	ErrPasswordError = &AppError{
		Code:    CodePasswordError,
		Message: "The password you entered is incorrect. Please try again.",
	}

	ErrUserNotFound = &AppError{
		Code:    CodeUserNotFound,
		Message: "No user found with this email address.",
	}

	ErrPostNotFound = &AppError{
		Code:    CodePostNotFound,
		Message: "The requested post does not exist or has been removed.",
	}

	ErrInternalParam = &AppError{
		Code:    CodeInternalError,
		Message: "An unexpected error occurred. Please try again later.",
	}
)
