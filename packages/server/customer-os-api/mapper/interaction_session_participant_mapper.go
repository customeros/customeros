package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"reflect"
)

func MapEntityToInteractionSessionParticipant(interactionSessionParticipantEntity *neo4jentity.InteractionSessionParticipant) any {
	switch (*interactionSessionParticipantEntity).EntityLabel() {
	case model2.NodeLabelEmail:
		emailEntity := (*interactionSessionParticipantEntity).(*neo4jentity.EmailEntity)
		return model.EmailParticipant{
			EmailParticipant: MapEntityToEmail(emailEntity),
			Type:             utils.StringPtrNillable(emailEntity.InteractionEventParticipantDetails.Type),
		}
	case model2.NodeLabelPhoneNumber:
		phoneNumberEntity := (*interactionSessionParticipantEntity).(*neo4jentity.PhoneNumberEntity)
		return model.PhoneNumberParticipant{
			PhoneNumberParticipant: MapEntityToPhoneNumber(phoneNumberEntity),
			Type:                   utils.StringPtrNillable(phoneNumberEntity.InteractionEventParticipantDetails.Type),
		}
	case model2.NodeLabelUser:
		userEntity := (*interactionSessionParticipantEntity).(*entity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
			Type:            utils.StringPtrNillable(userEntity.InteractionEventParticipantDetails.Type),
		}
	case model2.NodeLabelContact:
		contactEntity := (*interactionSessionParticipantEntity).(*neo4jentity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
			Type:               utils.StringPtrNillable(contactEntity.InteractionEventParticipantDetails.Type),
		}
	}
	fmt.Errorf("participant of type %s not identified", reflect.TypeOf(interactionSessionParticipantEntity))
	return nil
}

func MapEntitiesToInteractionSessionParticipants(entities *neo4jentity.InteractionSessionParticipants) []model.InteractionSessionParticipant {
	var interactionSessionParticipants []model.InteractionSessionParticipant
	for _, interactionSessionParticipantEntity := range *entities {
		interactionSessionParticipant := MapEntityToInteractionSessionParticipant(&interactionSessionParticipantEntity)
		if interactionSessionParticipant != nil {
			interactionSessionParticipants = append(interactionSessionParticipants, interactionSessionParticipant.(model.InteractionSessionParticipant))
		}
	}
	return interactionSessionParticipants
}
