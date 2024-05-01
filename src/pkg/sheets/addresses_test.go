package sheets

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{":A45", Range{}, "invalid range: expected A23:B54 found ':A45'"},
		{"A:23", Range{}, "invalid range: expected A23:B54 found 'A:23'"},
		{"23:A", Range{}, "invalid range: expected A23:B54 found '23:A'"},
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

func TestRange_Contains(t *testing.T) {
	r := Range{
		StartCol: 15, StartRow: 12,
		EndCol: 25, EndRow: 33,
	}
	for _, tt := range []struct {
		pos      Pos
		expected bool
	}{
		{Pos{Row: 11, Col: 16}, false},
		{Pos{Row: 12, Col: 16}, true},
		{Pos{Row: 22, Col: 23}, true},
		{Pos{Row: 33, Col: 25}, true},
		{Pos{Row: 34, Col: 25}, false},
		{Pos{Row: 33, Col: 26}, false},
	} {
		t.Run(tt.pos.String(), func(t *testing.T) {
			assert.Equal(t, tt.expected, r.Contains(tt.pos))
		})
	}
}

func TestRange_NextPos(t *testing.T) {
	r := mustParseRange(t, "D23:F25")
	pos := r.StartPos()

	var positions []string
	for {
		positions = append(positions, pos.String())
		nextPos, hasNext := r.NextPos(pos)
		if !hasNext {
			break
		}

		pos = nextPos
	}

	assert.Equal(t, []string{
		"D23", "E23", "F23",
		"D24", "E24", "F24",
		"D25", "E25", "F25"}, positions)
}

func TestRange_ContainsRange(t *testing.T) {
	r := mustParseRange(t, "D23:Q33")
	for _, tt := range []struct {
		name     string
		r        Range
		expected bool
	}{
		{"above the top of the range", mustParseRange(t, "E19:F22"), false},
		{"to the right of the range", mustParseRange(t, "Z24:AA29"), false},
		{"below the bottom of the range", mustParseRange(t, "E55:L75"), false},
		{"to the left of the range", mustParseRange(t, "A24:L29"), false},
		{"fully inside the range", mustParseRange(t, "E26:L29"), true},
		{"at the top-left corner of the range", mustParseRange(t, "D23:L29"), true},
		{"at the bottom-left corner of the range", mustParseRange(t, "D26:D33"), true},
		{"at the top-right corner of the range", mustParseRange(t, "Q23:Q26"), true},
		{"at the bottom-right corner of the range", mustParseRange(t, "Q26:Q33"), true},
		{"exactly aligned with the range", mustParseRange(t, "D23:Q33"), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, r.ContainsRange(tt.r))
		})
	}
}

func TestRange_NumCells(t *testing.T) {
	assert.Equal(t, 42, mustParseRange(t, "C4:H10").NumCells())
	assert.Equal(t, 30, mustParseRange(t, "D4:H9").NumCells())
}

func mustParseRange(t *testing.T, s string) Range {
	r, err := ParseRange(s)
	require.NoError(t, err)
	return r
}
