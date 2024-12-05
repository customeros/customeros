package dto

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type WebhookEvent struct {
	ExternalSystemId enum.ExternalSystemId
	Name             commonenum.FlowEvent
	DataType         string
	Data             any
}

func getEventName(externalSystem enum.ExternalSystemId, resource, action string) (commonenum.FlowEvent, error) {
	name := fmt.Sprintf("%s.%s.%s", externalSystem.String(), resource, action)
	return commonenum.GetFlowEvent(name)
}
