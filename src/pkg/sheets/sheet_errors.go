package sheets

import (
	"errors"
	"fmt"
)

// Error is a sheet-specific error.
type Error interface {
	Error() string
	TypeName() string
}

// UnwrapError unwraps a sheet error.
func UnwrapError(err error) (Error, bool) {
	for err != nil {
		if sheetErr, ok := err.(Error); ok {
			return sheetErr, true
		}

		err = errors.Unwrap(err)
	}

	return nil, false
}

// ErrDivideByZero occurs when a number if divided by 0.
var ErrDivideByZero Error = &divideByZeroError{}

type divideByZeroError struct{}

func (e divideByZeroError) Error() string    { return "divide by zero" }
func (e divideByZeroError) TypeName() string { return "#DIV/0" }

// A NameError occurs when a formula or function is unable to
// find the data it needs to complete a calculation. This can
// happen for a number of reasons, including:
//
//   - A typo in the formula name
//   - The formula refers to a name that has not been defined
//   - The formula has a typo in the defined name
//   - The syntax is missing double quotation marks for text values
//   - A colon was omitted in a range reference
type NameError struct {
	Message string
}

func (e *NameError) TypeName() string {
	return "#NAME"
}

func (e *NameError) Error() string {
	return e.Message
}

// NameErrorf returns a NameError with a formnatted message.
func NameErrorf(msg string, args ...any) Error {
	return &NameError{
		Message: fmt.Sprintf(msg, args...),
	}
}

// A ValueError occurs when a formula contains cells with different data types,
// or when a formula references cells that contain text instead of numbers. The
// #VALUE! error can also occur when:
//
//   - A cell reference refers to an error value
//   - The function's syntax is incorrect
//   - IF statements are built in a complex way
//   - A cell contains a space character
//   - Cells contain hidden or non-printing characters
//   - Formula references a few ranges that are not of the same size or shape
type ValueError struct {
	Message string
}

func (e ValueError) TypeName() string {
	return "#VALUE"
}

func (e ValueError) Error() string {
	return e.Message
}

// ValueErrorf creates a new ValueError with a formatted message.
func ValueErrorf(msg string, args ...any) *ValueError {
	return &ValueError{
		Message: fmt.Sprintf(msg, args...),
	}
}

// NotAvailableError means "no value is available."
type NotAvailableError struct {
	Message string
}

func (e NotAvailableError) Error() string {
	return e.Message
}

func (e NotAvailableError) TypeName() string {
	return "#N/A"
}

// NotAvailableErrorf creates a new NotAvailableError with a formatted message.
func NotAvailableErrorf(msg string, args ...any) *NotAvailableError {
	return &NotAvailableError{
		Message: fmt.Sprintf(msg, args...),
	}
}

var (
	_ Error = &ValueError{}
	_ Error = &NotAvailableError{}
	_ Error = &NameError{}
	_ error = (Error)(nil)
)
