package service

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	commonModuleService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	msProto "github.com/openline-ai/openline-customer-os/packages/server/message-store-api/proto/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type MessageService struct {
	msProto.UnimplementedMessageStoreServiceServer
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
	customerOSService    *CustomerOSService
	commonStoreService   *commonStoreService
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

	tenant := "openline" //TODO get tenant from context
	initiatorId, err := s.getInitiatorId(input.SenderType, tenant, *input.InitiatorIdentifier)
	if err != nil {
		return nil, err
	}

	var conversation *Conversation
	if input.ConversationId != nil {
		if conv, err := s.customerOSService.GetConversationById(tenant, *input.ConversationId); err != nil {
			return nil, err
		} else {
			conversation = conv
		}
	} else {
		entityType := s.commonStoreService.ConvertMSTypeToEntityType(input.Type)
		if conv, err := s.customerOSService.GetActiveConversationOrCreate(tenant, *initiatorId, *input.InitiatorIdentifier, input.SenderType, entityType, *input.Content); err != nil {
			return nil, err
		} else {
			conversation = conv
		}
	}

	previewMessage := s.getMessagePreview(input)

	senderType := s.commonStoreService.ConvertMSSenderTypeToEntitySenderType(input.SenderType)
	if _, err := s.customerOSService.UpdateConversation(tenant, conversation.Id, *initiatorId, senderType, previewMessage); err != nil {
		return nil, err
	}

	conversationEvent := s.saveConversationEvent(tenant, conversation, input, *initiatorId, senderType)

	return s.commonStoreService.EncodeMessageIdToMs(conversationEvent), nil
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

func (s *MessageService) getInitiatorId(st msProto.SenderType, tenant string, email string) (*string, error) {
	var initiatorId *string
	if st == msProto.SenderType_CONTACT {
		if contact, err := s.customerOSService.GetContactWithEmailOrCreate(tenant, email); err != nil {
			return nil, err
		} else {
			initiatorId = &contact.Id
		}
	} else if st == msProto.SenderType_USER {
		if user, err := s.customerOSService.GetUserByEmail(email); err != nil {
			return nil, err
		} else {
			initiatorId = &user.Id
		}
	}

	return initiatorId, nil
}

func (s *MessageService) saveConversationEvent(tenant string, conversation *Conversation, input *msProto.InputMessage, initiatorId string, senderType entity.SenderType) entity.ConversationEvent {
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

		SenderId:       initiatorId,
		SenderUsername: *input.InitiatorIdentifier,
		SenderType:     senderType,
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
