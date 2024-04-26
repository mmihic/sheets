package sheets

import (
	"testing"
	"time"

	"github.com/mmihic/golib/src/pkg/timex"
	"github.com/stretchr/testify/assert"
)

const (
	rfc3339Milli = "2006-01-02T15:04:05.999Z07:00"
)

func TestToFromTime(t *testing.T) {
	for _, tt := range []struct {
		fraction float64
		expected time.Time
	}{
		{41_994.523784722230000, timex.MustParseTime(rfc3339Milli, "2014-12-21T12:34:15.00Z")},
		{45_401.259826388888889, timex.MustParseTime(rfc3339Milli, "2024-04-19T06:14:09.00Z")},
		{34_638.924664351851852, timex.MustParseTime(rfc3339Milli, "1994-10-31T22:11:31.00Z")},
	} {
		t.Run(tt.expected.String(), func(t *testing.T) {
			tm, err := FromExcelTime(tt.fraction)
			if !assert.NoError(t, err) {
				return
			}

			if !assert.Equal(t, tt.expected, tm) {
				return
			}

			fraction, err := ToExcelTime(tm)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.fraction, fraction)
		})
	}
}

func TestDate(t *testing.T) {
	for _, tt := range []struct {
		fraction float64
		expected time.Time
	}{
		{1.00000000000000, timex.MustParseTime(rfc3339Milli, "1899-12-31T12:34:56.987Z")},
		{2.00000000000000, timex.MustParseTime(rfc3339Milli, "1900-01-01T12:34:56.987Z")},
		{61.00000000000000, timex.MustParseTime(rfc3339Milli, "1900-03-01T12:34:56.987Z")},
		{60.00000000000000, timex.MustParseTime(rfc3339Milli, "1900-02-28T12:34:56.987Z")},
		{41_994.00000000000000, timex.MustParseTime(rfc3339Milli, "2014-12-21T12:34:56.987Z")},
		{45_401.00000000000000, timex.MustParseTime(rfc3339Milli, "2024-04-19T12:34:56.987Z")},
	} {
		t.Run(tt.expected.String(), func(t *testing.T) {
			day, month, year := fromExcelDate(tt.fraction)
			assert.Equal(t, tt.expected, time.Date(year, month, day,
				12, 34, 56, 987000000, time.UTC))

			fraction := toExcelDate(day, month, year)
			assert.Equal(t, tt.fraction, fraction)
		})
	}
}

func TestTimeOfDay(t *testing.T) {
	for _, tt := range []struct {
		fraction float64
		expected time.Time
	}{
		{0.523785000000000, timex.MustParseTime(rfc3339Milli, "2024-02-12T12:34:15.024Z")},
		{0.259826388888889, timex.MustParseTime(rfc3339Milli, "2024-02-12T06:14:09.00Z")},
		{0.924664351851852, timex.MustParseTime(rfc3339Milli, "2024-02-12T22:11:31.00Z")},
		{0.000000000000000, timex.MustParseTime(rfc3339Milli, "2024-02-12T00:00:00.00Z")},
		{0.999999988425926, timex.MustParseTime(rfc3339Milli, "2024-02-12T23:59:59.999Z")},
	} {
		t.Run(tt.expected.String(), func(t *testing.T) {
			hour, minute, second, nanos := fromExcelTimeOfDay(tt.fraction)
			assert.Equal(t, tt.expected, time.Date(2024, time.February, 12,
				hour, minute, second, nanos, time.UTC))
			fraction := toExcelTimeOfDay(
				tt.expected.Hour(),
				tt.expected.Minute(),
				tt.expected.Second(),
				tt.expected.Nanosecond())
			assert.Equal(t, tt.fraction, fraction)
		})
	}
}
