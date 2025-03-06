package utils

import (
	"testing"
	"time"
)

func TestDefaultBackoffConfig(t *testing.T) {
	config := DefaultBackoffConfig()

	if config.InitialDelay != 50*time.Millisecond {
		t.Errorf("InitialDelay = %v; want %v", config.InitialDelay, 50*time.Millisecond)
	}
	if config.MaxDelay != 5*time.Second {
		t.Errorf("MaxDelay = %v; want %v", config.MaxDelay, 5*time.Second)
	}
	if config.Factor != 2.0 {
		t.Errorf("Factor = %v; want %v", config.Factor, 2.0)
	}
	if config.Jitter != 0.1 {
		t.Errorf("Jitter = %v; want %v", config.Jitter, 0.1)
	}
}

func TestCalculateExponentialBackoffDelay(t *testing.T) {
	// Test with no jitter for deterministic results
	config := BackoffConfig{
		InitialDelay: 50 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Factor:       2.0,
		Jitter:       0, // No jitter for predictable tests
	}

	tests := []struct {
		name    string
		attempt int
		config  BackoffConfig
		wantMin time.Duration
		wantMax time.Duration
	}{
		{
			name:    "negative attempt",
			attempt: -1,
			config:  config,
			wantMin: 50 * time.Millisecond,
			wantMax: 50 * time.Millisecond,
		},
		{
			name:    "zero attempt",
			attempt: 0,
			config:  config,
			wantMin: 50 * time.Millisecond,
			wantMax: 50 * time.Millisecond,
		},
		{
			name:    "first attempt",
			attempt: 1,
			config:  config,
			wantMin: 50 * time.Millisecond,
			wantMax: 50 * time.Millisecond,
		},
		{
			name:    "second attempt",
			attempt: 2,
			config:  config,
			wantMin: 100 * time.Millisecond,
			wantMax: 100 * time.Millisecond,
		},
		{
			name:    "max delay cap",
			attempt: 10,
			config:  config,
			wantMin: 5 * time.Second,
			wantMax: 5 * time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateExponentialBackoffDelay(tc.attempt, tc.config)
			if got < tc.wantMin || got > tc.wantMax {
				t.Errorf("CalculateExponentialBackoffDelay(%d, config) = %v; want between %v and %v",
					tc.attempt, got, tc.wantMin, tc.wantMax)
			}
		})
	}
}

func TestBackOffExponentialDelay(t *testing.T) {
	tests := []struct {
		name    string
		attempt int
		wantMin time.Duration
		wantMax time.Duration
	}{
		{
			name:    "negative attempt",
			attempt: -1,
			wantMin: 50 * time.Millisecond,
			wantMax: 55 * time.Millisecond, // With 10% jitter
		},
		{
			name:    "first attempt",
			attempt: 1,
			wantMin: 50 * time.Millisecond,
			wantMax: 55 * time.Millisecond,
		},
		{
			name:    "second attempt",
			attempt: 2,
			wantMin: 100 * time.Millisecond,
			wantMax: 110 * time.Millisecond,
		},
		{
			name:    "max delay cap",
			attempt: 10,
			wantMin: 5 * time.Second,
			wantMax: 5 * time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BackOffExponentialDelay(tc.attempt)
			if got < tc.wantMin || got > tc.wantMax {
				t.Errorf("BackOffExponentialDelay(%d) = %v; want between %v and %v",
					tc.attempt, got, tc.wantMin, tc.wantMax)
			}
		})
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
