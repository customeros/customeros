package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_TimelineEvents(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	channel := "EMAIL"
	interactionEventId1 := neo4jt.CreateInteractionEventFromEntity(ctx, driver, tenantName, entity.InteractionEventEntity{
		EventIdentifier: "myExternalId",
		Content:         "IE text 1",
		ContentType:     "application/json",
		Channel:         &channel,
		CreatedAt:       utils.TimePtr(utils.Now()),
		Hide:            false,
	})
	interactionEventId2 := neo4jt.CreateInteractionEventFromEntity(ctx, driver, tenantName, entity.InteractionEventEntity{
		EventIdentifier: "myExternalId",
		Content:         "IE text 3",
		ContentType:     "application/json",
		Channel:         &channel,
		CreatedAt:       utils.TimePtr(utils.Now()),
		Hide:            true,
	})
	neo4jt.CreateInteractionEventFromEntity(ctx, driver, tenantName,
		entity.InteractionEventEntity{
			EventIdentifier: "myExternalId",
			Content:         "IE text 2",
			ContentType:     "application/json",
			Channel:         &channel,
			CreatedAt:       utils.TimePtr(utils.Now()),
			Hide:            false,
		})
	issueId1 := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{})

	rawResponse := callGraphQL(t, "timeline/get_timeline_events_with_ids", map[string]interface{}{"ids": []string{interactionEventId1, interactionEventId2, issueId1}})

	timelineEvents := rawResponse.Data.(map[string]interface{})["timelineEvents"].([]interface{})

	require.Equal(t, 3, len(timelineEvents))

	var interactionTimelineEventIds, issueTimelineEventIds []string
	for _, timelineEvent := range timelineEvents {
		localTimelineEvent := timelineEvent.(map[string]interface{})
		if localTimelineEvent["__typename"].(string) == "InteractionEvent" {
			interactionTimelineEventIds = append(interactionTimelineEventIds, localTimelineEvent["id"].(string))
		} else {
			issueTimelineEventIds = append(issueTimelineEventIds, localTimelineEvent["id"].(string))
		}
	}

	require.ElementsMatch(t, []string{interactionEventId1, interactionEventId2}, interactionTimelineEventIds)
	require.Equal(t, issueId1, issueTimelineEventIds[0])
}
