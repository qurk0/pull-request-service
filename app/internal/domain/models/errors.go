package models

import "errors"

var (
	ErrInternal      = errors.New("internal")
	ErrTimeout       = errors.New("timeout")
	ErrCanceled      = errors.New("canceled")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)
