package test

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func AssertTimeRecentlyChanged(t *testing.T, checkTime time.Time) {
	// Set the time difference to 5 seconds
	X := 5

	currentTime := utils.Now()

	// Calculate the time difference
	diff := currentTime.Sub(checkTime)

	// Use the require package to assert the time difference
	require.True(t, diff <= time.Duration(X)*time.Second, "The time is within the last %d seconds.", X)
}
