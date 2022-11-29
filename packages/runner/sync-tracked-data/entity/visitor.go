package entity

import "database/sql"

type Visitor struct {
	VisitorId          sql.NullString `gorm:"column:visitor_id"`
	AppId              string         `gorm:"column:app_id;primaryKey"`
	TrackerName        string         `gorm:"column:name_tracker;primaryKey"`
	Tenant             string         `gorm:"column:tenant"`
	DomainUserId       string         `gorm:"column:domain_userid;primaryKey"`
	NetworkUserId      string         `gorm:"column:network_userid"`
	NumPageViews       int            `gorm:"column:page_views"`
	NumSessions        int            `gorm:"column:sessions"`
	EngagedTime        int            `gorm:"column:engaged_time_in_s"`
	SyncedToCustomerOs bool           `gorm:"column:synced_to_customer_os"`
}

type Visitors []Visitor

func (Visitor) TableName() string {
	return schemaNameDerived + "." + tableVisitors
}
