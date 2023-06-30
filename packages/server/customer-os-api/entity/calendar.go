package entity

import (
	"fmt"
	"time"
)

type CalendarEntity struct {
	Id            string
	CalType       string
	Link          string
	Primary       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

func (calendar CalendarEntity) ToString() string {
	return fmt.Sprintf("id: %s\nlink: %s", calendar.Id, calendar.Link)
}

type CalendarEntities []CalendarEntity

func (calendar CalendarEntity) Labels(tenant string) []string {
	return []string{"Calendar", "Calendar_" + tenant}
}
