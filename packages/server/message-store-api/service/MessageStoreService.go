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

type messageService struct {
	msProto.UnimplementedMessageStoreServiceServer
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
	customerOSService    *customerOSService
	commonStoreService   *commonStoreService
}

func (s *messageService) GetMessage(ctx context.Context, msgId *msProto.MessageId) (*msProto.Message, error) {
	apiKeyValid := commonModuleService.ApiKeyCheckerGRPC(ctx, s.postgresRepositories.CommonRepositories.AppKeyRepo, commonModuleService.MESSAGE_STORE_API)
	if !apiKeyValid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API Key")
	}

	tenantName, err := commonModuleService.GetTenantForUsernameForGRPC(ctx, s.postgresRepositories.CommonRepositories.UserRepo)
	if err != nil {
		return nil, err
	}

	if msgId == nil || msgId.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Message ID must be specified")
	}

	queryResult := s.postgresRepositories.ConversationEventRepository.GetEventById(msgId.GetId())
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

func (s *messageService) GetMessagesForFeed(ctx context.Context, feedIdRequest *msProto.FeedId) (*msProto.MessageListResponse, error) {
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

func (s *messageService) GetFeeds(ctx context.Context, request *msProto.GetFeedsPagedRequest) (*msProto.FeedItemPagedResponse, error) {
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

	//if feedRequest.GetStateIn() != nil {
	//	stateIn := make([]genConversation.State, 0, len(feedRequest.GetStateIn()))
	//	for _, state := range feedRequest.GetStateIn() {
	//		stateIn = append(stateIn, encodeConversationState(state))
	//	}
	//	query.Where(genConversation.StateIn(stateIn...))
	//}
	//
	//limit := 100 // default to 100 if no pagination is specified
	//if feedRequest.GetPageSize() != 0 {
	//	limit = int(feedRequest.GetPageSize())
	//}
	//offset := limit * int(feedRequest.GetPage())

	fl := &msProto.FeedItemPagedResponse{FeedItems: make([]*msProto.FeedItem, len(conversations))}
	fl.TotalElements = int32(len(conversations))

	for i, conversation := range conversations {
		fl.FeedItems[i] = s.commonStoreService.EncodeConversationToMS(conversation)
	}

	return fl, nil
}

func (s *messageService) GetFeed(ctx context.Context, feedIdRequest *msProto.FeedId) (*msProto.FeedItem, error) {
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

func (s *messageService) SaveMessage(ctx context.Context, input *msProto.InputMessage) (*msProto.MessageId, error) {
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
		conv, err := s.customerOSService.GetConversationById(tenant, *input.ConversationId)
		if err != nil {
			return nil, err
		}
		conversation = conv
	}

	if input.SenderType == msProto.SenderType_CONTACT {
		contact, err := s.customerOSService.GetContactWithEmailOrCreate(tenant, *input.Email)
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
		entityType := s.commonStoreService.ConvertMSTypeToEntityType(input.Type)
		conv, err := s.customerOSService.GetActiveConversationOrCreate(tenant, participantId, participantFirstName, participantLastName, *input.Email, input.SenderType, entityType)
		if err != nil {
			return nil, err
		}
		conversation = conv
	}

	conversationEvent := entity.ConversationEvent{
		TenantName:     tenant,
		ConversationId: conversation.Id,
		Type:           s.commonStoreService.ConvertMSTypeToEntityType(input.Type),
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

	return s.commonStoreService.EncodeMessageIdToMs(conversationEvent), nil
}

func NewMessageService(driver *neo4j.Driver, postgresRepositories *repository.PostgresRepositories, customerOSService *customerOSService, commonStoreService *commonStoreService) *messageService {
	ms := new(messageService)
	ms.driver = driver
	ms.postgresRepositories = postgresRepositories
	ms.customerOSService = customerOSService
	ms.commonStoreService = commonStoreService
	return ms
}
