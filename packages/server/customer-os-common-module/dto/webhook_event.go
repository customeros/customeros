package dto

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type WebhookEvent[T any] struct {
	ExternalSystemId enum.ExternalSystemId
	Resource         string
	Action           string
	Name             commonenum.FlowEvent
	Data             *T
}

func NewWebhookEvent(externalSystem enum.ExternalSystemId, resource, action string, data any) (WebhookEvent[any], error) {
	eventName, err := getEventName(externalSystem, resource, action)
	if err != nil {
		return WebhookEvent[any]{}, err
	}
	event := WebhookEvent[any]{
		ExternalSystemId: externalSystem,
		Resource:         resource,
		Action:           action,
		Name:             eventName,
		Data:             &data,
	}
	return event, nil
}

func getEventName(externalSystem enum.ExternalSystemId, resource, action string) (commonenum.FlowEvent, error) {
	name := fmt.Sprintf("%s.%s.%s", externalSystem.String(), resource, action)
	return commonenum.GetFlowEvent(name)
}
