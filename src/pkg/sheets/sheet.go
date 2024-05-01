package sheets

import (
	"context"
	"fmt"
)

// InvalidPosError is returned when attempting to access a position outside
// the bounds of a sheet.
type InvalidPosError struct {
	Pos Pos
}

// Error returns the error message.
func (e InvalidPosError) Error() string {
	return fmt.Sprintf("position '%s' outside of bounds", e.Pos)
}

// A ValueRange is an iterator that can report its position in the sheet
type ValueRange interface {
	ValueIter
	Pos() Pos
}

type valueRange struct {
	bounds     Range
	currentPos Pos
	index      int
	value      Value
	sheet      Sheet
	err        error
}

func (v *valueRange) Next(ctx context.Context) bool {
	if v.err != nil {
		return false
	}

	if v.value == nil {
		// This is the first retrieval
		v.value, v.err = v.sheet.Get(ctx, v.currentPos)
		return v.err == nil
	}

	nextPos, hasNext := v.bounds.NextPos(v.currentPos)
	if !hasNext {
		return false
	}

	v.currentPos = nextPos
	v.index++
	v.value, v.err = v.sheet.Get(ctx, v.currentPos)
	return v.err == nil
}

func (v *valueRange) Err() error {
	return v.err
}

func (v *valueRange) Value() Value {
	return v.value
}

func (v *valueRange) Index() int {
	return v.index
}

func (v *valueRange) Len() int {
	return v.bounds.NumCells()
}

func (v *valueRange) Pos() Pos {
	return v.currentPos
}

var (
	_ ValueRange = &valueRange{}
)

// Dimensions are the dimensions of a sheet.
type Dimensions struct {
	EndRow, EndCol int
}

// A Sheet is allows arbitrary access to a matrix of cells.
type Sheet interface {
	Get(ctx context.Context, pos Pos) (Value, error)
	Range(ctx context.Context, r Range) (ValueRange, error)
	Dimensions() Dimensions
}

// NewInMemorySheet creates a sheet that wraps a two-dimensional matrix.
func NewInMemorySheet(values [][]Value) (Sheet, error) {
	endCol := 0
	for _, row := range values {
		if len(row) >= endCol {
			endCol = len(row) - 1
		}
	}

	return &inMemorySheet{
		dims: Dimensions{
			EndRow: len(values) - 1,
			EndCol: endCol,
		},
		values: values,
	}, nil
}

type inMemorySheet struct {
	dims   Dimensions
	values [][]Value
}

func (s *inMemorySheet) Dimensions() Dimensions {
	return s.dims
}

func (s *inMemorySheet) Get(_ context.Context, pos Pos) (Value, error) {
	if !s.fullRange().Contains(pos) {
		return nil, InvalidPosError{pos}
	}

	row := s.values[pos.Row]
	if pos.Col >= len(row) {
		return StringValue(""), nil
	}

	return row[pos.Col], nil
}

func (s *inMemorySheet) Range(_ context.Context, r Range) (ValueRange, error) {
	fullRange := s.fullRange()
	if !fullRange.Contains(r.StartPos()) {
		return nil, InvalidPosError{r.StartPos()}
	}

	if !fullRange.Contains(r.EndPos()) {
		return nil, InvalidPosError{r.EndPos()}
	}

	return &valueRange{
		sheet:  s,
		bounds: r,
		currentPos: Pos{
			Row: r.StartRow,
			Col: r.StartCol,
		},
	}, nil
}

func (s *inMemorySheet) fullRange() Range {
	return Range{
		StartRow: 0, StartCol: 0,
		EndRow: s.dims.EndRow, EndCol: s.dims.EndCol,
	}
}

// A DataSet is a set of sheets, functions, and named ranges.
type DataSet interface {
	Sheet(name string) Sheet
}
