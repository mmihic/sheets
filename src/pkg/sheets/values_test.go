package sheets

import (
	"context"
	"testing"
	"time"

	"github.com/mmihic/golib/src/pkg/timex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToValue(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected Value
	}{
		{"This is a string", StringValue("This is a string")},
		{`"this is a string"`, StringValue(`"this is a string"`)},
		{"true", BoolValue(true)},
		{"True", BoolValue(true)},
		{"FALSE", BoolValue(false)},
		{"False", BoolValue(false)},
		{"1", Float64Value(1)},
		{"0", Float64Value(0)},
		{"3.1462", Float64Value(3.1462)},
		{"-0.03784", Float64Value(-0.03784)},
		{"12/31/24", TimeValue(timex.MustParseTime(time.RFC3339, "2024-12-31T00:00:00Z"))},
		{"12/24", TimeValue(timex.MustParseTime(time.RFC3339, "2024-12-01T00:00:00Z"))},
		{"2024-09-13T12:36:45Z", TimeValue(timex.MustParseTime(time.RFC3339, "2024-09-13T12:36:45Z"))},
	} {
		t.Run(tt.input, func(t *testing.T) {
			actual := StringToValue(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSliceValueIter(t *testing.T) {
	iter := SliceValueIter([]Value{
		StringValue("A"),
		StringValue("B"),
		StringValue("C"),
		StringValue("D"),
	})
	require.Equal(t, 4, iter.Len())

	var (
		indexes []int
		values  []string
	)

	for iter.Next(context.TODO()) {
		require.NoError(t, iter.Err())

		indexes = append(indexes, iter.Index())
		values = append(values, iter.Value().String())
	}

	require.NoError(t, iter.Err())
	require.Equal(t, 4, iter.Len())
	require.Equal(t, []int{0, 1, 2, 3}, indexes)
	require.Equal(t, []string{"A", "B", "C", "D"}, values)
}

func TestSingleValueIter(t *testing.T) {
	iter := SingleValueIter(StringValue("Foo"))
	require.Equal(t, 1, iter.Len())
	require.True(t, iter.Next(context.TODO()))
	require.Equal(t, 0, iter.Index())
	require.Equal(t, StringValue("Foo"), iter.Value())
	require.NoError(t, iter.Err())

	for i := 0; i < 10; i++ {
		require.False(t, iter.Next(context.TODO()))
		require.NoError(t, iter.Err())
		require.Equal(t, 1, iter.Len())
	}
}
