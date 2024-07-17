package config

type Config struct {
	CronScheduleShouldIdentifyTrackingRecords string `env:"CRON_SCHEDULE_SHOULD_IDENTIFY_TRACKING_RECORDS" envDefault:"*/10 * * * * *"`
	CronScheduleIdentifyTrackingRecords       string `env:"CRON_SCHEDULE_IDENTIFY_TRACKING_RECORDS" envDefault:"*/10 * * * * *"`
}
