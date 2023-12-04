package cron_config

type Config struct {
	// Defaults to each 15 minutes
	CronScheduleUpdateContract string `env:"CRON_SCHEDULE_UPDATE_CONTRACT" envDefault:"0 */15 * * * *"`
}
