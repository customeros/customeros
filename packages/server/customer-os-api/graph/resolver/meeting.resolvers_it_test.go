package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestMutationResolver_Meeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId, entity.ContactEntity{
		Prefix:        "MR",
		FirstName:     "first",
		LastName:      "last",
		Source:        entity.DataSourceHubspot,
		SourceOfTruth: entity.DataSourceHubspot,
	})

	// create meeting
	createRawResponse, err := c.RawPost(getQuery("meeting/create_meeting"))
	require.Nil(t, err)
	assertRawResponseSuccess(t, createRawResponse, err)
	var meetingCreate struct {
		Meeting_Create struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			CreatedAt     string `json:"createdAt"`
			UpdatedAt     string `json:"updatedAt"`
			Start         string `json:"start"`
			End           string `json:"end"`
			AppSource     string `json:"appSource"`
			Source        string `json:"source"`
			SourceOfTruth string `json:"sourceOfTruth"`
			Note          struct {
				ID string `json:"id"`
			}
			AttendedBy []map[string]interface{}
			CreatedBy  []map[string]interface{}
			Recording  string `json:"recording"`
		}
	}
	err = decode.Decode(createRawResponse.Data.(map[string]interface{}), &meetingCreate)
	require.Nil(t, err)
	require.NotNil(t, meetingCreate.Meeting_Create.ID)
	require.NotNil(t, meetingCreate.Meeting_Create.Note.ID)
	require.Equal(t, "", meetingCreate.Meeting_Create.Recording)

	for _, attendedBy := range append(meetingCreate.Meeting_Create.AttendedBy, meetingCreate.Meeting_Create.CreatedBy...) {
		if attendedBy["__typename"].(string) == "ContactParticipant" {
			contactParticipant, _ := attendedBy["contactParticipant"].(map[string]interface{})
			require.Equal(t, testContactId, contactParticipant["id"])
		} else if attendedBy["__typename"].(string) == "UserParticipant" {
			userParticipant, _ := attendedBy["userParticipant"].(map[string]interface{})
			require.Equal(t, testUserId, userParticipant["id"])
		} else {
			t.Error("Unexpected participant type: " + attendedBy["__typename"].(string))
		}
	}

	// get meeting
	getRawResponse, err := c.RawPost(getQuery("meeting/get_meeting"), client.Var("meetingId", meetingCreate.Meeting_Create.ID))
	assertRawResponseSuccess(t, getRawResponse, err)
	var meetingGet struct {
		Meeting_Get struct {
			ID            string `json:"id"`
			AppSource     string `json:"appSource"`
			Name          string `json:"name"`
			Start         string `json:"start"`
			End           string `json:"end"`
			Recoding      string `json:"recording"`
			Source        string `json:"source"`
			SourceOfTruth string `json:"sourceOfTruth"`
		}
	}
	err = decode.Decode(getRawResponse.Data.(map[string]interface{}), &meetingGet)
	require.Nil(t, err)
	require.NotNil(t, meetingGet.Meeting_Get.ID)

	// update meeting
	rawResponse, err := c.RawPost(getQuery("meeting/update_meeting"), client.Var("meetingId", meetingCreate.Meeting_Create.ID))
	assertRawResponseSuccess(t, rawResponse, err)

	var meeting struct {
		Meeting_Update struct {
			ID                string `json:"id"`
			AppSource         string `json:"appSource"`
			Name              string `json:"name"`
			ConferenceUrl     string `json:"conferenceUrl"`
			Agenda            string `json:"agenda"`
			AgendaContentType string `json:"agendaContentType"`
			Start             string `json:"start"`
			End               string `json:"end"`
			Recording         string `json:"recording"`
			Source            string `json:"source"`
			SourceOfTruth     string `json:"sourceOfTruth"`
		}
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &meeting)
	require.Nil(t, err)
	require.NotNil(t, meeting.Meeting_Update.ID)
	require.Equal(t, "test-app-source", meeting.Meeting_Update.AppSource)
	require.Equal(t, "test-name-updated", meeting.Meeting_Update.Name)
	require.Equal(t, "test-conference-url-updated", meeting.Meeting_Update.ConferenceUrl)
	require.Equal(t, "2022-01-01T00:00:00Z", meeting.Meeting_Update.Start)
	require.Equal(t, "2022-02-01T00:00:00Z", meeting.Meeting_Update.End)
	require.Equal(t, "test-agenda-updated", meeting.Meeting_Update.Agenda)
	require.Equal(t, "text/plain", meeting.Meeting_Update.AgendaContentType)
	require.Equal(t, "OPENLINE", meeting.Meeting_Update.Source)
	require.Equal(t, "OPENLINE", meeting.Meeting_Update.SourceOfTruth)
	require.Equal(t, "test-recording-id", meeting.Meeting_Update.Recording)
}

func TestMutationResolver_MergeContactsWithMeetings(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, "test_user_id")
	neo4jt.CreateContactWithId(ctx, driver, tenantName, "test_contact_id_1", entity.ContactEntity{
		Prefix:        "MR",
		FirstName:     "first",
		LastName:      "last",
		Source:        entity.DataSourceHubspot,
		SourceOfTruth: entity.DataSourceHubspot,
	})

	neo4jt.CreateContactWithId(ctx, driver, tenantName, "test_contact_id_2", entity.ContactEntity{
		Prefix:        "MR",
		FirstName:     "first",
		LastName:      "last",
		Source:        entity.DataSourceHubspot,
		SourceOfTruth: entity.DataSourceHubspot,
	})

	// create meeting
	meeting1RawResponse, err := c.RawPost(getQuery("meeting/create_meeting_contact"),
		client.Var("createdById", "test_user_id"),
		client.Var("attendedById", "test_contact_id_1"))
	require.Nil(t, err)

	meeting2RawResponse, err := c.RawPost(getQuery("meeting/create_meeting_contact"),
		client.Var("createdById", "test_user_id"),
		client.Var("attendedById", "test_contact_id_2"))
	require.Nil(t, err)

	assertRawResponseSuccess(t, meeting1RawResponse, err)
	assertRawResponseSuccess(t, meeting2RawResponse, err)

	var meeting1Create struct {
		Meeting_Create struct {
			ID string `json:"id"`
		}
	}

	var meeting2Create struct {
		Meeting_Create struct {
			ID string `json:"id"`
		}
	}

	err = decode.Decode(meeting1RawResponse.Data.(map[string]interface{}), &meeting1Create)
	err = decode.Decode(meeting2RawResponse.Data.(map[string]interface{}), &meeting2Create)

	require.NotNil(t, meeting1Create.Meeting_Create.ID)
	require.NotNil(t, meeting2Create.Meeting_Create.ID)

	// merge contacts.$parentContactId: ID!, $mergedContactId1: ID!
	mergeRawResponse, err := c.RawPost(getQuery("meeting/merge_contacts"),
		client.Var("parentContactId", "test_contact_id_1"),
		client.Var("mergedContactId", "test_contact_id_2"))
	require.Nil(t, err)
	assertRawResponseSuccess(t, mergeRawResponse, err)

	getRawResponse, err := c.RawPost(getQuery("contact/get_contact_with_timeline_events"),
		client.Var("contactId", "test_contact_id_1"),
		client.Var("from", utils.Now()),
		client.Var("size", 2))
	require.Nil(t, err)
	assertRawResponseSuccess(t, getRawResponse, err)

	contact := getRawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, "test_contact_id_1", contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 2, len(timelineEvents))
}

func TestMutationResolver_AddAttachmentToMeeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		MimeType:  "text/plain",
		Name:      "readme.txt",
		Extension: "txt",
		Size:      123,
	})

	rawResponse, err := c.RawPost(getQuery("meeting/add_attachment_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId))
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "INCLUDES"))

	var meeting struct {
		Meeting_LinkAttachment model.Meeting
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Meeting_LinkAttachment.ID)
	require.Len(t, meeting.Meeting_LinkAttachment.Includes, 1)
	require.Equal(t, meeting.Meeting_LinkAttachment.Includes[0].ID, attachmentId)
}

func TestMutationResolver_RemoveAttachmentFromMeeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	attachmentId1 := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		MimeType:  "text/plain",
		Name:      "readme1.txt",
		Extension: "txt",
		Size:      1,
	})

	attachmentId2 := neo4jt.CreateAttachment(ctx, driver, tenantName, entity.AttachmentEntity{
		MimeType:  "text/plain",
		Name:      "readme2.txt",
		Extension: "txt",
		Size:      2,
	})

	addAttachment1Response, err := c.RawPost(getQuery("meeting/add_attachment_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId1))
	assertRawResponseSuccess(t, addAttachment1Response, err)

	addAttachment2Response, err := c.RawPost(getQuery("meeting/add_attachment_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId2))
	assertRawResponseSuccess(t, addAttachment2Response, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "INCLUDES"))

	removeAttachmentResponse, err := c.RawPost(getQuery("meeting/remove_attachment_from_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId2))
	assertRawResponseSuccess(t, removeAttachmentResponse, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "INCLUDES"))

	var meeting struct {
		Meeting_UnlinkAttachment model.Meeting
	}

	err = decode.Decode(removeAttachmentResponse.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Meeting_UnlinkAttachment.ID)
	require.Len(t, meeting.Meeting_UnlinkAttachment.Includes, 1)
	require.Equal(t, meeting.Meeting_UnlinkAttachment.Includes[0].ID, attachmentId1)
}
