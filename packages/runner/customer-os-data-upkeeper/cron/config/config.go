package cron_config

type Config struct {
	// Contracts
	// Defaults to each 15 minutes
	CronScheduleUpdateContract string `env:"CRON_SCHEDULE_UPDATE_CONTRACT" envDefault:"0 */15 * * * *"`

	// Organizations
	// Defaults to each 6 hours
	CronScheduleRefreshLastTouchpoint string `env:"CRON_SCHEDULE_REFRESH_LAST_TOUCHPOINT" envDefault:"0 35 */6 * * *"`
	// Defaults to each 15 minutes
	CronScheduleUpdateOrganization string `env:"CRON_SCHEDULE_UPDATE_ORGANIZATION" envDefault:"0 */15 * * * *"`

	// Defaults to each 15 minutes
	CronScheduleUpkeepContacts          string `env:"CRON_SCHEDULE_UPKEEP_CONTACTS" envDefault:"0 */15 * * * *"`
	CronScheduleEnrichContactsFindEmail string `env:"CRON_SCHEDULE_ENRICH_CONTACTS_FIND_EMAIL" envDefault:"0 */5 * * * *"`
	CronScheduleEnrichContacts          string `env:"CRON_SCHEDULE_ENRICH_CONTACTS" envDefault:"0 */2 * * * *"`

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
	// Defaults to each 8 hours
	CronScheduleCleanupInvoices string `env:"CRON_SCHEDULE_CLEANUP_INVOICES" envDefault:"0 40 */8 * * *"`
	// Defaults to each 30 mins
	CronScheduleAdjustInvoiceStatus string `env:"CRON_SCHEDULE_ADJUST_INVOICE_STATUS" envDefault:"0 */30 * * * *"`

	// Issues
	// Defaults to each 2 minutes
	CronScheduleLinkUnthreadIssues string `env:"CRON_LINK_UNTHREAD_ISSUES" envDefault:"0 */2 * * * *"`

	// Other
	// Defaults to each 15 minutes between 15-16 hours on working days
	CronScheduleGetCurrencyRatesECB string `env:"CRON_SCHEDULE_GET_CURRENCY_RATES_ECB" envDefault:"0 15 14-16 * * 1-5"`

	// Defaults to each 1 hour
	CronScheduleRefreshApiCache string `env:"CRON_SCHEDULE_REFRESH_API_CACHE" envDefault:"* * */1 * * *"`
}
