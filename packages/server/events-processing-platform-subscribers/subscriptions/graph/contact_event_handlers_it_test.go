package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	eventcompletionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact"
	contactevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
	"time"
)

func TestGraphContactEventHandler_OnLocationLinkToContact(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	contactName := "test_contact_name"
	contactId := neo4jtest.CreateContact(ctx, testDatabase.Driver, tenantName, neo4jentity.ContactEntity{
		Name: contactName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contact": 1, "Contact_" + tenantName: 1})
	dbNodeAfterContactCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterContactCreate)
	propsAfterContactCreate := utils.GetPropsFromNode(*dbNodeAfterContactCreate)
	require.Equal(t, contactId, utils.GetStringPropOrEmpty(propsAfterContactCreate, "id"))

	locationName := "test_location_name"
	locationId := neo4jtest.CreateLocation(ctx, testDatabase.Driver, tenantName, neo4jentity.LocationEntity{
		Name: locationName,
	})

	dbNodeAfterLocationCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterLocationCreate)
	propsAfterLocationCreate := utils.GetPropsFromNode(*dbNodeAfterLocationCreate)
	require.Equal(t, locationName, utils.GetStringPropOrEmpty(propsAfterLocationCreate, "name"))

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	contactEventHandler := &ContactEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	orgAggregate := contact.NewContactAggregateWithTenantAndID(tenantName, contactId)
	now := utils.Now()
	event, err := contactevent.NewContactLinkLocationEvent(orgAggregate, locationId, now)
	require.Nil(t, err)
	err = contactEventHandler.OnLocationLinkToContact(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "ASSOCIATED_WITH"), "Incorrect number of ASSOCIATED_WITH relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contactId, "ASSOCIATED_WITH", locationId)
}

func TestGraphContactEventHandler_OnPhoneNumberLinkToContact(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	contactName := "test_contact_name"
	now := utils.Now()
	contactId := neo4jtest.CreateContact(ctx, testDatabase.Driver, tenantName, neo4jentity.ContactEntity{
		Name:      contactName,
		UpdatedAt: now,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contact": 1, "Contact_" + tenantName: 1})
	dbNodeAfterContactCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterContactCreate)
	propsAfterContactCreate := utils.GetPropsFromNode(*dbNodeAfterContactCreate)
	require.Equal(t, contactId, utils.GetStringPropOrEmpty(propsAfterContactCreate, "id"))

	validated := false
	e164 := "+0123456789"
	phoneNumberId := neo4jtest.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, neo4jentity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterPhoneNumberCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterPhoneNumberCreate)
	propsAfterPhoneNumberCreate := utils.GetPropsFromNode(*dbNodeAfterPhoneNumberCreate)
	require.Equal(t, false, utils.GetBoolPropOrFalse(propsAfterPhoneNumberCreate, "validated"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, &creationTime, utils.GetTimePropOrNil(propsAfterPhoneNumberCreate, "updatedAt"))

	// prepare grpc mock
	callbacks := mocked_grpc.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	mocked_grpc.SetEventCompletionServiceCallbacks(&callbacks)

	contactEventHandler := &ContactEventHandler{
		services:    testDatabase.Services,
		grpcClients: testMockedGrpcClient,
	}
	contactAggregate := contact.NewContactAggregateWithTenantAndID(tenantName, contactId)
	phoneNumberLabel := "phoneNumberLabel"
	updateTime := utils.Now()
	event, err := contactevent.NewContactLinkPhoneNumberEvent(contactAggregate, phoneNumberId, phoneNumberLabel, true, updateTime)
	require.Nil(t, err)
	err = contactEventHandler.OnPhoneNumberLinkToContact(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contactId, "HAS", phoneNumberId)

	dbContactNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)
	require.Less(t, now, *utils.GetTimePropOrNil(contactProps, "updatedAt"))
}
