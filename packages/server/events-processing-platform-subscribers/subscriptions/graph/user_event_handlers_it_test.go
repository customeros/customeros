package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test/neo4j"
	job_role_aggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/aggregate"
	job_role_model "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	job_role_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	user_aggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	user_events "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	user_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	RoleAdmin                   string = "ADMIN"
	RoleCustomerOsPlatformOwner string = "CUSTOMER_OS_PLATFORM_OWNER"
	RoleOwner                   string = "OWNER"
	RoleUser                    string = "USER"
)

func TestGraphUserEventHandler_OnUserCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &UserEventHandler{
		services: testDatabase.Services,
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
	},
		cmnmod.Source{
			Source:        "N/A",
			SourceOfTruth: "N/A",
			AppSource:     "event-processing-platform",
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
	require.Equal(t, "openline", utils.GetStringPropOrEmpty(props, "source"))
	require.Equal(t, "openline", utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	require.Equal(t, "event-processing-platform", utils.GetStringPropOrEmpty(props, "appSource"))
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
		services: testDatabase.Services,
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
		AppSource:     "event-processing-platform",
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
	require.Equal(t, "openline", utils.GetStringPropOrEmpty(props, "source"))
	require.Equal(t, "openline", utils.GetStringPropOrEmpty(props, "sourceOfTruth"))
	require.Equal(t, "event-processing-platform", utils.GetStringPropOrEmpty(props, "appSource"))
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
		services: testDatabase.Services,
	}
	jobRoleEventHandler := &JobRoleEventHandler{
		services: testDatabase.Services,
	}
	myUserId, _ := uuid.NewUUID()
	myJobRoleId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	jobRoleAggregate := job_role_aggregate.NewJobRoleAggregateWithTenantAndID(tenantName, myJobRoleId.String())

	curTime := utils.Now()

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
			AppSource:     "event-processing-platform",
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
			false, "N/A", "N/A", "event-processing-platform", &now, nil, &curTime))

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
	require.Equal(t, "openline", utils.GetStringPropOrEmpty(userProps, "source"))
	require.Equal(t, "openline", utils.GetStringPropOrEmpty(userProps, "sourceOfTruth"))
	require.Equal(t, "event-processing-platform", utils.GetStringPropOrEmpty(userProps, "appSource"))
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
	require.Equal(t, "event-processing-platform", utils.GetStringPropOrEmpty(jobRoleProps, "appSource"))
}

func TestGraphUserEventHandler_OnUserCreateWithJobRoleOutOfOrder(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	userEventHandler := &UserEventHandler{
		services: testDatabase.Services,
	}
	jobRoleEventHandler := &JobRoleEventHandler{
		services: testDatabase.Services,
	}
	myUserId, _ := uuid.NewUUID()
	myJobRoleId, _ := uuid.NewUUID()
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, myUserId.String())
	jobRoleAggregate := job_role_aggregate.NewJobRoleAggregateWithTenantAndID(tenantName, myJobRoleId.String())

	curTime := utils.Now()

	description := "I clean things"

	userCreateEvent, err := user_events.NewUserCreateEvent(userAggregate, user_models.UserDataFields{
		FirstName: "Bob",
		LastName:  "Dole",
		Name:      "Bob Dole",
	}, cmnmod.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "event-processing-platform",
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
			false, "N/A", "N/A", "event-processing-platform", nil, nil, &curTime))

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
	require.Equal(t, "openline", utils.GetStringPropOrEmpty(userProps, "source"))
	require.Equal(t, "openline", utils.GetStringPropOrEmpty(userProps, "sourceOfTruth"))
	require.Equal(t, "event-processing-platform", utils.GetStringPropOrEmpty(userProps, "appSource"))

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
	require.Equal(t, "event-processing-platform", utils.GetStringPropOrEmpty(jobRoleProps, "appSource"))
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
		AppSource:        constants.AppSourceEventProcessingPlatformSubscribers,
		Roles:            []string{RoleUser, RoleOwner},
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
		services: testDatabase.Services,
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
	require.Equal(t, 15, len(userProps))

	userId = utils.GetStringPropOrEmpty(userProps, "id")
	require.NotNil(t, userId)
	require.Equal(t, userNameUpdate, utils.GetStringPropOrEmpty(userProps, "name"))
	require.Equal(t, userFirstNameUpdate, utils.GetStringPropOrEmpty(userProps, "firstName"))
	require.Equal(t, userLastNameUpdate, utils.GetStringPropOrEmpty(userProps, "lastName"))
	require.Equal(t, userTimezoneUpdate, utils.GetStringPropOrEmpty(userProps, "timezone"))
	require.Equal(t, userProfilePhotoUrlUpdate, utils.GetStringPropOrEmpty(userProps, "profilePhotoUrl"))
	require.Less(t, userCreateTime, utils.GetTimePropOrNow(userProps, "updatedAt"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(userProps, "sourceOfTruth"))
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
		AppSource:        constants.AppSourceEventProcessingPlatformSubscribers,
		Roles:            []string{RoleUser, RoleOwner},
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

	phoneNumberId := neo4jtest.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, neo4jentity.PhoneNumberEntity{
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
		services: testDatabase.Services,
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
		AppSource:        constants.AppSourceEventProcessingPlatformSubscribers,
		Roles:            []string{RoleUser, RoleOwner},
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

	jobRoleId := neo4jt.CreateJobRole(ctx, testDatabase.Driver, tenantName, neo4jentity.JobRoleEntity{})

	dbNodeAfterjobRoleCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "JobRole_"+tenantName, jobRoleId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterjobRoleCreate)
	propsAfterJobRoleCreate := utils.GetPropsFromNode(*dbNodeAfterjobRoleCreate)
	require.NotNil(t, utils.GetStringPropOrEmpty(propsAfterJobRoleCreate, "id"))

	userEventHandler := &UserEventHandler{
		services: testDatabase.Services,
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
		AppSource:        constants.AppSourceEventProcessingPlatformSubscribers,
		Roles:            []string{RoleUser, RoleOwner},
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
		services: testDatabase.Services,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	roleAddedTime := utils.Now()

	roleEvent, err := user_events.NewUserAddRoleEvent(userAggregate, RoleCustomerOsPlatformOwner, roleAddedTime)
	require.Nil(t, err)
	err = userEventHandler.OnAddRole(context.Background(), roleEvent)
	require.Nil(t, err)
	user, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "User_"+tenantName)
	require.Nil(t, err)

	propsAfterAddRole := utils.GetPropsFromNode(*user)
	require.Equal(t, 13, len(propsAfterAddRole))
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
		AppSource:        constants.AppSourceEventProcessingPlatformSubscribers,
		Roles:            []string{RoleUser, RoleOwner},
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
		services: testDatabase.Services,
	}
	userAggregate := user_aggregate.NewUserAggregateWithTenantAndID(tenantName, userId)
	roleAddedTime := utils.Now()

	roleEvent, err := user_events.NewUserRemoveRoleEvent(userAggregate, RoleUser, roleAddedTime)
	require.Nil(t, err)
	err = userEventHandler.OnRemoveRole(context.Background(), roleEvent)
	require.Nil(t, err)
	user, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "User_"+tenantName)
	require.Nil(t, err)

	propsAfterAddRole := utils.GetPropsFromNode(*user)
	require.Equal(t, 13, len(propsAfterAddRole))
	require.Equal(t, 1, len(utils.GetListStringPropOrEmpty(propsAfterAddRole, "roles")))
	require.Equal(t, "OWNER", utils.GetListStringPropOrEmpty(propsAfterAddRole, "roles")[0])
	require.Less(t, userCreateTime, utils.GetTimePropOrNow(propsAfterAddRole, "updatedAt"))
}
