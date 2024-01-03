package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"reflect"
)

func MapEntityToIssueParticipant(issueParticipantEntity *entity.IssueParticipant) any {
	if issueParticipantEntity == nil {
		return nil
	}
	switch (*issueParticipantEntity).ParticipantLabel() {
	case neo4jentity.NodeLabel_User:
		userEntity := (*issueParticipantEntity).(*entity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
		}
	case neo4jentity.NodeLabel_Contact:
		contactEntity := (*issueParticipantEntity).(*entity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
		}
	case neo4jentity.NodeLabel_Organization:
		organizationEntity := (*issueParticipantEntity).(*entity.OrganizationEntity)
		return model.OrganizationParticipant{
			OrganizationParticipant: MapEntityToOrganization(organizationEntity),
		}
	}

	fmt.Errorf("participant of type %s not identified", reflect.TypeOf(issueParticipantEntity))
	return nil
}

func MapEntitiesToIssueParticipants(entities *entity.IssueParticipants) []model.IssueParticipant {
	var issueParticipants []model.IssueParticipant
	for _, issueParticipantEntity := range *entities {
		issueParticipant := MapEntityToIssueParticipant(&issueParticipantEntity)
		if issueParticipant != nil {
			issueParticipants = append(issueParticipants, issueParticipant.(model.IssueParticipant))
		}
	}
	return issueParticipants
}
