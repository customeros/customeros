package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	job_role_aggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/aggregate"
	job_role_model "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	job_role_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	user_aggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	user_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	user_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphUserEventHandler_OnUserCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	curTime := time.Now().UTC()

	event, err := user_events.NewUserCreateEvent(userAggregate, user_models.UserDataFields{
		FirstName:       "Bob",
		LastName:        "Dole",
		Name:            "Bob Dole",
		Internal:        true,
		Bot:             true,
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

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "User_"+tenantName), "Incorrect number of User_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "USER_BELONGS_TO_TENANT"), "Incorrect number of USER_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
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
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "bot"))
	require.Equal(t, "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", utils.GetStringPropOrEmpty(props, "profilePhotoUrl"))
	require.Equal(t, "Europe/Paris", utils.GetStringPropOrEmpty(props, "timezone"))
}

func TestGraphUserEventHandler_OnUserCreate_WithExternalSystem(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	neo4jt.CreateExternalSystem(ctx, testDatabase.Driver, tenantName, "sf")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 0, "ExternalSystem": 1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"IS_LINKED_WITH": 0,
	})

	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	curTime := utils.Now()

	event, err := user_events.NewUserCreateEvent(userAggregate, user_models.UserDataFields{
		FirstName:       "Bob",
		LastName:        "Dole",
		Name:            "Bob Dole",
		Internal:        true,
		Bot:             true,
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

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"User":               1,
		"User_" + tenantName: 1,
		"ExternalSystem":     1})
	neo4jtest.AssertNeo4jRelationCount(ctx, t, testDatabase.Driver, map[string]int{
		"IS_LINKED_WITH":         1,
		"USER_BELONGS_TO_TENANT": 1,
	})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
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
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "bot"))
	require.Equal(t, "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", utils.GetStringPropOrEmpty(props, "profilePhotoUrl"))
	require.Equal(t, "Europe/Paris", utils.GetStringPropOrEmpty(props, "timezone"))
}

func TestGraphUserEventHandler_OnUserCreateWithJobRole(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	jobRoleEventHandler := &JobRoleEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	myJobRoleId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	jobRoleAggregate := job_role_aggregate.NewJobRoleAggregateWithTenantAndID(tenantName, myJobRoleId.String())

	curTime := time.Now().UTC()

	description := "I clean things"

	userCreateEvent, err := user_events.NewUserCreateEvent(userAggregate, user_models.UserDataFields{
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

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "User_"+tenantName), "Incorrect number of User_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "USER_BELONGS_TO_TENANT"), "Incorrect number of USER_BELONGS_TO_TENANT relationships in Neo4j")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole_"+tenantName), "Incorrect number of JobRole_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "WORKS_AS"), "Incorrect number of WORKS_AS relationships in Neo4j")

	dbUserNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
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

	dbJobRoleNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "JobRole_"+tenantName, myJobRoleId.String())
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
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	jobRoleEventHandler := &JobRoleEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myUserId, _ := uuid.NewUUID()
	myJobRoleId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	jobRoleAggregate := job_role_aggregate.NewJobRoleAggregateWithTenantAndID(tenantName, myJobRoleId.String())

	curTime := time.Now().UTC()

	description := "I clean things"

	userCreateEvent, err := user_events.NewUserCreateEvent(userAggregate, user_models.UserDataFields{
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

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "User_"+tenantName), "Incorrect number of User_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "USER_BELONGS_TO_TENANT"), "Incorrect number of USER_BELONGS_TO_TENANT relationships in Neo4j")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "JobRole_"+tenantName), "Incorrect number of JobRole_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "WORKS_AS"), "Incorrect number of WORKS_AS relationships in Neo4j")

	dbUserNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, myUserId.String())
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

	dbJobRoleNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "JobRole_"+tenantName, myJobRoleId.String())
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

func TestGraphUserEventHandler_OnUserUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userCreateTime := utils.Now()
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName:        "UserFirstNameCreate",
		LastName:         "UserLastNameCreate",
		CreatedAt:        userCreateTime,
		UpdatedAt:        userCreateTime,
		Source:           constants.SourceOpenline,
		SourceOfTruth:    constants.SourceOpenline,
		AppSource:        constants.AppSourceEventProcessingPlatform,
		Roles:            []string{model.RoleUser.String(), model.RoleOwner.String()},
		ProfilePhotoUrl:  "www.photo.com/create",
		Timezone:         "userTimezoneCreate",
		Internal:         false,
		Bot:              false,
		DefaultForPlayer: false,
		Tenant:           tenantName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "User_" + tenantName: 1})
	dbNodeAfterUserCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, userId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterUserCreate)
	propsAfterUserCreate := utils.GetPropsFromNode(*dbNodeAfterUserCreate)
	require.Equal(t, userId, utils.GetStringPropOrEmpty(propsAfterUserCreate, "id"))

	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	userUpdateTime := utils.Now()
	userNameUpdate := "UserNameUpdate"
	userFirstNameUpdate := "UserFirstNameUpdate"
	userLastNameUpdate := "UserlastNameUpdate"
	userProfilePhotoUrlUpdate := "www.photo.com/update"
	userTimezoneUpdate := "userTimezoneUpdate"
	dataFields := user_models.UserDataFields{
		Name:            userNameUpdate,
		FirstName:       userFirstNameUpdate,
		LastName:        userLastNameUpdate,
		Internal:        true,
		Bot:             true,
		ProfilePhotoUrl: userProfilePhotoUrlUpdate,
		Timezone:        userTimezoneUpdate,
	}

	event, err := user_events.NewUserUpdateEvent(userAggregate, dataFields, constants.SourceOpenline, userUpdateTime, cmnmod.ExternalSystem{})
	require.Nil(t, err)
	err = userEventHandler.OnUserUpdate(context.Background(), event)
	require.Nil(t, err)
	user, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "User_"+tenantName)
	require.Nil(t, err)

	userProps := utils.GetPropsFromNode(*user)
	require.Equal(t, 10, len(userProps))

	userId = utils.GetStringPropOrEmpty(userProps, "id")
	require.NotNil(t, userId)
	require.Equal(t, userNameUpdate, utils.GetStringPropOrEmpty(userProps, "name"))
	require.Equal(t, userFirstNameUpdate, utils.GetStringPropOrEmpty(userProps, "firstName"))
	require.Equal(t, userLastNameUpdate, utils.GetStringPropOrEmpty(userProps, "lastName"))
	require.Equal(t, userTimezoneUpdate, utils.GetStringPropOrEmpty(userProps, "timezone"))
	require.Equal(t, userProfilePhotoUrlUpdate, utils.GetStringPropOrEmpty(userProps, "profilePhotoUrl"))
	require.Less(t, userCreateTime, utils.GetTimePropOrNow(userProps, "updatedAt"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(userProps, "sourceOfTruth"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(userProps, "syncedWithEventStore"))
	require.Equal(t, 2, len(utils.GetListStringPropOrEmpty(userProps, "roles")))
	require.Contains(t, utils.GetListStringPropOrEmpty(userProps, "roles"), "OWNER")
	require.Contains(t, utils.GetListStringPropOrEmpty(userProps, "roles"), "USER")

}

func TestGraphUserEventHandler_OnPhoneNumberLinkedToUser(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userCreateTime := utils.Now()
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName:        "UserFirstNameCreate",
		LastName:         "UserLastNameCreate",
		CreatedAt:        userCreateTime,
		UpdatedAt:        userCreateTime,
		Source:           constants.SourceOpenline,
		SourceOfTruth:    constants.SourceOpenline,
		AppSource:        constants.AppSourceEventProcessingPlatform,
		Roles:            []string{model.RoleUser.String(), model.RoleOwner.String()},
		ProfilePhotoUrl:  "www.photo.com/create",
		Timezone:         "userTimezoneCreate",
		Internal:         false,
		Bot:              false,
		DefaultForPlayer: false,
		Tenant:           tenantName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "User_" + tenantName: 1})
	dbNodeAfterUserCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, userId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterUserCreate)
	propsAfterUserCreate := utils.GetPropsFromNode(*dbNodeAfterUserCreate)
	require.Equal(t, userId, utils.GetStringPropOrEmpty(propsAfterUserCreate, "id"))

	validated := false
	e164 := "+0123456789"

	phoneNumberId := neo4jt.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, entity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterPhoneNumberCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterPhoneNumberCreate)

	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	phoneNumberLabel := "phoneNumberLabel"
	userLinkPhoneNumberTime := utils.Now()
	userLinkPhoneNumberEvent, err := user_events.NewUserLinkPhoneNumberEvent(userAggregate, tenantName, phoneNumberId, phoneNumberLabel, true, userLinkPhoneNumberTime)
	require.Nil(t, err)
	err = userEventHandler.OnPhoneNumberLinkedToUser(context.Background(), userLinkPhoneNumberEvent)
	require.Nil(t, err)
	userNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "User_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, userNode)
	userProps := utils.GetPropsFromNode(*userNode)
	require.Less(t, userCreateTime, utils.GetTimePropOrNow(userProps, "updatedAt"))

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS"), "Incorrect number of PHONE_NUMBER_BELONGS_TO_TENANT relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, userId, "HAS", phoneNumberId)
	userPhoneRelation, err := neo4jtest.GetRelationship(ctx, testDatabase.Driver, userId, phoneNumberId)
	require.Nil(t, err)
	userPhoneRelationProps := utils.GetPropsFromRelationship(*userPhoneRelation)
	require.Equal(t, true, utils.GetBoolPropOrFalse(userPhoneRelationProps, "primary"))
	require.Equal(t, phoneNumberLabel, utils.GetStringPropOrEmpty(userPhoneRelationProps, "label"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(userPhoneRelationProps, "syncedWithEventStore"))
}

func TestGraphUserEventHandler_OnEmailLinkedToUser(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userCreateTime := utils.Now()
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName:        "UserFirstNameCreate",
		LastName:         "UserLastNameCreate",
		CreatedAt:        userCreateTime,
		UpdatedAt:        userCreateTime,
		Source:           constants.SourceOpenline,
		SourceOfTruth:    constants.SourceOpenline,
		AppSource:        constants.AppSourceEventProcessingPlatform,
		Roles:            []string{model.RoleUser.String(), model.RoleOwner.String()},
		ProfilePhotoUrl:  "www.photo.com/create",
		Timezone:         "userTimezoneCreate",
		Internal:         false,
		Bot:              false,
		DefaultForPlayer: false,
		Tenant:           tenantName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "User_" + tenantName: 1})
	dbNodeAfterUserCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, userId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterUserCreate)
	propsAfterUserCreate := utils.GetPropsFromNode(*dbNodeAfterUserCreate)
	require.Equal(t, userId, utils.GetStringPropOrEmpty(propsAfterUserCreate, "id"))

	primary := true
	email := "email@website.com"
	emailId := neo4jt.CreateEmail(ctx, testDatabase.Driver, tenantName, entity.EmailEntity{
		Email:         email,
		RawEmail:      email,
		Primary:       primary,
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
		AppSource:     constants.SourceOpenline,
	})

	dbNodeAfterEmailCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, false, utils.GetBoolPropOrFalse(propsAfterEmailCreate, "primary"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, &creationTime, utils.GetTimePropOrNil(propsAfterEmailCreate, "updatedAt"))

	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	emailLabel := "emailLabel"
	userLinkEmailTime := utils.Now()
	userLinkEmailEvent, err := user_events.NewUserLinkEmailEvent(userAggregate, tenantName, emailId, emailLabel, true, userLinkEmailTime)
	require.Nil(t, err)
	err = userEventHandler.OnEmailLinkedToUser(context.Background(), userLinkEmailEvent)
	require.Nil(t, err)
	userNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "User_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, userNode)
	userProps := utils.GetPropsFromNode(*userNode)
	require.Less(t, userCreateTime, utils.GetTimePropOrNow(userProps, "updatedAt"))

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS"), "Incorrect number of PHONE_NUMBER_BELONGS_TO_TENANT relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, userId, "HAS", emailId)
	userEmailRelation, err := neo4jtest.GetRelationship(ctx, testDatabase.Driver, userId, emailId)
	require.Nil(t, err)
	userEmailRelationProps := utils.GetPropsFromRelationship(*userEmailRelation)
	require.Equal(t, 3, len(userEmailRelationProps))
	require.Equal(t, true, utils.GetBoolPropOrFalse(userEmailRelationProps, "primary"))
	require.Equal(t, emailLabel, utils.GetStringPropOrEmpty(userEmailRelationProps, "label"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(userEmailRelationProps, "syncedWithEventStore"))
}

func TestGraphUserEventHandler_OnJobRoleLinkedToUser(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userCreateTime := utils.Now()
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName:        "UserFirstNameCreate",
		LastName:         "UserLastNameCreate",
		CreatedAt:        userCreateTime,
		UpdatedAt:        userCreateTime,
		Source:           constants.SourceOpenline,
		SourceOfTruth:    constants.SourceOpenline,
		AppSource:        constants.AppSourceEventProcessingPlatform,
		Roles:            []string{model.RoleUser.String(), model.RoleOwner.String()},
		ProfilePhotoUrl:  "www.photo.com/create",
		Timezone:         "userTimezoneCreate",
		Internal:         false,
		Bot:              false,
		DefaultForPlayer: false,
		Tenant:           tenantName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "User_" + tenantName: 1})
	dbNodeAfterUserCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, userId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterUserCreate)
	propsAfterUserCreate := utils.GetPropsFromNode(*dbNodeAfterUserCreate)
	require.Equal(t, userId, utils.GetStringPropOrEmpty(propsAfterUserCreate, "id"))

	jobRoleId := neo4jt.CreateJobRole(ctx, testDatabase.Driver, tenantName, entity.JobRoleEntity{})

	dbNodeAfterjobRoleCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "JobRole_"+tenantName, jobRoleId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterjobRoleCreate)
	propsAfterJobRoleCreate := utils.GetPropsFromNode(*dbNodeAfterjobRoleCreate)
	require.NotNil(t, utils.GetStringPropOrEmpty(propsAfterJobRoleCreate, "id"))

	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	userLinkJobRoleTime := utils.Now()
	userLinkJobRoleEvent, err := user_events.NewUserLinkJobRoleEvent(userAggregate, tenantName, jobRoleId, userLinkJobRoleTime)
	require.Nil(t, err)
	err = userEventHandler.OnJobRoleLinkedToUser(context.Background(), userLinkJobRoleEvent)
	require.Nil(t, err)

	userNode, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "User_"+tenantName)
	require.Nil(t, err)
	require.NotNil(t, userNode)
	userProps := utils.GetPropsFromNode(*userNode)
	require.Less(t, userCreateTime, utils.GetTimePropOrNow(userProps, "updatedAt"))

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "WORKS_AS"), "Incorrect number of WORKS_AS relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, userId, "WORKS_AS", jobRoleId)
}

func TestGraphUserEventHandler_OnAddPlayer(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userCreateTime := utils.Now()
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName:        "UserFirstNameCreate",
		LastName:         "UserLastNameCreate",
		CreatedAt:        userCreateTime,
		UpdatedAt:        userCreateTime,
		Source:           constants.SourceOpenline,
		SourceOfTruth:    constants.SourceOpenline,
		AppSource:        constants.AppSourceEventProcessingPlatform,
		Roles:            []string{model.RoleUser.String(), model.RoleOwner.String()},
		ProfilePhotoUrl:  "www.photo.com/create",
		Timezone:         "userTimezoneCreate",
		Internal:         false,
		Bot:              false,
		DefaultForPlayer: false,
		Tenant:           tenantName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "User_" + tenantName: 1})
	dbNodeAfterUserCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, userId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterUserCreate)
	propsAfterUserCreate := utils.GetPropsFromNode(*dbNodeAfterUserCreate)
	require.Equal(t, userId, utils.GetStringPropOrEmpty(propsAfterUserCreate, "id"))

	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	playerCreateTime := utils.Now()
	playerInfoDataFields := user_models.PlayerInfo{
		Provider:   "PlayerInfoProvider",
		AuthId:     "PlayerInfoAuthId",
		IdentityId: "PlayerInfIdentityId",
	}

	playerInfoEvent, err := user_events.NewUserAddPlayerInfoEvent(userAggregate, playerInfoDataFields, cmnmod.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "unit-test",
	},
		playerCreateTime)
	require.Nil(t, err)
	err = userEventHandler.OnAddPlayer(context.Background(), playerInfoEvent)
	require.Nil(t, err)
	player, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Player")
	require.Nil(t, err)

	playerProps := utils.GetPropsFromNode(*player)
	require.Equal(t, 9, len(playerProps))

	playerId := utils.GetStringPropOrEmpty(playerProps, "id")
	require.NotNil(t, playerId)
	require.Equal(t, playerCreateTime, utils.GetTimePropOrNow(playerProps, "createdAt"))
	require.Less(t, playerCreateTime, utils.GetTimePropOrNow(playerProps, "updatedAt"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(playerProps, "sourceOfTruth"))
	require.Equal(t, "N/A", utils.GetStringPropOrEmpty(playerProps, "source"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(playerProps, "appSource"))
	require.Equal(t, "PlayerInfoProvider", utils.GetStringPropOrEmpty(playerProps, "provider"))
	require.Equal(t, "PlayerInfoAuthId", utils.GetStringPropOrEmpty(playerProps, "authId"))
	require.Equal(t, "PlayerInfIdentityId", utils.GetStringPropOrEmpty(playerProps, "identityId"))

	identifiesRelationCount := neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "IDENTIFIES")
	require.Equal(t, 1, identifiesRelationCount)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, playerId, "IDENTIFIES", userId)
}

func TestGraphUserEventHandler_OnAddRole(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userCreateTime := utils.Now()
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName:        "UserFirstNameCreate",
		LastName:         "UserLastNameCreate",
		CreatedAt:        userCreateTime,
		UpdatedAt:        userCreateTime,
		Source:           constants.SourceOpenline,
		SourceOfTruth:    constants.SourceOpenline,
		AppSource:        constants.AppSourceEventProcessingPlatform,
		Roles:            []string{model.RoleUser.String(), model.RoleOwner.String()},
		ProfilePhotoUrl:  "www.photo.com/create",
		Timezone:         "userTimezoneCreate",
		Internal:         false,
		Bot:              false,
		DefaultForPlayer: false,
		Tenant:           tenantName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "User_" + tenantName: 1})
	dbNodeAfterUserCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, userId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterUserCreate)
	propsAfterUserCreate := utils.GetPropsFromNode(*dbNodeAfterUserCreate)
	require.Equal(t, userId, utils.GetStringPropOrEmpty(propsAfterUserCreate, "id"))
	require.Equal(t, 2, len(utils.GetListStringPropOrEmpty(propsAfterUserCreate, "roles")))

	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	roleAddedTime := utils.Now()

	roleEvent, err := user_events.NewUserAddRoleEvent(userAggregate, model.RoleCustomerOsPlatformOwner.String(), roleAddedTime)
	require.Nil(t, err)
	err = userEventHandler.OnAddRole(context.Background(), roleEvent)
	require.Nil(t, err)
	user, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "User_"+tenantName)
	require.Nil(t, err)

	propsAfterAddRole := utils.GetPropsFromNode(*user)
	require.Equal(t, 5, len(propsAfterAddRole))
	require.Equal(t, 3, len(utils.GetListStringPropOrEmpty(propsAfterAddRole, "roles")))
	require.Equal(t, "CUSTOMER_OS_PLATFORM_OWNER", utils.GetListStringPropOrEmpty(propsAfterAddRole, "roles")[2])
	require.Less(t, userCreateTime, utils.GetTimePropOrNow(propsAfterAddRole, "updatedAt"))
}

func TestGraphUserEventHandler_OnRemoveRole(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userCreateTime := utils.Now()
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, neo4jentity.UserEntity{
		FirstName:        "UserFirstNameCreate",
		LastName:         "UserLastNameCreate",
		CreatedAt:        userCreateTime,
		UpdatedAt:        userCreateTime,
		Source:           constants.SourceOpenline,
		SourceOfTruth:    constants.SourceOpenline,
		AppSource:        constants.AppSourceEventProcessingPlatform,
		Roles:            []string{model.RoleUser.String(), model.RoleOwner.String()},
		ProfilePhotoUrl:  "www.photo.com/create",
		Timezone:         "userTimezoneCreate",
		Internal:         false,
		Bot:              false,
		DefaultForPlayer: false,
		Tenant:           tenantName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"User": 1, "User_" + tenantName: 1})
	dbNodeAfterUserCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "User_"+tenantName, userId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterUserCreate)
	propsAfterUserCreate := utils.GetPropsFromNode(*dbNodeAfterUserCreate)
	require.Equal(t, userId, utils.GetStringPropOrEmpty(propsAfterUserCreate, "id"))
	require.Equal(t, 2, len(utils.GetListStringPropOrEmpty(propsAfterUserCreate, "roles")))

	userEventHandler := &UserEventHandler{
		repositories: testDatabase.Repositories,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	roleAddedTime := utils.Now()

	roleEvent, err := user_events.NewUserRemoveRoleEvent(userAggregate, model.RoleUser.String(), roleAddedTime)
	require.Nil(t, err)
	err = userEventHandler.OnRemoveRole(context.Background(), roleEvent)
	require.Nil(t, err)
	user, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "User_"+tenantName)
	require.Nil(t, err)

	propsAfterAddRole := utils.GetPropsFromNode(*user)
	require.Equal(t, 5, len(propsAfterAddRole))
	require.Equal(t, 1, len(utils.GetListStringPropOrEmpty(propsAfterAddRole, "roles")))
	require.Equal(t, "OWNER", utils.GetListStringPropOrEmpty(propsAfterAddRole, "roles")[0])
	require.Less(t, userCreateTime, utils.GetTimePropOrNow(propsAfterAddRole, "updatedAt"))
}
