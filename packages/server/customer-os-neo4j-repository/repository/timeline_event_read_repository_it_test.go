package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTimelineEventRepository_CalculateAndGetLastTouchpoint_LastTouchpointIsLogEntry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{Name: "org 1"})
	logEntryId := neo4jtest.CreateLogEntryForOrganization(ctx, driver, tenantName, organizationId, entity.LogEntryEntity{Content: "test content", StartedAt: utils.Now()})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		model.NodeLabelOrganization:  1,
		model.NodeLabelLogEntry:      1,
		model.NodeLabelTimelineEvent: 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, driver, map[string]int{
		"CREATED_BY": 0,
		"LOGGED":     1,
	})

	when, timelineEventId, err := repositories.TimelineEventReadRepository.CalculateAndGetLastTouchPoint(ctx, tenantName, organizationId)
	if err != nil {
		t.Fatal(err)
	}

	if when == nil || timelineEventId == "" {
		t.Fatal("touchpoint should not be nil")
	}

	require.Equal(t, timelineEventId, logEntryId)
}
