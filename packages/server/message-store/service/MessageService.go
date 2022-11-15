package service

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/message-store/config"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/messagefeed"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/messageitem"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	pb "github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"regexp"
	"strings"
	time "time"
)

type messageService struct {
	proto.UnimplementedMessageStoreServiceServer
	client        *gen.Client
	graphqlClient *graphql.Client
	config        *c.Config
}

type ContactInfo struct {
	firstName string
	lastName  string
	id        string
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

func parseEmail(email string) (string, string) {
	re := regexp.MustCompile("^\"{0,1}([^\"]*)\"{0,1}[ ]*<(.*)>$")
	matches := re.FindStringSubmatch(strings.Trim(email, " "))
	if matches != nil {
		return strings.Trim(matches[1], " "), matches[2]
	}
	return "", email
}

func getContactByEmail(graphqlClient *graphql.Client, email string) (*ContactInfo, error) {

	graphqlRequest := graphql.NewRequest(fmt.Sprintf(`
  				query {
  					contactByEmail(email: "%s"){firstName,lastName,id}
  				}
    `, email))

	var graphqlResponse map[string]map[string]string
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return &ContactInfo{firstName: graphqlResponse["contactByEmail"]["firstName"],
		lastName: graphqlResponse["contactByEmail"]["lastName"],
		id:       graphqlResponse["contactByEmail"]["id"]}, nil
}

func createContact(graphqlClient *graphql.Client, firstName string, lastName string, email string) (string, error) {
	graphqlRequest := graphql.NewRequest(fmt.Sprintf(`
		mutation CreateContact {
		  createContact(input: {firstName: "%s",
			lastName: "%s",
		  email:{email:  "%s", label: WORK}}) {
			id
		  }
		}
    `, firstName, lastName, email))

	var graphqlResponse map[string]map[string]string
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}
	return graphqlResponse["createContact"]["id"], nil
}

func createConversation(graphqlClient *graphql.Client, userId string, contactId string) (string, error) {
	graphqlRequest := graphql.NewRequest(fmt.Sprintf(`
			mutation CreateConversation {
				createConversation(input: {
					userId: "%s"
					contactId: "%s"
				}) {
					id
					startedAt
				}
			}
    `, userId, contactId))

	var graphqlResponse map[string]map[string]string
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}
	return graphqlResponse["createConversation"]["id"], nil
}

func (s *messageService) SaveMessage(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	var contact *ContactInfo
	var err error
	var feed *gen.MessageFeed

	if message.Contact == nil {
		if message.GetChannel() == pb.MessageChannel_MAIL || message.GetChannel() == pb.MessageChannel_WIDGET {
			displayName, email := parseEmail(message.Username)
			contact, err = getContactByEmail(s.graphqlClient, email)

			if err != nil {

				log.Printf("Contact %s creating a new contact", email)
				firstName, lastName := "Unknown", "User"
				if displayName != "" {
					parts := strings.SplitN(displayName, " ", 2)
					firstName = parts[0]
					if len(parts) > 1 {
						lastName = parts[1]
					}
				}
				log.Printf("Making a contact, firstName=%s lastName=%s email=%s", firstName, lastName, email)
				contactId, err := createContact(s.graphqlClient, firstName, lastName, email)
				contact = &ContactInfo{
					firstName: firstName,
					lastName:  lastName,
					id:        contactId,
				}
				if err != nil {
					log.Printf("Unable to create contact! %v", err)
					return nil, err
				}
			}
			feed, err = s.client.MessageFeed.
				Query().
				Where(messagefeed.ContactId(contact.id)).
				First(ctx)

		} else {
			return nil, status.Errorf(codes.Unimplemented, "Contact mapping not implemented yet for %v", message.GetChannel())
		}
	} else {
		feed, err = s.client.MessageFeed.
			Get(ctx, int(*message.Contact.Id))
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "Couldn't find a feed for id of %d", *message.Contact.Id)
		}

		contact = &ContactInfo{firstName: feed.FirstName,
			lastName: feed.LastName,
			id:       feed.ContactId,
		}
	}

	if err != nil {
		// can only reach here if message.Contact is nil & the contactId found in neo4j doesn't match a message any feed
		se, _ := status.FromError(err)
		if se.Code() != codes.Unknown {
			return nil, status.Errorf(se.Code(), "Error upserting Feed")
		} else {
			feed, err = s.client.MessageFeed.
				Create().
				SetFirstName(contact.firstName).
				SetLastName(contact.lastName).
				SetContactId(contact.id).
				Save(ctx)
			if err != nil {
				se, _ = status.FromError(err)
				return nil, status.Errorf(se.Code(), "Error inserting new Feed")
			}

			conv, err := createConversation(s.graphqlClient, s.config.Identity.DefaultUserId, contact.id)

			if err != nil {
				log.Printf("Error making conversation %v", err)
				return nil, err
			}

			log.Printf("Created conversation %s", conv)

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
	var feedId int64 = int64(feed.ID)
	mi := &pb.Message{
		Type:      decodeType(msg.Type),
		Message:   msg.Message,
		Direction: decodeDirection(msg.Direction),
		Channel:   decodeChannel(msg.Channel),
		Username:  msg.Username,
		Id:        &id,
		Contact:   &pb.Contact{ContactId: contact.id, Id: &feedId, FirstName: contact.firstName, LastName: contact.lastName},
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

	mf, err := s.client.MessageItem.QueryMessageFeed(mi).First(ctx)
	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error finding Feed")
	}
	m := &pb.Message{
		Type:      decodeType(mi.Type),
		Message:   mi.Message,
		Direction: decodeDirection(mi.Direction),
		Channel:   decodeChannel(mi.Channel),
		Username:  mi.Username,
		Id:        msg.Id,
		Time:      timestamppb.New(mi.Time),
		Contact:   &pb.Contact{ContactId: mf.ContactId, FirstName: mf.FirstName, LastName: mf.LastName},
	}
	return m, nil
}

func (s *messageService) GetMessages(ctx context.Context, pc *pb.PagedContact) (*pb.MessageList, error) {
	var messages []*gen.MessageItem
	var err error
	var mf *gen.MessageFeed
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
		log.Printf("Looking up messages for Contact name %s", contact.GetContactId())
		mf, err = s.client.MessageFeed.Query().
			Where(messagefeed.ContactId(contact.GetContactId())).
			First(ctx)
		if err != nil {
			se, _ := status.FromError(err)
			return nil, status.Errorf(se.Code(), "Error finding Feed")
		}
	}

	if pageInfo.Before == nil {
		messages, err = s.client.MessageFeed.QueryMessageItem(mf).
			Order(gen.Desc(messageitem.FieldTime)).
			Limit(int(pageInfo.PageSize)).
			All(ctx)
	} else {
		messages, err = s.client.MessageFeed.QueryMessageItem(mf).
			Order(gen.Desc(messageitem.FieldTime)).
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
			Contact:   &pb.Contact{ContactId: contact.ContactId},
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
		fl.Contact[i] = &pb.Contact{ContactId: contact.ContactId, FirstName: contact.FirstName, LastName: contact.LastName, Id: &id}
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
		return &pb.Contact{FirstName: mf.FirstName, LastName: mf.LastName, ContactId: mf.ContactId, Id: &id}, nil
	} else {
		log.Printf("Looking up messages for Contact name %s", contact.GetContactId())
		mf, err := s.client.MessageFeed.Query().
			Where(messagefeed.ContactId(contact.GetContactId())).
			First(ctx)
		if err != nil {
			se, _ := status.FromError(err)
			return nil, status.Errorf(se.Code(), "Error finding Feed")
		}
		var id int64 = int64(mf.ID)
		return &pb.Contact{FirstName: mf.FirstName, LastName: mf.LastName, ContactId: mf.ContactId, Id: &id}, nil
	}
}

func NewMessageService(client *gen.Client, graphqlClient *graphql.Client, config *c.Config) *messageService {
	ms := new(messageService)
	ms.client = client
	ms.graphqlClient = graphqlClient
	ms.config = config
	return ms
}
