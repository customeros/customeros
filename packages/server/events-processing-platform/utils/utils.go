package utils

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"strings"
	"time"
)

func CheckErrMessages(err error, messages ...string) bool {
	for _, message := range messages {
		if strings.Contains(strings.TrimSpace(strings.ToLower(err.Error())), strings.TrimSpace(strings.ToLower(message))) {
			return true
		}
	}
	return false
}

func TimestampProtoToTime(pbTime *timestamp.Timestamp) *time.Time {
	if pbTime == nil {
		return nil
	}
	t := pbTime.AsTime()
	return &t
}
