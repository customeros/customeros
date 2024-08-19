package graph

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/country"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestCountryEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	eventHandler := &CountryEventHandler{
		log:      testLogger,
		services: testDatabase.Services,
	}

	id := uuid.New().String()
	timeNow := utils.Now()

	aggregate := country.NewCountryAggregateWithID(id)
	createEvent, err := country.NewCountryCreateEvent(
		aggregate,
		"A",
		"B",
		"C",
		"D",
		timeNow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnCreate(context.Background(), createEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelCountry: 1})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model.NodeLabelCountry, id)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	countryEntity := neo4jmapper.MapDbNodeToCountryEntity(dbNode)
	require.Equal(t, id, countryEntity.Id)
	require.Equal(t, timeNow, countryEntity.CreatedAt)
	test.AssertRecentTime(t, countryEntity.UpdatedAt)
	require.Equal(t, "A", countryEntity.Name)
	require.Equal(t, "B", countryEntity.CodeA2)
	require.Equal(t, "C", countryEntity.CodeA3)
	require.Equal(t, "D", countryEntity.PhoneCode)
}
