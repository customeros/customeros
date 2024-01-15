package test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
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

func CreateWorkspace(ctx context.Context, driver *neo4j.DriverWithContext, workspace string, provider string, tenant string) {
	query := `MATCH (t:Tenant {name: $tenant})
			  MERGE (t)-[:HAS_WORKSPACE]->(w:Workspace {name:$workspace, provider:$provider})`

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":    tenant,
		"provider":  provider,
		"workspace": workspace,
	})
}

func CreateDefaultUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateUser(ctx, driver, tenant, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateDefaultUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string) string {
	return CreateUser(ctx, driver, tenant, entity.UserEntity{
		Id:            userId,
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, user entity.UserEntity) string {
	now := utils.Now()
	createdAt := user.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := user.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = now
	}

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
		"createdAt":       createdAt,
		"updatedAt":       updatedAt,
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

func CreateInvoice(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, entity entity.InvoiceEntity) string {
	id := utils.NewUUIDIfEmpty(entity.Id)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s {id:$organizationId})
							MERGE (t)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id}) 
							ON CREATE SET 
								i:Invoice_%s,
								i.createdAt=$createdAt,
								i.updatedAt=$updatedAt,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.date=$date,
								i.dueDate=$dueDate,
								i.dryRun=$dryRun,
								i.amount=$amount,
								i.vat=$vat,
								i.total=$total,
								i.repositoryFileId=$repositoryFileId,
								i.pdfGenerated=$pdfGenerated
							WITH o, i 
							MERGE (o)-[:HAS_INVOICE]->(i) 
					`, tenant, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"organizationId":   organizationId,
		"id":               id,
		"createdAt":        entity.CreatedAt,
		"updatedAt":        entity.UpdatedAt,
		"source":           entity.Source,
		"sourceOfTruth":    entity.SourceOfTruth,
		"appSource":        entity.AppSource,
		"dryRun":           entity.DryRun,
		"date":             entity.Date,
		"dueDate":          entity.DueDate,
		"amount":           entity.Amount,
		"vat":              entity.Vat,
		"total":            entity.Total,
		"repositoryFileId": entity.RepositoryFileId,
		"pdfGenerated":     entity.PdfGenerated,
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

func CreateLocation(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, location entity.LocationEntity) string {
	locationId := utils.NewUUIDIfEmpty(location.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:LOCATION_BELONGS_TO_TENANT]-(i:Location {id:$id})
				SET i:Location_%s,
					i.name=$name,
					i.createdAt=$createdAt,
					i.updatedAt=$updatedAt,
					i.country=$country,
					i.region=$region,    
					i.locality=$locality,    
					i.address=$address,    
					i.address2=$address2,    
					i.zip=$zip,    
					i.addressType=$addressType,    
					i.houseNumber=$houseNumber,    
					i.postalCode=$postalCode,    
					i.plusFour=$plusFour,    
					i.commercial=$commercial,    
					i.predirection=$predirection,    
					i.district=$district,    
					i.street=$street,    
					i.rawAddress=$rawAddress,    
					i.latitude=$latitude,    
					i.longitude=$longitude,    
					i.timeZone=$timeZone,    
					i.utcOffset=$utcOffset,    
					i.sourceOfTruth=$sourceOfTruth,
					i.source=$source,
					i.appSource=$appSource`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"id":            locationId,
		"name":          location.Name,
		"createdAt":     location.CreatedAt,
		"updatedAt":     location.UpdatedAt,
		"country":       location.Country,
		"region":        location.Region,
		"locality":      location.Locality,
		"address":       location.Address,
		"address2":      location.Address2,
		"zip":           location.Zip,
		"addressType":   location.AddressType,
		"houseNumber":   location.HouseNumber,
		"postalCode":    location.PostalCode,
		"plusFour":      location.PlusFour,
		"commercial":    location.Commercial,
		"predirection":  location.Predirection,
		"district":      location.District,
		"street":        location.Street,
		"rawAddress":    location.RawAddress,
		"latitude":      location.Latitude,
		"longitude":     location.Longitude,
		"timeZone":      location.TimeZone,
		"utcOffset":     location.UtcOffset,
		"sourceOfTruth": location.SourceOfTruth,
		"source":        location.Source,
		"appSource":     location.AppSource,
	})
	return locationId
}

func CreateContractForOrganization(ctx context.Context, driver *neo4j.DriverWithContext, tenant, orgId string, contract entity.ContractEntity) string {
	contractId := utils.NewUUIDIfEmpty(contract.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}), (o:Organization {id:$orgId})
				MERGE (t)<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$id})<-[:HAS_CONTRACT]-(o)
				SET 
					c:Contract_%s,
					c.name=$name,
					c.contractUrl=$contractUrl,
					c.source=$source,
					c.sourceOfTruth=$sourceOfTruth,
					c.appSource=$appSource,
					c.status=$status,
					c.renewalCycle=$renewalCycle,
					c.renewalPeriods=$renewalPeriods,
					c.signedAt=$signedAt,
					c.serviceStartedAt=$serviceStartedAt,
					c.endedAt=$endedAt,
					c.createdAt=$createdAt,
					c.updatedAt=$updatedAt
				`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":               contractId,
		"orgId":            orgId,
		"tenant":           tenant,
		"name":             contract.Name,
		"contractUrl":      contract.ContractUrl,
		"source":           contract.Source,
		"sourceOfTruth":    contract.SourceOfTruth,
		"appSource":        contract.AppSource,
		"status":           contract.ContractStatus,
		"renewalCycle":     contract.RenewalCycle,
		"renewalPeriods":   contract.RenewalPeriods,
		"signedAt":         utils.TimePtrFirstNonNilNillableAsAny(contract.SignedAt),
		"serviceStartedAt": utils.TimePtrFirstNonNilNillableAsAny(contract.ServiceStartedAt),
		"endedAt":          utils.TimePtrFirstNonNilNillableAsAny(contract.EndedAt),
		"createdAt":        contract.CreatedAt,
		"updatedAt":        contract.UpdatedAt,
	})
	return contractId
}

func CreateOpportunityForContract(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, opportunity entity.OpportunityEntity) string {
	opportunityId := utils.NewUUIDIfEmpty(opportunity.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}), (c:Contract {id:$contractId})
				MERGE (t)<-[:OPPORTUNITY_BELONGS_TO_TENANT]-(op:Opportunity {id:$id})<-[:HAS_OPPORTUNITY]-(c)
				SET 
                    op:Opportunity_%s,
					op.name=$name,
					op.source=$source,
					op.sourceOfTruth=$sourceOfTruth,
					op.appSource=$appSource,
					op.amount=$amount,
					op.maxAmount=$maxAmount,
                    op.internalType=$internalType,
					op.externalType=$externalType,
					op.internalStage=$internalStage,
					op.externalStage=$externalStage,
					op.estimatedClosedAt=$estimatedClosedAt,
					op.generalNotes=$generalNotes,
                    op.comments=$comments,
                    op.renewedAt=$renewedAt,
                    op.renewalLikelihood=$renewalLikelihood,
                    op.renewalUpdatedByUserId=$renewalUpdatedByUserId,
                    op.renewalUpdateByUserAt=$renewalUpdateByUserAt,
					op.nextSteps=$nextSteps,
					op.createdAt=$createdAt,
					op.updatedAt=$updatedAt
				`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":                     opportunityId,
		"contractId":             contractId,
		"tenant":                 tenant,
		"name":                   opportunity.Name,
		"source":                 opportunity.Source,
		"sourceOfTruth":          opportunity.SourceOfTruth,
		"appSource":              opportunity.AppSource,
		"amount":                 opportunity.Amount,
		"maxAmount":              opportunity.MaxAmount,
		"internalType":           opportunity.InternalType,
		"externalType":           opportunity.ExternalType,
		"internalStage":          opportunity.InternalStage,
		"externalStage":          opportunity.ExternalStage,
		"estimatedClosedAt":      opportunity.EstimatedClosedAt,
		"generalNotes":           opportunity.GeneralNotes,
		"nextSteps":              opportunity.NextSteps,
		"comments":               opportunity.Comments,
		"renewedAt":              opportunity.RenewedAt,
		"renewalLikelihood":      opportunity.RenewalLikelihood,
		"renewalUpdatedByUserId": opportunity.RenewalUpdatedByUserId,
		"renewalUpdateByUserAt":  opportunity.RenewalUpdatedByUserAt,
		"createdAt":              opportunity.CreatedAt,
		"updatedAt":              opportunity.UpdatedAt,
	})
	return opportunityId
}

func ActiveRenewalOpportunityForContract(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId, opportunityId string) string {
	query := fmt.Sprintf(`
				MATCH (c:Contract_%s {id:$contractId}), (op:Opportunity_%s {id:$opportunityId})
				MERGE (c)-[:ACTIVE_RENEWAL]->(op)
				`, tenant, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"opportunityId": opportunityId,
		"contractId":    contractId,
	})
	return opportunityId
}

func CreateServiceLineItemForContract(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, serviceLineItem entity.ServiceLineItemEntity) string {
	serviceLineItemId := utils.NewUUIDIfEmpty(serviceLineItem.ID)
	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}), (c:Contract {id:$contractId})
				MERGE (t)<-[:CONTRACT_BELONGS_TO_TENANT]-(sli:ServiceLineItem {id:$id})<-[:HAS_SERVICE]-(c)
				SET 
					sli:ServiceLineItem_%s,
					sli.name=$name,
					sli.source=$source,
					sli.sourceOfTruth=$sourceOfTruth,
					sli.appSource=$appSource,
					sli.isCanceled=$isCanceled,	
					sli.billed=$billed,	
					sli.quantity=$quantity,	
					sli.price=$price,
					sli.previousBilled=$previousBilled,	
					sli.previousQuantity=$previousQuantity,	
					sli.previousPrice=$previousPrice,
                    sli.comments=$comments,
					sli.startedAt=$startedAt,
					sli.endedAt=$endedAt,
					sli.createdAt=$createdAt,
					sli.updatedAt=$updatedAt,
	                sli.parentId=$parentId
				`, tenant)

	params := map[string]any{
		"id":               serviceLineItemId,
		"contractId":       contractId,
		"tenant":           tenant,
		"name":             serviceLineItem.Name,
		"source":           serviceLineItem.Source,
		"sourceOfTruth":    serviceLineItem.SourceOfTruth,
		"appSource":        serviceLineItem.AppSource,
		"isCanceled":       serviceLineItem.IsCanceled,
		"billed":           serviceLineItem.Billed,
		"quantity":         serviceLineItem.Quantity,
		"price":            serviceLineItem.Price,
		"previousBilled":   serviceLineItem.PreviousBilled,
		"previousQuantity": serviceLineItem.PreviousQuantity,
		"previousPrice":    serviceLineItem.PreviousPrice,
		"startedAt":        serviceLineItem.StartedAt,
		"comments":         serviceLineItem.Comments,
		"createdAt":        serviceLineItem.CreatedAt,
		"updatedAt":        serviceLineItem.UpdatedAt,
		"parentId":         serviceLineItem.ParentID,
	}

	if serviceLineItem.EndedAt != nil {
		params["endedAt"] = *serviceLineItem.EndedAt
	} else {
		params["endedAt"] = nil
	}

	ExecuteWriteQuery(ctx, driver, query, params)
	return serviceLineItemId
}

func InsertContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, contract entity.ContractEntity, opportunity entity.OpportunityEntity) string {
	contractId := CreateContractForOrganization(ctx, driver, tenant, organizationId, contract)
	CreateOpportunityForContract(ctx, driver, tenant, contractId, opportunity)
	return contractId
}

func InsertContractWithActiveRenewalOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, tenant, organizationId string, contract entity.ContractEntity, opportunity entity.OpportunityEntity) string {
	contractId := CreateContractForOrganization(ctx, driver, tenant, organizationId, contract)
	opportunityId := CreateOpportunityForContract(ctx, driver, tenant, contractId, opportunity)
	ActiveRenewalOpportunityForContract(ctx, driver, tenant, contractId, opportunityId)
	return contractId
}

func InsertServiceLineItem(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, startedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		StartedAt: startedAt,
	})
	return id
}

func InsertServiceLineItemEnded(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		StartedAt: startedAt,
		EndedAt:   &endedAt,
	})
	return id
}

func InsertServiceLineItemCanceled(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:         id,
		ParentID:   id,
		Billed:     billedType,
		Price:      price,
		Quantity:   quantity,
		IsCanceled: true,
		StartedAt:  startedAt,
		EndedAt:    &endedAt,
	})
	return id
}

func InsertServiceLineItemWithParent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, previousBilledType enum.BilledType, previousPrice float64, previousQuantity int64, startedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		StartedAt:        startedAt,
	})
}

func InsertServiceLineItemEndedWithParent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, previousBilledType enum.BilledType, previousPrice float64, previousQuantity int64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		StartedAt:        startedAt,
		EndedAt:          &endedAt,
	})
}

func InsertServiceLineItemCanceledWithParent(ctx context.Context, driver *neo4j.DriverWithContext, tenant, contractId string, billedType enum.BilledType, price float64, quantity int64, previousBilledType enum.BilledType, previousPrice float64, previousQuantity int64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	CreateServiceLineItemForContract(ctx, driver, tenant, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		IsCanceled:       true,
		StartedAt:        startedAt,
		EndedAt:          &endedAt,
	})
}
