package cron

import (
	"os"
	"testing"

	"github.com/customeros/customeros/packages/server/core-crm/internal/config"
	cron_config "github.com/customeros/customeros/packages/server/core-crm/internal/cron/config"
	"github.com/customeros/customeros/packages/server/core-crm/internal/logger"
	cronv3 "github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes"
)

type mockKubernetesInterface struct {
	kubernetes.Interface
	mock.Mock
}

func getLogger() logger.Logger {
	appLogger := logger.NewAppLogger(&logger.Config{
		DevMode: true,
	})
	appLogger.InitLogger()
	return appLogger
}

func TestNewCronManager(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		AppConfig: &config.AppConfig{},
	}
	log := getLogger()
	k8s := &mockKubernetesInterface{}

	// Act
	cm := NewCronManager(cfg, log, k8s, nil)

	// Assert
	assert.NotNil(t, cm)
	assert.Equal(t, cfg, cm.cfg)
	assert.Equal(t, log, cm.log)
	assert.Equal(t, k8s, cm.k8s)
	assert.NotNil(t, cm.jobIDs)
}

func TestCronManager_StartCron(t *testing.T) {
	// Set environment variable for testing
	os.Setenv("CRON_SCHEDULE_OUTBOX", "0 * * * * *")
	defer os.Unsetenv("CRON_SCHEDULE_OUTBOX")

	// Arrange
	cfg := &config.Config{
		AppConfig: &config.AppConfig{},
	}
	log := getLogger()
	k8s := &mockKubernetesInterface{}
	cm := NewCronManager(cfg, log, k8s, nil)

	// Create a mock cron for testing
	mockCron := cronv3.New()

	// Register jobs directly
	var cronConfig cron_config.Config
	cronConfig.CronScheduleProcessOutboxEvents = "0 * * * * *"

	// Act - register jobs manually
	id, err := mockCron.AddFunc(cronConfig.CronScheduleProcessOutboxEvents, func() {})
	assert.NoError(t, err)
	cm.jobIDs["outbox"] = id

	cm.cron = mockCron

	// Assert

	cm.cron = mockCron

	// Assert
	assert.NotNil(t, cm.cron)
	assert.Equal(t, 3, len(cm.jobIDs))
}

func TestCronManager_Stop(t *testing.T) {
	// Set environment variable for testing
	os.Setenv("CRON_SCHEDULE_OUTBOX", "0 * * * * *")
	defer os.Unsetenv("CRON_SCHEDULE_OUTBOX")

	// Arrange
	cfg := &config.Config{
		AppConfig: &config.AppConfig{},
	}
	log := getLogger()
	k8s := &mockKubernetesInterface{}
	cm := NewCronManager(cfg, log, k8s, nil)

	// Create a mock cron for testing
	mockCron := cronv3.New()
	mockCron.Start()
	cm.cron = mockCron

	// Act
	cm.Stop()

	// Assert
	select {
	case <-cm.stopCh:
		// Channel is closed as expected
	default:
		t.Error("Stop channel was not closed")
	}
}
