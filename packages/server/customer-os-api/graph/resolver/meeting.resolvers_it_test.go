package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
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
		}
	}
	err = decode.Decode(createRawResponse.Data.(map[string]interface{}), &meetingCreate)
	require.Nil(t, err)
	require.NotNil(t, meetingCreate.Meeting_Create.ID)
	require.NotNil(t, meetingCreate.Meeting_Create.Note.ID)

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
			Location          string `json:"location"`
			Agenda            string `json:"agenda"`
			AgendaContentType string `json:"agendaContentType"`
			Start             string `json:"start"`
			End               string `json:"end"`
			Source            string `json:"source"`
			SourceOfTruth     string `json:"sourceOfTruth"`
		}
	}
	err = decode.Decode(rawResponse.Data.(map[string]interface{}), &meeting)
	require.Nil(t, err)
	require.NotNil(t, meeting.Meeting_Update.ID)
	require.Equal(t, "test-app-source", meeting.Meeting_Update.AppSource)
	require.Equal(t, "test-name-updated", meeting.Meeting_Update.Name)
	require.Equal(t, "test-location-updated", meeting.Meeting_Update.Location)
	require.Equal(t, "2022-01-01T00:00:00Z", meeting.Meeting_Update.Start)
	require.Equal(t, "2022-02-01T00:00:00Z", meeting.Meeting_Update.End)
	require.Equal(t, "test-agenda-updated", meeting.Meeting_Update.Agenda)
	require.Equal(t, "text/plain", meeting.Meeting_Update.AgendaContentType)
	require.Equal(t, "OPENLINE", meeting.Meeting_Update.Source)
	require.Equal(t, "OPENLINE", meeting.Meeting_Update.SourceOfTruth)
}
