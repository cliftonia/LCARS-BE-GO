// Package domain defines core business entities and repository interfaces.
package domain

import "errors"

// Common domain errors
var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrUserNameRequired  = errors.New("user name is required")
	ErrUserEmailRequired = errors.New("user email is required")

	// Message errors
	ErrMessageNotFound     = errors.New("message not found")
	ErrInvalidMessageID    = errors.New("invalid message ID")
	ErrMessageContentEmpty = errors.New("message content cannot be empty")
	ErrMessageUserIDEmpty  = errors.New("message user ID cannot be empty")
	ErrInvalidMessageKind  = errors.New("invalid message kind")

	// Pagination errors
	ErrInvalidLimit  = errors.New("invalid limit: must be between 1 and 100")
	ErrInvalidOffset = errors.New("invalid offset: must be non-negative")

	// Validation errors
	ErrNameTooLong    = errors.New("name exceeds maximum length")
	ErrEmailTooLong   = errors.New("email exceeds maximum length")
	ErrContentTooLong = errors.New("content exceeds maximum length")
)
