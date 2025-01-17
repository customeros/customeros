package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

var meetingStatusByModel = map[model.MeetingStatus]neo4jenum.MeetingStatus{
	model.MeetingStatusUndefined: neo4jenum.MeetingStatusUndefined,
	model.MeetingStatusAccepted:  neo4jenum.MeetingStatusAccepted,
	model.MeetingStatusCanceled:  neo4jenum.MeetingStatusCanceled,
}

var meetingStatusByValue = utils.ReverseMap(meetingStatusByModel)

func MapMeetingStatusFromModel(input model.MeetingStatus) neo4jenum.MeetingStatus {
	return meetingStatusByModel[input]
}

func MapMeetingStatusToModel(input neo4jenum.MeetingStatus) model.MeetingStatus {
	return meetingStatusByValue[input]
}
