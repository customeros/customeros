package utils

import "testing"

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		decimals int
		want     string
	}{
		{
			name:     "Simple Test",
			amount:   345.01,
			decimals: 2,
			want:     "345.01",
		},
		{
			name:     "Large Number",
			amount:   23567.89,
			decimals: 2,
			want:     "23,567.89",
		},
		{
			name:     "No Decimals",
			amount:   1000,
			decimals: 0,
			want:     "1,000",
		},
		{
			name:     "Negative Number",
			amount:   -1234.56,
			decimals: 2,
			want:     "-1,234.56",
		},
		{
			name:     "Zero",
			amount:   0,
			decimals: 2,
			want:     "0",
		},
		{
			name:     "Round Down",
			amount:   99.999,
			decimals: 2,
			want:     "99.99",
		},
		{
			name:     "Very Large Number",
			amount:   123456789.12345,
			decimals: 3,
			want:     "123,456,789.123",
		},
		{
			name:     "No fraction 1",
			amount:   123.00001,
			decimals: 3,
			want:     "123",
		},
		{
			name:     "No fraction 2",
			amount:   123,
			decimals: 10,
			want:     "123",
		},
		{
			name:     "No fraction 3",
			amount:   0,
			decimals: 2,
			want:     "0",
		},
		{
			name:     "No fraction 4",
			amount:   1234.1234,
			decimals: 0,
			want:     "1,234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatAmount(tt.amount, tt.decimals)
			if got != tt.want {
				t.Errorf("FormatAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundHalfUpFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		decimals int
		want     float64
	}{
		{"Zero value", 0, 0, 0},
		{"Zero value", 0, 1, 0},
		{"Zero value", 0, 2, 0},

		{"Whole number - no decimals", 123, 0, 123},
		{"Whole number - one decimal", 123, 1, 123},
		{"Whole number - two decimals", 123, 2, 123},

		{"Four rounding - no decimal", 123.4, 0, 123},
		{"Four rounding - one decimal", 123.44, 1, 123.4},
		{"Four rounding - two decimals", 123.444, 2, 123.44},

		{"Five rounding - no decimal", 123.5, 0, 124},
		{"Five rounding - one decimal", 123.45, 1, 123.5},
		{"Five rounding - two decimals", 123.455, 2, 123.46},

		{"Nine rounding - no decimal", 123.9, 0, 124},
		{"Nine rounding - one decimal", 123.99, 1, 124.0},
		{"Nine rounding - two decimals", 123.999, 2, 124.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundHalfUpFloat64(tt.input, tt.decimals); got != tt.want {
				t.Errorf("RoundHalfUpFloat64(%v, %d) = %v, want %v", tt.input, tt.decimals, got, tt.want)
			}
		})
	}
}
