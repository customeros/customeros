package graph

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/currency"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestCurrencyEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	eventHandler := &CurrencyEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	id := uuid.New().String()
	timeNow := utils.Now()

	aggregate := currency.NewCurrencyAggregateWithID(id)
	createEvent, err := currency.NewCurrencyCreateEvent(
		aggregate,
		"USD",
		"$",
		timeNow,
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelCurrency: 1})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelCurrency, id)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	currencyEntity := neo4jmapper.MapDbNodeToCurrencyEntity(dbNode)
	require.Equal(t, id, currencyEntity.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), currencyEntity.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, currencyEntity.AppSource)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), currencyEntity.SourceOfTruth)
	require.Equal(t, timeNow, currencyEntity.CreatedAt)
	test.AssertRecentTime(t, currencyEntity.UpdatedAt)
	require.Equal(t, "USD", currencyEntity.Name)
	require.Equal(t, "$", currencyEntity.Symbol)
}
