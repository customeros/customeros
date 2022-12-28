package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type CompanyRepository interface {
	MergeCompany(tenant string, syncDate time.Time, company entity.CompanyData) (string, error)
	MergeCompanyAddress(tenant, companyId string, company entity.CompanyData) error
}

type companyRepository struct {
	driver *neo4j.Driver
}

func (r *companyRepository) MergeCompanyAddress(tenant, companyId string, company entity.CompanyData) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (co:Company {id:$companyId})-[:COMPANY_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (co)-[:LOCATED_AT]->(a:Address {source:$source}) " +
		" ON CREATE SET a.id=randomUUID(), a.source=$source, " +
		"	a.country=$country, a.state=$state, a.city=$city, a.address=$address, " +
		"	a.address2=$address2, a.zip=$zip, a.phone=$phone, a:%s " +
		" ON MATCH SET 	a.country=$country, a.state=$state, a.city=$city, a.address=$address, " +
		"	a.address2=$address2, a.zip=$zip, a.phone=$phone"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Address_"+tenant),
			map[string]interface{}{
				"tenant":    tenant,
				"companyId": companyId,
				"source":    company.ExternalSystem,
				"country":   company.Country,
				"state":     company.State,
				"city":      company.City,
				"address":   company.Address,
				"address2":  company.Address2,
				"zip":       company.Zip,
				"phone":     company.Phone,
			})
		return nil, err
	})
	return err
}

func NewCompanyRepository(driver *neo4j.Driver) CompanyRepository {
	return &companyRepository{
		driver: driver,
	}
}

func (r *companyRepository) MergeCompany(tenant string, syncDate time.Time, company entity.CompanyData) (string, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MERGE (c:Company)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET r.externalId=$externalId, c.id=randomUUID(), c.createdAt=$createdAt, " +
		"               c.name=$name, c.description=$description, r.syncDate=$syncDate, c.readonly=$readonly, " +
		"               c.domain=$domain, c.website=$website, c.industry=$industry, c.isPublic=$isPublic, c:%s " +
		" ON MATCH SET c.name=$name, c.description=$description, r.syncDate=$syncDate, c.readonly=$readonly, " +
		"              c.domain=$domain, c.website=$website, c.industry=$industry, c.isPublic=$isPublic " +
		" WITH c, t " +
		" MERGE (c)-[:COMPANY_BELONGS_TO_TENANT]->(t) " +
		" RETURN c.id"

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "Company_"+tenant),
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
