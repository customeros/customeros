package cron_config

type Config struct {
	// Heartbeat check, every minute
	CronScheduleHeartbeat           string `env:"CRON_SCHEDULE_HEARTBEAT" envDefault:"0 * * * * *"`
	CronScheduleProcessOutboxEvents string `env:"CRON_SCHEDULE_PROCESS_OUTBOX_EVENTS" envDefault:"0 */1 * * * *"`
	CronScheduleProcessWebSessions  string `env:"CRON_PROCESS_WEB_SESSIONS" envDefault:"0 */2 * * * *"`
}
