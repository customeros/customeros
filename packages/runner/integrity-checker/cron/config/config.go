package cron_config

type Config struct {
	// Default to every 4 hours
	CronScheduleNeo4jIntegrityChecker string `env:"CRON_SCHEDULE_NEO4J_INTEGRITY_CHECKER" envDefault:"0 0 */4 * * *"`
}
