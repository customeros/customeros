package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	c "github.com/openline-ai/openline-customer-os/packages/server/message-store/config"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/enttest"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/messageitem"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/resolver"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/message-store/test/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

var graphqlClient *graphql.Client
var messageStoreService *messageService
var neo4jContainer testcontainers.Container
var driver *neo4j.Driver

func TestMain(m *testing.M) {
	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.Driver, ctx context.Context) {
		neo4jt.Close(driver, "Driver")
		neo4jt.Terminate(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	os.Exit(m.Run())
}

func tearDownTestCase() func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(driver)
	}
}

func NewWebServer(t *testing.T) (*httptest.Server, *graphql.Client, *resolver.Resolver, *gen.Client) {
	router := gin.Default()
	server := httptest.NewServer(router)
	handler, resolver := graph.GraphqlHandler()
	router.POST("/query", handler)
	dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	graphqlClient = graphql.NewClient(server.URL + "/query")
	conf := c.Config{}
	conf.Identity.DefaultUserId = "agentsmith"
	messageStoreService = NewMessageService(dbClient, driver, graphqlClient, &conf)
	t.Cleanup(func() {
		log.Printf("Shutting down webserver")
		server.Close()
		dbClient.Close()
	})
	return server, graphqlClient, resolver, dbClient
}

func messageStoreDialer() (*grpc.ClientConn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	proto.RegisterMessageStoreServiceServer(server, messageStoreService)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	dialFunc := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
	ctx := context.Background()
	return grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialFunc))
}

func Test_SaveMessageNewContact(t *testing.T) {
	defer tearDownTestCase()(t)
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()
	createCalled := false
	createConversationCalled := false

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()

	resolver.GetContactByEmail = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "x@x.org") {
			return nil, status.Error(500, "Unexpected email address")
		}
		return nil, status.Error(404, "Address Not Found")
	}
	resolver.ContactCreate = func(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
		createCalled = true
		if !assert.Equal(t, input.FirstName, "Torrey") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.LastName, "Searle") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.Email.Email, "x@x.org") {
			return nil, status.Error(500, "Email")
		}
		return &model.Contact{
			FirstName: "Torrey",
			LastName:  "Searle",
			ID:        "12345678",
		}, nil
	}
	resolver.ConversationCreate = func(ctx context.Context, input model.ConversationInput) (*model.Conversation, error) {
		log.Print("Inside Conversation Create!")
		createConversationCalled = true
		if !assert.Equal(t, input.UserID, "agentsmith") {
			return nil, status.Error(500, "Unknown userid")
		}
		if !assert.Equal(t, input.ContactID, "12345678") {
			return nil, status.Error(500, "Unknown ContactID")
		}
		return &model.Conversation{
			ID:        "7",
			Contact:   &model.Contact{ID: "12345678", FirstName: "Torrey", LastName: "Searle"},
			StartedAt: time.Now(),
			User:      &model.User{ID: "agentsmith", FirstName: "Agent", LastName: "Smith"},
		}, nil
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)
	neo4jt.CreateConversation(driver, "1")

	result, err := grpcClient.SaveMessage(ctx, &proto.Message{Message: "Hello Torrey",
		Direction: proto.MessageDirection_INBOUND,
		Channel:   proto.MessageChannel_WIDGET,
		Type:      proto.MessageType_MESSAGE,
		Username:  "Torrey Searle <x@x.org>"})
	if err != nil {
		t.Fatalf("Error getting message: %s", err.Error())
	}
	log.Printf("Message saved with ID %d", *result.Id)

	mi, err := dbClient.MessageItem.Get(ctx, int(*result.Id))

	if err != nil {
		t.Fatalf("Error finding Message %s", err.Error())
	}
	mf, err := dbClient.MessageItem.QueryMessageFeed(mi).First(ctx)

	assert.True(t, createCalled)
	assert.True(t, createConversationCalled)
	assert.Equal(t, mi.Message, "Hello Torrey")
	assert.Equal(t, mi.Direction, messageitem.DirectionINBOUND)
	assert.Equal(t, mi.Channel, messageitem.ChannelCHAT)
	assert.Equal(t, mi.Type, messageitem.TypeMESSAGE)
	assert.Equal(t, mi.Username, "Torrey Searle <x@x.org>")
	assert.Equal(t, mf.ContactId, "12345678")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Message"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Action"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "CONSISTS_OF"))
}

func Test_SaveMessageNewPhone(t *testing.T) {
	defer tearDownTestCase()(t)
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()
	createCalled := false
	createConversationCalled := false

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()
	resolver.GetContactByPhone = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "+3228080000") {
			return nil, status.Error(500, "Unexpected email address")
		}
		return nil, status.Error(404, "Phone Number Not Found")
	}
	resolver.ContactCreate = func(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
		createCalled = true
		if !assert.Equal(t, input.FirstName, "Unknown") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.LastName, "Caller") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.PhoneNumber.E164, "+3228080000") {
			return nil, status.Error(500, "Email")
		}
		return &model.Contact{
			FirstName: "Unknown",
			LastName:  "Caller",
			ID:        "12345678",
		}, nil
	}
	resolver.ConversationCreate = func(ctx context.Context, input model.ConversationInput) (*model.Conversation, error) {
		log.Print("Inside Conversation Create!")
		createConversationCalled = true
		if !assert.Equal(t, input.UserID, "agentsmith") {
			return nil, status.Error(500, "Unknown userid")
		}
		if !assert.Equal(t, input.ContactID, "12345678") {
			return nil, status.Error(500, "Unknown ContactID")
		}
		return &model.Conversation{
			ID:        "7",
			Contact:   &model.Contact{ID: "12345678", FirstName: "Torrey", LastName: "Searle"},
			StartedAt: time.Now(),
			User:      &model.User{ID: "agentsmith", FirstName: "Agent", LastName: "Smith"},
		}, nil
	}

	neo4jt.CreateConversation(driver, "1")

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)

	result, err := grpcClient.SaveMessage(ctx, &proto.Message{Message: "Call, duration 5 Minutes",
		Direction: proto.MessageDirection_INBOUND,
		Channel:   proto.MessageChannel_VOICE,
		Type:      proto.MessageType_MESSAGE,
		Username:  "+3228080000"})
	if err != nil {
		t.Fatalf("Error getting message: %s", err.Error())
	}
	log.Printf("Message saved with ID %d", *result.Id)

	mi, err := dbClient.MessageItem.Get(ctx, int(*result.Id))

	if err != nil {
		t.Fatalf("Error finding Message %s", err.Error())
	}
	mf, err := dbClient.MessageItem.QueryMessageFeed(mi).First(ctx)

	assert.True(t, createCalled)
	assert.True(t, createConversationCalled)
	assert.Equal(t, mi.Message, "Call, duration 5 Minutes")
	assert.Equal(t, mi.Direction, messageitem.DirectionINBOUND)
	assert.Equal(t, mi.Channel, messageitem.ChannelVOICE)
	assert.Equal(t, mi.Type, messageitem.TypeMESSAGE)
	assert.Equal(t, mi.Username, "+3228080000")
	assert.Equal(t, mf.ContactId, "12345678")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Message"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Action"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "CONSISTS_OF"))
}

func Test_SaveMessageMissingGraphContact(t *testing.T) {
	defer tearDownTestCase()(t)
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()
	createCalled := false

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()

	feed1, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	_, err = dbClient.MessageFeed.
		Create().
		SetFirstName("Echo").
		SetLastName("Test").
		SetContactId("echotest").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}
	resolver.GetContactByEmail = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "x@x.org") {
			return nil, status.Error(500, "Unexpected email address")
		}
		return nil, status.Error(404, "Address Not Found")
	}
	resolver.ContactCreate = func(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
		createCalled = true
		if !assert.Equal(t, input.FirstName, "Torrey") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.LastName, "Searle") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.Email.Email, "x@x.org") {
			return nil, status.Error(500, "Email")
		}
		return &model.Contact{
			FirstName: "Torrey",
			LastName:  "Searle",
			ID:        "12345678",
		}, nil
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)

	result, err := grpcClient.SaveMessage(ctx, &proto.Message{Message: "Hello Torrey",
		Direction: proto.MessageDirection_INBOUND,
		Channel:   proto.MessageChannel_WIDGET,
		Type:      proto.MessageType_MESSAGE,
		Username:  "Torrey Searle <x@x.org>"})
	if err != nil {
		t.Fatalf("Error getting message: %s", err.Error())
	}
	log.Printf("Message saved with ID %d", *result.Id)

	mi, err := dbClient.MessageItem.Get(ctx, int(*result.Id))

	if err != nil {
		t.Fatalf("Error finding Message %s", err.Error())
	}
	mf, err := dbClient.MessageItem.QueryMessageFeed(mi).First(ctx)

	assert.True(t, createCalled)
	assert.Equal(t, mi.Message, "Hello Torrey")
	assert.Equal(t, mi.Direction, messageitem.DirectionINBOUND)
	assert.Equal(t, mi.Channel, messageitem.ChannelCHAT)
	assert.Equal(t, mi.Type, messageitem.TypeMESSAGE)
	assert.Equal(t, mi.Username, "Torrey Searle <x@x.org>")
	assert.Equal(t, mf.ID, feed1.ID)
	assert.Equal(t, mf.ContactId, "12345678")
}

func Test_SaveMessageExistingContact(t *testing.T) {
	defer tearDownTestCase()(t)
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()

	feed1, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}
	neo4jt.CreateConversation(driver, strconv.Itoa(feed1.ID))

	_, err = dbClient.MessageFeed.
		Create().
		SetFirstName("Echo").
		SetLastName("Test").
		SetContactId("echotest").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}
	resolver.GetContactByEmail = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "x@x.org") {
			return nil, status.Error(500, "Unexpected email address")
		}
		return &model.Contact{ID: "12345678", FirstName: "Torrey", LastName: "Searle"}, nil
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)

	result, err := grpcClient.SaveMessage(ctx, &proto.Message{Message: "Hello Torrey",
		Direction: proto.MessageDirection_INBOUND,
		Channel:   proto.MessageChannel_WIDGET,
		Type:      proto.MessageType_MESSAGE,
		Username:  "Torrey Searle <x@x.org>"})
	if err != nil {
		t.Fatalf("Error getting message: %s", err.Error())
	}
	log.Printf("Message saved with ID %d", *result.Id)

	mi, err := dbClient.MessageItem.Get(ctx, int(*result.Id))

	if err != nil {
		t.Fatalf("Error finding Message %s", err.Error())
	}
	mf, err := dbClient.MessageItem.QueryMessageFeed(mi).First(ctx)

	assert.Equal(t, mi.Message, "Hello Torrey")
	assert.Equal(t, mi.Direction, messageitem.DirectionINBOUND)
	assert.Equal(t, mi.Channel, messageitem.ChannelCHAT)
	assert.Equal(t, mi.Type, messageitem.TypeMESSAGE)
	assert.Equal(t, mi.Username, "Torrey Searle <x@x.org>")
	assert.Equal(t, mf.ID, feed1.ID)
	assert.Equal(t, mf.ContactId, "12345678")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Message"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Action"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "CONSISTS_OF"))
}

func Test_SaveMessageExistingPhoneContact(t *testing.T) {
	defer tearDownTestCase()(t)
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()

	feed1, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	_, err = dbClient.MessageFeed.
		Create().
		SetFirstName("Echo").
		SetLastName("Test").
		SetContactId("echotest").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}
	resolver.GetContactByPhone = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "+3228080000") {
			return nil, status.Error(500, "Unexpected email address")
		}
		return &model.Contact{ID: "12345678", FirstName: "Torrey", LastName: "Searle"}, nil
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)

	result, err := grpcClient.SaveMessage(ctx, &proto.Message{Message: "Call, duration 5 Minutes",
		Direction: proto.MessageDirection_INBOUND,
		Channel:   proto.MessageChannel_VOICE,
		Type:      proto.MessageType_MESSAGE,
		Username:  "+3228080000"})
	if err != nil {
		t.Fatalf("Error getting message: %s", err.Error())
	}
	log.Printf("Message saved with ID %d", *result.Id)

	mi, err := dbClient.MessageItem.Get(ctx, int(*result.Id))

	if err != nil {
		t.Fatalf("Error finding Message %s", err.Error())
	}
	mf, err := dbClient.MessageItem.QueryMessageFeed(mi).First(ctx)

	assert.Equal(t, mi.Message, "Call, duration 5 Minutes")
	assert.Equal(t, mi.Direction, messageitem.DirectionINBOUND)
	assert.Equal(t, mi.Channel, messageitem.ChannelVOICE)
	assert.Equal(t, mi.Type, messageitem.TypeMESSAGE)
	assert.Equal(t, mi.Username, "+3228080000")
	assert.Equal(t, mf.ID, feed1.ID)
	assert.Equal(t, mf.ContactId, "12345678")
}

func Test_SaveMessageSpecifiedContact(t *testing.T) {
	defer tearDownTestCase()(t)
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()

	feed1, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	id := int64(feed1.ID)

	_, err = dbClient.MessageFeed.
		Create().
		SetFirstName("Echo").
		SetLastName("Test").
		SetContactId("echotest").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	resolver.GetContactById = func(ctx context.Context, id string) (*model.Contact, error) {
		if id == "12345678" {
			return &model.Contact{
				ID:        id,
				FirstName: "Torrey",
				LastName:  "Searle",
			}, nil
		} else if id == "echotest" {
			return nil, status.Error(404, "id not found")
		} else {
			t.Errorf("Unexpected ID %s", id)
			return nil, status.Error(500, "id incorrect")

		}
	}
	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.PhoneNumber{&model.PhoneNumber{E164: "+3228080000"}}, nil
	}

	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{&model.Email{Email: "x@x.org"}}, nil
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)

	result, err := grpcClient.SaveMessage(ctx, &proto.Message{Message: "Hello Torrey",
		Contact:   &proto.Contact{Id: &id},
		Direction: proto.MessageDirection_INBOUND,
		Channel:   proto.MessageChannel_WIDGET,
		Type:      proto.MessageType_MESSAGE,
		Username:  "x@x.org"})
	if err != nil {
		t.Fatalf("Error getting message: %s", err.Error())
	}
	log.Printf("Message saved with ID %d", *result.Id)

	mi, err := dbClient.MessageItem.Get(ctx, int(*result.Id))

	if err != nil {
		t.Fatalf("Error finding Message %s", err.Error())
	}
	mf, err := dbClient.MessageItem.QueryMessageFeed(mi).First(ctx)

	assert.Equal(t, mi.Message, "Hello Torrey")
	assert.Equal(t, mi.Direction, messageitem.DirectionINBOUND)
	assert.Equal(t, mi.Channel, messageitem.ChannelCHAT)
	assert.Equal(t, mi.Type, messageitem.TypeMESSAGE)
	assert.Equal(t, mi.Username, "x@x.org")
	assert.Equal(t, mf.ID, feed1.ID)
	assert.Equal(t, mf.ContactId, "12345678")
}

func Test_GetMessage(t *testing.T) {
	_, _, _, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}
	ctx := context.Background()

	feed1, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	oldTime := time.Now()
	oldTime = oldTime.Add(-24 * time.Hour) // old time is yesterday
	msg1, err := dbClient.MessageItem.
		Create().
		SetMessage("Hello Torrey").
		SetMessageFeed(feed1).
		SetChannel(messageitem.ChannelCHAT).
		SetUsername("x@x.org").
		SetNillableTime(&oldTime).
		SetDirection(messageitem.DirectionINBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}

	msg2, err := dbClient.MessageItem.
		Create().
		SetMessage("How may I help you").
		SetMessageFeed(feed1).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("x@x.org").
		SetDirection(messageitem.DirectionOUTBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)
	msg1id := int64(msg1.ID)
	message, err := grpcClient.GetMessage(ctx, &proto.Message{Id: &msg1id})
	if err != nil {
		t.Fatalf("Error getting message: %s", err.Error())
	}
	assert.Equal(t, "Hello Torrey", message.Message)
	assert.Equal(t, "12345678", message.Contact.ContactId)
	assert.Equal(t, int64(feed1.ID), *message.Id)
	assert.Equal(t, proto.MessageDirection_INBOUND, message.Direction)
	assert.Equal(t, proto.MessageChannel_WIDGET, message.Channel)
	assert.Equal(t, "x@x.org", message.Username)

	msg2id := int64(msg2.ID)
	message, err = grpcClient.GetMessage(ctx, &proto.Message{Id: &msg2id})
	if err != nil {
		t.Fatalf("Error getting message: %s", err.Error())
	}

	assert.Equal(t, "How may I help you", message.Message)
	assert.Equal(t, "12345678", message.Contact.ContactId)
	assert.Equal(t, int64(feed1.ID), *message.Contact.Id)
	assert.Equal(t, proto.MessageDirection_OUTBOUND, message.Direction)
	assert.Equal(t, proto.MessageChannel_MAIL, message.Channel)
	assert.Equal(t, "x@x.org", message.Username)

}

func Test_GetMessagesWithLimitBefore(t *testing.T) {
	defer tearDownTestCase()(t)
	_, _, _, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}
	ctx := context.Background()

	feed1, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	id := int64(feed1.ID)

	feed2, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Echo").
		SetLastName("Test").
		SetContactId("echotest").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	oldTime := time.Now()
	oldTime = oldTime.Add(-24 * time.Hour) // old time is yesterday
	_, err = dbClient.MessageItem.
		Create().
		SetMessage("Hello Torrey").
		SetMessageFeed(feed1).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("x@x.org").
		SetNillableTime(&oldTime).
		SetDirection(messageitem.DirectionINBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}

	_, err = dbClient.MessageItem.
		Create().
		SetMessage("How may I help you").
		SetMessageFeed(feed1).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("x@x.org").
		SetDirection(messageitem.DirectionOUTBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}
	_, err = dbClient.MessageItem.
		Create().
		SetMessage("Call Me").
		SetMessageFeed(feed2).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("echo@oasis.openline.ai").
		SetDirection(messageitem.DirectionINBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}
	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)
	pageTime := time.Now()
	pageTime = pageTime.Add(-12 * time.Hour) // page between now and first message
	before := timestamppb.New(pageTime)
	messageList, err := grpcClient.GetMessages(ctx, &proto.PagedContact{Page: &proto.PageInfo{PageSize: 1, Before: before}, Contact: &proto.Contact{Id: &id}})

	assert.Equal(t, 1, len(messageList.Message))
	assert.Equal(t, "Hello Torrey", messageList.Message[0].Message)

}

func Test_GetMessagesWithLimit(t *testing.T) {
	_, _, _, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}
	ctx := context.Background()

	feed1, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	id := int64(feed1.ID)

	feed2, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Echo").
		SetLastName("Test").
		SetContactId("echotest").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	oldTime := time.Now()
	oldTime = oldTime.Add(-24 * time.Hour) // old time is yesterday
	_, err = dbClient.MessageItem.
		Create().
		SetMessage("Hello Torrey").
		SetMessageFeed(feed1).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("x@x.org").
		SetNillableTime(&oldTime).
		SetDirection(messageitem.DirectionINBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}

	_, err = dbClient.MessageItem.
		Create().
		SetMessage("How may I help you").
		SetMessageFeed(feed1).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("x@x.org").
		SetDirection(messageitem.DirectionOUTBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}
	_, err = dbClient.MessageItem.
		Create().
		SetMessage("Call Me").
		SetMessageFeed(feed2).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("echo@oasis.openline.ai").
		SetDirection(messageitem.DirectionINBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}
	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)
	messageList, err := grpcClient.GetMessages(ctx, &proto.PagedContact{Page: &proto.PageInfo{PageSize: 1}, Contact: &proto.Contact{Id: &id}})

	assert.Equal(t, 1, len(messageList.Message))
	assert.Equal(t, "How may I help you", messageList.Message[0].Message)
}

func Test_GetMessages(t *testing.T) {
	_, _, _, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}
	ctx := context.Background()

	feed1, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	id := int64(feed1.ID)

	feed2, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Echo").
		SetLastName("Test").
		SetContactId("echotest").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	_, err = dbClient.MessageItem.
		Create().
		SetMessage("Hello Torrey").
		SetMessageFeed(feed1).
		SetChannel(messageitem.ChannelCHAT).
		SetUsername("x@x.org").
		SetDirection(messageitem.DirectionINBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}

	_, err = dbClient.MessageItem.
		Create().
		SetMessage("How may I help you").
		SetMessageFeed(feed1).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("x@x.org").
		SetDirection(messageitem.DirectionOUTBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}
	_, err = dbClient.MessageItem.
		Create().
		SetMessage("Call Me").
		SetMessageFeed(feed2).
		SetChannel(messageitem.ChannelMAIL).
		SetUsername("echo@oasis.openline.ai").
		SetDirection(messageitem.DirectionINBOUND).
		SetType(messageitem.TypeMESSAGE).
		Save(ctx)

	if err != nil {
		t.Fatalf("Error inserting message: %s", err.Error())
	}
	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)
	messageList, err := grpcClient.GetMessages(ctx, &proto.PagedContact{Contact: &proto.Contact{Id: &id}})

	assert.Equal(t, 2, len(messageList.Message))
	assert.Equal(t, "Hello Torrey", messageList.Message[0].Message)
	assert.Equal(t, "12345678", messageList.Message[0].Contact.ContactId)
	assert.Equal(t, int64(feed1.ID), *messageList.Message[0].Contact.Id)
	assert.Equal(t, proto.MessageDirection_INBOUND, messageList.Message[0].Direction)
	assert.Equal(t, proto.MessageChannel_WIDGET, messageList.Message[0].Channel)
	assert.Equal(t, "x@x.org", messageList.Message[0].Username)

	assert.Equal(t, "How may I help you", messageList.Message[1].Message)
	assert.Equal(t, "12345678", messageList.Message[1].Contact.ContactId)
	assert.Equal(t, int64(feed1.ID), *messageList.Message[1].Contact.Id)
	assert.Equal(t, proto.MessageDirection_OUTBOUND, messageList.Message[1].Direction)
	assert.Equal(t, proto.MessageChannel_MAIL, messageList.Message[1].Channel)
	assert.Equal(t, "x@x.org", messageList.Message[1].Username)

}

func Test_GetFeeds(t *testing.T) {
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()

	feed, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	id := int64(feed.ID)

	feed, err = dbClient.MessageFeed.
		Create().
		SetFirstName("Echo").
		SetLastName("Test").
		SetContactId("echotest").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}

	resolver.GetContactById = func(ctx context.Context, id string) (*model.Contact, error) {
		if id == "12345678" {
			return &model.Contact{
				ID:        id,
				FirstName: "Torrey",
				LastName:  "Searle",
			}, nil
		} else if id == "echotest" {
			return nil, status.Error(404, "id not found")
		} else {
			t.Errorf("Unexpected ID %s", id)
			return nil, status.Error(500, "id incorrect")

		}
	}
	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.PhoneNumber{&model.PhoneNumber{E164: "+3228080000"}}, nil
	}

	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{&model.Email{Email: "x@x.org"}}, nil
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)
	contacts, err := grpcClient.GetFeeds(ctx, &proto.Empty{})
	assert.Equal(t, 2, len(contacts.Contact))
	if err != nil {
		t.Fatalf("Error getting Feed %s", err.Error())
	}

	for i := range contacts.Contact {
		if *contacts.Contact[i].Id == id {
			assert.Equal(t, "12345678", contacts.Contact[i].ContactId)
			assert.Equal(t, "Torrey", contacts.Contact[i].FirstName)
			assert.Equal(t, "Searle", contacts.Contact[i].LastName)
			assert.Equal(t, "x@x.org", *contacts.Contact[i].Email)
			assert.Equal(t, "+3228080000", *contacts.Contact[i].Phone)
		} else {
			assert.Equal(t, "echotest", contacts.Contact[i].ContactId)
			assert.Equal(t, "Echo", contacts.Contact[i].FirstName)
			assert.Equal(t, "Test", contacts.Contact[i].LastName)
		}
	}
}

func Test_GetFeed(t *testing.T) {
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()

	feed, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}
	resolver.GetContactById = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return &model.Contact{
			ID:        id,
			FirstName: "Torrey",
			LastName:  "Searle",
		}, nil
	}

	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.PhoneNumber{&model.PhoneNumber{E164: "+3228080000"}}, nil
	}

	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{&model.Email{Email: "x@x.org"}}, nil
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)
	var id int64 = int64(feed.ID)
	contact, err := grpcClient.GetFeed(ctx, &proto.Contact{Id: &id})
	if err != nil {
		t.Fatalf("Error getting Feed %s", err.Error())
	}
	assert.Equal(t, "12345678", contact.ContactId)
	assert.Equal(t, "Torrey", contact.FirstName)
	assert.Equal(t, "Searle", contact.LastName)
	assert.Equal(t, id, *contact.Id)
	assert.Equal(t, "x@x.org", *contact.Email)
	assert.Equal(t, "+3228080000", *contact.Phone)
}

func Test_GetFeedByContactId(t *testing.T) {
	_, _, resolver, dbClient := NewWebServer(t)
	grpcConn, err := messageStoreDialer()

	if err != nil {
		t.Fatalf("Unable to make GRPC service, %s", err.Error())
	}

	ctx := context.Background()

	feed, err := dbClient.MessageFeed.
		Create().
		SetFirstName("Torrey").
		SetLastName("Searle").
		SetContactId("12345678").
		Save(ctx)
	if err != nil {
		t.Fatalf("Error inserting new Feed: %s", err.Error())
	}
	resolver.GetContactById = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return &model.Contact{
			ID:        id,
			FirstName: "Torrey",
			LastName:  "Searle",
		}, nil
	}
	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.PhoneNumber{&model.PhoneNumber{E164: "+3228080000"}}, nil
	}

	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{&model.Email{Email: "x@x.org"}}, nil
	}

	grpcClient := proto.NewMessageStoreServiceClient(grpcConn)
	var id int64 = int64(feed.ID)
	contact, err := grpcClient.GetFeed(ctx, &proto.Contact{ContactId: "12345678"})
	if err != nil {
		t.Fatalf("Error getting Feed %s", err.Error())
	}
	assert.Equal(t, "12345678", contact.ContactId)
	assert.Equal(t, "Torrey", contact.FirstName)
	assert.Equal(t, "Searle", contact.LastName)
	assert.Equal(t, id, *contact.Id)
	assert.Equal(t, "x@x.org", *contact.Email)
	assert.Equal(t, "+3228080000", *contact.Phone)
}

func Test_parseEmail(t *testing.T) {
	type args struct {
		email string
	}

	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{name: "unquoted displayname",
			args:  args{email: "Torrey Searle <tsearle@invalid.domain>"},
			want:  "Torrey Searle",
			want1: "tsearle@invalid.domain",
		},
		{name: "quoted displayname",
			args:  args{email: "\"Torrey Searle\" <tsearle@invalid.domain>"},
			want:  "Torrey Searle",
			want1: "tsearle@invalid.domain",
		},
		{name: "no display name",
			args:  args{email: "tsearle@invalid.domain"},
			want:  "",
			want1: "tsearle@invalid.domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseEmail(tt.args.email)
			if got != tt.want {
				t.Errorf("parseEmail() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseEmail() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_createContactWithEmail(t *testing.T) {
	_, client, resolver, _ := NewWebServer(t)
	resolver.ContactCreate = func(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
		if !assert.Equal(t, input.FirstName, "Torrey") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.LastName, "Searle") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.Email.Email, "x@x.org") {
			return nil, status.Error(500, "Email")
		}
		return &model.Contact{
			FirstName: "Torrey",
			LastName:  "Searle",
			ID:        "12345678",
		}, nil
	}
	result, err := createContactWithEmail(client, "Torrey", "Searle", "x@x.org")
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, "12345678", result)
}

func Test_createContactWithPhone(t *testing.T) {
	_, client, resolver, _ := NewWebServer(t)
	resolver.ContactCreate = func(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
		if !assert.Equal(t, input.FirstName, "Torrey") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.LastName, "Searle") {
			return nil, status.Error(500, "Unknown Firstname")
		}
		if !assert.Equal(t, input.PhoneNumber.E164, "+328080000") {
			return nil, status.Error(500, "Email")
		}
		return &model.Contact{
			FirstName: "Torrey",
			LastName:  "Searle",
			ID:        "12345678",
		}, nil
	}
	result, err := createContactWithPhone(client, "Torrey", "Searle", "+328080000")
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, "12345678", result)
}

func Test_getContact(t *testing.T) {
	_, client, resolver, _ := NewWebServer(t)
	resolver.GetContactById = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return &model.Contact{
			ID:        id,
			FirstName: "Torrey",
			LastName:  "Searle",
		}, nil
	}

	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.PhoneNumber{&model.PhoneNumber{E164: "+3228080000"}}, nil
	}

	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, "12345678") {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{&model.Email{Email: "x@x.org"}}, nil
	}

	result, err := getContactById(client, "12345678")
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, "12345678", result.id)
	assert.Equal(t, "+3228080000", *result.phone)
	assert.Equal(t, "x@x.org", *result.email)
}

func Test_createConversation(t *testing.T) {
	defer tearDownTestCase()(t)
	_, client, resolver, _ := NewWebServer(t)
	resolver.ConversationCreate = func(ctx context.Context, input model.ConversationInput) (*model.Conversation, error) {
		log.Print("Inside Conversation Create!")
		if !assert.Equal(t, input.UserID, "agentsmith") {
			return nil, status.Error(500, "Unknown userid")
		}
		if !assert.Equal(t, input.ContactID, "12345678") {
			return nil, status.Error(500, "Unknown ContactID")
		}
		if !assert.Equal(t, *input.ID, "7") {
			return nil, status.Error(500, "Unknown feedId")
		}
		return &model.Conversation{
			ID:        "7",
			Contact:   &model.Contact{ID: "12345678", FirstName: "Torrey", LastName: "Searle"},
			StartedAt: time.Now(),
			User:      &model.User{ID: "agentsmith", FirstName: "Agent", LastName: "Smith"},
		}, nil
	}
	result, err := createConversation(client, "agentsmith", "12345678", 7)
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, "7", result)
}

func Test_getConversationByEmail(t *testing.T) {
	_, client, resolver, _ := NewWebServer(t)
	resolver.GetContactByEmail = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "x@x.org") {
			return nil, status.Error(500, "Unexpected email address")
		}
		return &model.Contact{ID: "12345678", FirstName: "Torrey", LastName: "Searle"}, nil
	}

	result, err := getContactByEmail(client, "x@x.org")
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, result.firstName, "Torrey")
	assert.Equal(t, result.lastName, "Searle")
	assert.Equal(t, result.id, "12345678")

}

func Test_getConversationByPhone(t *testing.T) {
	_, client, resolver, _ := NewWebServer(t)
	resolver.GetContactByPhone = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, "+3228080000") {
			return nil, status.Error(500, "Unexpected email address")
		}
		return &model.Contact{ID: "12345678", FirstName: "Torrey", LastName: "Searle"}, nil
	}

	result, err := getContactByPhone(client, "+3228080000")
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, result.firstName, "Torrey")
	assert.Equal(t, result.lastName, "Searle")
	assert.Equal(t, result.id, "12345678")

}
