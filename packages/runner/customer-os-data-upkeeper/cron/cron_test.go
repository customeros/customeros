package cron

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/container"
	cron_config "github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/cron/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getLogger() logger.Logger {
	appLogger := logger.NewAppLogger(&logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()
	return appLogger
}

func TestStartCron(t *testing.T) {
	// Arrange
	cfg := config.Config{
		Cron: cron_config.Config{CronScheduleUpdateOrgNextCycleDate: "0 0 */1 * * *"},
	}

	// Act
	cron := StartCron(&container.Container{
		Cfg: &cfg,
		Log: getLogger(),
	})

	// Assert
	assert.NotNil(t, cron)
	assert.Equal(t, getNextHourStartTime(), cron.Entries()[0].Schedule.Next(time.Now()))
}

func getNextHourStartTime() time.Time {
	now := time.Now()

	year, month, day := now.Date()
	hour, _, _ := now.Clock()

	nextHour := time.Date(year, month, day, hour+1, 0, 0, 0, now.Location())

	return nextHour
}

func TestStopCron(t *testing.T) {
	// Arrange
	c := cron.New()
	c.Start()

	// Act
	err := StopCron(getLogger(), c)
	if err != nil {
		t.Fatalf("Error stopping cron: %v", err.Error())
	}

	// Assert
	assert.Equal(t, 0, len(c.Entries()))
}
