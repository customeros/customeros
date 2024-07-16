package config

type Config struct {
	CronScheduleIdentifyTrackingRecords string `env:"CRON_SCHEDULE_IDENTIFY_TRACKING_RECORDS" envDefault:"*/10 * * * * *"`
}
