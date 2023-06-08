package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphContactEventHandler_OnContactCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contactEventHandler := &GraphContactEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myContactId, _ := uuid.NewUUID()
	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(tenantName, myContactId.String())
	curTime := time.Now().UTC()

	contactDto := &models.ContactDto{
		ID:          myContactId.String(),
		Tenant:      tenantName,
		FirstName:   "Bob",
		LastName:    "Smith",
		Prefix:      "Mr.",
		Description: "This is a test contact.",
		Source:      commonModels.Source{Source: "N/A", SourceOfTruth: "N/A", AppSource: "unit-test"},
		CreatedAt:   &curTime,
		UpdatedAt:   &curTime,
	}
	event, err := events.NewContactCreateEvent(contactAggregate, contactDto, curTime, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnContactCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact_"+tenantName), "Incorrect number of Contact_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "CONTACT_BELONGS_TO_TENANT"), "Incorrect number of CONTACT_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(props, "firstName"))
	require.Equal(t, "Smith", utils.GetStringPropOrEmpty(props, "lastName"))
	require.Equal(t, "Mr.", utils.GetStringPropOrEmpty(props, "prefix"))
	require.Equal(t, "This is a test contact.", utils.GetStringPropOrEmpty(props, "description"))

	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(props, "appSource"))

}
