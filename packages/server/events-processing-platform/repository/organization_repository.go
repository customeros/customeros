package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, organizationId string, event events.OrganizationCreateEvent) error
	UpdateOrganization(ctx context.Context, organizationId string, event events.OrganizationUpdateEvent) error
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
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

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"id":            organizationId,
				"name":          event.Name,
				"description":   event.Description,
				"website":       event.Website,
				"industry":      event.Industry,
				"isPublic":      event.IsPublic,
				"tenant":        event.Tenant,
				"employees":     event.Employees,
				"market":        event.Market,
				"source":        event.Source,
				"sourceOfTruth": event.SourceOfTruth,
				"appSource":     event.AppSource,
				"createdAt":     event.CreatedAt,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
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
				org.isPublic = $isPublic,
				org.employees = $employees,
				org.market = $market,	
				org.sourceOfTruth = $sourceOfTruth,
				org.updatedAt = $updatedAt,
				org.syncedWithEventStore = true`, event.Tenant)

	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"id":            organizationId,
				"tenant":        event.Tenant,
				"name":          event.Name,
				"description":   event.Description,
				"website":       event.Website,
				"industry":      event.Industry,
				"isPublic":      event.IsPublic,
				"employees":     event.Employees,
				"market":        event.Market,
				"sourceOfTruth": event.SourceOfTruth,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}
