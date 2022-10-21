package entity

import "time"

const (
	SessionColumnName_Country        string = "geo_country"
	SessionColumnName_City           string = "geo_city"
	SessionColumnName_RegionName     string = "geo_region_name"
	SessionColumnName_ReferrerSource string = "refr_source"
	SessionColumnName_UtmCampaign    string = "mkt_campaign"
	SessionColumnName_UtmContent     string = "mkt_content"
	SessionColumnName_UtmMedium      string = "mkt_medium"
	SessionColumnName_UtmSource      string = "mkt_source"
	SessionColumnName_UtmNetwork     string = "mkt_network"
	SessionColumnName_UtmTerm        string = "mkt_term"
	SessionColumnName_DeviceName     string = "device_name"
	SessionColumnName_DeviceBrand    string = "device_brand"
	SessionColumnName_DeviceClass    string = "device_class"
	SessionColumnName_AgentName      string = "agent_name"
	SessionColumnName_AgentVersion   string = "agent_version_major"
	SessionColumnName_OsFamily       string = "os_family"
	SessionColumnName_OsVersionMajor string = "os_major"
	SessionColumnName_OsVersionMinor string = "os_minor"
	SessionColumnName_FirstPagePath  string = "first_page_urlpath"
	SessionColumnName_LastPagePath   string = "last_page_urlpath"
)

type SessionEntity struct {
	ID             string    `gorm:"column:domain_sessionid;type:varchar(128);NOT NULL" json:"sessionId" binding:"required"`
	AppId          string    `gorm:"column:app_id;type:varchar(255);NOT NULL" json:"appName" binding:"required"`
	TrackerName    string    `gorm:"column:name_tracker;type:varchar(128);NOT NULL" json:"trackerName" binding:"required"`
	Tenant         string    `gorm:"column:tenant;type:varchar(64);NOT NULL" json:"tenant" binding:"required"`
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

type SessionEntities []SessionEntity

func (SessionEntity) TableName() string {
	return "derived.sessions"
}
