package data_fields

import (
	"fmt"
)

// Add all supported flow events here //

const (
	// Cal.com events
	FlowEventCalComBookingCreated     FlowEvent = "cal_com.booking.created"
	FlowEventCalComBookingRescheduled FlowEvent = "cal_com.booking.rescheduled"
	FlowEventCalComBookingCancelled   FlowEvent = "cal_com.booking.cancelled"

	FlowEventFathomMeetingSummaryCreated FlowEvent = "fathom.meeting_summary.created"
	FlowEventGrainMeetingSummaryCreated  FlowEvent = "grain.meeting_summary.created"
)

// Define schema for all events here //

var eventRegistry = func() map[FlowEvent]FlowEventDefinition {
	events := []FlowEventDefinition{
		{
			Event:       FlowEventFathomMeetingSummaryCreated,
			Integration: IntegrationFathom,
			Schema:      MeetingSummaryCreatedEvent{},
			Description: "When a new Fathom meeting summary is created",
		},
		{
			Event:       FlowEventGrainMeetingSummaryCreated,
			Integration: IntegrationGrain,
			Schema:      MeetingSummaryCreatedEvent{},
			Description: "When a new Grain meeting summary is created",
		},
	}

	m := make(map[FlowEvent]FlowEventDefinition)
	for _, e := range events {
		m[e.Event] = e
	}
	return m
}()

func ParseEventType(s string) (FlowEvent, error) {
	if event, ok := eventRegistry[FlowEvent(s)]; ok {
		return event.Event, nil
	}
	return "", fmt.Errorf("invalid event type: %s", s)
}
