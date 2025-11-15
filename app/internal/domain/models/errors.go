package models

import "errors"

var (
	ErrInternal   = errors.New("internal")
	ErrTimeout    = errors.New("timeout")
	ErrCanceled   = errors.New("canceled")
	ErrNotFound   = errors.New("not found")
	ErrTeamExists = errors.New("team already exists")
	ErrPRExists   = errors.New("pr already exists")
)
