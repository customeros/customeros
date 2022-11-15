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
	UpdateCompanyPosition(tenant, contactId, companyPositionId, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error)
	DeleteCompanyPosition(tenant, contactId, companyPositionId string) error
	GetCompanyPositionsForContact(tenant, contactId string) ([]*CompanyWithPositionNodes, error)

	GetPaginatedCompaniesWithNameLike(tenant, companyName string, skip, limit int) (*utils.DbNodesWithTotalCount, error)
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

type CompanyWithPositionNodes struct {
	Company  *dbtype.Node
	Position *dbtype.Relationship
}

func (r *companyRepository) LinkNewCompanyToContact(tenant, contactId, companyName, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
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
	session := utils.NewNeo4jWriteSession(*r.driver)
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

func (r *companyRepository) UpdateCompanyPosition(tenant, contactId, companyPositionId, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
				  (c)-[r:WORKS_AT {id:$companyPositionId}]->(co:Company)
			SET r.jobTitle=$jobTitle
			RETURN co, r`,
			map[string]any{
				"tenant":            tenant,
				"contactId":         contactId,
				"companyPositionId": companyPositionId,
				"jobTitle":          jobTitle,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Single()
		}
	})
	if err != nil {
		return nil, nil, err
	}
	return utils.NodePtr(dbRecord.(*neo4j.Record).Values[0].(dbtype.Node)), utils.RelationshipPtr(dbRecord.(*neo4j.Record).Values[1].(dbtype.Relationship)), err
}

func (r *companyRepository) DeleteCompanyPosition(tenant, contactId, companyPositionId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)-[r:WORKS_AT {id:$companyPositionId}]->(co:Company)
			DELETE r`,
			map[string]any{
				"tenant":            tenant,
				"contactId":         contactId,
				"companyPositionId": companyPositionId,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *companyRepository) GetCompanyPositionsForContact(tenant, contactId string) ([]*CompanyWithPositionNodes, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)-[r:WORKS_AT]->(co:Company)
			RETURN co, r ORDER BY co.name, r.jobTitle`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	} else if len(dbRecords.([]*neo4j.Record)) == 0 {
		return nil, nil
	} else {
		companyWithPositionNodes := []*CompanyWithPositionNodes{}
		for _, v := range dbRecords.([]*neo4j.Record) {
			singleCompanyWithPositionNodes := new(CompanyWithPositionNodes)
			singleCompanyWithPositionNodes.Company = utils.NodePtr(v.Values[0].(dbtype.Node))
			singleCompanyWithPositionNodes.Position = utils.RelationshipPtr(v.Values[1].(dbtype.Relationship))
			companyWithPositionNodes = append(companyWithPositionNodes, singleCompanyWithPositionNodes)
		}
		return companyWithPositionNodes, nil
	}
}

func (r *companyRepository) GetPaginatedCompaniesWithNameLike(tenant, companyName string, skip, limit int) (*utils.DbNodesWithTotalCount, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[:COMPANY_BELONGS_TO_TENANT]-(co:Company) 
				WHERE toLower(co.name) CONTAINS toLower($companyName)
				RETURN count(co) as count`,
			map[string]interface{}{
				"tenant":      tenant,
				"companyName": companyName,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[:COMPANY_BELONGS_TO_TENANT]-(co:Company) 
                WHERE toLower(co.name) CONTAINS toLower($companyName)
				RETURN co ORDER BY co.name SKIP $skip LIMIT $limit`,
			map[string]interface{}{
				"tenant":      tenant,
				"companyName": companyName,
				"skip":        skip,
				"limit":       limit,
			})
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}

	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}

	return dbNodesWithTotalCount, nil
}
