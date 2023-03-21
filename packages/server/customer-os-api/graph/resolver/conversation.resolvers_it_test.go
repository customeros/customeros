package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestMutationResolver_ConversationCreate_Min(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("conversation/create_conversation_min"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		Conversation_Create model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.NotNil(t, conversation.Conversation_Create.ID)
	require.NotNil(t, conversation.Conversation_Create.StartedAt)
	require.Nil(t, conversation.Conversation_Create.EndedAt)
	require.Equal(t, model.ConversationStatusActive, conversation.Conversation_Create.Status)
	require.Equal(t, "", *conversation.Conversation_Create.Channel)
	require.Equal(t, int64(0), conversation.Conversation_Create.MessageCount)
	require.Empty(t, conversation.Conversation_Create.Users)
	require.Equal(t, contactId, conversation.Conversation_Create.Contacts[0].ID)
	require.Equal(t, model.DataSourceOpenline, conversation.Conversation_Create.Source)
	require.Equal(t, "func test", *conversation.Conversation_Create.AppSource)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Conversation"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Conversation_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "PARTICIPATES"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName,
		"Conversation", "Conversation_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName})
}

func TestMutationResolver_ConversationCreate_WithGivenIdAndMultipleParticipants(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	conversationId := "Some given conversation ID"
	userId1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("conversation/create_conversation_with_multiple_participants"),
		client.Var("contactId1", contactId1),
		client.Var("contactId2", contactId2),
		client.Var("userId1", userId1),
		client.Var("userId2", userId2),
		client.Var("conversationId", conversationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		Conversation_Create model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.Equal(t, conversationId, conversation.Conversation_Create.ID)
	require.NotNil(t, conversation.Conversation_Create.StartedAt)
	require.Equal(t, "2023-01-02 03:04:05 +0000 UTC", conversation.Conversation_Create.StartedAt.String())
	require.Nil(t, conversation.Conversation_Create.EndedAt)
	require.Equal(t, model.ConversationStatusClosed, conversation.Conversation_Create.Status)
	require.Equal(t, "EMAIL", *conversation.Conversation_Create.Channel)
	require.Equal(t, int64(0), conversation.Conversation_Create.MessageCount)
	require.ElementsMatch(t, []string{contactId1, contactId2},
		[]string{conversation.Conversation_Create.Contacts[0].ID, conversation.Conversation_Create.Contacts[1].ID})
	require.ElementsMatch(t, []string{userId1, userId2},
		[]string{conversation.Conversation_Create.Users[0].ID, conversation.Conversation_Create.Users[1].ID})

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Conversation"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Conversation_"+tenantName))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "PARTICIPATES"))
	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "User", "User_" + tenantName,
		"Conversation", "Conversation_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName})
}

func TestMutationResolver_ConversationCreate_WithoutParticipants_ShouldFail(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("conversation/create_conversation_without_participants"))

	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)
	require.Contains(t, string(rawResponse.Errors), "Missing participants for new conversation")

	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Conversation"))
}

func TestMutationResolver_ConversationClose(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	conversationId := neo4jt.CreateConversation(ctx, driver, tenantName, userId, contactId, "subject", utils.Now())

	rawResponse, err := c.RawPost(getQuery("conversation/close_conversation"),
		client.Var("conversationId", conversationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		Conversation_Close model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.Equal(t, conversationId, conversation.Conversation_Close.ID)
	require.NotNil(t, conversation.Conversation_Close.StartedAt)
	require.NotNil(t, conversation.Conversation_Close.EndedAt)
	require.Equal(t, model.ConversationStatusClosed, conversation.Conversation_Close.Status)
	require.Equal(t, contactId, conversation.Conversation_Close.Contacts[0].ID)
	require.Equal(t, userId, conversation.Conversation_Close.Users[0].ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Conversation"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "PARTICIPATES"))
}

func TestMutationResolver_ConversationUpdate_NoChanges(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	conversationId := neo4jt.CreateConversation(ctx, driver, tenantName, userId, contactId, "subject", utils.Now())

	rawResponse, err := c.RawPost(getQuery("conversation/update_conversation_no_changes"),
		client.Var("conversationId", conversationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		Conversation_Update model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.Equal(t, conversationId, conversation.Conversation_Update.ID)
	require.NotNil(t, conversation.Conversation_Update.StartedAt)
	require.Nil(t, conversation.Conversation_Update.EndedAt)
	require.Equal(t, model.ConversationStatusActive, conversation.Conversation_Update.Status)
	require.Equal(t, "VOICE", *conversation.Conversation_Update.Channel)
	require.Equal(t, int64(0), conversation.Conversation_Update.MessageCount)
	require.Equal(t, contactId, conversation.Conversation_Update.Contacts[0].ID)
	require.Equal(t, userId, conversation.Conversation_Update.Users[0].ID)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Conversation"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "PARTICIPATES"))
}

func TestMutationResolver_ConversationUpdate_ChangeAllFieldsAndAddNewParticipants(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	userId1 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	conversationId := neo4jt.CreateConversation(ctx, driver, tenantName, userId1, contactId1, "subject", utils.Now())

	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "PARTICIPATES"))

	rawResponse, err := c.RawPost(getQuery("conversation/update_conversation_new_participants"),
		client.Var("conversationId", conversationId),
		client.Var("contactId1", contactId1),
		client.Var("contactId2", contactId2),
		client.Var("userId1", userId1),
		client.Var("userId2", userId2))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		Conversation_Update model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.Equal(t, conversationId, conversation.Conversation_Update.ID)
	require.NotNil(t, conversation.Conversation_Update.StartedAt)
	require.Nil(t, conversation.Conversation_Update.EndedAt)
	require.Equal(t, model.ConversationStatusClosed, conversation.Conversation_Update.Status)
	require.Equal(t, "SMS", *conversation.Conversation_Update.Channel)
	require.Equal(t, int64(1), conversation.Conversation_Update.MessageCount)
	require.ElementsMatch(t, []string{contactId1, contactId2},
		[]string{conversation.Conversation_Update.Contacts[0].ID, conversation.Conversation_Update.Contacts[1].ID})
	require.ElementsMatch(t, []string{userId1, userId2},
		[]string{conversation.Conversation_Update.Users[0].ID, conversation.Conversation_Update.Users[1].ID})

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Conversation"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "PARTICIPATES"))
}
