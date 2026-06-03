package coreerr

import "errors"

// Sentinel errors representing the structured taxonomy for core operations.
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrTimeout      = errors.New("operation timed out")
	ErrCancelled    = errors.New("operation cancelled")
	ErrPermission   = errors.New("permission denied")
	ErrExport       = errors.New("export failed")
	ErrPartial      = errors.New("partial failure")
)
