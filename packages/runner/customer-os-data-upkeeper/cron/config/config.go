package cron_config

type Config struct {
	// Defaults to first minute after midnight & midday every day
	CronScheduleUpdateOrgNextCycleDate string `env:"CRON_SCHEDULE_UPDATE_ORGANIZATION_NEXT_CYCLE_DATE" envDefault:"0 1 0,12 * * *"`
	// Defaults to each 15 minutes
	CronScheduleUpdateContract string `env:"CRON_SCHEDULE_UPDATE_CONTRACT" envDefault:"0 */15 * * * *"`
}
