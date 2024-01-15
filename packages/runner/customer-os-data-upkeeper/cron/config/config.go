package cron_config

type Config struct {
	// Defaults to each 15 minutes
	CronScheduleUpdateContract string `env:"CRON_SCHEDULE_UPDATE_CONTRACT" envDefault:"0 */15 * * * *"`
	// Defaults to each 30 minutes
	CronScheduleWebScrapeOrganization string `env:"CRON_SCHEDULE_WEB_SCRAPE_ORGANIZATION" envDefault:"0 */30 * * * *"`
	// Defaults to each day at 8:00 AM
	CronScheduleGenerateInvoice string `env:"CRON_SCHEDULE_GENERATE_INVOICE" envDefault:"0 0 8 * * *"`
}
