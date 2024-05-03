package sheets

import (
	"errors"
	"fmt"
)

// ErrDivideByZero occurs when a number if divided by 0.
var ErrDivideByZero = errors.New("divide by zero")

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
	Name string
}

func (e NameError) Error() string {
	return fmt.Sprintf("invalid name '%s'", e.Name)
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

func (e ValueError) Error() string {
	return e.Message
}

// NotAvailableError means "no value is available."
type NotAvailableError struct {
	Message string
}

func (e NotAvailableError) Error() string {
	return e.Message
}
