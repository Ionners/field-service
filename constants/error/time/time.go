package error

import "errors"

var (
	ErrTimeNotFound = errors.New("field not found")
)

var TimeErrors = []error{
	ErrTimeNotFound,
}
