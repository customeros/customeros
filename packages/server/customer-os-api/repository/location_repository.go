package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type LocationRepository interface {
	GetAllForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error)
	GetAllForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error)
	GetAllForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	CreateLocationForEntity(ctx context.Context, fromContext string, entityType entity.EntityType, id string, source entity.SourceFields) (*dbtype.Node, error)
}

type locationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewLocationRepository(driver *neo4j.DriverWithContext) LocationRepository {
	return &locationRepository{
		driver: driver,
	}
}

func (r *locationRepository) GetAllForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(:Contact {id:$contactId})-[:ASSOCIATED_WITH]->(loc:Location)
			RETURN loc`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *locationRepository) GetAllForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:ASSOCIATED_WITH]->(loc:Location)-[:LOCATION_BELONGS_TO_TENANT]->(t)
			WHERE c.id IN $contactIds
			RETURN loc, c.id as contactId ORDER BY loc.name`,
			map[string]any{
				"tenant":     tenant,
				"contactIds": contactIds,
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

func (r *locationRepository) GetAllForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$organizationId})-[:ASSOCIATED_WITH]->(loc:Location)
			RETURN loc`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *locationRepository) GetAllForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:ASSOCIATED_WITH]->(loc:Location)-[:LOCATION_BELONGS_TO_TENANT]->(t)
			WHERE o.id IN $organizationIds
			RETURN loc, o.id as organizationId ORDER BY loc.name`,
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

func (r *locationRepository) CreateLocationForEntity(ctx context.Context, tenant string, entityType entity.EntityType, entityId string, source entity.SourceFields) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (e:%s {id:$entityId}) 
		 MERGE (e)-[:ASSOCIATED_WITH]->(loc:Location {id:randomUUID()}) 
		 ON CREATE SET 
		  loc.createdAt=$now, 
		  loc.updatedAt=$now, 
		  loc.source=$source, 
		  loc.sourceOfTruth=$sourceOfTruth, 
		  loc.appSource=$appSource, 
		  loc:%s
		 RETURN loc`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, entityType.Neo4jLabel()+"_"+tenant, "Location_"+tenant),
			map[string]any{
				"now":           utils.Now(),
				"entityId":      entityId,
				"source":        source.Source,
				"sourceOfTruth": source.SourceOfTruth,
				"appSource":     source.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
