package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var meetingBookingAssignmentMethodByModel = map[model.MeetingBookingAssignmentMethod]enum.MeetingBookingAssignmentMethod{
	model.MeetingBookingAssignmentMethodRoundRobinMaxFairness:     enum.MeetingBookingAssignmentMethodRoundRobinMaxFairness,
	model.MeetingBookingAssignmentMethodRoundRobinMaxAvailability: enum.MeetingBookingAssignmentMethodRoundRobinMaxAvailability,
	model.MeetingBookingAssignmentMethodCustom:                    enum.MeetingBookingAssignmentMethodCustom,
}

var meetingBookingAssignmentMethodByValue = utils.ReverseMap(meetingBookingAssignmentMethodByModel)

func MapMeetingBookingAssignmentMethodFromModel(method model.MeetingBookingAssignmentMethod) enum.MeetingBookingAssignmentMethod {
	if method == "" {
		return enum.DefaultMeetingBookingAssignmentMethod
	}
	return meetingBookingAssignmentMethodByModel[method]
}

func MapMeetingBookingAssignmentMethodToModel(method enum.MeetingBookingAssignmentMethod) model.MeetingBookingAssignmentMethod {
	if method == "" {
		return MapMeetingBookingAssignmentMethodToModel(enum.DefaultMeetingBookingAssignmentMethod)
	}
	return meetingBookingAssignmentMethodByValue[method]
}
