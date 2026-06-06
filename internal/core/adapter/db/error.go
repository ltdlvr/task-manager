package db

import "errors"

var (
	ErrEntityNotFound = errors.New("entity not found")
	ErrForbidden      = errors.New("forbidden")
)

// TODO - зарефакторить ошибки, отделить db-ошибки от сервисных и т.д.
