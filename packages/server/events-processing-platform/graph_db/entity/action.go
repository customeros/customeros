package entity

import (
	"time"
)

type ActionEntity struct {
	Id            string
	CreatedAt     time.Time
	Content       string
	Metadata      string
	Type          ActionType
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type ActionType string

const (
	ActionNA                          ActionType = ""
	ActionCreated                     ActionType = "CREATED"
	ActionContractStatusUpdated       ActionType = "CONTRACT_STATUS_UPDATED"
	ActionServiceLineItemPriceUpdated ActionType = "SERVICE_LINE_ITEM_PRICE_UPDATED"
)

var AllActionType = []ActionType{
	ActionCreated,
	ActionContractStatusUpdated,
	ActionServiceLineItemPriceUpdated,
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

func (ActionEntity) IsTimelineEvent() {
}

func (ActionEntity) TimelineEventLabel() string {
	return NodeLabel_Action
}
