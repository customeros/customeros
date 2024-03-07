package resolver

import (
	"context"
	"testing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
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
	uid := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{Name: "TEST USER"})
	rid := neo4jtest.CreateReminder(ctx, driver, tenantName, uid, orgId, now, neo4jentity.ReminderEntity{
		Content:       "TEST CONTENT",
		DueDate:       now,
		Dismissed:     false,
		Source:        neo4jentity.DataSourceOpenline,
		AppSource:     "TEST APP SOURCE",
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	})

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
}

// func TestMutationResolver_ReminderCreate(t *testing.T) {
// 	ctx := context.Background()
// 	defer tearDownTestCase(ctx)(t)
// 	neo4jtest.CreateTenant(ctx, driver, tenantName)
// 	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "TEST ORG"})
// 	uid := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{Name: "TEST USER"})

// 	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))

// 	rawResponse := callGraphQL(t, "reminder/create_reminder", map[string]interface{}{
// 		"input": map[string]interface{}{
// 			"userId":  uid,
// 			"orgId":   orgId,
// 			"content": "TEST CONTENT",
// 			"dueDate": utils.Now().Format("2006-01-02"),
// 		},
// 	})

// 	var reminderStruct struct {
// 		Reminder model.Reminder
// 	}
// 	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
// 	require.Nil(t, err)
// 	require.NotNil(t, reminderStruct.Reminder)
// 	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))
// 	require.Equal(t, "TEST CONTENT", reminderStruct.Reminder.Content)
// 	require.Equal(t, uid, reminderStruct.Reminder.Owner.ID)
// }

// func TestMutationResolver_ReminderUpdate(t *testing.T) {
// 	ctx := context.Background()
// 	defer tearDownTestCase(ctx)(t)
// 	now := utils.Now()
// 	neo4jtest.CreateTenant(ctx, driver, tenantName)
// 	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "TEST ORG"})
// 	uid := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{Name: "TEST USER"})
// 	rid := neo4jtest.CreateReminder(ctx, driver, tenantName, uid, orgId, now, neo4jentity.ReminderEntity{
// 		Content:       "TEST CONTENT",
// 		DueDate:       now,
// 		Dismissed:     false,
// 		Source:        neo4jentity.DataSourceOpenline,
// 		AppSource:     "TEST APP SOURCE",
// 		SourceOfTruth: neo4jentity.DataSourceOpenline,
// 	})

// 	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))

// 	rawResponse := callGraphQL(t, "reminder/update_reminder", map[string]interface{}{
// 		"input": map[string]interface{}{
// 			"id":        rid,
// 			"content":   "UPDATED CONTENT",
// 			"dueDate":   now.Format("2006-01-02"),
// 			"dismissed": true,
// 		},
// 	})

// 	var reminderStruct struct {
// 		Reminder model.Reminder
// 	}
// 	err := decode.Decode(rawResponse.Data.(map[string]any), &reminderStruct)
// 	require.Nil(t, err)
// 	require.NotNil(t, reminderStruct.Reminder)
// 	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Reminder"))
// 	require.Equal(t, "UPDATED CONTENT", reminderStruct.Reminder.Content)
// 	require.Equal(t, true, reminderStruct.Reminder.Dismissed)
// 	require.Equal(t, uid, reminderStruct.Reminder.Owner.ID)
// }
