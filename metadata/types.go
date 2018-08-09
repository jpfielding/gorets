package metadata

import (
	"strconv"
	"strings"
	"time"
)

//------------------------------------------

// AlphaNum ...
type AlphaNum string

// Text ...
type Text string

// PlainText ...
type PlainText string

// StringList ...
type StringList string

// List ...
func (s StringList) List() []string {
	return strings.Split(string(s), ",")
}

//------------------------------------------

// RETSID ...
type RETSID string

// RETSName ...
type RETSName string

// RETSNames ...
type RETSNames string

// Parse ...
func (r RETSNames) Parse() []string {
	return strings.Split(string(r), ",")
}

// List ...
func (r RETSNames) List() []RETSName {
	list := r.Parse()
	tmp := make([]RETSName, len(list))
	for i, n := range list {
		tmp[i] = RETSName(n)
	}
	return tmp
}

// ResourceClassName ...
type ResourceClassName string

// Parse ...
func (rc ResourceClassName) Parse() (RETSID, RETSName) {
	split := strings.Split(string(rc), ":")
	return RETSID(split[0]), RETSName(split[1])
}

//------------------------------------------

// Version ...
type Version string

// Parse ...
func (v Version) Parse(index int) Numeric {
	return Numeric(strings.Split(string(v), ".")[index])
}

// Major ...
func (v Version) Major() Numeric {
	return v.Parse(0)
}

// Minor ...
func (v Version) Minor() Numeric {
	return v.Parse(1)
}

// Release ...
func (v Version) Release() Numeric {
	return v.Parse(2)
}

//------------------------------------------

// TODO consider a high order function to hide the format and tz to process DateTime
// TODO consider a func that produces a time based on matching the best time format
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

// Numeric decodes the numeric value
type Numeric string

// Parse ...
func (n Numeric) Parse() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// NumericList ...
type NumericList string

// Parse ...
func (n NumericList) Parse() []string {
	return strings.Split(string(n), ",")
}

// List ...
func (n NumericList) List() []Numeric {
	list := n.Parse()
	tmp := make([]Numeric, len(list))
	for i, n := range list {
		tmp[i] = Numeric(n)
	}
	return tmp
}
