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
		{3, 200, "65.7×"},
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
