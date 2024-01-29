package notifications

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Boostport/mjml-go"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/pkg/errors"

	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/notifications"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
)

type MockNotificationProvider struct {
	called           bool
	emailContent     string
	notificationText string
	TemplatePath     string
	emailRawContent  string
}

func (m *MockNotificationProvider) SendNotification(ctx context.Context, u *notifications.NotifiableUser, payload map[string]interface{}, workflowId string) error {
	m.called = true
	switch workflowId {
	case notifications.WorkflowIdOrgOwnerUpdateEmail:
		m.emailContent = payload["html"].(string)
	case notifications.WorkflowIdOrgOwnerUpdateAppNotification:
		m.notificationText = payload["notificationText"].(string)
	}
	return nil
}
func (m *MockNotificationProvider) LoadEmailBody(ctx context.Context, workflowId string) error {
	switch workflowId {
	case notifications.WorkflowIdOrgOwnerUpdateEmail:
		if _, err := os.Stat(m.TemplatePath); os.IsNotExist(err) {
			return fmt.Errorf("(MockProvider.LoadEmailBody) error: %s", err.Error())
		}
		emailPath := fmt.Sprintf("%s/ownership.single.mjml", m.TemplatePath)
		if _, err := os.Stat(emailPath); err != nil {
			return fmt.Errorf("(MockProvider.LoadEmailBody) error: %s", err.Error())
		}

		rawMjml, err := os.ReadFile(emailPath)
		if err != nil {
			return fmt.Errorf("(MockProvider.LoadEmailBody) error: %s", err.Error())
		}
		m.emailRawContent = string(rawMjml[:])
	}
	return nil
}
func (m *MockNotificationProvider) Template(ctx context.Context, replace map[string]string) (string, error) {
	_, ok := replace["{{userFirstName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing userFirstName")
	}
	_, ok = replace["{{actorFirstName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing actorFirstName")
	}
	_, ok = replace["{{actorLastName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing actorLastName")
	}
	_, ok = replace["{{orgName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing orgName")
	}
	_, ok = replace["{{orgLink}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing orgLink")
	}
	mjmlf := m.emailRawContent
	for k, v := range replace {
		mjmlf = strings.Replace(mjmlf, k, v, -1)
	}
	m.emailRawContent = mjmlf
	// mjmlf := strings.Replace(string(np.emailRawContent[:]), "{{userFirstName}}", userFirstName, -1)
	// mjmlf = strings.Replace(mjmlf, "{{actorFirstName}}", actorFirstName, -1)
	// mjmlf = strings.Replace(mjmlf, "{{actorLastName}}", actorLastName, -1)
	// mjmlf = strings.Replace(mjmlf, "{{orgName}}", orgName, -1)
	// mjmlf = strings.Replace(mjmlf, "{{orgLink}}", orgLink, -1)

	html, err := mjml.ToHTML(context.Background(), mjmlf)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", mjmlError.Message)
	}
	m.emailContent = html
	return html, err
}
func (m *MockNotificationProvider) GetRawContent() string {
	return m.emailRawContent
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
	neo4jt.CreateEmailForUser(ctx, testDatabase.Driver, tenantName, actorUserId, entity.EmailEntity{
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
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
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
