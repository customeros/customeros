package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"reflect"
)

func MapEntityToMeetingParticipant(meetingParticipantEntity *entity.MeetingParticipant) any {
	switch (*meetingParticipantEntity).MeetingParticipantLabel() {
	case neo4jentity.NodeLabel_User:
		userEntity := (*meetingParticipantEntity).(*entity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
		}
	case neo4jentity.NodeLabel_Contact:
		contactEntity := (*meetingParticipantEntity).(*entity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
		}
	case neo4jentity.NodeLabel_Organization:
		organizationEntity := (*meetingParticipantEntity).(*entity.OrganizationEntity)
		return model.OrganizationParticipant{
			OrganizationParticipant: MapEntityToOrganization(organizationEntity),
		}
	case neo4jentity.NodeLabel_Email:
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
