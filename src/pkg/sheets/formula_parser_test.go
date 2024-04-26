package sheets

import (
	"testing"
	"time"

	"github.com/mmihic/golib/src/pkg/timex"
	"github.com/stretchr/testify/assert"
)

func TestParseFormula(t *testing.T) {
	for _, tt := range []struct {
		input       string
		expected    Formula
		expectedErr string
	}{
		// Functions
		{
			"no_args()", &FunctionCall{
				FunctionName: "NO_ARGS",
			}, "",
		},
		{
			"median(`My Sheet`!A:A)", &FunctionCall{
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
			}, "",
		},
		{
			"VLOOKUP( M23, `Other Sheet`!A1:C45, 1, FALSE )", &FunctionCall{
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
			}, "",
		},
		{
			`split("This is a set of \"quoted\" words")`, &FunctionCall{
				FunctionName: "SPLIT",
				Args: []Formula{
					&Constant{Value: StringValue(`This is a set of "quoted" words`)},
				},
			}, "",
		},

		// Constants
		{
			`100.45`, &Constant{
				Value: Float64Value(100.45),
			}, "",
		},
		{
			`True`, &Constant{
				Value: Float64Value(1),
			}, "",
		},
		{
			`FALSE`, &Constant{
				Value: Float64Value(0),
			}, "",
		},
		{
			`"This is a string"`, &Constant{
				Value: StringValue("This is a string"),
			}, "",
		},
		{
			`"2024-01-14T12:34:56Z"`, &Constant{
				Value: TimeValue(timex.MustParseTime(time.RFC3339, "2024-01-14T12:34:56Z")),
			}, "",
		},
		{
			`"2024/01/14"`, &Constant{
				Value: TimeValue(timex.MustParseTime(time.RFC3339, "2024-01-14T00:00:00Z")),
			}, "",
		},
		{
			`"2019.3746"`, &Constant{
				Value: Float64Value(2019.3746),
			}, "",
		},

		// References
		{
			`A34:C72`, &CellRangeReference{
				Range: Range{
					StartCol: 0, StartRow: 33,
					EndCol: 2, EndRow: 71,
				},
			}, "",
		},
		{
			`"Another Sheet"!A34:C72`, &CellRangeReference{
				Sheet: "Another Sheet",
				Range: Range{
					StartCol: 0, StartRow: 33,
					EndCol: 2, EndRow: 71,
				},
			}, "",
		},
		{
			`B34`, &CellReference{
				Pos: Pos{
					Col: 1,
					Row: 33,
				},
			}, "",
		},
		{
			`MyNamedRange`, &NamedRangeReference{
				NamedRange: "MyNamedRange",
			}, "",
		},
		{
			"YetAnotherSheet!B45", &CellReference{
				Sheet: "YetAnotherSheet",
				Pos: Pos{
					Col: 1,
					Row: 44,
				},
			}, "",
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			ref, err := ParseFormula(tt.input)
			if tt.expectedErr != "" {
				if !assert.Error(t, err) {
					return
				}
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.expected, ref)
		})
	}

}
