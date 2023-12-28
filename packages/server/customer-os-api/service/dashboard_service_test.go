package service

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_DashboardService_ComputeNumbersDisplay(t *testing.T) {
	testCases := [][]any{
		{1, 2, "+100%"},
		{0, 1, "+100%"},
		{0, 2, "+2"},
		{0, 10, "+10"},
		{10, 0, "-100%"},
		{100, 0, "-100%"},
		{1000, 0, "-100%"},
		{10, 12, "+20%"},
		{100, 12, "-88%"},
		{1000, 12, "-99%"},
		{10000, 12, "-100%"},
		{10, 75, "6.5×"},
		{3, 200, "65.67×"},
		{10, 30, "2×"},
		{10, 200, "19×"},
		{6, 2, "-67%"},
		{2, 6, "2×"},
		{0, 208, "+208"},
		{0, 0, "0%"},
		{150, 92, "-39%"},
		{100, 2, "-98%"},
		{4, 8, "+100%"},
		{8, 2, "-75%"},
		{30, 4, "-87%"},
		{950, 45, "-95%"},
		{12000, 10, "-100%"},
		{1200, 52, "-96%"},
		{1200, 12, "-99%"},
		{30, 0, "-100%"},
		{2, 0, "-100%"},
	}

	for _, testCase := range testCases {
		previousCount, currentCount := float64(testCase[0].(int)), float64(testCase[1].(int))
		displayValue := ComputeNumbersDisplay(previousCount, currentCount)
		fmt.Printf("Previous: %f, Current: %f, Display: %s\n", previousCount, currentCount, displayValue)
		require.Equal(t, testCase[2], displayValue)
	}
}
func Test_DashboardService_ComputePercentagesDisplay(t *testing.T) {
	testCases := [][]any{
		{0, 0, "0"},
		{0, 50, "+50"},
		{0, 100, "+100"},
		{1, 25, "+24"},
		{10, 50, "+40"},
		{25, 100, "+75"},
		{50, 100, "+50"},
		{84, 86, "+2"},
		{86, 84, "-2"},
		{100, 75, "-25"},

		{-100, -100, "0"},
		{-100, 100, "+100"},
		{-100, 50, "+100"},
		{-100, 0, "+100"},
		{-50, 0, "+50"},
		{-25, 0, "+25"},
		{-100, -50, "+50"},

		{-1000, -100, "+100"},
		{-1000, -50, "+100"},
		{-1000, 0, "+100"},
		{-1000, 100, "+100"},

		{200, -100, "-100"},
		{200, -50, "-100"},
		{200, 0, "-100"},
		{200, 100, "-100"},
		{200, 150, "-50"},
		{200, 200, "0"},
		{200, 250, "+50"},
		{200, 300, "+100"},
		{200, 400, "+100"},
	}

	for _, testCase := range testCases {
		previousCount, currentCount := float64(testCase[0].(int)), float64(testCase[1].(int))
		displayValue := ComputePercentagesDisplay(previousCount, currentCount)
		fmt.Printf("Previous: %f, Current: %f, Display: %s\n", previousCount, currentCount, displayValue)
		require.Equal(t, testCase[2], displayValue)
	}
}

func Test_DashboardService_PrintFloatValue(t *testing.T) {
	testCases := [][]any{
		{0, "0"},
		{0.0, "0"},
		{0.1, "0.1"},
		{0.11, "0.11"},
		{0.111, "0.11"},
		{0.5, "0.5"},
		{1, "1"},
		{1.0, "1"},
		{1.1, "1.1"},
		{1.11, "1.11"},
		{1.111, "1.11"},
		{1.4999999, "1.5"},
		{1.5, "1.5"},
		{10, "10"},
		{10.0, "10"},
		{10.1, "10.1"},
		{10.11, "10.11"},
		{10.111, "10.11"},
		{10.114999999, "10.11"},
		{10.115, "10.12"},
		{10.5, "10.5"},
		{100, "100"},
		{100.0, "100"},
		{100.1, "100"},
		{100.11, "100"},
		{100.111, "100"},
		{100.4999999, "100"},
		{100.5, "101"},
		{1000, "1000"},
		{1000.0, "1000"},
		{1000.1, "1000"},
		{1000.11, "1000"},
		{1000.111, "1000"},
		{1000.4999999, "1000"},
		{1000.5, "1001"},
	}

	for _, testCase := range testCases {
		val := getCorrectValueType(testCase[0])
		displayValue := PrintFloatValue(val, false)
		fmt.Printf("Val: %f, Display: %s\n", val, displayValue)
		require.Equal(t, testCase[1], displayValue)
	}
}

func getCorrectValueType(valueToExtract any) float64 {
	var v float64

	switch val := valueToExtract.(type) {
	case int:
		v = float64(val)
	case int64:
		v = float64(val)
	case float64:
		v = val
	default:
		fmt.Errorf("unexpected type %T", val)
	}

	return v
}
