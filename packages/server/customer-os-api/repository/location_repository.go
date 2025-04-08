package repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go/log"
)

type LocationRepository interface {
	CreateLocationForEntity(ctx context.Context, fromContext string, entityType model.EntityType, id string, source entity.SourceFields) (*dbtype.Node, error)
	Update(ctx context.Context, tenant string, locationEntity neo4jentity.LocationEntity) (*dbtype.Node, error)
	RemoveRelationshipAndDeleteOrphans(ctx context.Context, entityType model.EntityType, entityId, locationId string) error
}

type locationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewLocationRepository(driver *neo4j.DriverWithContext) LocationRepository {
	return &locationRepository{
		driver: driver,
	}
}

func (r *locationRepository) CreateLocationForEntity(ctx context.Context, tenant string, entityType model.EntityType, entityId string, source entity.SourceFields) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LocationRepository.CreateLocationForEntity")
	defer spans.Finish()

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (e:%s {id:$entityId}), (t:Tenant {name:$tenant})
		 MERGE (e)-[:ASSOCIATED_WITH]->(loc:Location {id:randomUUID()})-[:LOCATION_BELONGS_TO_TENANT]->(t)
		 ON CREATE SET 
		  loc.createdAt=$now, 
		  loc.updatedAt=datetime(), 
		  loc.source=$source, 
		  loc.appSource=$appSource, 
		  loc:%s
		 RETURN loc`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, entityType.Neo4jLabel()+"_"+tenant, "Location_"+tenant),
			map[string]any{
				"tenant":    tenant,
				"now":       utils.Now(),
				"entityId":  entityId,
				"source":    source.Source,
				"appSource": source.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *locationRepository) Update(ctx context.Context, tenant string, locationEntity neo4jentity.LocationEntity) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LocationRepository.Update")
	defer spans.Finish()

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(loc:Location {id:$id})
			SET loc.updatedAt=datetime(),
				loc.name=$name,
				loc.rawAddress=$rawAddress,
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
				"tenant":       tenant,
				"now":          utils.Now(),
				"id":           locationEntity.Id,
				"name":         locationEntity.Name,
				"rawAddress":   locationEntity.RawAddress,
				"country":      locationEntity.Country,
				"region":       locationEntity.Region,
				"locality":     locationEntity.Locality,
				"address":      locationEntity.Address,
				"address2":     locationEntity.Address2,
				"zip":          locationEntity.Zip,
				"addressType":  locationEntity.AddressType,
				"houseNumber":  locationEntity.HouseNumber,
				"postalCode":   locationEntity.PostalCode,
				"plusFour":     locationEntity.PlusFour,
				"commercial":   locationEntity.Commercial,
				"predirection": locationEntity.Predirection,
				"district":     locationEntity.District,
				"street":       locationEntity.Street,
				"latitude":     locationEntity.Latitude,
				"longitude":    locationEntity.Longitude,
				"timeZone":     locationEntity.TimeZone,
				"utcOffset":    locationEntity.UtcOffset,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *locationRepository) RemoveRelationshipAndDeleteOrphans(ctx context.Context, entityType model.EntityType, entityId, locationId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LocationRepository.RemoveRelationshipAndDeleteOrphans")
	defer spans.Finish()

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:LOCATION_BELONGS_TO_TENANT]-(loc:Location {id:$locationId}),
										(loc)<-[rel:ASSOCIATED_WITH]-(entity:%s {id:$entityId})
								DELETE rel 
								WITH t, loc
								WHERE NOT EXISTS {
  									MATCH (loc)--(node)
									WHERE NOT (node:Tenant)
								}
								DETACH DELETE loc`, entityType.Neo4jLabel()+"_"+common.GetTenantFromContext(ctx))
	spans.LogKV(log.String("query", query))

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
