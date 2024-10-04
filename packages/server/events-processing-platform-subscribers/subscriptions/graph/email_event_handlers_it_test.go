package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	emailAggregate "github.com/openline-ai/openline-customer-os/packages/server/events/event/email"
	event2 "github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphEmailEventHandler_OnEmailCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	emailEventHandler := &EmailEventHandler{
		services: testDatabase.Services,
	}
	myMailId, _ := uuid.NewUUID()
	agg := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, myMailId.String())
	email := "test@test.com"
	curTime := utils.Now()
	event, err := event2.NewEmailCreateEvent(agg, tenantName, email, common.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "event-processing-platform",
	}, curTime, curTime, nil, nil)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "Email"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "EMAIL_ADDRESS_BELONGS_TO_TENANT"), "Incorrect number of EMAIL_ADDRESS_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, myMailId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myMailId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, email, utils.GetStringPropOrEmpty(props, "rawEmail"))
	require.Equal(t, "event-processing-platform", utils.GetStringPropOrEmpty(props, "appSource"))
}

func TestGraphEmailEventHandler_OnEmailUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	emailCreate := "email@create.com"
	rawEmailCreate := "email@create.com"
	creationTime := utils.Now()
	emailId := neo4jtest.CreateEmail(ctx, testDatabase.Driver, tenantName, neo4jentity.EmailEntity{
		Email:     emailCreate,
		RawEmail:  rawEmailCreate,
		CreatedAt: creationTime,
		UpdatedAt: creationTime,
		IsRisky:   utils.BoolPtr(true),
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	dbNodeAfterEmailCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, emailId, utils.GetStringPropOrEmpty(propsAfterEmailCreate, "id"))

	emailEventHandler := &EmailEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}

	agg := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, emailId)
	tenant := agg.GetTenant()
	rawEmailUpdate := "email@update.com"
	sourceUpdate := constants.Anthropic
	updateTime := utils.Now()

	event, err := event2.NewEmailUpdateEvent(agg, tenant, rawEmailUpdate, sourceUpdate, updateTime)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailUpdate(context.Background(), event)
	require.Nil(t, err)
	email, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Email_"+tenantName)
	require.Nil(t, err)

	emailProps := utils.GetPropsFromNode(*email)
	require.Equal(t, 6, len(emailProps))
	emailId = utils.GetStringPropOrEmpty(emailProps, "id")
	require.NotNil(t, emailId)
	require.Nil(t, utils.GetBoolPropOrNil(emailProps, "isRisky"))
	require.Equal(t, "", utils.GetStringPropOrEmpty(emailProps, "email"))
	require.Equal(t, rawEmailUpdate, utils.GetStringPropOrEmpty(emailProps, "rawEmail"))
	require.Equal(t, creationTime, utils.GetTimePropOrNow(emailProps, "createdAt"))
	require.Less(t, creationTime, utils.GetTimePropOrNow(emailProps, "updatedAt"))
}
