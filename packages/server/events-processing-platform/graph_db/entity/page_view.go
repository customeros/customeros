package entity

import (
	"time"
)

type PageViewEntity struct {
	Id             string
	Application    string
	TrackerName    string
	SessionId      string
	PageUrl        string
	PageTitle      string
	OrderInSession int64
	EngagedTime    int64
	StartedAt      time.Time
	EndedAt        time.Time
	Source         DataSource
	SourceOfTruth  DataSource
	AppSource      string

	DataloaderKey string
}

func (PageViewEntity) IsTimelineEvent() {
}

func (PageViewEntity) TimelineEventLabel() string {
	return NodeLabel_PageView
}
