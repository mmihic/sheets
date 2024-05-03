package sheets

import (
	"testing"
	"time"

	"github.com/mmihic/golib/src/pkg/timex"
	"github.com/stretchr/testify/assert"
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
