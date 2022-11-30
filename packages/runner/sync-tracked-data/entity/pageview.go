package entity

import (
	"database/sql"
	"time"
)

type PageView struct {
	ID                 string         `gorm:"column:page_view_id;primaryKey"`
	TrackerName        string         `gorm:"column:name_tracker;primaryKey"`
	AppId              string         `gorm:"column:app_id"`
	Tenant             string         `gorm:"column:tenant"`
	VisitorID          sql.NullString `gorm:"column:visitor_id"`
	SessionID          string         `gorm:"column:domain_sessionid"`
	OrderInSession     int            `gorm:"column:page_view_in_session_index"`
	EngagedTime        int            `gorm:"column:engaged_time_in_s"`
	Url                string         `gorm:"column:page_url"`
	Title              string         `gorm:"column:page_title"`
	Start              time.Time      `gorm:"column:start_tstamp"`
	End                time.Time      `gorm:"column:end_tstamp"`
	SyncedToCustomerOs bool           `gorm:"column:synced_to_customer_os"`
	ContactID          sql.NullString `gorm:"column:customer_os_contact_id"`
}

type PageViews []PageView

func (PageView) TableName() string {
	return schemaNameDerived + "." + tablePageViews
}
