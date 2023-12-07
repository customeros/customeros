package entity

import (
	"fmt"
	"time"
)

type ActionEntity struct {
	Id            string
	CreatedAt     time.Time
	Content       string
	Metadata      string
	Type          ActionType
	Source        DataSource
	AppSource     string
	DataloaderKey string
}

type ActionType string

const (
	ActionNA                                        ActionType = ""
	ActionCreated                                   ActionType = "CREATED"
	ActionRenewalLikelihoodUpdated                  ActionType = "RENEWAL_LIKELIHOOD_UPDATED"
	ActionRenewalForecastUpdated                    ActionType = "RENEWAL_FORECAST_UPDATED"
	ActionContractStatusUpdated                     ActionType = "CONTRACT_STATUS_UPDATED"
	ActionServiceLineItemPriceUpdated               ActionType = "SERVICE_LINE_ITEM_PRICE_UPDATED"
	ActionServiceLineItemQuantityUpdated            ActionType = "SERVICE_LINE_ITEM_QUANTITY_UPDATED"
	ActionServiceLineItemBilledTypeUpdated          ActionType = "SERVICE_LINE_ITEM_BILLED_TYPE_UPDATED"
	ActionServiceLineItemBilledTypeRecurringCreated ActionType = "SERVICE_LINE_ITEM_BILLED_TYPE_RECURRING_CREATED"
	ActionServiceLineItemBilledTypeOnceCreated      ActionType = "SERVICE_LINE_ITEM_BILLED_TYPE_ONCE_CREATED"
	ActionServiceLineItemBilledTypeUsageCreated     ActionType = "SERVICE_LINE_ITEM_BILLED_TYPE_USAGE_CREATED"
	ActionContractRenewed                           ActionType = "CONTRACT_RENEWED"
	ActionServiceLineItemRemoved                    ActionType = "SERVICE_LINE_ITEM_REMOVED"
)

var AllActionType = []ActionType{
	ActionCreated,
	ActionRenewalLikelihoodUpdated,
	ActionRenewalForecastUpdated,
	ActionContractStatusUpdated,
	ActionServiceLineItemPriceUpdated,
	ActionServiceLineItemQuantityUpdated,
	ActionServiceLineItemBilledTypeUpdated,
	ActionServiceLineItemBilledTypeRecurringCreated,
	ActionServiceLineItemBilledTypeOnceCreated,
	ActionServiceLineItemBilledTypeUsageCreated,
	ActionContractRenewed,
	ActionServiceLineItemRemoved,
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
