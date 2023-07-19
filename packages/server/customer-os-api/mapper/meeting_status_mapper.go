package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var meetingStatusByModel = map[model.MeetingStatus]entity.MeetingStatus{
	model.MeetingStatusUndefined: entity.MeetingStatusUndefined,
	model.MeetingStatusAccepted:  entity.MeetingStatusAccepted,
	model.MeetingStatusCanceled:  entity.MeetingStatusCanceled,
}

var meetingStatusByValue = utils.ReverseMap(meetingStatusByModel)

func MapMeetingStatusFromModel(input model.MeetingStatus) entity.MeetingStatus {
	return meetingStatusByModel[input]
}

func MapMeetingStatusToModel(input entity.MeetingStatus) model.MeetingStatus {
	return meetingStatusByValue[input]
}
