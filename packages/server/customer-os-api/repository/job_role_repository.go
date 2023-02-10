package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type JobRoleRepository interface {
	GetJobRolesForContact(session neo4j.Session, tenant, contactId string) ([]*dbtype.Node, error)
	GetJobRolesForOrganization(session neo4j.Session, tenant, organizationId string) ([]*dbtype.Node, error)
	DeleteJobRoleInTx(tx neo4j.Transaction, tenant, contactId, roleId string) error
	SetOtherJobRolesForContactNonPrimaryInTx(tx neo4j.Transaction, tenant, contactId, skipRoleId string) error
	CreateJobRole(tx neo4j.Transaction, tenant, contactId string, input entity.JobRoleEntity) (*dbtype.Node, error)
	UpdateJobRoleDetails(tx neo4j.Transaction, tenant, contactId, roleId string, input entity.JobRoleEntity) (*dbtype.Node, error)
	LinkWithOrganization(tx neo4j.Transaction, tenant, roleId, organizationId string) error
}

type jobRoleRepository struct {
	driver *neo4j.Driver
}

func NewJobRoleRepository(driver *neo4j.Driver) JobRoleRepository {
	return &jobRoleRepository{
		driver: driver,
	}
}

func (r *jobRoleRepository) GetJobRolesForContact(session neo4j.Session, tenant, contactId string) ([]*dbtype.Node, error) {
	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[:WORKS_AS]->(r:JobRole) 
				RETURN r ORDER BY r.jobTitle`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}

func (r *jobRoleRepository) GetJobRolesForOrganization(session neo4j.Session, tenant, organizationId string) ([]*dbtype.Node, error) {
	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(org)<-[:ROLE_IN]-(r:JobRole) 
				RETURN r ORDER BY r.jobTitle`,
			map[string]interface{}{
				"organizationId": organizationId,
				"tenant":         tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}

func (r *jobRoleRepository) CreateJobRole(tx neo4j.Transaction, tenant string, contactId string, input entity.JobRoleEntity) (*dbtype.Node, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[:WORKS_AS]->(r:JobRole) " +
		" ON CREATE SET r.id=randomUUID(), " +
		"				r.jobTitle=$jobTitle, " +
		"				r.primary=$primary, " +
		"				r.responsibilityLevel=$responsibilityLevel, " +
		"				r.source=$source, " +
		"				r.sourceOfTruth=$sourceOfTruth, " +
		"				r.appSource=$appSource, " +
		"				r.createdAt=$now, " +
		"				r.updatedAt=$now, " +
		"				r:%s " +
		" RETURN r"

	if queryResult, err := tx.Run(fmt.Sprintf(query, "JobRole_"+tenant),
		map[string]interface{}{
			"tenant":              tenant,
			"contactId":           contactId,
			"jobTitle":            input.JobTitle,
			"primary":             input.Primary,
			"responsibilityLevel": input.ResponsibilityLevel,
			"source":              input.Source,
			"sourceOfTruth":       input.SourceOfTruth,
			"appSource":           input.AppSource,
			"now":                 utils.Now(),
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}
}

func (r *jobRoleRepository) UpdateJobRoleDetails(tx neo4j.Transaction, tenant, contactId, roleId string, input entity.JobRoleEntity) (*dbtype.Node, error) {
	if queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[:WORKS_AS]->(r:JobRole {id:$roleId})
			SET r.jobTitle=$jobTitle, 
				r.primary=$primary,
				r.responsibilityLevel=$responsibilityLevel,
				r.sourceOfTruth=$sourceOfTruth,
				r.updatedAt=datetime({timezone: 'UTC'})
			RETURN r`,
		map[string]interface{}{
			"tenant":              tenant,
			"contactId":           contactId,
			"roleId":              roleId,
			"jobTitle":            input.JobTitle,
			"primary":             input.Primary,
			"responsibilityLevel": input.ResponsibilityLevel,
			"sourceOfTruth":       input.SourceOfTruth,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}
}

func (r *jobRoleRepository) DeleteJobRoleInTx(tx neo4j.Transaction, tenant, contactId, roleId string) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)-[:WORKS_AS]->(r:JobRole {id:$roleId})
			DETACH DELETE r`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
			"roleId":    roleId,
		})
	return err
}

func (r *jobRoleRepository) SetOtherJobRolesForContactNonPrimaryInTx(tx neo4j.Transaction, tenant, contactId, skipRoleId string) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[:WORKS_AS]->(r:JobRole)
			WHERE r.id <> $skipRoleId
            SET r.primary=false,
				r.updatedAt=datetime({timezone: 'UTC'})`,
		map[string]interface{}{
			"tenant":     tenant,
			"contactId":  contactId,
			"skipRoleId": skipRoleId,
		})
	return err
}

func (r *jobRoleRepository) LinkWithOrganization(tx neo4j.Transaction, tenant string, roleId string, organizationId string) error {
	_, err := tx.Run(`
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(r:JobRole {id:$roleId})<-[:WORKS_AS]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (r)-[rel:ROLE_IN]->(org2:Organization)
				WHERE org2.id <> org.id
			DELETE rel
			WITH r, org
			MERGE (r)-[:ROLE_IN]->(org)
			`,
		map[string]interface{}{
			"tenant":         tenant,
			"roleId":         roleId,
			"organizationId": organizationId,
		})
	return err
}
