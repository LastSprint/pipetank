package errors

import "errors"

var (
	ErrBadInput = errors.New("[sig] bad input")
	ErrTimeout  = errors.New("[sig] timeout")
	ErrCanceled = errors.New("[sig] canceled")
	ErrInternal = errors.New("[sig ]internal error")
	ErrNotFound = errors.New("[sig] not found")
)
