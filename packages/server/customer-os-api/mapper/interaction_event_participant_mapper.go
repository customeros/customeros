package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"reflect"
)

func MapEntityToInteractionEventParticipant(interactionEventParticipantEntity *neo4jentity.InteractionEventParticipant) any {
	switch (*interactionEventParticipantEntity).EntityLabel() {
	case model2.NodeLabelEmail:
		emailEntity := (*interactionEventParticipantEntity).(*neo4jentity.EmailEntity)
		return model.EmailParticipant{
			EmailParticipant: MapEntityToEmail(emailEntity),
			Type:             utils.StringPtrNillable(emailEntity.InteractionEventParticipantDetails.Type),
		}
	case model2.NodeLabelPhoneNumber:
		phoneNumberEntity := (*interactionEventParticipantEntity).(*neo4jentity.PhoneNumberEntity)
		return model.PhoneNumberParticipant{
			PhoneNumberParticipant: MapEntityToPhoneNumber(phoneNumberEntity),
			Type:                   utils.StringPtrNillable(phoneNumberEntity.InteractionEventParticipantDetails.Type),
		}
	case model2.NodeLabelUser:
		userEntity := (*interactionEventParticipantEntity).(*neo4jentity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
			Type:            utils.StringPtrNillable(userEntity.InteractionEventParticipantDetails.Type),
		}
	case model2.NodeLabelContact:
		contactEntity := (*interactionEventParticipantEntity).(*neo4jentity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
			Type:               utils.StringPtrNillable(contactEntity.InteractionEventParticipantDetails.Type),
		}
	case model2.NodeLabelOrganization:
		organizationEntity := (*interactionEventParticipantEntity).(*neo4jentity.OrganizationEntity)
		return model.OrganizationParticipant{
			OrganizationParticipant: MapEntityToOrganization(organizationEntity),
			Type:                    utils.StringPtrNillable(organizationEntity.InteractionEventParticipantDetails.Type),
		}
	case model2.NodeLabelJobRole:
		jobRoleEntity := (*interactionEventParticipantEntity).(*neo4jentity.JobRoleEntity)
		return model.JobRoleParticipant{
			JobRoleParticipant: MapEntityToJobRole(jobRoleEntity),
			Type:               utils.StringPtrNillable(jobRoleEntity.InteractionEventParticipantDetails.Type),
		}
	}

	fmt.Errorf("participant of type %s not identified", reflect.TypeOf(interactionEventParticipantEntity))
	return nil
}

func MapEntitiesToInteractionEventParticipants(entities *neo4jentity.InteractionEventParticipants) []model.InteractionEventParticipant {
	var interactionEventParticipants []model.InteractionEventParticipant
	for _, interactionEventParticipantEntity := range *entities {
		interactionEventParticipant := MapEntityToInteractionEventParticipant(&interactionEventParticipantEntity)
		if interactionEventParticipant != nil {
			interactionEventParticipants = append(interactionEventParticipants, interactionEventParticipant.(model.InteractionEventParticipant))
		}
	}
	return interactionEventParticipants
}
