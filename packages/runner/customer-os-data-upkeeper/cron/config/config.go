package cron_config

type Config struct {
	// Contracts
	// Defaults to each 15 minutes
	CronScheduleUpdateContract string `env:"CRON_SCHEDULE_UPDATE_CONTRACT" envDefault:"0 */15 * * * *"`

	// Organizations
	// Defaults to each 5 minutes
	CronScheduleRefreshLastTouchpoint string `env:"CRON_SCHEDULE_REFRESH_LAST_TOUCHPOINT" envDefault:"30 */1 * * * *"`
	// Defaults to each 15 minutes
	CronScheduleUpdateOrganization string `env:"CRON_SCHEDULE_UPDATE_ORGANIZATION" envDefault:"0 */15 * * * *"`

	// Contacts
	CronScheduleUpkeepContacts                            string `env:"CRON_SCHEDULE_UPKEEP_CONTACTS" envDefault:"0 */15 * * * *"`
	CronScheduleAskForWorkEmailOnBetterContact            string `env:"CRON_SCHEDULE_ASK_FOR_WORK_EMAIL_ON_BETTER_CONTACT" envDefault:"30 */2 * * * *"`
	CronScheduleEnrichWithWorkEmailFromBetterContact      string `env:"CRON_SCHEDULE_ENRICH_WITH_WORK_EMAIL_FROM_BETTER_CONTACT" envDefault:"0 */1 * * * *"`
	CronScheduleCheckBetterContactRequestsWithoutResponse string `env:"CRON_SCHEDULE_CHECK_BETTER_CONTACT_REQUESTS_WITHOUT_RESPONSE" envDefault:"40 */5 * * * *"`
	CronScheduleEnrichContacts                            string `env:"CRON_SCHEDULE_ENRICH_CONTACTS" envDefault:"0 */2 * * * *"`
	CronScheduleAskForLinkedInConnections                 string `env:"CRON_SCHEDULE_ASK_FOR_LINKEDIN_CONNECTIONS" envDefault:"*/5 * * * * *"`
	CronScheduleProcessLinkedInConnections                string `env:"CRON_SCHEDULE_PROCESS_LINKEDIN_CONNECTIONS" envDefault:"*/5 * * * * *"`

	// Invoices
	// Defaults to 8:15am
	CronScheduleGenerateInvoice string `env:"CRON_SCHEDULE_GENERATE_INVOICE" envDefault:"0 15 8 * * *"`
	// Defaults to each 9:15am and 3:15pm
	CronScheduleGenerateOffCycleInvoice string `env:"CRON_SCHEDULE_GENERATE_OFF_CYCLE_INVOICE" envDefault:"0 30 9,15 * * *"`
	// Defaults to each 10 min
	CronScheduleGenerateNextPreviewInvoice string `env:"CRON_SCHEDULE_GENERATE_NEXT_PREVIEW_INVOICE" envDefault:"30 */10 * * * *"`
	// Defaults to each 30 min
	CronScheduleSendPayInvoiceNotification string `env:"CRON_SCHEDULE_SEND_PAY_INVOICE_NOTIFICATION" envDefault:"0 */30 * * * *"`
	// Defaults to each 10 min
	CronScheduleGenerateInvoicePaymentLink string `env:"CRON_SCHEDULE_GENERATE_INVOICE_PAYMENT_LINK" envDefault:"30 */10 * * * *"`
	// Defaults to each 5 min
	CronScheduleCheckInvoiceFinalized string `env:"CRON_SCHEDULE_CHECK_INVOICE_FINALIZED" envDefault:"0 */5 * * * *"`
	// Defaults to each 8 hours
	CronScheduleCleanupInvoices string `env:"CRON_SCHEDULE_CLEANUP_INVOICES" envDefault:"0 40 */8 * * *"`
	// Defaults to each 30 mins
	CronScheduleAdjustInvoiceStatus string `env:"CRON_SCHEDULE_ADJUST_INVOICE_STATUS" envDefault:"0 */30 * * * *"`
	// Defaults to 9:15am
	CronScheduleSendRemindInvoiceNotification string `env:"CRON_SCHEDULE_SEND_REMIND_INVOICE_NOTIFICATION" envDefault:"0 15 9 * * *"`

	// Issues
	// Defaults to each 2 minutes
	CronScheduleLinkUnthreadIssues string `env:"CRON_LINK_UNTHREAD_ISSUES" envDefault:"0 */2 * * * *"`

	CronScheduleLinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn string `env:"CRON_SCHEDULE_LINK_ORPHAN_CONTACTS_TO_ORGANIZATION" envDefault:"0 */30 * * * *"`

	// Other
	// Defaults to each 15 minutes between 15-16 hours on working days
	CronScheduleGetCurrencyRatesECB string `env:"CRON_SCHEDULE_GET_CURRENCY_RATES_ECB" envDefault:"0 15 14-16 * * 1-5"`

	// Workflows
	// Defaults to each 5 minutes
	CronScheduleExecuteWorkflow string `env:"CRON_SCHEDULE_EXECUTE_WORKFLOW" envDefault:"0 */5 * * * *"`

	// Cache
	// Defaults to each hour at 15 minutes
	CronScheduleRefreshApiCache string `env:"CRON_SCHEDULE_REFRESH_API_CACHE" envDefault:"0 15 * * * *"`

	// Email
	// Defaults to each 2 minutes
	CronScheduleValidateCustomerOSEmails       string `env:"CRON_SCHEDULE_VALIDATE_EMAILS" envDefault:"0,30 * * * * *"`
	CronScheduleValidateEmailsFromBulkRequests string `env:"CRON_SCHEDULE_VALIDATE_EMAILS_BULK_REQUEST" envDefault:"*/10 * * * * *"`
	CronScheduleCheckScrubbyResult             string `env:"CRON_SCHEDULE_CHECK_SCRUBBY_RESULT" envDefault:"0 45 * * * *"`
	CronScheduleCheckEnrowResults              string `env:"CRON_SCHEDULE_CHECK_ENROW_RESULTS" envDefault:"0 */5 * * * *"`
	CronScheduleCleanEmails                    string `env:"CRON_SCHEDULE_CLEAN_EMAILS" envDefault:"0 0 */6 * * *"`

	// Flows
	CronScheduleRampUpMailboxes string `env:"CRON_SCHEDULE_RAMP_UP_MAILBOXES" envDefault:"* */1 * * * *"`
	CronScheduleFlowExecution   string `env:"CRON_SCHEDULE_FLOW_EXECUTION" envDefault:"*/5 * * * * *"`
	CronScheduleFlowStatistics  string `env:"CRON_SCHEDULE_FLOW_STATISTICS" envDefault:"*/5 * * * * *"`

	CronScheduleSendEmails        string `env:"CRON_SCHEDULE_SEND_EMAILS" envDefault:"*/5 * * * * *"`
	CronScheduleProcessSentEmails string `env:"CRON_SCHEDULE_PROCESS_SENT_EMAILS" envDefault:"*/5 * * * * *"`

	// Domains
	CronScheduleCheckDomains string `env:"CRON_SCHEDULE_CHECK_DOMAINS" envDefault:"0 */1 * * * *"`
}
