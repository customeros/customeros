package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	commonModuleService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	msProto "github.com/openline-ai/openline-customer-os/packages/server/message-store-api/proto/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type MessageService struct {
	msProto.UnimplementedMessageStoreServiceServer
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
	customerOSService    *CustomerOSService
	commonStoreService   *commonStoreService
}

type Participant struct {
	Id   string
	Type entity.SenderType
}

func (s *MessageService) GetMessage(ctx context.Context, msgId *msProto.MessageId) (*msProto.Message, error) {
	apiKeyValid := commonModuleService.ApiKeyCheckerGRPC(ctx, s.postgresRepositories.CommonRepositories.AppKeyRepo, commonModuleService.MESSAGE_STORE_API)
	if !apiKeyValid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API Key")
	}

	tenantName, err := commonModuleService.GetTenantForUsernameForGRPC(ctx, s.postgresRepositories.CommonRepositories.UserRepo)
	if err != nil {
		return nil, err
	}

	if msgId == nil || msgId.GetConversationEventId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Message ID must be specified")
	}

	queryResult := s.postgresRepositories.ConversationEventRepository.GetEventById(msgId.GetConversationEventId())
	if queryResult.Error != nil {
		return nil, status.Errorf(codes.Internal, queryResult.Error.Error())
	}

	conversationEvent := *queryResult.Result.(*entity.ConversationEvent)

	conversationExists, err := s.customerOSService.ConversationByIdExists(*tenantName, conversationEvent.ConversationId)
	if err != nil {
		return nil, err
	}
	if !conversationExists {
		return nil, status.Errorf(codes.NotFound, "Conversation not found")
	}

	return s.commonStoreService.EncodeConversationEventToMS(conversationEvent), nil
}
func (s *MessageService) GetParticipants(ctx context.Context, feedId *msProto.FeedId) (*msProto.ParticipantsListResponse, error) {
	apiKeyValid := commonModuleService.ApiKeyCheckerGRPC(ctx, s.postgresRepositories.CommonRepositories.AppKeyRepo, commonModuleService.MESSAGE_STORE_API)
	if !apiKeyValid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API Key")
	}

	tenantName, err := commonModuleService.GetTenantForUsernameForGRPC(ctx, s.postgresRepositories.CommonRepositories.UserRepo)
	if err != nil {
		return nil, err
	}

	if feedId == nil || feedId.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Feed ID must be specified")
	}

	emails, err := s.customerOSService.GetConversationParticipants(*tenantName, feedId.GetId())
	var participants []string

	for _, participant := range emails {
		participants = append(participants, participant)
	}
	if err != nil {
		return nil, err
	}

	return &msProto.ParticipantsListResponse{
		Participants: participants,
	}, nil
}

func (s *MessageService) GetMessagesForFeed(ctx context.Context, feedIdRequest *msProto.FeedId) (*msProto.MessageListResponse, error) {
	apiKeyValid := commonModuleService.ApiKeyCheckerGRPC(ctx, s.postgresRepositories.CommonRepositories.AppKeyRepo, commonModuleService.MESSAGE_STORE_API)
	if !apiKeyValid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API Key")
	}

	tenantName, err := commonModuleService.GetTenantForUsernameForGRPC(ctx, s.postgresRepositories.CommonRepositories.UserRepo)
	if err != nil {
		return nil, err
	}

	if feedIdRequest == nil || feedIdRequest.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Feed ID must be specified")
	}

	exists, err := s.customerOSService.ConversationByIdExists(*tenantName, feedIdRequest.GetId())
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Conversation not found")
	}

	queryResult := s.postgresRepositories.ConversationEventRepository.GetEventsForConversation(feedIdRequest.GetId())
	if queryResult.Error != nil {
		return nil, status.Errorf(codes.Internal, queryResult.Error.Error())
	}

	var messages []*msProto.Message

	for _, event := range *queryResult.Result.(*[]entity.ConversationEvent) {
		messages = append(messages, s.commonStoreService.EncodeConversationEventToMS(event))
	}

	return &msProto.MessageListResponse{
		Messages: messages,
	}, nil
}

func (s *MessageService) GetFeeds(ctx context.Context, request *msProto.GetFeedsPagedRequest) (*msProto.FeedItemPagedResponse, error) {
	apiKeyValid := commonModuleService.ApiKeyCheckerGRPC(ctx, s.postgresRepositories.CommonRepositories.AppKeyRepo, commonModuleService.MESSAGE_STORE_API)
	if !apiKeyValid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API Key")
	}

	tenantName, err := commonModuleService.GetTenantForUsernameForGRPC(ctx, s.postgresRepositories.CommonRepositories.UserRepo)
	if err != nil {
		return nil, err
	}

	conversations, err := s.customerOSService.GetConversations(*tenantName)
	if err != nil {
		return nil, err
	}

	fl := &msProto.FeedItemPagedResponse{FeedItems: make([]*msProto.FeedItem, len(conversations))}
	fl.TotalElements = int32(len(conversations))

	for i, conversation := range conversations {
		fl.FeedItems[i] = s.commonStoreService.EncodeConversationToMS(conversation)
	}

	return fl, nil
}

func (s *MessageService) GetFeed(ctx context.Context, feedIdRequest *msProto.FeedId) (*msProto.FeedItem, error) {
	apiKeyValid := commonModuleService.ApiKeyCheckerGRPC(ctx, s.postgresRepositories.CommonRepositories.AppKeyRepo, commonModuleService.MESSAGE_STORE_API)
	if !apiKeyValid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API Key")
	}

	tenantName, err := commonModuleService.GetTenantForUsernameForGRPC(ctx, s.postgresRepositories.CommonRepositories.UserRepo)
	if err != nil {
		return nil, err
	}

	if feedIdRequest == nil || feedIdRequest.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Feed ID must be specified")
	}

	conversation, err := s.customerOSService.GetConversationById(*tenantName, feedIdRequest.GetId())
	if err != nil {
		return nil, err
	}

	return s.commonStoreService.EncodeConversationToMS(*conversation), nil
}

func (s *MessageService) SaveMessage(ctx context.Context, input *msProto.InputMessage) (*msProto.MessageId, error) {
	if input.ConversationId == nil && input.InitiatorIdentifier == nil {
		return nil, errors.New("conversationId or email must be provided")
	}
	if input.Content == nil {
		return nil, errors.New("message must be provided")
	}
	participants := []Participant{}

	tenant := "openline" //TODO get tenant from context
	initiator, err := s.getParticipant(tenant, *input.InitiatorIdentifier)
	participants = append(participants, *initiator)
	if err != nil {
		return nil, err
	}

	var threadId = ""
	entityType := s.commonStoreService.ConvertMSTypeToEntityType(input.Type)
	if err := s.getThreadIdAndParticipantsFromMail(tenant, &threadId, &participants, entityType, input); err != nil {
		log.Printf("Error handleing email: %v", err)
	}

	var conversation *Conversation
	if input.ConversationId != nil {
		if conv, err := s.customerOSService.GetConversationById(tenant, *input.ConversationId); err != nil {
			return nil, err
		} else {
			conversation = conv
		}
	} else {
		if conv, err := s.customerOSService.GetActiveConversationOrCreate(tenant, initiator.Id, *input.InitiatorIdentifier, input.SenderType, entityType, threadId); err != nil {
			return nil, err
		} else {
			conversation = conv
		}
	}

	previewMessage := s.getMessagePreview(input)

	userIds := []string{}
	contactIds := []string{}
	for participantsIndex := range participants {
		if participants[participantsIndex].Type == entity.CONTACT {
			contactIds = append(contactIds, participants[participantsIndex].Id)
		} else if participants[participantsIndex].Type == entity.USER {
			userIds = append(userIds, participants[participantsIndex].Id)
		}
	}

	lastSenderType := s.getSenderTypeStr(initiator)
	if _, err := s.customerOSService.UpdateConversation(tenant, conversation.Id, initiator.Id, lastSenderType, contactIds, userIds, previewMessage); err != nil {
		return nil, err
	}

	conversationEvent := s.saveConversationEvent(tenant, conversation, input, initiator)

	return s.commonStoreService.EncodeMessageIdToMs(conversationEvent), nil
}
func (s *MessageService) getSenderTypeStr(initiator *Participant) string {
	var lastSenderType = ""
	if initiator.Type == entity.CONTACT {
		lastSenderType = "CONTACT"
	} else if initiator.Type == entity.USER {
		lastSenderType = "USER"
	}
	return lastSenderType
}

func (s *MessageService) getThreadIdAndParticipantsFromMail(tenant string, threadId *string, participants *[]Participant, entityType entity.EventType, input *msProto.InputMessage) error {
	if entityType == entity.EMAIL {
		var messageJson EmailContent
		if err := json.Unmarshal([]byte(*input.Content), &messageJson); err != nil {
			return err
		}

		refSize := len(messageJson.Reference)
		if refSize > 0 {
			threadId = &messageJson.Reference[0]
		} else {
			threadId = &messageJson.MessageId
		}

		for _, toAddress := range append(messageJson.To, messageJson.Cc...) {
			if participant, err := s.getParticipant(tenant, toAddress); err != nil {
				log.Printf("Error getting participant: %v", err)
			} else {
				*participants = append(*participants, *participant)
			}
		}
	}
	return nil
}

func (s *MessageService) getMessagePreview(input *msProto.InputMessage) string {
	previewContent := ""
	if input.Content != nil {
		str := *input.Content
		msgLen := len(str)
		if msgLen > 20 {
			msgLen = 20
		}
		previewContent = str[0:msgLen]
	}
	return previewContent
}

func (s *MessageService) getParticipant(tenant string, initiatorIdentifier string) (*Participant, error) {
	user, err := s.customerOSService.GetUserByEmail(initiatorIdentifier)
	if err != nil {
		contact, err := s.customerOSService.GetContactWithEmailOrCreate(tenant, initiatorIdentifier)
		if err != nil {
			return nil, err
		}
		return &Participant{Id: contact.Id, Type: entity.CONTACT}, nil
	} else {
		return &Participant{Id: user.Id, Type: entity.USER}, nil
	}
}

func (s *MessageService) saveConversationEvent(tenant string, conversation *Conversation, input *msProto.InputMessage, initiator *Participant) entity.ConversationEvent {
	conversationEvent := entity.ConversationEvent{
		TenantName:     tenant,
		ConversationId: conversation.Id,
		Type:           s.commonStoreService.ConvertMSTypeToEntityType(input.Type),
		Subtype:        s.commonStoreService.ConvertMSSubtypeToEntitySubtype(input.Subtype),
		Content:        *input.Content,
		Source:         entity.OPENLINE,
		Direction:      s.commonStoreService.ConvertMSDirectionToEntityDirection(input.Direction),
		CreateDate:     time.Now(),

		InitiatorUsername: conversation.InitiatorUsername,

		SenderId:       initiator.Id,
		SenderUsername: *input.InitiatorIdentifier,
		SenderType:     initiator.Type,
		OriginalJson:   "TODO",
	}

	s.postgresRepositories.ConversationEventRepository.Save(&conversationEvent)
	return conversationEvent
}

func NewMessageService(driver *neo4j.Driver, postgresRepositories *repository.PostgresRepositories, customerOSService *CustomerOSService, commonStoreService *commonStoreService) *MessageService {
	ms := new(MessageService)
	ms.driver = driver
	ms.postgresRepositories = postgresRepositories
	ms.customerOSService = customerOSService
	ms.commonStoreService = commonStoreService
	return ms
}
