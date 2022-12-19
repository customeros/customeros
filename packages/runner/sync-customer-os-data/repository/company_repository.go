package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type CompanyRepository interface {
	MergeCompany(tenant string, syncDate time.Time, contact entity.CompanyData) (string, error)
}

type companyRepository struct {
	driver *neo4j.Driver
}

func NewCompanyRepository(driver *neo4j.Driver) CompanyRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) MergeCompany(tenant string, syncDate time.Time, company entity.CompanyData) (string, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				MERGE (c:Company)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e)
				ON CREATE SET r.externalId=$externalId, c.id=randomUUID(), c.createdAt=$createdAt,
								c.name=$name, c.description=$description, r.syncDate=$syncDate, c.readonly=$readonly,
								c.domain=$domain, c.website=$website, c.industry=$industry, c.isPublic=$isPublic
				ON MATCH SET 	c.name=$name, c.description=$description, r.syncDate=$syncDate, c.readonly=$readonly,
								c.domain=$domain, c.website=$website, c.industry=$industry, c.isPublic=$isPublic
				WITH c, t
				MERGE (c)-[:COMPANY_BELONGS_TO_TENANT]->(t)
				RETURN c.id`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": company.ExternalSystem,
				"externalId":     company.ExternalId,
				"syncDate":       syncDate,
				"name":           company.Name,
				"description":    company.Description,
				"readonly":       company.Readonly,
				"createdAt":      company.CreatedAt,
				"domain":         company.Domain,
				"website":        company.Website,
				"industry":       company.Industry,
				"isPublic":       company.IsPublic,
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return "", err
	}
	return dbRecord.(string), nil
}
