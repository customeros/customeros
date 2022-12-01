package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/resolver"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
	"log"
	"net/http/httptest"
	"testing"
	"time"
)

var graphqlClient *graphql.Client

func NewWebServer(t *testing.T) (*httptest.Server, *graphql.Client, *resolver.Resolver) {
	router := gin.Default()
	server := httptest.NewServer(router)
	handler, resolver := graph.GraphqlHandler()
	router.POST("/query", handler)

	graphqlClient = graphql.NewClient(server.URL + "/query")

	t.Cleanup(func() {
		log.Printf("Shutting down webserver")
		server.Close()
	})
	return server, graphqlClient, resolver
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

func Test_createContact(t *testing.T) {
	_, client, resolver := NewWebServer(t)
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
	result, err := createContact(client, "Torrey", "Searle", "x@x.org")
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, "12345678", result)
}

func Test_getContact(t *testing.T) {
	_, client, resolver := NewWebServer(t)
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
	_, client, resolver := NewWebServer(t)
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
	_, client, resolver := NewWebServer(t)
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
