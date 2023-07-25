package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, organizationId string, event events.OrganizationCreateEvent) error
	UpdateOrganization(ctx context.Context, organizationId string, event events.OrganizationUpdateEvent) error
	UpdateOrganizationIgnoreEmptyInputParams(ctx context.Context, organizationId string, event events.OrganizationUpdateEvent) error
	LinkWithDomain(ctx context.Context, tenant, organizationId, domain string) error
	OrganizationWebscrapedForDomain(ctx context.Context, tenant, organizationId, domain string) (bool, error)
	GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func (r *organizationRepository) GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("organizationId", organizationId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id}) RETURN org`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, organizationId string, event events.OrganizationCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.CreateOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("organizationId", organizationId))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id}) 
		 ON CREATE SET 	org.name = $name,
						org.description = $description,
						org.website = $website,
						org.industry = $industry,
						org.subIndustry = $subIndustry,
						org.industryGroup = $industryGroup,
						org.targetAudience = $targetAudience,
						org.valueProposition = $valueProposition,
						org.lastFundingRound = $lastFundingRound,
						org.lastFundingAmount = $lastFundingAmount,
						org.isPublic = $isPublic,
						org.source = $source,
						org.sourceOfTruth = $sourceOfTruth,
						org.employees = $employees,
						org.market = $market,
						org.appSource = $appSource,
						org.createdAt = $createdAt,
						org.updatedAt = $updatedAt,
						org.syncedWithEventStore = true 
		 ON MATCH SET 	org.syncedWithEventStore = true`, event.Tenant)
	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
		"id":                organizationId,
		"name":              event.Name,
		"description":       event.Description,
		"website":           event.Website,
		"industry":          event.Industry,
		"subIndustry":       event.SubIndustry,
		"industryGroup":     event.IndustryGroup,
		"targetAudience":    event.TargetAudience,
		"valueProposition":  event.ValueProposition,
		"isPublic":          event.IsPublic,
		"tenant":            event.Tenant,
		"employees":         event.Employees,
		"market":            event.Market,
		"lastFundingRound":  event.LastFundingRound,
		"lastFundingAmount": event.LastFundingAmount,
		"source":            event.Source,
		"sourceOfTruth":     event.SourceOfTruth,
		"appSource":         event.AppSource,
		"createdAt":         event.CreatedAt,
		"updatedAt":         event.UpdatedAt,
	})
}

func (r *organizationRepository) UpdateOrganization(ctx context.Context, organizationId string, event events.OrganizationUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("organizationId", organizationId))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization:Organization_%s {id:$id})
		 SET	org.name = $name,
				org.description = $description,
				org.website = $website,
				org.industry = $industry,
				org.subIndustry = $subIndustry,
				org.industryGroup = $industryGroup,
				org.targetAudience = $targetAudience,
				org.valueProposition = $valueProposition,
				org.lastFundingRound = $lastFundingRound,
				org.lastFundingAmount = $lastFundingAmount,
				org.isPublic = $isPublic,
				org.employees = $employees,
				org.market = $market,	
				org.sourceOfTruth = $sourceOfTruth,
				org.updatedAt = $updatedAt,
				org.syncedWithEventStore = true`, event.Tenant)

	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
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
		"isPublic":          event.IsPublic,
		"employees":         event.Employees,
		"market":            event.Market,
		"lastFundingRound":  event.LastFundingRound,
		"lastFundingAmount": event.LastFundingAmount,
		"sourceOfTruth":     event.SourceOfTruth,
		"updatedAt":         event.UpdatedAt,
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
				org.market = CASE WHEN $market <> '' THEN $market ELSE org.market END, 
				org.employees = CASE WHEN $employees <> 0 THEN $employees ELSE org.employees END, 
				org.sourceOfTruth = $sourceOfTruth,
				org.isPublic = $isPublic,
				org.updatedAt = $updatedAt,
				org.syncedWithEventStore = true`, event.Tenant)

	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
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
		"isPublic":          event.IsPublic,
		"employees":         event.Employees,
		"market":            event.Market,
		"lastFundingRound":  event.LastFundingRound,
		"lastFundingAmount": event.LastFundingAmount,
		"sourceOfTruth":     event.SourceOfTruth,
		"updatedAt":         event.UpdatedAt,
	})
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

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
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

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
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

// Common database interaction method
func (r *organizationRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteQuery(ctx, *r.driver, query, params)
}
