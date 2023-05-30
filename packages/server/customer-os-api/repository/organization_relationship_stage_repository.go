package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type OrganizationRelationshipStageRepository interface {
	CreateDefaultStagesForNewTenant(ctx context.Context, tenant string) error
}

type organizationRelationshipStageRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRelationshipStageRepository(driver *neo4j.DriverWithContext) OrganizationRelationshipStageRepository {
	return &organizationRelationshipStageRepository{
		driver: driver,
	}
}

func (r *organizationRelationshipStageRepository) CreateDefaultStagesForNewTenant(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRelationshipStageRepository.CreateDefaultStagesForNewTenant")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentNeo4jRepository)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `WITH $stages AS stages
				UNWIND stages AS stage
				MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage {name: stage})
				ON CREATE SET 	s.id=randomUUID(), 
								s.createdAt=$now, 
								s:OrganizationRelationshipStage_%s
				WITH s
				MATCH (or:OrganizationRelationship)
				MERGE (s)<-[:HAS_STAGE]-(or)`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"tenant": tenant,
				"stages": entity.DefaultOrganizationRelationshipStageNames,
				"now":    utils.Now(),
			})
		return nil, err
	})
	return err
}
