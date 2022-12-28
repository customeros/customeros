package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type ContactRoleRepository interface {
	GetRolesForContact(session neo4j.Session, tenant, contactId string) ([]*dbtype.Node, error)
	DeleteContactRoleInTx(tx neo4j.Transaction, tenant, contactId, roleId string) error
	SetOtherRolesNonPrimaryInTx(tx neo4j.Transaction, tenant, contactId, skipRoleId string) error
	CreateContactRole(tx neo4j.Transaction, tenant, contactId string, input entity.ContactRoleEntity) (*dbtype.Node, error)
	UpdateContactRoleDetails(tx neo4j.Transaction, tenant, contactId, roleId string, input entity.ContactRoleEntity) (*dbtype.Node, error)
	LinkWithCompany(tx neo4j.Transaction, tenant, roleId, companyId string) error
}

type contactRoleRepository struct {
	driver *neo4j.Driver
}

func NewContactRoleRepository(driver *neo4j.Driver) ContactRoleRepository {
	return &contactRoleRepository{
		driver: driver,
	}
}

func (r *contactRoleRepository) GetRolesForContact(session neo4j.Session, tenant, contactId string) ([]*dbtype.Node, error) {
	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[:HAS_ROLE]->(r:Role) 
				RETURN r`,
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

func (r *contactRoleRepository) CreateContactRole(tx neo4j.Transaction, tenant string, contactId string, input entity.ContactRoleEntity) (*dbtype.Node, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[:HAS_ROLE]->(r:Role) " +
		" ON CREATE SET r.id=randomUUID(), r.jobTitle=$jobTitle, r.primary=$primary, r:%s " +
		" RETURN r"

	if queryResult, err := tx.Run(fmt.Sprintf(query, "Role_"+tenant),
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"jobTitle":  input.JobTitle,
			"primary":   input.Primary,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}
}

func (r *contactRoleRepository) UpdateContactRoleDetails(tx neo4j.Transaction, tenant, contactId, roleId string, input entity.ContactRoleEntity) (*dbtype.Node, error) {
	if queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[:HAS_ROLE]->(r:Role {id:$roleId})
			SET r.jobTitle=$jobTitle, r.primary=$primary
			RETURN r`,
		map[string]interface{}{
			"tenant":    tenant,
			"contactId": contactId,
			"roleId":    roleId,
			"jobTitle":  input.JobTitle,
			"primary":   input.Primary,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}
}

func (r *contactRoleRepository) DeleteContactRoleInTx(tx neo4j.Transaction, tenant, contactId, roleId string) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)-[:HAS_ROLE]->(r:Role {id:$roleId})
			DETACH DELETE r`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
			"roleId":    roleId,
		})
	return err
}

func (r *contactRoleRepository) SetOtherRolesNonPrimaryInTx(tx neo4j.Transaction, tenant, contactId, skipRoleId string) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[:HAS_ROLE]->(r:Role)
			WHERE r.id <> $skipRoleId
            SET r.primary=false`,
		map[string]interface{}{
			"tenant":     tenant,
			"contactId":  contactId,
			"skipRoleId": skipRoleId,
		})
	return err
}

func (r *contactRoleRepository) LinkWithCompany(tx neo4j.Transaction, tenant string, roleId string, companyId string) error {
	_, err := tx.Run(`
			MATCH (co:Company {id:$companyId})-[:COMPANY_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(r:Role {id:$roleId})<-[:HAS_ROLE]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (r)-[rel:WORKS]->(co2:Company)
				WHERE co2.id <> co.id
			DELETE rel
			WITH r, co
			MERGE (r)-[:WORKS]->(co)
			`,
		map[string]interface{}{
			"tenant":    tenant,
			"roleId":    roleId,
			"companyId": companyId,
		})
	return err
}
