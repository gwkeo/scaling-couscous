package errors

import "errors"

var (
	ErrTokenExpired    = errors.New("token is expired")
	ErrInvalidToken    = errors.New("invalid token")
	ErrAccessForbidden = errors.New("access forbidden")
)
