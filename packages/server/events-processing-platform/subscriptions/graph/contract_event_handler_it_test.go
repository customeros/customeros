package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContractEventHandler_OnCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{})
	userIdCreator := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{})
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "User": 1, "ExternalSystem": 1, "Contract": 0})

	// Prepare the event handler
	contractEventHandler := &ContractEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create a ContractCreateEvent
	contractId := uuid.New().String()
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenantName, contractId)
	timeNow := utils.Now()
	createEvent, err := event.NewContractCreateEvent(
		contractAggregate,
		model.ContractDataFields{
			Name:             "New Contract",
			OrganizationId:   orgId,
			CreatedByUserId:  userIdCreator,
			ServiceStartedAt: &timeNow,
			SignedAt:         &timeNow,
			RenewalCycle:     model.MonthlyRenewal,
			Status:           model.Live,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		commonmodel.ExternalSystem{
			ExternalSystemId: "sf",
			ExternalId:       "ext-id-1",
		},
		timeNow,
		timeNow,
	)
	require.Nil(t, err, "failed to create contract create event")

	// Execute
	err = contractEventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err, "failed to execute contract create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1,
		"User":         1,
		"Contract":     1, "Contract_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "CREATED_BY", userIdCreator)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, orgId, "HAS_CONTRACT", contractId)
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "IS_LINKED_WITH", "sf")

	// Assert relationships and contract data in Neo4j
	contractDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contract_"+tenantName, contractId)
	require.Nil(t, err)
	require.NotNil(t, contractDbNode)

	// Verify contract
	contract := graph_db.MapDbNodeToContractEntity(*contractDbNode)
	require.Equal(t, contractId, contract.Id)
	require.Equal(t, "New Contract", contract.Name)
	require.Equal(t, model.Live.String(), contract.Status)
	require.Equal(t, model.MonthlyRenewal.String(), contract.RenewalCycle)
	require.True(t, timeNow.Equal(contract.CreatedAt.UTC()))
	require.True(t, timeNow.Equal(contract.UpdatedAt.UTC()))
	require.True(t, timeNow.Equal(*contract.ServiceStartedAt))
	require.True(t, timeNow.Equal(*contract.SignedAt))
	require.Nil(t, contract.EndedAt)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), contract.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, contract.AppSource)
}
