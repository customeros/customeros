package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
)

func MapEntityToInteractionSessionParticipant(interactionSessionParticipantEntity *entity.InteractionSessionParticipant) any {
	switch (*interactionSessionParticipantEntity).ParticipantLabel() {
	case entity.NodeLabel_Email:
		emailEntity := (*interactionSessionParticipantEntity).(*entity.EmailEntity)
		return model.EmailParticipant{
			EmailParticipant: MapEntityToEmail(emailEntity),
			Type:             utils.StringPtrNillable(emailEntity.InteractionEventParticipantDetails.Type),
		}
	case entity.NodeLabel_PhoneNumber:
		phoneNumberEntity := (*interactionSessionParticipantEntity).(*entity.PhoneNumberEntity)
		return model.PhoneNumberParticipant{
			PhoneNumberParticipant: MapEntityToPhoneNumber(phoneNumberEntity),
			Type:                   utils.StringPtrNillable(phoneNumberEntity.InteractionEventParticipantDetails.Type),
		}
	case entity.NodeLabel_User:
		userEntity := (*interactionSessionParticipantEntity).(*entity.UserEntity)
		return model.UserParticipant{
			UserParticipant: MapEntityToUser(userEntity),
			Type:            utils.StringPtrNillable(userEntity.InteractionEventParticipantDetails.Type),
		}
	case entity.NodeLabel_Contact:
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
