package _registry

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder"
	reminderevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"reflect"
)

func InitAggregate(request event.BaseEvent) eventstore.Aggregate {
	switch request.EntityType {
	case model.CONTACT:
		return contact.NewContactAggregateWithTenantAndID(request.Tenant, request.EntityId)
	case model.REMINDER:
		return reminder.NewReminderAggregateWithTenantAndID(request.Tenant, request.EntityId)
	case model.OPPORTUNITY:
		return opportunityevent.NewOpportunityAggregateWithTenantAndID(request.Tenant, request.EntityId)
	}
	return nil
}

var eventsRegistry = map[string]reflect.Type{
	reminderevent.ReminderCreateV1:       reflect.TypeOf(reminderevent.ReminderCreateEvent{}),
	reminderevent.ReminderUpdateV1:       reflect.TypeOf(reminderevent.ReminderUpdateEvent{}),
	reminderevent.ReminderNotificationV1: reflect.TypeOf(reminderevent.ReminderNotificationEvent{}),
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
