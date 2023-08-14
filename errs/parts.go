package errs

import "errors"

var (
	CategoryCreationError  = errors.New("category creation error")
	SelectionCreationError = errors.New("selection creation error")
)
