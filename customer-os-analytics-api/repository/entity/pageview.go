package entity

type PageViewEntity struct {
	ID             string `gorm:"column:page_view_id"`
	SessionID      string `gorm:"column:domain_sessionid"`
	OrderInSession int    `gorm:"column:page_view_in_session_index"`
	EngagedTime    int    `gorm:"column:engaged_time_in_s"`
	Path           string `gorm:"column:page_urlpath"`
	Title          string `gorm:"column:page_title"`
}

type PageViewEntities []PageViewEntity

func (PageViewEntity) TableName() string {
	return "derived.page_views"
}
