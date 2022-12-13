package service

import (
	"context"
	"encoding/json"
	"github.com/machinebox/graphql"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	c "github.com/openline-ai/openline-customer-os/packages/server/message-store/config"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen"
	genConversation "github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/conversation"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/conversationitem"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	pb "github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"regexp"
	"strconv"
	"strings"
	time "time"
)

type messageService struct {
	proto.UnimplementedMessageStoreServiceServer
	client        *gen.Client
	graphqlClient *graphql.Client
	config        *c.Config
	driver        *neo4j.Driver
}

type ContactInfo struct {
	firstName string
	lastName  string
	email     *string
	phone     *string
	id        string
}

func encodeConversationState(feedState pb.FeedItemState) genConversation.State {
	switch feedState {
	case pb.FeedItemState_NEW:
		return genConversation.StateNEW
	case pb.FeedItemState_IN_PROGRESS:
		return genConversation.StateIN_PROGRESS
	case pb.FeedItemState_CLOSED:
		return genConversation.StateCLOSED
	default:
		return genConversation.StateNEW
	}
}

func decodeConversationState(feedState genConversation.State) pb.FeedItemState {
	switch feedState {
	case genConversation.StateNEW:
		return pb.FeedItemState_NEW
	case genConversation.StateIN_PROGRESS:
		return pb.FeedItemState_IN_PROGRESS
	case genConversation.StateCLOSED:
		return pb.FeedItemState_CLOSED

	default:
		return pb.FeedItemState_NEW
	}
}

func decodeSenderType(feedState genConversation.LastSenderType) pb.SenderType {
	switch feedState {
	case genConversation.LastSenderTypeCONTACT:
		return pb.SenderType_CONTACT
	case genConversation.LastSenderTypeUSER:
		return pb.SenderType_USER
	default:
		return pb.SenderType_CONTACT
	}
}

func encodeChannel(channel pb.MessageChannel) conversationitem.Channel {
	switch channel {
	case pb.MessageChannel_WIDGET:
		return conversationitem.ChannelCHAT
	case pb.MessageChannel_MAIL:
		return conversationitem.ChannelMAIL
	case pb.MessageChannel_WHATSAPP:
		return conversationitem.ChannelWHATSAPP
	case pb.MessageChannel_FACEBOOK:
		return conversationitem.ChannelFACEBOOK
	case pb.MessageChannel_TWITTER:
		return conversationitem.ChannelTWITTER
	case pb.MessageChannel_VOICE:
		return conversationitem.ChannelVOICE
	default:
		return conversationitem.ChannelCHAT
	}
}

func encodeDirection(direction pb.MessageDirection) conversationitem.Direction {
	switch direction {
	case pb.MessageDirection_INBOUND:
		return conversationitem.DirectionINBOUND
	case pb.MessageDirection_OUTBOUND:
		return conversationitem.DirectionOUTBOUND
	default:
		return conversationitem.DirectionOUTBOUND
	}
}

func encodeType(t pb.MessageType) conversationitem.Type {
	switch t {
	case pb.MessageType_MESSAGE:
		return conversationitem.TypeMESSAGE
	case pb.MessageType_FILE:
		return conversationitem.TypeFILE
	default:
		return conversationitem.TypeMESSAGE
	}
}

func decodeType(t conversationitem.Type) pb.MessageType {
	switch t {
	case conversationitem.TypeMESSAGE:
		return pb.MessageType_MESSAGE
	case conversationitem.TypeFILE:
		return pb.MessageType_FILE
	default:
		return pb.MessageType_MESSAGE
	}
}

func decodeDirection(direction conversationitem.Direction) pb.MessageDirection {
	switch direction {
	case conversationitem.DirectionINBOUND:
		return pb.MessageDirection_INBOUND
	case conversationitem.DirectionOUTBOUND:
		return pb.MessageDirection_OUTBOUND
	default:
		return pb.MessageDirection_OUTBOUND
	}
}

func decodeChannel(channel conversationitem.Channel) pb.MessageChannel {
	switch channel {
	case conversationitem.ChannelCHAT:
		return pb.MessageChannel_WIDGET
	case conversationitem.ChannelMAIL:
		return pb.MessageChannel_MAIL
	case conversationitem.ChannelWHATSAPP:
		return pb.MessageChannel_WHATSAPP
	case conversationitem.ChannelFACEBOOK:
		return pb.MessageChannel_FACEBOOK
	case conversationitem.ChannelTWITTER:
		return pb.MessageChannel_TWITTER
	case conversationitem.ChannelVOICE:
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

	graphqlRequest := graphql.NewRequest(`
  				query ($email: String!) {
  					contact_ByEmail(email: $email){firstName,lastName,id}
  				}
    `)

	graphqlRequest.Var("email", email)
	var graphqlResponse map[string]map[string]string
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return &ContactInfo{firstName: graphqlResponse["contact_ByEmail"]["firstName"],
		lastName: graphqlResponse["contact_ByEmail"]["lastName"],
		id:       graphqlResponse["contact_ByEmail"]["id"]}, nil
}

func getContactByPhone(graphqlClient *graphql.Client, e164 string) (*ContactInfo, error) {

	graphqlRequest := graphql.NewRequest(`
  				query ($e164: String!) {
  					contact_ByPhone(e164: $e164){firstName,lastName,id}
  				}
    `)

	graphqlRequest.Var("e164", e164)
	var graphqlResponse map[string]map[string]string
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return &ContactInfo{firstName: graphqlResponse["contact_ByPhone"]["firstName"],
		lastName: graphqlResponse["contact_ByPhone"]["lastName"],
		id:       graphqlResponse["contact_ByPhone"]["id"]}, nil
}

type contactResponse struct {
	Contact struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"LastName"`
		ID           string `json:"id"`
		PhoneNumbers []struct {
			E164 string `json:"e164"`
		} `json:"phoneNumbers"`
		Emails []struct {
			Email string `json:"email"`
		} `json:"emails"`
	} `json:"contact"`
}

func getContactById(graphqlClient *graphql.Client, id string) (*ContactInfo, error) {

	graphqlRequest := graphql.NewRequest(`
  				query ($id: ID!) {
  					contact(id: $id){
						firstName,
						lastName,
						id,
						phoneNumbers {
						   e164
						 }, emails {
						   email
						 }
      				} 
  				}
    `)

	graphqlRequest.Var("id", id)
	var graphqlResponse contactResponse
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		log.Printf("Grapql got error %s", err.Error())
		return nil, err
	}
	contactInfo := &ContactInfo{firstName: graphqlResponse.Contact.FirstName,
		lastName: graphqlResponse.Contact.LastName,
		id:       graphqlResponse.Contact.ID}
	if len(graphqlResponse.Contact.Emails) > 0 {
		contactInfo.email = &graphqlResponse.Contact.Emails[0].Email
	}
	if len(graphqlResponse.Contact.PhoneNumbers) > 0 {
		contactInfo.phone = &graphqlResponse.Contact.PhoneNumbers[0].E164
	}
	return contactInfo, nil
}
func createContactWithEmail(graphqlClient *graphql.Client, firstName string, lastName string, email string) (string, error) {
	graphqlRequest := graphql.NewRequest(`
		mutation CreateContact ($firstName: String!, $lastName: String! $email: String!) {
		  contact_Create(input: {firstName: $firstName,
			lastName: $lastName,
		  email:{email:  $email, label: WORK}}) {
			id
		  }
		}
    `)

	graphqlRequest.Var("firstName", firstName)
	graphqlRequest.Var("lastName", lastName)
	graphqlRequest.Var("email", email)
	var graphqlResponse map[string]map[string]string
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}
	return graphqlResponse["contact_Create"]["id"], nil
}

func createContactWithPhone(graphqlClient *graphql.Client, firstName string, lastName string, phone string) (string, error) {
	graphqlRequest := graphql.NewRequest(`
		mutation CreateContact ($firstName: String!, $lastName: String!, $e164: String!) {
		  contact_Create(input: {
            firstName: $firstName,
			lastName: $lastName,
		    phoneNumber:{e164:  $e164, label: WORK}
		  }) {
			  id
		  }
		}
    `)

	graphqlRequest.Var("firstName", firstName)
	graphqlRequest.Var("lastName", lastName)
	graphqlRequest.Var("e164", phone)
	var graphqlResponse map[string]map[string]string
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}
	return graphqlResponse["contact_Create"]["id"], nil
}

func createConversation(graphqlClient *graphql.Client, userId string, contactId string, feedId int) (string, error) {
	graphqlRequest := graphql.NewRequest(`
			mutation CreateConversation ($userId: ID!, $contactId: ID!, $feedId: ID!) {
				conversationCreate(input: {
					userId: $userId
					contactId: $contactId
					id: $feedId
				}) {
					id
					startedAt
				}
			}
    `)

	graphqlRequest.Var("userId", userId)
	graphqlRequest.Var("contactId", contactId)
	graphqlRequest.Var("feedId", feedId)
	var graphqlResponse map[string]map[string]string
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}
	return graphqlResponse["conversationCreate"]["id"], nil
}

func addMessageToConversationInGraphDb(driver *neo4j.Driver, conversationId string, msg *gen.ConversationItem, t *time.Time) error {
	session := (*driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	channel := msg.Channel.String()

	params := map[string]interface{}{
		"conversationId": conversationId,
		"messageId":      strconv.Itoa(msg.ID),
		"startedAt":      t.UTC(),
		"channel":        channel,
	}
	query := "MATCH (c:Conversation {id:$conversationId})" +
		" MERGE (c)-[:CONSISTS_OF]->(m:Message:Action {id:$messageId})" +
		" ON CREATE SET m.channel=$channel, m.startedAt=$startedAt, m.conversationId=$conversationId"
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(query, params)
		return nil, err
	})

	return err
}

func (s *messageService) SaveMessage(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	var contact *ContactInfo
	var err error
	var conversation *gen.Conversation

	if message.ContactId == nil {
		if message.GetChannel() == pb.MessageChannel_MAIL || message.GetChannel() == pb.MessageChannel_WIDGET {
			displayName, email := parseEmail(*message.Username)
			contact, err = getContactByEmail(s.graphqlClient, email)

			if err != nil {

				log.Printf("Contact %s creating a new contact", email)
				firstName, lastName := "", ""
				if displayName != "" {
					parts := strings.SplitN(displayName, " ", 2)
					firstName = parts[0]
					if len(parts) > 1 {
						lastName = parts[1]
					}
				}
				log.Printf("Making a contact, firstName=%s lastName=%s email=%s", firstName, lastName, email)
				contactId, err := createContactWithEmail(s.graphqlClient, firstName, lastName, email)
				contact = &ContactInfo{
					firstName: firstName,
					lastName:  lastName,
					id:        contactId,
				}
				if err != nil {
					log.Printf("Unable to create contact! %s", err.Error())
					return nil, err
				}
			}
			conversation, err = s.client.Conversation.
				Query().
				Where(genConversation.ContactId(contact.id)).
				First(ctx)
		} else if message.GetChannel() == pb.MessageChannel_VOICE {
			contact, err = getContactByPhone(s.graphqlClient, *message.Username)
			if err != nil {

				log.Printf("Contact %s creating a new contact", message.Username)
				firstName, lastName := "Unknown", "Caller"
				log.Printf("Making a contact, firstName=%s lastName=%s email=%s", firstName, lastName, message.Username)
				contactId, err := createContactWithPhone(s.graphqlClient, firstName, lastName, *message.Username)
				contact = &ContactInfo{
					firstName: firstName,
					lastName:  lastName,
					id:        contactId,
				}
				if err != nil {
					log.Printf("Unable to create contact! %s", err.Error())
					return nil, err
				}
			}
			conversation, err = s.client.Conversation.
				Query().
				Where(genConversation.ContactId(contact.id)).
				First(ctx)
		} else {
			return nil, status.Errorf(codes.Unimplemented, "Contact mapping not implemented yet for %v", message.GetChannel())
		}
	} else {
		contact, err = getContactById(s.graphqlClient, *message.ContactId)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "Couldn't find a contact for id of %d", *message.ContactId)
		}

		conversation, err = s.client.Conversation.Query().Where(genConversation.ContactId(contact.id)).First(ctx)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "Couldn't find a feed for id of %d", *message.ContactId)
		}

		contact = &ContactInfo{firstName: contact.firstName,
			lastName: contact.lastName,
			id:       contact.id,
		}
	}

	now := time.Now()
	if err != nil {
		// can only reach here if message.Contact is nil & the contactId found in neo4j doesn't match a message any feed
		se, _ := status.FromError(err)
		if se.Code() != codes.Unknown {
			return nil, status.Errorf(se.Code(), "Error upserting Feed: %s", err.Error())
		} else {

			conversationCreate := s.client.Conversation.
				Create().
				SetContactId(contact.id).
				SetState(genConversation.StateNEW).
				SetCreatedOn(now).
				SetUpdatedOn(now).
				SetLastMessage(message.Message).
				SetLastSenderId(contact.id).
				SetLastSenderType(genConversation.LastSenderTypeCONTACT)

			if message.GetDirection() == pb.MessageDirection_INBOUND {
				conversationCreate = conversationCreate.SetLastSenderId(contact.id)
				conversationCreate = conversationCreate.SetLastSenderType(genConversation.LastSenderTypeCONTACT)
			} else {
				conversationCreate = conversationCreate.SetLastSenderId(*message.UserId)
				conversationCreate = conversationCreate.SetLastSenderType(genConversation.LastSenderTypeUSER)
			}

			newConversation, err := conversationCreate.Save(ctx)
			conversation = newConversation

			if err != nil {
				se, _ = status.FromError(err)
				return nil, status.Errorf(se.Code(), "Error inserting new Feed %s", err.Error())
			}

			conv, err := createConversation(s.graphqlClient, s.config.Identity.DefaultUserId, contact.id, conversation.ID)

			if err != nil {
				log.Printf("Error making conversation %v", err)
				return nil, err
			}

			log.Printf("Created conversation %s", conv)

		}
	}

	conversationItemCreate := s.client.ConversationItem.
		Create().
		SetMessage(message.GetMessage()).
		SetConversationID(conversation.ID).
		SetChannel(encodeChannel(message.GetChannel())).
		SetDirection(encodeDirection(message.GetDirection())).
		SetType(encodeType(message.GetType())).
		SetTime(now)

	if message.GetDirection() == pb.MessageDirection_INBOUND {
		conversationItemCreate.SetSenderId(contact.id)
		conversationItemCreate.SetSenderType(conversationitem.SenderTypeCONTACT)
	} else {
		conversationItemCreate.SetSenderId(*message.UserId)
		conversationItemCreate.SetSenderType(conversationitem.SenderTypeUSER)
	}

	conversationItem, err := conversationItemCreate.Save(ctx)

	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error inserting message: %s", err.Error())
	}

	//on the first reply of a user, we mark the conversation as in progress
	if (conversation.State == genConversation.StateNEW) && (message.GetDirection() == pb.MessageDirection_OUTBOUND) {
		conversation.State = genConversation.StateIN_PROGRESS
	}

	conversation.LastMessage = message.GetMessage()

	if message.GetDirection() == pb.MessageDirection_INBOUND {
		conversation.LastSenderId = contact.id
	} else {
		conversation.LastSenderId = *message.UserId
	}

	conversation.UpdatedOn = message.Time.AsTime()
	conversation.Update()

	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error inserting message: %s", err.Error())
	}

	addMessageToConversationInGraphDb(s.driver, strconv.Itoa(conversation.ID), conversationItem, &now)
	if err != nil {
		log.Printf("Error saving message metadata in graph %v", err)
		return nil, err
	}

	var id = int64(conversationItem.ID)
	var conversationid = int64(conversation.ID)
	mi := &pb.Message{
		Type:      decodeType(conversationItem.Type),
		Message:   conversationItem.Message,
		Direction: decodeDirection(conversationItem.Direction),
		Channel:   decodeChannel(conversationItem.Channel),
		Username:  message.Username,
		UserId:    message.UserId,
		ContactId: message.ContactId,
		Id:        &id,
		FeedId:    &conversationid,
		Time:      timestamppb.New(now),
	}
	return mi, nil
}

func (s *messageService) GetMessage(ctx context.Context, msgId *pb.Id) (*pb.Message, error) {
	if msgId == nil || msgId.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Message ID must be specified")
	}

	mi, err := s.client.ConversationItem.Get(ctx, int(msgId.GetId()))
	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error finding Message")
	}

	mf, err := s.client.ConversationItem.QueryConversation(mi).First(ctx)
	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error finding Feed")
	}

	messageId := int64(mi.ID)
	conversationid := int64(mf.ID)

	contactById, err := getContactById(s.graphqlClient, mf.ContactId)

	m := &pb.Message{
		Type:      decodeType(mi.Type),
		Message:   mi.Message,
		Direction: decodeDirection(mi.Direction),
		Channel:   decodeChannel(mi.Channel),
		Username:  contactById.email,
		ContactId: &mf.ContactId,
		Id:        &messageId,
		FeedId:    &conversationid,
		Time:      timestamppb.New(mi.Time),
	}

	if mi.Direction == conversationitem.DirectionOUTBOUND {
		m.UserId = &mi.SenderId
	}

	return m, nil
}

func (s *messageService) GetMessages(ctx context.Context, messagesRequest *pb.GetMessagesRequest) (*pb.MessagePagedResponse, error) {
	var messages []*gen.ConversationItem
	var err error
	var conversation *gen.Conversation

	if messagesRequest != nil {
		log.Printf("Looking up messages for conversation with id %d", messagesRequest.GetConversationId())
		conversation, err = s.client.Conversation.Get(ctx, int(messagesRequest.GetConversationId()))
		if err != nil {
			se, _ := status.FromError(err)
			return nil, status.Errorf(se.Code(), "Error finding conversation with id  %d", messagesRequest.GetConversationId())
		}
	} else {
		log.Printf("Conversation id is required")
		return nil, status.Errorf(1, "Conversation id is required")
	}

	limit := 100 // default to 100 if no pagination is specified
	if messagesRequest.GetPageSize() != 0 {
		limit = int(messagesRequest.GetPageSize())
	}

	if messagesRequest.GetBefore() == nil {
		messages, err = s.client.Conversation.QueryConversationItem(conversation).
			Order(gen.Desc(conversationitem.FieldTime)).
			Limit(limit).
			All(ctx)
	} else {
		messages, err = s.client.Conversation.QueryConversationItem(conversation).
			Order(gen.Desc(conversationitem.FieldTime)).
			Where(conversationitem.TimeLT(messagesRequest.GetBefore().AsTime())).
			Limit(limit).
			All(ctx)
	}

	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error getting messages")
	}
	ml := &pb.MessagePagedResponse{Message: make([]*pb.Message, len(messages))}

	conversationId := int64(conversation.ID)

	for i, j := len(messages)-1, 0; i >= 0; i, j = i-1, j+1 {
		var mid = int64(messages[i].ID)
		mi := &pb.Message{
			Type:      decodeType(messages[i].Type),
			Message:   messages[i].Message,
			Direction: decodeDirection(messages[i].Direction),
			Channel:   decodeChannel(messages[i].Channel),
			Username:  nil,
			ContactId: &conversation.ContactId,
			Id:        &mid,
			FeedId:    &conversationId,
			Time:      timestamppb.New(messages[i].Time),
		}
		ml.Message[j] = mi
	}
	return ml, nil
}

func (s *messageService) GetFeeds(ctx context.Context, feedRequest *pb.GetFeedsPagedRequest) (*pb.FeedItemPagedResponse, error) {
	query := s.client.Conversation.Query()

	if feedRequest.GetStateIn() != nil {
		stateIn := make([]genConversation.State, 0, len(feedRequest.GetStateIn()))
		for _, state := range feedRequest.GetStateIn() {
			stateIn = append(stateIn, encodeConversationState(state))
		}
		query.Where(genConversation.StateIn(stateIn...))
	}

	limit := 100 // default to 100 if no pagination is specified
	if feedRequest.GetPageSize() != 0 {
		limit = int(feedRequest.GetPageSize())
	}
	offset := limit * int(feedRequest.GetPage())

	conversations, err := query.Limit(limit).Offset(offset).All(ctx)
	count, err2 := query.Count(ctx)

	if err != nil || err2 != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error getting messages: %s", err.Error())
	}
	fl := &pb.FeedItemPagedResponse{FeedItems: make([]*pb.FeedItem, len(conversations))}
	fl.TotalElements = int32(count)

	for i, conversation := range conversations {
		var id = int64(conversation.ID)
		log.Printf("Got a conversation id of %d", id)

		contactById, err := getContactById(s.graphqlClient, conversation.ContactId)

		if err != nil {
			se, _ := status.FromError(err)
			return nil, status.Errorf(se.Code(), "Error getting messages: %s", err.Error())
		}

		fl.FeedItems[i] = &pb.FeedItem{
			Id:               int64(conversation.ID),
			ContactId:        contactById.id,
			ContactFirstName: contactById.firstName,
			ContactLastName:  contactById.lastName,
			ContactEmail:     *contactById.email,
			State:            decodeConversationState(conversation.State),
			LastSenderId:     conversation.LastSenderId,
			LastSenderType:   decodeSenderType(conversation.LastSenderType),
			Message:          conversation.LastMessage,
			UpdatedOn:        timestamppb.New(conversation.UpdatedOn),
		}

		msg, _ := json.Marshal(fl.FeedItems[i])
		log.Printf("Got a feed item of %s", msg)
	}
	return fl, nil
}

func (s *messageService) GetFeed(ctx context.Context, feedIdRequest *pb.Id) (*pb.FeedItem, error) {
	if feedIdRequest == nil || feedIdRequest.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Feed ID must be specified")
	}

	conversation, err := s.client.Conversation.Get(ctx, int(feedIdRequest.GetId()))
	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error finding conversation")
	}

	contactById, err := getContactById(s.graphqlClient, conversation.ContactId)
	if err != nil {
		se, _ := status.FromError(err)
		return nil, status.Errorf(se.Code(), "Error getting messages: %s", err.Error())
	}

	return &pb.FeedItem{
		Id:               int64(conversation.ID),
		ContactId:        contactById.id,
		ContactFirstName: contactById.firstName,
		ContactLastName:  contactById.lastName,
		ContactEmail:     *contactById.email,
		State:            decodeConversationState(conversation.State),
		LastSenderId:     conversation.LastSenderId,
		LastSenderType:   decodeSenderType(conversation.LastSenderType),
		Message:          conversation.LastMessage,
		UpdatedOn:        timestamppb.New(conversation.UpdatedOn),
	}, nil
}

func NewMessageService(client *gen.Client, driver *neo4j.Driver, graphqlClient *graphql.Client, config *c.Config) *messageService {
	ms := new(messageService)
	ms.client = client
	ms.graphqlClient = graphqlClient
	ms.config = config
	ms.driver = driver
	return ms
}
