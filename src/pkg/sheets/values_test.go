package sheets

import (
	"testing"
	"time"

	"github.com/mmihic/golib/src/pkg/timex"
	"github.com/stretchr/testify/require"
)

func TestToFromStringTime(t *testing.T) {
	tm := timex.MustParseTime(time.RFC3339, "2024-02-12T14:35:46Z")
	tv := TimeValue(tm)

	s, err := tv.AsString()
	require.NoError(t, err)
	require.Equal(t, "2024-02-12T14:35:46Z", s)

	timeFromString, err := StringValue(s).AsTime()
	require.NoError(t, err)
	require.Equal(t, tm, timeFromString)
}

func TestToFromFloat64Time(t *testing.T) {
	tm := timex.MustParseTime(time.RFC3339, "2024-02-12T14:35:46Z")
	tv := TimeValue(tm)

	n, err := tv.AsFloat64()
	require.NoError(t, err)
	require.Equal(t, 45334.6081712963, n)

	timeFromFloat64, err := Float64Value(n).AsTime()
	require.NoError(t, err)
	require.Equal(t, tm, timeFromFloat64)

	tm2, err := tv.AsTime()
	require.NoError(t, err)
	require.Equal(t, tm, tm2)
}

func TestToFromStringFloat64(t *testing.T) {
	n := Float64Value(29834.238934)

	s, err := n.AsString()
	require.NoError(t, err)
	require.Equal(t, "29834.238934", s)

	float64FromString, err := StringValue(s).AsFloat64()
	require.NoError(t, err)
	require.Equal(t, float64(n), float64FromString)
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

	for iter.Next() {
		s, err := iter.Value().AsString()
		require.NoError(t, err)
		require.NoError(t, iter.Err())

		indexes = append(indexes, iter.Index())
		values = append(values, s)
	}

	require.NoError(t, iter.Err())
	require.Equal(t, 4, iter.Len())
	require.Equal(t, []int{0, 1, 2, 3}, indexes)
	require.Equal(t, []string{"A", "B", "C", "D"}, values)
}

func TestSingleValueIter(t *testing.T) {
	iter := SingleValueIter(StringValue("Foo"))
	require.Equal(t, 1, iter.Len())
	require.True(t, iter.Next())
	require.Equal(t, 0, iter.Index())
	require.Equal(t, StringValue("Foo"), iter.Value())
	require.NoError(t, iter.Err())

	for i := 0; i < 10; i++ {
		require.False(t, iter.Next())
		require.NoError(t, iter.Err())
		require.Equal(t, 1, iter.Len())
	}
}
