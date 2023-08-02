package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var myVconConfig = &config.Config{
	Service: struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
		FileStoreAPI     string `env:"FILE_STORE_API,required"`
		FileStoreAPIKey  string `env:"FILE_STORE_API_KEY,required"`
		ServerAddress    string `env:"COMMS_API_SERVER_ADDRESS,required"`
		CorsUrl          string `env:"COMMS_API_CORS_URL,required"`
	}{CustomerOsAPIKey: "my-key"},
	VCon: struct {
		ApiKey string `env:"COMMS_API_VCON_API_KEY,required"`
	}{ApiKey: "my-vcon-key"},
}

func init() {
}

const LIVE_TRANSCRIPTION = `{"vcon":"","uuid":"e061697f-673d-4756-a5f7-4f114e66a191","subject":"","parties":[{"mailto":"torrey@openline.ai"},{"tel":"+32485112970"}],"dialog":[{"type":"text","start":"2023-03-27T07:11:45.872099866Z","duration":0,"parties":[0,1],"mimetype":"text/plain","body":"Alors?","encoding":"None"}]}`
const LIVE_TRANSCRIPTON_APPENDED = `{"vcon":"","uuid":"54235d4f-f566-4d41-86a3-9726083c6aff","subject":"","appended":{"uuid":"e061697f-673d-4756-a5f7-4f114e66a191"},"parties":[{"mailto":"torrey@openline.ai"},{"tel":"+32485112970"}],"dialog":[{"type":"text","start":"2023-03-27T07:11:50.724873918Z","duration":0,"parties":[1,0],"mimetype":"text/plain","body":"some real time. Exit","encoding":"None"}]}`
const POST_TRANSCRIBE = `{"vcon":"","uuid":"6a50a92a-7322-4988-a85c-437e1d25ea45","subject":"","appended":{"uuid":"e061697f-673d-4756-a5f7-4f114e66a191"},"parties":[{"mailto":"torrey@openline.ai"},{"tel":"+32485112970"}],"analysis":[{"type":"transcript","dialog":null,"mimetype":"application/x-openline-transcript","body":"[{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Hello? Give me some real time\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Hello?\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Give me some real-\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Bye\"}]","encoding":""}]}`
const POST_TRANSCRIBE_AS_FIRST_EVENT = `{"vcon":"","uuid":"6a50a92a-7322-4988-a85c-437e1d25ea45","subject":"","parties":[{"mailto":"torrey@openline.ai"},{"tel":"+32485112970"}],"analysis":[{"type":"transcript","dialog":null,"mimetype":"application/x-openline-transcript","body":"[{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Hello? Give me some real time\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Hello?\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Give me some real-\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Bye\"}]","encoding":""}]}`

func Test_invalidApiKey(t *testing.T) {
	vconRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVconConfig)
	route := vconRouter.Group("/")

	addVconRoutes(myVconConfig, route, customerOs, nil)

	resolver.InteractionEventCreate = func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
		log.Printf("InteractionEventCreate: Got Event %v", event)
		require.Fail(t, "Should not have reached InteractionEventCreate")
		return nil, nil
	}

	resolver.InteractionSessionBySessionIdentifier = func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionBySessionIdentifier: Got Session Identifier %s", sessionIdentifier)
		assert.Fail(t, "Should not have reached InteractionSessionBySessionIdentifier")
		return nil, nil
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)

		assert.Fail(t, "Should not have reached InteractionSessionCreate")

		return nil, nil
	}

	resolver.SentBy = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentBy: Got Event %v", obj)
		assert.Fail(t, "Should not have reached SentBy")
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.SentTo = func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
		log.Printf("SentTo: Got Event %v", obj)
		assert.Fail(t, "Should not have reached SentTo")
		return []model.InteractionEventParticipant{}, nil
	}

	resolver.RepliesTo = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionEvent, error) {
		log.Printf("RepliesTo: Got Event %v", obj)
		assert.Fail(t, "Should not have reached RepliesTo")
		return &model.InteractionEvent{}, nil
	}

	resolver.AttendedBy = func(ctx context.Context, obj *model.InteractionSession) ([]model.InteractionSessionParticipant, error) {
		log.Printf("AttendedBy: Got Session %v", obj)
		assert.Fail(t, "Should not have reached AttendedBy")
		return []model.InteractionSessionParticipant{}, nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/vcon", strings.NewReader(LIVE_TRANSCRIPTION))
	req.Header.Add("X-Openline-VCon-Api-Key", "invalid")
	vconRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, w.Code, 403) {
		return
	}

}

func Test_vConDialogEvent(t *testing.T) {
	vconRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVconConfig)
	route := vconRouter.Group("/")
	reachedSessionCreate := false
	reachedSessionBySessionIdentifier := false
	var attendedBy []*model.InteractionSessionParticipantInput

	addVconRoutes(myVconConfig, route, customerOs, nil)

	resolver.InteractionEventCreate = func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
		log.Printf("InteractionEventCreate: Got Event %v", event)
		assert.Equal(t, "Alors?", *event.Content)
		assert.Equal(t, "text/plain", *event.ContentType)
		assert.Equal(t, "VOICE", *event.Channel)
		assert.Equal(t, "torrey@openline.ai", *event.SentBy[0].Email)
		assert.Equal(t, "+32485112970", *event.SentTo[0].PhoneNumber)
		msgTime, err := time.Parse(time.RFC3339, "2023-03-27T07:11:45.872099866Z")
		if err != nil {
			assert.Fail(t, "Could not parse time %v", err)
		}
		assert.NotNil(t, event.CreatedAt)
		assert.Equal(t, msgTime, event.CreatedAt.UTC())

		return &model.InteractionEvent{
			ID:                 "my-event-id",
			CreatedAt:          time.Now().UTC(),
			EventIdentifier:    event.EventIdentifier,
			Content:            event.Content,
			ContentType:        event.ContentType,
			Channel:            event.Channel,
			ChannelData:        event.ChannelData,
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
		require.Equal(t, "Outgoing call to +32485112970", session.Name)
		require.Equal(t, "ACTIVE", session.Status)
		require.Equal(t, "CALL", *session.Type)
		require.Equal(t, "VOICE", *session.Channel)
		require.Equal(t, "torrey@openline.ai", *session.AttendedBy[0].Email)
		require.Equal(t, "+32485112970", *session.AttendedBy[1].PhoneNumber)

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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/vcon", strings.NewReader(LIVE_TRANSCRIPTION))
	req.Header.Add("X-Openline-VCon-Api-Key", myVconConfig.VCon.ApiKey)
	vconRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, 200, w.Code) {
		return
	}

	assert.True(t, reachedSessionCreate)
	assert.True(t, reachedSessionBySessionIdentifier)
	assert.Len(t, attendedBy, 2)
}

func Test_vConDialogEventInExistingSession(t *testing.T) {
	vconRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVconConfig)
	route := vconRouter.Group("/")
	reachedSessionCreate := false
	reachedSessionBySessionIdentifier := false

	addVconRoutes(myVconConfig, route, customerOs, nil)

	resolver.InteractionEventCreate = func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
		log.Printf("InteractionEventCreate: Got Event %v", event)
		assert.Equal(t, "some real time. Exit", *event.Content)
		assert.Equal(t, "text/plain", *event.ContentType)
		assert.Equal(t, "VOICE", *event.Channel)
		assert.Equal(t, "torrey@openline.ai", *event.SentTo[0].Email)
		assert.Equal(t, "+32485112970", *event.SentBy[0].PhoneNumber)

		return &model.InteractionEvent{
			ID:                 "my-event-id",
			CreatedAt:          time.Now().UTC(),
			EventIdentifier:    event.EventIdentifier,
			Content:            event.Content,
			ContentType:        event.ContentType,
			Channel:            event.Channel,
			ChannelData:        event.ChannelData,
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
		SESSION_TYPE := "CALL"
		SESSION_CHANNEL := "VOICE"
		return &model.InteractionSession{
			ID:                "my-new-session-id",
			StartedAt:         time.Now().UTC(),
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
			SessionIdentifier: &sessionIdentifier,
			Name:              "Outgoing call to +32485112970",
			Status:            "ACTIVE",
			Type:              &SESSION_TYPE,
			Channel:           &SESSION_CHANNEL,
			ChannelData:       nil,
			Source:            "TEST",
			SourceOfTruth:     "TEST",
			AppSource:         "COMMS-API",
		}, nil
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)
		reachedSessionCreate = true
		require.Fail(t, "Interaction Session Create should not be called!")
		return nil, fmt.Errorf("interaction Session Create should not be called")
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

	resolver.InteractionSessionResolver = func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionResolver: Got Event %v", obj)
		return &model.InteractionSession{
			Name: "my-session-name",
		}, nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/vcon", strings.NewReader(LIVE_TRANSCRIPTON_APPENDED))
	req.Header.Add("X-Openline-VCon-Api-Key", myVconConfig.VCon.ApiKey)
	vconRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, w.Code, 200) {
		return
	}

	assert.False(t, reachedSessionCreate)
	assert.True(t, reachedSessionBySessionIdentifier)
}

func Test_vConAnalysisEventInExistingSession(t *testing.T) {
	vconRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVconConfig)
	route := vconRouter.Group("/")
	reachedSessionCreate := false
	reachedSessionBySessionIdentifier := false

	addVconRoutes(myVconConfig, route, customerOs, nil)

	resolver.AnalysisCreate = func(ctx context.Context, analysis model.AnalysisInput) (*model.Analysis, error) {
		log.Printf("InteractionEventCreate: Got Event %v", analysis)
		assert.Equal(t, "[{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Hello? Give me some real time\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Hello?\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Give me some real-\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Bye\"}]", *analysis.Content)
		assert.Equal(t, "application/x-openline-transcript", *analysis.ContentType)
		assert.Equal(t, "transcript", *analysis.AnalysisType)

		return &model.Analysis{
			ID:            "my-analysis-id",
			CreatedAt:     time.Now().UTC(),
			Content:       analysis.Content,
			ContentType:   analysis.ContentType,
			AnalysisType:  analysis.AnalysisType,
			Source:        "TEST",
			SourceOfTruth: "TEST",
			AppSource:     analysis.AppSource,
		}, nil
	}

	resolver.InteractionSessionBySessionIdentifier = func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionBySessionIdentifier: Got Session Identifier %s", sessionIdentifier)
		require.Equal(t, "e061697f-673d-4756-a5f7-4f114e66a191", sessionIdentifier)
		reachedSessionBySessionIdentifier = true
		SESSION_TYPE := "CALL"
		SESSION_CHANNEL := "VOICE"
		return &model.InteractionSession{
			ID:                "my-new-session-id",
			StartedAt:         time.Now().UTC(),
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
			SessionIdentifier: &sessionIdentifier,
			Name:              "Outgoing call to +32485112970",
			Status:            "ACTIVE",
			Type:              &SESSION_TYPE,
			Channel:           &SESSION_CHANNEL,
			ChannelData:       nil,
			Source:            "TEST",
			SourceOfTruth:     "TEST",
			AppSource:         "COMMS-API",
		}, nil
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)
		reachedSessionCreate = true
		require.Fail(t, "Interaction Session Create should not be called!")
		return nil, fmt.Errorf("interaction Session Create should not be called")
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/vcon", strings.NewReader(POST_TRANSCRIBE))
	req.Header.Add("X-Openline-VCon-Api-Key", myVconConfig.VCon.ApiKey)
	vconRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, w.Code, 200) {
		return
	}

	assert.False(t, reachedSessionCreate)
	assert.True(t, reachedSessionBySessionIdentifier)
}

func Test_vConAnalysisAsFirstEvent(t *testing.T) {
	vconRouter := gin.Default()

	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVconConfig)
	route := vconRouter.Group("/")
	reachedSessionCreate := false
	reachedSessionBySessionIdentifier := false
	var attendedBy []*model.InteractionSessionParticipantInput

	addVconRoutes(myVconConfig, route, customerOs, nil)

	resolver.AnalysisCreate = func(ctx context.Context, analysis model.AnalysisInput) (*model.Analysis, error) {
		log.Printf("InteractionEventCreate: Got Event %v", analysis)
		assert.Equal(t, "[{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Hello? Give me some real time\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Hello?\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Give me some real-\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Bye\"}]", *analysis.Content)
		assert.Equal(t, "application/x-openline-transcript", *analysis.ContentType)
		assert.Equal(t, "transcript", *analysis.AnalysisType)

		return &model.Analysis{
			ID:            "my-analysis-id",
			CreatedAt:     time.Now().UTC(),
			Content:       analysis.Content,
			ContentType:   analysis.ContentType,
			AnalysisType:  analysis.AnalysisType,
			Source:        "TEST",
			SourceOfTruth: "TEST",
			AppSource:     analysis.AppSource,
		}, nil
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)
		reachedSessionCreate = true
		attendedBy = session.AttendedBy
		require.Equal(t, "6a50a92a-7322-4988-a85c-437e1d25ea45", *session.SessionIdentifier)
		require.Equal(t, "Outgoing call to +32485112970", session.Name)
		require.Equal(t, "ACTIVE", session.Status)
		require.Equal(t, "CALL", *session.Type)
		require.Equal(t, "VOICE", *session.Channel)
		require.Equal(t, "torrey@openline.ai", *session.AttendedBy[0].Email)
		require.Equal(t, "+32485112970", *session.AttendedBy[1].PhoneNumber)

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

	resolver.InteractionSessionBySessionIdentifier = func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionBySessionIdentifier: Got Session Identifier %s", sessionIdentifier)
		require.Equal(t, "6a50a92a-7322-4988-a85c-437e1d25ea45", sessionIdentifier)
		reachedSessionBySessionIdentifier = true
		return nil, fmt.Errorf("Session not found: %s", sessionIdentifier)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/vcon", strings.NewReader(POST_TRANSCRIBE_AS_FIRST_EVENT))
	req.Header.Add("X-Openline-VCon-Api-Key", myVconConfig.VCon.ApiKey)
	vconRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, w.Code, 200) {
		return
	}

	assert.True(t, reachedSessionCreate)
	assert.True(t, reachedSessionBySessionIdentifier)
	assert.Len(t, attendedBy, 2)
}
