package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

const (
	NodeLabel_PageView           = "PageView"
	NodeLabel_InteractionSession = "InteractionSession"
	NodeLabel_Ticket             = "Ticket"
	NodeLabel_Conversation       = "Conversation"
	NodeLabel_Note               = "Note"
)

var NodeLabelsByActionType = map[string]string{
	model.ActionTypePageView.String():           NodeLabel_PageView,
	model.ActionTypeInteractionSession.String(): NodeLabel_InteractionSession,
	model.ActionTypeTicket.String():             NodeLabel_Ticket,
	model.ActionTypeConversation.String():       NodeLabel_Conversation,
	model.ActionTypeNote.String():               NodeLabel_Note,
}

type Action interface {
	Action()
	ActionName() string
}

type ActionEntities []Action
