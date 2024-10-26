package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestMutationResolver_Meeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
		Source:    neo4jentity.DataSourceHubspot,
	})
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test organization")
	neo4jt.CreateCalComExternalSystem(ctx, driver, tenantName)

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		RefreshLastTouchpoint: func(context context.Context, org *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: organizationId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	// create meeting
	createRawResponse, err := c.RawPost(getQuery("meeting/create_meeting"),
		client.Var("organizationId", organizationId))
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
			Note          []struct {
				ID string `json:"id"`
			}
			AttendedBy     []map[string]interface{}
			CreatedBy      []map[string]interface{}
			Recording      string                 `json:"recording"`
			ExternalSystem []model.ExternalSystem `json:"externalSystem"`
			Status         string                 `json:"status"`
		}
	}
	err = decode.Decode(createRawResponse.Data.(map[string]interface{}), &meetingCreate)
	require.Nil(t, err)
	require.NotNil(t, meetingCreate.Meeting_Create.ID)
	require.NotNil(t, meetingCreate.Meeting_Create.Note[0].ID)
	require.Equal(t, "", meetingCreate.Meeting_Create.Recording)
	require.Equal(t, "calcom", *meetingCreate.Meeting_Create.ExternalSystem[0].ExternalSource)
	require.Equal(t, model.ExternalSystemType("CALCOM"), meetingCreate.Meeting_Create.ExternalSystem[0].Type)
	require.Equal(t, "https://link-to-some-meeting.com", *meetingCreate.Meeting_Create.ExternalSystem[0].ExternalURL)
	require.Equal(t, "123", *meetingCreate.Meeting_Create.ExternalSystem[0].ExternalID)
	require.Equal(t, "ACCEPTED", meetingCreate.Meeting_Create.Status)

	for _, attendedBy := range append(meetingCreate.Meeting_Create.AttendedBy, meetingCreate.Meeting_Create.CreatedBy...) {
		if attendedBy["__typename"].(string) == "ContactParticipant" {
			contactParticipant, _ := attendedBy["contactParticipant"].(map[string]interface{})
			require.Equal(t, testContactId, contactParticipant["id"])
		} else if attendedBy["__typename"].(string) == "UserParticipant" {
			userParticipant, _ := attendedBy["userParticipant"].(map[string]interface{})
			require.Equal(t, testUserId, userParticipant["id"])
		} else if attendedBy["__typename"].(string) == "OrganizationParticipant" {
			organizationParticipant, _ := attendedBy["organizationParticipant"].(map[string]interface{})
			require.Equal(t, organizationId, organizationParticipant["id"])
		} else {
			t.Error("Unexpected participant type: " + attendedBy["__typename"].(string))
		}
	}

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)

	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "IE 1", "application/json", channel, secAgo10)
	neo4jt.InteractionEventPartOfMeeting(ctx, driver, interactionEventId1, meetingCreate.Meeting_Create.ID)

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
			Events        []struct {
				ID          string `json:"id"`
				ContentType string `json:"contentType"`
				Content     string `json:"content"`
				CreatedAt   string `json:"createdAt"`
			}
			Status string `json:"status"`
		}
	}
	err = decode.Decode(getRawResponse.Data.(map[string]interface{}), &meetingGet)
	log.Printf("meetingGet: %+v", getRawResponse.Data)
	require.Nil(t, err)
	require.NotNil(t, meetingGet.Meeting.ID)
	require.Equal(t, meetingGet.Meeting.ID, meetingGet.Meeting.ID)
	require.Equal(t, meetingGet.Meeting.Name, meetingGet.Meeting.Name)
	require.Equal(t, meetingGet.Meeting.AppSource, meetingGet.Meeting.AppSource)
	require.Equal(t, meetingGet.Meeting.StartedAt, meetingGet.Meeting.StartedAt)
	require.Equal(t, meetingGet.Meeting.EndedAt, meetingGet.Meeting.EndedAt)
	require.Equal(t, meetingGet.Meeting.Source, meetingGet.Meeting.Source)
	require.Equal(t, meetingGet.Meeting.SourceOfTruth, meetingGet.Meeting.SourceOfTruth)
	require.Equal(t, 1, len(meetingGet.Meeting.Events))
	require.Equal(t, interactionEventId1, meetingGet.Meeting.Events[0].ID)
	require.Equal(t, "application/json", meetingGet.Meeting.Events[0].ContentType)
	require.Equal(t, "IE 1", meetingGet.Meeting.Events[0].Content)
	require.Equal(t, "ACCEPTED", meetingGet.Meeting.Status)

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
			Status             string `json:"status"`
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
	require.Equal(t, "CANCELED", meeting.Meeting_Update.Status)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Note_"+tenantName))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent_"+tenantName))

	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Meeting", "Meeting_" + tenantName,
		"Note", "Note_" + tenantName,
		"Contact", "Contact_" + tenantName, "ExternalSystem", "ExternalSystem_" + tenantName, "TimelineEvent", "TimelineEvent_" + tenantName,
		"User", "User_" + tenantName, "Organization", "Organization_" + tenantName,
		"InteractionEvent", "InteractionEvent_" + tenantName})
}

func TestMutationResolver_MergeContactsWithMeetings(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	testUserId := "test_user_id"
	neo4jtest.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	testContactId1 := "test_contact_id_1"
	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId1, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
		Source:    neo4jentity.DataSourceHubspot,
	})

	testContactId2 := "test_contact_id_2"
	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId2, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
		Source:    neo4jentity.DataSourceHubspot,
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

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme.txt",
	})

	rawResponse, err := c.RawPost(getQuery("meeting/add_attachment_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId))
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES"))

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

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	attachmentId1 := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme1.txt",
	})

	attachmentId2 := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme2.txt",
	})

	addAttachment1Response, err := c.RawPost(getQuery("meeting/add_attachment_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId1))
	assertRawResponseSuccess(t, addAttachment1Response, err)

	addAttachment2Response, err := c.RawPost(getQuery("meeting/add_attachment_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId2))
	assertRawResponseSuccess(t, addAttachment2Response, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES"))

	removeAttachmentResponse, err := c.RawPost(getQuery("meeting/remove_attachment_from_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId2))
	assertRawResponseSuccess(t, removeAttachmentResponse, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES"))

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

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	attachmentId := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme.txt",
	})

	rawResponse, err := c.RawPost(getQuery("meeting/add_recording_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId))
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "RECORDING"))

	var meeting struct {
		Meeting_LinkRecording model.Meeting
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Meeting_LinkRecording.ID)
	require.NotNil(t, meeting.Meeting_LinkRecording.Recording)
	require.Len(t, meeting.Meeting_LinkRecording.Includes, 0)
}

func TestMutationResolver_RemoveRecordingFromMeeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	attachmentId1 := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme1.txt",
	})

	attachmentId2 := neo4jt.CreateAttachment(ctx, driver, tenantName, neo4jentity.AttachmentEntity{
		MimeType: "text/plain",
		FileName: "readme2.txt",
	})

	addAttachment1Response, err := c.RawPost(getQuery("meeting/add_recording_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId1))
	assertRawResponseSuccess(t, addAttachment1Response, err)
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "RECORDING"))

	addAttachment2Response, err := c.RawPost(getQuery("meeting/add_attachment_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId2))
	assertRawResponseSuccess(t, addAttachment2Response, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "RECORDING"))

	removeAttachmentResponse, err := c.RawPost(getQuery("meeting/remove_recording_from_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("attachmentId", attachmentId1))
	assertRawResponseSuccess(t, removeAttachmentResponse, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Attachment"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES"))
	require.Equal(t, 0, neo4jtest.GetCountOfRelationships(ctx, driver, "INCLUDES {nature: \"Recording\"}"))

	var meeting struct {
		Meeting_UnlinkRecording model.Meeting
	}

	err = decode.Decode(removeAttachmentResponse.Data.(map[string]any), &meeting)
	require.Nil(t, err)

	require.NotNil(t, meeting.Meeting_UnlinkRecording.ID)
	require.Len(t, meeting.Meeting_UnlinkRecording.Includes, 1)
	require.Equal(t, meeting.Meeting_UnlinkRecording.Includes[0].ID, attachmentId2)
	require.Nil(t, meeting.Meeting_UnlinkRecording.Recording)
}

func TestMutationResolver_AddAndRemoveContactAttendeeToMeeting(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	contactId1 := uuid.New().String()
	neo4jt.CreateContactWithId(ctx, driver, tenantName, contactId1, neo4jentity.ContactEntity{
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
	neo4jt.CreateContactWithId(ctx, driver, tenantName, contactId2, neo4jentity.ContactEntity{
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

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

	rawResponseRemove, err := c.RawPost(getQuery("meeting/remove_attendee_from_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			ContactID: &contactId2,
		}))
	assertRawResponseSuccess(t, rawResponseRemove, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

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

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "Meeting", time.Now().UTC())

	userId1 := uuid.New().String()
	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		Id:        userId1,
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
	neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		Id:        userId2,
		FirstName: "a",
		LastName:  "b",
	})
	rawResponse2, err := c.RawPost(getQuery("meeting/add_attendee_to_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			UserID: &userId2,
		}))
	assertRawResponseSuccess(t, rawResponse2, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

	rawResponseRemove, err := c.RawPost(getQuery("meeting/remove_attendee_from_meeting"),
		client.Var("meetingId", meetingId),
		client.Var("participant", model.MeetingParticipantInput{
			UserID: &userId2,
		}))
	assertRawResponseSuccess(t, rawResponseRemove, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

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

func TestQueryResolver_Contact_WithMultipleMeetingsInTimelineEvents(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	secondContactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	userId := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)

	neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId, "contact1@openline.ai", true, "WORK")
	neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, secondContactId, "contact2@openline.ai", true, "WORK")
	neo4jt.AddEmailTo(ctx, driver, commonModel.USER, tenantName, userId, "user@openline.ai", true, "WORK")

	now := time.Now().UTC()
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo20 := now.Add(time.Duration(-20) * time.Second)

	// prepare meeting
	meetingId1 := neo4jt.CreateMeeting(ctx, driver, tenantName, "firstMeeting", secAgo20)
	neo4jt.MeetingCreatedBy(ctx, driver, meetingId1, userId)
	neo4jt.MeetingAttendedBy(ctx, driver, meetingId1, contactId)
	neo4jt.MeetingAttendedBy(ctx, driver, meetingId1, secondContactId)
	meetingId2 := neo4jt.CreateMeeting(ctx, driver, tenantName, "secondMeeting", secAgo10)
	neo4jt.MeetingCreatedBy(ctx, driver, meetingId1, userId)
	neo4jt.MeetingAttendedBy(ctx, driver, meetingId2, contactId)
	neo4jt.MeetingAttendedBy(ctx, driver, meetingId2, secondContactId)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))

	rawResponse, err := c.RawPost(getQuery("meeting/get_multiple_meetings_in_timeline"),
		client.Var("contactId", contactId),
		client.Var("from", now),
		client.Var("size", 100))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 2, len(timelineEvents))
}

func TestMutationResolver_GetMeetings(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	testUserId := "test_user_id"
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	testContactId1 := "test_contact_id_1"
	neo4jt.CreateCalComExternalSystem(ctx, driver, tenantName)
	neo4jt.AddEmailTo(ctx, driver, commonModel.USER, tenantName, testUserId, "test-user-email", true, "MAIN")

	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId1, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
		Source:    neo4jentity.DataSourceHubspot,
	})

	testContactId2 := "test_contact_id_2"
	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId2, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
		Source:    neo4jentity.DataSourceHubspot,
	})

	// create meeting
	meeting1RawResponse, err := c.RawPost(getQuery("meeting/create_meeting_contact_external"),
		client.Var("createdById", testUserId),
		client.Var("attendedById", testContactId1))
	require.Nil(t, err)

	meeting2RawResponse, err := c.RawPost(getQuery("meeting/create_meeting_contact_external"),
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

	getMeetingResponse, err := c.RawPost(getQuery("meeting/get_meetings_basic_filters"))
	require.Nil(t, err)
	assertRawResponseSuccess(t, getMeetingResponse, err)
	var externalMeetingsResponse struct {
		ExternalMeetings struct {
			Content []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"content"`
			TotalElements int `json:"totalElements"`
			TotalPages    int `json:"totalPages"`
		} `json:"externalMeetings"`

		Errors     []interface{} `json:"errors"`
		Extensions struct {
		} `json:"extensions"`
	}

	err = decode.Decode(getMeetingResponse.Data.(map[string]interface{}), &externalMeetingsResponse)

	require.Equal(t, 2, externalMeetingsResponse.ExternalMeetings.TotalElements)
	require.Equal(t, 1, externalMeetingsResponse.ExternalMeetings.TotalPages)
	require.Equal(t, 2, len(externalMeetingsResponse.ExternalMeetings.Content))
	require.ElementsMatch(t, []string{meeting1Create.Meeting_Create.ID, meeting2Create.Meeting_Create.ID},
		[]string{externalMeetingsResponse.ExternalMeetings.Content[0].ID, externalMeetingsResponse.ExternalMeetings.Content[1].ID})
}

func TestMutationResolver_GetMeetingsWithExternalId(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	testUserId := "test_user_id"
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)
	testContactId1 := "test_contact_id_1"
	neo4jt.CreateCalComExternalSystem(ctx, driver, tenantName)
	neo4jt.AddEmailTo(ctx, driver, commonModel.USER, tenantName, testUserId, "test-user-email", true, "MAIN")

	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId1, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
		Source:    neo4jentity.DataSourceHubspot,
	})

	testContactId2 := "test_contact_id_2"
	neo4jt.CreateContactWithId(ctx, driver, tenantName, testContactId2, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "first",
		LastName:  "last",
		Source:    neo4jentity.DataSourceHubspot,
	})

	// create meeting
	meeting1RawResponse, err := c.RawPost(getQuery("meeting/create_meeting_contact_external"),
		client.Var("createdById", testUserId),
		client.Var("attendedById", testContactId1))
	require.Nil(t, err)

	assertRawResponseSuccess(t, meeting1RawResponse, err)

	var meeting1Create struct {
		Meeting_Create struct {
			ID string `json:"id"`
		}
	}

	err = decode.Decode(meeting1RawResponse.Data.(map[string]interface{}), &meeting1Create)

	require.NotNil(t, meeting1Create.Meeting_Create.ID)

	// merge contacts.$parentContactId: ID!, $mergedContactId1: ID!

	getMeetingResponse, err := c.RawPost(getQuery("meeting/get_meeting_by_external_id"))
	require.Nil(t, err)
	assertRawResponseSuccess(t, getMeetingResponse, err)
	var externalMeetingsResponse struct {
		ExternalMeetings struct {
			Content []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"content"`
			TotalElements int `json:"totalElements"`
			TotalPages    int `json:"totalPages"`
		} `json:"externalMeetings"`

		Errors     []interface{} `json:"errors"`
		Extensions struct {
		} `json:"extensions"`
	}

	err = decode.Decode(getMeetingResponse.Data.(map[string]interface{}), &externalMeetingsResponse)

	require.Equal(t, 1, externalMeetingsResponse.ExternalMeetings.TotalElements)
	require.Equal(t, 1, externalMeetingsResponse.ExternalMeetings.TotalPages)
	require.Equal(t, 1, len(externalMeetingsResponse.ExternalMeetings.Content))
	require.Equal(t, meeting1Create.Meeting_Create.ID, externalMeetingsResponse.ExternalMeetings.Content[0].ID)
}
