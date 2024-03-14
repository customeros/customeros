package resolver

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/stretchr/testify/require"
)

func TestQueryResolver_Reminder(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	now := utils.Now()
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "TEST ORG"})
	uid := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{Name: "TEST USER"})
	rid := neo4jtest.CreateReminder(ctx, driver, tenantName, uid, orgId, now, neo4jentity.ReminderEntity{
		Content:       "TEST CONTENT",
		DueDate:       utils.Now(),
		Dismissed:     false,
		Source:        neo4jentity.DataSourceOpenline,
		AppSource:     "TEST APP SOURCE",
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	})
	neo4jtest.LinkReminderToUser(ctx, driver, tenantName, rid, uid)
	neo4jtest.LinkReminderToOrganization(ctx, driver, tenantName, rid, orgId)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))

	rawResponse := callGraphQL(t, "reminder/get_reminder_by_id", map[string]interface{}{"reminderId": rid})

	var reminderStruct struct {
		Reminder model.Reminder
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
	require.Nil(t, err)
	require.NotNil(t, reminderStruct.Reminder)
	require.Equal(t, rid, reminderStruct.Reminder.Metadata.ID)
	require.Equal(t, "TEST CONTENT", reminderStruct.Reminder.Content)
	require.Equal(t, now.Format("2006-01-02"), reminderStruct.Reminder.DueDate.Format("2006-01-02"))
	require.Equal(t, false, reminderStruct.Reminder.Dismissed)
	require.Equal(t, uid, reminderStruct.Reminder.Owner.ID)
}

func TestQueryResolver_RemindersForOrg(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	now := utils.Now()
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "TEST ORG"})
	uid := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{Name: "TEST USER", FirstName: "TEST", LastName: "USER"})
	rid := neo4jtest.CreateReminder(ctx, driver, tenantName, uid, orgId, now, neo4jentity.ReminderEntity{
		Content:       "TEST CONTENT",
		DueDate:       now,
		Dismissed:     false,
		Source:        neo4jentity.DataSourceOpenline,
		AppSource:     "TEST APP SOURCE",
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	})
	neo4jtest.LinkReminderToUser(ctx, driver, tenantName, rid, uid)
	neo4jtest.LinkReminderToOrganization(ctx, driver, tenantName, rid, orgId)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))

	rawResponse := callGraphQL(t, "reminder/get_reminders_for_org", map[string]interface{}{"organizationId": orgId})
	var reminderStruct struct {
		RemindersForOrganization []*model.Reminder
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
	require.Nil(t, err)
	require.NotNil(t, reminderStruct.RemindersForOrganization)
	require.Equal(t, 1, len(reminderStruct.RemindersForOrganization))
	require.Equal(t, rid, reminderStruct.RemindersForOrganization[0].Metadata.ID)
	require.Equal(t, "TEST CONTENT", reminderStruct.RemindersForOrganization[0].Content)
	require.Equal(t, now.Format("2006-01-02"), reminderStruct.RemindersForOrganization[0].DueDate.Format("2006-01-02"))
	require.Equal(t, false, reminderStruct.RemindersForOrganization[0].Dismissed)
	require.Equal(t, uid, reminderStruct.RemindersForOrganization[0].Owner.ID)
	require.Equal(t, "TEST", reminderStruct.RemindersForOrganization[0].Owner.FirstName)

}

func TestMutationResolver_ReminderCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "TEST ORG"})
	uid := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: "TEST", LastName: "USER"})
	rid := uuid.New().String()
	dueDate := utils.Now()
	neo4jtest.CreateReminder(ctx, driver, tenantName, uid, orgId, utils.Now(), neo4jentity.ReminderEntity{
		Id:            rid,
		Content:       "TEST CONTENT",
		DueDate:       dueDate,
		Dismissed:     false,
		Source:        neo4jentity.DataSourceOpenline,
		AppSource:     "TEST APP SOURCE",
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	})
	neo4jtest.LinkReminderToUser(ctx, driver, tenantName, rid, uid)
	neo4jtest.LinkReminderToOrganization(ctx, driver, tenantName, rid, orgId)

	calledCreateReminder := false
	reminderServiceCallbacks := events_platform.MockReminderServiceCallbacks{
		ReminderCreate: func(context context.Context, request *reminderpb.CreateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error) {
			require.Equal(t, tenantName, request.Tenant)
			require.Equal(t, uid, request.UserId)
			require.Equal(t, orgId, request.OrganizationId)
			require.Equal(t, "TEST CONTENT", request.Content)
			reqDate := utils.TimestampProtoToTime(request.DueDate)
			require.Equal(t, dueDate, reqDate)

			calledCreateReminder = true

			return &reminderpb.ReminderGrpcResponse{
				Id: rid,
			}, nil
		},
	}
	events_platform.SetReminderCallbacks(&reminderServiceCallbacks)

	rawResponse := callGraphQL(t, "reminder/create_reminder", map[string]interface{}{
		"organizationId": orgId,
		"userId":         uid,
		"content":        "TEST CONTENT",
		"dueDate":        dueDate,
	})

	var reminderStruct struct {
		Reminder_Create model.Reminder
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
	require.Nil(t, err)
	require.True(t, calledCreateReminder)
	require.NotNil(t, reminderStruct.Reminder_Create)
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))
	require.Equal(t, "TEST CONTENT", reminderStruct.Reminder_Create.Content)
	require.Equal(t, uid, reminderStruct.Reminder_Create.Owner.ID)
	require.Equal(t, "TEST", reminderStruct.Reminder_Create.Owner.FirstName)
}

func TestMutationResolver_ReminderUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	now := utils.Now()
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "TEST ORG"})
	uid := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: "TEST", LastName: "USER"})
	rid := neo4jtest.CreateReminder(ctx, driver, tenantName, uid, orgId, now, neo4jentity.ReminderEntity{
		Content:       "TEST CONTENT",
		DueDate:       now,
		Dismissed:     false,
		Source:        neo4jentity.DataSourceOpenline,
		AppSource:     "TEST APP SOURCE",
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	})
	neo4jtest.LinkReminderToUser(ctx, driver, tenantName, rid, uid)
	neo4jtest.LinkReminderToOrganization(ctx, driver, tenantName, rid, orgId)

	dueDate := utils.Now()

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))

	calledUpdateReminder := false
	reminderServiceCallbacks := events_platform.MockReminderServiceCallbacks{
		ReminderUpdate: func(context context.Context, request *reminderpb.UpdateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error) {
			require.Equal(t, tenantName, request.Tenant)
			require.Equal(t, "UPDATED CONTENT", request.Content)
			reqDate := utils.TimestampProtoToTime(request.DueDate)
			require.Equal(t, dueDate, reqDate)
			require.Equal(t, true, request.Dismissed)

			// helper update the reminder in neo4j
			neo4jtest.UpdateReminder(ctx, driver, tenantName, rid, neo4jentity.ReminderEntity{
				Content:   "UPDATED CONTENT",
				DueDate:   dueDate,
				Dismissed: true,
				UpdatedAt: utils.Now(),
			})

			calledUpdateReminder = true

			return &reminderpb.ReminderGrpcResponse{
				Id: rid,
			}, nil
		},
	}
	events_platform.SetReminderCallbacks(&reminderServiceCallbacks)

	rawResponse := callGraphQL(t, "reminder/update_reminder", map[string]interface{}{
		"id":        rid,
		"content":   "UPDATED CONTENT",
		"dueDate":   dueDate,
		"dismissed": true,
	})

	var reminderStruct struct {
		Reminder_Update model.Reminder
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
	require.Nil(t, err)
	require.NotNil(t, reminderStruct.Reminder_Update)
	require.True(t, calledUpdateReminder)
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))
	require.Equal(t, "UPDATED CONTENT", reminderStruct.Reminder_Update.Content)
	require.Equal(t, true, reminderStruct.Reminder_Update.Dismissed)
	require.Equal(t, uid, reminderStruct.Reminder_Update.Owner.ID)
	require.Equal(t, "TEST", reminderStruct.Reminder_Update.Owner.FirstName)
}
