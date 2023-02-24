package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	commonModuleService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/resolver"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http/httptest"
	"testing"
	"time"
)

func NewWebServer(t *testing.T) (*httptest.Server, *graphql.Client, *resolver.Resolver, context.Context) {
	router := gin.Default()
	server := httptest.NewServer(router)
	handler, resolver := graph.GraphqlHandler()
	router.POST("/query", handler)
	graphqlClient := graphql.NewClient(server.URL + "/query")

	md := metadata.New(map[string]string{commonModuleService.UsernameHeader: "x@x.org"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	t.Cleanup(func() {
		log.Printf("Shutting down webserver")
		server.Close()
	})
	return server, graphqlClient, resolver, ctx
}

func Test_GetUserByEmail(t *testing.T) {
	_, graphqlClient, resolver, ctx := NewWebServer(t)
	service := &CustomerOSService{graphqlClient: graphqlClient,
		conf: &config.Config{

			Service: struct {
				ServerPort       int    `env:"MESSAGE_STORE_SERVER_PORT,required"`
				CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
				CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
			}{CustomerOsAPIKey: "Hello World"},
		},
	}
	resolver.UserByEmail = func(ctx context.Context, email string) (*model.User, error) {
		if !assert.Equal(t, email, "x@x.org") {
			return nil, status.Error(500, "Unexpected email address")
		}
		return &model.User{
			ID:            "my_user_id",
			FirstName:     "Bonnie",
			LastName:      "Clyde",
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
			Source:        "manual",
			Conversations: nil,
		}, nil
	}
	resolver.EmailsByUser = func(ctx context.Context, obj *model.User) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, "my_user_id") {
			return nil, status.Error(500, "id incorrect")
		}
		mail := "x@x.org"

		return []*model.Email{&model.Email{Email: &mail}}, nil
	}

	user, err := service.GetUserByEmail(ctx, "x@x.org")
	if assert.NoErrorf(t, err, "Unexpected error: %v", err) {
		assert.Equal(t, "my_user_id", user.Id)
	}
}

func Test_GetContactById(t *testing.T) {
	_, graphqlClient, resolver, ctx := NewWebServer(t)
	service := &CustomerOSService{graphqlClient: graphqlClient,
		conf: &config.Config{

			Service: struct {
				ServerPort       int    `env:"MESSAGE_STORE_SERVER_PORT,required"`
				CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
				CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
			}{CustomerOsAPIKey: "Hello World"},
		},
	}

	my_id := "my_contact_id"
	resolver.GetContactById = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, my_id) {
			return nil, status.Error(500, "Unexpected id")
		}
		firstName := "Bonnie"
		lastName := "Clyde"
		return &model.Contact{
			ID:        id,
			FirstName: &firstName,
			LastName:  &lastName,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Source:    "manual",
		}, nil
	}
	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}
		mail := "x@x.org"

		return []*model.Email{&model.Email{Email: &mail}}, nil
	}

	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}
		phoneNumber := "+123456"

		return []*model.PhoneNumber{&model.PhoneNumber{E164: &phoneNumber}}, nil
	}
	contact, err := service.GetContactById(ctx, my_id)
	if assert.NoErrorf(t, err, "Unexpected error: %v", err) {
		assert.Equal(t, my_id, contact.Id)
		assert.Equal(t, "x@x.org", contact.Emails[0].Email)
		assert.Equal(t, "+123456", contact.PhoneNumbers[0].E164)

	}
}

func Test_GetContactByEmail(t *testing.T) {
	_, graphqlClient, resolver, ctx := NewWebServer(t)
	service := &CustomerOSService{graphqlClient: graphqlClient,
		conf: &config.Config{

			Service: struct {
				ServerPort       int    `env:"MESSAGE_STORE_SERVER_PORT,required"`
				CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
				CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
			}{CustomerOsAPIKey: "Hello World"},
		},
	}

	my_email := "x@x.org"
	my_id := "my_contact_id"
	my_phone := "+123456"

	resolver.GetContactByEmail = func(ctx context.Context, email string) (*model.Contact, error) {
		if !assert.Equal(t, email, my_email) {
			return nil, status.Error(500, "Unexpected id")
		}
		firstName := "Bonnie"
		lastName := "Clyde"
		return &model.Contact{
			ID:        my_id,
			FirstName: &firstName,
			LastName:  &lastName,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Source:    "manual",
		}, nil
	}
	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{&model.Email{Email: &my_email}}, nil
	}

	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}

		return []*model.PhoneNumber{&model.PhoneNumber{E164: &my_phone}}, nil
	}
	contact, err := service.GetContactByEmail(ctx, my_email)
	if assert.NoErrorf(t, err, "Unexpected error: %v", err) {
		assert.Equal(t, my_id, contact.Id)
		assert.Equal(t, my_email, contact.Emails[0].Email)
		assert.Equal(t, my_phone, contact.PhoneNumbers[0].E164)

	}
}

func Test_GetContactByPhone(t *testing.T) {
	_, graphqlClient, resolver, ctx := NewWebServer(t)
	service := &CustomerOSService{graphqlClient: graphqlClient,
		conf: &config.Config{

			Service: struct {
				ServerPort       int    `env:"MESSAGE_STORE_SERVER_PORT,required"`
				CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
				CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
			}{CustomerOsAPIKey: "Hello World"},
		},
	}

	my_email := "x@x.org"
	my_id := "my_contact_id"
	my_phone := "+123456"
	resolver.GetContactByPhone = func(ctx context.Context, phone string) (*model.Contact, error) {
		if !assert.Equal(t, phone, my_phone) {
			return nil, status.Error(500, "Unexpected id")
		}
		firstName := "Bonnie"
		lastName := "Clyde"
		return &model.Contact{
			ID:        my_id,
			FirstName: &firstName,
			LastName:  &lastName,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Source:    "manual",
		}, nil
	}
	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{&model.Email{Email: &my_email}}, nil
	}

	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}

		return []*model.PhoneNumber{&model.PhoneNumber{E164: &my_phone}}, nil
	}
	contact, err := service.GetContactByPhone(ctx, my_phone)
	if assert.NoErrorf(t, err, "Unexpected error: %v", err) {
		assert.Equal(t, my_id, contact.Id)
		assert.Equal(t, my_email, contact.Emails[0].Email)
		assert.Equal(t, my_phone, contact.PhoneNumbers[0].E164)

	}
}

func Test_createContactWithEmail(t *testing.T) {
	_, graphqlClient, resolver, ctx := NewWebServer(t)
	service := &CustomerOSService{graphqlClient: graphqlClient,
		conf: &config.Config{

			Service: struct {
				ServerPort       int    `env:"MESSAGE_STORE_SERVER_PORT,required"`
				CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
				CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
			}{CustomerOsAPIKey: "Hello World"},
		},
	}

	my_email := "x@x.org"
	my_id := "my_contact_id"

	resolver.ContactCreate = func(ctx context.Context, input model.ContactInput) (*model.Contact, error) {

		if !assert.Equal(t, input.Email.Email, "x@x.org") {
			return nil, status.Error(500, "Email")
		}
		return &model.Contact{
			ID: my_id,
		}, nil
	}

	resolver.GetContactById = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, my_id) {
			return nil, status.Error(500, "Unexpected id")
		}
		return &model.Contact{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Source:    "manual",
		}, nil
	}
	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{&model.Email{Email: &my_email}}, nil
	}

	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}

		return []*model.PhoneNumber{}, nil
	}

	result, err := service.CreateContactWithEmail(ctx, "openline", my_email)
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, my_id, result.Id)
	assert.Equal(t, my_email, result.Emails[0].Email)

}

func Test_createContactWithPhone(t *testing.T) {
	_, graphqlClient, resolver, ctx := NewWebServer(t)
	service := &CustomerOSService{graphqlClient: graphqlClient,
		conf: &config.Config{

			Service: struct {
				ServerPort       int    `env:"MESSAGE_STORE_SERVER_PORT,required"`
				CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
				CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
			}{CustomerOsAPIKey: "Hello World"},
		},
	}

	my_phone := "+123456"
	my_id := "my_contact_id"

	resolver.ContactCreate = func(ctx context.Context, input model.ContactInput) (*model.Contact, error) {

		if !assert.Equal(t, input.PhoneNumber.PhoneNumber, my_phone) {
			return nil, status.Error(500, "Phone number")
		}
		return &model.Contact{
			ID: my_id,
		}, nil
	}

	resolver.GetContactById = func(ctx context.Context, id string) (*model.Contact, error) {
		if !assert.Equal(t, id, my_id) {
			return nil, status.Error(500, "Unexpected id")
		}
		return &model.Contact{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Source:    "manual",
		}, nil
	}
	resolver.EmailsByContact = func(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}
		return []*model.Email{}, nil
	}

	resolver.PhoneNumbersByContact = func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
		if !assert.Equal(t, obj.ID, my_id) {
			return nil, status.Error(500, "id incorrect")
		}

		return []*model.PhoneNumber{&model.PhoneNumber{E164: &my_phone}}, nil
	}

	result, err := service.CreateContactWithPhone(ctx, "openline", my_phone)
	if err != nil {
		log.Fatalf("Got an error: %s", err.Error())
	}
	assert.Equal(t, my_id, result.Id)
	assert.Equal(t, my_phone, result.PhoneNumbers[0].E164)

}
