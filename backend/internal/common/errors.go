package common

import "errors"

var (
	ErrEmailExists  = errors.New("email already exists")
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)
