package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"reflect"
)

func MapEntityToMeetingParticipant(meetingParticipantEntity *entity.MeetingParticipant) any {
	switch (*meetingParticipantEntity).MeetingParticipantLabel() {
	case entity.NodeLabel_User:
		userEntity := (*meetingParticipantEntity).(*entity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
		}
	case entity.NodeLabel_Contact:
		contactEntity := (*meetingParticipantEntity).(*entity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
		}
	case entity.NodeLabel_Organization:
		organizationEntity := (*meetingParticipantEntity).(*entity.OrganizationEntity)
		return model.OrganizationParticipant{
			OrganizationParticipant: MapEntityToOrganization(organizationEntity),
		}
	case entity.NodeLabel_Email:
		emailEntity := (*meetingParticipantEntity).(*entity.EmailEntity)
		return model.EmailParticipant{
			EmailParticipant: MapEntityToEmail(emailEntity),
		}
	}
	fmt.Errorf("participant of type %s not identified", reflect.TypeOf(meetingParticipantEntity))
	return nil
}

func MapEntitiesToMeetingParticipants(entities *entity.MeetingParticipants) []model.MeetingParticipant {
	var meetingParticipants []model.MeetingParticipant
	for _, meetingParticipantEntity := range *entities {
		meetingParticipant := MapEntityToMeetingParticipant(&meetingParticipantEntity)
		if meetingParticipant != nil {
			meetingParticipants = append(meetingParticipants, meetingParticipant.(model.MeetingParticipant))
		}
	}
	return meetingParticipants
}
