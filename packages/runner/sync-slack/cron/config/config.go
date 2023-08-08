package cron_config

type Config struct {
	// default to 1 min
	CronScheduleSyncFromSlack string `env:"CRON_SYNC_FROM_SLACK" envDefault:"0 */1 * * * *"`
}
