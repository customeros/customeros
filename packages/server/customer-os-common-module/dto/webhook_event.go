package dto

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

type WebhookEvent[T any] struct {
	ExternalSystemId enum.ExternalSystemId
	Resource         string
	Action           string
	Name             string
	Data             *T
}

func NewWebhookEvent(externalSystem enum.ExternalSystemId, resource, action string, data any) WebhookEvent[any] {
	return WebhookEvent[any]{
		ExternalSystemId: externalSystem,
		Resource:         resource,
		Action:           action,
		Name:             getEventName(externalSystem, resource, action),
		Data:             &data,
	}
}

func getEventName(externalSystem enum.ExternalSystemId, resource, action string) string {
	return fmt.Sprintf("%s.%s.%s", externalSystem.String(), resource, action)
}
