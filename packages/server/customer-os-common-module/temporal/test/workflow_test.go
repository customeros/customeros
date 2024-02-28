package test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/activity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/workflows"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
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
	s.env.OnActivity(activity.WebhookActivity, mock.Anything, mock.Anything).Return(
		"", errors.New("WebhookActivityFailure"))
	s.env.ExecuteWorkflow(workflows.WebhookWorkflow, "test_failure")

	s.True(s.env.IsWorkflowCompleted())

	err := s.env.GetWorkflowError()
	s.Error(err)
	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("WebhookActivityFailure", applicationErr.Error())
}

func (s *UnitTestSuite) Test_WebhookWorkflow_ActivityParamCorrect() {
	workflowOptions := client.StartWorkflowOptions{
		ID:                       "test-webhook-calls_" + uuid.New().String(),
		WorkflowExecutionTimeout: time.Hour * 24 * 365 * 10,
		TaskQueue:                "test_task_queue",
	}

	workflowParams := workflows.WHWorkflowParam{
		TargetUrl:       "test_URL",
		RequestBody:     bytes.NewBuffer([]byte("test_body")),
		AuthHeaderName:  "test_auth_header",
		AuthHeaderValue: "test_auth_value",
		RetryPolicy:     nil,
	}
	s.env.OnActivity(activity.WebhookActivity, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, targetUrl, authHeaderName, authHeaderValue string, reqBody *bytes.Buffer) error {
			//s.Equal("test_success", value)
			s.Equal("test_URL", targetUrl)
			s.Equal("test_auth_header", authHeaderName)
			s.Equal("test_auth_value", authHeaderValue)
			s.Equal("test_body", reqBody.String())
			return nil
		})

	s.env.ExecuteWorkflow(context.Background(), workflowOptions, workflows.WebhookWorkflow, workflowParams)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}
