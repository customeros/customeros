package resolver

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_DashboardHelper_GetPeriod(t *testing.T) {
	testCases := [][]any{
		{"2020-01-01", "2019-02-28", "2020-01-01"},
		{"2020-01-15", "2019-02-28", "2020-01-15"},
		{"2020-01-31", "2019-02-28", "2020-01-31"},
		{"2020-02-01", "2019-03-31", "2020-02-01"},
		{"2020-02-15", "2019-03-31", "2020-02-15"},
		{"2020-02-29", "2019-03-31", "2020-02-29"},
		{"2020-03-01", "2019-04-30", "2020-03-01"},

		{"2021-01-01", "2020-02-29", "2021-01-01"},
		{"2021-01-15", "2020-02-29", "2021-01-15"},
		{"2021-01-28", "2020-02-29", "2021-01-28"},
		{"2021-01-29", "2020-02-29", "2021-01-29"},
		{"2021-01-30", "2020-02-29", "2021-01-30"},
		{"2021-01-31", "2020-02-29", "2021-01-31"},
		{"2021-02-01", "2020-03-31", "2021-02-01"},
		{"2021-02-15", "2020-03-31", "2021-02-15"},
		{"2021-02-28", "2020-03-31", "2021-02-28"},
		{"2021-03-01", "2020-04-30", "2021-03-01"},
	}

	for _, testCase := range testCases {
		nowString := testCase[0].(string)
		now, err := time.Parse(time.DateOnly, nowString)
		require.NoError(t, err)

		start, stop, err := getPeriod(nil, now)
		require.NoError(t, err)

		startString := start.Format(time.DateOnly)
		stopString := stop.Format(time.DateOnly)

		fmt.Printf("now: %s, start: %s, stop: %s\n", nowString, startString, stopString)
		require.Equal(t, testCase[1], startString)
		require.Equal(t, testCase[2], stopString)
	}
}
