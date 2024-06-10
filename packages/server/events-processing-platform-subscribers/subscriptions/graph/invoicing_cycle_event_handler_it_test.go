package graph

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	invoicingcycle "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestInvoicingCycleEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	eventHandler := &InvoicingCycleEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	id := uuid.New().String()
	aggregate := invoicingcycle.NewInvoicingCycleAggregateWithTenantAndID(tenantName, id)
	timeNow := utils.Now()
	createEvent, err := invoicingcycle.NewInvoicingCycleCreateEvent(
		aggregate,
		string(invoicingcycle.InvoicingCycleTypeAnniversary),
		&timeNow,
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelInvoicingCycle:                    1,
		neo4jutil.NodeLabelInvoicingCycle + "_" + tenantName: 1})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelInvoicingCycle, id)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	invoicingCycle := neo4jmapper.MapDbNodeToInvoicingCycleEntity(dbNode)
	require.Equal(t, id, invoicingCycle.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), invoicingCycle.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatformSubscribers, invoicingCycle.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), invoicingCycle.SourceOfTruth)
	require.Equal(t, timeNow, invoicingCycle.CreatedAt)
	test.AssertRecentTime(t, invoicingCycle.UpdatedAt)
	require.Equal(t, neo4jentity.InvoicingCycleTypeAnniversary, invoicingCycle.Type)
}

func TestInvoicingCycleEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	id := neo4jtest.CreateInvoicingCycle(ctx, testDatabase.Driver, tenantName, neo4jentity.InvoicingCycleEntity{})

	// Prepare the event handler
	eventHandler := &InvoicingCycleEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	timeNow := utils.Now()

	aggregate := invoicingcycle.NewInvoicingCycleAggregateWithTenantAndID(tenantName, id)
	updateEvent, err := invoicingcycle.NewInvoicingCycleUpdateEvent(
		aggregate,
		string(invoicingcycle.InvoicingCycleTypeAnniversary),
		&timeNow,
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnUpdate(context.Background(), updateEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelInvoicingCycle:                    1,
		neo4jutil.NodeLabelInvoicingCycle + "_" + tenantName: 1})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelInvoicingCycle, id)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	invoicingCycle := neo4jmapper.MapDbNodeToInvoicingCycleEntity(dbNode)
	require.Equal(t, id, invoicingCycle.Id)
	test.AssertRecentTime(t, invoicingCycle.UpdatedAt)
	require.Equal(t, neo4jentity.InvoicingCycleTypeAnniversary, invoicingCycle.Type)
}
