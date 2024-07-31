package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/mocked_grpc"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
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
		repositories: testDatabase.Repositories,
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
	isReachable := "emailIsNotReachable"
	creationTime := utils.Now()
	emailId := neo4jtest.CreateEmail(ctx, testDatabase.Driver, tenantName, neo4jentity.EmailEntity{
		Email:       emailCreate,
		RawEmail:    rawEmailCreate,
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
		IsReachable: &isReachable,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	dbNodeAfterEmailCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, emailId, utils.GetStringPropOrEmpty(propsAfterEmailCreate, "id"))

	emailEventHandler := &EmailEventHandler{
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	agg := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, emailId)
	tenant := agg.GetTenant()
	rawEmailUpdate := "email@update.com"
	sourceUpdate := constants.Anthropic
	updateTime := utils.Now()

	// prepare grpc mock
	calledEmailValidateRequest := false
	emailCallbacks := mocked_grpc.MockEmailServiceCallbacks{
		RequestEmailValidation: func(context context.Context, op *emailpb.RequestEmailValidationGrpcRequest) (*emailpb.EmailIdGrpcResponse, error) {
			require.Equal(t, tenantName, op.Tenant)
			require.Equal(t, emailId, op.Id)
			calledEmailValidateRequest = true
			return &emailpb.EmailIdGrpcResponse{}, nil
		},
	}
	mocked_grpc.SetEmailCallbacks(&emailCallbacks)

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
	require.Equal(t, "", utils.GetStringPropOrEmpty(emailProps, "isReachable"))
	require.Equal(t, "", utils.GetStringPropOrEmpty(emailProps, "email"))
	require.Equal(t, rawEmailUpdate, utils.GetStringPropOrEmpty(emailProps, "rawEmail"))
	require.Equal(t, creationTime, utils.GetTimePropOrNow(emailProps, "createdAt"))
	require.Less(t, creationTime, utils.GetTimePropOrNow(emailProps, "updatedAt"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(emailProps, "syncedWithEventStore"))

	require.True(t, calledEmailValidateRequest)
}

func TestGraphEmailEventHandler_OnEmailValidationFailed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	emailCreate := "email@create.com"
	rawEmailCreate := "email@create.com"
	isReachable := "emailIsNotReachable"
	creationTime := utils.Now()
	emailId := neo4jtest.CreateEmail(ctx, testDatabase.Driver, tenantName, neo4jentity.EmailEntity{
		Email:       emailCreate,
		RawEmail:    rawEmailCreate,
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
		IsReachable: &isReachable,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	dbNodeAfterEmailCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailtCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, emailId, utils.GetStringPropOrEmpty(propsAfterEmailtCreate, "id"))

	agg := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, emailId)
	validationError := "Email validation failed with this custom message!"
	event, err := event2.NewEmailFailedValidationEvent(agg, tenantName, validationError)
	require.Nil(t, err)

	emailEventHandler := &EmailEventHandler{
		repositories: testDatabase.Repositories,
	}
	err = emailEventHandler.OnEmailValidationFailed(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	emailProps := utils.GetPropsFromNode(*dbNode)
	require.Equal(t, 8, len(emailProps))
	emailNotReachable := "emailIsNotReachable"
	require.Equal(t, false, utils.GetBoolPropOrFalse(emailProps, "validated"))
	require.Equal(t, validationError, utils.GetStringPropOrEmpty(emailProps, "validationError"))
	require.Less(t, creationTime, *utils.GetTimePropOrNil(emailProps, "updatedAt"))
	require.Equal(t, creationTime, *utils.GetTimePropOrNil(emailProps, "createdAt"))
	require.Equal(t, emailNotReachable, utils.GetStringPropOrEmpty(emailProps, "isReachable"))
	require.Equal(t, emailCreate, utils.GetStringPropOrEmpty(emailProps, "email"))
	require.Equal(t, rawEmailCreate, utils.GetStringPropOrEmpty(emailProps, "rawEmail"))
	require.Equal(t, false, utils.GetBoolPropOrFalse(emailProps, "syncedWithEventStore"))
}

func TestGraphEmailEventHandler_OnEmailValidated(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	emailCreate := "email@create.com"
	rawEmailCreate := "email@create.com"
	isReachable := "true"
	creationTime := utils.Now()
	acceptsMail := true
	canConnectSmtp := true
	hasFullInbox := true
	isCatchAll := true
	isDisabled := true
	isDeliverable := true
	isDisposable := true
	isRoleAccount := true
	emailId := neo4jtest.CreateEmail(ctx, testDatabase.Driver, tenantName, neo4jentity.EmailEntity{
		Email:          emailCreate,
		RawEmail:       rawEmailCreate,
		Primary:        false,
		CreatedAt:      creationTime,
		UpdatedAt:      creationTime,
		IsReachable:    &isReachable,
		CanConnectSMTP: &canConnectSmtp,
		AcceptsMail:    &acceptsMail,
		HasFullInbox:   &hasFullInbox,
		IsCatchAll:     &isCatchAll,
		IsDisabled:     &isDisabled,
		IsDeliverable:  &isDeliverable,
		IsDisposable:   &isDisposable,
		IsRoleAccount:  &isRoleAccount,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	dbNodeAfterEmailCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailtCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, emailId, utils.GetStringPropOrEmpty(propsAfterEmailtCreate, "id"))

	agg := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, emailId)
	validationError := "Email validation failed with this custom message!"
	domain := "emailUpdateDomain"
	username := "emailUsername"
	event, err := event2.NewEmailValidatedEvent(agg, tenantName, rawEmailCreate, isReachable, validationError, domain, username, emailCreate, acceptsMail, canConnectSmtp, hasFullInbox, isCatchAll, isDisabled, true, isDeliverable, isDisposable, isRoleAccount)
	require.Nil(t, err)

	emailEventHandler := &EmailEventHandler{
		repositories: testDatabase.Repositories,
	}
	err = emailEventHandler.OnEmailValidated(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "acceptsMail"))
	require.Equal(t, username, utils.GetStringPropOrEmpty(props, "username"))
	require.Equal(t, rawEmailCreate, utils.GetStringPropOrEmpty(props, "rawEmail"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "isValidSyntax"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "canConnectSmtp"))
	require.Equal(t, isReachable, utils.GetStringPropOrEmpty(props, "isReachable"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "isCatchAll"))
	require.Equal(t, validationError, utils.GetStringPropOrEmpty(props, "validationError"))
	require.Equal(t, emailCreate, utils.GetStringPropOrEmpty(props, "email"))
	require.Less(t, creationTime, utils.GetTimePropOrNow(props, "updatedAt"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "hasFullInbox"))
	require.Equal(t, creationTime, *utils.GetTimePropOrNil(props, "createdAt"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "validated"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "isDisabled"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "isDeliverable"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "isDisposable"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "isRoleAccount"))
}
