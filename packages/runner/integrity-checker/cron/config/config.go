package cron_config

type Config struct {
	CronScheduleNeo4jIntegrityChecker    string `env:"CRON_SCHEDULE_NEO4J_INTEGRITY_CHECKER" envDefault:"0 0 */2 * * *"`
	CronSchedulePostgresIntegrityChecker string `env:"CRON_SCHEDULE_POSTGRES_INTEGRITY_CHECKER" envDefault:"0 0 */2 * * *"`
}
