package handlers

import (
	"testing"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: nil,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: nil,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: domain.ErrUserEmailRequired,
		},
		{
			name:    "invalid email - no @",
			email:   "userexample.com",
			wantErr: domain.ErrInvalidEmail,
		},
		{
			name:    "invalid email - no domain",
			email:   "user@",
			wantErr: domain.ErrInvalidEmail,
		},
		{
			name:    "invalid email - no TLD",
			email:   "user@example",
			wantErr: domain.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if err != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUserName(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		wantErr  error
	}{
		{
			name:     "valid name",
			userName: "John Doe",
			wantErr:  nil,
		},
		{
			name:     "empty name",
			userName: "",
			wantErr:  domain.ErrUserNameRequired,
		},
		{
			name:     "whitespace only",
			userName: "   ",
			wantErr:  domain.ErrUserNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserName(tt.userName)
			if err != tt.wantErr {
				t.Errorf("validateUserName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessageContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr error
	}{
		{
			name:    "valid content",
			content: "Hello, World!",
			wantErr: nil,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: domain.ErrMessageContentEmpty,
		},
		{
			name:    "whitespace only",
			content: "   ",
			wantErr: domain.ErrMessageContentEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessageContent(tt.content)
			if err != tt.wantErr {
				t.Errorf("validateMessageContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessageKind(t *testing.T) {
	tests := []struct {
		name    string
		kind    domain.MessageKind
		wantErr error
	}{
		{
			name:    "valid kind - info",
			kind:    domain.MessageKindInfo,
			wantErr: nil,
		},
		{
			name:    "valid kind - warning",
			kind:    domain.MessageKindWarning,
			wantErr: nil,
		},
		{
			name:    "valid kind - error",
			kind:    domain.MessageKindError,
			wantErr: nil,
		},
		{
			name:    "valid kind - success",
			kind:    domain.MessageKindSuccess,
			wantErr: nil,
		},
		{
			name:    "invalid kind",
			kind:    domain.MessageKind("invalid"),
			wantErr: domain.ErrInvalidMessageKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessageKind(tt.kind)
			if err != tt.wantErr {
				t.Errorf("validateMessageKind() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		offset  int
		wantErr error
	}{
		{
			name:    "valid pagination",
			limit:   20,
			offset:  0,
			wantErr: nil,
		},
		{
			name:    "valid pagination with offset",
			limit:   50,
			offset:  100,
			wantErr: nil,
		},
		{
			name:    "limit too small",
			limit:   0,
			offset:  0,
			wantErr: domain.ErrInvalidLimit,
		},
		{
			name:    "limit too large",
			limit:   101,
			offset:  0,
			wantErr: domain.ErrInvalidLimit,
		},
		{
			name:    "negative offset",
			limit:   20,
			offset:  -1,
			wantErr: domain.ErrInvalidOffset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePagination(tt.limit, tt.offset)
			if err != tt.wantErr {
				t.Errorf("validatePagination() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
