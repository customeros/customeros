package notifications

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"io"
	"strings"
	"testing"

	"github.com/Boostport/mjml-go"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/aws_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/stretchr/testify/require"
)

type MockNotificationProvider struct {
	called           bool
	emailContent     string
	notificationText string
	s3               aws_client.S3ClientI
}

func (m *MockNotificationProvider) SendNotification(ctx context.Context, notification *notifications.NovuNotification, span opentracing.Span) error {
	m.called = true
	payload := notification.Payload
	workflowId := notification.WorkflowId
	rawEmailTemplate, err := m.LoadEmailBody(workflowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if rawEmailTemplate != "" {
		htmlEmailTemplate, err := m.FillTemplate(rawEmailTemplate, notification.TemplateData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		payload["html"] = htmlEmailTemplate
	}
	switch workflowId {
	case notifications.WorkflowIdOrgOwnerUpdateEmail:
		m.emailContent = payload["html"].(string)
	case notifications.WorkflowIdOrgOwnerUpdateAppNotification:
		m.notificationText = payload["notificationText"].(string)
	}
	return nil
}

func (np *MockNotificationProvider) LoadEmailBody(workflowId string) (string, error) {
	var fileName string
	switch workflowId {
	case notifications.WorkflowIdOrgOwnerUpdateEmail:
		fileName = "ownership.single.mjml"
	}

	if fileName == "" {
		return "", nil
	}

	np.s3.ChangeRegion("eu-west-1")
	return np.s3.Download("openline-production-mjml-templates", fileName)
}

func (np *MockNotificationProvider) FillTemplate(template string, replace map[string]string) (string, error) {
	mjmlf := template
	for k, v := range replace {
		mjmlf = strings.Replace(mjmlf, k, v, -1)
	}

	html, err := mjml.ToHTML(context.Background(), mjmlf)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		return "", fmt.Errorf("(NovuProvider.FillTemplate) error: %s", mjmlError.Message)
	}
	return html, err
}

type MockS3Client struct {
	aws_client.S3ClientI
}

func (m *MockS3Client) Download(bucketName string, fileName string) (string, error) {
	return EMAIL_TEMPLATE, nil
}

func (m *MockS3Client) ChangeRegion(region string) {}

func (m *MockS3Client) Upload(bucketName string, fileName string, fileContent io.Reader) error {
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
		notificationProvider: &MockNotificationProvider{
			s3: &MockS3Client{},
		},
		cfg: &config.Config{Subscriptions: config.Subscriptions{NotificationsSubscription: config.NotificationsSubscription{RedirectUrl: "https://app.openline.dev"}}},
	}

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
	require.NotEqual(t, "", orgEventHandler.notificationProvider.(*MockNotificationProvider).emailContent)
	require.True(t, emailContentHasCorrectData)
	require.True(t, emailContentIsHTML)
}

///////////////////////////////////////// email template mjml for test /////////////////////////////////////////

const EMAIL_TEMPLATE = `<mjml>
<mj-body>
  <mj-section>
	<mj-column>
	  <mj-text>
		<p>Hello {{userFirstName}},</p>
		<p>{{actorFirstName}} {{actorLastName}} made you the owner of the <a href="{{orgLink}}">{{orgName}}</a> account on CustomerOS.</p>
		<p>You’ll now see <a href="{{orgLink}}">{{orgName}}</a> in your Portfolio, and you can:
		<ul>
		  <li>Find out everything about the company and their timeline</li>
		  <li>Get familiar with the people at {{orgName}}</li>
		  <li>Manage their contract and service line items</li>
		  <li>Get their onboarding started and plan their success</li>
		</ul>
		</p>
		<p>If you have questions about this assignment of ownership, you can ask {{actorFirstName}} {{actorLastName}} about it by simply replying to this email.</p>
		<p>Have a great day!</p>
		<p>The CustomerOS Team</p>
	  </mj-text>
	</mj-column>
  </mj-section>
</mj-body>
</mjml>`
