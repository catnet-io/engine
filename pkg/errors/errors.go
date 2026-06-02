package errors

import "errors"

var (
	ErrInvalidTarget      = errors.New("invalid target format")
	ErrUnsupportedProfile = errors.New("unsupported scan profile")
	ErrPermissionDenied   = errors.New("permission denied: root/admin privileges required")
	ErrTimeout            = errors.New("scan operation timed out")
	ErrExportFailed       = errors.New("export failed")
	ErrScanInProgress     = errors.New("scan is already in progress")
)
