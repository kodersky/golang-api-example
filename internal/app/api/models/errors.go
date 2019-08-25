package models

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested item is not found")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given param is not valid")
	// ErrConflict will throw  409 if there is a conflict
	ErrConflict = errors.New("order already taken")
	// ErrTimeout will throw 408 if timeout is reached
	ErrTimeout = errors.New("please try again later")
)
