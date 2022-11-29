package entity

import "time"

type Session struct {
	ID             string    `gorm:"column:domain_sessionid"`
	AppId          string    `gorm:"column:app_id"`
	TrackerName    string    `gorm:"column:name_tracker"`
	Tenant         string    `gorm:"column:tenant"`
	Country        string    `gorm:"column:geo_country"`
	Region         string    `gorm:"column:geo_region_name"`
	City           string    `gorm:"column:geo_city"`
	ReferrerSource string    `gorm:"column:refr_source"`
	UtmCampaign    string    `gorm:"column:mkt_campaign"`
	UtmContent     string    `gorm:"column:mkt_content"`
	UtmMedium      string    `gorm:"column:mkt_medium"`
	UtmSource      string    `gorm:"column:mkt_source"`
	UtmNetwork     string    `gorm:"column:mkt_network"`
	UtmTerm        string    `gorm:"column:mkt_term"`
	DeviceBrand    string    `gorm:"column:device_brand"`
	DeviceName     string    `gorm:"column:device_name"`
	DeviceClass    string    `gorm:"column:device_class"`
	AgentName      string    `gorm:"column:agent_name"`
	AgentVersion   string    `gorm:"column:agent_version_major"`
	OsFamily       string    `gorm:"column:os_family"`
	OsVersionMajor string    `gorm:"column:os_major"`
	OsVersionMinor string    `gorm:"column:os_minor"`
	FirstPagePath  string    `gorm:"column:first_page_urlpath"`
	LastPagePath   string    `gorm:"column:last_page_urlpath"`
	Start          time.Time `gorm:"column:start_tstamp"`
	End            time.Time `gorm:"column:end_tstamp"`
	EngagedTime    int       `gorm:"column:engaged_time_in_s"`
}

type Sessions []Session

func (Session) TableName() string {
	return schemaNameDerived + "." + tableSessions
}
