package entity

import (
	"database/sql"
	"time"
)

type PageView struct {
	ID             string         `gorm:"column:page_view_id;primaryKey"`
	AppId          string         `gorm:"column:app_id"`
	TrackerName    string         `gorm:"column:name_tracker;primaryKey"`
	Tenant         string         `gorm:"column:tenant"`
	VisitorID      sql.NullString `gorm:"column:visitor_id"`
	SessionID      string         `gorm:"column:domain_sessionid"`
	OrderInSession int            `gorm:"column:page_view_in_session_index"`
	EngagedTime    int            `gorm:"column:engaged_time_in_s"`
	Path           string         `gorm:"column:page_urlpath"`
	Title          string         `gorm:"column:page_title"`
	Start          time.Time      `gorm:"column:start_tstamp"`
}

type PageViews []PageView

func (PageView) TableName() string {
	return schemaNameDerived + "." + tablePageViews
}
