package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/ent"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/ent/messagefeed"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/ent/messageitem"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/ent/proto"
	pb "github.com/openline-ai/openline-customer-os/packages/server/message-store/ent/proto"
	time "time"
)

type messageService struct {
	proto.UnimplementedMessageStoreServiceServer
	client *ent.Client
}

func encodeChannel(channel pb.MessageChannel) messageitem.Channel {
	switch channel {
	case pb.MessageChannel_WIDGET:
		return messageitem.ChannelCHAT
	case pb.MessageChannel_MAIL:
		return messageitem.ChannelMAIL
	case pb.MessageChannel_WHATSAPP:
		return messageitem.ChannelWHATSAPP
	case pb.MessageChannel_FACEBOOK:
		return messageitem.ChannelFACEBOOK
	case pb.MessageChannel_TWITTER:
		return messageitem.ChannelTWITTER
	case pb.MessageChannel_VOICE:
		return messageitem.ChannelVOICE
	default:
		return messageitem.ChannelCHAT
	}
}

func encodeDirection(direction pb.MessageDirection) messageitem.Direction {
	switch direction {
	case pb.MessageDirection_INBOUND:
		return messageitem.DirectionINBOUND
	case pb.MessageDirection_OUTBOUND:
		return messageitem.DirectionOUTBOUND
	default:
		return messageitem.DirectionOUTBOUND
	}
}

func encodeType(t pb.MessageType) messageitem.Type {
	switch t {
	case pb.MessageType_MESSAGE:
		return messageitem.TypeMESSAGE
	case pb.MessageType_FILE:
		return messageitem.TypeFILE
	default:
		return messageitem.TypeMESSAGE
	}
}

func decodeType(t messageitem.Type) pb.MessageType {
	switch t {
	case messageitem.TypeMESSAGE:
		return pb.MessageType_MESSAGE
	case messageitem.TypeFILE:
		return pb.MessageType_FILE
	default:
		return pb.MessageType_MESSAGE
	}
}

func decodeDirection(direction messageitem.Direction) pb.MessageDirection {
	switch direction {
	case messageitem.DirectionINBOUND:
		return pb.MessageDirection_INBOUND
	case messageitem.DirectionOUTBOUND:
		return pb.MessageDirection_OUTBOUND
	default:
		return pb.MessageDirection_OUTBOUND
	}
}

func decodeChannel(channel messageitem.Channel) pb.MessageChannel {
	switch channel {
	case messageitem.ChannelCHAT:
		return pb.MessageChannel_WIDGET
	case messageitem.ChannelMAIL:
		return pb.MessageChannel_MAIL
	case messageitem.ChannelWHATSAPP:
		return pb.MessageChannel_WHATSAPP
	case messageitem.ChannelFACEBOOK:
		return pb.MessageChannel_FACEBOOK
	case messageitem.ChannelTWITTER:
		return pb.MessageChannel_TWITTER
	case messageitem.ChannelVOICE:
		return pb.MessageChannel_VOICE
	default:
		return pb.MessageChannel_WIDGET
	}
}
func (s *messageService) SaveMessage(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	var contact string
	if message.Contact == nil {
		contact = message.Username // TODO: resolve address to contact
	} else {
		contact = message.Contact.Username
	}

	feed, err := s.client.MessageFeed.
		Query().
		Where(messagefeed.Username(contact)).
		First(ctx)

	if err != nil {
		se, _ := status.FromError(err)
		if se.Code() != codes.Unknown {
			return nil, status.Errorf(se.Code(), "Error upserting Feed")
		} else {
			feed, err = s.client.MessageFeed.
				Create().
				SetUsername(contact).
				Save(ctx)
			if err != nil {
				se, _ = status.FromError(err)
				return nil, status.Errorf(se.Code(), "Error inserting new Feed")
			}

		}
	}

	var t *time.Time = nil
	if message.GetTime() != nil {
		var timeref = message.GetTime().AsTime()
		t = &timeref
	}
	msg, err := s.client.MessageItem.
		Create().
		SetMessage(message.GetMessage()).
		SetMessageFeed(feed).
		SetChannel(encodeChannel(message.GetChannel())).
		SetNillableTime(t).
		SetUsername(message.GetUsername()).
		SetDirection(encodeDirection(message.GetDirection())).
		SetType(encodeType(message.GetType())).
		Save(ctx)

	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error inserting message")
	}

	if t == nil {
		var timeRef = time.Now()
		t = &timeRef
	}

	var id int64 = int64(msg.ID)
	mi := &pb.Message{
		Type:      decodeType(msg.Type),
		Message:   msg.Message,
		Direction: decodeDirection(msg.Direction),
		Channel:   decodeChannel(msg.Channel),
		Username:  msg.Username,
		Id:        &id,
		Contact:   &pb.Contact{Username: contact},
		Time:      timestamppb.New(*t),
	}
	return mi, nil
}

func (s *messageService) GetMessage(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	if msg.Id == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Message ID must be specified")
	}

	mi, err := s.client.MessageItem.Get(ctx, int(msg.GetId()))
	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error finding Message")
	}

	m := &pb.Message{
		Type:      decodeType(mi.Type),
		Message:   mi.Message,
		Direction: decodeDirection(mi.Direction),
		Channel:   decodeChannel(mi.Channel),
		Username:  mi.Username,
		Id:        msg.Id,
		Time:      timestamppb.New(mi.Time),
		Contact:   &pb.Contact{Username: mi.Username},
	}
	return m, nil
}

func (s *messageService) GetMessages(ctx context.Context, pc *pb.PagedContact) (*pb.MessageList, error) {
	var messages []*ent.MessageItem
	var err error
	var mf *ent.MessageFeed
	contact := pc.Contact
	pageInfo := pc.Page

	if contact.Id != nil {
		log.Printf("Looking up messages for Contact id %d", *contact.Id)
		mf, err = s.client.MessageFeed.Get(ctx, int(*contact.Id))
		if err != nil {
			se, _ := status.FromError(err)
			return nil, status.Errorf(se.Code(), "Error finding Feed")
		}
	} else {
		log.Printf("Looking up messages for Contact name %s", contact.GetUsername())
		mf, err = s.client.MessageFeed.Query().
			Where(messagefeed.Username(contact.GetUsername())).
			First(ctx)
		if err != nil {
			se, _ := status.FromError(err)
			return nil, status.Errorf(se.Code(), "Error finding Feed")
		}
	}

	if pageInfo.Before == nil {
		messages, err = s.client.MessageFeed.QueryMessageItem(mf).
			Order(ent.Desc(messageitem.FieldTime)).
			Limit(int(pageInfo.PageSize)).
			All(ctx)
	} else {
		messages, err = s.client.MessageFeed.QueryMessageItem(mf).
			Order(ent.Desc(messageitem.FieldTime)).
			Where(messageitem.TimeLT(pageInfo.Before.AsTime())).
			Limit(int(pageInfo.PageSize)).
			All(ctx)
	}

	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error getting messages")
	}
	ml := &pb.MessageList{Message: make([]*pb.Message, len(messages))}

	for i, j := len(messages)-1, 0; i >= 0; i, j = i-1, j+1 {
		message := messages[i]
		var id int64 = int64(message.ID)
		mi := &pb.Message{
			Type:      decodeType(message.Type),
			Message:   message.Message,
			Direction: decodeDirection(message.Direction),
			Channel:   decodeChannel(message.Channel),
			Username:  message.Username,
			Id:        &id,
			Time:      timestamppb.New(message.Time),
			Contact:   &pb.Contact{Username: contact.Username},
		}
		ml.Message[j] = mi
	}
	return ml, nil
}

func (s *messageService) GetFeeds(ctx context.Context, _ *pb.Empty) (*pb.FeedList, error) {
	contacts, err := s.client.MessageFeed.Query().All(ctx)
	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error getting messages")
	}
	fl := &pb.FeedList{Contact: make([]*pb.Contact, len(contacts))}

	for i, contact := range contacts {
		var id int64 = int64(contact.ID)
		log.Printf("Got an feed id of %d", id)
		fl.Contact[i] = &pb.Contact{Username: contact.Username, Id: &id}
	}
	return fl, nil
}
func (s *messageService) GetFeed(ctx context.Context, contact *pb.Contact) (*pb.Contact, error) {
	if contact.Id != nil {
		log.Printf("Looking up messages for Contact id %d", *contact.Id)
		mf, err := s.client.MessageFeed.Get(ctx, int(*contact.Id))
		if err != nil {
			se, _ := status.FromError(err)
			return nil, status.Errorf(se.Code(), "Error finding Feed")
		}
		var id int64 = int64(mf.ID)
		return &pb.Contact{Username: mf.Username, Id: &id}, nil
	} else {
		log.Printf("Looking up messages for Contact name %s", contact.GetUsername())
		mf, err := s.client.MessageFeed.Query().
			Where(messagefeed.Username(contact.GetUsername())).
			First(ctx)
		if err != nil {
			se, _ := status.FromError(err)
			return nil, status.Errorf(se.Code(), "Error finding Feed")
		}
		var id int64 = int64(mf.ID)
		return &pb.Contact{Username: mf.Username, Id: &id}, nil
	}
}

func NewMessageService(client *ent.Client) *messageService {
	ms := new(messageService)
	ms.client = client
	return ms
}
