package sheets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormulaString(t *testing.T) {
	for _, tt := range []struct {
		expected string
		input    Formula
	}{
		{
			"NO_ARGS()", &FunctionCall{
				FunctionName: "NO_ARGS",
			}},
		{
			"MEDIAN(`My Sheet`!A:A)", &FunctionCall{
				FunctionName: "MEDIAN",
				Args: []Formula{
					&CellRangeReference{
						Sheet: "My Sheet",
						Range: Range{
							StartCol: 0, StartRow: 0,
							EndCol: 0, EndRow: MaxRow,
						},
					},
				},
			}},
		{
			"VLOOKUP(M23, `Other Sheet`!A1:C45, 1, 0)", &FunctionCall{
				FunctionName: "VLOOKUP",
				Args: []Formula{
					&CellReference{
						Pos: Pos{
							Col: 12,
							Row: 22,
						},
					},
					&CellRangeReference{
						Sheet: "Other Sheet",
						Range: Range{
							StartCol: 0, StartRow: 0,
							EndCol: 2, EndRow: 44,
						},
					},
					&Constant{
						Value: Float64Value(1),
					},
					&Constant{
						Value: Float64Value(0),
					},
				},
			}},
		{
			`SPLIT("This is a set of \"quoted\" words")`, &FunctionCall{
				FunctionName: "SPLIT",
				Args: []Formula{
					&Constant{Value: StringValue(`This is a set of "quoted" words`)},
				},
			}},
	} {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.String())
		})
	}
}
