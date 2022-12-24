package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
)

type RoleRepository interface {
	MergeRole(tenant, contactId, companyExternalId, externalSystemId string) error
	MergePrimaryRole(tenant, contactId, jobTitle, companyExternalId, externalSystemId string) error
	RemoveOutdatedRoles(tenant, contactId, externalSystemId string, companiesExternalIds []string) error
}

type roleRepository struct {
	driver *neo4j.Driver
}

func NewRoleRepository(driver *neo4j.Driver) RoleRepository {
	return &roleRepository{
		driver: driver,
	}
}

func (r *roleRepository) MergeRole(tenant, contactId, companyExternalId, externalSystemId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MATCH (e:ExternalSystem {id:$externalSystemId})<-[:IS_LINKED_WITH {externalId:$companyExternalId}]-(co:Company)-[:COMPANY_BELONGS_TO_TENANT]->(t)
			MERGE (c)-[:HAS_ROLE]->(r:Role)-[:WORKS]->(co)
            ON CREATE SET r.primary=false, r.id=randomUUID()
			RETURN r`,
			map[string]interface{}{
				"tenant":            tenant,
				"contactId":         contactId,
				"externalSystemId":  externalSystemId,
				"companyExternalId": companyExternalId,
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single()
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *roleRepository) MergePrimaryRole(tenant, contactId, jobTitle, companyExternalId, externalSystemId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[:HAS_ROLE]->(r:Role)
			SET r.primary=false
			REMOVE r.jobTitle
			WITH distinct c
			MATCH (c)-[:HAS_ROLE]->(r:Role)-[:WORKS]->(co:Company)-[:IS_LINKED_WITH {externalId:$companyExternalId}]->(e:ExternalSystem {id:$externalSystemId})
			SET r.primary=true, r.jobTitle=$jobTitle
			return r`,
			map[string]interface{}{
				"tenant":            tenant,
				"contactId":         contactId,
				"jobTitle":          jobTitle,
				"externalSystemId":  externalSystemId,
				"companyExternalId": companyExternalId,
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single()
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *roleRepository) RemoveOutdatedRoles(tenant, contactId, externalSystemId string, companiesExternalIds []string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MATCH (c)-[:HAS_ROLE]->(r:Role)-[:WORKS]->(:Company)-[er:IS_LINKED_WITH]->(e:ExternalSystem {id:$externalSystemId})
			WHERE NOT er.externalId in $companiesExternalIds
			DETACH DELETE r`,
			map[string]interface{}{
				"tenant":               tenant,
				"contactId":            contactId,
				"externalSystemId":     externalSystemId,
				"companiesExternalIds": companiesExternalIds,
			})
		return nil, err
	})
	return err
}
