package mappers

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/customeros/customeros/packages/server/leads/dto"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

// ConvertToProtoWebTrackerEvent converts from your application struct to the protobuf generated struct
func ConvertToProtoWebTrackerEvent(event *dto.WebTrackerEvent) *pb.WebTrackerEvent {
	if event == nil {
		return nil
	}

	// Convert timestamp to protobuf timestamp
	pbTimestamp := &timestamppb.Timestamp{}
	pbTimestamp.Seconds = event.Timestamp.Unix()
	pbTimestamp.Nanos = int32(event.Timestamp.Nanosecond())

	return &pb.WebTrackerEvent{
		Id:               event.ID,
		VisitorId:        event.VisitorID,
		Ip:               event.IP,
		EventType:        ConvertToProtoEventType(event.EventType),
		EventData:        event.EventData,
		Timestamp:        pbTimestamp,
		Href:             event.Href,
		Referrer:         event.Referrer,
		UserAgent:        event.UserAgent,
		Language:         event.Language,
		CookiesEnabled:   event.CookiesEnabled,
		ScreenResolution: event.ScreenResolution,
	}
}

// ConvertFromProtoWebTrackerEvent converts from the protobuf generated struct to your application struct
func ConvertFromProtoWebTrackerEvent(pbEvent *pb.WebTrackerEvent) *dto.WebTrackerEvent {
	if pbEvent == nil {
		return nil
	}

	// Convert protobuf timestamp to time.Time
	timestamp := time.Unix(pbEvent.Timestamp.Seconds, int64(pbEvent.Timestamp.Nanos))

	return &dto.WebTrackerEvent{
		ID:               pbEvent.Id,
		VisitorID:        pbEvent.VisitorId,
		IP:               pbEvent.Ip,
		EventType:        ConvertFromProtoEventType(pbEvent.EventType),
		EventData:        pbEvent.EventData,
		Timestamp:        timestamp,
		Href:             pbEvent.Href,
		Referrer:         pbEvent.Referrer,
		UserAgent:        pbEvent.UserAgent,
		Language:         pbEvent.Language,
		CookiesEnabled:   pbEvent.CookiesEnabled,
		ScreenResolution: pbEvent.ScreenResolution,
	}
}

// These are the event type conversion functions from the previous answer
func ConvertToProtoEventType(eventType enum.WebTrackerEvent) pb.WebTrackerEventType {
	switch eventType {
	case enum.WebTrackerPageExit:
		return pb.WebTrackerEventType_WEB_TRACKER_PAGE_EXIT
	case enum.WebTrackerPageView:
		return pb.WebTrackerEventType_WEB_TRACKER_PAGE_VIEW
	case enum.WebTrackerClick:
		return pb.WebTrackerEventType_WEB_TRACKER_CLICK
	case enum.WebTrackerIdentify:
		return pb.WebTrackerEventType_WEB_TRACKER_IDENTIFY
	default:
		return pb.WebTrackerEventType_WEB_TRACKER_EVENT_UNSPECIFIED
	}
}

func ConvertFromProtoEventType(eventType pb.WebTrackerEventType) enum.WebTrackerEvent {
	switch eventType {
	case pb.WebTrackerEventType_WEB_TRACKER_PAGE_EXIT:
		return enum.WebTrackerPageExit
	case pb.WebTrackerEventType_WEB_TRACKER_PAGE_VIEW:
		return enum.WebTrackerPageView
	case pb.WebTrackerEventType_WEB_TRACKER_CLICK:
		return enum.WebTrackerClick
	case pb.WebTrackerEventType_WEB_TRACKER_IDENTIFY:
		return enum.WebTrackerIdentify
	default:
		return ""
	}
}
