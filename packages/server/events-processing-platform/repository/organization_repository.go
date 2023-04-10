package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"golang.org/x/net/context"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, aggregateId string, event events.OrganizationCreatedEvent) error
	UpdateOrganization(ctx context.Context, aggregateId string, event events.OrganizationUpdatedEvent) error
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, aggregateId string, event events.OrganizationCreatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(p:Organization:Organization_%s {id:$id}) 
		 ON CREATE SET 	p.name = $name,
						p.description = $description,
						p.website = $website,
						p.industry = $industry,
						p.isPublic = $isPublic,
						p.source = $source,
						p.sourceOfTruth = $sourceOfTruth,
						p.appSource = $appSource,
						p.createdAt = $createdAt,
						p.updatedAt = $updatedAt,
						p.syncedWithEventStore = true 
		 ON MATCH SET 	p.syncedWithEventStore = true
`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            aggregateId,
				"name":          event.Name,
				"description":   event.Description,
				"website":       event.Website,
				"industry":      event.Industry,
				"isPublic":      event.IsPublic,
				"tenant":        event.Tenant,
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

func (r *organizationRepository) UpdateOrganization(ctx context.Context, aggregateId string, event events.OrganizationUpdatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(p:Organization:Organization_%s {id:$id})
		 SET	p.name = $name,
				p.description = $description,
				p.website = $website,
				p.industry = $industry,
				p.isPublic = $isPublic,
				p.sourceOfTruth = $sourceOfTruth,
				p.updatedAt = $updatedAt,
				p.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            aggregateId,
				"tenant":        event.Tenant,
				"name":          event.Name,
				"description":   event.Description,
				"website":       event.Website,
				"industry":      event.Industry,
				"isPublic":      event.IsPublic,
				"sourceOfTruth": event.SourceOfTruth,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}
