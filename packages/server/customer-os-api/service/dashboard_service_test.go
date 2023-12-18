package service

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_DashboardService_ComputeNewCustomersDisplay(t *testing.T) {
	testCases := [][]float64{
		{1, 2},
		{0, 1},
		{0, 2},
		{0, 10},
		{10, 0},
		{100, 0},
		{1000, 0},
		{10, 12},
		{100, 12},
		{1000, 12},
		{10000, 12},
		{10, 75},
		{3, 200},
		{10, 30},
		{10, 200},
		{6, 2},
		{2, 6},
		{0, 208},
		{0, 0},
		{150, 92},
		{100, 2},
		{4, 8},
		{8, 2},
		{30, 4},
		{950, 45},
		{12000, 10},
		{1200, 52},
		{1200, 12},
		{30, 0},
		{2, 0},
	}

	expected := []string{
		"+100%",
		"+100%",
		"+2",
		"+10",
		"-100%",
		"-100%",
		"-100%",
		"+20%",
		"-88%",
		"-99%",
		"-100%",
		"6.5×",
		"65.7×",
		"2×",
		"19×",
		"-67%",
		"2×",
		"+208",
		"0%",
		"-39%",
		"-98%",
		"+100%",
		"-75%",
		"-87%",
		"-95%",
		"-100%",
		"-96%",
		"-99%",
		"-100%",
		"-100%",
	}

	for index, testCase := range testCases {
		previousCount, currentCount := testCase[0], testCase[1]
		displayValue := ComputeNewCustomersDisplay(previousCount, currentCount)
		fmt.Printf("Previous: %f, Current: %f, Display: %s\n", previousCount, currentCount, displayValue)
		require.Equal(t, expected[index], displayValue)
	}
}
