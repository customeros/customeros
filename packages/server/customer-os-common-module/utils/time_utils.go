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

func ZeroTime() time.Time {
	return time.Time{}
}

func Now() time.Time {
	return time.Now().UTC()
}

func NowAsPtr() *time.Time {
	return TimePtr(time.Now().UTC())
}

func ConvertTimeToTimestampPtr(input *time.Time) *timestamppb.Timestamp {
	if input == nil {
		return nil
	}
	return timestamppb.New(*input)
}

func ToDateNillable(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	y, m, d := t.Date()
	val := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return &val
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
	t, err = time.Parse(customLayout1, input)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(customLayout2, input)
	if err == nil {
		return &t, nil
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

func TimestampProtoToTime(pbTime *timestamppb.Timestamp) *time.Time {
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
