package cron_config

type Config struct {
	// Defaults to each 15 minutes
	CronScheduleUpdateContract string `env:"CRON_SCHEDULE_UPDATE_CONTRACT" envDefault:"0 */15 * * * *"`
	// Defaults to each 30 minutes
	CronScheduleWebScrapeOrganization string `env:"CRON_SCHEDULE_WEB_SCRAPE_ORGANIZATION" envDefault:"0 */30 * * * *"`
	// Defaults to 8:15am
	CronScheduleGenerateInvoice string `env:"CRON_SCHEDULE_GENERATE_INVOICE" envDefault:"0 15 8 * * *"`
	// Defaults to each 9:15am and 3:15pm
	CronScheduleGenerateOffCycleInvoice string `env:"CRON_SCHEDULE_GENERATE_OFF_CYCLE_INVOICE" envDefault:"0 30 9,15 * * *"`
	// Defaults to each 30 min
	CronScheduleSendPayInvoiceNotification string `env:"CRON_SCHEDULE_SEND_PAY_INVOICE_NOTIFICATION" envDefault:"0 */30 * * * *"`
	// Defaults to each 6 hours
	CronScheduleRefreshLastTouchpoint string `env:"CRON_SCHEDULE_REFRESH_LAST_TOUCHPOINT" envDefault:"0 35 */6 * * *"`
	// Defaults to each 15 minutes between 15-16 hours on working days
	CronScheduleGetCurrencyRatesECB string `env:"CRON_SCHEDULE_GET_CURRENCY_RATES_ECB" envDefault:"0 15 14-16 * * 1-5"`
	// Defaults to each 2 minutes
	CronScheduleLinkUnthreadIssues string `env:"CRON_LINK_UNTHREAD_ISSUES" envDefault:"0 */2 * * * *"`
	// Defaults to each 10 min
	CronScheduleGenerateInvoicePaymentLink string `env:"CRON_SCHEDULE_GENERATE_INVOICE_PAYMENT_LINK" envDefault:"30 */10 * * * *"`
}
