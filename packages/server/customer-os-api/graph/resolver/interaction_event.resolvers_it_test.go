package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestMutationResolver_InteractionSessionCreate_Min(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	rawResponse, err := c.RawPost(getQuery("interaction_event/create_interaction_session_min"))
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("interactionSession: %v", rawResponse.Data)
	var interactionSession struct {
		InteractionSession_Create struct {
			ID                string `json:"id"`
			Channel           string `json:"channel"`
			AppSource         string `json:"appSource"`
			SessionIdentifier string `json:"sessionIdentifier"`
			Type              string `json:"type"`
			Name              string `json:"name"`
			Status            string `json:"status"`
		}
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionSession)
	require.Nil(t, err)
	require.Equal(t, "ACTIVE", interactionSession.InteractionSession_Create.Status)
	require.Equal(t, "CHAT", interactionSession.InteractionSession_Create.Channel)
	require.Equal(t, "Oasis", interactionSession.InteractionSession_Create.AppSource)
}

func TestMutationResolver_InteractionSessionCreateWithAttachment(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	now := time.Now().UTC()

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, false)
	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		Id:            "",
		MimeType:      "text/plain",
		Name:          "readme.txt",
		Extension:     "txt",
		Size:          123,
		Source:        "",
		SourceOfTruth: "",
		AppSource:     "",
	})

	rawResponse, err := c.RawPost(getQuery("interaction_event/add_attachment_to_interaction_session"),
		client.Var("sessionId", interactionSession1),
		client.Var("attachmentId", attachmentId),
	)

	assertRawResponseSuccess(t, rawResponse, err)

	var interactionSession struct {
		InteractionSession_LinkAttachment model.InteractionSession
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionSession)
	require.Nil(t, err)
	require.Equal(t, "ACTIVE", interactionSession.InteractionSession_LinkAttachment.Status)
	require.Equal(t, "EMAIL", *interactionSession.InteractionSession_LinkAttachment.Channel)
	require.Equal(t, "test", interactionSession.InteractionSession_LinkAttachment.AppSource)
	require.Len(t, interactionSession.InteractionSession_LinkAttachment.Includes, 1)
	require.Equal(t, attachmentId, interactionSession.InteractionSession_LinkAttachment.Includes[0].ID)
	require.Equal(t, "text/plain", interactionSession.InteractionSession_LinkAttachment.Includes[0].MimeType)
}

func TestMutationResolver_InteractionSessionCreateWithPhone(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	userId := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "Agent",
		LastName:  "Smith",
	})

	emailId1 := neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId, "user1@openline.ai", true, "WORK")

	rawResponse, err := c.RawPost(getQuery("interaction_event/create_interaction_session_with_phone"),
		client.Var("sessionIdentifier", "My Session Identifier"),
		client.Var("name", "My Session Name"),
		client.Var("type", "THREAD"),
		client.Var("channel", "EMAIL"),
		client.Var("channelData", "{\"threading-info\":\"test\"}"),
		client.Var("status", "ACTIVE"),
	)

	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("interactionSession: %v", rawResponse.Data)
	var interactionSession struct {
		InteractionSession_Create struct {
			ID                string `json:"id"`
			Channel           string `json:"channel"`
			ChannelData       string `json:"channelData"`
			AppSource         string `json:"appSource"`
			SessionIdentifier string `json:"sessionIdentifier"`
			Type              string `json:"type"`
			Name              string `json:"name"`
			Status            string `json:"status"`
			AttendedBy        []map[string]interface{}
		}
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionSession)
	require.Nil(t, err)
	require.Equal(t, "ACTIVE", interactionSession.InteractionSession_Create.Status)
	require.Equal(t, "EMAIL", interactionSession.InteractionSession_Create.Channel)
	require.Equal(t, "{\"threading-info\":\"test\"}", interactionSession.InteractionSession_Create.ChannelData)
	require.Equal(t, "Oasis", interactionSession.InteractionSession_Create.AppSource)
	require.Equal(t, "My Session Identifier", interactionSession.InteractionSession_Create.SessionIdentifier)
	require.Equal(t, "My Session Name", interactionSession.InteractionSession_Create.Name)

	for _, attendedBy := range interactionSession.InteractionSession_Create.AttendedBy {
		if attendedBy["__typename"].(string) == "EmailParticipant" {
			emailParticipant, _ := attendedBy["emailParticipant"].(map[string]interface{})
			require.Equal(t, emailId1, emailParticipant["id"])
			require.Equal(t, "user1@openline.ai", emailParticipant["rawEmail"])
		} else if attendedBy["__typename"].(string) == "PhoneNumberParticipant" {
			phoneNumberParticipant, _ := attendedBy["phoneNumberParticipant"].(map[string]interface{})
			require.Equal(t, "+1234567890", phoneNumberParticipant["rawPhoneNumber"])
		} else {
			t.Error("Unexpected participant type: " + attendedBy["__typename"].(string))
		}
	}
}

func TestMutationResolver_InteractionSessionCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	rawResponse, err := c.RawPost(getQuery("interaction_event/create_interaction_session"),
		client.Var("sessionIdentifier", "My Session Identifier"),
		client.Var("name", "My Session Name"),
		client.Var("type", "THREAD"),
		client.Var("channel", "EMAIL"),
		client.Var("channelData", "{\"threading-info\":\"test\"}"),
		client.Var("status", "ACTIVE"),
	)
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("interactionSession: %v", rawResponse.Data)
	var interactionSession struct {
		InteractionSession_Create struct {
			ID                string `json:"id"`
			Channel           string `json:"channel"`
			ChannelData       string `json:"channelData"`
			AppSource         string `json:"appSource"`
			SessionIdentifier string `json:"sessionIdentifier"`
			Type              string `json:"type"`
			Name              string `json:"name"`
			Status            string `json:"status"`
		}
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionSession)
	require.Nil(t, err)
	require.Equal(t, "ACTIVE", interactionSession.InteractionSession_Create.Status)
	require.Equal(t, "EMAIL", interactionSession.InteractionSession_Create.Channel)
	require.Equal(t, "{\"threading-info\":\"test\"}", interactionSession.InteractionSession_Create.ChannelData)
	require.Equal(t, "Oasis", interactionSession.InteractionSession_Create.AppSource)
	require.Equal(t, "My Session Identifier", interactionSession.InteractionSession_Create.SessionIdentifier)
	require.Equal(t, "My Session Name", interactionSession.InteractionSession_Create.Name)
}

func TestMutationResolver_InteractionEventCreateWithAttachment(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := time.Now().UTC()

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", &channel, now)
	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		Id:            "",
		MimeType:      "text/plain",
		Name:          "readme.txt",
		Extension:     "txt",
		Size:          123,
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
		InteractionEvent_LinkAttachment model.InteractionEvent
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionEvent)
	require.Nil(t, err)
	require.Equal(t, "application/json", *interactionEvent.InteractionEvent_LinkAttachment.ContentType)
	require.Equal(t, "EMAIL", *interactionEvent.InteractionEvent_LinkAttachment.Channel)
	require.Equal(t, "test", interactionEvent.InteractionEvent_LinkAttachment.AppSource)
	require.Len(t, interactionEvent.InteractionEvent_LinkAttachment.Includes, 1)
	require.Equal(t, attachmentId, interactionEvent.InteractionEvent_LinkAttachment.Includes[0].ID)
	require.Equal(t, "text/plain", interactionEvent.InteractionEvent_LinkAttachment.Includes[0].MimeType)

}

func TestMutationResolver_InteractionEventCreate_Min(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	rawResponse, err := c.RawPost(getQuery("interaction_event/create_interaction_event_min"))
	assertRawResponseSuccess(t, rawResponse, err)

	var interactionEvent struct {
		InteractionEvent_Create struct {
			ID        string  `json:"id"`
			Channel   *string `json:"channel"`
			AppSource string  `json:"appSource"`
			SentTo    []struct {
				Typename         string `json:"__typename"`
				EmailParticipant struct {
					ID       string `json:"id"`
					RawEmail string `json:"rawEmail"`
				} `json:"emailParticipant"`
				Type string `json:"type"`
			} `json:"sentTo"`
			SentBy []struct {
				Typename         string `json:"__typename"`
				EmailParticipant struct {
					ID       string `json:"id"`
					RawEmail string `json:"rawEmail"`
				} `json:"emailParticipant"`
				Type string `json:"type"`
			} `json:"sentBy"`
			InteractionSession struct {
				ID                string `json:"id"`
				AppSource         string `json:"appSource"`
				Channel           string `json:"channel"`
				Name              string `json:"name"`
				SessionIdentifier string `json:"sessionIdentifier"`
			} `json:"interactionSession"`
			RepliesTo struct {
				ID string `json:"id"`
			} `json:"repliesTo"`
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionEvent)
	log.Printf("interactionEvent: %v", rawResponse.Data)

	require.Nil(t, err)
	require.NotNil(t, interactionEvent)
	require.Equal(t, *interactionEvent.InteractionEvent_Create.Channel, "CHAT")
	require.Equal(t, interactionEvent.InteractionEvent_Create.AppSource, "Oasis")
	require.Equal(t, len(interactionEvent.InteractionEvent_Create.SentBy), 1)
	require.Equal(t, interactionEvent.InteractionEvent_Create.SentBy[0].EmailParticipant.RawEmail, "email_1@openline.ai")

}

func TestMutationResolver_InteractionEventCreate_Email(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	now := time.Now().UTC()

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, false)

	rawResponse, err := c.RawPost(getQuery("interaction_event/create_interaction_event_email"),
		client.Var("content", "Content 1"),
		client.Var("contentType", "text/plain"),
		client.Var("sessionId", interactionSession1))
	assertRawResponseSuccess(t, rawResponse, err)

	type interactionEventType struct {
		InteractionEvent_Create struct {
			ID          string  `json:"id"`
			Channel     *string `json:"channel"`
			ChannelData *string `json:"channelData"`
			AppSource   string  `json:"appSource"`
			Content     string  `json:"content"`
			ContentType string  `json:"contentType"`
			SentTo      []struct {
				Typename         string `json:"__typename"`
				EmailParticipant struct {
					ID       string `json:"id"`
					RawEmail string `json:"rawEmail"`
				} `json:"emailParticipant"`
				Type string `json:"type"`
			} `json:"sentTo"`
			SentBy []struct {
				Typename         string `json:"__typename"`
				EmailParticipant struct {
					ID       string `json:"id"`
					RawEmail string `json:"rawEmail"`
				} `json:"emailParticipant"`
				Type string `json:"type"`
			} `json:"sentBy"`
			InteractionSession struct {
				ID                string `json:"id"`
				AppSource         string `json:"appSource"`
				Channel           string `json:"channel"`
				Name              string `json:"name"`
				SessionIdentifier string `json:"sessionIdentifier"`
			} `json:"interactionSession"`
			RepliesTo struct {
				ID string `json:"id"`
			} `json:"repliesTo"`
		}
	}

	var firstEvent interactionEventType
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &firstEvent)

	require.Nil(t, err)
	require.NotNil(t, firstEvent)
	require.Equal(t, *firstEvent.InteractionEvent_Create.Channel, "EMAIL")
	require.Equal(t, *firstEvent.InteractionEvent_Create.ChannelData, "{\"References\":[\"<CAJYQ2j8Q>\"],\"Replies-To\":\"<CAJYQ2j8Q>\"}")
	require.Equal(t, firstEvent.InteractionEvent_Create.AppSource, "Oasis")
	require.Equal(t, firstEvent.InteractionEvent_Create.Content, "Content 1")
	require.Equal(t, firstEvent.InteractionEvent_Create.ContentType, "text/plain")
	require.Equal(t, len(firstEvent.InteractionEvent_Create.SentBy), 1)

	require.Equal(t, firstEvent.InteractionEvent_Create.InteractionSession.ID, interactionSession1)
	require.Equal(t, firstEvent.InteractionEvent_Create.InteractionSession.Name, "session1")

	sentById := firstEvent.InteractionEvent_Create.SentBy[0].EmailParticipant.ID
	require.Equal(t, firstEvent.InteractionEvent_Create.SentBy[0].EmailParticipant.RawEmail, "sentBy@openline.ai")

	//send to 2 people
	require.Equal(t, len(firstEvent.InteractionEvent_Create.SentTo), 2)

	dest1 := ""
	dest2 := ""
	for _, sendTo := range firstEvent.InteractionEvent_Create.SentTo {
		if sendTo.EmailParticipant.RawEmail == "dest1@openline.ai" {
			dest1 = sendTo.EmailParticipant.ID
			require.Equal(t, sendTo.Type, "TO")
		}
		if sendTo.EmailParticipant.RawEmail == "dest2@openline.ai" {
			dest2 = sendTo.EmailParticipant.ID
			require.Equal(t, sendTo.Type, "CC")
		}
	}

	require.NotEmpty(t, dest1)
	require.NotEmpty(t, dest2)

	origMsgId := firstEvent.InteractionEvent_Create.ID

	var secondEvent interactionEventType

	rawResponse, err = c.RawPost(getQuery("interaction_event/create_interaction_event_email"),
		client.Var("content", "Content 2"),
		client.Var("contentType", "text/plain"),
		client.Var("sessionId", interactionSession1),
		client.Var("replyTo", origMsgId))
	assertRawResponseSuccess(t, rawResponse, err)
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &secondEvent)

	// check the email addresses are re-used
	require.Nil(t, err)
	require.Equal(t, secondEvent.InteractionEvent_Create.RepliesTo.ID, origMsgId)
	require.Equal(t, secondEvent.InteractionEvent_Create.SentBy[0].EmailParticipant.ID, sentById)

	for _, sendTo := range secondEvent.InteractionEvent_Create.SentTo {
		if sendTo.EmailParticipant.RawEmail == "dest1@openline.ai" {
			require.Equal(t, dest1, sendTo.EmailParticipant.ID)
		}
		if sendTo.EmailParticipant.RawEmail == "dest2@openline.ai" {
			require.Equal(t, dest2, sendTo.EmailParticipant.ID)
		}
	}
}

func TestMutationResolver_InteractionEventCreate_Meeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	now := time.Now().UTC()

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "meeting-name", now)

	rawResponse, err := c.RawPost(getQuery("interaction_event/create_interaction_event_meeting"),
		client.Var("content", "Content 1"),
		client.Var("contentType", "text/plain"),
		client.Var("meetingId", meetingId),
		client.Var("eventType", "meeting"))
	assertRawResponseSuccess(t, rawResponse, err)

	type interactionEventType struct {
		InteractionEvent_Create struct {
			ID          string `json:"id"`
			Channel     string `json:"channel"`
			AppSource   string `json:"appSource"`
			Content     string `json:"content"`
			ContentType string `json:"contentType"`
			EventType   string `json:"eventType"`
			SentTo      []struct {
				Typename         string `json:"__typename"`
				EmailParticipant struct {
					ID       string `json:"id"`
					RawEmail string `json:"rawEmail"`
				} `json:"emailParticipant"`
				PhoneNumberParticipant struct {
					ID             string `json:"id"`
					RawPhoneNumber string `json:"rawPhoneNumber"`
				} `json:"phoneNumberParticipant"`
			} `json:"sentTo"`
			SentBy []struct {
				Typename         string `json:"__typename"`
				EmailParticipant struct {
					ID       string `json:"id"`
					RawEmail string `json:"rawEmail"`
				} `json:"emailParticipant"`
				PhoneNumberParticipant struct {
					ID             string `json:"id"`
					RawPhoneNumber string `json:"rawPhoneNumber"`
				} `json:"phoneNumberParticipant"`
			} `json:"sentBy"`
			InteractionSession struct {
				ID                string `json:"id"`
				AppSource         string `json:"appSource"`
				Channel           string `json:"channel"`
				Name              string `json:"name"`
				SessionIdentifier string `json:"sessionIdentifier"`
			} `json:"interactionSession"`
			Meeting struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"meeting"`
			RepliesTo struct {
				ID string `json:"id"`
			} `json:"repliesTo"`
		}
	}

	var interactionEvent interactionEventType
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionEvent)
	log.Printf("interactionEvent: %v", rawResponse.Data)

	require.Nil(t, err)
	require.NotNil(t, interactionEvent)
	require.Equal(t, interactionEvent.InteractionEvent_Create.AppSource, "Oasis")
	require.Equal(t, interactionEvent.InteractionEvent_Create.EventType, "meeting")
	require.Equal(t, interactionEvent.InteractionEvent_Create.Content, "Content 1")
	require.Equal(t, interactionEvent.InteractionEvent_Create.ContentType, "text/plain")
	require.Equal(t, len(interactionEvent.InteractionEvent_Create.SentBy), 0)

	require.Equal(t, len(interactionEvent.InteractionEvent_Create.SentTo), 0)

	require.Equal(t, interactionEvent.InteractionEvent_Create.Meeting.ID, meetingId)
	require.Equal(t, interactionEvent.InteractionEvent_Create.Meeting.Name, "meeting-name")

}

func TestMutationResolver_InteractionEventCreate_Voice(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	userId := neo4jt.CreateUser(ctx, driver, tenantName, entity.UserEntity{
		FirstName: "Agent",
		LastName:  "Smith",
	})

	emailId1 := neo4jt.AddEmailTo(ctx, driver, entity.USER, tenantName, userId, "user1@openline.ai", true, "WORK")

	now := time.Now().UTC()

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "VOICE", now, false)

	rawResponse, err := c.RawPost(getQuery("interaction_event/create_interaction_event_call"),
		client.Var("content", "Content 1"),
		client.Var("contentType", "text/plain"),
		client.Var("sessionId", interactionSession1))
	assertRawResponseSuccess(t, rawResponse, err)

	type interactionEventType struct {
		InteractionEvent_Create struct {
			ID          string `json:"id"`
			Channel     string `json:"channel"`
			AppSource   string `json:"appSource"`
			Content     string `json:"content"`
			ContentType string `json:"contentType"`
			SentTo      []struct {
				Typename         string `json:"__typename"`
				EmailParticipant struct {
					ID       string `json:"id"`
					RawEmail string `json:"rawEmail"`
				} `json:"emailParticipant"`
				PhoneNumberParticipant struct {
					ID             string `json:"id"`
					RawPhoneNumber string `json:"rawPhoneNumber"`
				} `json:"phoneNumberParticipant"`
			} `json:"sentTo"`
			SentBy []struct {
				Typename         string `json:"__typename"`
				EmailParticipant struct {
					ID       string `json:"id"`
					RawEmail string `json:"rawEmail"`
				} `json:"emailParticipant"`
				PhoneNumberParticipant struct {
					ID             string `json:"id"`
					RawPhoneNumber string `json:"rawPhoneNumber"`
				} `json:"phoneNumberParticipant"`
			} `json:"sentBy"`
			InteractionSession struct {
				ID                string `json:"id"`
				AppSource         string `json:"appSource"`
				Channel           string `json:"channel"`
				Name              string `json:"name"`
				SessionIdentifier string `json:"sessionIdentifier"`
			} `json:"interactionSession"`
			RepliesTo struct {
				ID string `json:"id"`
			} `json:"repliesTo"`
		}
	}

	var interactionEvent interactionEventType
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionEvent)
	log.Printf("interactionEvent: %v", rawResponse.Data)

	require.Nil(t, err)
	require.NotNil(t, interactionEvent)
	require.Equal(t, interactionEvent.InteractionEvent_Create.Channel, "VOICE")
	require.Equal(t, interactionEvent.InteractionEvent_Create.AppSource, "Oasis")
	require.Equal(t, interactionEvent.InteractionEvent_Create.Content, "Content 1")
	require.Equal(t, interactionEvent.InteractionEvent_Create.ContentType, "text/plain")
	require.Equal(t, len(interactionEvent.InteractionEvent_Create.SentBy), 1)
	require.Equal(t, interactionEvent.InteractionEvent_Create.SentBy[0].PhoneNumberParticipant.RawPhoneNumber, "+1234567890")

	phoneNumberId := interactionEvent.InteractionEvent_Create.SentBy[0].PhoneNumberParticipant.ID

	require.Equal(t, len(interactionEvent.InteractionEvent_Create.SentTo), 1)
	require.Equal(t, interactionEvent.InteractionEvent_Create.SentTo[0].EmailParticipant.RawEmail, "user1@openline.ai")
	require.Equal(t, interactionEvent.InteractionEvent_Create.SentTo[0].EmailParticipant.ID, emailId1)

	require.Equal(t, interactionEvent.InteractionEvent_Create.InteractionSession.ID, interactionSession1)
	require.Equal(t, interactionEvent.InteractionEvent_Create.InteractionSession.Name, "session1")

	rawResponse, err = c.RawPost(getQuery("interaction_event/create_interaction_event_call2"),
		client.Var("content", "Content 2"),
		client.Var("contentType", "text/plain"),
		client.Var("sessionId", interactionSession1))
	assertRawResponseSuccess(t, rawResponse, err)

	var interactionEvent2 interactionEventType
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &interactionEvent2)
	log.Printf("interactionEvent: %v", rawResponse.Data)

	require.Nil(t, err)
	require.Equal(t, interactionEvent2.InteractionEvent_Create.Channel, "VOICE")
	require.Equal(t, interactionEvent2.InteractionEvent_Create.AppSource, "Oasis")
	require.Equal(t, interactionEvent2.InteractionEvent_Create.Content, "Content 2")
	require.Equal(t, interactionEvent2.InteractionEvent_Create.ContentType, "text/plain")
	require.Equal(t, len(interactionEvent2.InteractionEvent_Create.SentBy), 1)
	require.Equal(t, interactionEvent2.InteractionEvent_Create.SentBy[0].EmailParticipant.RawEmail, "user1@openline.ai")
	require.Equal(t, interactionEvent2.InteractionEvent_Create.SentBy[0].EmailParticipant.ID, emailId1)

	require.Equal(t, len(interactionEvent2.InteractionEvent_Create.SentTo), 1)
	require.Equal(t, interactionEvent2.InteractionEvent_Create.SentTo[0].PhoneNumberParticipant.RawPhoneNumber, "+1234567890")
	require.Equal(t, interactionEvent2.InteractionEvent_Create.SentTo[0].PhoneNumberParticipant.ID, phoneNumberId)
}

func TestQueryResolver_InteractionEvent(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "some@email.com", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo40 := now.Add(time.Duration(-10) * time.Second)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", &channel, secAgo10)
	interactionEventId4_WithoutSession := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId4", "IE 4", "application/json", &channel, secAgo40)

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4_WithoutSession, emailId, "")

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, false)

	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId1, interactionSession1)
	neo4jt.InteractionEventRepliesToInteractionEvent(ctx, driver, tenantName, interactionEventId1, interactionEventId4_WithoutSession)

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
	require.Equal(t, "test", timelineEvent1["appSource"].(string))

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
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["startedAt"].(string))
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["endedAt"].(string))

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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	minAgo10 := now.Add(time.Duration(-10) * time.Minute)

	// prepare interaction event
	interactionEventId := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", nil, secAgo10)
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

func TestQueryResolver_InteractionEvent_ByEventIdentifier(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "some@email.com", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo40 := now.Add(time.Duration(-10) * time.Second)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", &channel, secAgo10)

	interactionEventId4_WithoutSession := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId4", "IE 4", "application/json", &channel, secAgo40)

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4_WithoutSession, emailId, "")

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, false)

	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId1, interactionSession1)
	neo4jt.InteractionEventRepliesToInteractionEvent(ctx, driver, tenantName, interactionEventId1, interactionEventId4_WithoutSession)
	rawResponse, err := c.RawPost(getQuery("interaction_event/get_interaction_event_by_event_identifier"),
		client.Var("eventId", "myExternalId1"))
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("response: %v", rawResponse.Data)
	responseMap, ok := rawResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("response is not a map")
	}
	interactionEventInterface := responseMap["interactionEvent_ByEventIdentifier"]
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
	require.Equal(t, "test", timelineEvent1["appSource"].(string))

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
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["startedAt"].(string))
	require.NotNil(t, timelineEvent1["interactionSession"].(map[string]interface{})["endedAt"].(string))

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

func TestQueryResolver_InteractionSession(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "some@email.com", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo40 := now.Add(time.Duration(-10) * time.Second)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", &channel, secAgo10)

	interactionEventId4_WithoutSession := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId4", "IE 4", "application/json", &channel, secAgo40)

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4_WithoutSession, emailId, "")

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, false)

	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId1, interactionSession1)
	neo4jt.InteractionEventRepliesToInteractionEvent(ctx, driver, tenantName, interactionEventId1, interactionEventId4_WithoutSession)

	analysis1 := neo4jt.CreateAnalysis(ctx, driver, tenantName, "This is a summary of the conversation", "text/plain", "SUMMARY", now)
	neo4jt.ActionDescribes(ctx, driver, tenantName, analysis1, interactionSession1, entity.DESCRIBES_TYPE_INTERACTION_SESSION)

	rawResponse, err := c.RawPost(getQuery("interaction_event/get_interaction_session"),
		client.Var("sessionId", interactionSession1))
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("response: %v", rawResponse.Data)
	responseMap, ok := rawResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("response is not a map")
	}
	interactionEventInterface := responseMap["interactionSession"]
	timelineEvent1, ok := interactionEventInterface.(map[string]interface{})
	if !ok {
		t.Fatalf("timelineEventInterface is not a map")
	}

	require.Equal(t, interactionSession1, timelineEvent1["id"].(string))
	require.Equal(t, "session1", timelineEvent1["name"].(string))
	require.Equal(t, "THREAD", timelineEvent1["type"].(string))
	require.Equal(t, "ACTIVE", timelineEvent1["status"].(string))
	require.Equal(t, "EMAIL", timelineEvent1["channel"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent1["appSource"].(string))
	require.NotNil(t, timelineEvent1["startedAt"].(string))
	require.NotNil(t, timelineEvent1["endedAt"].(string))

	events := timelineEvent1["events"].([]interface{})
	require.NotEmpty(t, events)
	event := events[0].(map[string]interface{})
	require.Equal(t, interactionEventId1, event["id"].(string))
	require.NotNil(t, event["createdAt"].(string))
	require.Equal(t, "IE 1", event["content"].(string))
	require.Equal(t, "myExternalId1", event["eventIdentifier"].(string))
	require.Equal(t, "application/json", event["contentType"].(string))
	require.Equal(t, "EMAIL", event["channel"].(string))
	require.Equal(t, "OPENLINE", event["source"].(string))
	require.Equal(t, "OPENLINE", event["sourceOfTruth"].(string))
	require.Equal(t, "test", event["appSource"].(string))

	analyses := timelineEvent1["describedBy"].([]interface{})
	require.NotEmpty(t, analyses)
	analysis := analyses[0].(map[string]interface{})

	require.Equal(t, analysis1, analysis["id"].(string))
	require.Equal(t, "This is a summary of the conversation", analysis["content"].(string))
	require.Equal(t, "text/plain", analysis["contentType"].(string))
	require.Equal(t, "SUMMARY", analysis["analysisType"].(string))

	require.Equal(t, interactionEventId4_WithoutSession, event["repliesTo"].(map[string]interface{})["id"].(string))
	require.Equal(t, "IE 4", event["repliesTo"].(map[string]interface{})["content"].(string))
	require.Equal(t, "myExternalId4", event["repliesTo"].(map[string]interface{})["eventIdentifier"].(string))

}

func TestQueryResolver_InteractionSession_BySessionIdentifier(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	emailId := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId, "some@email.com", false, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo40 := now.Add(time.Duration(-10) * time.Second)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", &channel, secAgo10)

	interactionEventId4_WithoutSession := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId4", "IE 4", "application/json", &channel, secAgo40)

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4_WithoutSession, emailId, "")

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, false)

	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId1, interactionSession1)
	neo4jt.InteractionEventRepliesToInteractionEvent(ctx, driver, tenantName, interactionEventId1, interactionEventId4_WithoutSession)
	rawResponse, err := c.RawPost(getQuery("interaction_event/get_interaction_session_by_identifier"),
		client.Var("sessionId", "mySessionIdentifier"))
	assertRawResponseSuccess(t, rawResponse, err)
	log.Printf("response: %v", rawResponse.Data)
	responseMap, ok := rawResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("response is not a map")
	}
	interactionEventInterface := responseMap["interactionSession_BySessionIdentifier"]
	timelineEvent1, ok := interactionEventInterface.(map[string]interface{})
	if !ok {
		t.Fatalf("timelineEventInterface is not a map")
	}

	require.Equal(t, interactionSession1, timelineEvent1["id"].(string))
	require.Equal(t, "session1", timelineEvent1["name"].(string))
	require.Equal(t, "THREAD", timelineEvent1["type"].(string))
	require.Equal(t, "ACTIVE", timelineEvent1["status"].(string))
	require.Equal(t, "EMAIL", timelineEvent1["channel"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent1["sourceOfTruth"].(string))
	require.Equal(t, "test", timelineEvent1["appSource"].(string))
	require.Equal(t, "mySessionIdentifier", timelineEvent1["sessionIdentifier"].(string))
	require.NotNil(t, timelineEvent1["startedAt"].(string))
	require.NotNil(t, timelineEvent1["endedAt"].(string))

	events := timelineEvent1["events"].([]interface{})
	require.NotEmpty(t, events)
	event := events[0].(map[string]interface{})
	require.Equal(t, interactionEventId1, event["id"].(string))
	require.NotNil(t, event["createdAt"].(string))
	require.Equal(t, "IE 1", event["content"].(string))
	require.Equal(t, "myExternalId1", event["eventIdentifier"].(string))
	require.Equal(t, "application/json", event["contentType"].(string))
	require.Equal(t, "EMAIL", event["channel"].(string))
	require.Equal(t, "OPENLINE", event["source"].(string))
	require.Equal(t, "OPENLINE", event["sourceOfTruth"].(string))
	require.Equal(t, "test", event["appSource"].(string))

	require.Equal(t, interactionEventId4_WithoutSession, event["repliesTo"].(map[string]interface{})["id"].(string))
	require.Equal(t, "IE 4", event["repliesTo"].(map[string]interface{})["content"].(string))
	require.Equal(t, "myExternalId4", event["repliesTo"].(map[string]interface{})["eventIdentifier"].(string))

}

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
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 1", "application/json", &channel, secAgo10)
	interactionEventId2 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 2", "application/json", &channel, secAgo20)
	interactionEventId3 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 3", "application/json", &channel, secAgo30)
	interactionEventId4_WithoutSession := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 4", "application/json", &channel, secAgo40)

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4_WithoutSession, emailId, "")

	interactionSession1 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "THREAD", "ACTIVE", "EMAIL", now, true)
	interactionSession2 := neo4jt.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session2", "THREAD", "INACTIVE", "EMAIL", now, true)

	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId1, interactionSession1)
	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId2, interactionSession2)
	neo4jt.InteractionEventPartOfInteractionSession(ctx, driver, interactionEventId3, interactionSession2)

	neo4jt.CreateActionItemLinkedWith(ctx, driver, tenantName, string(repository.LINKED_WITH_INTERACTION_EVENT), interactionEventId1, "test action item 1", now)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "InteractionSession"))
	require.Equal(t, 6, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))
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
	require.NotNil(t, timelineEvent1["actionItems"].([]interface{}))
	require.Equal(t, "test action item 1", timelineEvent1["actionItems"].([]interface{})[0].(map[string]interface{})["content"].(string))
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
	require.NotNil(t, timelineEvent2["actionItems"].([]interface{}))
	require.Equal(t, 0, len(timelineEvent2["actionItems"].([]interface{})))
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
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "testOrg")

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
	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 1", "application/json", &channel, secAgo10)
	interactionEventId2 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 2", "application/json", &channel, secAgo20)
	interactionEventId3 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE 3", "application/json", &channel, secAgo30)

	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId1, phoneNumberId1, "CC")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId1, phoneNumberId2, "CC")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId1, "FROM")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId2, "FROM")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId3, "FROM")

	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, userId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, contactId, "TO")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, organizationId, "COLLABORATOR")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "SENT_BY"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "SENT_TO"))

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
