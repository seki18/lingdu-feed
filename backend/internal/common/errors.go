package common

type Code int

const (
	CodeSuccess Code = 0

	CodeInvalidParam Code = 40001
	CodeUnauthorized Code = 40100
	CodeForbidden    Code = 40300

	CodeNotFound    Code = 40400
	CodeEmailExists Code = 40901

	CodeInternalError Code = 50000
)

type AppError struct {
	Code    Code
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrEmailExists = &AppError{
		Code:    CodeEmailExists,
		Message: "email already exists",
	}

	ErrInvalidParam = &AppError{
		Code:    CodeInvalidParam,
		Message: "invalid parameter",
	}

	ErrInternalParam = &AppError{
		Code:    CodeInternalError,
		Message: "internal server error",
	}
)
