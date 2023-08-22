package entity

import (
	"fmt"
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

func (pageView PageViewEntity) ToString() string {
	return fmt.Sprintf("id: %s\napplication: %s\nurl: %s", pageView.Id, pageView.Application, pageView.PageUrl)
}

type PageViewEntities []PageViewEntity

func (PageViewEntity) IsTimelineEvent() {
}

func (pageView *PageViewEntity) SetDataloaderKey(key string) {
	pageView.DataloaderKey = key
}

func (pageView PageViewEntity) GetDataloaderKey() string {
	return pageView.DataloaderKey
}

func (PageViewEntity) TimelineEventLabel() string {
	return NodeLabel_PageView
}

func (pageView PageViewEntity) Labels() []string {
	return []string{"PageView", "TimelineEvent"}
}
