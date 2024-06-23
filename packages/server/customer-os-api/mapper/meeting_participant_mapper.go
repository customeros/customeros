package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"reflect"
)

func MapEntityToMeetingParticipant(meetingParticipantEntity *neo4jentity.MeetingParticipant) any {
	switch (*meetingParticipantEntity).EntityLabel() {
	case neo4jutil.NodeLabelUser:
		userEntity := (*meetingParticipantEntity).(*entity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
		}
	case neo4jutil.NodeLabelContact:
		contactEntity := (*meetingParticipantEntity).(*neo4jentity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
		}
	case neo4jutil.NodeLabelOrganization:
		organizationEntity := (*meetingParticipantEntity).(*neo4jentity.OrganizationEntity)
		return model.OrganizationParticipant{
			OrganizationParticipant: MapEntityToOrganization(organizationEntity),
		}
	case neo4jutil.NodeLabelEmail:
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
