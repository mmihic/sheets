package sheets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type invalidValue string

func (v invalidValue) valueMarker()                {}
func (v invalidValue) ToFloat64() (float64, error) { return 0, nil }
func (v invalidValue) String() string              { return "invalid value" }

var (
	_ Value = invalidValue("")
)

func TestMath(t *testing.T) {
	type operators map[Operator]Value

	for _, tt := range []struct {
		name      string
		first     Value
		second    Value
		operators operators
	}{
		{"string against same string",
			StringValue("A"), StringValue("A"), operators{
				Eq:       BoolValue(true),
				Neq:      BoolValue(false),
				Gt:       BoolValue(false),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(true),
				Add:      ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Multiply: ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Subtract: ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Divide:   ErrorValue{ValueErrorf("unable to convert 'A' to float")},
			}},
		{"string against greater string",
			StringValue("A"), StringValue("B"), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(false),
				Geq:      BoolValue(false),
				Lt:       BoolValue(true),
				Leq:      BoolValue(true),
				Add:      ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Multiply: ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Subtract: ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Divide:   ErrorValue{ValueErrorf("unable to convert 'A' to float")},
			}},
		{"string against float",
			StringValue("A"), Float64Value(10.5), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(true),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(false),
				Add:      ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Multiply: ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Subtract: ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Divide:   ErrorValue{ValueErrorf("unable to convert 'A' to float")},
			}},
		{"float against string",
			Float64Value(10.5), StringValue("A"), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(false),
				Geq:      BoolValue(false),
				Lt:       BoolValue(true),
				Leq:      BoolValue(true),
				Add:      ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Multiply: ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Subtract: ErrorValue{ValueErrorf("unable to convert 'A' to float")},
				Divide:   ErrorValue{ValueErrorf("unable to convert 'A' to float")},
			}},
		{"float against same float",
			Float64Value(10.5), Float64Value(10.5), operators{
				Eq:       BoolValue(true),
				Neq:      BoolValue(false),
				Gt:       BoolValue(false),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(true),
				Add:      Float64Value(10.5 + 10.5),
				Multiply: Float64Value(10.5 * 10.5),
				Subtract: Float64Value(10.5 - 10.5),
				Divide:   Float64Value(10.5 / 10.5),
			}},
		{"float against smaller float",
			Float64Value(20.5), Float64Value(10.5), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(true),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(false),
				Add:      Float64Value(20.5 + 10.5),
				Multiply: Float64Value(20.5 * 10.5),
				Subtract: Float64Value(20.5 - 10.5),
				Divide:   Float64Value(20.5 / 10.5),
			}},
		{"float against 0",
			Float64Value(20.5), Float64Value(0), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(true),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(false),
				Add:      Float64Value(20.5),
				Multiply: Float64Value(0),
				Subtract: Float64Value(20.5),
				Divide:   ErrorValue{ErrDivideByZero},
			}},
		{"float against bool true",
			Float64Value(20.5), BoolValue(true), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(false),
				Geq:      BoolValue(false),
				Lt:       BoolValue(true),
				Leq:      BoolValue(true),
				Add:      Float64Value(20.5 + 1),
				Multiply: Float64Value(20.5),
				Subtract: Float64Value(20.5 - 1),
				Divide:   Float64Value(20.5 / 1),
			}},
		{"float against bool false",
			Float64Value(20.5), BoolValue(false), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(false),
				Geq:      BoolValue(false),
				Lt:       BoolValue(true),
				Leq:      BoolValue(true),
				Add:      Float64Value(20.5),
				Multiply: Float64Value(0),
				Subtract: Float64Value(20.5),
				Divide:   ErrorValue{ErrDivideByZero},
			}},
		{"true against float",
			BoolValue(true), Float64Value(20.5), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(true),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(false),
				Add:      Float64Value(20.5 + 1),
				Multiply: Float64Value(20.5),
				Subtract: Float64Value(1 - 20.5),
				Divide:   Float64Value(1 / 20.5),
			}},
		{"false against float",
			BoolValue(false), Float64Value(20.5), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(true),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(false),
				Add:      Float64Value(20.5),
				Multiply: Float64Value(0),
				Subtract: Float64Value(0 - 20.5),
				Divide:   Float64Value(0),
			}},
		{"true against false",
			BoolValue(true), BoolValue(false), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(true),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(false),
				Add:      Float64Value(1),
				Multiply: Float64Value(0),
				Subtract: Float64Value(1),
				Divide:   ErrorValue{ErrDivideByZero},
			}},
		{"false against true",
			BoolValue(false), BoolValue(true), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(false),
				Geq:      BoolValue(false),
				Lt:       BoolValue(true),
				Leq:      BoolValue(true),
				Add:      Float64Value(1),
				Multiply: Float64Value(0),
				Subtract: Float64Value(-1),
				Divide:   Float64Value(0),
			}},
		{"true against true",
			BoolValue(true), BoolValue(true), operators{
				Eq:       BoolValue(true),
				Neq:      BoolValue(false),
				Gt:       BoolValue(false),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(true),
				Add:      Float64Value(2),
				Multiply: Float64Value(1),
				Subtract: Float64Value(0),
				Divide:   Float64Value(1),
			}},
		{"false against false",
			BoolValue(false), BoolValue(false), operators{
				Eq:       BoolValue(true),
				Neq:      BoolValue(false),
				Gt:       BoolValue(false),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(true),
				Add:      Float64Value(0),
				Multiply: Float64Value(0),
				Subtract: Float64Value(0),
				Divide:   ErrorValue{ErrDivideByZero},
			}},
		{"numeric string against float",
			StringValue("1.75"), Float64Value(1.5), operators{
				Eq:       BoolValue(false),
				Neq:      BoolValue(true),
				Gt:       BoolValue(true),
				Geq:      BoolValue(true),
				Lt:       BoolValue(false),
				Leq:      BoolValue(false),
				Add:      Float64Value(1.75 + 1.5),
				Multiply: Float64Value(1.75 * 1.5),
				Subtract: Float64Value(1.75 - 1.5),
				Divide:   Float64Value(1.75 / 1.5),
			}},
		{"float against numeric string",
			Float64Value(1.75), StringValue("1.5"), operators{
				Eq:       BoolValue(false), // not compared numerically
				Neq:      BoolValue(true),
				Gt:       BoolValue(false),
				Geq:      BoolValue(false),
				Lt:       BoolValue(true),
				Leq:      BoolValue(true),
				Add:      Float64Value(1.75 + 1.5),
				Multiply: Float64Value(1.75 * 1.5),
				Subtract: Float64Value(1.75 - 1.5),
				Divide:   Float64Value(1.75 / 1.5),
			}},
		{
			"float against error",
			Float64Value(1.75), ErrorValue{ValueErrorf("bad!")}, operators{
				Eq:       ErrorValue{ValueErrorf("bad!")},
				Neq:      ErrorValue{ValueErrorf("bad!")},
				Gt:       ErrorValue{ValueErrorf("bad!")},
				Geq:      ErrorValue{ValueErrorf("bad!")},
				Lt:       ErrorValue{ValueErrorf("bad!")},
				Leq:      ErrorValue{ValueErrorf("bad!")},
				Add:      ErrorValue{ValueErrorf("bad!")},
				Multiply: ErrorValue{ValueErrorf("bad!")},
				Subtract: ErrorValue{ValueErrorf("bad!")},
				Divide:   ErrorValue{ValueErrorf("bad!")},
			}},
		{
			"error against float",
			ErrorValue{ValueErrorf("bad!")}, Float64Value(1.75), operators{
				Eq:       ErrorValue{ValueErrorf("bad!")},
				Neq:      ErrorValue{ValueErrorf("bad!")},
				Gt:       ErrorValue{ValueErrorf("bad!")},
				Geq:      ErrorValue{ValueErrorf("bad!")},
				Lt:       ErrorValue{ValueErrorf("bad!")},
				Leq:      ErrorValue{ValueErrorf("bad!")},
				Add:      ErrorValue{ValueErrorf("bad!")},
				Multiply: ErrorValue{ValueErrorf("bad!")},
				Subtract: ErrorValue{ValueErrorf("bad!")},
				Divide:   ErrorValue{ValueErrorf("bad!")},
			}},
		{
			"string against error",
			StringValue("my major string"), ErrorValue{ValueErrorf("bad!")}, operators{
				Eq:       ErrorValue{ValueErrorf("bad!")},
				Neq:      ErrorValue{ValueErrorf("bad!")},
				Gt:       ErrorValue{ValueErrorf("bad!")},
				Geq:      ErrorValue{ValueErrorf("bad!")},
				Lt:       ErrorValue{ValueErrorf("bad!")},
				Leq:      ErrorValue{ValueErrorf("bad!")},
				Add:      ErrorValue{ValueErrorf("bad!")},
				Multiply: ErrorValue{ValueErrorf("bad!")},
				Subtract: ErrorValue{ValueErrorf("bad!")},
				Divide:   ErrorValue{ValueErrorf("bad!")},
			}},
		{
			"error against string",
			ErrorValue{ValueErrorf("bad!")}, StringValue("my major string"), operators{
				Eq:       ErrorValue{ValueErrorf("bad!")},
				Neq:      ErrorValue{ValueErrorf("bad!")},
				Gt:       ErrorValue{ValueErrorf("bad!")},
				Geq:      ErrorValue{ValueErrorf("bad!")},
				Lt:       ErrorValue{ValueErrorf("bad!")},
				Leq:      ErrorValue{ValueErrorf("bad!")},
				Add:      ErrorValue{ValueErrorf("bad!")},
				Multiply: ErrorValue{ValueErrorf("bad!")},
				Subtract: ErrorValue{ValueErrorf("bad!")},
				Divide:   ErrorValue{ValueErrorf("bad!")},
			}},
		{
			"bool against error",
			BoolValue(true), ErrorValue{ValueErrorf("bad!")}, operators{
				Eq:       ErrorValue{ValueErrorf("bad!")},
				Neq:      ErrorValue{ValueErrorf("bad!")},
				Gt:       ErrorValue{ValueErrorf("bad!")},
				Geq:      ErrorValue{ValueErrorf("bad!")},
				Lt:       ErrorValue{ValueErrorf("bad!")},
				Leq:      ErrorValue{ValueErrorf("bad!")},
				Add:      ErrorValue{ValueErrorf("bad!")},
				Multiply: ErrorValue{ValueErrorf("bad!")},
				Subtract: ErrorValue{ValueErrorf("bad!")},
				Divide:   ErrorValue{ValueErrorf("bad!")},
			}},
		{
			"error against bool",
			ErrorValue{ValueErrorf("bad!")}, BoolValue(true), operators{
				Eq:       ErrorValue{ValueErrorf("bad!")},
				Neq:      ErrorValue{ValueErrorf("bad!")},
				Gt:       ErrorValue{ValueErrorf("bad!")},
				Geq:      ErrorValue{ValueErrorf("bad!")},
				Lt:       ErrorValue{ValueErrorf("bad!")},
				Leq:      ErrorValue{ValueErrorf("bad!")},
				Add:      ErrorValue{ValueErrorf("bad!")},
				Multiply: ErrorValue{ValueErrorf("bad!")},
				Subtract: ErrorValue{ValueErrorf("bad!")},
				Divide:   ErrorValue{ValueErrorf("bad!")},
			}},
	} {
		for op, expected := range tt.operators {
			name := fmt.Sprintf("%s %s", tt.name, op)
			t.Run(name, func(t *testing.T) {
				out := op.Apply(tt.first, tt.second)
				assert.Equal(t, expected, out)
			})
		}
	}
}
