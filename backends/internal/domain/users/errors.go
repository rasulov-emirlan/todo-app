package users

import "errors"

var (
	ErrNoSuchUser        = errors.New("no such user")
	ErrEmailIsTaken      = errors.New("email is taken")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidUsername   = errors.New("username has to be longer than 6 and shorter than 20 characters")
	ErrInvalidPassword   = errors.New("password has to be longer than 6 and shorter than 60 characters")
	ErrWrongPassword     = errors.New("wrong password")
	ErrInvalidRefreshKey = errors.New("invalid refresh key")
)
