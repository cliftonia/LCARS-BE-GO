package constants

import "time"

// Pagination constants
const (
	DefaultPageLimit = 20
	MaxPageLimit     = 100
	MinPageLimit     = 1
)

// Server timeout constants
const (
	ServerReadTimeout  = 15 * time.Second
	ServerWriteTimeout = 15 * time.Second
	ServerIdleTimeout  = 60 * time.Second
	ShutdownTimeout    = 30 * time.Second
)

// Validation constants
const (
	MaxNameLength    = 255
	MaxEmailLength   = 255
	MaxContentLength = 5000
	MinNameLength    = 1
	MinEmailLength   = 3
)
