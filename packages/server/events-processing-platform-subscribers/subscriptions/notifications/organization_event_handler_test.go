package notifications

import (
	"context"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"testing"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/aws_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/stretchr/testify/require"
)

type MockNotificationProvider struct {
	called           bool
	emailContent     string
	notificationText string
	s3               aws_client.S3ClientI
}

func (m *MockNotificationProvider) SendNotification(ctx context.Context, notification *commonService.NovuNotification) error {
	m.called = true
	return nil
}

func TestGraphOrganizationEventHandler_OnOrganizationUpdateOwner(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	newOwnerUserId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "owner",
		LastName:  "user",
	})
	neo4jtest.CreateEmailForEntity(ctx, testDatabase.Driver, tenantName, newOwnerUserId, neo4jentity.EmailEntity{
		Email: "owner.email@email.test",
	})

	actorUserId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "actor",
		LastName:  "user",
	})
	neo4jtest.CreateEmailForEntity(ctx, testDatabase.Driver, tenantName, actorUserId, neo4jentity.EmailEntity{
		Email: "actor.email@email.test",
	})
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "test org",
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1,
		"User":         2, "User_" + tenantName: 2,
		"Action": 0, "TimelineEvent": 0})

	// prepare event handler
	orgEventHandler := &OrganizationEventHandler{
		services: testDatabase.Services,
		log:      testLogger,
		cfg:      &config.Config{Subscriptions: config.Subscriptions{NotificationsSubscription: config.NotificationsSubscription{RedirectUrl: "https://app.openline.dev"}}},
	}

	orgEventHandler.services.CommonServices.NovuService = &MockNotificationProvider{}

	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationOwnerUpdateEvent(orgAggregate, newOwnerUserId, actorUserId, orgId, now)
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnOrganizationUpdateOwner(context.Background(), event)
	require.Nil(t, err)

	// verify no new nodes created nor changed, our handler just sends notification
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User": 2, "User_" + tenantName: 2,
		"Organization": 1, "Organization_" + tenantName: 1,
		"Action": 0, "Action_" + tenantName: 0,
		"TimelineEvent": 0, "TimelineEvent_" + tenantName: 0})

	orgDbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, orgDbNode)

	// verify organization
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, "test org", organization.Name)
	require.NotNil(t, organization.CreatedAt)
	require.NotNil(t, organization.UpdatedAt)
	require.Nil(t, organization.OnboardingDetails.SortingOrder)
}
