package sheets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var matrix = [][]Value{
	{StringValue("cat"), StringValue("hat"), StringValue("mat"), StringValue("gat")},
	{StringValue("apple"), StringValue("banana")},
	{StringValue("mork"), StringValue("bork"), StringValue("rork"), StringValue("dork"), StringValue("zork")},
}

func TestInMemorySheet_Dimensions(t *testing.T) {
	s, err := NewInMemorySheet(matrix)
	require.NoError(t, err)
	assert.Equal(t, Dimensions{
		EndRow: 2,
		EndCol: 4,
	}, s.Dimensions())
}

func TestInMemorySheet_Range(t *testing.T) {
	s, err := NewInMemorySheet(matrix)
	require.NoError(t, err)

	iter, err := s.Range(context.TODO(), mustParseRange(t, "B1:C3"))
	require.NoError(t, err)
	assert.Equal(t, 6, iter.Len())

	var (
		indexes   []int
		positions []string
		values    []Value
	)

	for iter.Next(context.TODO()) {
		indexes = append(indexes, iter.Index())
		positions = append(positions, iter.Pos().String())
		values = append(values, iter.Value())
	}

	require.NoError(t, iter.Err())
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, indexes)
	assert.Equal(t, []string{"B1", "C1", "B2", "C2", "B3", "C3"}, positions)
	assert.Equal(t, []Value{
		StringValue("hat"), StringValue("mat"),
		StringValue("banana"), StringValue(""),
		StringValue("bork"), StringValue("rork")},
		values)
	assert.Equal(t, 6, iter.Len())

	// invalid end of range
	_, err = s.Range(context.TODO(), mustParseRange(t, "A1:Z3"))
	require.Error(t, err)
	assert.Contains(t, "position 'Z3' outside of bounds", err.Error())

	_, ok := err.(InvalidPosError)
	require.True(t, ok)

	// invalid start of range
	_, err = s.Range(context.TODO(), mustParseRange(t, "A4:C4"))
	require.Error(t, err)
	assert.Contains(t, "position 'A4' outside of bounds", err.Error())

	_, ok = err.(InvalidPosError)
	require.True(t, ok)

}

func TestInMemorySheet_Get(t *testing.T) {
	s, err := NewInMemorySheet(matrix)
	require.NoError(t, err)

	v, err := s.Get(context.TODO(), mustParsePos(t, "A1"))
	require.NoError(t, err)
	require.Equal(t, StringValue("cat"), v)

	// Doesn't explicitly exist
	v, err = s.Get(context.TODO(), mustParsePos(t, "C2"))
	require.NoError(t, err)
	require.Equal(t, StringValue(""), v)

	v, err = s.Get(context.TODO(), mustParsePos(t, "E3"))
	require.NoError(t, err)
	require.Equal(t, StringValue("zork"), v)

	// Invalid position
	_, err = s.Get(context.TODO(), mustParsePos(t, "Q3"))
	require.Error(t, err)
	require.Contains(t, "position 'Q3' outside of bounds", err.Error())

	_, ok := err.(InvalidPosError)
	require.True(t, ok)
}

func mustParsePos(t *testing.T, s string) Pos {
	pos, err := ParsePos(s)
	require.NoError(t, err)
	return pos
}
