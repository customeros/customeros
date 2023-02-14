package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

const (
	NodeLabel_PageView = "PageView"
)

var NodeLabelsByActionType = map[string]string{
	model.ActionTypePageView.String(): NodeLabel_PageView,
}

type Action interface {
	Action()
	ActionName() string
}

type ActionEntities []Action
