package errors

import "errors"

var (
	ErrFailedToPrepareStmt    = errors.New("failed to prepare stmt")
	ErrFailedToInsertUser     = errors.New("failed to insert user")
	ErrFailedToGetLastId      = errors.New("failed to get last id")
	ErrFailedToUpdateUser     = errors.New("failed to update user")
	ErrFailedToDeleteUser     = errors.New("failed to delete user")
	ErrFailedToGetUser        = errors.New("failed to get user")
	ErrFailedToGetUserByEmail = errors.New("failed to get user by email")
	ErrUserNotFound           = errors.New("user not found")
)
