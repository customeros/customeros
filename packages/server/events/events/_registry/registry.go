package _registry

import (
	"encoding/json"
	"fmt"
	events "github.com/openline-ai/openline-customer-os/packages/server/events/events"
	event "github.com/openline-ai/openline-customer-os/packages/server/events/events/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/generic"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"reflect"
)

func InitAggregate(request events.BaseEvent) eventstore.Aggregate {
	switch request.EntityType {
	case events.CONTACT:
		return event.NewContactAggregateWithTenantAndID(request.Tenant, request.EntityId)
	}

	return nil
}

var eventsRegistry = map[string]reflect.Type{
	generic.UpsertEmailToEntityV1: reflect.TypeOf(generic.UpsertEmailToEntityEvent{}),
}

func UnmarshalEventPayload(structName string, jsonString string) (interface{}, error) {
	// Look up the type in the registry
	t, found := eventsRegistry[structName]
	if !found {
		return nil, fmt.Errorf("type %s not found in registry", structName)
	}

	// Create a new instance of the type
	v := reflect.New(t).Interface()

	// Unmarshal the JSON into the new instance
	err := json.Unmarshal([]byte(jsonString), v)
	if err != nil {
		return nil, err
	}

	// Return the unmarshaled struct
	return v, nil
}
