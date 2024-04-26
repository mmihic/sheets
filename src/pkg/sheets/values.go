package sheets

import (
	"fmt"
	"strconv"
	"time"
)

// Value is a cell value.
type Value interface {
	AsString() (string, error)
	AsTime() (time.Time, error)
	AsFloat64() (float64, error)
}

// Various types of values
type (
	Float64Value float64
	TimeValue    time.Time
	StringValue  string
)

// AsString converts the value to a string.
func (v Float64Value) AsString() (string, error) {
	return strconv.FormatFloat(float64(v), 'g', -1, 64), nil
}

// AsTime converts the value to a time.Time.
func (v Float64Value) AsTime() (time.Time, error) {
	// We use Excel formatting - the date portion is an integer value of the number of
	// days since January 1, 1900, which stored as the number 1. The time portion is a
	// decimal between .0 and .99999 representing the fraction of the day, with .0 as
	// 00:00:00 and .99999 as 23:59:59.
	return FromExcelTime(float64(v))
}

// AsFloat64 converts the value to a float64.
func (v Float64Value) AsFloat64() (float64, error) {
	return float64(v), nil
}

// AsString converts the value to a string.
func (v StringValue) AsString() (string, error) {
	return string(v), nil
}

// AsTime converts the value to a time.Time.
func (v StringValue) AsTime() (time.Time, error) {
	for _, layout := range supportedTimeLayouts {
		if t, err := time.Parse(layout, string(v)); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("'%s' cannot be parsed as a date or time", v)
}

// AsFloat64 converts the value to a float64.
func (v StringValue) AsFloat64() (float64, error) {
	return strconv.ParseFloat(string(v), 64)
}

// AsFloat64 converts the value to a float64.
func (v TimeValue) AsFloat64() (float64, error) {
	return ToExcelTime(time.Time(v))
}

// AsString converts the value to a string.
func (v TimeValue) AsString() (string, error) {
	return time.Time(v).Format(time.RFC3339), nil
}

// AsTime converts the value to a time.Time.
func (v TimeValue) AsTime() (time.Time, error) {
	return time.Time(v), nil
}

// A ValueIter is an iterator over a set of values.
type ValueIter interface {
	Next() bool
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

func (iter *sliceValueIter) Next() bool {
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

func (iter *singleValueIter) Next() bool {
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

var (
	supportedTimeLayouts = []string{
		time.RFC3339,
		time.DateOnly,
		time.DateTime,
		time.UnixDate,
		time.Kitchen,
		"2006/01/02",
	}

	_ Value = Float64Value(100)
	_ Value = TimeValue(time.Time{})
	_ Value = StringValue("")

	_ ValueIter = &sliceValueIter{}
	_ ValueIter = &singleValueIter{}
)
