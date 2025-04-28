package proto_mappers

import (
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

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
