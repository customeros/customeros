package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type CompanyPositionService interface {
	FindAllForContact(ctx context.Context, obj *model.Contact) (*entity.CompanyPositionEntities, error)
}

type companyPositionService struct {
	repository *repository.RepositoryContainer
}

func NewCompanyPositionService(repository *repository.RepositoryContainer) CompanyPositionService {
	return &companyPositionService{
		repository: repository,
	}
}

func (s *companyPositionService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *companyPositionService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.CompanyPositionEntities, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
             			(c:Contact {id:$id})-[r:WORKS_AT]->(co:Company)
				RETURN co, r`,
			map[string]interface{}{
				"id":     contact.ID,
				"tenant": common.GetContext(ctx).Tenant})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	companyPositionEntities := entity.CompanyPositionEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		companyPositionEntity := s.mapDbNodeToEntity(dbRecord.Values[0].(dbtype.Node))
		s.addDbRelationshipToEntity(dbRecord.Values[1].(dbtype.Relationship), companyPositionEntity)
		companyPositionEntities = append(companyPositionEntities, *companyPositionEntity)
	}

	return &companyPositionEntities, nil
}

func setCompanyPositionRelationshipInTx(ctx context.Context, contactId string, input entity.CompanyPositionEntity, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
            MERGE (co:Company {name: $company})-[:COMPANY_BELONGS_TO_TENANT]->(t)
			MERGE (c)-[:WORKS_AT {jobTitle:$jobTitle}]->(co)`,
		map[string]interface{}{
			"tenant":    common.GetContext(ctx).Tenant,
			"contactId": contactId,
			"company":   input.Company,
			"jobTitle":  input.JobTitle,
		})
	return err
}

func (s *companyPositionService) mapDbNodeToEntity(node dbtype.Node) *entity.CompanyPositionEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.CompanyPositionEntity{
		Company: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &result
}

func (s *companyPositionService) addDbRelationshipToEntity(relationship dbtype.Relationship, companyPositionEntity *entity.CompanyPositionEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	companyPositionEntity.JobTitle = utils.GetStringPropOrEmpty(props, "jobTitle")
}
