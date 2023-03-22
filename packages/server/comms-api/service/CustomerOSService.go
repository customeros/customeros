package service

import (
	"encoding/json"
	"fmt"
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	commonModuleService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"golang.org/x/net/context"
)

type CustomerOSService struct {
	graphqlClient *graphql.Client
	conf          *c.Config
}

type CustomerOSServiceInterface interface {
	CreateInteractionEvent(ctx context.Context, options ...EventOption) (*string, error)
	CreateInteractionSession(ctx context.Context, options ...SessionOption) (*string, error)
	GetInteractionSession(ctx context.Context, sessionIdentifier string, tenant string) (*string, error)
}

type InteractionEventParticipantInput struct {
	Email           *string `json:"email,omitempty"`
	PhoneNumber     *string `json:"phoneNumber,omitempty"`
	ContactID       *string `json:"contactID,omitempty"`
	UserID          *string `json:"userID,omitempty"`
	ParticipantType *string `json:"type,omitempty"`
}

func (s *CustomerOSService) addHeadersToGraphRequest(req *graphql.Request, ctx context.Context, tenant string) error {
	req.Header.Add("X-Openline-API-KEY", s.conf.Service.CustomerOsAPIKey)
	user, err := commonModuleService.GetUsernameMetadataForGRPC(ctx)
	if err != nil && user != nil {
		req.Header.Add("X-Openline-USERNAME", *user)
	}

	req.Header.Add("X-Openline-TENANT", tenant)
	return nil
}

func (s *CustomerOSService) CreateInteractionEvent(ctx context.Context, options ...EventOption) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateInteractionEvent(
				$sessionId: ID!, $channel: String, $sentBy: [InteractionEventParticipantInput!]!, $sentTo: [InteractionEventParticipantInput!]!, $appSource: String!, $repliesTo: ID, $content: String, $contentType: String) {
  				interactionEvent_Create(
    				event: {interactionSession: $sessionId, channel: $channel, sentBy: $sentBy, sentTo: $sentTo, appSource: $appSource, repliesTo: $repliesTo, content: $content, contentType: $contentType}
  				) {
    			id
    			createdAt
    			content
    			contentType
    			channel
    			interactionSession {
      				id
      				startedAt
      				endedAt
      				sessionIdentifier
      				name
      				status
      				type
      				channel
      				source
      				sourceOfTruth
      				appSource
    			}
				sentBy {
      				__typename
      				... on EmailParticipant {
        			emailParticipant {
          				id
          				rawEmail
        			}
        			type
				}
			    	... on UserParticipant {
        			userParticipant {
          				id
          				firstName
        			}
        			type
      			}
      				... on PhoneNumberParticipant {
        			phoneNumberParticipant {
          				id
          				rawPhoneNumber
        		}
        			type
      			}
      				... on ContactParticipant {
        			contactParticipant {
          				id
          				firstName
        		}
        			type
      			}
    		}
    		sentTo {
      			__typename
				... on EmailParticipant {
        			emailParticipant {
          				id
          				rawEmail
        			}
        			type
      			}
      			... on UserParticipant {
        			userParticipant {
          				id
          				firstName
        			}
        			type
      			}
      			... on PhoneNumberParticipant {
        			phoneNumberParticipant {
          				id
          				rawPhoneNumber
        			}
        			type
      			}
      			... on ContactParticipant {
        			contactParticipant {
          				id
          				firstName
        			}
        			type
      			}
    		}
    		repliesTo {
      			id
      			eventIdentifier
				content
      			contentType
      			channel
    		}
    		source
    		sourceOfTruth
    		appSource
  		}
	}`)

	params := EventOptions{}
	for _, opt := range options {
		opt(&params)
	}

	graphqlRequest.Var("sessionId", params.sessionId)
	graphqlRequest.Var("content", params.content)
	graphqlRequest.Var("contentType", params.contentType)
	graphqlRequest.Var("channel", params.channel)
	graphqlRequest.Var("sentBy", params.sentBy)
	graphqlRequest.Var("sentTo", params.sentTo)
	graphqlRequest.Var("appSource", params.appSource)

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx, params.tenant)

	if err != nil {
		return nil, fmt.Errorf("CreateContactWithPhone: %w", err)
	}

	var graphqlResponse interface{}
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("CreateInteractionEvent: %w", err)
	}
	//id := graphqlResponse["contact_Create"]["id"]
	bytes, _ := json.Marshal(graphqlResponse)
	print(string(bytes))
	//log.Printf("CreateContactWithPhone: phoneNumber=%s graphqlResponse = %s", phoneNumber, bytes)
	//return s.GetContactById(ctx, id, tenant)
	return nil, nil
}

func (s *CustomerOSService) GetInteractionSession(ctx context.Context, sessionIdentifier string, tenant string) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`query GetInteractionSession($sessionIdentifier: String!) {
  					interactionSession_BySessionIdentifier(sessionIdentifier: $sessionIdentifier) {
       					id
       					startedAt
       					endedAt
       					sessionIdentifier
       					name
       					status
       					type
       					channel
       					source
       					sourceOfTruth
       					appSource
       					events {
         					id
         					createdAt
         					eventIdentifier
         					content
         					contentType
         					channel
         					source
         					sourceOfTruth
         					appSource
							repliesTo {
           						id
           						eventIdentifier
           						content
           						contentType
           						channel
         					}
       					}
  					}
				}`)

	graphqlRequest.Var("sessionId", sessionIdentifier)

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx, tenant)

	if err != nil {
		return nil, fmt.Errorf("CreateContactWithPhone: %w", err)
	}

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("CreateInteractionEvent: %w", err)
	}
	id := graphqlResponse["contact_Create"]["id"]
	return &id, nil
}

func (s *CustomerOSService) CreateInteractionSession(ctx context.Context, options ...SessionOption) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateInteractionSession($sessionIdentifier: String, $channel: String, $name: String!, $status: String!, $appSource: String!) {
				interactionSession_Create(
				session: {
					sessionIdentifier: $sessionIdentifier
        			channel: $channel
        			name: $name
        			status: $status
        			appSource: $appSource
    			}
  			) {
				id
       			createdAt
       			updatedAt
       			sessionIdentifier
       			name
       			status
       			type
       			channel
       			source
       			sourceOfTruth
       			appSource
   			}
		}
	`)

	params := SessionOptions{}
	for _, opt := range options {
		opt(&params)
	}

	graphqlRequest.Var("sessionIdentifier", params.sessionIdentifier)
	graphqlRequest.Var("channel", params.channel)
	graphqlRequest.Var("name", params.name)
	graphqlRequest.Var("status", params.status)
	graphqlRequest.Var("appSource", params.appSource)

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx, params.tenant)

	if err != nil {
		return nil, fmt.Errorf("CreateContactWithPhone: %w", err)
	}

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("CreateInteractionEvent: %w", err)
	}
	id := graphqlResponse["interactionSession_Create"]["id"]
	return &id, nil

}

type EventOptions struct {
	tenant      string
	sessionId   string
	repliesTo   string
	content     string
	contentType string
	channel     string
	sentBy      []InteractionEventParticipantInput
	sentTo      []InteractionEventParticipantInput
	appSource   string
}

type SessionOptions struct {
	channel           string
	name              string
	status            string
	appSource         string
	tenant            string
	sessionIdentifier string
}

type EventOption func(*EventOptions)
type SessionOption func(*SessionOptions)

func WithTenant(value string) EventOption {
	return func(options *EventOptions) {
		options.tenant = value
	}
}

func WithSessionId(value string) EventOption {
	return func(options *EventOptions) {
		options.sessionId = value
	}
}

func WithRepliesTo(value string) EventOption {
	return func(options *EventOptions) {
		options.repliesTo = value
	}
}

func WithContent(value string) EventOption {
	return func(options *EventOptions) {
		options.content = value
	}
}

func WithContentType(value string) EventOption {
	return func(options *EventOptions) {
		options.contentType = value
	}
}

func WithChannel(value string) EventOption {
	return func(options *EventOptions) {
		options.channel = value
	}
}

func WithSentBy(value []InteractionEventParticipantInput) EventOption {
	return func(options *EventOptions) {
		options.sentBy = value
	}
}

func WithSentTo(value []InteractionEventParticipantInput) EventOption {
	return func(options *EventOptions) {
		options.sentTo = value
	}
}

func WithSessionIdentifier(value string) SessionOption {
	return func(options *SessionOptions) {
		options.sessionIdentifier = value
	}
}

func WithSessionChannel(value string) SessionOption {
	return func(options *SessionOptions) {
		options.channel = value
	}
}

func WithSessionName(value string) SessionOption {
	return func(options *SessionOptions) {
		options.name = value
	}
}

func WithSessionStatus(value string) SessionOption {
	return func(options *SessionOptions) {
		options.status = value
	}
}

func WithSessionAppSource(value string) SessionOption {
	return func(options *SessionOptions) {
		options.appSource = value
	}
}
func WithSessionTenant(value string) SessionOption {
	return func(options *SessionOptions) {
		options.tenant = value
	}
}

func NewCustomerOSService(graphqlClient *graphql.Client, config *c.Config) *CustomerOSService {
	customerOsService := new(CustomerOSService)
	customerOsService.graphqlClient = graphqlClient
	customerOsService.conf = config
	return customerOsService
}
