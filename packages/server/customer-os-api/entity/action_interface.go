package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

const (
	// FIXME alexb check where it's used
	NodeLabel_PageView           = "PageView"
	NodeLabel_InteractionSession = "InteractionSession"
)

var NodeLabelsByActionType = map[string]string{
	model.ActionTypePageView.String():           NodeLabel_PageView,
	model.ActionTypeInteractionSession.String(): NodeLabel_InteractionSession,
}

type Action interface {
	Action()
	ActionName() string
}

type ActionEntities []Action
