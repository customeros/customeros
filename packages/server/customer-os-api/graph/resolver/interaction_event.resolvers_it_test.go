package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestQueryResolver_Contact_WithTimelineEvents_InteractionEvents_With_InteractionSession(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "some@email.com", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo20 := now.Add(time.Duration(-20) * time.Second)
	secAgo30 := now.Add(time.Duration(-30) * time.Second)
	secAgo40 := now.Add(time.Duration(-40) * time.Second)

	// prepare interaction events
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "IE 1", "application/json", "EMAIL", secAgo10)
	interactionEventId2 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "IE 2", "application/json", "EMAIL", secAgo20)
	interactionEventId3 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "IE 3", "application/json", "EMAIL", secAgo30)
	interactionEventId4_WithoutSession := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "IE 4", "application/json", "EMAIL", secAgo40)

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4_WithoutSession, emailId, "")

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "session1", "THREAD", "ACTIVE", "EMAIL", now)
	interactionSession2 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "session2", "THREAD", "INACTIVE", "EMAIL", now)

	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId1, interactionSession1)
	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId2, interactionSession2)
	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId3, interactionSession2)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "InteractionSession"))
	require.Equal(t, 6, neo4jt.GetCountOfNodes(ctx, driver, "Action"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "PART_OF"))

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
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["startedAt"].(string))
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["endedAt"].(string))

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
	require.NotNil(t, timelineEvent2["interactionSession"].(map[string]interface{})["startedAt"].(string))
	require.NotNil(t, timelineEvent2["interactionSession"].(map[string]interface{})["endedAt"].(string))

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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	userId := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "Agent",
		LastName:  "Smith",
	})

	emailId1 := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "email_1@email.com", false, "WORK")
	emailId2 := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "email_2@email.com", false, "WORK")
	emailId3 := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "email_3@email.com", false, "WORK")
	phoneNumberId1 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1111", false, "WORK")
	phoneNumberId2 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+2222", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo20 := now.Add(time.Duration(-20) * time.Second)
	secAgo30 := now.Add(time.Duration(-30) * time.Second)

	// prepare interaction events
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "IE 1", "application/json", "EMAIL", secAgo10)
	interactionEventId2 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "IE 2", "application/json", "EMAIL", secAgo20)
	interactionEventId3 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "IE 3", "application/json", "EMAIL", secAgo30)

	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId1, phoneNumberId1, "CC")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId1, phoneNumberId2, "CC")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId1, "FROM")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId2, "FROM")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId3, "FROM")

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, userId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, contactId, "TO")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "SENT_BY"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "SENT_TO"))

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

	require.Equal(t, 1, len(timelineEvent3["sentTo"].([]interface{})))
	require.Equal(t, "TO", timelineEvent3["sentTo"].([]interface{})[0].(map[string]interface{})["type"].(string))
	require.Equal(t, contactId, timelineEvent3["sentTo"].([]interface{})[0].(map[string]interface{})["contactParticipant"].(map[string]interface{})["id"].(string))
	require.Equal(t, "first", timelineEvent3["sentTo"].([]interface{})[0].(map[string]interface{})["contactParticipant"].(map[string]interface{})["firstName"].(string))
}
