package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"reflect"
)

func MapEntityToIssueParticipant(issueParticipantEntity *neo4jentity.IssueParticipant) any {
	if issueParticipantEntity == nil {
		return nil
	}
	switch (*issueParticipantEntity).EntityLabel() {
	case neo4jutil.NodeLabelUser:
		userEntity := (*issueParticipantEntity).(*entity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
		}
	case neo4jutil.NodeLabelContact:
		contactEntity := (*issueParticipantEntity).(*entity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapLocalEntityToContact(contactEntity),
		}
	case neo4jutil.NodeLabelOrganization:
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
