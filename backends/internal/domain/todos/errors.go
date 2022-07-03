package todos

import "errors"

var (
	ErrInvalidTitle = errors.New("todos: title can't be less than 6 characters and more than 100 characters")
	ErrInvalidBody  = errors.New("todos: body can't be more than 2000 characters")

	ErrInvalidDeadline = errors.New("todos: deadline can't be in the past")
	ErrNotAllowed      = errors.New("todos: only admins are allowed to update todos that dont belong to them")
)
