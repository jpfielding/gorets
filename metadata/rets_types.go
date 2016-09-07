package metadata

import (
	"strconv"
	"strings"
	"time"
)

//------------------------------------------

// StringList ...
type StringList string

// Parse ...
func (sl StringList) Parse() []string {
	return strings.Split(string(sl), ",")
}

//------------------------------------------

// Version ...
type Version string

// Major ...
func (mv Version) Major() string {
	return strings.Split(string(mv), ".")[0]
}

// Minor ...
func (mv Version) Minor() string {
	return strings.Split(string(mv), ".")[1]
}

// Release ...
func (mv Version) Release() string {
	return strings.Split(string(mv), ".")[2]
}

//------------------------------------------

// TOOD consider a high order function to hide the format and tz to process DateTime
// TOOD consider a func that produces a time based on matching the best time format
const (
	// RETSDateTimeFormat is the simple date format for most rets dates
	RETSDateTimeFormat = "2006-01-02T15:04:05Z"
	// RETSDateTimeMiilliFormat is the date format for rets dates with millis
	RETSDateTimeMiilliFormat = "2006-01-02T15:04:05.000Z"
)

// DateTime ...
type DateTime string

// Parse ...
func (dt DateTime) Parse(format string, tz *time.Location) (time.Time, error) {
	return time.ParseInLocation(format, string(dt), tz)
}

//------------------------------------------

// Boolean decodes the 1 or 0 to t/f values
type Boolean string

// Parse ...
func (b Boolean) Parse() (bool, error) {
	// TODO consider some error conditions
	return string(b) == "1", nil
}

//------------------------------------------

// Number decodes the numeric value
type Number string

// Parse ...
func (n Number) Parse() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}
