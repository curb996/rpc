package driver

import "errors"

var (
	ErrUnsupported = errors.New("driver: unsupported operation")
	ErrBadConfig   = errors.New("driver: bad config")
)
