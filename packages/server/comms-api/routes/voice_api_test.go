package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	voiceModel "github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var myVoiceApiConfig = &config.Config{
	Service: struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
		ServerAddress    string `env:"COMMS_API_SERVER_ADDRESS,required"`
		CorsUrl          string `env:"COMMS_API_CORS_URL,required"`
	}{CustomerOsAPIKey: "my-key"},
	VCon: struct {
		ApiKey          string `env:"COMMS_API_VCON_API_KEY,required"`
		AwsAccessKey    string `env:"AWS_ACCESS_KEY"`
		AwsAccessSecret string `env:"AWS_ACCESS_SECRET"`
		AwsRegion       string `env:"AWS_REGION"`
		AwsBucket       string `env:"AWS_BUCKET"`
	}{ApiKey: "my-vcon-key"},
}

func init() {
}

func Test_eventCallStarted(t *testing.T) {
	voiceApiRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVoiceApiConfig)
	route := voiceApiRouter.Group("/")
	reachedSessionCreate := false
	reachedSessionBySessionIdentifier := false
	var attendedBy []*model.InteractionSessionParticipantInput

	tenantApiKey := "my-tenant-key"
	testRedisDatabase := utils.NewTestRedisService()
	testRedisDatabase.KeyMap[tenantApiKey] = utils.DatabaseValues{
		Active: true,
		Tenant: "my-tenant",
	}
	addVoiceApiRoutes(myVconConfig, route, customerOs, nil, testRedisDatabase)

	from := "AgentSmith@openline.ai"
	to := "+32485111000"
	startTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:45.872099866Z")
	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	resolver.InteractionEventCreate = func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
		log.Printf("InteractionEventCreate: Got Event %v", event)
		var callData OpenlineCallProgressData
		err := json.Unmarshal([]byte(*event.Content), &callData)
		assert.Nil(t, err)
		assert.Equal(t, startTime, *callData.StartTime)
		assert.Equal(t, "WEBRTC", *callData.SentByType)
		assert.Equal(t, "PSTN", *callData.SentToType)
		assert.Equal(t, "application/x-openline-call-progress", *event.ContentType)
		assert.Equal(t, "VOICE", *event.Channel)
		assert.Equal(t, from, *event.SentBy[0].Email)
		assert.Equal(t, to, *event.SentTo[0].PhoneNumber)
		assert.NotNil(t, event.CreatedAt)
		assert.Equal(t, startTime, event.CreatedAt.UTC())
		assert.Equal(t, "CALL_START", *event.EventType)

		return &model.InteractionEvent{
			ID:                 "my-event-id",
			CreatedAt:          startTime.UTC(),
			EventIdentifier:    event.EventIdentifier,
			Content:            event.Content,
			ContentType:        event.ContentType,
			Channel:            event.Channel,
			ChannelData:        event.ChannelData,
			EventType:          event.EventType,
			InteractionSession: nil,
			SentBy:             nil,
			SentTo:             nil,
			RepliesTo:          nil,
			Source:             "TEST",
			SourceOfTruth:      "TEST",
			AppSource:          event.AppSource,
		}, nil
	}

	resolver.InteractionSessionBySessionIdentifier = func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionBySessionIdentifier: Got Session Identifier %s", sessionIdentifier)
		require.Equal(t, "e061697f-673d-4756-a5f7-4f114e66a191", sessionIdentifier)
		reachedSessionBySessionIdentifier = true
		return nil, fmt.Errorf("Session not found: %s", sessionIdentifier)
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)
		reachedSessionCreate = true
		attendedBy = session.AttendedBy
		require.Equal(t, "e061697f-673d-4756-a5f7-4f114e66a191", *session.SessionIdentifier)
		require.Equal(t, fmt.Sprintf("Outgoing call to %s", to), session.Name)
		require.Equal(t, "ACTIVE", session.Status)
		require.Equal(t, "CALL", *session.Type)
		require.Equal(t, "VOICE", *session.Channel)
		require.Equal(t, from, *session.AttendedBy[0].Email)
		require.Equal(t, to, *session.AttendedBy[1].PhoneNumber)

		return &model.InteractionSession{
			ID:                "my-new-session-id",
			StartedAt:         time.Now().UTC(),
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
			SessionIdentifier: session.SessionIdentifier,
			Name:              session.Name,
			Status:            session.Status,
			Type:              session.Type,
			Channel:           session.Channel,
			ChannelData:       session.ChannelData,
			Source:            "TEST",
			SourceOfTruth:     "TEST",
			AppSource:         session.AppSource,
			Events:            nil,
			AttendedBy:        nil,
		}, nil
	}
	resolver.InteractionSessionResolver = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionResolver: Got Event %v", obj)
		return &model.InteractionSession{
			Name: "my-session",
		}, nil
	}
	resolver.SentBy = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentBy: Got Event %v", obj)
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.SentTo = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentTo: Got Event %v", obj)
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.RepliesTo = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionEvent, error) {
		log.Printf("RepliesTo: Got Event %v", obj)
		return &model.InteractionEvent{}, nil
	}

	resolver.AttendedBy = func(ctx context.Context, obj *model.InteractionSession) ([]model.InteractionSessionParticipant, error) {
		log.Printf("AttendedBy: Got Session %v", obj)
		return []model.InteractionSessionParticipant{}, nil
	}

	eventStart := voiceModel.CallEventStart{
		CallEvent: voiceModel.CallEvent{
			Version:       "1.0",
			CorrelationId: "e061697f-673d-4756-a5f7-4f114e66a191",
			Event:         "CALL_START",
			From:          &voiceModel.CallEventParty{Mailto: &from, Type: voiceModel.CALL_EVENT_TYPE_WEBTRC},
			To:            &voiceModel.CallEventParty{Tel: &to, Type: voiceModel.CALL_EVENT_TYPE_PSTN},
		},
		StartTime: startTime,
	}

	w := httptest.NewRecorder()
	msgBytes, err := json.Marshal(eventStart)
	req, _ := http.NewRequest("POST", "/call_progress", bytes.NewReader(msgBytes))
	req.Header.Add("X-API-KEY", tenantApiKey)
	voiceApiRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, 200, w.Code) {
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.Nil(t, err)
	ids, ok := result["ids"]
	assert.True(t, ok)
	createdIds, ok := ids.([]interface{})
	assert.True(t, ok)
	assert.Len(t, createdIds, 1)

	assert.True(t, reachedSessionCreate)
	assert.True(t, reachedSessionBySessionIdentifier)
	assert.Len(t, attendedBy, 2)

}

func Test_eventCallAnswered(t *testing.T) {
	voiceApiRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVoiceApiConfig)
	route := voiceApiRouter.Group("/")
	reachedSessionBySessionIdentifier := false

	tenantApiKey := "my-tenant-key"
	testRedisDatabase := utils.NewTestRedisService()
	testRedisDatabase.KeyMap[tenantApiKey] = utils.DatabaseValues{
		Active: true,
		Tenant: "my-tenant",
	}
	addVoiceApiRoutes(myVconConfig, route, customerOs, nil, testRedisDatabase)

	from := "AgentSmith@openline.ai"
	to := "+32485111000"
	sip := "test001@openline.ai"
	startTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:45.872099866Z")
	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	answerTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:46.872099866Z")
	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	resolver.InteractionEventCreate = func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
		log.Printf("InteractionEventCreate: Got Event %v", event)
		var callData OpenlineCallProgressData
		err := json.Unmarshal([]byte(*event.Content), &callData)
		assert.Nil(t, err)
		assert.Equal(t, startTime, *callData.StartTime)
		assert.Equal(t, answerTime, *callData.AnsweredTime)
		assert.Equal(t, "PSTN", *callData.SentByType)
		assert.Equal(t, "ESIM", *callData.SentToType)
		assert.Equal(t, "application/x-openline-call-progress", *event.ContentType)
		assert.Equal(t, "VOICE", *event.Channel)
		assert.Equal(t, from, *event.SentTo[0].Email)
		assert.Equal(t, to, *event.SentBy[0].PhoneNumber)
		assert.NotNil(t, event.CreatedAt)
		assert.Equal(t, answerTime, event.CreatedAt.UTC())
		assert.Equal(t, "CALL_ANSWERED", *event.EventType)

		return &model.InteractionEvent{
			ID:                 "my-event-id",
			CreatedAt:          answerTime.UTC(),
			EventIdentifier:    event.EventIdentifier,
			Content:            event.Content,
			ContentType:        event.ContentType,
			Channel:            event.Channel,
			ChannelData:        event.ChannelData,
			EventType:          event.EventType,
			InteractionSession: nil,
			SentBy:             nil,
			SentTo:             nil,
			RepliesTo:          nil,
			Source:             "TEST",
			SourceOfTruth:      "TEST",
			AppSource:          event.AppSource,
		}, nil
	}

	resolver.InteractionSessionBySessionIdentifier = func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionBySessionIdentifier: Got Session Identifier %s", sessionIdentifier)
		require.Equal(t, "e061697f-673d-4756-a5f7-4f114e66a191", sessionIdentifier)
		reachedSessionBySessionIdentifier = true
		sessionType := "CALL"
		sessionChannel := "VOICE"
		return &model.InteractionSession{
			ID:                "my-new-session-id",
			StartedAt:         startTime.UTC(),
			CreatedAt:         startTime.UTC(),
			UpdatedAt:         startTime.UTC(),
			SessionIdentifier: &sessionIdentifier,
			Name:              fmt.Sprintf("Outgoing call to %s", to),
			Status:            "ACTIVE",
			Type:              &sessionType,
			Channel:           &sessionChannel,
			Source:            "TEST",
			SourceOfTruth:     "TEST",
			AppSource:         "TEST",
			Events:            nil,
			AttendedBy:        nil,
		}, nil
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)
		require.Fail(t, "Session should not be created")
		return nil, fmt.Errorf("Session already exists: %s", *session.SessionIdentifier)
	}
	resolver.InteractionSessionResolver = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionResolver: Got Event %v", obj)
		return &model.InteractionSession{
			Name: "my-session",
		}, nil
	}
	resolver.SentBy = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentBy: Got Event %v", obj)
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.SentTo = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentTo: Got Event %v", obj)
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.RepliesTo = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionEvent, error) {
		log.Printf("RepliesTo: Got Event %v", obj)
		return &model.InteractionEvent{}, nil
	}

	resolver.AttendedBy = func(ctx context.Context, obj *model.InteractionSession) ([]model.InteractionSessionParticipant, error) {
		log.Printf("AttendedBy: Got Session %v", obj)
		return []model.InteractionSessionParticipant{}, nil
	}

	eventStart := voiceModel.CallEventAnswered{
		CallEvent: voiceModel.CallEvent{
			Version:       "1.0",
			CorrelationId: "e061697f-673d-4756-a5f7-4f114e66a191",
			Event:         "CALL_ANSWERED",
			From:          &voiceModel.CallEventParty{Mailto: &from, Type: voiceModel.CALL_EVENT_TYPE_SIP, Sip: &sip},
			To:            &voiceModel.CallEventParty{Tel: &to, Type: voiceModel.CALL_EVENT_TYPE_PSTN},
		},
		StartTime:    startTime,
		AnsweredTime: answerTime,
	}

	w := httptest.NewRecorder()
	msgBytes, err := json.Marshal(eventStart)
	req, _ := http.NewRequest("POST", "/call_progress", bytes.NewReader(msgBytes))
	req.Header.Add("X-API-KEY", tenantApiKey)
	voiceApiRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, 200, w.Code) {
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.Nil(t, err)
	ids, ok := result["ids"]
	assert.True(t, ok)
	createdIds, ok := ids.([]interface{})
	assert.True(t, ok)
	assert.Len(t, createdIds, 1)

	assert.True(t, reachedSessionBySessionIdentifier)
}

func Test_eventCallCalledHangup(t *testing.T) {
	voiceApiRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVoiceApiConfig)
	route := voiceApiRouter.Group("/")
	reachedSessionBySessionIdentifier := false

	tenantApiKey := "my-tenant-key"
	testRedisDatabase := utils.NewTestRedisService()
	testRedisDatabase.KeyMap[tenantApiKey] = utils.DatabaseValues{
		Active: true,
		Tenant: "my-tenant",
	}
	addVoiceApiRoutes(myVconConfig, route, customerOs, nil, testRedisDatabase)

	from := "AgentSmith@openline.ai"
	to := "+32485111000"
	sip := "test001@openline.ai"
	startTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:45.872099866Z")
	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	answerTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:46.872099866Z")
	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	endTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:47.872099866Z")

	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	resolver.InteractionEventCreate = func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
		log.Printf("InteractionEventCreate: Got Event %v", event)
		var callData OpenlineCallProgressData
		err := json.Unmarshal([]byte(*event.Content), &callData)
		assert.Nil(t, err)
		assert.Equal(t, startTime, *callData.StartTime)
		assert.Equal(t, answerTime, *callData.AnsweredTime)
		assert.Equal(t, endTime, *callData.EndTime)
		assert.Equal(t, endTime.Sub(answerTime).Milliseconds(), *callData.Duration)
		assert.Equal(t, "PSTN", *callData.SentByType)
		assert.Equal(t, "ESIM", *callData.SentToType)

		assert.Equal(t, "application/x-openline-call-progress", *event.ContentType)
		assert.Equal(t, "VOICE", *event.Channel)
		assert.Equal(t, from, *event.SentTo[0].Email)
		assert.Equal(t, to, *event.SentBy[0].PhoneNumber)
		assert.NotNil(t, event.CreatedAt)
		assert.Equal(t, endTime, event.CreatedAt.UTC())
		assert.Equal(t, "CALL_END", *event.EventType)

		return &model.InteractionEvent{
			ID:                 "my-event-id",
			CreatedAt:          answerTime.UTC(),
			EventIdentifier:    event.EventIdentifier,
			Content:            event.Content,
			ContentType:        event.ContentType,
			Channel:            event.Channel,
			ChannelData:        event.ChannelData,
			EventType:          event.EventType,
			InteractionSession: nil,
			SentBy:             nil,
			SentTo:             nil,
			RepliesTo:          nil,
			Source:             "TEST",
			SourceOfTruth:      "TEST",
			AppSource:          event.AppSource,
		}, nil
	}

	resolver.InteractionSessionBySessionIdentifier = func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionBySessionIdentifier: Got Session Identifier %s", sessionIdentifier)
		require.Equal(t, "e061697f-673d-4756-a5f7-4f114e66a191", sessionIdentifier)
		reachedSessionBySessionIdentifier = true
		sessionType := "CALL"
		sessionChannel := "VOICE"
		return &model.InteractionSession{
			ID:                "my-new-session-id",
			StartedAt:         startTime.UTC(),
			CreatedAt:         startTime.UTC(),
			UpdatedAt:         startTime.UTC(),
			SessionIdentifier: &sessionIdentifier,
			Name:              fmt.Sprintf("Outgoing call to %s", to),
			Status:            "ACTIVE",
			Type:              &sessionType,
			Channel:           &sessionChannel,
			Source:            "TEST",
			SourceOfTruth:     "TEST",
			AppSource:         "TEST",
			Events:            nil,
			AttendedBy:        nil,
		}, nil
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)
		require.Fail(t, "Session should not be created")
		return nil, fmt.Errorf("Session already exists: %s", *session.SessionIdentifier)
	}
	resolver.InteractionSessionResolver = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionResolver: Got Event %v", obj)
		return &model.InteractionSession{
			Name: "my-session",
		}, nil
	}
	resolver.SentBy = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentBy: Got Event %v", obj)
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.SentTo = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentTo: Got Event %v", obj)
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.RepliesTo = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionEvent, error) {
		log.Printf("RepliesTo: Got Event %v", obj)
		return &model.InteractionEvent{}, nil
	}

	resolver.AttendedBy = func(ctx context.Context, obj *model.InteractionSession) ([]model.InteractionSessionParticipant, error) {
		log.Printf("AttendedBy: Got Session %v", obj)
		return []model.InteractionSessionParticipant{}, nil
	}

	eventStart := voiceModel.CallEventEnd{
		CallEvent: voiceModel.CallEvent{
			Version:       "1.0",
			CorrelationId: "e061697f-673d-4756-a5f7-4f114e66a191",
			Event:         "CALL_END",
			From:          &voiceModel.CallEventParty{Mailto: &from, Sip: &sip, Type: voiceModel.CALL_EVENT_TYPE_SIP},
			To:            &voiceModel.CallEventParty{Tel: &to, Type: voiceModel.CALL_EVENT_TYPE_PSTN},
		},
		StartTime:    &startTime,
		AnsweredTime: &answerTime,
		EndTime:      endTime,
		Duration:     endTime.Sub(answerTime).Milliseconds(),
		FromCaller:   false,
	}

	w := httptest.NewRecorder()
	msgBytes, err := json.Marshal(eventStart)
	req, _ := http.NewRequest("POST", "/call_progress", bytes.NewReader(msgBytes))
	req.Header.Add("X-API-KEY", tenantApiKey)
	voiceApiRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, 200, w.Code) {
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.Nil(t, err)
	ids, ok := result["ids"]
	assert.True(t, ok)
	createdIds, ok := ids.([]interface{})
	assert.True(t, ok)
	assert.Len(t, createdIds, 1)

	assert.True(t, reachedSessionBySessionIdentifier)
}

func Test_eventCallCallingHangup(t *testing.T) {
	voiceApiRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVoiceApiConfig)
	route := voiceApiRouter.Group("/")
	reachedSessionBySessionIdentifier := false

	tenantApiKey := "my-tenant-key"
	testRedisDatabase := utils.NewTestRedisService()
	testRedisDatabase.KeyMap[tenantApiKey] = utils.DatabaseValues{
		Active: true,
		Tenant: "my-tenant",
	}
	addVoiceApiRoutes(myVconConfig, route, customerOs, nil, testRedisDatabase)

	from := "AgentSmith@openline.ai"
	to := "+32485111000"
	startTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:45.872099866Z")
	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	answerTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:46.872099866Z")
	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	endTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:47.872099866Z")

	if err != nil {
		assert.Fail(t, "Could not parse time %v", err)
	}
	resolver.InteractionEventCreate = func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
		log.Printf("InteractionEventCreate: Got Event %v", event)
		var callData OpenlineCallProgressData
		err := json.Unmarshal([]byte(*event.Content), &callData)
		assert.Nil(t, err)
		assert.Equal(t, startTime, *callData.StartTime)
		assert.Equal(t, answerTime, *callData.AnsweredTime)
		assert.Equal(t, endTime, *callData.EndTime)
		assert.Equal(t, endTime.Sub(answerTime).Milliseconds(), *callData.Duration)
		assert.Equal(t, "WEBRTC", *callData.SentByType)
		assert.Equal(t, "PSTN", *callData.SentToType)

		assert.Equal(t, "application/x-openline-call-progress", *event.ContentType)
		assert.Equal(t, "VOICE", *event.Channel)
		assert.Equal(t, from, *event.SentBy[0].Email)
		assert.Equal(t, to, *event.SentTo[0].PhoneNumber)
		assert.NotNil(t, event.CreatedAt)
		assert.Equal(t, endTime, event.CreatedAt.UTC())
		assert.Equal(t, "CALL_END", *event.EventType)

		return &model.InteractionEvent{
			ID:                 "my-event-id",
			CreatedAt:          answerTime.UTC(),
			EventIdentifier:    event.EventIdentifier,
			Content:            event.Content,
			ContentType:        event.ContentType,
			Channel:            event.Channel,
			ChannelData:        event.ChannelData,
			EventType:          event.EventType,
			InteractionSession: nil,
			SentBy:             nil,
			SentTo:             nil,
			RepliesTo:          nil,
			Source:             "TEST",
			SourceOfTruth:      "TEST",
			AppSource:          event.AppSource,
		}, nil
	}

	resolver.InteractionSessionBySessionIdentifier = func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionBySessionIdentifier: Got Session Identifier %s", sessionIdentifier)
		require.Equal(t, "e061697f-673d-4756-a5f7-4f114e66a191", sessionIdentifier)
		reachedSessionBySessionIdentifier = true
		sessionType := "CALL"
		sessionChannel := "VOICE"
		return &model.InteractionSession{
			ID:                "my-new-session-id",
			StartedAt:         startTime.UTC(),
			CreatedAt:         startTime.UTC(),
			UpdatedAt:         startTime.UTC(),
			SessionIdentifier: &sessionIdentifier,
			Name:              fmt.Sprintf("Outgoing call to %s", to),
			Status:            "ACTIVE",
			Type:              &sessionType,
			Channel:           &sessionChannel,
			Source:            "TEST",
			SourceOfTruth:     "TEST",
			AppSource:         "TEST",
			Events:            nil,
			AttendedBy:        nil,
		}, nil
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)
		require.Fail(t, "Session should not be created")
		return nil, fmt.Errorf("Session already exists: %s", *session.SessionIdentifier)
	}
	resolver.InteractionSessionResolver = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionResolver: Got Event %v", obj)
		return &model.InteractionSession{
			Name: "my-session",
		}, nil
	}
	resolver.SentBy = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentBy: Got Event %v", obj)
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.SentTo = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentTo: Got Event %v", obj)
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.RepliesTo = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionEvent, error) {
		log.Printf("RepliesTo: Got Event %v", obj)
		return &model.InteractionEvent{}, nil
	}

	resolver.AttendedBy = func(ctx context.Context, obj *model.InteractionSession) ([]model.InteractionSessionParticipant, error) {
		log.Printf("AttendedBy: Got Session %v", obj)
		return []model.InteractionSessionParticipant{}, nil
	}

	eventStart := voiceModel.CallEventEnd{
		CallEvent: voiceModel.CallEvent{
			Version:       "1.0",
			CorrelationId: "e061697f-673d-4756-a5f7-4f114e66a191",
			Event:         "CALL_END",
			From:          &voiceModel.CallEventParty{Mailto: &from, Type: voiceModel.CALL_EVENT_TYPE_WEBTRC},
			To:            &voiceModel.CallEventParty{Tel: &to, Type: voiceModel.CALL_EVENT_TYPE_PSTN},
		},
		StartTime:    &startTime,
		AnsweredTime: &answerTime,
		EndTime:      endTime,
		Duration:     endTime.Sub(answerTime).Milliseconds(),
		FromCaller:   true,
	}

	w := httptest.NewRecorder()
	msgBytes, err := json.Marshal(eventStart)
	req, _ := http.NewRequest("POST", "/call_progress", bytes.NewReader(msgBytes))
	req.Header.Add("X-API-KEY", tenantApiKey)
	voiceApiRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, 200, w.Code) {
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.Nil(t, err)
	ids, ok := result["ids"]
	assert.True(t, ok)
	createdIds, ok := ids.([]interface{})
	assert.True(t, ok)
	assert.Len(t, createdIds, 1)

	assert.True(t, reachedSessionBySessionIdentifier)
}
