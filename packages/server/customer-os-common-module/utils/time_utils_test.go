package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestZeroTime(t *testing.T) {
	expected := time.Time{}
	actual := ZeroTime()

	if !actual.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestNow(t *testing.T) {
	now := Now()
	if now.Location() != time.UTC {
		t.Errorf("Now() should be in UTC, but got %s", now.Location())
	}

	if time.Since(now) > time.Second {
		t.Errorf("Now() is not returning the current time")
	}
}

func TestNowPtr(t *testing.T) {
	nowPtr := NowPtr()
	if nowPtr == nil {
		t.Fatal("NowPtr() returned nil")
	}

	if nowPtr.Location() != time.UTC {
		t.Errorf("NowPtr() should be in UTC, but got %s", nowPtr.Location())
	}

	if time.Since(*nowPtr) > time.Second {
		t.Errorf("NowPtr() is not returning the current time")
	}
}

func TestConvertTimeToTimestampPtr(t *testing.T) {
	// Test with a non-nil time
	testTime := time.Now()
	result := ConvertTimeToTimestampPtr(&testTime)
	if result == nil {
		t.Fatal("ConvertTimeToTimestampPtr returned nil for non-nil input")
	}
	if result.Seconds != testTime.Unix() {
		t.Errorf("Expected seconds %v, got %v", testTime.Unix(), result.Seconds)
	}
	if result.Nanos != int32(testTime.Nanosecond()) {
		t.Errorf("Expected nanos %v, got %v", testTime.Nanosecond(), result.Nanos)
	}

	// Test with a nil time
	resultNil := ConvertTimeToTimestampPtr(nil)
	if resultNil != nil {
		t.Fatal("ConvertTimeToTimestampPtr should return nil for nil input")
	}
}

func TestToDatePtr(t *testing.T) {
	// Test with a non-nil time
	now := time.Now()
	datePtr := ToDatePtr(&now)
	if datePtr == nil {
		t.Fatal("ToDatePtr returned nil for non-nil input")
	}
	if !datePtr.Equal(now.Truncate(24 * time.Hour).UTC()) {
		t.Errorf("Expected %v, got %v", now.Truncate(24*time.Hour).UTC(), *datePtr)
	}

	// Test with a nil time
	nilDatePtr := ToDatePtr(nil)
	if nilDatePtr != nil {
		t.Fatal("ToDatePtr should return nil for nil input")
	}
}

func TestUnmarshalDateTime(t *testing.T) {
	customLayout1 := "2006-01-02 15:04:05"
	customLayout2 := "2006-01-02T15:04:05.000-0700"
	customLayout3 := "2006-01-02T15:04:05-07:00"

	// Test valid RFC3339 input
	rfc3339Input := "2006-01-02T15:04:05Z"
	dt, err := UnmarshalDateTime(rfc3339Input)
	if err != nil {
		t.Errorf("UnmarshalDateTime returned an error for valid RFC3339 input: %v", err)
	}
	if dt == nil || dt.Format(time.RFC3339) != rfc3339Input {
		t.Errorf("Expected %s, got %v", rfc3339Input, dt)
	}

	// Test with custom layout 1
	custom1Input := "2006-01-02 15:04:05"
	custom1Dt, custom1Err := UnmarshalDateTime(custom1Input)
	if custom1Err != nil || custom1Dt == nil || custom1Dt.Format(customLayout1) != custom1Input {
		t.Errorf("UnmarshalDateTime failed for custom layout 1: %v", custom1Err)
	}

	// Test with custom layout 2
	custom2Input := "2006-01-02T15:04:05.000-0700"
	custom2Dt, custom2Err := UnmarshalDateTime(custom2Input)
	if custom2Err != nil || custom2Dt == nil || custom2Dt.Format(customLayout2) != custom2Input {
		t.Errorf("UnmarshalDateTime failed for custom layout 2: %v", custom2Err)
	}

	// Test with custom layout 3
	custom3Input := "2006-01-02T15:04:05-07:00"
	custom3Dt, custom3Err := UnmarshalDateTime(custom3Input)
	if custom3Err != nil || custom3Dt == nil || custom3Dt.Format(customLayout3) != custom3Input {
		t.Errorf("UnmarshalDateTime failed for custom layout 3: %v", custom3Err)
	}

	// Test with empty input
	emptyDt, emptyErr := UnmarshalDateTime("")
	if emptyErr != nil || emptyDt != nil {
		t.Errorf("Expected nil for empty input, got %v and error %v", emptyDt, emptyErr)
	}

	// Test with invalid input
	invalidInput := "invalid-date"
	invalidDt, invalidErr := UnmarshalDateTime(invalidInput)
	if invalidErr == nil {
		t.Errorf("Expected error for invalid input, got %v", invalidDt)
	}
}

func TestTimestampProtoToTimePtr(t *testing.T) {
	// Test with a non-nil timestamp
	testTimestamp := timestamppb.New(time.Now())
	result := TimestampProtoToTimePtr(testTimestamp)
	if result == nil {
		t.Fatal("TimestampProtoToTimePtr returned nil for non-nil input")
	}
	if !result.Equal(testTimestamp.AsTime()) {
		t.Errorf("Expected %v, got %v", testTimestamp.AsTime(), *result)
	}

	// Test with a nil timestamp
	resultNil := TimestampProtoToTimePtr(nil)
	if resultNil != nil {
		t.Fatal("TimestampProtoToTimePtr should return nil for nil input")
	}
}

func TestIsEqualTimePtr(t *testing.T) {
	now := time.Now()

	// Both pointers are nil
	if !IsEqualTimePtr(nil, nil) {
		t.Error("IsEqualTimePtr should return true for two nil pointers")
	}

	// One pointer is nil, the other is not
	if IsEqualTimePtr(&now, nil) {
		t.Error("IsEqualTimePtr should return false when only one pointer is nil")
	}
	if IsEqualTimePtr(nil, &now) {
		t.Error("IsEqualTimePtr should return false when only one pointer is nil")
	}

	// Both pointers are non-nil and equal
	timeCopy := now
	if !IsEqualTimePtr(&now, &timeCopy) {
		t.Error("IsEqualTimePtr should return true for pointers to equal times")
	}

	// Both pointers are non-nil and not equal
	differentTime := now.Add(time.Hour)
	if IsEqualTimePtr(&now, &differentTime) {
		t.Error("IsEqualTimePtr should return false for pointers to different times")
	}
}

func TestBackOffExponentialDelay(t *testing.T) {
	tests := []struct {
		attempt   int
		wantDelay time.Duration
	}{
		{-1, 100 * time.Millisecond},
		{0, 100 * time.Millisecond},
		{1, 100 * time.Millisecond},
		{2, 200 * time.Millisecond},
		{3, 400 * time.Millisecond},
		{4, 800 * time.Millisecond},
		{5, 1600 * time.Millisecond},
		{6, 3200 * time.Millisecond},
		{7, 5 * time.Second},  // Cap at 5 seconds
		{10, 5 * time.Second}, // Cap at 5 seconds
	}

	for _, tc := range tests {
		got := BackOffExponentialDelay(tc.attempt)
		if got != tc.wantDelay {
			t.Errorf("BackOffExponentialDelay(%d) = %v; want %v", tc.attempt, got, tc.wantDelay)
		}
	}
}

func TestBackOffIncrementalDelay(t *testing.T) {
	tests := []struct {
		attempt   int
		wantDelay time.Duration
	}{
		{-1, 50 * time.Millisecond},
		{0, 50 * time.Millisecond},
		{1, 50 * time.Millisecond},
		{2, 100 * time.Millisecond},
		{3, 150 * time.Millisecond},
		{10, 500 * time.Millisecond},
		{20, 1000 * time.Millisecond},
		{40, 2 * time.Second}, // Cap at 2 seconds
		{50, 2 * time.Second}, // Cap at 2 seconds
	}

	for _, tc := range tests {
		got := BackOffIncrementalDelay(tc.attempt)
		if got != tc.wantDelay {
			t.Errorf("BackOffIncrementalDelay(%d) = %v; want %v", tc.attempt, got, tc.wantDelay)
		}
	}
}
