package sheets

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

// An Operator is a comparison or arithmetic operator in an expression
type Operator string

// Various Operators
const (
	Add      Operator = "+"
	Subtract Operator = "-"
	Multiply Operator = "*"
	Divide   Operator = "/"
	Exp      Operator = "^"
	Gt       Operator = ">"
	Lt       Operator = "<"
	Geq      Operator = ">="
	Leq      Operator = "<="
	Eq       Operator = "="
	Neq      Operator = "<>"
)

// Apply applies the operator to two values, returning the results of the operator.
func (op Operator) Apply(v1, v2 Value) Value {
	if ok := op.isArithmetic(); ok {
		return op.applyArithmetic(v1, v2)
	}

	if ok := op.isComparison(); ok {
		cmpResult, err := op.compare(v1, v2)
		if err != nil {
			return ErrorValue{err}
		}

		return op.applyComparison(cmpResult)
	}

	return ErrorValue{NameErrorf("unsupported operator '%s'", op)}
}

func (op Operator) applyComparison(cmpResult int) Value {
	switch op {
	case Eq:
		return BoolValue(cmpResult == 0)
	case Lt:
		return BoolValue(cmpResult < 0)
	case Leq:
		return BoolValue(cmpResult <= 0)
	case Gt:
		return BoolValue(cmpResult > 0)
	case Geq:
		return BoolValue(cmpResult >= 0)
	case Neq:
		return BoolValue(cmpResult != 0)
	default:
		return ErrorValue{&ValueError{
			Message: fmt.Sprintf("cannot interpret %d as a bool", cmpResult),
		}}
	}
}

func (op Operator) applyArithmetic(v1, v2 Value) Value {
	// If either value is an error, return the error
	if errVal, ok := v1.(ErrorValue); ok {
		return errVal
	}

	if errVal, ok := v2.(ErrorValue); ok {
		return errVal
	}

	n1, err := v1.ToFloat64()
	if err != nil {
		return ErrorValue{err}
	}

	n2, err := v2.ToFloat64()
	if err != nil {
		return ErrorValue{err}
	}

	switch op {
	case Add:
		return Float64Value(n1 + n2)
	case Subtract:
		return Float64Value(n1 - n2)
	case Divide:
		if n2 == 0 {
			return ErrorValue{ErrDivideByZero}
		}

		return Float64Value(n1 / n2)
	case Multiply:
		return Float64Value(n1 * n2)
	case Exp:
		return Float64Value(math.Pow(n1, n2))
	default:
		// This is an internal coding error
		panic(&ValueError{
			Message: fmt.Sprintf("'%s' is not an arithmetic operator", op),
		})
	}
}

func (op Operator) isArithmetic() bool {
	switch op {
	case Add, Subtract, Multiply, Divide, Exp:
		return true
	case Gt, Lt, Geq, Leq, Eq, Neq:
		return false
	default:
		panic(&NameError{
			fmt.Sprintf("invalid operator '%s'", op),
		})
	}
}

func (op Operator) compare(v1, v2 Value) (int, error) {
	switch tv1 := v1.(type) {
	case Float64Value:
		return op.compareFloats(float64(tv1), v2)
	case StringValue:
		return op.compareStrings(string(tv1), v2)
	case TimeValue:
		return op.compareFloats(ToExcelTime(time.Time(tv1)), v2)
	case BoolValue:
		return op.compareBools(bool(tv1), v2)
	case ErrorValue:
		return 0, tv1.Err
	default:
		// This is an internal coding error
		panic(&ValueError{
			fmt.Sprintf("unsupported type %s: '%s'", reflect.TypeOf(v1), v1),
		})
	}
}

type ordered interface {
	string | float64
}

func compare[T ordered](v1, v2 T) int {
	if v1 == v2 {
		return 0
	}

	if v1 < v2 {
		return -1
	}

	return 1
}

func (op Operator) compareFloats(n1 float64, v2 Value) (int, error) {
	switch tv2 := v2.(type) {
	case StringValue, BoolValue:
		return -1, nil
	case TimeValue:
		return compare(n1, ToExcelTime(time.Time(tv2))), nil
	case Float64Value:
		return compare(n1, float64(tv2)), nil
	case ErrorValue:
		return 0, tv2.Err
	default:
		// This is an internal coding error
		panic(&ValueError{
			fmt.Sprintf("unsupported type %s: '%s'", reflect.TypeOf(v2), v2),
		})
	}
}

func (op Operator) compareStrings(s1 string, v2 Value) (int, error) {
	switch tv2 := v2.(type) {
	case Float64Value, TimeValue:
		return 1, nil
	case BoolValue:
		return -1, nil
	case StringValue:
		return compare(s1, string(tv2)), nil
	case ErrorValue:
		return 0, tv2.Err
	default:
		// This is an internal coding error
		panic(&ValueError{
			fmt.Sprintf("unsupported type %s: '%s'", reflect.TypeOf(v2), v2),
		})
	}
}

func (op Operator) compareBools(b1 bool, v2 Value) (int, error) {
	switch tv2 := v2.(type) {
	case Float64Value, StringValue, TimeValue:
		return 1, nil
	case BoolValue:
		b2 := bool(tv2)
		if (b1 && b2) || (!b1 && !b2) {
			return 0, nil
		}

		if b1 && !b2 {
			return 1, nil
		}

		return -1, nil
	case ErrorValue:
		return 0, tv2.Err
	default:
		// This is an internal coding error
		panic(&ValueError{
			fmt.Sprintf("unsupported type %s: '%s'", reflect.TypeOf(v2), v2),
		})
	}
}

func (op Operator) isComparison() bool {
	switch op {
	case Add, Subtract, Multiply, Divide, Exp:
		return false
	case Gt, Lt, Geq, Leq, Eq, Neq:
		return true
	default:
		panic(&NameError{
			fmt.Sprintf("invalid operator '%s'", op),
		})
	}
}
