package mapper

import (
	"fmt"
	"reflect"

	model2 "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToMeetingParticipant(meetingParticipantEntity *neo4jentity.MeetingParticipant) any {
	switch (*meetingParticipantEntity).EntityLabel() {
	case model2.NodeLabelUser:
		userEntity := (*meetingParticipantEntity).(*neo4jentity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
		}
	case model2.NodeLabelContact:
		contactEntity := (*meetingParticipantEntity).(*neo4jentity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
		}
	case model2.NodeLabelOrganization:
		organizationEntity := (*meetingParticipantEntity).(*neo4jentity.OrganizationEntity)
		return model.OrganizationParticipant{
			OrganizationParticipant: MapEntityToOrganization(organizationEntity),
		}
	case model2.NodeLabelEmail:
		emailEntity := (*meetingParticipantEntity).(*neo4jentity.EmailEntity)
		return model.EmailParticipant{
			EmailParticipant: MapEntityToEmail(emailEntity),
		}
	}
	fmt.Errorf("participant of type %s not identified", reflect.TypeOf(meetingParticipantEntity))
	return nil
}

func MapEntitiesToMeetingParticipants(entities *neo4jentity.MeetingParticipants) []model.MeetingParticipant {
	var meetingParticipants []model.MeetingParticipant
	for _, meetingParticipantEntity := range *entities {
		meetingParticipant := MapEntityToMeetingParticipant(&meetingParticipantEntity)
		if meetingParticipant != nil {
			meetingParticipants = append(meetingParticipants, meetingParticipant.(model.MeetingParticipant))
		}
	}
	return meetingParticipants
}
