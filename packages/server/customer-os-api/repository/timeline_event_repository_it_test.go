package repository

import (
	"context"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimelineEventRepository_CalculateAndGetLastTouchpoint_ExcludeTranscriptionElement(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	userId := neo4jt.CreateDefaultUser(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org 1")

	neo4jt.ContactWorksForOrganization(ctx, driver, contactId, organizationId, "job role", true)

	now := utils.Now()
	oneHourAgo := now.Add(time.Duration(-3600) * time.Second)
	tenMinutesAgo := now.Add(time.Duration(-600) * time.Second)

	meetingId1 := neo4jt.CreateMeeting(ctx, driver, tenantName, "meeting name", oneHourAgo)
	neo4jt.MeetingCreatedBy(ctx, driver, meetingId1, userId)
	neo4jt.MeetingAttendedBy(ctx, driver, meetingId1, contactId)

	//insert transcription element after the meeting creation date
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId1", "transcription content", "x-openline-transcript-element", nil, tenMinutesAgo)
	neo4jt.InteractionEventPartOfMeeting(ctx, driver, interactionEventId1, meetingId1)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "CREATED_BY"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "ATTENDED_BY"))

	when, timelineEventId, err := repositories.TimelineEventRepository.CalculateAndGetLastTouchpoint(ctx, tenantName, organizationId)
	if err != nil {
		t.Fatal(err)
	}

	if when == nil || timelineEventId == "" {
		t.Fatal("touchpoint should not be nil")
	}

	require.Equal(t, timelineEventId, meetingId1)
}
