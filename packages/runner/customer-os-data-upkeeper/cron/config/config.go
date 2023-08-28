package cron_config

type Config struct {
	// Defaults to first minute after midnight & midday every day
	CronScheduleUpdateOrgNextCycleDate string `env:"CRON_SCHEDULE_UPDATE_ORGANIZATION_NEXT_CYCLE_DATE" envDefault:"0 1 0,12 * * *"`
}
