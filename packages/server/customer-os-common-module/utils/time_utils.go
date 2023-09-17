package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

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
