package service

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	msProto "github.com/openline-ai/openline-customer-os/packages/server/message-store/proto/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository/entity"
	"time"
)

type webChatMessageStoreService struct {
	msProto.UnimplementedWebChatMessageStoreServiceServer
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
	customerOSService    *customerOSService
	commonStoreService   *commonStoreService
}

func (s *webChatMessageStoreService) SaveMessage(ctx context.Context, input *msProto.WebChatInputMessage) (*msProto.Message, error) {
	//var err error
	//var conversation *gen.Conversation
	//
	var conversation *Conversation
	var participantId string
	var participantFirstName string
	var participantLastName string
	var participantUsername string

	if input.ConversationId == nil && input.Email == nil {
		return nil, errors.New("conversationId or email must be provided")
	}
	if input.Message == nil && input.Bytes == nil {
		return nil, errors.New("message or bytes must be provided")
	}

	tenant := "openline" //TODO get tenant from context

	if input.ConversationId != nil {
		conv, err := s.customerOSService.GetConversationById(*input.ConversationId)
		if err != nil {
			return nil, err
		}
		conversation = conv
	}

	if input.SenderType == msProto.SenderType_CONTACT {
		contact, err := s.getContactWithEmailOrCreate(tenant, *input.Email)
		if err != nil {
			return nil, err
		}
		participantId = contact.Id
	} else if input.SenderType == msProto.SenderType_USER {
		user, err := s.customerOSService.GetUserByEmail(*input.Email)
		if err != nil {
			return nil, err
		}
		participantId = user.Id
	}

	if participantId == "" {
		return nil, errors.New("participant not found")
	}

	participantUsername = *input.Email

	if input.ConversationId == nil {
		conv, err := s.getActiveConversationOrCreate(tenant, participantId, participantFirstName, participantLastName, *input.Email, input.SenderType)
		if err != nil {
			return nil, err
		}
		conversation = conv
	}

	conversationEvent := entity.ConversationEvent{
		TenantName:     tenant,
		ConversationId: conversation.Id,
		Type:           entity.WEB_CHAT,
		Subtype:        s.commonStoreService.ConvertMSSubtypeToEntitySubtype(input.Subtype),
		Content:        *input.Message,
		Source:         entity.OPENLINE,
		Direction:      s.commonStoreService.ConvertMSDirectionToEntityDirection(input.Direction),
		CreateDate:     time.Now(),

		InitiatorUsername: conversation.InitiatorUsername,

		SenderId:       participantId,
		SenderUsername: participantUsername,

		OriginalJson: "TODO",
	}

	if input.GetDirection() == msProto.MessageDirection_INBOUND {
		conversationEvent.SenderType = entity.CONTACT
	} else {
		conversationEvent.SenderType = entity.USER
	}

	s.postgresRepositories.ConversationEventRepository.Save(&conversationEvent)

	senderType := s.commonStoreService.ConvertMSSenderTypeToEntitySenderType(input.SenderType)

	previewMessage := ""
	if input.Message != nil {
		s := *input.Message
		len := len(s)
		if len > 20 {
			len = 20
		}
		previewMessage = s[0:len]
	}
	_, err := s.customerOSService.UpdateConversation(tenant, conversation.Id, participantId, senderType, participantFirstName, participantLastName, previewMessage)
	if err != nil {
		return nil, err
	}

	return s.commonStoreService.EncodeConversationEventToMS(conversationEvent), nil
}

func (s *webChatMessageStoreService) getContactWithEmailOrCreate(tenant string, email string) (Contact, error) {
	contact, err := s.customerOSService.GetContactByEmail(email)
	if err != nil {
		contact, err = s.customerOSService.CreateContactWithEmail(tenant, email)
		if err != nil {
			return Contact{}, err
		}
		if contact == nil {
			return Contact{}, errors.New("contact not found and could not be created")
		}
		return *contact, nil
	} else {
		return *contact, nil
	}
}

func (s *webChatMessageStoreService) getActiveConversationOrCreate(
	tenant string,
	participantId string,
	firstName string,
	lastname string,
	username string,
	senderType msProto.SenderType,
) (*Conversation, error) {
	var conversation *Conversation
	var err error

	if senderType == msProto.SenderType_CONTACT {
		conversation, err = s.customerOSService.GetWebChatConversationWithContactInitiator(tenant, participantId)
	} else if senderType == msProto.SenderType_USER {
		conversation, err = s.customerOSService.GetWebChatConversationWithUserInitiator(tenant, participantId)
	}

	if err != nil {
		return nil, err
	}

	if conversation == nil {
		conversation, err = s.customerOSService.CreateConversation(tenant, participantId, firstName, lastname, username, s.commonStoreService.ConvertMSSenderTypeToEntitySenderType(senderType), entity.WEB_CHAT)
	}
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func NewWebChatMessageStoreService(driver *neo4j.Driver, postgresRepositories *repository.PostgresRepositories, customerOSService *customerOSService, commonStoreService *commonStoreService) *webChatMessageStoreService {
	ms := new(webChatMessageStoreService)
	ms.driver = driver
	ms.postgresRepositories = postgresRepositories
	ms.customerOSService = customerOSService
	ms.commonStoreService = commonStoreService
	return ms
}
