package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type WebhookEvent struct {
	ExternalSystemId enum.Source
	Name             enum.FlowListenerEvent
	DataType         string
	Data             interface{}
}
