package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestQueryResolver_TimelineEvents(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 1", "application/json", "EMAIL", utils.Now())
	neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 2", "application/json", "EMAIL", utils.Now())
	issueId1 := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{})

	rawResponse := callGraphQL(t, "timeline/get_timeline_events_with_ids", map[string]interface{}{"ids": []string{interactionEventId1, issueId1}})

	timelineEvents := rawResponse.Data.(map[string]interface{})["timelineEvents"].([]interface{})

	require.Equal(t, 2, len(timelineEvents))

	var interactionTimelineEvent, issueTimelineEvent map[string]interface{}
	for _, timelineEvent := range timelineEvents {
		localTimelineEvent := timelineEvent.(map[string]interface{})
		if localTimelineEvent["__typename"].(string) == "InteractionEvent" {
			interactionTimelineEvent = localTimelineEvent
		} else {
			issueTimelineEvent = localTimelineEvent
		}
	}

	require.Equal(t, interactionEventId1, interactionTimelineEvent["id"].(string))
	require.Equal(t, issueId1, issueTimelineEvent["id"].(string))
}
