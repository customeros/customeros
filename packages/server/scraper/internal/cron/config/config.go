package cron_config

type Config struct {
	// Heartbeat check, every minute
	CronScheduleHeartbeat string `env:"CRON_SCHEDULE_HEARTBEAT" envDefault:"0 * * * * *"`
}
