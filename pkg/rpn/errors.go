package rpn

import "errors"

var (
	ErrDividingByZero             = errors.New("dividing by zero")
	ErrDuplicateOpertaionsSigns   = errors.New("some of operations are dublicated")
	ErrOpertaionsSigns            = errors.New("expression couldn't starts/ends with operation sign")
	ErrExpressionStringEmpty      = errors.New("expression is empty")
	ErrExpressionStringParetheses = errors.New("expression has problem with parentheses")
)
