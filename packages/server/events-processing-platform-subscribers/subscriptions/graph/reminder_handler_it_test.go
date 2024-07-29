package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder"
	reminderevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReminderEventHandler_OnCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{Name: "ORG"})
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{Name: "USER"})
	reminderEvtHdlr := &ReminderEventHandler{
		repositories: testDatabase.Repositories,
		log:          testLogger,
	}
	reminderId := uuid.New().String()
	reminderAgg := reminder.NewReminderAggregateWithTenantAndID(tenantName, reminderId)

	createdAt := utils.Now()
	dueDate := createdAt.AddDate(0, 0, 1)

	e := reminderevent.ReminderCreateEvent{
		BaseEvent: event.BaseEvent{
			Tenant:     tenantName,
			EventName:  reminderevent.ReminderCreateV1,
			CreatedAt:  createdAt,
			AppSource:  "test",
			Source:     constants.SourceOpenline,
			EntityType: model.REMINDER,
		},
		Content:        "content",
		DueDate:        dueDate,
		UserId:         userId,
		OrganizationId: orgId,
		Dismissed:      false,
	}

	evt := eventstore.NewBaseEvent(reminderAgg, e.EventName)

	err := evt.SetJsonData(&e)
	require.Nil(t, err)

	err = reminderEvtHdlr.OnCreate(ctx, evt)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, model.NodeLabelReminder))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, model.NodeLabelReminder+"_"+tenantName), "Incorrect number of Reminder_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "REMINDER_BELONGS_TO_TENANT"), "Incorrect number of REMINDER_BELONGS_TO_TENANT relationships in Neo4j")
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "REMINDER_BELONGS_TO_ORGANIZATION"), "Incorrect number of REMINDER_BELONGS_TO_ORGANIZATION relationships in Neo4j")
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "REMINDER_BELONGS_TO_USER"), "Incorrect number of REMINDER_BELONGS_TO_USER relationships in Neo4j")

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Reminder_"+tenantName, reminderId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, reminderId, utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, "content", utils.GetStringPropOrEmpty(props, "content"))
	require.Equal(t, "test", utils.GetStringPropOrEmpty(props, "appSource"))
}

func TestReminderEventHandler_OnUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, entity.OrganizationEntity{Name: "ORG"})
	userId := neo4jtest.CreateUser(ctx, testDatabase.Driver, tenantName, entity.UserEntity{Name: "USER"})
	reminderEvtHdlr := &ReminderEventHandler{
		repositories: testDatabase.Repositories,
		log:          testLogger,
	}
	reminderId := uuid.New().String()
	reminderAgg := reminder.NewReminderAggregateWithTenantAndID(tenantName, reminderId)

	createdAt := utils.Now()
	dueDate := createdAt.AddDate(0, 0, 1)

	e := reminderevent.ReminderCreateEvent{
		BaseEvent: event.BaseEvent{
			Tenant:     tenantName,
			EventName:  reminderevent.ReminderCreateV1,
			CreatedAt:  createdAt,
			AppSource:  "test",
			Source:     constants.SourceOpenline,
			EntityType: model.REMINDER,
		},
		Content:        "content",
		DueDate:        dueDate,
		UserId:         userId,
		OrganizationId: orgId,
		Dismissed:      false,
	}

	evt1 := eventstore.NewBaseEvent(reminderAgg, e.EventName)

	err := evt1.SetJsonData(&e)
	require.Nil(t, err)

	err = reminderEvtHdlr.OnCreate(ctx, evt1)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, model.NodeLabelReminder))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, model.NodeLabelReminder+"_"+tenantName), "Incorrect number of Reminder_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "REMINDER_BELONGS_TO_TENANT"), "Incorrect number of REMINDER_BELONGS_TO_TENANT relationships in Neo4j")
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "REMINDER_BELONGS_TO_ORGANIZATION"), "Incorrect number of REMINDER_BELONGS_TO_ORGANIZATION relationships in Neo4j")
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "REMINDER_BELONGS_TO_USER"), "Incorrect number of REMINDER_BELONGS_TO_USER relationships in neo4j")

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Reminder_"+tenantName, reminderId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, reminderId, utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, "content", utils.GetStringPropOrEmpty(props, "content"))
	require.Equal(t, "test", utils.GetStringPropOrEmpty(props, "appSource"))

	e2 := reminderevent.ReminderUpdateEvent{
		BaseEvent: event.BaseEvent{
			Tenant:     tenantName,
			EventName:  reminderevent.ReminderUpdateV1,
			CreatedAt:  createdAt,
			AppSource:  "test",
			Source:     constants.SourceOpenline,
			EntityId:   reminderId,
			EntityType: model.REMINDER,
		},
		Content:   "NEW_CONTENT",
		DueDate:   dueDate.AddDate(0, 0, 1),
		Dismissed: true,
	}

	evt2 := eventstore.NewBaseEvent(reminderAgg, e2.EventName)

	err = evt2.SetJsonData(&e2)
	require.Nil(t, err)

	err = reminderEvtHdlr.OnUpdate(ctx, evt2)
	require.Nil(t, err)

	dbNode, err = neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Reminder_"+tenantName, reminderId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props = utils.GetPropsFromNode(*dbNode)

	require.Equal(t, reminderId, utils.GetStringPropOrEmpty(props, "id"))
	require.Equal(t, "NEW_CONTENT", utils.GetStringPropOrEmpty(props, "content"))
	require.Equal(t, "test", utils.GetStringPropOrEmpty(props, "appSource"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(props, "dismissed"))
	require.Equal(t, dueDate.AddDate(0, 0, 1).Format("2006-01-02"), utils.GetTimePropOrZeroTime(props, "dueDate").Format("2006-01-02"))
}
