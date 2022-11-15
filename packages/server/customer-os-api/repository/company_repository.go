package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type CompanyRepository interface {
	LinkNewCompanyToContact(tenant, contactId, companyName, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error)
	LinkExistingCompanyToContact(tenant, contactId, companyId, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error)
}

type companyRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewCompanyRepository(driver *neo4j.Driver, repos *RepositoryContainer) CompanyRepository {
	return &companyRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *companyRepository) LinkNewCompanyToContact(tenant, contactId, companyName, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error) {
	session := (*r.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MERGE (co:Company {id:randomUUID(), name: $companyName})-[:COMPANY_BELONGS_TO_TENANT]->(t)
			MERGE (c)-[r:WORKS_AT {id:randomUUID(), jobTitle:$jobTitle}]->(co)
			RETURN co, r`,
			map[string]any{
				"tenant":      tenant,
				"contactId":   contactId,
				"companyName": companyName,
				"jobTitle":    jobTitle,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Single()
		}
	})
	if err != nil {
		return nil, nil, err
	}
	return utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node)), utils.RelationshipPtr(dbRecord.(*db.Record).Values[1].(dbtype.Relationship)), err
}

func (r *companyRepository) LinkExistingCompanyToContact(tenant, contactId, companyId, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error) {
	session := (*r.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
				  (co:Company {id:$companyId})-[:COMPANY_BELONGS_TO_TENANT]->(t)
			MERGE (c)-[r:WORKS_AT {id:randomUUID(), jobTitle:$jobTitle}]->(co)
			RETURN co, r`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"companyId": companyId,
				"jobTitle":  jobTitle,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Single()
		}
	})
	if err != nil {
		return nil, nil, err
	}
	return utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node)), utils.RelationshipPtr(dbRecord.(*db.Record).Values[1].(dbtype.Relationship)), err
}
