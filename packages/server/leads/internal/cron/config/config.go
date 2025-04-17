package cron_config

type Config struct {
	// Heartbeat check, every minute
	CronScheduleHeartbeat          string `env:"CRON_SCHEDULE_HEARTBEAT" envDefault:"0 * * * * *"`
	CronScheduleProcessWebSessions string `env:"CRON_PROCESS_WEB_SESSIONS" envDefault:"0 */2 * * * *"`
}
