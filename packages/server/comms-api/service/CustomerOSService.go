package service

import (
	"fmt"
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	"golang.org/x/net/context"
)

type CustomerOSService struct {
	graphqlClient *graphql.Client
	conf          *c.Config
}

type CustomerOSServiceInterface interface {
	CreateAnalysis(options ...AnalysisOption) (*string, error)
	CreateInteractionEvent(options ...EventOption) (*model.InteractionEventCreateResponse, error)
	CreateInteractionSession(options ...SessionOption) (*string, error)

	GetInteractionEvent(interactionEventId *string, user *string) (*model.InteractionEventGetResponse, error)
	GetInteractionSession(sessionIdentifier *string, tenant *string, user *string) (*string, error)
}

func (s *CustomerOSService) GetInteractionEvent(interactionEventId *string, user *string) (*model.InteractionEventGetResponse, error) {
	graphqlRequest := graphql.NewRequest(
		`query GetInteractionEvent($id: ID!) {
			interactionEvent(id: $id) {
				eventIdentifier,
				channelData,
				subject,
				sessionId
			}
		}`)
	graphqlRequest.Var("id", interactionEventId)

	err := s.addHeadersToGraphRequest(graphqlRequest, nil, user)

	if err != nil {
		return nil, fmt.Errorf("GetInteractionSession: %w", err)
	}

	var graphqlResponse model.InteractionEventGetResponse
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("GetInteractionSession: %w", err)
	}

	return &graphqlResponse, nil
}

func (s *CustomerOSService) addHeadersToGraphRequest(req *graphql.Request, tenant *string, user *string) error {
	req.Header.Add("X-Openline-API-KEY", s.conf.Service.CustomerOsAPIKey)
	if user != nil {
		req.Header.Add("X-Openline-USERNAME", *user)
	}
	if tenant != nil {
		req.Header.Add("X-Openline-TENANT", *tenant)
	}

	return nil
}

func (s *CustomerOSService) CreateInteractionEvent(options ...EventOption) (*model.InteractionEventCreateResponse, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateInteractionEvent(
				$sessionId: ID!, 
				$eventIdentifier: String,
				$channel: String,
				$channelData: String,
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
	graphqlRequest.Var("eventIdentifier", params.eventIdentifier)
	graphqlRequest.Var("content", params.content)
	graphqlRequest.Var("contentType", params.contentType)
	graphqlRequest.Var("channelData", params.channelData)
	graphqlRequest.Var("channel", params.channel)
	graphqlRequest.Var("sentBy", params.sentBy)
	graphqlRequest.Var("sentTo", params.sentTo)
	graphqlRequest.Var("appSource", params.appSource)

	err := s.addHeadersToGraphRequest(graphqlRequest, params.tenant, params.username)

	if err != nil {
		return nil, fmt.Errorf("CreateContactWithPhone: %w", err)
	}

	var graphqlResponse model.InteractionEventCreateResponse
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("CreateInteractionEvent: %w", err)
	}

	return &graphqlResponse, nil
}

func (s *CustomerOSService) GetInteractionSession(sessionIdentifier *string, tenant *string, user *string) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`query GetInteractionSession($sessionIdentifier: String!) {
  					interactionSession_BySessionIdentifier(sessionIdentifier: $sessionIdentifier) {
       					id
				}
			}`)

	graphqlRequest.Var("sessionIdentifier", sessionIdentifier)

	err := s.addHeadersToGraphRequest(graphqlRequest, tenant, user)

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

func (s *CustomerOSService) CreateInteractionSession(options ...SessionOption) (*string, error) {
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
	graphqlRequest.Var("attendedBy", params.attendedBy)

	err := s.addHeadersToGraphRequest(graphqlRequest, params.tenant, params.username)

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

func (s *CustomerOSService) CreateAnalysis(options ...AnalysisOption) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateAnalysis($content: String, $contentType: String, $analysisType: String, $appSource: String!, $describes: [AnalysisDescriptionInput!]!) {
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

	err := s.addHeadersToGraphRequest(graphqlRequest, params.tenant, params.username)

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
	tenant          *string
	username        *string
	sessionId       *string
	eventIdentifier *string
	repliesTo       *string
	content         *string
	contentType     *string
	channel         *string
	channelData     *string
	sentBy          []model.InteractionEventParticipantInput
	sentTo          []model.InteractionEventParticipantInput
	appSource       *string
}

type SessionOptions struct {
	channel           *string
	name              *string
	status            *string
	appSource         *string
	tenant            *string
	username          *string
	sessionIdentifier *string
	attendedBy        []model.InteractionSessionParticipantInput
}

type AnalysisOptions struct {
	analysisType *string
	content      *string
	contentType  *string
	appSource    *string
	tenant       *string
	username     *string
	describes    *model.AnalysisDescriptionInput
}

type EventOption func(*EventOptions)
type SessionOption func(*SessionOptions)
type AnalysisOption func(options *AnalysisOptions)

func (s *CustomerOSService) WithTenant(value *string) EventOption {
	return func(options *EventOptions) {
		options.tenant = value
	}
}

func (s *CustomerOSService) WithUsername(value *string) EventOption {
	return func(options *EventOptions) {
		options.username = value
	}
}

func (s *CustomerOSService) WithSessionId(value *string) EventOption {
	return func(options *EventOptions) {
		options.sessionId = value
	}
}

func (s *CustomerOSService) WithRepliesTo(value *string) EventOption {
	return func(options *EventOptions) {
		options.repliesTo = value
	}
}

func (s *CustomerOSService) WithContent(value *string) EventOption {
	return func(options *EventOptions) {
		options.content = value
	}
}

func (s *CustomerOSService) WithContentType(value *string) EventOption {
	return func(options *EventOptions) {
		options.contentType = value
	}
}

func (s *CustomerOSService) WithChannel(value *string) EventOption {
	return func(options *EventOptions) {
		options.channel = value
	}
}

func (s *CustomerOSService) WithSentBy(value []model.InteractionEventParticipantInput) EventOption {
	return func(options *EventOptions) {
		options.sentBy = value
	}
}

func (s *CustomerOSService) WithSentTo(value []model.InteractionEventParticipantInput) EventOption {
	return func(options *EventOptions) {
		options.sentTo = value
	}
}

func (s *CustomerOSService) WithSessionIdentifier(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.sessionIdentifier = value
	}
}

func (s *CustomerOSService) WithSessionChannel(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.channel = value
	}
}

func (s *CustomerOSService) WithSessionName(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.name = value
	}
}

func (s *CustomerOSService) WithSessionStatus(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.status = value
	}
}

func (s *CustomerOSService) WithSessionAttendedBy(value []model.InteractionSessionParticipantInput) SessionOption {
	return func(options *SessionOptions) {
		options.attendedBy = value
	}
}

func (s *CustomerOSService) WithSessionAppSource(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.appSource = value
	}
}

func (s *CustomerOSService) WithSessionTenant(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.tenant = value
	}
}

func (s *CustomerOSService) WithSessionUsername(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.username = value
	}
}

func (s *CustomerOSService) WithAnalysisType(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.analysisType = value
	}
}

func (s *CustomerOSService) WithAnalysisContent(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.content = value
	}
}

func (s *CustomerOSService) WithAnalysisContentType(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.contentType = value
	}
}

func (s *CustomerOSService) WithAnalysisAppSource(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.appSource = value
	}
}

func (s *CustomerOSService) WithAnalysisTenant(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.tenant = value
	}
}

func (s *CustomerOSService) WithAnalysisUsername(value *string) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.username = value
	}
}

func (s *CustomerOSService) WithAnalysisDescribes(value *model.AnalysisDescriptionInput) AnalysisOption {
	return func(options *AnalysisOptions) {
		options.describes = value
	}
}

func (s *CustomerOSService) WithEventIdentifier(eventIdentifier string) EventOption {
	return func(options *EventOptions) {
		options.eventIdentifier = &eventIdentifier
	}
}

func (s *CustomerOSService) WithChannelData(ChannelData *string) EventOption {
	return func(options *EventOptions) {
		options.channelData = ChannelData
	}
}

func NewCustomerOSService(graphqlClient *graphql.Client, config *c.Config) *CustomerOSService {
	customerOsService := new(CustomerOSService)
	customerOsService.graphqlClient = graphqlClient
	customerOsService.conf = config
	return customerOsService
}
