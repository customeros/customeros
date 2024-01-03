package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTimelineEventRepository_CalculateAndGetLastTouchpoint_LastTouchpointIsLogEntry(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{Name: "org 1"})
	logEntryId := neo4jt.CreateLogEntryForOrg(ctx, driver, tenantName, organizationId, entity.LogEntryEntity{Content: "test content", StartedAt: utils.Now()})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization":  1,
		"LogEntry":      1,
		"TimelineEvent": 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{
		"CREATED_BY": 0,
		"LOGGED":     1,
	})

	when, timelineEventId, err := repositories.TimelineEventRepository.CalculateAndGetLastTouchpoint(ctx, tenantName, organizationId)
	if err != nil {
		t.Fatal(err)
	}

	if when == nil || timelineEventId == "" {
		t.Fatal("touchpoint should not be nil")
	}

	require.Equal(t, timelineEventId, logEntryId)
}
