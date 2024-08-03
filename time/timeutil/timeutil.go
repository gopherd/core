package timeutil

import (
	"errors"
	"time"
)

// ErrUnrecognizedTime is the error that parsing failed
var ErrUnrecognizedTime = errors.New("time: unrecognized time")

var (
	offset time.Duration
)

// SetOffset sets offset to current time
func SetOffset(off time.Duration) {
	offset = off
}

// StdNow returns standard library current time
func StdNow() time.Time {
	return time.Now()
}

// Now returns current time with offset
func Now() time.Time {
	return StdNow().Add(offset)
}

// supported layouts for parsing time from string
var layouts = []string{
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"15:04:05 2006-01-02",
	"15:04:05 2006/01/02",
	"01-02-2006 15:04:05",
	"01/02/2006 15:04:05",
	"15:04:05 01-02-2006",
	"15:04:05 01/02/2006",
	"2006-01-02",
	"2006/01/02",
	"2006:01:02",
	"15:04:05",
	"2006-1-2 15:4:5",
	"2006/1/2 15:4:5",
	"15:4:5 2006-1-2",
	"15:4:5 2006/1/2",
	"1-2-2006 15:4:5",
	"1/2/2006 15:4:5",
	"15:4:5 1-2-2006",
	"15:4:5 1/2/2006",
	"2006-1-2",
	"2006/1/2",
	"2006:1:2",
	"15:4:5",
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	"2006-01-02T15:04:05.999Z07:00",
	"2006-01-02T15:04:05.999999Z07:00",
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

// Parse parse string as time. If the layout of time is known,
// please use standard library time.Parse instead.
func Parse(value string) (time.Time, error) {
	return parse(value, time.Local)
}

// ParseInLocation is like Parse but differs in two important ways.
// First, in the absence of time zone information, Parse interprets a time as UTC;
// ParseInLocation interprets the time as in the given location.
// Second, when given a zone offset or abbreviation, Parse tries to match it
// against the Local location; ParseInLocation uses the given location.
func ParseInLocation(value string, loc *time.Location) (time.Time, error) {
	return parse(value, loc)
}

func parse(value string, loc *time.Location) (time.Time, error) {
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, value, loc)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, ErrUnrecognizedTime
}

// Timestamp is unix timestamp milliseconds
type Timestamp int64

// Seconds returns seconds of timestamp
func (t Timestamp) Seconds() int64 {
	return int64(t) / 1e3
}

// Milliseconds returns milliseconds of timestamp
func (t Timestamp) Milliseconds() int64 {
	return int64(t)
}

// Time returns time of timestamp
func (t Timestamp) Time() time.Time {
	return time.Unix(t.Seconds(), int64(t)%1e3*1e6)
}

// String returns string of timestamp
func (t Timestamp) String() string {
	return t.Time().Format("2006-01-02T15:04:05.999Z07:00")
}

// GetTimestamp returns timestamp of t
func GetTimestamp(t time.Time) Timestamp {
	return Timestamp(t.UnixNano() / 1e6)
}
