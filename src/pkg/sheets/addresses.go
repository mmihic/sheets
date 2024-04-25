package sheet

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Pos is the position of a cell in a sheet. Uses 0-based
// indexing.
type Pos struct {
	Row, Col int
}

// String returns the string form of a Pos.
func (pos Pos) String() string {
	return columnToString(pos.Col) + rowToString(pos.Row)
}

// ParsePos parses a position in the form "AA23" where "AA"
// is the column and 23 is the row in 1-based indexing.
func ParsePos(s string) (Pos, error) {
	elts := rePosition.FindStringSubmatch(s)
	if len(elts) != 3 {
		return Pos{}, fmt.Errorf("invalid range: expected A23 found '%s'", s)
	}

	return Pos{
		Col: columnOffset(elts[1]),
		Row: rowOffset(elts[2]),
	}, nil
}

// A Range is a description of a collection of cells in a sheet.
// Uses 0-based indexing.
type Range struct {
	StartRow, EndRow int
	StartCol, EndCol int
}

// ParseRange parses a range in the form "AA23:BC45" where
// "AA23" is the starting position of the range and "BC45"
// is the ending position of the range. Supports four formats:
//
// * AA23:BC54 - covers all cells in columns AA-BC and rows 23-54
// * AA:CBC    - covers all rows in columns AA-BC
// * 23:54     - covers rows 23-54 in every column
func ParseRange(s string) (Range, error) {
	elts := reRange.FindStringSubmatch(s)
	if len(elts) != 5 {
		return Range{}, fmt.Errorf("invalid range: expected A23:B54 found '%s'", s)
	}

	startCol, endCol, startRow, endRow := elts[1], elts[3], elts[2], elts[4]
	r := Range{
		StartCol: 0,
		EndCol:   MaxColumn,
		StartRow: 0,
		EndRow:   MaxRow,
	}

	if len(startCol) != 0 {
		r.StartCol = columnOffset(startCol)
	}

	if len(endCol) != 0 {
		r.EndCol = columnOffset(endCol)
	}

	if len(startRow) != 0 {
		r.StartRow = rowOffset(startRow)
	}

	if len(endRow) != 0 {
		r.EndRow = rowOffset(endRow)
	}

	return r, nil
}

// String returns the string form of the range
func (r Range) String() string {
	var (
		startRow, endRow string
		startCol, endCol string
	)

	// Only show the start row if this is not a column-only range
	if r.StartRow != 0 || r.EndRow != MaxRow {
		startRow = rowToString(r.StartRow)
	}

	// Only show the end row if it doesn't cover the entire sheet
	if r.EndRow != MaxRow {
		endRow = rowToString(r.EndRow)
	}

	// Only show the start column if this is not a row-only range
	if r.StartCol != 0 || r.EndCol != MaxColumn {
		startCol = columnToString(r.StartCol)
	}

	// Only show the end column if it doesn't cover the entire sheet
	if r.EndCol != MaxColumn {
		endCol = columnToString(r.EndCol)
	}

	return fmt.Sprintf("%s%s:%s%s", startCol, startRow, endCol, endRow)
}

const (
	// MaxRow is used as the value of EndRow to indicate that the range
	// covers up to the maximum row in the sheet.
	MaxRow = -1

	// MaxColumn is used as the value of EndColumn to indicate that the range
	// covers up to the maximum column in the sheet.
	MaxColumn = -1
)

func columnOffset(colText string) int {
	colText = strings.ToUpper(colText)

	col := 0
	for i, r := range colText {
		if i != 0 {
			col = (col + 1) * 26
		}
		col += int(r - 'A')
	}

	return col
}

func rowToString(row int) string {
	return strconv.Itoa(row + 1)
}

func columnToString(column int) string {
	var r []rune
	for column >= 0 {
		r = append(r, 'A'+rune(column%26))
		column /= 26
		column--
	}

	slices.Reverse(r)
	return string(r)
}

func rowOffset(rowText string) int {
	// NB(mmihic): This is only ever called for row values that have already been validated
	row1Index, _ := strconv.Atoi(rowText)
	return row1Index - 1 // convert to 0-based index
}

var (
	rePosition = regexp.MustCompile(`^([A-Za-z]{1,3})(\d+)$`)
	reRange    = regexp.MustCompile(`^([A-Za-z]{1,3})?(\d+)?\s*:\s*([A-Za-z]{1,3})?(\d+)?$`)
)
