package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	contactAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	contactEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	contactModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphContactEventHandler_OnContactCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contactEventHandler := &ContactEventHandler{
		repositories: testDatabase.Repositories,
	}
	myContactId, _ := uuid.NewUUID()
	contactAggregate := contactAggregate.NewContactAggregateWithTenantAndID(tenantName, myContactId.String())
	curTime := time.Now().UTC()

	dataFields := contactModels.ContactDataFields{
		FirstName:   "Bob",
		LastName:    "Smith",
		Prefix:      "Mr.",
		Description: "This is a test contact.",
	}
	source :=
		cmnmod.Source{Source: "N/A", SourceOfTruth: "N/A", AppSource: "unit-test"}
	event, err := contactEvents.NewContactCreateEvent(contactAggregate, dataFields, source, cmnmod.ExternalSystem{}, curTime, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnContactCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact_"+tenantName), "Incorrect number of Contact_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "CONTACT_BELONGS_TO_TENANT"), "Incorrect number of CONTACT_BELONGS_TO_TENANT relationships in Neo4j")

	dbContactNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(contactProps, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(contactProps, "firstName"))
	require.Equal(t, "Smith", utils.GetStringPropOrEmpty(contactProps, "lastName"))
	require.Equal(t, "Mr.", utils.GetStringPropOrEmpty(contactProps, "prefix"))
	require.Equal(t, "This is a test contact.", utils.GetStringPropOrEmpty(contactProps, "description"))

	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(contactProps, "appSource"))
}
