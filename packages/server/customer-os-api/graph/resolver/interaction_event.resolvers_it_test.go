package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestMutationResolver_InteractionEventCreateWithAttachment(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	now := time.Now().UTC()

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", channel, now)
	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		Id:            "",
		MimeType:      "text/plain",
		FileName:      "readme.txt",
		Source:        "",
		SourceOfTruth: "",
		AppSource:     "",
	})

	rawResponse, err := c.RawPost(getQuery("interaction_event/add_attachment_to_interaction_event"),
		client.Var("eventId", interactionEventId1),
		client.Var("attachmentId", attachmentId),
	)

	assertRawResponseSuccess(t, rawResponse, err)

	var interactionEvent struct {
		InteractionEvent_LinkAttachment model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionEvent)
	require.Nil(t, err)
	require.Equal(t, true, interactionEvent.InteractionEvent_LinkAttachment.Result)
}

func TestQueryResolver_InteractionEvent(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId, "some@email.com", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo40 := now.Add(time.Duration(-10) * time.Second)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", channel, secAgo10)
	interactionEventId4_WithoutSession := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId4", "IE 4", "application/json", channel, secAgo40)

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4_WithoutSession, emailId, "")

	interactionSession1 := neo4jtest.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, false)

	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId1, interactionSession1)
	neo4jt.InteractionEventRepliesToInteractionEvent(ctx, driver, tenantName, interactionEventId1, interactionEventId4_WithoutSession)

	neo4jt.CreateActionForInteractionEvent(ctx, driver, tenantName, interactionEventId1, neo4jenum.ActionInteractionEventRead, now)

	rawResponse, err := c.RawPost(getQuery("interaction_event/get_interaction_event"),
		client.Var("eventId", interactionEventId1))
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("response: %v", rawResponse.Data)
	responseMap, ok := rawResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("response is not a map")
	}
	interactionEventInterface := responseMap["interactionEvent"]
	timelineEvent1, ok := interactionEventInterface.(map[string]interface{})
	if !ok {
		t.Fatalf("timelineEventInterface is not a map")
	}
	require.Equal(t, interactionEventId1, timelineEvent1["id"].(string))
	require.NotNil(t, timelineEvent1["createdAt"].(string))
	require.Equal(t, "IE 1", timelineEvent1["content"].(string))
	require.Equal(t, "myExternalId1", timelineEvent1["eventIdentifier"].(string))
	require.Equal(t, "application/json", timelineEvent1["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent1["channel"].(string))
	require.NotNil(t, timelineEvent1["createdAt"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["sourceOfTruth"].(string))
	require.Equal(t, 1, len(timelineEvent1["actions"].([]interface{})))

	require.Equal(t, interactionEventId4_WithoutSession, timelineEvent1["repliesTo"].(map[string]interface{})["id"].(string))
	require.Equal(t, "IE 4", timelineEvent1["repliesTo"].(map[string]interface{})["content"].(string))
	require.Equal(t, "myExternalId4", timelineEvent1["repliesTo"].(map[string]interface{})["eventIdentifier"].(string))

	require.Equal(t, interactionSession1, timelineEvent1["interactionSession"].(map[string]interface{})["id"].(string))
	require.Equal(t, "session1", timelineEvent1["interactionSession"].(map[string]interface{})["name"].(string))
	require.Equal(t, "THREAD", timelineEvent1["interactionSession"].(map[string]interface{})["type"].(string))
	require.Equal(t, "ACTIVE", timelineEvent1["interactionSession"].(map[string]interface{})["status"].(string))
	require.Equal(t, "EMAIL", timelineEvent1["interactionSession"].(map[string]interface{})["channel"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["interactionSession"].(map[string]interface{})["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["interactionSession"].(map[string]interface{})["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent1["interactionSession"].(map[string]interface{})["appSource"].(string))
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["createdAt"].(string))
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["updatedAt"].(string))

	rawResponse, err = c.RawPost(getQuery("interaction_event/get_interaction_event"),
		client.Var("eventId", interactionEventId4_WithoutSession))
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("response: %v", rawResponse.Data)

	responseMap, ok = rawResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("response is not a map")
	}
	interactionEventInterface = responseMap["interactionEvent"]
	timelineEvent4, ok := interactionEventInterface.(map[string]interface{})
	require.Equal(t, interactionEventId4_WithoutSession, timelineEvent4["id"].(string))
	require.NotNil(t, timelineEvent4["createdAt"].(string))
	require.Equal(t, "IE 4", timelineEvent4["content"].(string))
	require.Equal(t, "application/json", timelineEvent4["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent4["channel"].(string))
	require.NotNil(t, timelineEvent4["createdAt"].(string))
	require.Equal(t, "OPENLINE", timelineEvent4["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent4["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent4["appSource"].(string))
}

func TestQueryResolver_InteractionEvent_WithIssue(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	minAgo10 := now.Add(time.Duration(-10) * time.Minute)

	// prepare interaction event
	interactionEventId := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", "", secAgo10)
	issueId := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{
		Subject:   "subject",
		CreatedAt: minAgo10,
	})

	neo4jt.InteractionEventPartOfIssue(ctx, driver, interactionEventId, issueId)

	rawResponse := callGraphQL(t, "interaction_event/get_interaction_event", map[string]interface{}{"eventId": interactionEventId})

	responseMap, ok := rawResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("response is not a map")
	}

	interactionEventInterface := responseMap["interactionEvent"]
	timelineEvent1, ok := interactionEventInterface.(map[string]interface{})
	if !ok {
		t.Fatalf("timelineEventInterface is not a map")
	}
	require.Equal(t, interactionEventId, timelineEvent1["id"].(string))

	require.Equal(t, issueId, timelineEvent1["issue"].(map[string]interface{})["id"].(string))
	require.Equal(t, "subject", timelineEvent1["issue"].(map[string]interface{})["subject"].(string))
}

func TestQueryResolver_Contact_WithTimelineEvents_InteractionEvents_With_InteractionSession(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId, "some@email.com", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo20 := now.Add(time.Duration(-20) * time.Second)
	secAgo30 := now.Add(time.Duration(-30) * time.Second)
	secAgo40 := now.Add(time.Duration(-40) * time.Second)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 1", "application/json", channel, secAgo10)
	interactionEventId2 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 2", "application/json", channel, secAgo20)
	interactionEventId3 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 3", "application/json", channel, secAgo30)
	interactionEventId4_WithoutSession := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 4", "application/json", channel, secAgo40)

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4_WithoutSession, emailId, "")

	interactionSession1 := neo4jtest.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, true)
	interactionSession2 := neo4jtest.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session2", "THREAD", "INACTIVE", "EMAIL", now, true)

	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId1, interactionSession1)
	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId2, interactionSession2)
	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId3, interactionSession2)

	neo4jt.CreateActionItemLinkedWith(ctx, driver, tenantName, string(repository.LINKED_WITH_INTERACTION_EVENT), interactionEventId1, "test action item 1", now)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "InteractionSession"))
	require.Equal(t, 6, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "ActionItem"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "PART_OF"))

	rawResponse, err := c.RawPost(getQuery("interaction_event/get_interaction_events_with_session_in_timeline_event"),
		client.Var("contactId", contactId),
		client.Var("from", now),
		client.Var("size", 100))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 4, len(timelineEvents))

	timelineEvent1 := timelineEvents[0].(map[string]interface{})
	require.Equal(t, interactionEventId1, timelineEvent1["id"].(string))
	require.NotNil(t, timelineEvent1["createdAt"].(string))

	require.Equal(t, "IE 1", timelineEvent1["content"].(string))
	require.Equal(t, "application/json", timelineEvent1["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent1["channel"].(string))
	require.NotNil(t, timelineEvent1["createdAt"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent1["appSource"].(string))
	require.Equal(t, interactionSession1, timelineEvent1["interactionSession"].(map[string]interface{})["id"].(string))
	require.Equal(t, "session1", timelineEvent1["interactionSession"].(map[string]interface{})["name"].(string))
	require.Equal(t, "THREAD", timelineEvent1["interactionSession"].(map[string]interface{})["type"].(string))
	require.Equal(t, "ACTIVE", timelineEvent1["interactionSession"].(map[string]interface{})["status"].(string))
	require.Equal(t, "EMAIL", timelineEvent1["interactionSession"].(map[string]interface{})["channel"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["interactionSession"].(map[string]interface{})["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["interactionSession"].(map[string]interface{})["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent1["interactionSession"].(map[string]interface{})["appSource"].(string))
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["createdAt"].(string))
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["updatedAt"].(string))
	require.NotNil(t, timelineEvent1["actionItems"].([]interface{}))
	require.Equal(t, "test action item 1", timelineEvent1["actionItems"].([]interface{})[0].(map[string]interface{})["content"].(string))

	timelineEvent2 := timelineEvents[1].(map[string]interface{})
	require.Equal(t, "IE 2", timelineEvent2["content"].(string))
	require.Equal(t, interactionEventId2, timelineEvent2["id"].(string))
	require.Equal(t, "application/json", timelineEvent2["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent2["channel"].(string))
	require.NotNil(t, timelineEvent2["createdAt"].(string))
	require.Equal(t, "OPENLINE", timelineEvent2["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent2["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent2["appSource"].(string))
	require.Equal(t, interactionSession2, timelineEvent2["interactionSession"].(map[string]interface{})["id"].(string))
	require.Equal(t, "session2", timelineEvent2["interactionSession"].(map[string]interface{})["name"].(string))
	require.Equal(t, "THREAD", timelineEvent2["interactionSession"].(map[string]interface{})["type"].(string))
	require.Equal(t, "INACTIVE", timelineEvent2["interactionSession"].(map[string]interface{})["status"].(string))
	require.Equal(t, "EMAIL", timelineEvent2["interactionSession"].(map[string]interface{})["channel"].(string))
	require.Equal(t, "OPENLINE", timelineEvent2["interactionSession"].(map[string]interface{})["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent2["interactionSession"].(map[string]interface{})["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent2["interactionSession"].(map[string]interface{})["appSource"].(string))
	require.NotNil(t, timelineEvent2["interactionSession"].(map[string]interface{})["createdAt"].(string))
	require.NotNil(t, timelineEvent2["interactionSession"].(map[string]interface{})["updatedAt"].(string))
	require.Nil(t, timelineEvent2["actionItems"])

	timelineEvent3 := timelineEvents[2].(map[string]interface{})
	require.Equal(t, "IE 3", timelineEvent3["content"].(string))
	require.Equal(t, interactionEventId3, timelineEvent3["id"].(string))
	require.Equal(t, "application/json", timelineEvent3["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent3["channel"].(string))
	require.NotNil(t, timelineEvent3["createdAt"].(string))
	require.Equal(t, "OPENLINE", timelineEvent3["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent3["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent3["appSource"].(string))
	require.Equal(t, interactionSession2, timelineEvent3["interactionSession"].(map[string]interface{})["id"].(string))

	timelineEvent4 := timelineEvents[3].(map[string]interface{})
	require.Equal(t, "IE 4", timelineEvent4["content"].(string))
	require.Equal(t, interactionEventId4_WithoutSession, timelineEvent4["id"].(string))
	require.Equal(t, "application/json", timelineEvent3["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent3["channel"].(string))
	require.NotNil(t, timelineEvent4["createdAt"].(string))
	require.Equal(t, "OPENLINE", timelineEvent4["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent4["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent4["appSource"].(string))
	require.Nil(t, timelineEvent4["interactionSession"])
}

func TestQueryResolver_Contact_WithTimelineEvents_InteractionEvents_With_MultipleParticipants(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "testOrg")

	userId := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		FirstName: "Agent",
		LastName:  "Smith",
	})

	emailId1 := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId, "email_1@email.com", false, "WORK")
	emailId2 := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId, "email_2@email.com", false, "WORK")
	emailId3 := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId, "email_3@email.com", false, "WORK")
	phoneNumberId1 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1111", false, "WORK")
	phoneNumberId2 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+2222", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo20 := now.Add(time.Duration(-20) * time.Second)
	secAgo30 := now.Add(time.Duration(-30) * time.Second)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 1", "application/json", channel, secAgo10)
	interactionEventId2 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 2", "application/json", channel, secAgo20)
	interactionEventId3 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 3", "application/json", channel, secAgo30)

	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId1, phoneNumberId1, "CC")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId1, phoneNumberId2, "CC")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId1, "FROM")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId2, "FROM")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId3, "FROM")

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, userId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, contactId, "TO")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, organizationId, "COLLABORATOR")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "SENT_BY"))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "SENT_TO"))

	rawResponse, err := c.RawPost(getQuery("interaction_event/get_interaction_events_with_participants_in_timeline_event"),
		client.Var("contactId", contactId),
		client.Var("from", now),
		client.Var("size", 100))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 3, len(timelineEvents))

	timelineEvent1 := timelineEvents[0].(map[string]interface{})
	require.Equal(t, interactionEventId1, timelineEvent1["id"].(string))
	require.Equal(t, 0, len(timelineEvent1["sentBy"].([]interface{})))
	require.Equal(t, 2, len(timelineEvent1["sentTo"].([]interface{})))
	require.Equal(t, "CC", timelineEvent1["sentTo"].([]interface{})[0].(map[string]interface{})["type"].(string))
	require.Equal(t, "CC", timelineEvent1["sentTo"].([]interface{})[1].(map[string]interface{})["type"].(string))
	require.ElementsMatch(t, []string{phoneNumberId1, phoneNumberId2},
		[]string{
			timelineEvent1["sentTo"].([]interface{})[0].(map[string]interface{})["phoneNumberParticipant"].(map[string]interface{})["id"].(string),
			timelineEvent1["sentTo"].([]interface{})[1].(map[string]interface{})["phoneNumberParticipant"].(map[string]interface{})["id"].(string),
		})
	require.ElementsMatch(t, []string{"+1111", "+2222"},
		[]string{
			timelineEvent1["sentTo"].([]interface{})[0].(map[string]interface{})["phoneNumberParticipant"].(map[string]interface{})["rawPhoneNumber"].(string),
			timelineEvent1["sentTo"].([]interface{})[1].(map[string]interface{})["phoneNumberParticipant"].(map[string]interface{})["rawPhoneNumber"].(string),
		})

	timelineEvent2 := timelineEvents[1].(map[string]interface{})
	require.Equal(t, interactionEventId2, timelineEvent2["id"].(string))
	require.Equal(t, 0, len(timelineEvent2["sentTo"].([]interface{})))
	require.Equal(t, 3, len(timelineEvent2["sentBy"].([]interface{})))
	require.Equal(t, "FROM", timelineEvent2["sentBy"].([]interface{})[0].(map[string]interface{})["type"].(string))
	require.Equal(t, "FROM", timelineEvent2["sentBy"].([]interface{})[1].(map[string]interface{})["type"].(string))
	require.Equal(t, "FROM", timelineEvent2["sentBy"].([]interface{})[2].(map[string]interface{})["type"].(string))
	require.ElementsMatch(t, []string{emailId1, emailId2, emailId3},
		[]string{
			timelineEvent2["sentBy"].([]interface{})[0].(map[string]interface{})["emailParticipant"].(map[string]interface{})["id"].(string),
			timelineEvent2["sentBy"].([]interface{})[1].(map[string]interface{})["emailParticipant"].(map[string]interface{})["id"].(string),
			timelineEvent2["sentBy"].([]interface{})[2].(map[string]interface{})["emailParticipant"].(map[string]interface{})["id"].(string),
		})
	require.ElementsMatch(t, []string{"email_1@email.com", "email_2@email.com", "email_3@email.com"},
		[]string{
			timelineEvent2["sentBy"].([]interface{})[0].(map[string]interface{})["emailParticipant"].(map[string]interface{})["rawEmail"].(string),
			timelineEvent2["sentBy"].([]interface{})[1].(map[string]interface{})["emailParticipant"].(map[string]interface{})["rawEmail"].(string),
			timelineEvent2["sentBy"].([]interface{})[2].(map[string]interface{})["emailParticipant"].(map[string]interface{})["rawEmail"].(string),
		})

	timelineEvent3 := timelineEvents[2].(map[string]interface{})
	require.Equal(t, interactionEventId3, timelineEvent3["id"].(string))

	require.Equal(t, 1, len(timelineEvent3["sentBy"].([]interface{})))
	require.Nil(t, timelineEvent3["sentBy"].([]interface{})[0].(map[string]interface{})["type"])
	require.Equal(t, userId, timelineEvent3["sentBy"].([]interface{})[0].(map[string]interface{})["userParticipant"].(map[string]interface{})["id"].(string))
	require.Equal(t, "Agent", timelineEvent3["sentBy"].([]interface{})[0].(map[string]interface{})["userParticipant"].(map[string]interface{})["firstName"].(string))

	var contactParticipant, organizationParticipant map[string]interface{}

	require.Equal(t, 2, len(timelineEvent3["sentTo"].([]interface{})))
	if timelineEvent3["sentTo"].([]interface{})[0].(map[string]interface{})["__typename"].(string) == "ContactParticipant" {
		contactParticipant = timelineEvent3["sentTo"].([]interface{})[0].(map[string]interface{})
		organizationParticipant = timelineEvent3["sentTo"].([]interface{})[1].(map[string]interface{})
	} else {
		contactParticipant = timelineEvent3["sentTo"].([]interface{})[1].(map[string]interface{})
		organizationParticipant = timelineEvent3["sentTo"].([]interface{})[0].(map[string]interface{})
	}
	require.Equal(t, "TO", contactParticipant["type"].(string))
	require.Equal(t, contactId, contactParticipant["contactParticipant"].(map[string]interface{})["id"].(string))
	require.Equal(t, "first", contactParticipant["contactParticipant"].(map[string]interface{})["firstName"].(string))

	require.Equal(t, "COLLABORATOR", organizationParticipant["type"].(string))
	require.Equal(t, organizationId, organizationParticipant["organizationParticipant"].(map[string]interface{})["id"].(string))
	require.Equal(t, "testOrg", organizationParticipant["organizationParticipant"].(map[string]interface{})["name"].(string))
}

func TestQueryResolver_InteractionEvent_WithExternalLinks(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	interactionEventId := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "event1", "content", "text/plain", "slack", utils.Now())

	neo4jt.CreateSlackExternalSystem(ctx, driver, tenantName)
	syncDate1 := utils.Now()
	syncDate2 := syncDate1.Add(time.Hour * 1)
	neo4jt.LinkWithSlackExternalSystem(ctx, driver, interactionEventId, "111", utils.StringPtr("www.external1.com"), nil, syncDate1)
	neo4jt.LinkWithSlackExternalSystem(ctx, driver, interactionEventId, "222", utils.StringPtr("www.external2.com"), nil, syncDate2)

	rawResponse := callGraphQL(t, "interaction_event/get_interaction_event_with_external_links",
		map[string]interface{}{"interactionEventId": interactionEventId})

	var interactionEventStruct struct {
		InteractionEvent model.InteractionEvent
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &interactionEventStruct)
	require.Nil(t, err)
	require.NotNil(t, interactionEventStruct)

	interactionEvent := interactionEventStruct.InteractionEvent
	require.Equal(t, interactionEventId, interactionEvent.ID)
	require.Equal(t, 2, len(interactionEvent.ExternalLinks))
	require.Equal(t, "111", *interactionEvent.ExternalLinks[0].ExternalID)
	require.Equal(t, "222", *interactionEvent.ExternalLinks[1].ExternalID)
	require.Equal(t, "www.external1.com", *interactionEvent.ExternalLinks[0].ExternalURL)
	require.Equal(t, "www.external2.com", *interactionEvent.ExternalLinks[1].ExternalURL)
	require.Nil(t, interactionEvent.ExternalLinks[0].ExternalSource)
	require.Nil(t, interactionEvent.ExternalLinks[1].ExternalSource)
	require.Equal(t, syncDate1, *interactionEvent.ExternalLinks[0].SyncDate)
	require.Equal(t, syncDate2, *interactionEvent.ExternalLinks[1].SyncDate)
}
