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
	var participantUsername string

	if input.ConversationId == nil && input.Email == nil {
		return nil, errors.New("conversationId or email must be provided")
	}
	if input.Message == nil && input.Bytes == nil {
		return nil, errors.New("message or bytes must be provided")
	}

	tenant := "openline" //TODO get tenant from context

	if input.ConversationId != nil {
		conv, err := s.customerOSService.GetConversationById(tenant, *input.ConversationId)
		if err != nil {
			return nil, err
		}
		conversation = conv
	}

	if input.SenderType == msProto.SenderType_CONTACT {
		contactId, err := s.getContactIdWithEmailOrCreate(tenant, *input.Email)
		if err != nil {
			return nil, err
		}
		participantId = contactId
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
		conv, err := s.getActiveConversationOrCreate(tenant, participantId, *input.Email, input.SenderType)
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

	_, err := s.customerOSService.UpdateConversation(tenant, conversation.Id, participantId, s.commonStoreService.ConvertMSSenderTypeToEntitySenderType(input.SenderType))
	if err != nil {
		return nil, err
	}

	return s.commonStoreService.EncodeConversationEventToMS(conversationEvent), nil
}

func (s *webChatMessageStoreService) getContactIdWithEmailOrCreate(tenant string, email string) (string, error) {
	contact, err := s.customerOSService.GetContactByEmail(email)
	if err != nil {
		contactId, err := s.customerOSService.CreateContactWithEmail(tenant, email)
		if err != nil {
			return "", err
		}
		return contactId, nil
	} else {
		return contact.Id, nil
	}
}

func (s *webChatMessageStoreService) getActiveConversationOrCreate(tenant string, participantId string, initiatorUsername string, senderType msProto.SenderType) (*Conversation, error) {
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
		conversation, err = s.customerOSService.CreateConversation(tenant, participantId, initiatorUsername, s.commonStoreService.ConvertMSSenderTypeToEntitySenderType(senderType), entity.WEB_CHAT)
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
