package errs

import "errors"

// some errors for user module
var (
	ErrNameRequired        = errors.New("name and password are required")
	ErrNameUsed            = errors.New("name already exists")
	ErrWrongNameOrPassword = errors.New("wrong name or password")
	ErrPath                = errors.New("wrong request path")
	ErrWrongUID            = errors.New("wrong uid")
)
