package service

import (
	"fmt"
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"golang.org/x/net/context"
)

type CustomerOSService struct {
	graphqlClient *graphql.Client
	conf          *c.Config
}

type CustomerOSServiceInterface interface {
	CreateInteractionEvent(ctx context.Context, options ...EventOption) (*InteractionEventCreateResponse, error)
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

type InteractionSessionParticipantInput struct {
	Email           *string `json:"email,omitempty"`
	PhoneNumber     *string `json:"phoneNumber,omitempty"`
	ContactID       *string `json:"contactID,omitempty"`
	UserID          *string `json:"userID,omitempty"`
	ParticipantType *string `json:"type,omitempty"`
}

type InteractionEventParticipant struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	RawEmail       string `json:"rawEmail,omitempty"`
	FirstName      string `json:"firstName,omitempty"`
	RawPhoneNumber string `json:"rawPhoneNumber,omitempty"`
}
type AnalysisDescriptionInput struct {
	InteractionEventId   *string `json:"interactionEventId,omitempty"`
	InteractionSessionId *string `json:"interactionSessionId,omitempty"`
}

type InteractionEventCreateResponse struct {
	InteractionEventCreate struct {
		Id     string `json:"id"`
		SentBy []struct {
			Typename         string `json:"__typename"`
			EmailParticipant struct {
				Id       string `json:"id"`
				RawEmail string `json:"rawEmail"`
			} `json:"emailParticipant"`
			PhoneNumberParticipant struct {
				ID             string `json:"id"`
				RawPhoneNumber string `json:"rawPhoneNumber"`
			} `json:"phoneNumberParticipant"`
			Type string `json:"type"`
		} `json:"sentBy"`
		SentTo []struct {
			Typename         string `json:"__typename"`
			EmailParticipant struct {
				Id       string `json:"id"`
				RawEmail string `json:"rawEmail"`
			} `json:"emailParticipant"`
			PhoneNumberParticipant struct {
				ID             string `json:"id"`
				RawPhoneNumber string `json:"rawPhoneNumber"`
			} `json:"phoneNumberParticipant"`
			Type string `json:"type"`
		} `json:"sentTo"`
	} `json:"interactionEvent_Create"`
}

func (s *CustomerOSService) addHeadersToGraphRequest(req *graphql.Request, ctx context.Context, tenant *string, user *string) error {
	req.Header.Add("X-Openline-API-KEY", s.conf.Service.CustomerOsAPIKey)
	if user != nil {
		req.Header.Add("X-Openline-USERNAME", *user)
	}
	if tenant != nil {
		req.Header.Add("X-Openline-TENANT", tenant)
	}

	return nil
}

func (s *CustomerOSService) CreateInteractionEvent(ctx context.Context, options ...EventOption) (*InteractionEventCreateResponse, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateInteractionEvent(
				$sessionId: ID!, 
				$channel: String,
				$sentBy: [InteractionEventParticipantInput!]!, 
				$sentTo: [InteractionEventParticipantInput!]!, 
				$appSource: String!, 
				$repliesTo: ID, 
				$content: String, 
				$contentType: String) {
  					interactionEvent_Create(
    					event: {interactionSession: $sessionId, 
								channel: $channel, 
								sentBy: $sentBy, 
								sentTo: $sentTo, 
								appSource: $appSource, 
								repliesTo: $repliesTo, 
								content: $content, 
								contentType: $contentType}
  					) {
						id
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

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx, params.tenant, params.username)

	if err != nil {
		return nil, fmt.Errorf("CreateContactWithPhone: %w", err)
	}

	var graphqlResponse InteractionEventCreateResponse
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("CreateInteractionEvent: %w", err)
	}

	return &graphqlResponse, nil
}

func (s *CustomerOSService) GetInteractionSession(ctx context.Context, sessionIdentifier string, tenant *string, user *string) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`query GetInteractionSession($sessionIdentifier: String!) {
  					interactionSession_BySessionIdentifier(sessionIdentifier: $sessionIdentifier) {
       					id
				}
			}`)

	graphqlRequest.Var("sessionIdentifier", sessionIdentifier)

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx, tenant, user)

	if err != nil {
		return nil, fmt.Errorf("GetInteractionSession: %w", err)
	}

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("GetInteractionSession: %w", err)
	}
	id := graphqlResponse["interactionSession_BySessionIdentifier"]["id"]
	return &id, nil
}

func (s *CustomerOSService) CreateInteractionSession(ctx context.Context, options ...SessionOption) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateInteractionSession($sessionIdentifier: String, $channel: String, $name: String!, $status: String!, $appSource: String!, $attendedBy: [InteractionSessionParticipantInput!]) {
				interactionSession_Create(
				session: {
					sessionIdentifier: $sessionIdentifier
        			channel: $channel
        			name: $name
        			status: $status
        			appSource: $appSource
                    attendedBy: $attendedBy
    			}
  			) {
				id
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

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx, params.tenant, params.username)

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

func (s *CustomerOSService) CreateAnalysis(ctx context.Context, options ...AnalysisOption) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateAnalysis($content: String, $contentType: String, $analysisType: String, $appSource: String!, $describes: [AnalysisDescriptionInput!]) {
				analysis_Create(
					analysis: {
						content: $content
						contentType: $contentType
						analysisType: $analysisType
						describes: $describes
						appSource: $appSource
					}
				  ) {
					  id
				}
			}
	`)

	params := AnalysisOptions{}
	for _, opt := range options {
		opt(&params)
	}

	graphqlRequest.Var("content", params.content)
	graphqlRequest.Var("contentType", params.contentType)
	graphqlRequest.Var("analysisType", params.analysisType)
	graphqlRequest.Var("appSource", params.appSource)

	if params.describes != nil {
		graphqlRequest.Var("describes", params.describes)
	}

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx, params.tenant, params.username)

	if err != nil {
		return nil, fmt.Errorf("CreateAnalysis: %w", err)
	}

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("CreateAnalysis: %w", err)
	}
	id := graphqlResponse["analysis_Create"]["id"]
	return &id, nil

}

type EventOptions struct {
	tenant      *string
	username    *string
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
	tenant            *string
	username          *string
	sessionIdentifier string
	attendedBy        []InteractionSessionParticipantInput
}

type AnalysisOptions struct {
	analysisType string
	content      string
	contentType  string
	appSource    string
	tenant       *string
	username     *string
	describes    *AnalysisDescriptionInput
}

type EventOption func(*EventOptions)
type SessionOption func(*SessionOptions)
type AnalysisOption func(options *AnalysisOptions)

func WithTenant(value string) EventOption {
	return func(options *EventOptions) {
		options.tenant = &value
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
		options.tenant = &value
	}
}

func WithAnalysisType(value string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.analysisType = value
	}
}

func WithAnalysisContent(value string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.content = value
	}
}

func WithAnalysisContentType(value string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.contentType = value
	}
}

func WithAnalysisAppSource(value string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.appSource = value
	}
}

func WithAnalysisTenant(value string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.tenant = &value
	}
}

func WithAnalysisDescribes(value *AnalysisDescriptionInput) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.describes = value
	}
}

func NewCustomerOSService(graphqlClient *graphql.Client, config *c.Config) *CustomerOSService {
	customerOsService := new(CustomerOSService)
	customerOsService.graphqlClient = graphqlClient
	customerOsService.conf = config
	return customerOsService
}
