package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type RoleRepository interface {
	MergeRole(tenant, contactId, organizationExternalId, externalSystemId string) error
	MergePrimaryRole(tenant, contactId, jobTitle, organizationExternalId, externalSystemId string) error
	RemoveOutdatedRoles(tenant, contactId, externalSystemId string, organizationsExternalIds []string) error
}

type roleRepository struct {
	driver *neo4j.Driver
}

func NewRoleRepository(driver *neo4j.Driver) RoleRepository {
	return &roleRepository{
		driver: driver,
	}
}

func (r *roleRepository) MergeRole(tenant, contactId, organizationExternalId, externalSystemId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MATCH (e:ExternalSystem {id:$externalSystemId})<-[:IS_LINKED_WITH {externalId:$organizationExternalId}]-(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t) " +
		" MERGE (c)-[:HAS_ROLE]->(r:Role)-[:WORKS]->(org) " +
		" ON CREATE SET r.primary=false, r.id=randomUUID(), r.source=$source, r.sourceOfTruth=$sourceOfTruth, r.appSource=$appSource, " +
		"				r.createdAt=$now, r.updatedAt=$now, r:%s " +
		" RETURN r"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "Role_"+tenant),
			map[string]interface{}{
				"tenant":                 tenant,
				"contactId":              contactId,
				"externalSystemId":       externalSystemId,
				"organizationExternalId": organizationExternalId,
				"source":                 externalSystemId,
				"sourceOfTruth":          externalSystemId,
				"appSource":              externalSystemId,
				"now":                    time.Now().UTC(),
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

func (r *roleRepository) MergePrimaryRole(tenant, contactId, jobTitle, organizationExternalId, externalSystemId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[:HAS_ROLE]->(r:Role)
			SET r.primary=false
			REMOVE r.jobTitle
			WITH distinct c
			MATCH (c)-[:HAS_ROLE]->(r:Role)-[:WORKS]->(org:Organization)-[:IS_LINKED_WITH {externalId:$organizationExternalId}]->(e:ExternalSystem {id:$externalSystemId})
			SET r.primary=true, r.jobTitle=$jobTitle, r.updatedAt=$now
			return r`,
			map[string]interface{}{
				"tenant":                 tenant,
				"contactId":              contactId,
				"jobTitle":               jobTitle,
				"externalSystemId":       externalSystemId,
				"organizationExternalId": organizationExternalId,
				"now":                    time.Now().UTC(),
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

func (r *roleRepository) RemoveOutdatedRoles(tenant, contactId, externalSystemId string, organizationsExternalIds []string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MATCH (c)-[:HAS_ROLE]->(r:Role)-[:WORKS]->(:Organization)-[er:IS_LINKED_WITH]->(e:ExternalSystem {id:$externalSystemId})
			WHERE NOT er.externalId in $organizationsExternalIds
			DETACH DELETE r`,
			map[string]interface{}{
				"tenant":                   tenant,
				"contactId":                contactId,
				"externalSystemId":         externalSystemId,
				"organizationsExternalIds": organizationsExternalIds,
			})
		return nil, err
	})
	return err
}
