package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

const (
	NodeLabel_PageView           = "PageView"
	NodeLabel_InteractionSession = "InteractionSession"
	NodeLabel_Ticket             = "Ticket"
)

var NodeLabelsByActionType = map[string]string{
	model.ActionTypePageView.String():           NodeLabel_PageView,
	model.ActionTypeInteractionSession.String(): NodeLabel_InteractionSession,
	model.ActionTypeTicket.String():             NodeLabel_Ticket,
}

type Action interface {
	Action()
	ActionName() string
}

type ActionEntities []Action
