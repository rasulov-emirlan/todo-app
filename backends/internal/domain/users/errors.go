package users

import "errors"

var (
	ErrNoSuchUser        = errors.New("users: no such user")
	ErrEmailIsTaken      = errors.New("users: email is taken")
	ErrInvalidEmail      = errors.New("users: invalid email")
	ErrInvalidUsername   = errors.New("users: username has to be longer than 6 and shorter than 20 characters")
	ErrInvalidPassword   = errors.New("users: password has to be longer than 6 and shorter than 60 characters")
	ErrWrongPassword     = errors.New("users: wrong password")
	ErrInvalidRefreshKey = errors.New("users: invalid refresh key")
)
