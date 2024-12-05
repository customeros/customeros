package dto

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type WebhookEvent struct {
	ExternalSystemId enum.ExternalSystemId
	Name             commonenum.FlowEvent
	DataType         string
	Data             interface{}
}
