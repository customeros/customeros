package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"reflect"
)

func MapEntityToIssueParticipant(issueParticipantEntity *neo4jentity.IssueParticipant) any {
	if issueParticipantEntity == nil {
		return nil
	}
	switch (*issueParticipantEntity).EntityLabel() {
	case model2.NodeLabelUser:
		userEntity := (*issueParticipantEntity).(*neo4jentity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
		}
	case model2.NodeLabelContact:
		contactEntity := (*issueParticipantEntity).(*neo4jentity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
		}
	case model2.NodeLabelOrganization:
		organizationEntity := (*issueParticipantEntity).(*neo4jentity.OrganizationEntity)
		return model.OrganizationParticipant{
			OrganizationParticipant: MapEntityToOrganization(organizationEntity),
		}
	}

	fmt.Errorf("participant of type %s not identified", reflect.TypeOf(issueParticipantEntity))
	return nil
}

func MapEntitiesToIssueParticipants(entities *neo4jentity.IssueParticipants) []model.IssueParticipant {
	var issueParticipants []model.IssueParticipant
	for _, issueParticipantEntity := range *entities {
		issueParticipant := MapEntityToIssueParticipant(&issueParticipantEntity)
		if issueParticipant != nil {
			issueParticipants = append(issueParticipants, issueParticipant.(model.IssueParticipant))
		}
	}
	return issueParticipants
}
