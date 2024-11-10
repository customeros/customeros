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
		Cron: cron_config.Config{
			CronScheduleUpdateContract:                                        "0 0 */1 * * *",
			CronScheduleUpdateOrganization:                                    "0 0 */1 * * *",
			CronScheduleGenerateInvoice:                                       "0 0 */1 * * *",
			CronScheduleGenerateOffCycleInvoice:                               "0 0 */1 * * *",
			CronScheduleGenerateNextPreviewInvoice:                            "0 0 */1 * * *",
			CronScheduleSendPayInvoiceNotification:                            "0 0 */1 * * *",
			CronScheduleSendRemindInvoiceNotification:                         "0 0 */1 * * *",
			CronScheduleRefreshLastTouchpoint:                                 "0 0 */1 * * *",
			CronScheduleGetCurrencyRatesECB:                                   "0 0 */1 * * *",
			CronScheduleLinkUnthreadIssues:                                    "0 0 */1 * * *",
			CronScheduleGenerateInvoicePaymentLink:                            "0 0 */1 * * *",
			CronScheduleCheckInvoiceFinalized:                                 "0 0 */1 * * *",
			CronScheduleCleanupInvoices:                                       "0 0 */1 * * *",
			CronScheduleAdjustInvoiceStatus:                                   "0 0 */1 * * *",
			CronScheduleUpkeepContacts:                                        "0 0 */1 * * *",
			CronScheduleAskForWorkEmailOnBetterContact:                        "0 0 */1 * * *",
			CronScheduleEnrichWithWorkEmailFromBetterContact:                  "0 0 */1 * * *",
			CronScheduleCheckBetterContactRequestsWithoutResponse:             "0 0 */1 * * *",
			CronScheduleEnrichContacts:                                        "0 0 */1 * * *",
			CronScheduleRefreshApiCache:                                       "0 0 */1 * * *",
			CronScheduleExecuteWorkflow:                                       "0 0 */1 * * *",
			CronScheduleAskForLinkedInConnections:                             "0 0 */1 * * *",
			CronScheduleProcessLinkedInConnections:                            "0 0 */1 * * *",
			CronScheduleLinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn: "0 0 */1 * * *",
			CronScheduleValidateEmails:                                        "0 0 */1 * * *",
			CronScheduleValidateEmailsFromBulkRequests:                        "0 0 */1 * * *",
			CronScheduleCheckScrubbyResult:                                    "0 0 */1 * * *",
			CronScheduleCheckEnrowResults:                                     "0 0 */1 * * *",
			CronScheduleRampUpMailboxes:                                       "0 0 */1 * * *",
			CronScheduleFlowExecution:                                         "0 0 */1 * * *",
			CronScheduleFlowStatistics:                                        "0 0 */1 * * *",
			CronScheduleCleanEmails:                                           "0 0 */1 * * *",
			CronScheduleSendEmails:                                            "0 0 */1 * * *",
			CronScheduleProcessSentEmails:                                     "0 0 */1 * * *",
			CronScheduleCheckDomains:                                          "0 0 */1 * * *",
		},
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
