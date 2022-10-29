package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
)

type TextCustomPropertyService interface {
}

type textCustomPropertyService struct {
	driver *neo4j.Driver
}

func NewTextCustomPropertyService(driver *neo4j.Driver) TextCustomPropertyService {
	return &textCustomPropertyService{
		driver: driver,
	}
}

func addTextCustomFieldToContact(ctx context.Context, contactId string, input entity.TextCustomFieldEntity, tx neo4j.Transaction) (interface{}, error) {
	result, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})
			CREATE (f:TextCustomField {
				  group: $group,
				  name: $name,
				  value: $value
			})<-[:HAS_PROPERTY]-(c)
			RETURN f`,
		map[string]interface{}{
			"contactId": contactId,
			"group":     input.Group,
			"name":      input.Name,
			"value":     input.Value,
		})

	record, err := result.Single()
	if err != nil {
		return nil, err
	}
	return record, nil
}
