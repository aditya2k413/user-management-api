package errors

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidDate  = errors.New("invalid date format")
)
