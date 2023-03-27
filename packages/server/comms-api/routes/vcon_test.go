package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var vconRouter *gin.Engine
var myVconConfig = &config.Config{
	Service: struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
		ServerAddress    string `env:"COMMS_API_SERVER_ADDRESS,required"`
		CorsUrl          string `env:"COMMS_API_CORS_URL,required"`
	}{CustomerOsAPIKey: "my-key"},
	VCon: struct {
		ApiKey          string `env:"VCON_API_KEY,required"`
		AwsAccessKey    string `env:"AWS_ACCESS_KEY"`
		AwsAccessSecret string `env:"AWS_ACCESS_SECRET"`
		AwsRegion       string `env:"AWS_REGION"`
		AwsBucket       string `env:"AWS_BUCKET"`
	}{ApiKey: "my-vcon-key"},
}

func init() {
	vconRouter = gin.Default()
}

const LIVE_TRANSCRIPTION = `{"vcon":"","uuid":"e061697f-673d-4756-a5f7-4f114e66a191","subject":"","parties":[{"mailto":"torrey@openline.ai"},{"tel":"+32485112970"}],"dialog":[{"type":"text","start":"2023-03-27T07:11:45.872099866Z","duration":0,"parties":[0,1],"mimetype":"text/plain","body":"Alors?","encoding":"None"}]}`
const LIVE_TRANSCRIPTON_APPENDED = `{"vcon":"","uuid":"54235d4f-f566-4d41-86a3-9726083c6aff","subject":"","appended":{"uuid":"e061697f-673d-4756-a5f7-4f114e66a191"},"parties":[{"mailto":"torrey@openline.ai"},{"tel":"+32485112970"}],"dialog":[{"type":"text","start":"2023-03-27T07:11:50.724873918Z","duration":0,"parties":[1,0],"mimetype":"text/plain","body":"some real time. Exit","encoding":"None"}]}`
const POST_TRANSCRIBE = `{"vcon":"","uuid":"6a50a92a-7322-4988-a85c-437e1d25ea45","subject":"","appended":{"uuid":"e061697f-673d-4756-a5f7-4f114e66a191"},"parties":[{"mailto":"torrey@openline.ai"},{"tel":"+32485112970"}],"analysis":[{"type":"transcript","dialog":null,"mimetype":"application/x-openline-transcript","body":"[{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Hello? Give me some real time\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Hello?\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Give me some real-\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Bye\"}]","encoding":""}]}`
const POST_TRANSCRIBE_AS_FIRST_EVENT = `{"vcon":"","uuid":"6a50a92a-7322-4988-a85c-437e1d25ea45","subject":"","parties":[{"mailto":"torrey@openline.ai"},{"tel":"+32485112970"}],"analysis":[{"type":"transcript","dialog":null,"mimetype":"application/x-openline-transcript","body":"[{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Hello? Give me some real time\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Hello?\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Give me some real-\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Excellent\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"tel\":\"+32485112970\"},\"text\":\"Now I'll hang up\"},{\"party\":{\"mailto\":\"torrey@openline.ai\"},\"text\":\"Bye\"}]","encoding":""}]}`

func Test_vConDialogEvent(t *testing.T) {
	_, client, resolver := utils.NewWebServer(t)
	customerOs := service.NewCustomerOSService(client, myVconConfig)
	route := vconRouter.Group("/")

	AddVconRoutes(myVconConfig, route, customerOs)

	resolver.InteractionEventCreate = func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
		log.Printf("InteractionEventCreate: Got Event %v", event)
		return &model.InteractionEvent{}, nil
	}

	resolver.InteractionSessionBySessionIdentifier = func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionBySessionIdentifier: Got Session Identifier %s", sessionIdentifier)
		return &model.InteractionSession{}, nil
	}

	resolver.InteractionSessionCreate = func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
		log.Printf("InteractionSessionCreate: Got Session %v", session)
		return &model.InteractionSession{}, nil
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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/vcon", strings.NewReader(LIVE_TRANSCRIPTION))
	req.Header.Add("X-Openline-VCon-Api-Key", myVconConfig.VCon.ApiKey)
	vconRouter.ServeHTTP(w, req)
	log.Printf("Got Body %s", w.Body)
	if !assert.Equal(t, w.Code, 200) {
		return
	}
}
