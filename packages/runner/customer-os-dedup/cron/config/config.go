package cron_config

type Config struct {
	// Default to hourly
	CronScheduleOrganizationDedup string `env:"CRON_SCHEDULE_ORGANIZATION_DEDUP" envDefault:"0 0 */1 * * *"`
}
