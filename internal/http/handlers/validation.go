package handlers

import (
	"regexp"
	"strings"

	"github.com/cliftonbaggerman/subspace-backend/internal/constants"
	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// validateEmail validates email format
func validateEmail(email string) error {
	if email == "" {
		return domain.ErrUserEmailRequired
	}

	if len(email) > constants.MaxEmailLength {
		return domain.ErrEmailTooLong
	}

	if !emailRegex.MatchString(email) {
		return domain.ErrInvalidEmail
	}

	return nil
}

// validateUserName validates user name
func validateUserName(name string) error {
	if strings.TrimSpace(name) == "" {
		return domain.ErrUserNameRequired
	}

	if len(name) > constants.MaxNameLength {
		return domain.ErrNameTooLong
	}

	return nil
}

// validateMessageContent validates message content
func validateMessageContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return domain.ErrMessageContentEmpty
	}

	if len(content) > constants.MaxContentLength {
		return domain.ErrContentTooLong
	}

	return nil
}

// validateMessageKind validates message kind
func validateMessageKind(kind domain.MessageKind) error {
	switch kind {
	case domain.MessageKindInfo, domain.MessageKindWarning, domain.MessageKindError, domain.MessageKindSuccess:
		return nil
	default:
		return domain.ErrInvalidMessageKind
	}
}

// validatePagination validates pagination parameters
func validatePagination(limit, offset int) error {
	if limit < constants.MinPageLimit || limit > constants.MaxPageLimit {
		return domain.ErrInvalidLimit
	}

	if offset < 0 {
		return domain.ErrInvalidOffset
	}

	return nil
}
