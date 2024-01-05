package notifications

import (
	"context"
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"strings"
	"testing"

	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
)

type MockNotificationProvider struct {
	called           bool
	emailContent     string
	notificationText string
}

func (m *MockNotificationProvider) SendNotification(ctx context.Context, u *NotifiableUser, payload map[string]interface{}, workflowId string) error {
	m.called = true
	switch workflowId {
	case WorkflowIdOrgOwnerUpdateEmail:
		m.emailContent = payload["html"].(string)
	case WorkflowIdOrgOwnerUpdateAppNotification:
		m.notificationText = payload["notificationText"].(string)
	}
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
	neo4jt.CreateEmailForUser(ctx, testDatabase.Driver, tenantName, newOwnerUserId, entity.EmailEntity{
		Email: "owner.email@email.test",
	})

	actorUserId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName: "actor",
		LastName:  "user",
	})
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
	})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Organization": 1,
		"User":         2, "User_" + tenantName: 2,
		"Action": 0, "TimelineEvent": 0})

	// prepare event handler
	orgEventHandler := &OrganizationEventHandler{
		repositories:         testDatabase.Repositories,
		log:                  testLogger,
		notificationProvider: &MockNotificationProvider{},
		cfg: &config.Config{Services: config.Services{MJML: struct {
			ApplicationId string "env:\"MJML_APPLICATION_ID,required\" envDefault:\"\""
			SecretKey     string "env:\"MJML_SECRET_KEY,required\" envDefault:\"\""
		}{ApplicationId: "", SecretKey: ""}}, Subscriptions: config.Subscriptions{NotificationsSubscription: config.NotificationsSubscription{RedirectUrl: "https://app.openline.dev", EmailTemplatePath: "./email_templates"}}},
	}

	require.Equal(t, "", orgEventHandler.cfg.Services.MJML.ApplicationId)
	require.Equal(t, "", orgEventHandler.cfg.Services.MJML.SecretKey)

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
	organization := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, "test org", organization.Name)
	require.NotNil(t, organization.CreatedAt)
	require.NotNil(t, organization.UpdatedAt)
	require.Nil(t, organization.OnboardingDetails.SortingOrder)

	// verify we call send notification
	expectedInAppNotification := fmt.Sprintf("%s %s made you the owner of %s", "actor", "user", "test org")
	expectedSubString := fmt.Sprintf(`<p>%s %s made you the owner of the <a href="https://app.openline.dev/organization/%s">%s</a> account on CustomerOS.</p>`, "actor", "user", orgId, "test org")
	emailContentHasCorrectData := strings.Contains(orgEventHandler.notificationProvider.(*MockNotificationProvider).emailContent, expectedSubString)
	emailContentIsHTML := strings.Contains(orgEventHandler.notificationProvider.(*MockNotificationProvider).emailContent, "<!doctype html>")
	require.True(t, orgEventHandler.notificationProvider.(*MockNotificationProvider).called)
	require.Equal(t, orgEventHandler.notificationProvider.(*MockNotificationProvider).notificationText, expectedInAppNotification)
	require.True(t, emailContentHasCorrectData)
	require.True(t, emailContentIsHTML)
}
