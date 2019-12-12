package errs

import "errors"

// some errors for file module
var (
	ErrNoSuchFile    = errors.New("no such file")
	ErrWrongFID      = errors.New("wrong fid")
	ErrTokenRequired = errors.New("token is required")
	ErrInvalidToken  = errors.New("invalid token")
	ErrPostOnly      = errors.New("upload only by POST")
)
