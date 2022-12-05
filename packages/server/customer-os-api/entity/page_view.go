package entity

import (
	"fmt"
	"time"
)

const LabelName_PageView = NodeLabel_PageView

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
}

func (pageView PageViewEntity) ToString() string {
	return fmt.Sprintf("id: %s\napplication: %s\nurl: %s", pageView.Id, pageView.Application, pageView.PageUrl)
}

type PageViewEntities []PageViewEntity

func (pageView PageViewEntity) Action() {
}

func (pageView PageViewEntity) ActionName() string {
	return NodeLabel_PageView
}
