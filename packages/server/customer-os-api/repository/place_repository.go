package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type PlaceRepository interface {
	GetAnyForLocation(tenant, locationId string) ([]*dbtype.Node, error)
}

type placeRepository struct {
	driver *neo4j.Driver
}

func NewPlaceRepository(driver *neo4j.Driver) PlaceRepository {
	return &placeRepository{
		driver: driver,
	}
}

func (r *placeRepository) GetAnyForLocation(tenant, locationId string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (loc:Location_%s {id:$locationId})-[:LOCATED_AT]->(place:Place) RETURN place"

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(fmt.Sprintf(query, tenant),
			map[string]any{
				"locationId": locationId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsNodePtrs(queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}
