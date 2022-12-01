package entity

import (
	"fmt"
	"time"
)

const LabelName_PageViewAction = ActionName_PageView

type PageViewActionEntity struct {
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

func (pageViewAction PageViewActionEntity) ToString() string {
	return fmt.Sprintf("id: %s\napplication: %s\nurl: %s", pageViewAction.Id, pageViewAction.Application, pageViewAction.PageUrl)
}

type PageViewActionEntities []PageViewActionEntity

func (pageViewAction PageViewActionEntity) Action() {
}

func (pageViewAction PageViewActionEntity) ActionName() string {
	return ActionName_PageView
}
