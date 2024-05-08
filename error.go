package flowconf

import "errors"

var (
	NotAPtrErr = errors.New("config needs to be a pointer")
	IsNilErr   = errors.New("config is nil")
)
