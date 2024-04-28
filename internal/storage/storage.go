package storage

import "errors"

var (
	ErrUserNotFound      = errors.New("user not fount")
	ErrAppNotFound       = errors.New("app not fount")
	ErrUserAlreadyExists = errors.New("user already exist")
)
