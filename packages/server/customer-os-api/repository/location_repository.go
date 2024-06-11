package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LocationRepository interface {
	GetAllForContact(ctx context.Context, tenant, contactId string) ([]*dbtype.Node, error)
	GetAllForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error)
	GetAllForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	CreateLocationForEntity(ctx context.Context, fromContext string, entityType entity.EntityType, id string, source entity.SourceFields) (*dbtype.Node, error)
	Update(ctx context.Context, tenant string, locationEntity entity.LocationEntity) (*dbtype.Node, error)
	RemoveRelationshipAndDeleteOrphans(ctx context.Context, entityType entity.EntityType, entityId, locationId string) error
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.GetAllForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *locationRepository) GetAllForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.GetAllForContacts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.GetAllForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *locationRepository) GetAllForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.GetAllForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.CreateLocationForEntity")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (e:%s {id:$entityId}), (t:Tenant {name:$tenant})
		 MERGE (e)-[:ASSOCIATED_WITH]->(loc:Location {id:randomUUID()})-[:LOCATION_BELONGS_TO_TENANT]->(t)
		 ON CREATE SET 
		  loc.createdAt=$now, 
		  loc.updatedAt=datetime(), 
		  loc.source=$source, 
		  loc.sourceOfTruth=$sourceOfTruth, 
		  loc.appSource=$appSource, 
		  loc:%s
		 RETURN loc`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, entityType.Neo4jLabel()+"_"+tenant, "Location_"+tenant),
			map[string]any{
				"tenant":        tenant,
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

func (r *locationRepository) Update(ctx context.Context, tenant string, locationEntity entity.LocationEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(loc:Location {id:$id})
			SET loc.updatedAt=datetime(),
				loc.name=$name,
				loc.rawAddress=$rawAddress,
				loc.sourceOfTruth=$sourceOfTruth,
				loc.country=$country,	
				loc.region=$region,
				loc.locality=$locality,
				loc.address=$address,
				loc.address2=$address2,
				loc.zip=$zip,
				loc.addressType=$addressType,
				loc.houseNumber=$houseNumber,
				loc.postalCode=$postalCode,
				loc.plusFour=$plusFour,
				loc.commercial=$commercial,
				loc.predirection=$predirection,
				loc.district=$district,
				loc.street=$street,	
				loc.latitude=$latitude,
				loc.longitude=$longitude,	
				loc.timeZone=$timeZone,
				loc.utcOffset=$utcOffset
			RETURN loc`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":        tenant,
				"now":           utils.Now(),
				"id":            locationEntity.Id,
				"name":          locationEntity.Name,
				"rawAddress":    locationEntity.RawAddress,
				"sourceOfTruth": locationEntity.SourceOfTruth,
				"country":       locationEntity.Country,
				"region":        locationEntity.Region,
				"locality":      locationEntity.Locality,
				"address":       locationEntity.Address,
				"address2":      locationEntity.Address2,
				"zip":           locationEntity.Zip,
				"addressType":   locationEntity.AddressType,
				"houseNumber":   locationEntity.HouseNumber,
				"postalCode":    locationEntity.PostalCode,
				"plusFour":      locationEntity.PlusFour,
				"commercial":    locationEntity.Commercial,
				"predirection":  locationEntity.Predirection,
				"district":      locationEntity.District,
				"street":        locationEntity.Street,
				"latitude":      locationEntity.Latitude,
				"longitude":     locationEntity.Longitude,
				"timeZone":      locationEntity.TimeZone,
				"utcOffset":     locationEntity.UtcOffset,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *locationRepository) RemoveRelationshipAndDeleteOrphans(ctx context.Context, entityType entity.EntityType, entityId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationRepository.RemoveRelationshipAndDeleteOrphans")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(loc:Location {id:$locationId}),
										(loc)<-[rel:ASSOCIATED_WITH]-(entity:%s {id:$entityId})
								DELETE rel 
								WITH t, loc
								WHERE NOT EXISTS {
  									MATCH (loc)--(node)
									WHERE NOT (node:Tenant)
								}
								DETACH DELETE loc`, entityType.Neo4jLabel()+"_"+common.GetTenantFromContext(ctx))
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"entityId":   entityId,
				"locationId": locationId,
				"tenant":     common.GetTenantFromContext(ctx),
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}
