package service

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactGroupService interface {
	Create(contactGroup *entity.ContactGroupNode) (*entity.ContactGroupNode, error)
	FindAll() (*entity.ContactGroupNodes, error)
}

type contactGroupService struct {
	driver *neo4j.Driver
}

func NewContactGroupService(driver *neo4j.Driver) ContactGroupService {
	return &contactGroupService{
		driver: driver,
	}
}

func (s *contactGroupService) Create(newContactGroup *entity.ContactGroupNode) (*entity.ContactGroupNode, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			CREATE (g:ContactGroup {
				  id: randomUUID(),
				  name: $name
			})
			RETURN g`,
			map[string]interface{}{
				"name": newContactGroup.Name,
			})

		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}
	contactGroup := entity.ContactGroupNode{}
	mapstructure.Decode(utils.GetPropsFromNode(queryResult.(dbtype.Node)), &contactGroup)

	return &contactGroup, nil
}

func (s *contactGroupService) FindAll() (*entity.ContactGroupNodes, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`MATCH (g:ContactGroup) RETURN g`, map[string]interface{}{})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	contactGroups := entity.ContactGroupNodes{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		contactGroup := entity.ContactGroupNode{}
		mapstructure.Decode(utils.GetPropsFromNode(dbRecord.Values[0].(dbtype.Node)), &contactGroup)
		contactGroups = append(contactGroups, contactGroup)
	}

	return &contactGroups, nil
}
