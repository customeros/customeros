package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_ConversationCreate_AutogenerateID(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	userId := createDefaultUser(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)

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
}

func TestMutationResolver_ConversationCreate_WithGivenID(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	conversationId := "Some conversation ID"
	userId := createDefaultUser(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)

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
}

func TestMutationResolver_ConversationAddMessage(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)
	messageId := "A message ID"
	userId := createDefaultUser(driver, tenantName)
	contactId := createDefaultContact(driver, tenantName)
	conversationId := createConversation(driver, userId, contactId)

	require.Equal(t, 1, getCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, getCountOfNodes(driver, "User"))
	require.Equal(t, 1, getCountOfNodes(driver, "Conversation"))

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

	require.Equal(t, 1, getCountOfNodes(driver, "Message"))
	require.Equal(t, 1, getCountOfNodes(driver, "Action"))
}
