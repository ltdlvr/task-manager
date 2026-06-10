package db

import "errors"

var (
	ErrEntityNotFound = errors.New("entity not found")
	ErrForbidden      = errors.New("forbidden")
	ErrOutOfRange     = errors.New("out of range")
)

// TODO - зарефакторить ошибки, отделить db-ошибки от сервисных и т.д.
