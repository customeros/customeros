package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_ConversationCreate_AutogenerateID(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	userId := neo4jt.CreateDefaultUser(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_conversation"),
		client.Var("contactId", contactId),
		client.Var("userId", userId))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		ConversationCreate model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.NotNil(t, conversation.ConversationCreate.ID)
	require.NotNil(t, conversation.ConversationCreate.StartedAt)

	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Conversation", "User"})
}

func TestMutationResolver_ConversationCreate_WithGivenID(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	conversationId := "Some conversation ID"
	userId := neo4jt.CreateDefaultUser(driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_conversation_with_id"),
		client.Var("contactId", contactId),
		client.Var("userId", userId),
		client.Var("conversationId", conversationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var conversation struct {
		ConversationCreate model.Conversation
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &conversation)
	require.Nil(t, err)
	require.NotNil(t, conversation)
	require.NotNil(t, conversation.ConversationCreate.StartedAt)
	require.Equal(t, conversationId, conversation.ConversationCreate.ID)

	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Conversation", "User"})
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
