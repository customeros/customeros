package enum

type ActionType string

const (
	ActionNA                             ActionType = ""
	ActionCreated                        ActionType = "CREATED"
	ActionContractStatusUpdated          ActionType = "CONTRACT_STATUS_UPDATED"
	ActionServiceLineItemPriceUpdated    ActionType = "SERVICE_LINE_ITEM_PRICE_UPDATED"
	ActionServiceLineItemQuantityUpdated ActionType = "SERVICE_LINE_ITEM_QUANTITY_UPDATED"
	// Deprecated
	ActionServiceLineItemBilledTypeUpdated          ActionType = "SERVICE_LINE_ITEM_BILLED_TYPE_UPDATED"
	ActionServiceLineItemBilledTypeRecurringCreated ActionType = "SERVICE_LINE_ITEM_BILLED_TYPE_RECURRING_CREATED"
	ActionServiceLineItemBilledTypeOnceCreated      ActionType = "SERVICE_LINE_ITEM_BILLED_TYPE_ONCE_CREATED"
	ActionServiceLineItemBilledTypeUsageCreated     ActionType = "SERVICE_LINE_ITEM_BILLED_TYPE_USAGE_CREATED"
	ActionContractRenewed                           ActionType = "CONTRACT_RENEWED"
	ActionServiceLineItemRemoved                    ActionType = "SERVICE_LINE_ITEM_REMOVED"
	ActionOnboardingStatusChanged                   ActionType = "ONBOARDING_STATUS_CHANGED"
	ActionRenewalLikelihoodUpdated                  ActionType = "RENEWAL_LIKELIHOOD_UPDATED"
	ActionRenewalForecastUpdated                    ActionType = "RENEWAL_FORECAST_UPDATED"
	ActionInvoiceIssued                             ActionType = "INVOICE_ISSUED"
	ActionInvoicePaid                               ActionType = "INVOICE_PAID"
	ActionInvoiceVoided                             ActionType = "INVOICE_VOIDED"
	ActionInvoiceOverdue                            ActionType = "INVOICE_OVERDUE"
	ActionInvoiceSent                               ActionType = "INVOICE_SENT"
	ActionInteractionEventRead                      ActionType = "INTERACTION_EVENT_READ"
)

var AllActionType = []ActionType{
	ActionCreated,
	ActionContractStatusUpdated,
	ActionServiceLineItemPriceUpdated,
	ActionServiceLineItemQuantityUpdated,
	ActionServiceLineItemBilledTypeUpdated,
	ActionServiceLineItemBilledTypeRecurringCreated,
	ActionServiceLineItemBilledTypeOnceCreated,
	ActionServiceLineItemBilledTypeUsageCreated,
	ActionContractRenewed,
	ActionServiceLineItemRemoved,
	ActionOnboardingStatusChanged,
	ActionRenewalLikelihoodUpdated,
	ActionInvoiceIssued,
	ActionInvoicePaid,
	ActionInvoiceVoided,
	ActionInvoiceOverdue,
	ActionInvoiceSent,
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
