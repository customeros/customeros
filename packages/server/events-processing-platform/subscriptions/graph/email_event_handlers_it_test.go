package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	emailAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	emailEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphEmailEventHandler_OnEmailCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	emailEventHandler := &EmailEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myMailId, _ := uuid.NewUUID()
	emailAggregate := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, myMailId.String())
	email := "test@test.com"
	curTime := time.Now().UTC()
	event, err := emailEvents.NewEmailCreateEvent(emailAggregate, tenantName, email, cmnmod.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "unit-test",
	}, curTime, curTime)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "EMAIL_ADDRESS_BELONGS_TO_TENANT"), "Incorrect number of EMAIL_ADDRESS_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, myMailId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myMailId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, email, utils.GetStringPropOrEmpty(props, "rawEmail"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(props, "appSource"))
}

func TestGraphEmailEventHandler_OnEmailUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)

	emailCreate := "email@create.com"
	rawEmailCreate := "email@create.com"
	isReachable := "emailIsNotReachable"
	creationTime := utils.Now()
	emailId := neo4jt.CreateEmail(ctx, testDatabase.Driver, tenantName, entity.EmailEntity{
		Email:       emailCreate,
		RawEmail:    rawEmailCreate,
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
		IsReachable: &isReachable,
	})

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	dbNodeAfterEmailCreate, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, emailId, utils.GetStringPropOrEmpty(propsAfterEmailCreate, "id"))

	emailEventHandler := &EmailEventHandler{
		Repositories: testDatabase.Repositories,
	}
	emailAggregate := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, emailId)
	tenant := emailAggregate.GetTenant()
	rawEmailUpdate := "email@update.com"
	sourceUpdate := constants.Anthropic
	updateTime := utils.Now()
	event, err := emailEvents.NewEmailUpdateEvent(emailAggregate, rawEmailUpdate, tenant, sourceUpdate, updateTime)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailUpdate(context.Background(), event)
	require.Nil(t, err)
	email, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Email_"+tenantName)
	require.Nil(t, err)

	emailProps := utils.GetPropsFromNode(*email)
	require.Equal(t, 7, len(emailProps))
	emailId = utils.GetStringPropOrEmpty(emailProps, "id")
	require.NotNil(t, emailId)
	emailNotReachable := "emailIsNotReachable"
	require.Equal(t, emailNotReachable, utils.GetStringPropOrEmpty(emailProps, "isReachable"))
	require.Equal(t, emailCreate, utils.GetStringPropOrEmpty(emailProps, "email"))
	require.Equal(t, rawEmailCreate, utils.GetStringPropOrEmpty(emailProps, "rawEmail"))
	require.Equal(t, creationTime, utils.GetTimePropOrNow(emailProps, "createdAt"))
	require.Less(t, creationTime, utils.GetTimePropOrNow(emailProps, "updatedAt"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(emailProps, "syncedWithEventStore"))
}

func TestGraphEmailEventHandler_OnEmailValidationFailed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)

	emailCreate := "email@create.com"
	rawEmailCreate := "email@create.com"
	isReachable := "emailIsNotReachable"
	creationTime := utils.Now()
	emailId := neo4jt.CreateEmail(ctx, testDatabase.Driver, tenantName, entity.EmailEntity{
		Email:       emailCreate,
		RawEmail:    rawEmailCreate,
		CreatedAt:   creationTime,
		UpdatedAt:   creationTime,
		IsReachable: &isReachable,
	})

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	dbNodeAfterEmailCreate, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailtCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, emailId, utils.GetStringPropOrEmpty(propsAfterEmailtCreate, "id"))

	emailAggregate := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, emailId)
	validationError := "Email validation failed with this custom message!"
	event, err := emailEvents.NewEmailFailedValidationEvent(emailAggregate, tenantName, validationError)
	require.Nil(t, err)

	emailEventHandler := &EmailEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = emailEventHandler.OnEmailValidationFailed(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
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

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)

	emailCreate := "email@create.com"
	rawEmailCreate := "email@create.com"
	isReachable := "true"
	creationTime := utils.Now()
	acceptsMail := true
	canConnectSmtp := true
	hasFullInbox := true
	isCatchAll := true
	IsDeliverable := true
	isDisabled := true
	emailId := neo4jt.CreateEmail(ctx, testDatabase.Driver, tenantName, entity.EmailEntity{
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
		IsDeliverable:  &IsDeliverable,
		IsDisabled:     &isDisabled,
	})

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Email": 1, "Email_" + tenantName: 1})
	dbNodeAfterEmailCreate, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailtCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, emailId, utils.GetStringPropOrEmpty(propsAfterEmailtCreate, "id"))

	emailAggregate := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, emailId)
	validationError := "Email validation failed with this custom message!"
	domain := "emailUpdateDomain"
	username := "emailUsername"
	isValidSyntax := true
	event, err := emailEvents.NewEmailValidatedEvent(emailAggregate, tenantName, rawEmailCreate, isReachable, validationError, domain, username, emailCreate, acceptsMail, canConnectSmtp, hasFullInbox, isCatchAll, IsDeliverable, isDisabled, isValidSyntax)
	require.Nil(t, err)

	emailEventHandler := &EmailEventHandler{
		Repositories: testDatabase.Repositories,
	}
	err = emailEventHandler.OnEmailValidated(context.Background(), event)
	require.Nil(t, err)

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "acceptsMail"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "isDeliverable"))
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
}
