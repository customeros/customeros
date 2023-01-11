package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	pb "github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository/entity"
	"time"
)

type webChatMessageStoreService struct {
	proto.UnimplementedWebChatMessageStoreServiceServer
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
	customerOSService    *customerOSService
}

// sender == contact -> find contact by email
// sender == user -> find user by email

// sender == contact -> find conversation by initiator = contact and channel = webchat
// sender == user -> find conversation by id
func (s *webChatMessageStoreService) SaveMessage(ctx context.Context, input *pb.WebChatInputMessage) (*pb.Message, error) {
	//var err error
	//var conversation *gen.Conversation
	//
	var contactId string
	var userId string
	var initiatorId string

	if input.Email == "" {
		return nil, nil // TODO: return error
	}
	if input.Message == nil && input.Bytes == nil {
		return nil, nil // TODO: return error
	}

	if input.SenderType == pb.SenderType_CONTACT {
		contact, err := s.customerOSService.GetContactByEmail(input.Email)
		if err != nil {
			contactId, err = s.customerOSService.CreateContactWithEmail("openline", input.Email)
			if err != nil {
				return nil, err
			}
		} else {
			contactId = contact.Id
		}
		initiatorId = contactId
	} else if input.SenderType == pb.SenderType_USER {
		user, err := s.customerOSService.GetUserByEmail(input.Email)
		if err != nil {
			return nil, err
		} else {
			userId = user.Id
		}
		initiatorId = userId
	}

	//todo
	//conversationId, err := s.customerOSService.GetWebChatConversationIdWithContactInitiator(contactId)
	//if err != nil {

	//todo
	conversationId, err := s.customerOSService.CreateConversation("openline", initiatorId, convertSenderTypeToConversationSenderType(input.SenderType), entity.WEB_CHAT)
	if err != nil {
		return nil, err
	}
	//}

	conversationEvent := entity.ConversationEvent{
		TenantId:       "openline", //todo
		ConversationId: conversationId,
		Type:           entity.WEB_CHAT,
		Content:        *input.Message,
		Source:         entity.OPENLINE,
		Direction:      encodeConversationEventDirection(input.Direction),
		CreateDate:     time.Time{},
	}

	if input.GetDirection() == pb.MessageDirection_INBOUND {
		conversationEvent.SenderId = contactId
		conversationEvent.SenderType = entity.CONTACT
	} else {
		conversationEvent.SenderId = userId
		conversationEvent.SenderType = entity.USER
	}

	s.postgresRepositories.ConversationEventRepository.Save(&conversationEvent)

	//id := int64(conversationEvent.ID)
	//mi := &pb.Message{
	//	Id:        id,
	//	FeedId:    &conversationid,
	//	Type:      pb.MessageType_MESSAGE,
	//	Message:   conversationItem.Message,
	//	Direction: decodeDirection(conversationItem.Direction),
	//	Channel:   decodeChannel(conversationItem.Channel),
	//	Username:  message.Username,
	//	UserId:    userId,
	//	ContactId: contactId,
	//	Time:      timestamppb.New(now),
	//}
	return nil, nil
}

func convertSenderTypeToConversationSenderType(senderType pb.SenderType) entity.SenderType {
	switch senderType {
	case pb.SenderType_CONTACT:
		return entity.CONTACT
	case pb.SenderType_USER:
		return entity.USER
	default:
		return entity.CONTACT
	}
}

func encodeConversationEventDirection(direction pb.MessageDirection) entity.Direction {
	switch direction {
	case pb.MessageDirection_INBOUND:
		return entity.INBOUND
	case pb.MessageDirection_OUTBOUND:
		return entity.OUTBOUND
	default:
		return entity.OUTBOUND
	}
}

func NewWebChatMessageStoreService(driver *neo4j.Driver, postgresRepositories *repository.PostgresRepositories, customerOSService *customerOSService) *webChatMessageStoreService {
	ms := new(webChatMessageStoreService)
	ms.driver = driver
	ms.postgresRepositories = postgresRepositories
	ms.customerOSService = customerOSService
	return ms
}
