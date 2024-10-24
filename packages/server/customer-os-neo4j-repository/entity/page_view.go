package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type PageViewEntity struct {
	DataLoaderKey
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
}

func (PageViewEntity) IsTimelineEvent() {
}

func (PageViewEntity) TimelineEventLabel() string {
	return model.NodeLabelPageView
}

func (e *PageViewEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *PageViewEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
