package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"time"
)

type CalendarEntity struct {
	Id            string
	CalType       string
	Link          string
	Primary       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
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
