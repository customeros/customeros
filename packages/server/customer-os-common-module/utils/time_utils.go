package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"strings"
	"time"
)

const customLayout1 = "2006-01-02 15:04:05"
const customLayout2 = "2006-01-02T15:04:05.000-0700"
const customLayout3 = "2006-01-02T15:04:05-07:00"
const customLayout4 = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
const customLayout5 = "Mon, 2 Jan 2006 15:04:05 MST"
const customLayout6 = "Mon, 2 Jan 2006 15:04:05 -0700"
const customLayout7 = "Mon, 2 Jan 2006 15:04:05 +0000 (GMT)"
const customLayout8 = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
const customLayout9 = "2 Jan 2006 15:04:05 -0700"

type YearMonth struct {
	Year  int
	Month time.Month
}

func ZeroTime() time.Time {
	return time.Time{}
}

func Now() time.Time {
	return time.Now().UTC()
}

func NowIfZero(t time.Time) time.Time {
	if t.IsZero() {
		return Now()
	}
	return t
}

func TimeOrNowFromPtr(t *time.Time) time.Time {
	if t == nil {
		return Now()
	}
	if t.IsZero() {
		return Now()
	}
	return *t
}

func Today() time.Time {
	return ToDate(Now())
}

func NowPtr() *time.Time {
	return TimePtr(time.Now().UTC())
}

func ConvertTimeToTimestampPtr(input *time.Time) *timestamppb.Timestamp {
	if input == nil {
		return nil
	}
	return timestamppb.New(*input)
}

func ToDate(t time.Time) time.Time {
	val := t.UTC().Truncate(24 * time.Hour)
	return val
}

func ToDatePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	val := t.UTC().Truncate(24 * time.Hour)
	return &val
}

func ToDateAsAny(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	val := t.UTC().Truncate(24 * time.Hour)
	return val
}

func UnmarshalDateTime(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, input)
	if err == nil {
		// Parsed as RFC3339
		return &t, nil
	}

	// Try custom layouts
	customLayouts := []string{customLayout1, customLayout2, customLayout4, customLayout5, customLayout6, customLayout7, customLayout8, customLayout9}

	for _, layout := range customLayouts {
		t, err = time.Parse(layout, input)
		if err == nil {
			return &t, nil
		}
	}
	inputForLayout3 := input
	if !strings.Contains(input, "[UTC]") {
		index := strings.Index(input, "[")
		// If found, strip off the timezone information
		if index != -1 {
			inputForLayout3 = input[:index]
		}
	}
	t, err = time.Parse(customLayout3, inputForLayout3)
	if err == nil {
		return &t, nil
	}

	return nil, errors.New(fmt.Sprintf("cannot parse input as date time %s", input))
}

func TimestampProtoToTime(pbTime *timestamppb.Timestamp) time.Time {
	if pbTime == nil {
		return ZeroTime()
	}
	t := pbTime.AsTime()
	return t
}

func TimestampProtoToTimePtr(pbTime *timestamppb.Timestamp) *time.Time {
	if pbTime == nil {
		return nil
	}
	t := pbTime.AsTime()
	return &t
}

// IsEqualTimePtr compares two *time.Time values and returns true if both are nil or if both point to the same time.
func IsEqualTimePtr(t1, t2 *time.Time) bool {
	// if both are nil, return true
	if t1 == nil && t2 == nil {
		return true
	}
	// if one is nil, return false
	if t1 == nil || t2 == nil {
		return false
	}
	// if both are not nil, compare the time values they point to
	return (*t1).Equal(*t2)
}

// Implement a backoffDelay function that calculates the delay before the next retry.
func BackOffExponentialDelay(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}
	// Calculate the delay with a simple exponential backoff formula
	delay := time.Duration(math.Pow(2, float64(attempt))) * time.Millisecond * 50
	// Cap the delay at 5 seconds
	maxDelay := 5 * time.Second
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}

func BackOffIncrementalDelay(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}
	// Calculate the delay with a simple exponential backoff formula
	delay := time.Duration(attempt) * time.Millisecond * 50
	// Cap the delay at 2 seconds
	maxDelay := 2 * time.Second
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}

func FirstTimeOfMonth(year, month int) time.Time {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
}

func MiddleTimeOfMonth(year, month int) time.Time {
	return FirstTimeOfMonth(year, month).AddDate(0, 0, 15)
}

func LastTimeOfMonth(year, month int) time.Time {
	return FirstTimeOfMonth(year, month).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func LastDayOfMonth(year, month int) time.Time {
	return FirstTimeOfMonth(year, month).AddDate(0, 1, 0).Add(-time.Hour * 24)
}

func StartOfDayInUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func EndOfDayInUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
}

func AddOneMonthFallbackToLastDayOfMonth(date time.Time) time.Time {
	// Calculate the next month
	nextMonth := date.AddDate(0, 1, 0)

	for i := 0; i < 4; i++ {
		if date.Day()-nextMonth.Day() > 27 {
			// Decrease the day by 1
			nextMonth = nextMonth.AddDate(0, 0, -1)
		}
	}

	// Keep the same day
	return nextMonth
}

func GenerateYearMonths(start, end time.Time) []YearMonth {
	yearMonths := []YearMonth{}

	// Set the day to the 15th of the month for the start date
	start = time.Date(start.Year(), start.Month(), 15, 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), 16, 0, 0, 0, 0, end.Location())
	current := start
	for current.Before(end) || current.Equal(end) {
		yearMonths = append(yearMonths, YearMonth{Year: current.Year(), Month: current.Month()})
		current = current.AddDate(0, 1, 0)
	}

	return yearMonths
}

func IsEndOfMonth(t time.Time) bool {
	return t.Day() == LastDayOfMonth(t.Year(), int(t.Month())).Day()
}

func GetCurrentTimeInTimeZone(timezone string) time.Time {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Now()
	}
	return time.Now().In(loc)
}
