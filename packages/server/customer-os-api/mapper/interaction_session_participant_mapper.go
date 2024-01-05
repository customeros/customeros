package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"reflect"
)

func MapEntityToInteractionSessionParticipant(interactionSessionParticipantEntity *entity.InteractionSessionParticipant) any {
	switch (*interactionSessionParticipantEntity).ParticipantLabel() {
	case neo4jentity.NodeLabelEmail:
		emailEntity := (*interactionSessionParticipantEntity).(*entity.EmailEntity)
		return model.EmailParticipant{
			EmailParticipant: MapEntityToEmail(emailEntity),
			Type:             utils.StringPtrNillable(emailEntity.InteractionEventParticipantDetails.Type),
		}
	case neo4jentity.NodeLabelPhoneNumber:
		phoneNumberEntity := (*interactionSessionParticipantEntity).(*entity.PhoneNumberEntity)
		return model.PhoneNumberParticipant{
			PhoneNumberParticipant: MapEntityToPhoneNumber(phoneNumberEntity),
			Type:                   utils.StringPtrNillable(phoneNumberEntity.InteractionEventParticipantDetails.Type),
		}
	case neo4jentity.NodeLabelUser:
		userEntity := (*interactionSessionParticipantEntity).(*entity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
			Type:            utils.StringPtrNillable(userEntity.InteractionEventParticipantDetails.Type),
		}
	case neo4jentity.NodeLabelContact:
		contactEntity := (*interactionSessionParticipantEntity).(*entity.ContactEntity)
		return model.ContactParticipant{
			ContactParticipant: MapEntityToContact(contactEntity),
			Type:               utils.StringPtrNillable(contactEntity.InteractionEventParticipantDetails.Type),
		}
	}
	fmt.Errorf("participant of type %s not identified", reflect.TypeOf(interactionSessionParticipantEntity))
	return nil
}

func MapEntitiesToInteractionSessionParticipants(entities *entity.InteractionSessionParticipants) []model.InteractionSessionParticipant {
	var interactionSessionParticipants []model.InteractionSessionParticipant
	for _, interactionSessionParticipantEntity := range *entities {
		interactionSessionParticipant := MapEntityToInteractionSessionParticipant(&interactionSessionParticipantEntity)
		if interactionSessionParticipant != nil {
			interactionSessionParticipants = append(interactionSessionParticipants, interactionSessionParticipant.(model.InteractionSessionParticipant))
		}
	}
	return interactionSessionParticipants
}
