package sheets

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Value is a cell value.
type Value interface {
	ToFloat64() (float64, error)

	valueMarker()
}

// Various types of values.
type (
	Float64Value float64
	TimeValue    time.Time
	StringValue  string
	BoolValue    bool
)

// StringToValue converts a string into a Value, using the most optimal Value
// representation.
func StringToValue(s string) Value {
	if tm, err := ParseTime(s); err == nil {
		return TimeValue(tm)
	}

	if b, err := ParseBool(s); err == nil {
		return BoolValue(b)
	}

	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return Float64Value(n)
	}

	return StringValue(s)
}

func (v Float64Value) valueMarker() {}

// ToFloat64 converts the value to a float64.
func (v Float64Value) ToFloat64() (float64, error) {
	return float64(v), nil
}

func (v TimeValue) valueMarker() {}

// ToFloat64 converts the value to a float64.
func (v TimeValue) ToFloat64() (float64, error) {
	return ToExcelTime(time.Time(v)), nil
}

func (v StringValue) valueMarker() {}

// ToFloat64 converts the value to a float64.
func (v StringValue) ToFloat64() (float64, error) {
	n, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0, &ValueError{fmt.Sprintf("unable to convert '%s' to float", v)}
	}

	return n, nil
}

func (v BoolValue) valueMarker() {}

// ToFloat64 converts the value to a float64.
func (v BoolValue) ToFloat64() (float64, error) {
	if v {
		return 1, nil
	}

	return 0, nil
}

// ErrorValue is a value that is an error.
type ErrorValue struct {
	Err error
}

func (v ErrorValue) valueMarker() {}

// ToFloat64 converts the value to a float64.
func (v ErrorValue) ToFloat64() (float64, error) {
	return 0, v.Err
}

// A ValueIter is an iterator over a set of values.
type ValueIter interface {
	Next(ctx context.Context) bool
	Err() error
	Value() Value
	Index() int
	Len() int
}

// SliceValueIter wraps a slice of Values in an iterator.
func SliceValueIter(vals []Value) ValueIter {
	return &sliceValueIter{
		vals: vals,
	}
}

type sliceValueIter struct {
	vals    []Value
	nextIdx int
}

func (iter *sliceValueIter) Next(_ context.Context) bool {
	if iter.nextIdx >= len(iter.vals) {
		return false
	}

	iter.nextIdx++
	return true
}

func (iter *sliceValueIter) Err() error {
	return nil
}

func (iter *sliceValueIter) Value() Value {
	return iter.vals[iter.nextIdx-1]
}

func (iter *sliceValueIter) Index() int {
	return iter.nextIdx - 1
}

func (iter *sliceValueIter) Len() int {
	return len(iter.vals)
}

// SingleValueIter wraps a single Value in an iterator.
func SingleValueIter(val Value) ValueIter {
	return &singleValueIter{val: val}
}

type singleValueIter struct {
	consumed bool
	val      Value
}

func (iter *singleValueIter) Next(_ context.Context) bool {
	if iter.consumed {
		return false
	}

	iter.consumed = true
	return true
}

func (iter *singleValueIter) Err() error {
	return nil
}

func (iter *singleValueIter) Value() Value {
	return iter.val
}

func (iter *singleValueIter) Len() int {
	return 1
}

func (iter *singleValueIter) Index() int {
	return 0
}

// ParseBool parses one of the recognized boolean strings (a subset
// of the boolean strings recognized by golang).
func ParseBool(s string) (bool, error) {
	switch strings.ToUpper(s) {
	case "TRUE":
		return true, nil
	case "FALSE":
		return false, nil
	default:
		return false, &ValueError{fmt.Sprintf("'%s' is not a valid boolean value", s)}
	}
}

var (
	_ Value = Float64Value(100)
	_ Value = TimeValue(time.Time{})
	_ Value = StringValue("")

	_ ValueIter = &sliceValueIter{}
	_ ValueIter = &singleValueIter{}
)
