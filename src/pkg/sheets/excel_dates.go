package sheets

import (
	"fmt"
	"math"
	"time"
)

const (
	microsInNanos  = int64(time.Microsecond / time.Nanosecond)
	secsInNanos    = int64(time.Second / time.Nanosecond)
	minutesInNanos = int64(time.Minute / time.Nanosecond)
	hoursInNanos   = int64(time.Hour / time.Nanosecond)
	daysInNanos    = int64((time.Hour * 24) / time.Nanosecond)

	roundEpsilon = 1e-9
)

var (
	epochStart = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
)

// ParseTime parses a string as a time, using the support date and time formats.
func ParseTime(s string) (time.Time, error) {
	for _, layout := range supportedTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("'%s' cannot be parsed as a date or time", s)
}

// FromExcelTime converts an Excel fractional datetime to the appropriate golang time.Time.
func FromExcelTime(n float64) time.Time {
	datePart := float64(int64(n))
	day, month, year := fromExcelDate(datePart)

	timePart := n - datePart + roundEpsilon
	hour, minute, second, nano := fromExcelTimeOfDay(timePart)
	return time.Date(year, month, day, hour, minute, second, nano, time.UTC).Truncate(time.Second)
}

// fromExcelDate converts a fraction whose integer portion is an Excel numeric
// date into the appropriate day, month, and year. Excel represents dates
// as the number of days since Dec 31, 1899 (where Dec 31, 1899 is the value 1).
func fromExcelDate(fraction float64) (day int, month time.Month, year int) {
	diffSinceEpoch := time.Duration(int64(fraction) * daysInNanos)
	date := epochStart.Add(diffSinceEpoch)
	return date.Day(), date.Month(), date.Year()
}

// fromExcelTimeOfDay whose decimal portion is an Excel numeric time into the appropriate hour,
// minute, second, and nanosecond. Excel represents times as a percentage of the day.
func fromExcelTimeOfDay(fraction float64) (hour, minute, second, nanosecond int) {
	// Convert the decimal portion of the fraction into the number of nanoseconds
	// since the start of the day
	nanosSinceStartOfDay := int64(float64(daysInNanos)*fraction + float64(microsInNanos)/2)
	nanosecond = int(((nanosSinceStartOfDay % secsInNanos) / microsInNanos) * microsInNanos)
	second = int((nanosSinceStartOfDay / secsInNanos) % 60)
	minute = int((nanosSinceStartOfDay / minutesInNanos) % 60)
	hour = int((nanosSinceStartOfDay / hoursInNanos) % 60)
	return
}

// ToExcelTime converts a golang time.Time to the Excel fractional date equivalent.
func ToExcelTime(tm time.Time) float64 {
	tm = tm.UTC()

	var (
		dayFraction  = toExcelDate(tm.Day(), tm.Month(), tm.Year())
		timeFraction = toExcelTimeOfDay(tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond())
	)

	return roundFloat(dayFraction+timeFraction, 15)
}

// toExcelDate converts a (day, month, year) into the internal numeric format
// used by Excel to store dates.
func toExcelDate(day int, month time.Month, year int) float64 {
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	diffSinceEpoch := int64(date.Sub(epochStart))
	return float64(diffSinceEpoch / daysInNanos)
}

// toExcelTimeOfDay converts the given (hour, minute, second, nanosecond) to a fraction
// whose decimal value is the percentage of the day at that timestamp.
func toExcelTimeOfDay(hour, minute, second, nanosecond int) float64 {
	return roundFloat(float64(int64(hour)*hoursInNanos+
		int64(minute)*minutesInNanos+
		int64(second)*secsInNanos+
		int64(nanosecond))/float64(daysInNanos), 15)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

var (
	supportedTimeLayouts = []string{
		time.RFC3339,
		time.DateOnly,
		time.DateTime,
		time.UnixDate,
		time.Kitchen,
		"2006/01/02",
		"01/02/2006",
		"01/02/06",
		"01/06",
	}
)
