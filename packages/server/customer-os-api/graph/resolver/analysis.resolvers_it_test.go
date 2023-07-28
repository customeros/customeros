package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestMutationResolver_AnalysisCreate_Session(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	now := utils.Now()
	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "CALL", "INACTIVE", "VOICE", now, false)

	rawResponse, err := c.RawPost(getQuery("analysis/create_analysis"),
		client.Var("contentType", "text/plain"),
		client.Var("content", "This is a summary of the conversation"),
		client.Var("analysisType", "SUMMARY"),
		client.Var("sessionId", interactionSession1),
	)

	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("interactionSession: %v", rawResponse.Data)
	var analysis struct {
		Analysis_Create struct {
			ID           string              `json:"id"`
			ContentType  string              `json:"contentType"`
			Content      string              `json:"content"`
			AnalysisType string              `json:"analysisType"`
			AppSource    string              `json:"appSource"`
			Describes    []map[string]string `json:"describes"`
		} `json:"analysis_Create"`
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &analysis)
	require.Nil(t, err)
	require.Equal(t, "text/plain", analysis.Analysis_Create.ContentType)
	require.Equal(t, "This is a summary of the conversation", analysis.Analysis_Create.Content)
	require.Equal(t, "Oasis", analysis.Analysis_Create.AppSource)

	require.Len(t, analysis.Analysis_Create.Describes, 1)
	log.Printf("Describe: %v", analysis.Analysis_Create.Describes[0])
	require.Equal(t, "InteractionSession", analysis.Analysis_Create.Describes[0]["__typename"])
	require.Equal(t, interactionSession1, analysis.Analysis_Create.Describes[0]["id"])
	require.Equal(t, "mySessionIdentifier", analysis.Analysis_Create.Describes[0]["sessionIdentifier"])
	require.Equal(t, "session1", analysis.Analysis_Create.Describes[0]["name"])

}

func TestMutationResolver_AnalysisCreate_Meeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	now := utils.Now()
	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "meeting-name", now)

	rawResponse, err := c.RawPost(getQuery("analysis/create_analysis"),
		client.Var("contentType", "text/plain"),
		client.Var("content", "This is a summary of the conversation"),
		client.Var("analysisType", "SUMMARY"),
		client.Var("meetingId", meetingId),
	)

	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("analysisCreate: %v", rawResponse.Data)
	var analysis struct {
		Analysis_Create struct {
			ID           string              `json:"id"`
			ContentType  string              `json:"contentType"`
			Content      string              `json:"content"`
			AnalysisType string              `json:"analysisType"`
			AppSource    string              `json:"appSource"`
			Describes    []map[string]string `json:"describes"`
		} `json:"analysis_Create"`
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &analysis)
	require.Nil(t, err)
	require.Equal(t, "text/plain", analysis.Analysis_Create.ContentType)
	require.Equal(t, "This is a summary of the conversation", analysis.Analysis_Create.Content)
	require.Equal(t, "Oasis", analysis.Analysis_Create.AppSource)

	require.Len(t, analysis.Analysis_Create.Describes, 1)
	log.Printf("Describe: %v", analysis.Analysis_Create.Describes[0])
	require.Equal(t, "Meeting", analysis.Analysis_Create.Describes[0]["__typename"])
	require.Equal(t, meetingId, analysis.Analysis_Create.Describes[0]["id"])
	require.Equal(t, "meeting-name", analysis.Analysis_Create.Describes[0]["meetingName"])

}

func TestMutationResolver_AnalysisCreate_Event(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	now := utils.Now()
	channel := "VOICE"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "Hello?", "text/plain", &channel, now)

	rawResponse, err := c.RawPost(getQuery("analysis/create_analysis"),
		client.Var("contentType", "application/x-openline-translation"),
		client.Var("content", "{\"lang\": \"fr\", \"text\": \"Bonjour?\"}"),
		client.Var("analysisType", "TRANSLATION"),
		client.Var("eventId", interactionEventId1),
	)

	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("interactionSession: %v", rawResponse.Data)
	var analysis struct {
		Analysis_Create struct {
			ID           string              `json:"id"`
			ContentType  string              `json:"contentType"`
			Content      string              `json:"content"`
			AnalysisType string              `json:"analysisType"`
			AppSource    string              `json:"appSource"`
			Describes    []map[string]string `json:"describes"`
		} `json:"analysis_Create"`
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &analysis)
	require.Nil(t, err)
	require.Equal(t, "application/x-openline-translation", analysis.Analysis_Create.ContentType)
	require.Equal(t, "{\"lang\": \"fr\", \"text\": \"Bonjour?\"}", analysis.Analysis_Create.Content)
	require.Equal(t, "Oasis", analysis.Analysis_Create.AppSource)

	require.Len(t, analysis.Analysis_Create.Describes, 1)
	log.Printf("Describe: %v", analysis.Analysis_Create.Describes[0])
	require.Equal(t, "InteractionEvent", analysis.Analysis_Create.Describes[0]["__typename"])
	require.Equal(t, interactionEventId1, analysis.Analysis_Create.Describes[0]["id"])
	require.Equal(t, "myExternalId1", analysis.Analysis_Create.Describes[0]["eventIdentifier"])
	require.Equal(t, "Hello?", analysis.Analysis_Create.Describes[0]["content"])
	require.Equal(t, "text/plain", analysis.Analysis_Create.Describes[0]["contentType"])

}

func TestQueryResolver_Analysis(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := time.Now().UTC()

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "CALL", "ACTIVE", "VOICE", now, false)

	analysis1 := neo4jt.CreateAnalysis(ctx, driver, tenantName, "This is a summary of the conversation", "text/plain", "SUMMARY", now)
	neo4jt.AnalysisDescribes(ctx, driver, tenantName, analysis1, interactionSession1, string(repository.LINKED_WITH_INTERACTION_SESSION))

	rawResponse, err := c.RawPost(getQuery("analysis/get_analysis"),
		client.Var("analysisId", analysis1))
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("response: %v", rawResponse.Data)

	var analysis struct {
		Analysis struct {
			ID           string              `json:"id"`
			ContentType  string              `json:"contentType"`
			Content      string              `json:"content"`
			AnalysisType string              `json:"analysisType"`
			AppSource    string              `json:"appSource"`
			Describes    []map[string]string `json:"describes"`
		} `json:"analysis"`
	}
	_, ok := rawResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("response is not a map")
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &analysis)
	require.Nil(t, err)
	require.Equal(t, "text/plain", analysis.Analysis.ContentType)
	require.Equal(t, "This is a summary of the conversation", analysis.Analysis.Content)
	require.Equal(t, "test", analysis.Analysis.AppSource)

	require.Len(t, analysis.Analysis.Describes, 1)
	log.Printf("Describe: %v", analysis.Analysis.Describes[0])
	require.Equal(t, "InteractionSession", analysis.Analysis.Describes[0]["__typename"])
	require.Equal(t, interactionSession1, analysis.Analysis.Describes[0]["id"])
	require.Equal(t, "mySessionIdentifier", analysis.Analysis.Describes[0]["sessionIdentifier"])
	require.Equal(t, "session1", analysis.Analysis.Describes[0]["name"])

}
