package test

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func CleanupAllData(ctx context.Context, driver *neo4j.DriverWithContext) {
	ExecuteWriteQuery(ctx, driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func CreateTenant(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant}) ON CREATE SET t.createdAt=$now`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
		"now":    utils.Now(),
	})
}

func CreateUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, user entity.UserEntity) string {
	userId := utils.NewUUIDIfEmpty(user.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
			MERGE (u:User {id: $userId})-[:USER_BELONGS_TO_TENANT]->(t)
			SET u:User_%s, 
				u.roles=$roles,
				u.internal=$internal,
				u.bot=$bot,
				u.firstName=$firstName,
				u.lastName=$lastName,
				u.profilePhotoUrl=$profilePhotoUrl,
				u.createdAt=$createdAt,
				u.updatedAt=$updatedAt,
				u.source=$source,
				u.sourceOfTruth=$sourceOfTruth,
				u.appSource=$appSource`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":          tenant,
		"userId":          userId,
		"firstName":       user.FirstName,
		"lastName":        user.LastName,
		"source":          user.Source,
		"sourceOfTruth":   user.SourceOfTruth,
		"appSource":       user.AppSource,
		"roles":           user.Roles,
		"internal":        user.Internal,
		"bot":             user.Bot,
		"profilePhotoUrl": user.ProfilePhotoUrl,
		"createdAt":       user.CreatedAt,
		"updatedAt":       user.UpdatedAt,
	})
	return userId
}

func CreateUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string) {
	CreateUser(ctx, driver, tenant, entity.UserEntity{
		Id: userId,
	})
}

func CreateInvoicingCycle(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, entity entity.InvoicingCycleEntity) string {
	id := utils.NewUUIDIfEmpty(entity.Id)

	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:INVOICING_CYCLE_BELONGS_TO_TENANT]-(ic:InvoicingCycle {id:$id}) 
				SET ic:InvoicingCycle_%s,
					ic.type=$type,
					ic.createdAt=$createdAt,
					ic.source=$source,
					ic.sourceOfTruth=$sourceOfTruth,
					ic.appSource=$appSource
					`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"id":            id,
		"type":          entity.Type,
		"createdAt":     entity.CreatedAt,
		"source":        entity.Source,
		"sourceOfTruth": entity.SourceOfTruth,
		"appSource":     entity.AppSource,
	})
	return id
}

func CreateMasterPlan(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, masterPlan entity.MasterPlanEntity) string {
	masterPlanId := utils.NewUUIDIfEmpty(masterPlan.Id)

	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan {id:$id})
				SET mp:MasterPlan_%s,
					mp.name=$name,
					mp.createdAt=$createdAt,
					mp.source=$source,
					mp.sourceOfTruth=$sourceOfTruth,
					mp.appSource=$appSource,
					mp.retired=$retired
					`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"id":            masterPlanId,
		"name":          masterPlan.Name,
		"createdAt":     masterPlan.CreatedAt,
		"source":        masterPlan.Source,
		"sourceOfTruth": masterPlan.SourceOfTruth,
		"appSource":     masterPlan.AppSource,
		"retired":       masterPlan.Retired,
	})
	return masterPlanId
}

func CreateMasterPlanMilestone(ctx context.Context, driver *neo4j.DriverWithContext, tenant, masterPlanId string, masterPlanMilestone entity.MasterPlanMilestoneEntity) string {
	masterPlanMilestoneId := utils.NewUUIDIfEmpty(masterPlanMilestone.Id)

	query := fmt.Sprintf(`MATCH (mp:MasterPlan {id: $masterPlanId})
			  MERGE (mp)-[:HAS_MILESTONE]->(m:MasterPlanMilestone {id:$id})
				SET m:MasterPlanMilestone_%s,
					m.name=$name,
					m.createdAt=$createdAt,
					m.order=$order,
					m.durationHours=$durationHours,
					m.optional=$optional,
					m.items=$items,
					m.retired=$retired`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"masterPlanId":  masterPlanId,
		"id":            masterPlanMilestoneId,
		"name":          masterPlanMilestone.Name,
		"createdAt":     masterPlanMilestone.CreatedAt,
		"order":         masterPlanMilestone.Order,
		"durationHours": masterPlanMilestone.DurationHours,
		"optional":      masterPlanMilestone.Optional,
		"items":         masterPlanMilestone.Items,
		"retired":       masterPlanMilestone.Retired,
	})
	return masterPlanMilestoneId
}

func CreateOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, organization entity.OrganizationEntity) string {
	orgId := utils.NewUUIDIfEmpty(organization.ID)
	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$id})
				SET o:Organization_%s,
					o.name=$name,
					o.hide=$hide
				`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
		"name":   organization.Name,
		"hide":   organization.Hide,
		"id":     orgId,
	})
	return orgId
}

func CreateLogEntry(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, logEntry entity.LogEntryEntity) string {
	logEntryId := utils.NewUUIDIfEmpty(logEntry.Id)
	query := fmt.Sprintf(`
			  MERGE (l:LogEntry {id:$id})
				SET l:LogEntry_%s,
					l:TimelineEvent,
					l:TimelineEvent_%s,
					l.content=$content,
					l.contentType=$contentType,
					l.startedAt=$startedAt
				`, tenant, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":      tenant,
		"id":          logEntryId,
		"content":     logEntry.Content,
		"contentType": logEntry.ContentType,
		"startedAt":   logEntry.StartedAt,
	})
	return logEntryId
}

func CreateBillingProfileForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, orgId string, billingProfile entity.BillingProfileEntity) string {
	billingProfileId := CreateBillingProfile(ctx, driver, tenant, billingProfile)
	LinkNodes(ctx, driver, orgId, billingProfileId, "HAS_BILLING_PROFILE")
	return billingProfileId
}

func CreateBillingProfile(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, billingProfile entity.BillingProfileEntity) string {
	billingProfileId := utils.NewUUIDIfEmpty(billingProfile.Id)
	query := fmt.Sprintf(`
			  MERGE (bp:BillingProfile {id:$id})
				SET bp:BillingProfile_%s
				`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
		"id":     billingProfileId,
	})
	return billingProfileId
}

func CreateLogEntryForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, orgId string, logEntry entity.LogEntryEntity) string {
	logEntryId := CreateLogEntry(ctx, driver, tenant, logEntry)
	LinkNodes(ctx, driver, orgId, logEntryId, "LOGGED")
	return logEntryId
}

func LinkNodes(ctx context.Context, driver *neo4j.DriverWithContext, fromId, toId string, relation string, properties ...map[string]any) {
	query := fmt.Sprintf(`
			  MATCH (from {id: $fromId})
			  MATCH (to {id: $toId})
			  MERGE (from)-[rel:%s]->(to)`, relation)
	params := map[string]any{
		"fromId": fromId,
		"toId":   toId,
	}
	if len(properties) > 0 {
		params["properties"] = properties[0]
		query += " SET rel += $properties"
	}

	ExecuteWriteQuery(ctx, driver, query, params)
}

func CreateEmail(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, entity entity.EmailEntity) string {
	emailId := utils.NewUUIDIfEmpty(entity.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
								MERGE (e:Email {id:$emailId})
								MERGE (e)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
								ON CREATE SET e:Email_%s,
									e.email=$email,
									e.rawEmail=$rawEmail,
									e.isReachable=$isReachable,
									e.createdAt=$createdAt,
									e.updatedAt=$updatedAt
							`, tenant)
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":      tenant,
		"emailId":     emailId,
		"email":       entity.Email,
		"rawEmail":    entity.RawEmail,
		"isReachable": entity.IsReachable,
		"createdAt":   entity.CreatedAt,
		"updatedAt":   entity.UpdatedAt,
	})
	return emailId
}
