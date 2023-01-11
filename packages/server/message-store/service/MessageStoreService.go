package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	pb "github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository"
)

type messageService struct {
	proto.UnimplementedMessageStoreServiceServer
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
}

//func encodeConversationState(feedState pb.FeedItemState) genConversation.State {
//	switch feedState {
//	case pb.FeedItemState_NEW:
//		return genConversation.StateNEW
//	case pb.FeedItemState_IN_PROGRESS:
//		return genConversation.StateIN_PROGRESS
//	case pb.FeedItemState_CLOSED:
//		return genConversation.StateCLOSED
//	default:
//		return genConversation.StateNEW
//	}
//}
//
//func decodeConversationState(feedState genConversation.State) pb.FeedItemState {
//	switch feedState {
//	case genConversation.StateNEW:
//		return pb.FeedItemState_NEW
//	case genConversation.StateIN_PROGRESS:
//		return pb.FeedItemState_IN_PROGRESS
//	case genConversation.StateCLOSED:
//		return pb.FeedItemState_CLOSED
//
//	default:
//		return pb.FeedItemState_NEW
//	}
//}
//
//func decodeSenderType(feedState genConversation.LastSenderType) pb.SenderType {
//	switch feedState {
//	case genConversation.LastSenderTypeCONTACT:
//		return pb.SenderType_CONTACT
//	case genConversation.LastSenderTypeUSER:
//		return pb.SenderType_USER
//	default:
//		return pb.SenderType_CONTACT
//	}
//}
//
//func encodeConversationEventType(channel pb.MessageChannel) entity.EventType {
//	switch channel {
//	case pb.MessageChannel_WIDGET:
//		return entity.WEB_CHAT
//	case pb.MessageChannel_MAIL:
//		return entity.EMAIL
//	case pb.MessageChannel_VOICE:
//		return entity.VOICE
//	default:
//		return entity.WEB_CHAT
//	}
//}
//
//func encodeChannel(channel pb.MessageChannel) conversationitem.Channel {
//	switch channel {
//	case pb.MessageChannel_WIDGET:
//		return conversationitem.ChannelCHAT
//	case pb.MessageChannel_MAIL:
//		return conversationitem.ChannelMAIL
//	case pb.MessageChannel_VOICE:
//		return conversationitem.ChannelVOICE
//	default:
//		return conversationitem.ChannelCHAT
//	}
//}
//
//func encodeDirection(direction pb.MessageDirection) conversationitem.Direction {
//	switch direction {
//	case pb.MessageDirection_INBOUND:
//		return conversationitem.DirectionINBOUND
//	case pb.MessageDirection_OUTBOUND:
//		return conversationitem.DirectionOUTBOUND
//	default:
//		return conversationitem.DirectionOUTBOUND
//	}
//}
//
//func encodeConversationEventDirection(direction pb.MessageDirection) entity.Direction {
//	switch direction {
//	case pb.MessageDirection_INBOUND:
//		return entity.INBOUND
//	case pb.MessageDirection_OUTBOUND:
//		return entity.OUTBOUND
//	default:
//		return entity.OUTBOUND
//	}
//}
//
//func encodeType(t pb.MessageType) conversationitem.Type {
//	switch t {
//	case pb.MessageType_MESSAGE:
//		return conversationitem.TypeMESSAGE
//	case pb.MessageType_FILE:
//		return conversationitem.TypeFILE
//	default:
//		return conversationitem.TypeMESSAGE
//	}
//}
//
//func decodeType(t conversationitem.Type) pb.MessageType {
//	switch t {
//	case conversationitem.TypeMESSAGE:
//		return pb.MessageType_MESSAGE
//	case conversationitem.TypeFILE:
//		return pb.MessageType_FILE
//	default:
//		return pb.MessageType_MESSAGE
//	}
//}
//
//func decodeDirection(direction conversationitem.Direction) pb.MessageDirection {
//	switch direction {
//	case conversationitem.DirectionINBOUND:
//		return pb.MessageDirection_INBOUND
//	case conversationitem.DirectionOUTBOUND:
//		return pb.MessageDirection_OUTBOUND
//	default:
//		return pb.MessageDirection_OUTBOUND
//	}
//}
//
//func decodeChannel(channel conversationitem.Channel) pb.MessageChannel {
//	switch channel {
//	case conversationitem.ChannelCHAT:
//		return pb.MessageChannel_WIDGET
//	case conversationitem.ChannelMAIL:
//		return pb.MessageChannel_MAIL
//	case conversationitem.ChannelVOICE:
//		return pb.MessageChannel_VOICE
//	default:
//		return pb.MessageChannel_WIDGET
//	}
//}

func (s *messageService) GetMessage(ctx context.Context, msgId *pb.Id) (*pb.Message, error) {
	//if msgId == nil || msgId.GetId() == 0 {
	//	return nil, status.Errorf(codes.InvalidArgument, "Message ID must be specified")
	//}
	//
	//mi, err := s.client.ConversationItem.Get(ctx, int(msgId.GetId()))
	//if err != nil {
	//	se, _ := status.FromError(err)
	//	return nil, status.Errorf(se.Code(), "Error finding Message")
	//}
	//
	//mf, err := s.client.ConversationItem.QueryConversation(mi).First(ctx)
	//if err != nil {
	//	se, _ := status.FromError(err)
	//	return nil, status.Errorf(se.Code(), "Error finding Feed")
	//}
	//
	//messageId := int64(mi.ID)
	//conversationid := int64(mf.ID)
	//
	//contactById, err := getContactById(s.graphqlClient, mf.ContactId, s.config.Service.CustomerOsAPIKey)
	//
	//m := &pb.Message{
	//	Type:      decodeType(mi.Type),
	//	Message:   mi.Message,
	//	Direction: decodeDirection(mi.Direction),
	//	Channel:   decodeChannel(mi.Channel),
	//	Username:  contactById.email,
	//	ContactId: &mf.ContactId,
	//	Id:        &messageId,
	//	FeedId:    &conversationid,
	//	Time:      timestamppb.New(mi.Time),
	//}
	//
	//if mi.Direction == conversationitem.DirectionOUTBOUND {
	//	m.UserId = &mi.SenderId
	//}

	return nil, nil
}

func (s *messageService) GetMessages(ctx context.Context, messagesRequest *pb.GetMessagesRequest) (*pb.MessagePagedResponse, error) {
	//var messages []*gen.ConversationItem
	//var err error
	//var conversation *gen.Conversation
	//
	//if messagesRequest != nil {
	//	log.Printf("Looking up messages for conversation with id %d", messagesRequest.GetConversationId())
	//	conversation, err = s.client.Conversation.Get(ctx, int(messagesRequest.GetConversationId()))
	//	if err != nil {
	//		se, _ := status.FromError(err)
	//		return nil, status.Errorf(se.Code(), "Error finding conversation with id  %d", messagesRequest.GetConversationId())
	//	}
	//} else {
	//	log.Printf("Conversation id is required")
	//	return nil, status.Errorf(1, "Conversation id is required")
	//}
	//
	//limit := 100 // default to 100 if no pagination is specified
	//if messagesRequest.GetPageSize() != 0 {
	//	limit = int(messagesRequest.GetPageSize())
	//}
	//
	//if messagesRequest.GetBefore() == nil {
	//	messages, err = s.client.Conversation.QueryConversationItem(conversation).
	//		Order(gen.Desc(conversationitem.FieldTime)).
	//		Limit(limit).
	//		All(ctx)
	//} else {
	//	messages, err = s.client.Conversation.QueryConversationItem(conversation).
	//		Order(gen.Desc(conversationitem.FieldTime)).
	//		Where(conversationitem.TimeLT(messagesRequest.GetBefore().AsTime())).
	//		Limit(limit).
	//		All(ctx)
	//}
	//
	//if err != nil {
	//	se, _ := status.FromError(err)
	//	return nil, status.Errorf(se.Code(), "Error getting messages")
	//}
	//ml := &pb.MessagePagedResponse{Message: make([]*pb.Message, len(messages))}
	//
	//conversationId := int64(conversation.ID)
	//
	//for i, j := len(messages)-1, 0; i >= 0; i, j = i-1, j+1 {
	//	var mid = int64(messages[i].ID)
	//	mi := &pb.Message{
	//		Type:      decodeType(messages[i].Type),
	//		Message:   messages[i].Message,
	//		Direction: decodeDirection(messages[i].Direction),
	//		Channel:   decodeChannel(messages[i].Channel),
	//		Username:  nil,
	//		ContactId: &conversation.ContactId,
	//		Id:        &mid,
	//		FeedId:    &conversationId,
	//		Time:      timestamppb.New(messages[i].Time),
	//	}
	//	ml.Message[j] = mi
	//}
	return nil, nil
}

func (s *messageService) GetFeeds(ctx context.Context, feedRequest *pb.GetFeedsPagedRequest) (*pb.FeedItemPagedResponse, error) {
	//query := s.client.Conversation.Query()
	//
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
	//
	//conversations, err := query.Limit(limit).Offset(offset).All(ctx)
	//count, err2 := query.Count(ctx)
	//
	//if err != nil || err2 != nil {
	//	se, _ := status.FromError(err)
	//	return nil, status.Errorf(se.Code(), "Error getting messages: %s", err.Error())
	//}
	//fl := &pb.FeedItemPagedResponse{FeedItems: make([]*pb.FeedItem, len(conversations))}
	//fl.TotalElements = int32(count)
	//
	//for i, conversation := range conversations {
	//	var id = int64(conversation.ID)
	//	log.Printf("Got a conversation id of %d", id)
	//
	//	contactById, err := getContactById(s.graphqlClient, conversation.ContactId, s.config.Service.CustomerOsAPIKey)
	//
	//	if err != nil {
	//		se, _ := status.FromError(err)
	//		return nil, status.Errorf(se.Code(), "Error getting messages: %s", err.Error())
	//	}
	//
	//	fl.FeedItems[i] = &pb.FeedItem{
	//		Id:               int64(conversation.ID),
	//		ContactId:        contactById.id,
	//		ContactFirstName: contactById.firstName,
	//		ContactLastName:  contactById.lastName,
	//		ContactEmail:     *contactById.email,
	//		State:            decodeConversationState(conversation.State),
	//		LastSenderId:     conversation.LastSenderId,
	//		LastSenderType:   decodeSenderType(conversation.LastSenderType),
	//		Message:          conversation.LastMessage,
	//		UpdatedOn:        timestamppb.New(conversation.UpdatedOn),
	//	}
	//
	//	msg, _ := json.Marshal(fl.FeedItems[i])
	//	log.Printf("Got a feed item of %s", msg)
	//}
	return nil, nil
}

func (s *messageService) GetFeed(ctx context.Context, feedIdRequest *pb.Id) (*pb.FeedItem, error) {
	//if feedIdRequest == nil || feedIdRequest.GetId() == 0 {
	//	return nil, status.Errorf(codes.InvalidArgument, "Feed ID must be specified")
	//}
	//
	//conversation, err := s.client.Conversation.Get(ctx, int(feedIdRequest.GetId()))
	//if err != nil {
	//	se, _ := status.FromError(err)
	//	return nil, status.Errorf(se.Code(), "Error finding conversation")
	//}
	//
	//contactById, err := getContactById(s.graphqlClient, conversation.ContactId, s.config.Service.CustomerOsAPIKey)
	//if err != nil {
	//	se, _ := status.FromError(err)
	//	return nil, status.Errorf(se.Code(), "Error getting messages: %s", err.Error())
	//}
	//
	//return &pb.FeedItem{
	//	Id:               int64(conversation.ID),
	//	ContactId:        contactById.id,
	//	ContactFirstName: contactById.firstName,
	//	ContactLastName:  contactById.lastName,
	//	ContactEmail:     *contactById.email,
	//	State:            decodeConversationState(conversation.State),
	//	LastSenderId:     conversation.LastSenderId,
	//	LastSenderType:   decodeSenderType(conversation.LastSenderType),
	//	Message:          conversation.LastMessage,
	//	UpdatedOn:        timestamppb.New(conversation.UpdatedOn),
	//}, nil
	return nil, nil
}

func NewMessageService(driver *neo4j.Driver, postgresRepositories *repository.PostgresRepositories) *messageService {
	ms := new(messageService)
	ms.driver = driver
	ms.postgresRepositories = postgresRepositories
	return ms
}
