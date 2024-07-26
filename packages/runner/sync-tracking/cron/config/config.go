package config

type Config struct {
	CronScheduleProcessNewRecords                         string `env:"CRON_SCHEDULE_PROCESS_NEW_RECORDS" envDefault:"*/10 * * * * *"`
	CronScheduleProcessIPDataRequests                     string `env:"CRON_SCHEDULE_PROCESS_IP_DATA_REQUESTS" envDefault:"*/10 * * * * *"`
	CronScheduleProcessIPDataResponses                    string `env:"CRON_SCHEDULE_PROCESS_IP_DATA_RESPONSES" envDefault:"*/10 * * * * *"`
	CronScheduleIdentifyTrackingRecords                   string `env:"CRON_SCHEDULE_IDENTIFY_TRACKING_RECORDS" envDefault:"*/10 * * * * *"`
	CronScheduleCreateOrganizationsFromTrackedDataRecords string `env:"CRON_SCHEDULE_CREATE_ORGANIZATIONS_FROM_TRACKED_DATA" envDefault:"*/10 * * * * *"`
	CronScheduleNotifyOnSlack                             string `env:"CRON_SCHEDULE_NOTIFY_ON_SLACK" envDefault:"*/10 * * * * *"`
}
