package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphJobRoleEventHandler_OnJobRoleCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	jobRoleEventHandler := &JobRoleEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myJobRoleId, _ := uuid.NewUUID()
	curTime := time.Now().UTC()

	description := "I clean things"
	jobRoleAggregate := aggregate.NewJobRoleAggregateWithTenantAndID(tenantName, myJobRoleId.String())
	createCommand, err :=
		events.NewJobRoleCreateEvent(jobRoleAggregate, model.NewCreateJobRoleCommand(myJobRoleId.String(),
			tenantName, "Chief Janitor", &description,
			false, "N/A", "N/A", "unit-test", nil, nil, &curTime))
	require.Nil(t, err)
	err = jobRoleEventHandler.OnJobRoleCreate(context.Background(), createCommand)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole_"+tenantName), "Incorrect number of JobRole_%s nodes in Neo4j", tenantName)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "JobRole_"+tenantName, myJobRoleId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myJobRoleId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, "Chief Janitor", utils.GetStringPropOrEmpty(props, "jobTitle"))
	require.Equal(t, description, utils.GetStringPropOrEmpty(props, "description"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(props, "appSource"))
}
