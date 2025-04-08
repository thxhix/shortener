package customErrors

import "errors"

var ErrDuplicate = errors.New("такая ссылка уже сжата")

type DuplicateError struct {
	Err error
}

func (e *DuplicateError) Error() string {
	return e.Err.Error()
}

func (e *DuplicateError) Unwrap() error {
	return e.Err
}

func NewDuplicateError() error {
	return &DuplicateError{Err: ErrDuplicate}
}
