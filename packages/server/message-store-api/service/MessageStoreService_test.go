package service

import (
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	msProto "github.com/openline-ai/openline-customer-os/packages/server/message-store-api/proto/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

func NewDatabaseServer(t *testing.T) *MessageService {
	postgresDb := embeddedpostgres.NewDatabase()
	err := postgresDb.Start()
	dbClient, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		log.Printf("Shutting down database")
		postgresDb.Stop()
	})

	return &MessageService{postgresRepositories: repository.InitRepositories(dbClient, nil)}
}

func TestGetMessages(t *testing.T) {
	s := NewDatabaseServer(t)

	message1 := &entity.ConversationEvent{

		ConversationId: "blah",
		SenderId:       "bob",

		Content:    "recent message",
		CreateDate: time.Now().Add(-1 * time.Hour),
	}

	message2 := &entity.ConversationEvent{

		ConversationId: "blah",
		SenderId:       "bob",

		Content:    "older message",
		CreateDate: time.Now().Add(-2 * time.Hour),
	}

	queryResult := s.postgresRepositories.ConversationEventRepository.Save(message1)
	if queryResult.Error != nil {
		t.Fatal(queryResult.Error)
	}

	queryResult = s.postgresRepositories.ConversationEventRepository.Save(message2)

	if queryResult.Error != nil {
		t.Fatal(queryResult.Error)
	}

	// default pagination
	queryResult = s.postgresRepositories.ConversationEventRepository.GetEventsForConversation("blah", nil, 10)
	if queryResult.Error != nil {
		t.Fatal(queryResult.Error)
	}
	var messages []*msProto.MessageDeprecate
	for _, event := range *queryResult.Result.(*[]entity.ConversationEvent) {
		messages = append(messages, s.commonStoreService.EncodeConversationEventToMS(event))
	}
	assert.Equal(t, 2, len(messages))
	t.Logf("Got %d messages", len(messages))

	// first page
	queryResult = s.postgresRepositories.ConversationEventRepository.GetEventsForConversation("blah", nil, 1)
	if queryResult.Error != nil {
		t.Fatal(queryResult.Error)
	}
	messages = make([]*msProto.MessageDeprecate, 0)
	for _, event := range *queryResult.Result.(*[]entity.ConversationEvent) {
		messages = append(messages, s.commonStoreService.EncodeConversationEventToMS(event))
	}
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, "recent message", messages[0].Content)
	t.Logf("Got %d messages", len(messages))

	nextTime := messages[0].Time.AsTime()
	queryResult = s.postgresRepositories.ConversationEventRepository.GetEventsForConversation("blah", &nextTime, 1)
	if queryResult.Error != nil {
		t.Fatal(queryResult.Error)
	}
	messages = make([]*msProto.MessageDeprecate, 0)
	for _, event := range *queryResult.Result.(*[]entity.ConversationEvent) {
		messages = append(messages, s.commonStoreService.EncodeConversationEventToMS(event))
	}
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, "older message", messages[0].Content)
	t.Logf("Got %d messages", len(messages))
	// ...
}
