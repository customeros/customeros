package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type CompanyRepository interface {
	LinkNewCompanyToContactInTx(tx neo4j.Transaction, tenant, contactId, companyName, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error)
	LinkExistingCompanyToContactInTx(tx neo4j.Transaction, tenant, contactId, companyId, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error)
	UpdateCompanyPositionInTx(tx neo4j.Transaction, tenant, contactId, companyPositionId, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error)
	DeleteCompanyPositionInTx(tx neo4j.Transaction, tenant, contactId, companyPositionId string) error
	GetCompanyPositionsForContact(session neo4j.Session, tenant, contactId string) ([]*CompanyWithPositionNode, error)
	GetCompanyPositionForContact(session neo4j.Session, tenant, contactId, companyPositionId string) (*CompanyWithPositionNode, error)

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

type CompanyWithPositionNode struct {
	Company  *dbtype.Node
	Position *dbtype.Relationship
}

func (r *companyRepository) LinkNewCompanyToContactInTx(tx neo4j.Transaction, tenant, contactId, companyName, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error) {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MERGE (co:Company {id:randomUUID(), name: $companyName})-[:COMPANY_BELONGS_TO_TENANT]->(t)
			MERGE (c)-[r:WORKS_AT {id:randomUUID(), jobTitle:$jobTitle}]->(co)
			RETURN co, r`,
		map[string]any{
			"tenant":      tenant,
			"contactId":   contactId,
			"companyName": companyName,
			"jobTitle":    jobTitle,
		})
	if err != nil {
		return nil, nil, err
	}
	dbRecord, err := queryResult.Single()
	return utils.NodePtr(dbRecord.Values[0].(dbtype.Node)), utils.RelationshipPtr(dbRecord.Values[1].(dbtype.Relationship)), err
}

func (r *companyRepository) LinkExistingCompanyToContactInTx(tx neo4j.Transaction, tenant, contactId, companyId, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error) {
	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
				  (co:Company {id:$companyId})-[:COMPANY_BELONGS_TO_TENANT]->(t)
			MERGE (c)-[r:WORKS_AT {id:randomUUID(), jobTitle:$jobTitle}]->(co)
			RETURN co, r`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
			"companyId": companyId,
			"jobTitle":  jobTitle,
		})
	if err != nil {
		return nil, nil, err
	}
	dbRecord, err := queryResult.Single()

	return utils.NodePtr(dbRecord.Values[0].(dbtype.Node)), utils.RelationshipPtr(dbRecord.Values[1].(dbtype.Relationship)), err
}

func (r *companyRepository) UpdateCompanyPositionInTx(tx neo4j.Transaction, tenant, contactId, companyPositionId, jobTitle string) (*dbtype.Node, *dbtype.Relationship, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
				  (c)-[r:WORKS_AT {id:$companyPositionId}]->(co:Company)
			SET r.jobTitle=$jobTitle
			RETURN co, r`,
		map[string]any{
			"tenant":            tenant,
			"contactId":         contactId,
			"companyPositionId": companyPositionId,
			"jobTitle":          jobTitle,
		})
	if err != nil {
		return nil, nil, err
	}
	dbRecord, err := queryResult.Single()
	return utils.NodePtr(dbRecord.Values[0].(dbtype.Node)), utils.RelationshipPtr(dbRecord.Values[1].(dbtype.Relationship)), err
}

func (r *companyRepository) DeleteCompanyPositionInTx(tx neo4j.Transaction, tenant, contactId, companyPositionId string) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)-[r:WORKS_AT {id:$companyPositionId}]->(co:Company)
			DELETE r`,
		map[string]any{
			"tenant":            tenant,
			"contactId":         contactId,
			"companyPositionId": companyPositionId,
		})
	return err
}

func (r *companyRepository) GetCompanyPositionsForContact(session neo4j.Session, tenant, contactId string) ([]*CompanyWithPositionNode, error) {
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
		companyWithPositionNodes := []*CompanyWithPositionNode{}
		for _, v := range dbRecords.([]*neo4j.Record) {
			singleCompanyWithPositionNodes := new(CompanyWithPositionNode)
			singleCompanyWithPositionNodes.Company = utils.NodePtr(v.Values[0].(dbtype.Node))
			singleCompanyWithPositionNodes.Position = utils.RelationshipPtr(v.Values[1].(dbtype.Relationship))
			companyWithPositionNodes = append(companyWithPositionNodes, singleCompanyWithPositionNodes)
		}
		return companyWithPositionNodes, nil
	}
}

func (r *companyRepository) GetCompanyPositionForContact(session neo4j.Session, tenant, contactId, companyPositionId string) (*CompanyWithPositionNode, error) {
	dbRecord, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)-[r:WORKS_AT {id:$companyPositionId}]->(co:Company)
			RETURN co, r ORDER BY co.name, r.jobTitle`,
			map[string]any{
				"tenant":            tenant,
				"contactId":         contactId,
				"companyPositionId": companyPositionId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Single()
		}
	})
	if err != nil {
		return nil, err
	}
	singleCompanyWithPositionNodes := new(CompanyWithPositionNode)
	singleCompanyWithPositionNodes.Company = utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node))
	singleCompanyWithPositionNodes.Position = utils.RelationshipPtr(dbRecord.(*db.Record).Values[1].(dbtype.Relationship))
	return singleCompanyWithPositionNodes, nil
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
