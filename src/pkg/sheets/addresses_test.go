package sheets

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPos(t *testing.T) {
	for _, tt := range []struct {
		input       string
		expected    Pos
		expectedErr string
	}{
		{"A23", Pos{Row: 22, Col: 0}, ""},
		{"CD45", Pos{Row: 44, Col: 81}, ""},
		{"45", Pos{}, "invalid range: expected A23 found '45'"},
		{"AA", Pos{}, "invalid range: expected A23 found 'AA'"},
		{"#$^A34", Pos{}, "invalid range: expected A23 found '#$^A34'"},
		{"", Pos{}, "invalid range: expected A23 found ''"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			pos, err := ParsePos(tt.input)

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
			assert.Equal(t, tt.expected, pos)
			assert.Equal(t, strings.ToUpper(tt.input), pos.String())
		})
	}
}

func TestRange_String(t *testing.T) {
	for _, tt := range []struct {
		input    Range
		expected string
	}{
		{Range{
			StartCol: 26, StartRow: 22,
			EndCol: 54, EndRow: 44,
		}, "AA23:BC45"},
		{Range{
			StartCol: 2, StartRow: 0,
			EndCol: 4, EndRow: MaxRow,
		}, "C:E"},
		{Range{
			StartCol: 0, StartRow: 22,
			EndCol: MaxColumn, EndRow: 44,
		}, "23:45"},
		{Range{
			StartCol: 3, StartRow: 22,
			EndCol: MaxColumn, EndRow: 44,
		}, "D23:45"},
		{Range{
			StartCol: 0, StartRow: 22,
			EndCol: 2, EndRow: 44,
		}, "A23:C45"},
	} {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.String())
		})
	}

}

func TestParseRange(t *testing.T) {
	for _, tt := range []struct {
		input       string
		expected    Range
		expectedErr string
	}{
		{"AA23:BC45", Range{
			StartCol: 26, StartRow: 22,
			EndCol: 54, EndRow: 44,
		}, ""},
		{"C:E", Range{
			StartCol: 2, StartRow: 0,
			EndCol: 4, EndRow: MaxRow,
		}, ""},
		{"23:45", Range{
			StartCol: 0, StartRow: 22,
			EndCol: MaxColumn, EndRow: 44,
		}, ""},
		{"d23:45", Range{
			StartCol: 3, StartRow: 22,
			EndCol: MaxColumn, EndRow: 44,
		}, ""},
		{"23:c45", Range{
			StartCol: 0, StartRow: 22,
			EndCol: 2, EndRow: 44,
		}, ""},
		{"A23", Range{}, "invalid range: expected A23:B54 found 'A23'"},
		{":A45", Range{}, "invalid range: expected A23:B54 found ':A43'"},
		{":", Range{}, "invalid range: expected A23:B54 found ':'"},
		{"25A:A45", Range{}, "invalid range: expected A23:B54 found '25A:A45'"},
		{"A25:45B", Range{}, "invalid range: expected A23:B54 found 'A25:45B'"},
		{"#%$a76:54c", Range{}, "invalid range: expected A23:B54 found '#%$a76:54c'"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			r, err := ParseRange(tt.input)

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
			assert.Equal(t, tt.expected, r)
		})
	}
}
