package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func TruncateFloat64(input float64, decimals int) float64 {
	if input == 0 {
		return 0
	}
	multiplier := math.Pow(10, float64(decimals))
	truncated := math.Trunc(input*multiplier) / multiplier
	return truncated
}

func RoundHalfUpFloat64(input float64, decimals int) float64 {
	if input == 0 {
		return 0
	}
	multiplier := math.Pow(10, float64(decimals))
	rounded := math.Round(input*multiplier) / multiplier
	return rounded
}

func Float64PtrEquals(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b != nil {
		return *a == *b
	}
	return false
}

func ParseStringToFloat(input string) *float64 {
	if input == "" {
		return nil
	}

	parsedFloat, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Error parsing string to float: %v\n", err)
		return nil
	}
	return &parsedFloat
}

func FloatToString(num *float64) string {
	if num == nil {
		return ""
	}
	return fmt.Sprintf("%f", *num)
}

func FormatAmount(amount float64, decimals int) string {
	// Split the amount into its whole and fractional parts.
	truncatedAmount := TruncateFloat64(amount, decimals)

	// Convert the whole part to an integer.
	wholePart := int64(truncatedAmount)

	// Format the whole part with commas.
	wholeStr := fmt.Sprintf("%d", wholePart)
	if len(wholeStr) > 3 {
		for i := len(wholeStr) - 3; i > 0; i -= 3 {
			wholeStr = wholeStr[:i] + "," + wholeStr[i:]
		}
	}

	if decimals == 0 {
		return wholeStr
	}
	// Format the fractional part by converting the truncated amount to a string
	// and then extracting the fractional part.
	amountStr := fmt.Sprintf("%.*f", decimals, truncatedAmount)
	amountStr = strings.TrimRight(amountStr, "0")
	amountStr = strings.TrimRight(amountStr, ".")
	split := strings.Split(amountStr, ".")
	fractionalStr := ""
	if len(split) > 1 {
		fractionalStr = "." + split[1]
	}

	// Concatenate the whole and fractional parts.
	return wholeStr + fractionalStr
}
