package resolver

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-api/utils/decode"
)

func TestQueryResolver_Reminder(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})
	userId := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{})
	reminderId := neo4jtest.CreateReminder(ctx, driver, tenantName, userId, organizationId, now, neo4jentity.ReminderEntity{
		Content:   "TEST CONTENT",
		DueDate:   utils.Now(),
		Dismissed: false,
	})

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))

	rawResponse := callGraphQL(t, "reminder/get_reminder_by_id", map[string]interface{}{"reminderId": reminderId})

	var reminderStruct struct {
		Reminder model.Reminder
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
	require.Nil(t, err)
	require.NotNil(t, reminderStruct.Reminder)
	require.Equal(t, reminderId, reminderStruct.Reminder.Metadata.ID)
	require.Equal(t, "TEST CONTENT", *reminderStruct.Reminder.Content)
	require.Equal(t, now.Format("2006-01-02"), reminderStruct.Reminder.DueDate.Format("2006-01-02"))
	require.Equal(t, false, *reminderStruct.Reminder.Dismissed)
	require.Equal(t, userId, reminderStruct.Reminder.Owner.ID)
}

func TestQueryResolver_RemindersForOrg(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	now := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "TEST ORG"})
	userId := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{Name: "TEST USER", FirstName: "TEST", LastName: "USER"})
	reminderId := neo4jtest.CreateReminder(ctx, driver, tenantName, userId, organizationId, now, neo4jentity.ReminderEntity{
		Content:   "TEST CONTENT",
		DueDate:   now,
		Dismissed: false,
	})

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))

	// TODO - move this in organization
	rawResponse := callGraphQL(t, "reminder/get_reminders_for_org", map[string]interface{}{"organizationId": organizationId})
	var reminderStruct struct {
		RemindersForOrganization []*model.Reminder
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
	require.Nil(t, err)
	require.NotNil(t, reminderStruct.RemindersForOrganization)
	require.Equal(t, 1, len(reminderStruct.RemindersForOrganization))
	require.Equal(t, reminderId, reminderStruct.RemindersForOrganization[0].Metadata.ID)
	require.Equal(t, "TEST CONTENT", *reminderStruct.RemindersForOrganization[0].Content)
	require.Equal(t, now.Format("2006-01-02"), reminderStruct.RemindersForOrganization[0].DueDate.Format("2006-01-02"))
	require.Equal(t, false, *reminderStruct.RemindersForOrganization[0].Dismissed)
	require.Equal(t, userId, reminderStruct.RemindersForOrganization[0].Owner.ID)
	require.Equal(t, "TEST", reminderStruct.RemindersForOrganization[0].Owner.FirstName)
}

func TestMutationResolver_ReminderCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	dueDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "TEST ORG"})
	userId := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{FirstName: "TEST", LastName: "USER"})

	rawResponse := callGraphQL(t, "reminder/create_reminder", map[string]interface{}{
		"organizationId": organizationId,
		"userId":         userId,
		"content":        "TEST CONTENT",
		"dueDate":        dueDate,
	})

	var reminderStruct struct {
		Reminder_Create *string
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
	require.Nil(t, err)
	require.NotNil(t, reminderStruct.Reminder_Create)
}

func TestMutationResolver_ReminderUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	dueDate := utils.Now()

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	reminderId := uuid.New().String()

	rawResponse := callGraphQL(t, "reminder/update_reminder", map[string]interface{}{
		"id":        reminderId,
		"content":   "UPDATED CONTENT",
		"dueDate":   dueDate,
		"dismissed": true,
	})

	var reminderStruct struct {
		Reminder_Update *string
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
	require.Nil(t, err)
	require.NotNil(t, reminderStruct.Reminder_Update)
	require.Equal(t, reminderId, *reminderStruct.Reminder_Update)
}
