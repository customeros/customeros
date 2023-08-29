package entity

import (
	"fmt"
	"time"
)

type ActionEntity struct {
	Id            string
	CreatedAt     time.Time
	Content       string
	Type          ActionType
	Source        DataSource
	AppSource     string
	DataloaderKey string
}

type ActionType string

const (
	ActionNA                       ActionType = ""
	ActionCreated                  ActionType = "CREATED"
	ActionRenewalLikelihoodUpdated ActionType = "RENEWAL_LIKELIHOOD_UPDATED"
)

var AllActionType = []ActionType{
	ActionCreated,
	ActionRenewalLikelihoodUpdated,
}

func GetActionType(s string) ActionType {
	if IsValidActionType(s) {
		return ActionType(s)
	}
	return ActionNA
}

func IsValidActionType(s string) bool {
	for _, ds := range AllActionType {
		if ds == ActionType(s) {
			return true
		}
	}
	return false
}

func (action ActionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", action.Id, action.Type)
}

func (action ActionEntity) GetDataloaderKey() string {
	return action.DataloaderKey
}

func (action *ActionEntity) SetDataloaderKey(key string) {
	action.DataloaderKey = key
}

func (ActionEntity) IsTimelineEvent() {
}

func (ActionEntity) TimelineEventLabel() string {
	return NodeLabel_Action
}

type ActionEntities []ActionEntity

func (action ActionEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Action,
		NodeLabel_Action + "_" + tenant,
		NodeLabel_TimelineEvent,
		NodeLabel_TimelineEvent + "_" + tenant,
	}
}
