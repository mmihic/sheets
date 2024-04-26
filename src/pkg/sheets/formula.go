package sheets

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// A Formula resolves / transforms values from a DataSet.
type Formula interface {
	marker() // Ensures implementation of the interface
	fmt.Stringer
}

// A Constant is a constant value.
type Constant struct {
	Value Value
}

func (c *Constant) marker() {}

// String returns the constant in string format.
func (c *Constant) String() string {
	switch ct := c.Value.(type) {
	case StringValue:
		return strconv.Quote(string(ct))
	case Float64Value:
		return strconv.FormatFloat(float64(ct), 'g', -1, 64)
	case TimeValue:
		return fmt.Sprintf(`"%s"`, time.Time(ct).Format(time.RFC3339))
	default:
		return fmt.Sprintf("%s", ct)
	}
}

// A CellReference is a reference to a cell in a sheet.
type CellReference struct {
	Sheet string
	Pos   Pos
}

func (r *CellReference) marker() {}

// String returns the string form of the reference.
func (r *CellReference) String() string {
	if r.Sheet != "" {
		return fmt.Sprintf("`%s`!%s", r.Sheet, r.Pos)
	}

	return r.Pos.String()
}

// A CellRangeReference is a reference to a range of cells in a sheet.
type CellRangeReference struct {
	Sheet string
	Range Range
}

func (r *CellRangeReference) marker() {}

// String returns the string form of the reference.
func (r *CellRangeReference) String() string {
	if r.Sheet != "" {
		return fmt.Sprintf("`%s`!%s", r.Sheet, r.Range)
	}

	return r.Range.String()
}

// A NamedRangeReference is a reference to a named range.
type NamedRangeReference struct {
	NamedRange string
}

func (r *NamedRangeReference) marker() {}

// String returns the named range reference.
func (r *NamedRangeReference) String() string {
	return r.NamedRange
}

// A FunctionCall is a call of a function with a set of arguments.
type FunctionCall struct {
	FunctionName string
	Args         []Formula
}

func (fc *FunctionCall) marker() {}

// String returns the function call in string form.
func (fc *FunctionCall) String() string {
	var sb strings.Builder
	sb.WriteString(fc.FunctionName)
	sb.WriteRune('(')
	for i, arg := range fc.Args {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.String())
	}
	sb.WriteRune(')')
	return sb.String()
}

var (
	_ Formula = &CellReference{}
	_ Formula = &CellRangeReference{}
	_ Formula = &FunctionCall{}
	_ Formula = &Constant{}
	_ Formula = &NamedRangeReference{}
)
