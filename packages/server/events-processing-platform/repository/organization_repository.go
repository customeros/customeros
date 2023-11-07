package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, organizationId string, event events.OrganizationCreateEvent) error
	CreateOrganizationInTx(ctx context.Context, tx neo4j.ManagedTransaction, organizationId string, event events.OrganizationCreateEvent) error
	UpdateOrganization(ctx context.Context, organizationId string, event events.OrganizationUpdateEvent) error
	UpdateOrganizationIgnoreEmptyInputParams(ctx context.Context, organizationId string, event events.OrganizationUpdateEvent) error
	LinkWithDomain(ctx context.Context, tenant, organizationId, domain string) error
	OrganizationWebscrapedForDomain(ctx context.Context, tenant, organizationId, domain string) (bool, error)
	GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	UpdateRenewalLikelihood(ctx context.Context, orgId string, event events.OrganizationUpdateRenewalLikelihoodEvent) error
	UpdateRenewalForecast(ctx context.Context, orgId string, event events.OrganizationUpdateRenewalForecastEvent) error
	UpdateBillingDetails(ctx context.Context, orgId string, event events.OrganizationUpdateBillingDetailsEvent) error
	ReplaceOwner(ctx context.Context, tenant, organizationId, userId string) error
	SetVisibility(ctx context.Context, tenant, organizationId string, hide bool) error
	UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error
	SetCustomerOsIdIfMissing(ctx context.Context, tenant, organizationId, customerOsId string) error
	LinkWithParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId, subOrganizationType string) error
	UnlinkParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId string) error
	GetOrganizationIdsConnectedToInteractionEvent(ctx context.Context, tenant, interactionEventId string) ([]string, error)
}

type organizationRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext, database string) OrganizationRepository {
	return &organizationRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, organizationId string, event events.OrganizationCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.CreateOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.CreateOrganizationInTx(ctx, tx, organizationId, event)
	})
	return err
}

func (r *organizationRepository) CreateOrganizationInTx(ctx context.Context, tx neo4j.ManagedTransaction, organizationId string, event events.OrganizationCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.CreateOrganizationInTx")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("organizationId", organizationId))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id}) 
		 ON CREATE SET 	org.name = $name,
						org.description = $description,
						org.hide = $hide,
						org.website = $website,
						org.industry = $industry,
						org.subIndustry = $subIndustry,
						org.industryGroup = $industryGroup,
						org.targetAudience = $targetAudience,
						org.valueProposition = $valueProposition,
						org.lastFundingRound = $lastFundingRound,
						org.lastFundingAmount = $lastFundingAmount,
						org.referenceId = $referenceId,
						org.note = $note,
						org.isPublic = $isPublic,
						org.isCustomer = $isCustomer,
						org.source = $source,
						org.sourceOfTruth = $sourceOfTruth,
						org.employees = $employees,
						org.market = $market,
						org.appSource = $appSource,
						org.createdAt = $createdAt,
						org.updatedAt = $updatedAt,
						org.syncedWithEventStore = true 
		 ON MATCH SET 	org.name = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.name is null OR org.name = '' THEN $name ELSE org.name END,
						org.description = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.description is null OR org.description = '' THEN $description ELSE org.description END,
						org.hide = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$sourceOfTruth AND $hide = false) THEN $hide ELSE org.hide END,
						org.isCustomer = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$sourceOfTruth AND $isCustomer = true) THEN $isCustomer ELSE org.isCustomer END,
						org.website = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.website is null OR org.website = '' THEN $website ELSE org.website END,
						org.industry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.industry is null OR org.industry = '' THEN $industry ELSE org.industry END,
						org.subIndustry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.subIndustry is null OR org.subIndustry = '' THEN $subIndustry ELSE org.subIndustry END,
						org.industryGroup = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.industryGroup is null OR org.industryGroup = '' THEN $industryGroup ELSE org.industryGroup END,
						org.targetAudience = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.targetAudience is null OR org.targetAudience = '' THEN $targetAudience ELSE org.targetAudience END,
						org.valueProposition = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.valueProposition is null OR org.valueProposition = '' THEN $valueProposition ELSE org.valueProposition END,
						org.lastFundingRound = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.lastFundingRound is null OR org.lastFundingRound = '' THEN $lastFundingRound ELSE org.lastFundingRound END,
						org.lastFundingAmount = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.lastFundingAmount is null OR org.lastFundingAmount = '' THEN $lastFundingAmount ELSE org.lastFundingAmount END,
						org.referenceId = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.referenceId is null OR org.referenceId = '' THEN $referenceId ELSE org.referenceId END,
						org.note = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.note is null OR org.note = '' THEN $note ELSE org.note END,
						org.isPublic = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.isPublic is null THEN $isPublic ELSE org.isPublic END,
						org.employees = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.employees is null THEN $employees ELSE org.employees END,
						org.market = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.market is null OR org.market = '' THEN $market ELSE org.market END,
						org.updatedAt=$updatedAt,
						org.syncedWithEventStore = true`, event.Tenant)
	params := map[string]any{
		"id":                organizationId,
		"name":              event.Name,
		"hide":              event.Hide,
		"description":       event.Description,
		"website":           event.Website,
		"industry":          event.Industry,
		"subIndustry":       event.SubIndustry,
		"industryGroup":     event.IndustryGroup,
		"targetAudience":    event.TargetAudience,
		"valueProposition":  event.ValueProposition,
		"isPublic":          event.IsPublic,
		"isCustomer":        event.IsCustomer,
		"tenant":            event.Tenant,
		"employees":         event.Employees,
		"market":            event.Market,
		"lastFundingRound":  event.LastFundingRound,
		"lastFundingAmount": event.LastFundingAmount,
		"referenceId":       event.ReferenceId,
		"note":              event.Note,
		"source":            helper.GetSource(event.Source),
		"sourceOfTruth":     helper.GetSource(event.SourceOfTruth),
		"appSource":         helper.GetSource(event.AppSource),
		"createdAt":         event.CreatedAt,
		"updatedAt":         event.UpdatedAt,
		"overwrite":         helper.GetSource(event.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return utils.ExecuteQueryInTx(ctx, tx, query, params)
}

func (r *organizationRepository) UpdateOrganization(ctx context.Context, organizationId string, event events.OrganizationUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("organizationId", organizationId))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id})
		 SET	org.name = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.name is null OR org.name = '' THEN $name ELSE org.name END,
				org.description = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.description is null OR org.description = '' THEN $description ELSE org.description END,
				org.hide = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$sourceOfTruth AND $hide = false) THEN $hide ELSE org.hide END,
				org.isCustomer = CASE WHEN $overwrite=true OR (org.sourceOfTruth=$sourceOfTruth AND $isCustomer = true) THEN $isCustomer ELSE org.isCustomer END,				
				org.website = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.website is null OR org.website = '' THEN $website ELSE org.website END,
				org.industry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.industry is null OR org.industry = '' THEN $industry ELSE org.industry END,
				org.subIndustry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.subIndustry is null OR org.subIndustry = '' THEN $subIndustry ELSE org.subIndustry END,
				org.industryGroup = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.industryGroup is null OR org.industryGroup = '' THEN $industryGroup ELSE org.industryGroup END,
				org.targetAudience = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.targetAudience is null OR org.targetAudience = '' THEN $targetAudience ELSE org.targetAudience END,
				org.valueProposition = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.valueProposition is null OR org.valueProposition = '' THEN $valueProposition ELSE org.valueProposition END,
				org.lastFundingRound = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.lastFundingRound is null OR org.lastFundingRound = '' THEN $lastFundingRound ELSE org.lastFundingRound END,
				org.lastFundingAmount = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.lastFundingAmount is null OR org.lastFundingAmount = '' THEN $lastFundingAmount ELSE org.lastFundingAmount END,
				org.referenceId = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.referenceId is null OR org.referenceId = '' THEN $referenceId ELSE org.referenceId END,
				org.note = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.note is null OR org.note = '' THEN $note ELSE org.note END,
				org.isPublic = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.isPublic is null THEN $isPublic ELSE org.isPublic END,
				org.employees = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.employees is null THEN $employees ELSE org.employees END,
				org.market = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR org.market is null OR org.market = '' THEN $market ELSE org.market END,
				org.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE org.sourceOfTruth END,
				org.updatedAt = $updatedAt,
				org.syncedWithEventStore = true`, event.Tenant)

	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
		"id":                organizationId,
		"tenant":            event.Tenant,
		"name":              event.Name,
		"hide":              event.Hide,
		"description":       event.Description,
		"website":           event.Website,
		"industry":          event.Industry,
		"subIndustry":       event.SubIndustry,
		"industryGroup":     event.IndustryGroup,
		"targetAudience":    event.TargetAudience,
		"valueProposition":  event.ValueProposition,
		"isPublic":          event.IsPublic,
		"isCustomer":        event.IsCustomer,
		"employees":         event.Employees,
		"market":            event.Market,
		"lastFundingRound":  event.LastFundingRound,
		"lastFundingAmount": event.LastFundingAmount,
		"referenceId":       event.ReferenceId,
		"note":              event.Note,
		"sourceOfTruth":     helper.GetSource(event.Source),
		"updatedAt":         event.UpdatedAt,
		"overwrite":         helper.GetSource(event.Source) == constants.SourceOpenline,
	})
}

func (r *organizationRepository) UpdateOrganizationIgnoreEmptyInputParams(ctx context.Context, organizationId string, event events.OrganizationUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateOrganizationIgnoreEmptyInputParams")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("organizationId", organizationId))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id})
		 SET	org.name = CASE WHEN $name <> '' THEN $name ELSE org.name END, 
				org.description = CASE WHEN $description <> '' THEN $description ELSE org.description END, 
				org.website = CASE WHEN $website <> '' THEN $website ELSE org.website END, 
				org.industry = CASE WHEN $industry <> '' THEN $industry ELSE org.industry END, 
				org.subIndustry = CASE WHEN $subIndustry <> '' THEN $subIndustry ELSE org.subIndustry END, 
				org.industryGroup = CASE WHEN $industryGroup <> '' THEN $industryGroup ELSE org.industryGroup END, 
				org.targetAudience = CASE WHEN $targetAudience <> '' THEN $targetAudience ELSE org.targetAudience END, 
				org.valueProposition = CASE WHEN $valueProposition <> '' THEN $valueProposition ELSE org.valueProposition END, 
				org.lastFundingRound = CASE WHEN $lastFundingRound <> '' THEN $lastFundingRound ELSE org.lastFundingRound END, 
				org.lastFundingAmount = CASE WHEN $lastFundingAmount <> '' THEN $lastFundingAmount ELSE org.lastFundingAmount END, 
				org.referenceId = CASE WHEN $referenceId <> '' THEN $referenceId ELSE org.referenceId END, 
				org.note = CASE WHEN $note <> '' THEN $note ELSE org.note END, 
				org.market = CASE WHEN $market <> '' THEN $market ELSE org.market END, 
				org.employees = CASE WHEN $employees <> 0 THEN $employees ELSE org.employees END, 
				org.isCustomer = CASE WHEN $isCustomer = true THEN $isCustomer ELSE org.isCustomer END, 
				org.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE org.sourceOfTruth END,
				org.updatedAt = $updatedAt,
				org.syncedWithEventStore = true`, event.Tenant)
	params := map[string]any{
		"id":                organizationId,
		"tenant":            event.Tenant,
		"name":              event.Name,
		"description":       event.Description,
		"website":           event.Website,
		"industry":          event.Industry,
		"subIndustry":       event.SubIndustry,
		"industryGroup":     event.IndustryGroup,
		"targetAudience":    event.TargetAudience,
		"valueProposition":  event.ValueProposition,
		"employees":         event.Employees,
		"market":            event.Market,
		"lastFundingRound":  event.LastFundingRound,
		"lastFundingAmount": event.LastFundingAmount,
		"referenceId":       event.ReferenceId,
		"note":              event.Note,
		"isCustomer":        event.IsCustomer,
		"sourceOfTruth":     helper.GetSource(event.Source),
		"updatedAt":         event.UpdatedAt,
		"overwrite":         helper.GetSource(event.Source) == constants.SourceOpenline,
	}

	span.LogFields(log.String("query", query), log.Object("params", params))

	return r.executeQuery(ctx, query, params)
}

func (r *organizationRepository) LinkWithDomain(ctx context.Context, tenant, organizationId, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeOrganizationDomain")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId))

	query := `MERGE (d:Domain {domain:$domain}) 
				ON CREATE SET 	d.id=randomUUID(), 
								d.createdAt=$now, 
								d.updatedAt=$now,
								d.appSource=$appSource
				WITH d
				MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				MERGE (org)-[rel:HAS_DOMAIN]->(d)
				SET rel.syncedWithEventStore = true
				RETURN rel`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"domain":         strings.ToLower(domain),
				"appSource":      constants.AppSourceEventProcessingPlatform,
				"now":            utils.Now(),
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *organizationRepository) OrganizationWebscrapedForDomain(ctx context.Context, tenant, organizationId, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.OrganizationWebscrapedForDomain")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("domain", domain))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})-[:HAS_DOMAIN]->(d:Domain {domain:$domain})
				WHERE org.sourceOfTruth = $webscrape
				RETURN org`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"webscrape":      constants.SourceWebscrape,
				"tenant":         tenant,
				"organizationId": organizationId,
				"domain":         strings.ToLower(domain),
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	return len(dbRecords.([]*neo4j.Record)) > 0, err
}

func (r *organizationRepository) GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id}) RETURN org`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"id":     organizationId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) UpdateRenewalLikelihood(ctx context.Context, orgId string, event events.OrganizationUpdateRenewalLikelihoodEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateRenewalLikelihood")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("organizationId", orgId), log.Object("event", event))

	query := ` MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			 SET 	org.renewalLikelihoodPrevious = CASE WHEN org.renewalLikelihood <> $renewalLikelihood THEN org.renewalLikelihood ELSE org.renewalLikelihoodPrevious END,
					org.renewalLikelihood=$renewalLikelihood,					
					org.renewalLikelihoodComment=$comment, 
			 		org.renewalLikelihoodUpdatedBy=$updatedBy, 
					org.renewalLikelihoodUpdatedAt=$updatedAt,
					org.updatedAt=$now, 
					org.sourceOfTruth=$source`
	span.LogFields(log.String("query", query))

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":            event.Tenant,
		"organizationId":    orgId,
		"renewalLikelihood": event.GetRenewalLikelihoodAsStringForGraphDb(),
		"comment":           event.Comment,
		"updatedBy":         event.UpdatedBy,
		"updatedAt":         event.UpdatedAt,
		"source":            constants.SourceOpenline,
		"now":               utils.Now(),
	})
}

func (r *organizationRepository) UpdateRenewalForecast(ctx context.Context, orgId string, event events.OrganizationUpdateRenewalForecastEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateRenewalForecast")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("organizationId", orgId), log.Object("event", event))

	query := ` MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			 SET 	org.renewalForecastAmount=$renewalForecast, 
					org.renewalForecastPotentialAmount=CASE WHEN $updatedBy = '' THEN $renewalForecastPotential ELSE org.renewalForecastPotentialAmount END, 
					org.renewalForecastComment=$comment, 
			 		org.renewalForecastUpdatedBy=$updatedBy, 
					org.renewalForecastUpdatedAt=$updatedAt,
					org.updatedAt=$now, 
					org.sourceOfTruth=$source`
	span.LogFields(log.String("query", query))

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":                   event.Tenant,
		"organizationId":           orgId,
		"renewalForecast":          event.Amount,
		"renewalForecastPotential": event.PotentialAmount,
		"comment":                  event.Comment,
		"updatedBy":                event.UpdatedBy,
		"updatedAt":                event.UpdatedAt,
		"source":                   constants.SourceOpenline,
		"now":                      utils.Now(),
	})
}

func (r *organizationRepository) UpdateBillingDetails(ctx context.Context, orgId string, event events.OrganizationUpdateBillingDetailsEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateBillingDetails")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("organizationId", orgId), log.Object("event", event))

	query := ` MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			 SET 	org.billingDetailsAmount=$amount, 
					org.billingDetailsFrequency=$frequency, 
					org.billingDetailsRenewalCycle=$renewalCycle, 
			 		org.billingDetailsRenewalCycleStart=$renewalCycleStart,
					org.billingDetailsRenewalCycleNext = CASE WHEN $updatedBy = '' THEN $renewalCycleNext ELSE org.billingDetailsRenewalCycleNext END,
					org.updatedAt=$now, 
					org.sourceOfTruth=$source`
	span.LogFields(log.String("query", query))

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":            event.Tenant,
		"organizationId":    orgId,
		"amount":            event.Amount,
		"frequency":         event.Frequency,
		"renewalCycle":      event.RenewalCycle,
		"renewalCycleStart": utils.TimePtrFirstNonNilNillableAsAny(event.RenewalCycleStart),
		"renewalCycleNext":  utils.TimePtrFirstNonNilNillableAsAny(event.RenewalCycleNext),
		"source":            constants.SourceOpenline,
		"now":               utils.Now(),
		"updatedBy":         event.UpdatedBy,
	})
}

func (r *organizationRepository) ReplaceOwner(ctx context.Context, tenant, organizationId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.ReplaceOwner")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("userId", userId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (:User)-[rel:OWNS]->(org)
			DELETE rel
			WITH org, t
			MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			WHERE (u.internal=false OR u.internal is null) AND (u.bot=false OR u.bot is null)
			MERGE (u)-[:OWNS]->(org)
			SET org.updatedAt=$now, org.sourceOfTruth=$source`

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"userId":         userId,
		"source":         constants.SourceOpenline,
		"now":            utils.Now(),
	})
}

func (r *organizationRepository) SetVisibility(ctx context.Context, tenant, organizationId string, hide bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.SetVisibility")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.Bool("hide", hide))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id})
		 SET	org.hide = $hide,
				org.updatedAt = $now`, tenant)

	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
		"id":     organizationId,
		"tenant": tenant,
		"hide":   hide,
		"now":    utils.Now(),
	})
}

func (r *organizationRepository) UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.SetVisibility")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("touchpointId", touchpointId), log.Object("touchpointAt", touchpointAt))

	query := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.lastTouchpointAt=$touchpointAt, org.lastTouchpointId=$touchpointId`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"touchpointAt":   touchpointAt,
		"touchpointId":   touchpointId,
	})
}

func (r *organizationRepository) SetCustomerOsIdIfMissing(ctx context.Context, tenant, organizationId, customerOsId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.SetCustomerOsIdIfMissing")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("customerOsId", customerOsId))

	query := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.customerOsId = CASE WHEN (org.customerOsId IS NULL OR org.customerOsId = '') AND $customerOsId <> '' THEN $customerOsId ELSE org.customerOsId END`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"customerOsId":   customerOsId,
	})
}

func (r *organizationRepository) LinkWithParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId, subOrganizationType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.LinkWithParentOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("parentOrganizationId", parentOrganizationId), log.String("subOrganizationType", subOrganizationType))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$parentOrganizationId}),
		 			(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(sub:Organization {id:$subOrganizationId}) 
		 	MERGE (sub)-[rel:SUBSIDIARY_OF]->(parent) 
		 		ON CREATE SET rel.type=$type 
		 		ON MATCH SET rel.type=$type`
	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":               tenant,
		"subOrganizationId":    organizationId,
		"parentOrganizationId": parentOrganizationId,
		"type":                 subOrganizationType,
	})
}

func (r *organizationRepository) UnlinkParentOrganization(ctx context.Context, tenant, organizationId, parentOrganizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UnlinkParentOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId), log.String("parentOrganizationId", parentOrganizationId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$parentOrganizationId})<-[rel:SUBSIDIARY_OF]-(sub:Organization {id:$subOrganizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
		 		DELETE rel`
	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":               tenant,
		"subOrganizationId":    organizationId,
		"parentOrganizationId": parentOrganizationId,
	})
}

func (r *organizationRepository) GetOrganizationIdsConnectedToInteractionEvent(ctx context.Context, tenant, interactionEventId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdsConnectedToInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	query := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId}),
				(t:Tenant {name:$tenant})
				CALL {
					WITH ie, t 
					MATCH (ie)-[:PART_OF]->(is:Issue)-[:REPORTED_BY]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
					RETURN org.id as orgId
				UNION 
					WITH ie, t 
					MATCH (ie)-[:PART_OF]->(is:Issue)-[:SUBMITTED_BY]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
					RETURN org.id as orgId
				}
				RETURN distinct orgId`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]string), err
}

// Common database interaction method
func (r *organizationRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
