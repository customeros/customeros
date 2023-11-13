package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceLineItemEventHandler_OnCreate(t *testing.T) {
	ctx := context.TODO()

	// Setup test environment
	serviceLineItemId := "service-line-item-id-1"
	contractId := "contract-id-1"

	// Prepare Neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateContract(ctx, testDatabase.Driver, tenantName, entity.ContractEntity{
		Id: contractId,
	})

	// Prepare the event handler
	serviceLineItemEventHandler := &ServiceLineItemEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	// Create a ServiceLineItemCreateEvent
	timeNow := utils.Now()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenantName, serviceLineItemId)
	createEvent, err := event.NewServiceLineItemCreateEvent(
		serviceLineItemAggregate,
		model.ServiceLineItemDataFields{
			Billed:      model.MonthlyBilled,
			Licenses:    10,
			Price:       100.50,
			Description: "Test service line item",
			ContractId:  contractId,
		},
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		timeNow,
		timeNow,
	)
	require.Nil(t, err, "failed to create service line item create event")

	// Execute the event handler
	err = serviceLineItemEventHandler.OnCreate(ctx, createEvent)
	require.Nil(t, err, "failed to execute service line item create event handler")

	// Assert Neo4j Node Counts
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Contract":        1,
		"ServiceLineItem": 1, "ServiceLineItem_" + tenantName: 1,
	})
	neo4jt.AssertRelationship(ctx, t, testDatabase.Driver, contractId, "HAS_SERVICE", serviceLineItemId)

	// Validate that the service line item is saved in the repository
	serviceLineItemDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "ServiceLineItem_"+tenantName, serviceLineItemId)
	require.Nil(t, err)
	require.NotNil(t, serviceLineItemDbNode)

	serviceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	require.Equal(t, serviceLineItemId, serviceLineItem.Id)
	require.Equal(t, model.MonthlyBilled.String(), serviceLineItem.Billed)
	require.Equal(t, int64(10), serviceLineItem.Licenses)
	require.Equal(t, float64(100.50), serviceLineItem.Price)
	require.Equal(t, "Test service line item", serviceLineItem.Description)
}
