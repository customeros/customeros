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

type OrganizationRelationshipRepository interface {
	GetOrganizationRelationshipsForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	GetOrganizationRelationshipsWithStagesForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodePairAndId, error)
	CreateDefaultStagesForNewTenant(ctx context.Context, tenant string) error
}

type organizationRelationshipRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRelationshipRepository(driver *neo4j.DriverWithContext) OrganizationRelationshipRepository {
	return &organizationRelationshipRepository{
		driver: driver,
	}
}

func (r *organizationRelationshipRepository) GetOrganizationRelationshipsForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRelationshipRepository.GetOrganizationRelationshipsForOrganizations")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentNeo4jRepository)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:IS]->(or:OrganizationRelationship)
			WHERE o.id IN $organizationIds
			RETURN or, o.id`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":          tenant,
				"organizationIds": organizationIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationRelationshipRepository) GetOrganizationRelationshipsWithStagesForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodePairAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRelationshipRepository.GetOrganizationRelationshipsForOrganizations")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentNeo4jRepository)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:IS]->(or:OrganizationRelationship)
			OPTIONAL MATCH (or)-[:HAS_STAGE]->(ors:OrganizationRelationshipStage)<-[:HAS_STAGE]-(org)
			WHERE o.id IN $organizationIds
			RETURN or, ors, o.id`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":          tenant,
				"organizationIds": organizationIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodePairAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodePairAndId), err
}

func (r *organizationRelationshipRepository) CreateDefaultStagesForNewTenant(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRelationshipRepository.CreateDefaultStagesForNewTenant")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, constants.ComponentNeo4jRepository)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `WITH $stages AS stages
				UNWIND stages AS stage
				MATCH (t:Tenant {name:$tenant}), (or:OrganizationRelationship)
				MERGE (t)<-[:STAGE_BELONGS_TO_TENANT]-(s:OrganizationRelationshipStage {name: stage})<-[:HAS_STAGE]-(or)
				ON CREATE SET 	s.id=randomUUID(), 
								s.createdAt=$now, 
								s:OrganizationRelationshipStage_%s`

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
