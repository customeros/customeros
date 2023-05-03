package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"log"
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
			StartedAt     string `json:"startedAt"`
			EndedAt       string `json:"endedAt"`
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

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)

	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", "EMAIL", secAgo10)

	neo4jt.InteractionEventPartOfMeeting(ctx, driver, interactionEventId1, meetingCreate.Meeting_Create.ID)

	analysis1 := neo4jt.CreateAnalysis(ctx, driver, tenantName, "This is a summary of the conversation", "text/plain", "SUMMARY", now)
	neo4jt.ActionDescribes(ctx, driver, tenantName, analysis1, meetingCreate.Meeting_Create.ID, repository.DESCRIBES_TYPE_MEETING)

	// get meeting
	getRawResponse, err := c.RawPost(getQuery("meeting/get_meeting"), client.Var("meetingId", meetingCreate.Meeting_Create.ID))
	assertRawResponseSuccess(t, getRawResponse, err)
	var meetingGet struct {
		Meeting struct {
			ID            string `json:"id"`
			AppSource     string `json:"appSource"`
			Name          string `json:"name"`
			StartedAt     string `json:"startedAt"`
			EndedAt       string `json:"endedAt"`
			Recoding      string `json:"recording"`
			Source        string `json:"source"`
			SourceOfTruth string `json:"sourceOfTruth"`
			DescribedBy   []struct {
				ID           string `json:"id"`
				ContentType  string `json:"contentType"`
				Content      string `json:"content"`
				CreatedAt    string `json:"createdAt"`
				AnalysisType string `json:"analysisType"`
			}
			Events []struct {
				ID          string `json:"id"`
				ContentType string `json:"contentType"`
				Content     string `json:"content"`
				CreatedAt   string `json:"createdAt"`
			}
		}
	}
	err = decode.Decode(getRawResponse.Data.(map[string]interface{}), &meetingGet)
	log.Printf("meetingGet: %+v", getRawResponse.Data)
	require.Nil(t, err)
	require.NotNil(t, meetingGet.Meeting.ID)
	require.Equal(t, meetingCreate.Meeting_Create.ID, meetingGet.Meeting.ID)
	require.Equal(t, meetingCreate.Meeting_Create.Name, meetingGet.Meeting.Name)
	require.Equal(t, meetingCreate.Meeting_Create.AppSource, meetingGet.Meeting.AppSource)
	require.Equal(t, meetingCreate.Meeting_Create.StartedAt, meetingGet.Meeting.StartedAt)
	require.Equal(t, meetingCreate.Meeting_Create.EndedAt, meetingGet.Meeting.EndedAt)
	require.Equal(t, meetingCreate.Meeting_Create.Recording, meetingGet.Meeting.Recoding)
	require.Equal(t, meetingCreate.Meeting_Create.Source, meetingGet.Meeting.Source)
	require.Equal(t, meetingCreate.Meeting_Create.SourceOfTruth, meetingGet.Meeting.SourceOfTruth)
	require.Equal(t, 1, len(meetingGet.Meeting.DescribedBy))
	require.Equal(t, analysis1, meetingGet.Meeting.DescribedBy[0].ID)
	require.Equal(t, "text/plain", meetingGet.Meeting.DescribedBy[0].ContentType)
	require.Equal(t, "This is a summary of the conversation", meetingGet.Meeting.DescribedBy[0].Content)
	require.Equal(t, "SUMMARY", meetingGet.Meeting.DescribedBy[0].AnalysisType)
	require.Equal(t, 1, len(meetingGet.Meeting.Events))
	require.Equal(t, interactionEventId1, meetingGet.Meeting.Events[0].ID)
	require.Equal(t, "application/json", meetingGet.Meeting.Events[0].ContentType)
	require.Equal(t, "IE 1", meetingGet.Meeting.Events[0].Content)

	// update meeting
	rawResponse, err := c.RawPost(getQuery("meeting/update_meeting"), client.Var("meetingId", meetingCreate.Meeting_Create.ID))
	assertRawResponseSuccess(t, rawResponse, err)

	var meeting struct {
		Meeting_Update struct {
			ID                 string `json:"id"`
			AppSource          string `json:"appSource"`
			Name               string `json:"name"`
			ConferenceUrl      string `json:"conferenceUrl"`
			MeetingExternalUrl string `json:"meetingExternalUrl"`
			Agenda             string `json:"agenda"`
			AgendaContentType  string `json:"agendaContentType"`
			StartedAt          string `json:"startedAt"`
			EndedAt            string `json:"endedAt"`
			Recording          string `json:"recording"`
			Source             string `json:"source"`
			SourceOfTruth      string `json:"sourceOfTruth"`
		}
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &meeting)
	require.Nil(t, err)
	require.NotNil(t, meeting.Meeting_Update.ID)
	require.Equal(t, "test-app-source", meeting.Meeting_Update.AppSource)
	require.Equal(t, "test-name-updated", meeting.Meeting_Update.Name)
	require.Equal(t, "test-conference-url-updated", meeting.Meeting_Update.ConferenceUrl)
	require.Equal(t, "test-meeting-external-url-updated", meeting.Meeting_Update.MeetingExternalUrl)
	require.Equal(t, "2022-01-01T00:00:00Z", meeting.Meeting_Update.StartedAt)
	require.Equal(t, "2022-02-01T00:00:00Z", meeting.Meeting_Update.EndedAt)
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
	testUserId := "test_user_id"
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	testContactId1 := "test_contact_id_1"
	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId1, entity.ContactEntity{
		Prefix:        "MR",
		FirstName:     "first",
		LastName:      "last",
		Source:        entity.DataSourceHubspot,
		SourceOfTruth: entity.DataSourceHubspot,
	})

	testContactId2 := "test_contact_id_2"
	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId2, entity.ContactEntity{
		Prefix:        "MR",
		FirstName:     "first",
		LastName:      "last",
		Source:        entity.DataSourceHubspot,
		SourceOfTruth: entity.DataSourceHubspot,
	})

	// create meeting
	meeting1RawResponse, err := c.RawPost(getQuery("meeting/create_meeting_contact"),
		client.Var("createdById", testUserId),
		client.Var("attendedById", testContactId1))
	require.Nil(t, err)

	meeting2RawResponse, err := c.RawPost(getQuery("meeting/create_meeting_contact"),
		client.Var("createdById", testUserId),
		client.Var("attendedById", testContactId2))
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
		client.Var("parentContactId", testContactId1),
		client.Var("mergedContactId", testContactId2))
	require.Nil(t, err)
	assertRawResponseSuccess(t, mergeRawResponse, err)

	getRawResponse, err := c.RawPost(getQuery("contact/get_contact_with_timeline_events"),
		client.Var("contactId", testContactId1),
		client.Var("from", utils.Now()),
		client.Var("size", 2))
	require.Nil(t, err)
	assertRawResponseSuccess(t, getRawResponse, err)

	contact := getRawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, testContactId1, contact.(map[string]interface{})["id"])

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

func TestMutationResolver_AddRecordingToMeeting(t *testing.T) {
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

	rawResponse, err := c.RawPost(getQuery("meeting/add_recording_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId))
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "INCLUDES {nature: \"Recording\"}"))

	var meeting struct {
		Meeting_LinkRecording model.Meeting
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Meeting_LinkRecording.ID)
	require.NotNil(t, meeting.Meeting_LinkRecording.Recording)
	require.Len(t, meeting.Meeting_LinkRecording.Includes, 0)
}

func TestMutationResolver_AddAndRemoveContactAttendeeToMeeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	contactId1 := uuid.New().String()
	neo4jt.CreateContactWithId(ctx, driver, tenantName, contactId1, entity.ContactEntity{
		Prefix:    "MR",
		FirstName: "a",
		LastName:  "b",
	})
	addAttendeeToMeeting1, err := c.RawPost(getQuery("meeting/add_attendee_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			ContactID: &contactId1,
		}))
	assertRawResponseSuccess(t, addAttendeeToMeeting1, err)

	contactId2 := uuid.New().String()
	neo4jt.CreateContactWithId(ctx, driver, tenantName, contactId2, entity.ContactEntity{
		Prefix:    "MR",
		FirstName: "c",
		LastName:  "d",
	})
	addAttendeeToMeeting2, err := c.RawPost(getQuery("meeting/add_attendee_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			ContactID: &contactId2,
		}))
	assertRawResponseSuccess(t, addAttendeeToMeeting2, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

	rawResponseRemove, err := c.RawPost(getQuery("meeting/remove_attendee_from_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			ContactID: &contactId2,
		}))
	assertRawResponseSuccess(t, rawResponseRemove, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

	var meeting struct {
		Meeting_UnlinkAttendedBy struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			AttendedBy []map[string]interface{}
		}
	}

	err = decode.Decode(rawResponseRemove.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Meeting_UnlinkAttendedBy.ID)
	require.Len(t, meeting.Meeting_UnlinkAttendedBy.AttendedBy, 1)

	for _, attendedBy := range meeting.Meeting_UnlinkAttendedBy.AttendedBy {
		contactParticipant, _ := attendedBy["contactParticipant"].(map[string]interface{})
		require.Equal(t, contactId1, contactParticipant["id"])
	}
}

func TestMutationResolver_AddAndRemoveUserAttendeeToMeeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	userId1 := uuid.New().String()
	neo4jt.CreateUserWithId(ctx, driver, tenantName, userId1, entity.UserEntity{
		FirstName: "a",
		LastName:  "b",
	})
	rawResponse1, err := c.RawPost(getQuery("meeting/add_attendee_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			UserID: &userId1,
		}))
	assertRawResponseSuccess(t, rawResponse1, err)

	userId2 := uuid.New().String()
	neo4jt.CreateUserWithId(ctx, driver, tenantName, userId2, entity.UserEntity{
		FirstName: "c",
		LastName:  "d",
	})
	rawResponse2, err := c.RawPost(getQuery("meeting/add_attendee_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			UserID: &userId2,
		}))
	assertRawResponseSuccess(t, rawResponse2, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

	rawResponseRemove, err := c.RawPost(getQuery("meeting/remove_attendee_from_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			UserID: &userId2,
		}))
	assertRawResponseSuccess(t, rawResponseRemove, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

	var meeting struct {
		Meeting_UnlinkAttendedBy struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			AttendedBy []map[string]interface{}
		}
	}

	err = decode.Decode(rawResponseRemove.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Meeting_UnlinkAttendedBy.ID)
	require.Len(t, meeting.Meeting_UnlinkAttendedBy.AttendedBy, 1)

	for _, attendedBy := range meeting.Meeting_UnlinkAttendedBy.AttendedBy {
		userParticipant, _ := attendedBy["userParticipant"].(map[string]interface{})
		require.Equal(t, userId1, userParticipant["id"])
	}
}
