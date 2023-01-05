package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_ConversationCreate_Min(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_conversation_min"),
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
	require.Equal(t, int64(0), conversation.Conversation_Create.ItemCount)

	//FIXME alexb check contacts / users

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Conversation"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Conversation_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "PARTICIPATES"))
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Conversation", "Conversation_" + tenantName})
}

func TestMutationResolver_ConversationCreate_WithGivenIdAndMultipleParticipants(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	conversationId := "Some given conversation ID"
	userId1 := neo4jt.CreateDefaultUser(driver, tenantName)
	userId2 := neo4jt.CreateDefaultUser(driver, tenantName)
	contactId1 := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_conversation_with_multiple_participants"),
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
	require.Equal(t, int64(0), conversation.Conversation_Create.ItemCount)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Conversation"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Conversation_"+tenantName))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(driver, "PARTICIPATES"))
	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "User", "Conversation", "Conversation_" + tenantName})
}

func TestMutationResolver_ConversationAddMessage(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	messageId := "A message ID"
	userId := neo4jt.CreateDefaultUser(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	conversationId := neo4jt.CreateConversation(driver, userId, contactId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Conversation"))

	rawResponse, err := c.RawPost(getQuery("add_message_to_conversation"),
		client.Var("conversationId", conversationId),
		client.Var("messageId", messageId),
		client.Var("channel", model.MessageChannelFacebook))
	assertRawResponseSuccess(t, rawResponse, err)

	var message struct {
		ConversationAddMessage model.Message
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &message)
	require.Nil(t, err)
	require.NotNil(t, message)
	createdMessage := message.ConversationAddMessage
	require.NotNil(t, createdMessage.StartedAt)
	require.Equal(t, messageId, createdMessage.ID)
	require.Equal(t, model.MessageChannelFacebook, createdMessage.Channel)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Message"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Action"))
}
