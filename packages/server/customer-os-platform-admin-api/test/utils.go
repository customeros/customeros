package test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func AssertRecentTime(t *testing.T, checkTime time.Time) {
	x := 5 // Set the time difference to 5 seconds

	diff := time.Since(checkTime)

	require.True(t, diff <= time.Duration(x)*time.Second, "The time is within the last %d seconds.", x)
}
