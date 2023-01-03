package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Companies_FilterByNameLike(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateCompany(driver, tenantName, "A closed company")
	neo4jt.CreateCompany(driver, tenantName, "OPENLINE")
	neo4jt.CreateCompany(driver, tenantName, "the openline")
	neo4jt.CreateCompany(driver, tenantName, "some other open company")
	neo4jt.CreateCompany(driver, tenantName, "OpEnLiNe")

	require.Equal(t, 5, neo4jt.GetCountOfNodes(driver, "Company"))

	rawResponse, err := c.RawPost(getQuery("get_companies"),
		client.Var("page", 1),
		client.Var("limit", 3),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var companies struct {
		Companies model.CompanyPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &companies)
	require.Nil(t, err)
	require.NotNil(t, companies)
	pagedCompanies := companies.Companies
	require.Equal(t, 2, pagedCompanies.TotalPages)
	require.Equal(t, int64(4), pagedCompanies.TotalElements)
	require.Equal(t, "OPENLINE", pagedCompanies.Content[0].Name)
	require.Equal(t, "OpEnLiNe", pagedCompanies.Content[1].Name)
	require.Equal(t, "some other open company", pagedCompanies.Content[2].Name)
}

func TestQueryResolver_Company(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	companyInput := entity.CompanyEntity{
		Name:        "Company name",
		Description: "Company description",
		Domain:      "Company domain",
		Website:     "Company_website.com",
		Industry:    "tech",
		IsPublic:    true,
	}
	companyId1 := neo4jt.CreateFullCompany(driver, tenantName, companyInput)
	neo4jt.CreateCompany(driver, tenantName, "otherCompany")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Company"))

	rawResponse, err := c.RawPost(getQuery("get_company_by_id"),
		client.Var("companyId", companyId1),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var company struct {
		Company model.Company
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &company)
	require.Nil(t, err)
	require.NotNil(t, company)
	require.Equal(t, companyId1, company.Company.ID)
	require.Equal(t, companyInput.Name, company.Company.Name)
	require.Equal(t, companyInput.Description, *company.Company.Description)
	require.Equal(t, companyInput.Domain, *company.Company.Domain)
	require.Equal(t, companyInput.Website, *company.Company.Website)
	require.Equal(t, companyInput.IsPublic, *company.Company.IsPublic)
	require.Equal(t, companyInput.Industry, *company.Company.Industry)
	require.NotNil(t, company.Company.CreatedAt)
}

func TestQueryResolver_Companies_WithAddresses(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	company1 := neo4jt.CreateCompany(driver, tenantName, "OPENLINE")
	company2 := neo4jt.CreateCompany(driver, tenantName, "some other company")
	addressInput := entity.AddressEntity{
		Source:   "hubspot",
		Country:  "testCountry",
		State:    "testState",
		City:     "testCity",
		Address:  "testAddress",
		Address2: "testAddress2",
		Zip:      "testZip",
		Phone:    "testPhone",
		Fax:      "testFax",
	}
	address1 := neo4jt.CreateAddress(driver, addressInput)
	address2 := neo4jt.CreateAddress(driver, entity.AddressEntity{
		Source: "manual",
	})
	neo4jt.CompanyHasAddress(driver, company1, address1)
	neo4jt.CompanyHasAddress(driver, company2, address2)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Company"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Address"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "LOCATED_AT"))

	rawResponse, err := c.RawPost(getQuery("get_companies_with_addresses"),
		client.Var("page", 1),
		client.Var("limit", 3),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var companies struct {
		Companies model.CompanyPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &companies)
	require.Nil(t, err)
	require.NotNil(t, companies)
	pagedCompanies := companies.Companies
	require.Equal(t, int64(1), pagedCompanies.TotalElements)
	require.Equal(t, 1, len(companies.Companies.Content[0].Addresses))
	address := companies.Companies.Content[0].Addresses[0]
	require.Equal(t, address1, address.ID)
	require.Equal(t, addressInput.Source, *address.Source)
	require.Equal(t, addressInput.Country, *address.Country)
	require.Equal(t, addressInput.City, *address.City)
	require.Equal(t, addressInput.State, *address.State)
	require.Equal(t, addressInput.Address, *address.Address)
	require.Equal(t, addressInput.Address2, *address.Address2)
	require.Equal(t, addressInput.Fax, *address.Fax)
	require.Equal(t, addressInput.Phone, *address.Phone)
	require.Equal(t, addressInput.Zip, *address.Zip)
}

func TestMutationResolver_CompanyCreate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetTotalCountOfNodes(driver))

	rawResponse, err := c.RawPost(getQuery("create_company"))
	assertRawResponseSuccess(t, rawResponse, err)

	var company struct {
		Company_Create model.Company
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &company)
	require.Nil(t, err)
	require.NotNil(t, company)
	require.NotNil(t, company.Company_Create.ID)
	require.NotNil(t, company.Company_Create.CreatedAt)
	require.Equal(t, "company name", company.Company_Create.Name)
	require.Equal(t, "company description", *company.Company_Create.Description)
	require.Equal(t, "company domain", *company.Company_Create.Domain)
	require.Equal(t, "company website", *company.Company_Create.Website)
	require.Equal(t, "company industry", *company.Company_Create.Industry)
	require.Equal(t, true, *company.Company_Create.IsPublic)
	require.Equal(t, false, *company.Company_Create.Readonly)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Company"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Company_"+tenantName))

	assertNeo4jLabels(t, driver, []string{"Tenant", "Company", "Company_" + tenantName})
}

func TestMutationResolver_CompanyUpdate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	companyId := neo4jt.CreateCompany(driver, tenantName, "some company")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Company"))

	rawResponse, err := c.RawPost(getQuery("update_company"),
		client.Var("companyId", companyId))
	assertRawResponseSuccess(t, rawResponse, err)

	var company struct {
		Company_Update model.Company
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &company)
	require.Nil(t, err)
	require.NotNil(t, company)
	require.Equal(t, companyId, company.Company_Update.ID)
	require.NotNil(t, company.Company_Update.CreatedAt)
	require.Equal(t, "updated name", company.Company_Update.Name)
	require.Equal(t, "updated description", *company.Company_Update.Description)
	require.Equal(t, "updated domain", *company.Company_Update.Domain)
	require.Equal(t, "updated website", *company.Company_Update.Website)
	require.Equal(t, "updated industry", *company.Company_Update.Industry)
	require.Equal(t, true, *company.Company_Update.IsPublic)

	// Check still single company node exists after update, no new node created
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Company"))
}

func TestMutationResolver_CompanyDelete(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	companyId := neo4jt.CreateCompany(driver, tenantName, "LLC LLC")
	addressId := neo4jt.CreateAddress(driver, entity.AddressEntity{
		Source: "manual",
	})
	neo4jt.CompanyHasAddress(driver, companyId, addressId)

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "LOCATED_AT"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Address"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Company"))

	rawResponse, err := c.RawPost(getQuery("delete_company"),
		client.Var("companyId", companyId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		Company_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.Company_Delete.Result)

	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "LOCATED_AT"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Address"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Company"))

	assertNeo4jLabels(t, driver, []string{"Tenant"})
}
