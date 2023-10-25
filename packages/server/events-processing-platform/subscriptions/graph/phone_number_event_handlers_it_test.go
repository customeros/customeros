package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphPhoneNumberEventHandler_OnPhoneNumberCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	phoneNumberEventHandler := &GraphPhoneNumberEventHandler{
		Repositories: testDatabase.Repositories,
	}
	phoneNumberId, _ := uuid.NewUUID()
	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(tenantName, phoneNumberId.String())
	phoneNumber := "+0123456789"
	curTime := time.Now().UTC()
	event, err := events.NewPhoneNumberCreateEvent(phoneNumberAggregate, tenantName, phoneNumber, cmnmod.Source{
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
		AppSource:     "test",
	}, curTime, curTime)
	require.Nil(t, err)
	err = phoneNumberEventHandler.OnPhoneNumberCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName), "Incorrect number of PhoneNumber_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "PHONE_NUMBER_BELONGS_TO_TENANT"), "Incorrect number of PHONE_NUMBER_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, phoneNumberId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, phoneNumber, utils.GetStringPropOrEmpty(props, "rawPhoneNumber"))
	require.Equal(t, "test", utils.GetStringPropOrEmpty(props, "appSource"))

}
