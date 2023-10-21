package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	job_role_aggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/aggregate"
	job_role_model "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	job_role_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	user_aggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	user_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphUserEventHandler_OnUserCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &GraphUserEventHandler{
		repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	curTime := time.Now().UTC()

	event, err := user_events.NewUserCreateEvent(userAggregate, models.UserDataFields{
		FirstName:       "Bob",
		LastName:        "Dole",
		Name:            "Bob Dole",
		Internal:        true,
		ProfilePhotoUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
		Timezone:        "Europe/Paris",
	},
		cmnmod.Source{
			Source:        "N/A",
			SourceOfTruth: "N/A",
			AppSource:     "unit-test",
		},
		cmnmod.ExternalSystem{},
		curTime, curTime)
	require.Nil(t, err)
	err = userEventHandler.OnUserCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "User_"+tenantName), "Incorrect number of User_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "USER_BELONGS_TO_TENANT"), "Incorrect number of USER_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myUserId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(props, "firstName"))
	require.Equal(t, "Dole", utils.GetStringPropOrEmpty(props, "lastName"))
	require.Equal(t, "Bob Dole", utils.GetStringPropOrEmpty(props, "name"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(props, "source"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(props, "appSource"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "syncedWithEventStore"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "internal"))
	require.Equal(t, "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", utils.GetStringPropOrEmpty(props, "profilePhotoUrl"))
	require.Equal(t, "Europe/Paris", utils.GetStringPropOrEmpty(props, "timezone"))
}

func TestGraphUserEventHandler_OnUserCreate_WithExternalSystem(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 0, "ExternalSystem": 1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"IS_LINKED_WITH": 0,
	})

	userEventHandler := &GraphUserEventHandler{
		repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	curTime := utils.Now()

	event, err := user_events.NewUserCreateEvent(userAggregate, models.UserDataFields{
		FirstName:       "Bob",
		LastName:        "Dole",
		Name:            "Bob Dole",
		Internal:        true,
		ProfilePhotoUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
		Timezone:        "Europe/Paris",
	}, cmnmod.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "unit-test",
	}, cmnmod.ExternalSystem{
		ExternalSystemId: "sf",
		ExternalId:       "123",
		ExternalIdSecond: "ABC",
	}, curTime, curTime)
	require.Nil(t, err)
	err = userEventHandler.OnUserCreate(context.Background(), event)
	require.Nil(t, err)

	neo4jt.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User":               1,
		"User_" + tenantName: 1,
		"ExternalSystem":     1})
	neo4jt.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"IS_LINKED_WITH":         1,
		"USER_BELONGS_TO_TENANT": 1,
	})

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myUserId.String(), utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(props, "firstName"))
	require.Equal(t, "Dole", utils.GetStringPropOrEmpty(props, "lastName"))
	require.Equal(t, "Bob Dole", utils.GetStringPropOrEmpty(props, "name"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(props, "source"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(props, "appSource"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "syncedWithEventStore"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "internal"))
	require.Equal(t, "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", utils.GetStringPropOrEmpty(props, "profilePhotoUrl"))
	require.Equal(t, "Europe/Paris", utils.GetStringPropOrEmpty(props, "timezone"))
}

func TestGraphUserEventHandler_OnUserCreateWithJobRole(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &GraphUserEventHandler{
		repositories: testDatabase.Repositories,
	}
	jobRoleEventHandler := &GraphJobRoleEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	myJobRoleId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	jobRoleAggregate := job_role_aggregate.NewJobRoleAggregateWithTenantAndID(tenantName, myJobRoleId.String())

	curTime := time.Now().UTC()

	description := "I clean things"

	userCreateEvent, err := user_events.NewUserCreateEvent(userAggregate, models.UserDataFields{
		FirstName:       "Bob",
		LastName:        "Dole",
		Name:            "Bob Dole",
		ProfilePhotoUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
		Timezone:        "Africa/Abidjan",
	},
		cmnmod.Source{
			Source:        "N/A",
			SourceOfTruth: "N/A",
			AppSource:     "unit-test",
		},
		cmnmod.ExternalSystem{},
		curTime, curTime)
	require.Nil(t, err)
	err = userEventHandler.OnUserCreate(context.Background(), userCreateEvent)
	require.Nil(t, err)

	now := utils.Now()
	jobRoleCreateEvent, err := job_role_events.NewJobRoleCreateEvent(jobRoleAggregate,
		job_role_model.NewCreateJobRoleCommand(myJobRoleId.String(),
			tenantName, "Chief Janitor", &description,
			false, "N/A", "N/A", "unit-test", &now, nil, &curTime))

	require.Nil(t, err)
	err = jobRoleEventHandler.OnJobRoleCreate(context.Background(), jobRoleCreateEvent)
	require.Nil(t, err)

	linkJobRoleEvent, err := user_events.NewUserLinkJobRoleEvent(userAggregate, tenantName, myJobRoleId.String(), curTime)
	require.Nil(t, err)
	err = userEventHandler.OnJobRoleLinkedToUser(context.Background(), linkJobRoleEvent)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "User_"+tenantName), "Incorrect number of User_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "USER_BELONGS_TO_TENANT"), "Incorrect number of USER_BELONGS_TO_TENANT relationships in Neo4j")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole_"+tenantName), "Incorrect number of JobRole_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "WORKS_AS"), "Incorrect number of WORKS_AS relationships in Neo4j")

	dbUserNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
	require.Nil(t, err)
	require.NotNil(t, dbUserNode)
	userProps := utils.GetPropsFromNode(*dbUserNode)

	require.Equal(t, myUserId.String(), utils.GetStringPropOrEmpty(userProps, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(userProps, "firstName"))
	require.Equal(t, "Dole", utils.GetStringPropOrEmpty(userProps, "lastName"))
	require.Equal(t, "Bob Dole", utils.GetStringPropOrEmpty(userProps, "name"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(userProps, "source"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(userProps, "sourceOfTruth"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(userProps, "appSource"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(userProps, "syncedWithEventStore"))
	require.Equal(t, "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", utils.GetStringPropOrEmpty(userProps, "profilePhotoUrl"))
	require.Equal(t, "Africa/Abidjan", utils.GetStringPropOrEmpty(userProps, "timezone"))

	dbJobRoleNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "JobRole_"+tenantName, myJobRoleId.String())
	if err != nil {
		t.Fatalf("Error getting JobRole node from Neo4j: %s", err.Error())
	}
	require.Nil(t, err)
	require.NotNil(t, dbJobRoleNode)
	jobRoleProps := utils.GetPropsFromNode(*dbJobRoleNode)

	require.Equal(t, myJobRoleId.String(), utils.GetStringPropOrEmpty(jobRoleProps, "id"))
	require.Equal(t, "Chief Janitor", utils.GetStringPropOrEmpty(jobRoleProps, "jobTitle"))
	require.Equal(t, description, utils.GetStringPropOrEmpty(jobRoleProps, "description"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(jobRoleProps, "appSource"))
}

func TestGraphUserEventHandler_OnUserCreateWithJobRoleOutOfOrder(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &GraphUserEventHandler{
		repositories: testDatabase.Repositories,
	}
	jobRoleEventHandler := &GraphJobRoleEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	myJobRoleId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	jobRoleAggregate := job_role_aggregate.NewJobRoleAggregateWithTenantAndID(tenantName, myJobRoleId.String())

	curTime := time.Now().UTC()

	description := "I clean things"

	userCreateEvent, err := user_events.NewUserCreateEvent(userAggregate, models.UserDataFields{
		FirstName: "Bob",
		LastName:  "Dole",
		Name:      "Bob Dole",
	}, cmnmod.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "unit-test",
	}, cmnmod.ExternalSystem{}, curTime, curTime)
	require.Nil(t, err)
	err = userEventHandler.OnUserCreate(context.Background(), userCreateEvent)
	require.Nil(t, err)

	linkJobRoleEvent, err := user_events.NewUserLinkJobRoleEvent(userAggregate, tenantName, myJobRoleId.String(), curTime)
	require.Nil(t, err)
	err = userEventHandler.OnJobRoleLinkedToUser(context.Background(), linkJobRoleEvent)
	require.Nil(t, err)

	jobRoleCreateEvent, err := job_role_events.NewJobRoleCreateEvent(jobRoleAggregate,
		job_role_model.NewCreateJobRoleCommand(myJobRoleId.String(),
			tenantName, "Chief Janitor", &description,
			false, "N/A", "N/A", "unit-test", nil, nil, &curTime))

	require.Nil(t, err)
	err = jobRoleEventHandler.OnJobRoleCreate(context.Background(), jobRoleCreateEvent)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "User"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "User_"+tenantName), "Incorrect number of User_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "USER_BELONGS_TO_TENANT"), "Incorrect number of USER_BELONGS_TO_TENANT relationships in Neo4j")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole_"+tenantName), "Incorrect number of JobRole_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "WORKS_AS"), "Incorrect number of WORKS_AS relationships in Neo4j")

	dbUserNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
	require.Nil(t, err)
	require.NotNil(t, dbUserNode)
	userProps := utils.GetPropsFromNode(*dbUserNode)

	require.Equal(t, myUserId.String(), utils.GetStringPropOrEmpty(userProps, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(userProps, "firstName"))
	require.Equal(t, "Dole", utils.GetStringPropOrEmpty(userProps, "lastName"))
	require.Equal(t, "Bob Dole", utils.GetStringPropOrEmpty(userProps, "name"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(userProps, "source"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(userProps, "sourceOfTruth"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(userProps, "appSource"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(userProps, "syncedWithEventStore"))

	dbJobRoleNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "JobRole_"+tenantName, myJobRoleId.String())
	if err != nil {
		t.Fatalf("Error getting JobRole node from Neo4j: %s", err.Error())
	}
	require.Nil(t, err)
	require.NotNil(t, dbJobRoleNode)
	jobRoleProps := utils.GetPropsFromNode(*dbJobRoleNode)

	require.Equal(t, myJobRoleId.String(), utils.GetStringPropOrEmpty(jobRoleProps, "id"))
	require.Equal(t, "Chief Janitor", utils.GetStringPropOrEmpty(jobRoleProps, "jobTitle"))
	require.Equal(t, description, utils.GetStringPropOrEmpty(jobRoleProps, "description"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(jobRoleProps, "appSource"))
}
