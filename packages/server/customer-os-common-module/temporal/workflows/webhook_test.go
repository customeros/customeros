package workflows

import (
	"context"
	"errors"
	"testing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/activity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_WebhookWorkflow_ActivityFails() {
	s.env.OnActivity(activity.WebhookActivity,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		func(ctx context.Context, targetUrl, authHeaderName, authHeaderValue string, reqBody string, notifyFailure bool, notification *notifications.NovuNotification, provider notifications.NotificationProvider) error {
			return errors.New("WebhookActivityFailure")
		})
	workflowParams := WHWorkflowParam{
		TargetUrl:            "test_URL",
		RequestBody:          `{"data":{"amountDue":28.16,"amountDueInSmallestUnit":2816,"amountPaid":0,"amountRemaining":28.16,"currency":"GBP","due":"2024-02-29T19:59:12.032628504Z","invoiceNumber":"GVD-88528","invoicePeriodEnd":"2024-01-31T00:00:00Z","invoicePeriodStart":"2024-01-01T00:00:00Z","invoiceUrl":"https://fs.customeros.ai/file//download","note":"invoiceNote","paid":false,"status":"DUE","subtotal":28.16,"taxDue":0,"contract":{"contractName":"BCC Contract Name 002","contractStatus":"LIVE","metadata":{"id":"172c636c-bb9c-4abe-9a93-ff4cb6dc466a"}},"invoiceLineItems":[{"description":"Service line 002 for BCC Contract Name 001","metadata":{"id":"94163956-2af9-402f-b4f1-df73c25a6454"}},{"description":"Service line 003 for BCC Contract Name 001","metadata":{"id":"b0f36efc-c8c3-4a7c-981f-2523aeceacec"}},{"description":"Service line 004 for BCC Contract Name 001","metadata":{"id":"dd72fce7-b724-4ba4-aab9-afd5d8650c77"}}],"metadata":{"created":"2024-02-29T19:59:12.032628504Z","id":"980f37be-04c7-42f6-9fd4-852456fd450e"},"organization":{"customerOsId":"C-PPK-7RN","metadata":{"id":"efb6a64a-3222-4b46-acaa-c09cc0749752"},"name":"BCC"}},"event":"invoice.finalized"}`,
		AuthHeaderName:       "test_auth_header",
		AuthHeaderValue:      "test_auth_value",
		RetryPolicy:          nil,
		Notification:         nil,
		NotificationProvider: nil,
		NotifyFailure:        false,
	}
	s.env.ExecuteWorkflow(WebhookWorkflow, workflowParams)

	s.True(s.env.IsWorkflowCompleted())

	err := s.env.GetWorkflowError()
	s.Error(err)
	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("WebhookActivityFailure", applicationErr.Error())
}

func (s *UnitTestSuite) Test_WebhookWorkflow_ActivityParamCorrect() {
	payload := `{"data":{"amountDue":28.16,"amountDueInSmallestUnit":2816,"amountPaid":0,"amountRemaining":28.16,"currency":"GBP","due":"2024-02-29T19:59:12.032628504Z","invoiceNumber":"GVD-88528","invoicePeriodEnd":"2024-01-31T00:00:00Z","invoicePeriodStart":"2024-01-01T00:00:00Z","invoiceUrl":"https://fs.customeros.ai/file//download","note":"invoiceNote","paid":false,"status":"DUE","subtotal":28.16,"taxDue":0,"contract":{"contractName":"BCC Contract Name 002","contractStatus":"LIVE","metadata":{"id":"172c636c-bb9c-4abe-9a93-ff4cb6dc466a"}},"invoiceLineItems":[{"description":"Service line 002 for BCC Contract Name 001","metadata":{"id":"94163956-2af9-402f-b4f1-df73c25a6454"}},{"description":"Service line 003 for BCC Contract Name 001","metadata":{"id":"b0f36efc-c8c3-4a7c-981f-2523aeceacec"}},{"description":"Service line 004 for BCC Contract Name 001","metadata":{"id":"dd72fce7-b724-4ba4-aab9-afd5d8650c77"}}],"metadata":{"created":"2024-02-29T19:59:12.032628504Z","id":"980f37be-04c7-42f6-9fd4-852456fd450e"},"organization":{"customerOsId":"C-PPK-7RN","metadata":{"id":"efb6a64a-3222-4b46-acaa-c09cc0749752"},"name":"BCC"}},"event":"invoice.finalized"}`
	workflowParams := WHWorkflowParam{
		TargetUrl:            "test_URL",
		RequestBody:          payload,
		AuthHeaderName:       "test_auth_header",
		AuthHeaderValue:      "test_auth_value",
		RetryPolicy:          nil,
		Notification:         nil,
		NotificationProvider: nil,
		NotifyFailure:        false,
	}
	s.env.OnActivity(activity.WebhookActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, targetUrl, authHeaderName, authHeaderValue string, reqBody string, notifyFailure bool, notification *notifications.NovuNotification, provider notifications.NotificationProvider) error {
			s.Equal("test_URL", targetUrl)
			s.Equal("test_auth_header", authHeaderName)
			s.Equal("test_auth_value", authHeaderValue)
			s.Equal(payload, reqBody)
			return nil
		})

	s.env.ExecuteWorkflow(WebhookWorkflow, workflowParams)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}
