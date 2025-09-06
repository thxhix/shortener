package errors

import "errors"

// ErrDuplicate is returned when hash exist for the given link.
var ErrDuplicate = errors.New("такая ссылка уже сжата")
