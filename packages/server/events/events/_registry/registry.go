package _registry

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	events "github.com/openline-ai/openline-customer-os/packages/server/events/events"
	event "github.com/openline-ai/openline-customer-os/packages/server/events/events/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/generic"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"reflect"
)

func InitAggregate(request events.BaseEvent) eventstore.Aggregate {
	switch request.EntityType {
	case model.CONTACT:
		return event.NewContactAggregateWithTenantAndID(request.Tenant, request.EntityId)
	}

	return nil
}

var eventsRegistry = map[string]reflect.Type{
	generic.LinkEntityWithEntityV1: reflect.TypeOf(generic.LinkEntityWithEntity{}),
}

func UnmarshalBaseEventPayload(eventDataBytes []byte) (interface{}, error) {
	// Create a new instance of the type
	var vv interface{}

	// Unmarshal the JSON into the new instance
	err := json.Unmarshal(eventDataBytes, &vv)
	if err != nil {
		return nil, err
	}

	// Look up the type in the registry
	eventName := vv.(map[string]interface{})["eventName"].(string)

	if eventName == "" {
		return nil, fmt.Errorf("eventName not found in event data")
	}

	t, found := eventsRegistry[eventName]
	if !found {
		return nil, fmt.Errorf("type %s not found in registry", eventName)
	}

	v := reflect.New(t).Interface()

	err = json.Unmarshal(eventDataBytes, &v)
	if err != nil {
		return nil, err
	}

	// Return the unmarshaled struct
	return v, nil
}
