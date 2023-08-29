package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphOrganizationEventHandler_OnRenewalLikelihoodUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	// prepare neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "new",
		LastName:  "user",
	})
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
		RenewalLikelihood: entity.RenewalLikelihood{
			RenewalLikelihood:         string(entity.RenewalLikelihoodZero),
			PreviousRenewalLikelihood: string(entity.RenewalLikelihoodHigh),
			Comment:                   utils.StringPtr("old comment"),
			UpdatedBy:                 "old user",
		},
	})
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"User": 1, "User_" + tenantName: 1, "Action": 0, "TimelineEvent": 0})

	// prepare event handler
	orgEventHandler := &GraphOrganizationEventHandler{
		Repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewOrganizationCommands(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationUpdateRenewalLikelihoodEvent(orgAggregate, models.RenewalLikelihoodLOW, models.RenewalLikelihoodHIGH, userId, utils.StringPtr("new comment"), now)
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnRenewalLikelihoodUpdate(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"User": 1, "User_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	orgDbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, orgDbNode)

	// verify organization
	organization := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, string(entity.RenewalLikelihoodZero), organization.RenewalLikelihood.PreviousRenewalLikelihood)
	require.Equal(t, string(entity.RenewalLikelihoodLow), organization.RenewalLikelihood.RenewalLikelihood)
	require.Equal(t, now, *organization.RenewalLikelihood.UpdatedAt)
	require.Equal(t, "new comment", *organization.RenewalLikelihood.Comment)
	require.Equal(t, userId, organization.RenewalLikelihood.UpdatedBy)
	require.Equal(t, entity.DataSourceOpenline, organization.SourceOfTruth)
	require.NotNil(t, organization.UpdatedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, now, action.CreatedAt)
	require.Equal(t, entity.ActionRenewalLikelihoodUpdated, action.Type)
	require.Equal(t, "Renewal likelihood set to Low by new user", action.Content)

	// Check request was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[orgAggregate.ID]
	require.Equal(t, 1, len(eventList))
	generatedEvent := eventList[0]
	require.Equal(t, events.OrganizationRequestRenewalForecastV1, generatedEvent.EventType)
	var eventData events.OrganizationRequestRenewalForecastEvent
	err = generatedEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	test.AssertRecentTime(t, eventData.RequestedAt)
	require.Equal(t, tenantName, eventData.Tenant)
}

func TestGraphOrganizationEventHandler_OnRenewalForecastUpdate_ByUser(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	// create neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "new",
		LastName:  "user",
	})
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
		RenewalForecast: entity.RenewalForecast{
			Amount:          utils.Float64Ptr(100),
			PotentialAmount: utils.Float64Ptr(200),
			Comment:         utils.StringPtr("old comment"),
		},
	})

	// prepare event handler
	orgEventHandler := &GraphOrganizationEventHandler{
		Repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewOrganizationCommands(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationUpdateRenewalForecastEvent(orgAggregate, utils.Float64Ptr(50), utils.Float64Ptr(60), utils.Float64Ptr(10), userId, utils.StringPtr("new comment"), now, "")
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnRenewalForecastUpdate(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"User": 1, "User_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify organization
	organization := graph_db.MapDbNodeToOrganizationEntity(*dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, float64(50), *organization.RenewalForecast.Amount)
	// potential should not be updated
	require.Equal(t, float64(200), *organization.RenewalForecast.PotentialAmount)
	require.Equal(t, now, *organization.RenewalForecast.UpdatedAt)
	require.Equal(t, "new comment", *organization.RenewalForecast.Comment)
	require.Equal(t, userId, organization.RenewalForecast.UpdatedBy)
	require.Equal(t, entity.DataSourceOpenline, organization.SourceOfTruth)
	require.NotNil(t, organization.UpdatedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, now, action.CreatedAt)
	require.Equal(t, entity.ActionRenewalForecastUpdated, action.Type)
	require.Equal(t, "Renewal forecast set to $50 by new user", action.Content)

	// Check request was not generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 0, len(eventsMap))
}

func TestGraphOrganizationEventHandler_OnRenewalForecastUpdate_ByInternalProcess(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	// create neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
		RenewalForecast: entity.RenewalForecast{
			Amount:          utils.Float64Ptr(100),
			PotentialAmount: utils.Float64Ptr(200),
			Comment:         utils.StringPtr("old comment"),
			UpdatedBy:       "old-user",
		},
	})

	// prepare event handler
	orgEventHandler := &GraphOrganizationEventHandler{
		Repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewOrganizationCommands(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationUpdateRenewalForecastEvent(orgAggregate, utils.Float64Ptr(5000), utils.Float64Ptr(10000), nil, "", utils.StringPtr("new comment"), now, models.RenewalLikelihoodMEDIUM)
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnRenewalForecastUpdate(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"Action": 1, "Action_" + tenantName: 1,
		"TimelineEvent": 1, "TimelineEvent_" + tenantName: 1})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify organization
	organization := graph_db.MapDbNodeToOrganizationEntity(*dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, float64(5000), *organization.RenewalForecast.Amount)
	require.Equal(t, float64(10000), *organization.RenewalForecast.PotentialAmount)
	require.Equal(t, now, *organization.RenewalForecast.UpdatedAt)
	require.Equal(t, "new comment", *organization.RenewalForecast.Comment)
	require.Equal(t, "", organization.RenewalForecast.UpdatedBy)
	require.Equal(t, entity.DataSourceOpenline, organization.SourceOfTruth)
	require.NotNil(t, organization.UpdatedAt)

	// verify action
	actionDbNode, err := neo4jt.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Action_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, actionDbNode)
	action := graph_db.MapDbNodeToActionEntity(*actionDbNode)
	require.NotNil(t, action.Id)
	require.Equal(t, entity.DataSource(constants.SourceOpenline), action.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, action.AppSource)
	require.Equal(t, now, action.CreatedAt)
	require.Equal(t, entity.ActionRenewalForecastUpdated, action.Type)
	require.Equal(t, "Renewal forecast set by default to $5,000, by discounting the billing amount using the renewal likelihood", action.Content)

	// Check request was not generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 0, len(eventsMap))
}

func TestGraphOrganizationEventHandler_OnRenewalForecastUpdate_ResetAmount(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	// create neo4j data
	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userId := neo4jt.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{
		FirstName: "new",
		LastName:  "user",
	})
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
		RenewalForecast: entity.RenewalForecast{
			Amount:          utils.Float64Ptr(100),
			PotentialAmount: utils.Float64Ptr(200),
			Comment:         utils.StringPtr("old comment"),
			UpdatedBy:       "old-user",
		},
	})

	// prepare event handler
	orgEventHandler := &GraphOrganizationEventHandler{
		Repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewOrganizationCommands(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	event, err := events.NewOrganizationUpdateRenewalForecastEvent(orgAggregate, nil, nil, nil, userId, utils.StringPtr("new comment"), now, "")
	require.Nil(t, err)

	// EXECUTE
	err = orgEventHandler.OnRenewalForecastUpdate(context.Background(), event)
	require.Nil(t, err)

	// no actions created
	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1,
		"User": 1, "User_" + tenantName: 1,
		"Action": 0, "Action_" + tenantName: 0,
		"TimelineEvent": 0, "TimelineEvent_" + tenantName: 0})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	organization := graph_db.MapDbNodeToOrganizationEntity(*dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Nil(t, organization.RenewalForecast.Amount)
	require.Equal(t, float64(200), *organization.RenewalForecast.PotentialAmount)
	require.Equal(t, now, *organization.RenewalForecast.UpdatedAt)
	require.Equal(t, "new comment", *organization.RenewalForecast.Comment)
	require.Equal(t, userId, organization.RenewalForecast.UpdatedBy)
	require.Equal(t, entity.DataSourceOpenline, organization.SourceOfTruth)
	require.NotNil(t, organization.UpdatedAt)

	// Check request was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[orgAggregate.ID]
	require.Equal(t, 1, len(eventList))
	generatedEvent := eventList[0]
	require.Equal(t, events.OrganizationRequestRenewalForecastV1, generatedEvent.EventType)
	var eventData events.OrganizationRequestRenewalForecastEvent
	err = generatedEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	test.AssertRecentTime(t, eventData.RequestedAt)
	require.Equal(t, tenantName, eventData.Tenant)
}

func TestGraphOrganizationEventHandler_OnBillingDetailsUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	hourAgo := utils.Now().Add(time.Duration(-1) * time.Hour)
	minAgo := utils.Now().Add(time.Duration(-1) * time.Minute)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
		BillingDetails: entity.BillingDetails{
			Amount:            utils.Float64Ptr(100),
			Frequency:         "WEEKLY",
			RenewalCycle:      "MONTHLY",
			RenewalCycleStart: utils.TimePtr(hourAgo),
			RenewalCycleNext:  utils.TimePtr(minAgo),
		},
	})
	orgEventHandler := &GraphOrganizationEventHandler{
		Repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewOrganizationCommands(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()

	event, err := events.NewOrganizationUpdateBillingDetailsEvent(orgAggregate, utils.Float64Ptr(50), "MONTHLY", "ANNUALLY", "new user", utils.TimePtr(now), utils.TimePtr(now))
	require.Nil(t, err)
	err = orgEventHandler.OnBillingDetailsUpdate(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	organization := graph_db.MapDbNodeToOrganizationEntity(*dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, float64(50), *organization.BillingDetails.Amount)
	require.Equal(t, "MONTHLY", organization.BillingDetails.Frequency)
	require.Equal(t, "ANNUALLY", organization.BillingDetails.RenewalCycle)
	require.Equal(t, now, *organization.BillingDetails.RenewalCycleStart)
	require.Equal(t, minAgo, *organization.BillingDetails.RenewalCycleNext)
	require.Equal(t, entity.DataSourceOpenline, organization.SourceOfTruth)
	require.NotNil(t, organization.UpdatedAt)

	// Check request was generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[orgAggregate.ID]
	require.Equal(t, 2, len(eventList))

	generatedEvent1 := eventList[0]
	require.Equal(t, events.OrganizationRequestRenewalForecastV1, generatedEvent1.EventType)
	var eventData1 events.OrganizationRequestRenewalForecastEvent
	err = generatedEvent1.GetJsonData(&eventData1)
	require.Nil(t, err)
	test.AssertRecentTime(t, eventData1.RequestedAt)
	require.Equal(t, tenantName, eventData1.Tenant)

	generatedEvent2 := eventList[1]
	require.Equal(t, events.OrganizationRequestNextCycleDateV1, generatedEvent2.EventType)
	var eventData2 events.OrganizationRequestNextCycleDateEvent
	err = generatedEvent2.GetJsonData(&eventData2)
	require.Nil(t, err)
	test.AssertRecentTime(t, eventData2.RequestedAt)
	require.Equal(t, tenantName, eventData2.Tenant)
}

func TestGraphOrganizationEventHandler_OnBillingDetailsUpdate_SetNotByUser(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()

	hourAgo := utils.Now().Add(time.Duration(-1) * time.Hour)
	minAgo := utils.Now().Add(time.Duration(-1) * time.Minute)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{
		Name: "test org",
		BillingDetails: entity.BillingDetails{
			Amount:            utils.Float64Ptr(100),
			Frequency:         "WEEKLY",
			RenewalCycle:      "MONTHLY",
			RenewalCycleStart: utils.TimePtr(hourAgo),
			RenewalCycleNext:  utils.TimePtr(minAgo),
		},
	})
	orgEventHandler := &GraphOrganizationEventHandler{
		Repositories:         testDatabase.Repositories,
		organizationCommands: command_handler.NewOrganizationCommands(testLogger, &config.Config{}, aggregateStore, testDatabase.Repositories),
	}
	orgAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenantName, orgId)
	now := utils.Now()
	tomorrow := now.Add(time.Duration(24) * time.Hour)

	event, err := events.NewOrganizationUpdateBillingDetailsEvent(orgAggregate, utils.Float64Ptr(50), "MONTHLY", "ANNUALLY", "", utils.TimePtr(now), utils.TimePtr(tomorrow))
	require.Nil(t, err)
	err = orgEventHandler.OnBillingDetailsUpdate(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Organization": 1, "Organization_" + tenantName: 1})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, orgId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	organization := graph_db.MapDbNodeToOrganizationEntity(*dbNode)
	require.Equal(t, orgId, organization.ID)
	require.Equal(t, float64(50), *organization.BillingDetails.Amount)
	require.Equal(t, "MONTHLY", organization.BillingDetails.Frequency)
	require.Equal(t, "ANNUALLY", organization.BillingDetails.RenewalCycle)
	require.Equal(t, now, *organization.BillingDetails.RenewalCycleStart)
	require.Equal(t, tomorrow, *organization.BillingDetails.RenewalCycleNext)
	require.Equal(t, entity.DataSourceOpenline, organization.SourceOfTruth)
	require.NotNil(t, organization.UpdatedAt)

	// Check request was not generated
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 0, len(eventsMap))
}
