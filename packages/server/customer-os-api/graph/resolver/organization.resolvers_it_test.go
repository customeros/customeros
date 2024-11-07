package resolver

import (
	"context"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/stretchr/testify/require"
)

func TestQueryResolver_Organizations_FilterByNameLike(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "A closed organization"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "OPENLINE"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "the openline"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "some other open organization"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "OpEnLiNe"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organizations"),
		client.Var("page", 1),
		client.Var("limit", 3),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizations struct {
		Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizations)
	require.Nil(t, err)
	require.NotNil(t, organizations)
	pagedOrganizations := organizations.Organizations
	require.Equal(t, 2, pagedOrganizations.TotalPages)
	require.Equal(t, int64(4), pagedOrganizations.TotalElements)
	require.Equal(t, "OPENLINE", pagedOrganizations.Content[0].Name)
	require.Equal(t, "OpEnLiNe", pagedOrganizations.Content[1].Name)
	require.Equal(t, "some other open organization", pagedOrganizations.Content[2].Name)
}

func TestQueryResolver_Organizations_FilterByName_IsEmpty(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "A closed organization"})
	org2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: ""})
	org3 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: ""})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organizations_filter_is_empty"),
		client.Var("page", 1),
		client.Var("limit", 3),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizations struct {
		Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizations)
	require.Nil(t, err)
	require.NotNil(t, organizations)
	pagedOrganizations := organizations.Organizations
	require.Equal(t, 1, pagedOrganizations.TotalPages)
	require.Equal(t, int64(2), pagedOrganizations.TotalElements)
	require.ElementsMatch(t, []string{org2, org3}, []string{pagedOrganizations.Content[0].ID, pagedOrganizations.Content[1].ID})
}

func TestQueryResolver_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	now := utils.NowPtr()
	inputOrganizationEntity := neo4jentity.OrganizationEntity{
		Name:               "Organization name",
		CustomerOsId:       "C-123-ABC",
		ReferenceId:        "100/200",
		Description:        "Organization description",
		Website:            "Organization_website.com",
		Industry:           "tech",
		SubIndustry:        "tech-sub",
		IndustryGroup:      "tech-group",
		TargetAudience:     "tech-audience",
		ValueProposition:   "value-proposition",
		LastFundingRound:   "Seed",
		LastFundingAmount:  "10k",
		Note:               "Some note",
		IsPublic:           true,
		YearFounded:        utils.ToPtr(int64(2019)),
		Headquarters:       "San Francisco, CA",
		EmployeeGrowthRate: "10%",
		LogoUrl:            "https://www.openline.ai/logo.png",
		IconUrl:            "https://www.openline.ai/icon.png",
		Relationship:       neo4jenum.Customer,
		Stage:              neo4jenum.Lead,
		OnboardingDetails: neo4jentity.OnboardingDetails{
			Status:    "DONE",
			Comments:  "some comments",
			UpdatedAt: now,
		},
	}
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, inputOrganizationEntity)
	neo4jt.AddDomainToOrg(ctx, driver, organizationId, "domain1.com")
	neo4jt.AddDomainToOrg(ctx, driver, organizationId, "domain2.com")
	contact1Id := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contact2Id := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.ContactWorksForOrganization(ctx, driver, contact1Id, organizationId, "CTO", false)
	neo4jt.ContactWorksForOrganization(ctx, driver, contact2Id, organizationId, "CTO", false)
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "otherOrganization"})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	rawResponse := callGraphQL(t, "organization/get_organization_by_id", map[string]interface{}{"organizationId": organizationId})

	var organizationStruct struct {
		Organization model.Organization
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)
	require.Equal(t, organizationId, organizationStruct.Organization.Metadata.ID)
	require.Equal(t, inputOrganizationEntity.CustomerOsId, organizationStruct.Organization.CustomerOsID)
	require.Equal(t, inputOrganizationEntity.ReferenceId, *organizationStruct.Organization.CustomID)
	require.Equal(t, inputOrganizationEntity.Name, organizationStruct.Organization.Name)
	require.Equal(t, inputOrganizationEntity.Description, *organizationStruct.Organization.Description)
	require.Equal(t, []string{"domain1.com", "domain2.com"}, organizationStruct.Organization.Domains)
	require.Equal(t, inputOrganizationEntity.Website, *organizationStruct.Organization.Website)
	require.Equal(t, inputOrganizationEntity.IsPublic, *organizationStruct.Organization.Public)
	require.Equal(t, inputOrganizationEntity.Industry, *organizationStruct.Organization.Industry)
	require.Equal(t, inputOrganizationEntity.SubIndustry, *organizationStruct.Organization.SubIndustry)
	require.Equal(t, inputOrganizationEntity.IndustryGroup, *organizationStruct.Organization.IndustryGroup)
	require.Equal(t, inputOrganizationEntity.TargetAudience, *organizationStruct.Organization.TargetAudience)
	require.Equal(t, inputOrganizationEntity.ValueProposition, *organizationStruct.Organization.ValueProposition)
	require.Equal(t, model.FundingRoundSeed, *organizationStruct.Organization.LastFundingRound)
	require.Equal(t, inputOrganizationEntity.LastFundingAmount, *organizationStruct.Organization.LastFundingAmount)
	require.Equal(t, "Some note", *organizationStruct.Organization.Note)
	require.Equal(t, inputOrganizationEntity.YearFounded, organizationStruct.Organization.YearFounded)
	require.Equal(t, inputOrganizationEntity.Headquarters, *organizationStruct.Organization.Headquarters)
	require.Equal(t, inputOrganizationEntity.EmployeeGrowthRate, *organizationStruct.Organization.EmployeeGrowthRate)
	require.Equal(t, inputOrganizationEntity.LogoUrl, *organizationStruct.Organization.Logo)
	require.Equal(t, inputOrganizationEntity.IconUrl, *organizationStruct.Organization.Icon)
	require.NotNil(t, organizationStruct.Organization.CreatedAt)
	require.Equal(t, inputOrganizationEntity.OnboardingDetails.UpdatedAt, organizationStruct.Organization.AccountDetails.Onboarding.UpdatedAt)
	require.Equal(t, model.OnboardingStatusDone, organizationStruct.Organization.AccountDetails.Onboarding.Status)
	require.Equal(t, inputOrganizationEntity.OnboardingDetails.Comments, *organizationStruct.Organization.AccountDetails.Onboarding.Comments)
	require.Equal(t, model.OrganizationRelationshipCustomer, *organizationStruct.Organization.Relationship)
	require.Equal(t, model.OrganizationStageLead, *organizationStruct.Organization.Stage)
	require.Equal(t, int64(2), organizationStruct.Organization.ContactCount)
}

func TestQueryResolver_OrganizationByCustomerOsId(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Name:         "Organization name",
		CustomerOsId: "C-123-ABC",
	})
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelOrganization))

	rawResponse := callGraphQL(t, "organization/get_organization_by_customer_os_id", map[string]interface{}{"customerOsId": "C-123-ABC"})

	var organizationStruct struct {
		Organization_ByCustomerOsId model.Organization
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)
	require.Equal(t, organizationId, organizationStruct.Organization_ByCustomerOsId.ID)
	require.Equal(t, "C-123-ABC", organizationStruct.Organization_ByCustomerOsId.CustomerOsID)
	require.Equal(t, "Organization name", organizationStruct.Organization_ByCustomerOsId.Name)
}

func TestQueryResolver_OrganizationByCustomerOsId_NotFound(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Name:         "Organization name",
		CustomerOsId: "C-123-ABC",
	})
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelOrganization))

	response := callGraphQLExpectError(t, "organization/get_organization_by_customer_os_id", map[string]interface{}{"customerOsId": "C-999-JJJ"})

	require.Equal(t, "Organization not found by customerOsId C-999-JJJ", response.Message)
	require.Equal(t, "organization_ByCustomerOsId", response.Path[0])
}

func TestQueryResolver_OrganizationByCustomId(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Name:        "Organization name",
		ReferenceId: "R-123",
	})
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelOrganization))

	rawResponse := callGraphQL(t, "organization/get_organization_by_custom_id", map[string]interface{}{"customId": "R-123"})

	var organizationStruct struct {
		Organization_ByCustomId model.Organization
	}
	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)
	require.Equal(t, organizationId, organizationStruct.Organization_ByCustomId.ID)
	require.Equal(t, "Organization name", organizationStruct.Organization_ByCustomId.Name)
	require.Equal(t, "R-123", *organizationStruct.Organization_ByCustomId.CustomID)
}

func TestQueryResolver_OrganizationByCustomId_NotFound(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Name:        "Organization name",
		ReferenceId: "R-123-ABC",
	})
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelOrganization))

	response := callGraphQLExpectError(t, "organization/get_organization_by_custom_id", map[string]interface{}{"customId": "R-0000"})

	require.Equal(t, "Organization not found by customId R-0000", response.Message)
	require.Equal(t, "organization_ByCustomId", response.Path[0])
}

func TestQueryResolver_Organizations_WithLocations(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "OPENLINE")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "some other organization")
	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{
		Name:      "WORK",
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test",
		Country:   "testCountry",
		Region:    "testRegion",
		Locality:  "testLocality",
		Address:   "testAddress",
		Address2:  "testAddress2",
		Zip:       "testZip",
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{
		Name:      "UNKNOWN",
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test",
	})
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId2)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organizations_with_locations"),
		client.Var("page", 1),
		client.Var("limit", 3),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationsStruct struct {
		Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationsStruct)
	require.Nil(t, err)

	organizations := organizationsStruct.Organizations
	require.NotNil(t, organizations)
	require.Equal(t, int64(1), organizations.TotalElements)
	require.Equal(t, 2, len(organizations.Content[0].Locations))

	var locationWithAddressDtls, locationWithoutAddressDtls *model.Location
	if organizations.Content[0].Locations[0].ID == locationId1 {
		locationWithAddressDtls = organizations.Content[0].Locations[0]
		locationWithoutAddressDtls = organizations.Content[0].Locations[1]
	} else {
		locationWithAddressDtls = organizations.Content[0].Locations[1]
		locationWithoutAddressDtls = organizations.Content[0].Locations[0]
	}

	require.Equal(t, locationId1, locationWithAddressDtls.ID)
	require.Equal(t, "WORK", *locationWithAddressDtls.Name)
	require.NotNil(t, locationWithAddressDtls.CreatedAt)
	require.NotNil(t, locationWithAddressDtls.UpdatedAt)
	require.Equal(t, "test", locationWithAddressDtls.AppSource)
	require.Equal(t, model.DataSourceOpenline, locationWithAddressDtls.Source)
	require.Equal(t, "testCountry", *locationWithAddressDtls.Country)
	require.Equal(t, "testLocality", *locationWithAddressDtls.Locality)
	require.Equal(t, "testRegion", *locationWithAddressDtls.Region)
	require.Equal(t, "testAddress", *locationWithAddressDtls.Address)
	require.Equal(t, "testAddress2", *locationWithAddressDtls.Address2)
	require.Equal(t, "testZip", *locationWithAddressDtls.Zip)

	require.Equal(t, locationId2, locationWithoutAddressDtls.ID)
	require.Equal(t, "UNKNOWN", *locationWithoutAddressDtls.Name)
	require.NotNil(t, locationWithoutAddressDtls.CreatedAt)
	require.NotNil(t, locationWithoutAddressDtls.UpdatedAt)
	require.Equal(t, "test", locationWithoutAddressDtls.AppSource)
	require.Equal(t, model.DataSourceOpenline, locationWithoutAddressDtls.Source)
	require.Equal(t, "", *locationWithoutAddressDtls.Country)
	require.Equal(t, "", *locationWithoutAddressDtls.Region)
	require.Equal(t, "", *locationWithoutAddressDtls.Locality)
	require.Equal(t, "", *locationWithoutAddressDtls.Address)
	require.Equal(t, "", *locationWithoutAddressDtls.Address2)
	require.Equal(t, "", *locationWithoutAddressDtls.Zip)
}

func TestQueryResolver_Organizations_WithTags(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "Org 1 with 2 tags")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "Org 2 with 1 tag")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "Org 3 with 0 tags")
	tag1 := neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{
		Name: "tag1",
	})
	tag2 := neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{
		Name: "tag2",
	})

	neo4jt.TagOrganization(ctx, driver, organizationId1, tag1)
	neo4jt.TagOrganization(ctx, driver, organizationId1, tag2)
	neo4jt.TagOrganization(ctx, driver, organizationId2, tag1)

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "TAGGED"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organizations_with_tags"),
		client.Var("page", 1),
		client.Var("limit", 10),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationsStruct struct {
		Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationsStruct)
	require.Nil(t, err)

	organizations := organizationsStruct.Organizations
	require.NotNil(t, organizations)
	require.Equal(t, int64(3), organizations.TotalElements)
	require.Equal(t, 2, len(organizations.Content[0].Tags))
	require.ElementsMatch(t, []string{tag1, tag2},
		[]string{organizations.Content[0].Tags[0].ID, organizations.Content[0].Tags[1].ID})
	require.Equal(t, 1, len(organizations.Content[1].Tags))
	require.Equal(t, tag1, organizations.Content[1].Tags[0].ID)
	require.Equal(t, 0, len(organizations.Content[2].Tags))
}

func TestQueryResolver_Organization_WithRoles_ById(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "some organization")
	role1 := neo4jt.ContactWorksForOrganization(ctx, driver, contactId1, organizationId, "CTO", false)
	role2 := neo4jt.ContactWorksForOrganization(ctx, driver, contactId2, organizationId, "CEO", true)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_job_roles_by_id"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedOrganization struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedOrganization)
	require.Nil(t, err)
	require.Equal(t, organizationId, searchedOrganization.Organization.ID)

	roles := searchedOrganization.Organization.JobRoles
	require.Equal(t, 2, len(roles))
	var cto, ceo *model.JobRole
	ceo = roles[0]
	cto = roles[1]
	require.Equal(t, role1, cto.ID)
	require.Equal(t, "CTO", *cto.JobTitle)
	require.Equal(t, false, cto.Primary)
	require.Equal(t, contactId1, cto.Contact.ID)

	require.Equal(t, role2, ceo.ID)
	require.Equal(t, "CEO", *ceo.JobTitle)
	require.Equal(t, true, ceo.Primary)
	require.Equal(t, contactId2, ceo.Contact.ID)
}

func TestQueryResolver_Organization_WithContacts_ById(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "organization1"})
	organizationId2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "organization2"})
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId3 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId4 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId1, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId2, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId3, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId4, organizationId2)

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelContact))
	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelJobRole))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelOrganization))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_contacts_by_id"),
		client.Var("organizationId", organizationId),
		client.Var("limit", 1),
		client.Var("page", 1),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedOrganization struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedOrganization)
	require.Nil(t, err)
	require.Equal(t, organizationId, searchedOrganization.Organization.ID)
	require.Equal(t, 3, searchedOrganization.Organization.Contacts.TotalPages)
	require.Equal(t, int64(3), searchedOrganization.Organization.Contacts.TotalElements)

	contacts := searchedOrganization.Organization.Contacts.Content
	require.Equal(t, 1, len(contacts))
	require.Equal(t, contactId1, contacts[0].ID)
}

func TestQueryResolver_Organization_WithTimelineEvents_DirectAndFromMultipleContacts(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	jobRoleId1 := neo4jt.LinkContactWithOrganization(ctx, driver, contactId1, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId2, organizationId)

	now := time.Now().UTC()
	secInFuture10 := now.Add(time.Duration(10) * time.Second)
	secAgo5 := now.Add(time.Duration(-5) * time.Second)
	secAgo50 := now.Add(time.Duration(-50) * time.Second)
	secAgo60 := now.Add(time.Duration(-60) * time.Second)
	secAgo70 := now.Add(time.Duration(-70) * time.Second)
	secAgo80 := now.Add(time.Duration(-80) * time.Second)
	secAgo100 := now.Add(time.Duration(-100) * time.Second)
	secAgo110 := now.Add(time.Duration(-110) * time.Second)
	secAgo120 := now.Add(time.Duration(-120) * time.Second)
	secAgo1000 := now.Add(time.Duration(-1000) * time.Second)

	actionId1 := neo4jt.CreateActionForOrganization(ctx, driver, tenantName, organizationId, neo4jenum.ActionCreated, secAgo5)

	// prepare contact and org interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 1", "application/json", channel, secAgo50)
	interactionEventId2 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 2", "application/json", channel, secAgo60)
	interactionEventId3 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 3", "application/json", channel, secAgo70)
	emailIdContact := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId1, "email1", false, "WORK")
	emailIdOrg := neo4jt.AddEmailTo(ctx, driver, commonModel.ORGANIZATION, tenantName, organizationId, "email2", false, "WORK")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId2, "+1234", false, "WORK")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailIdContact, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId2, phoneNumberId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, emailIdOrg, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, phoneNumberId, "")

	// prepare direct interaction events
	interactionEventId4 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 4", "application/json", "EMAIL", secAgo100)
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4, organizationId, "TO")

	// prepare direct interaction events linked to job role
	interactionEventId5 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 5", "application/json", "EMAIL", secAgo110)
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId5, jobRoleId1, "")

	// prepare log entry for organization
	logEntryId := neo4jtest.CreateLogEntryForOrganization(ctx, driver, tenantName, organizationId, neo4jentity.LogEntryEntity{
		StartedAt:   secAgo120,
		Content:     "log entry content",
		ContentType: "text/plain",
	})
	userId := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)
	neo4jt.LogEntryCreatedByUser(ctx, driver, logEntryId, userId)

	// prepare issue with tags
	issueId1 := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{
		Subject:     "subject 1",
		CreatedAt:   secAgo80,
		Priority:    "P1",
		Status:      "OPEN",
		Description: "description 1",
	})
	tagId1 := neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{
		Name: "tag1",
	})
	tagId2 := neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{
		Name: "tag2",
	})
	neo4jt.TagIssue(ctx, driver, issueId1, tagId1)
	neo4jt.TagIssue(ctx, driver, issueId1, tagId2)
	neo4jt.IssueReportedBy(ctx, driver, issueId1, organizationId)

	issueId2 := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{CreatedAt: secAgo1000})
	neo4jt.IssueReportedBy(ctx, driver, issueId2, organizationId)
	issueId3 := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{CreatedAt: secInFuture10})
	neo4jt.IssueReportedBy(ctx, driver, issueId3, organizationId)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Issue"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Action"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "LogEntry"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 10, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_timeline_events_direct_and_via_contacts"),
		client.Var("organizationId", organizationId),
		client.Var("from", now),
		client.Var("size", 12))
	assertRawResponseSuccess(t, rawResponse, err)

	organization := rawResponse.Data.(map[string]interface{})["organization"]
	require.Equal(t, organizationId, organization.(map[string]interface{})["id"])

	timelineEvents := organization.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 9, len(timelineEvents))

	timelineEvent1 := timelineEvents[0].(map[string]interface{})
	require.Equal(t, "Action", timelineEvent1["__typename"].(string))
	require.Equal(t, actionId1, timelineEvent1["id"].(string))
	require.NotNil(t, timelineEvent1["createdAt"].(string))
	require.Equal(t, "CREATED", timelineEvent1["actionType"].(string))

	timelineEvent2 := timelineEvents[1].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent2["__typename"].(string))
	require.Equal(t, interactionEventId1, timelineEvent2["id"].(string))
	require.NotNil(t, timelineEvent2["createdAt"].(string))
	require.Equal(t, "IE text 1", timelineEvent2["content"].(string))

	timelineEvent3 := timelineEvents[2].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent3["__typename"].(string))
	require.Equal(t, interactionEventId2, timelineEvent3["id"].(string))
	require.NotNil(t, timelineEvent3["createdAt"].(string))
	require.Equal(t, "IE text 2", timelineEvent3["content"].(string))

	timelineEvent4 := timelineEvents[3].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent4["__typename"].(string))
	require.Equal(t, interactionEventId3, timelineEvent4["id"].(string))
	require.NotNil(t, timelineEvent4["createdAt"].(string))
	require.Equal(t, "IE text 3", timelineEvent4["content"].(string))

	timelineEvent5 := timelineEvents[4].(map[string]interface{})
	require.Equal(t, "Issue", timelineEvent5["__typename"].(string))
	require.Equal(t, issueId1, timelineEvent5["id"].(string))
	require.NotNil(t, timelineEvent5["createdAt"].(string))
	require.Equal(t, "subject 1", timelineEvent5["subject"].(string))
	require.Equal(t, "P1", timelineEvent5["priority"].(string))
	require.Equal(t, "OPEN", timelineEvent5["status"].(string))
	require.Equal(t, "description 1", timelineEvent5["description"].(string))
	require.Equal(t, "test", timelineEvent5["appSource"].(string))
	require.Equal(t, "OPENLINE", timelineEvent5["source"].(string))
	require.Equal(t, "OPENLINE", timelineEvent5["sourceOfTruth"].(string))
	require.ElementsMatch(t, []string{tagId1, tagId2},
		[]string{
			timelineEvent5["tags"].([]interface{})[0].(map[string]interface{})["id"].(string),
			timelineEvent5["tags"].([]interface{})[1].(map[string]interface{})["id"].(string)})
	require.ElementsMatch(t, []string{"tag1", "tag2"},
		[]string{
			timelineEvent5["tags"].([]interface{})[0].(map[string]interface{})["name"].(string),
			timelineEvent5["tags"].([]interface{})[1].(map[string]interface{})["name"].(string)})

	timelineEvent6 := timelineEvents[5].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent6["__typename"].(string))
	require.Equal(t, interactionEventId4, timelineEvent6["id"].(string))
	require.NotNil(t, timelineEvent6["createdAt"].(string))
	require.Equal(t, "IE text 4", timelineEvent6["content"].(string))

	timelineEvent7 := timelineEvents[6].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent7["__typename"].(string))
	require.Equal(t, interactionEventId5, timelineEvent7["id"].(string))
	require.NotNil(t, timelineEvent7["createdAt"].(string))
	require.Equal(t, "IE text 5", timelineEvent7["content"].(string))

	timelineEvent8 := timelineEvents[7].(map[string]interface{})
	require.Equal(t, "LogEntry", timelineEvent8["__typename"].(string))
	require.Equal(t, logEntryId, timelineEvent8["id"].(string))
	require.NotNil(t, timelineEvent8["startedAt"].(string))
	require.Equal(t, "log entry content", timelineEvent8["content"].(string))
	require.Equal(t, "text/plain", timelineEvent8["contentType"].(string))
	require.Equal(t, userId, timelineEvent8["createdBy"].(map[string]interface{})["id"].(string))
}

func TestQueryResolver_Organization_WithTimelineEventsTotalCount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId1, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId2, organizationId)

	now := time.Now().UTC()

	neo4jt.CreateActionForOrganization(ctx, driver, tenantName, organizationId, neo4jenum.ActionCreated, now)

	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId1, "contact note 1", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId2, "contact note 2", "text/plain", now)
	neo4jt.CreateNoteForOrganization(ctx, driver, tenantName, organizationId, "org note 1", now)

	// prepare contact and org interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 1", "application/json", channel, now)
	interactionEventId2 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 2", "application/json", channel, now)
	interactionEventId3 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 3", "application/json", channel, now)
	interactionEventId4Hidden := neo4jtest.CreateInteractionEventFromEntity(ctx, driver, tenantName, neo4jentity.InteractionEventEntity{
		Identifier:  "myExternalId",
		Content:     "IE text 4",
		ContentType: "application/json",
		Channel:     channel,
		CreatedAt:   now,
		Hide:        true,
	})
	emailIdContact := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId1, "email1", false, "WORK")
	emailIdOrg := neo4jt.AddEmailTo(ctx, driver, commonModel.ORGANIZATION, tenantName, organizationId, "email2", false, "WORK")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId2, "+1234", false, "WORK")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailIdContact, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId2, phoneNumberId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, emailIdOrg, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, phoneNumberId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId4Hidden, emailIdContact, "")

	issueId1 := neo4jt.CreateIssue(ctx, driver, tenantName, entity.IssueEntity{
		Subject:     "subject 1",
		CreatedAt:   now,
		Priority:    "P1",
		Status:      "OPEN",
		Description: "description 1",
	})
	neo4jt.IssueReportedBy(ctx, driver, issueId1, organizationId)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Issue"))
	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Action"))
	require.Equal(t, 9, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent"))
	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, model2.NodeLabelContract))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_timeline_events_total_count"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	organization := rawResponse.Data.(map[string]interface{})["organization"]
	require.Equal(t, organizationId, organization.(map[string]interface{})["id"])
	require.Equal(t, float64(5), organization.(map[string]interface{})["timelineEventsTotalCount"].(float64))
}

func TestQueryResolver_Organization_WithEmails(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org")
	emailId1 := neo4jt.AddEmailTo(ctx, driver, commonModel.ORGANIZATION, tenantName, organizationId, "email1", true, "MAIN")
	emailId2 := neo4jt.AddEmailTo(ctx, driver, commonModel.ORGANIZATION, tenantName, organizationId, "email2", false, "WORK")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_emails"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationStruct struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	organization := organizationStruct.Organization

	require.Equal(t, organizationId, organization.ID)
	emails := organization.Emails
	require.Equal(t, 2, len(emails))
	var emailA, emailB *model.Email
	if emailId1 == emails[0].ID {
		emailA = emails[0]
		emailB = emails[1]
	} else {
		emailA = emails[1]
		emailB = emails[0]
	}
	require.Equal(t, emailId1, emailA.ID)
	require.NotNil(t, emailA.CreatedAt)
	require.Equal(t, true, emailA.Primary)
	require.Equal(t, "email1", *emailA.RawEmail)
	require.Equal(t, "email1", *emailA.Email)

	require.Equal(t, emailId2, emailB.ID)
	require.NotNil(t, emailB.CreatedAt)
	require.Equal(t, false, emailB.Primary)
	require.Equal(t, "email2", *emailB.RawEmail)
	require.Equal(t, "email2", *emailB.Email)
}

func TestQueryResolver_Organization_WithPhoneNumbers(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org")
	phoneNumberId1 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId, "+1111", true, "MAIN")
	phoneNumberId2 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId, "+2222", false, "WORK")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_phone_numbers"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationStruct struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	organization := organizationStruct.Organization

	require.Equal(t, organizationId, organization.ID)
	phoneNumbers := organization.PhoneNumbers
	require.Equal(t, 2, len(phoneNumbers))
	var phoneNumber1, phoneNumber2 *model.PhoneNumber
	if phoneNumberId1 == phoneNumbers[0].ID {
		phoneNumber1 = phoneNumbers[0]
		phoneNumber2 = phoneNumbers[1]
	} else {
		phoneNumber1 = phoneNumbers[1]
		phoneNumber2 = phoneNumbers[0]
	}
	require.Equal(t, phoneNumberId1, phoneNumber1.ID)
	require.NotNil(t, phoneNumber1.CreatedAt)
	require.Equal(t, true, phoneNumber1.Primary)
	require.Equal(t, "+1111", *phoneNumber1.RawPhoneNumber)
	require.Equal(t, "+1111", *phoneNumber1.E164)
	require.Equal(t, model.PhoneNumberLabelMain, *phoneNumber1.Label)

	require.Equal(t, phoneNumberId2, phoneNumber2.ID)
	require.NotNil(t, phoneNumber2.CreatedAt)
	require.Equal(t, false, phoneNumber2.Primary)
	require.Equal(t, "+2222", *phoneNumber2.RawPhoneNumber)
	require.Equal(t, "+2222", *phoneNumber2.E164)
	require.Equal(t, model.PhoneNumberLabelWork, *phoneNumber2.Label)
}

func TestQueryResolver_Organization_WithSubsidiaries(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	parentOrganizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "parent org")
	subsidiaryOrganizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "sub org 1")
	subsidiaryOrganizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "sub org 2")
	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parentOrganizationId, subsidiaryOrganizationId1, "shop")
	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parentOrganizationId, subsidiaryOrganizationId2, "station")

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "SUBSIDIARY_OF"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_subsidiaries"),
		client.Var("organizationId", parentOrganizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationStruct struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	parentOrganization := organizationStruct.Organization

	require.Equal(t, parentOrganizationId, parentOrganization.ID)
	subsidiaries := parentOrganization.Subsidiaries
	require.Equal(t, 2, len(subsidiaries))
	require.Equal(t, subsidiaryOrganizationId1, subsidiaries[0].Organization.ID)
	require.Equal(t, "shop", *subsidiaries[0].Type)
	require.Equal(t, subsidiaryOrganizationId2, subsidiaries[1].Organization.ID)
	require.Equal(t, "station", *subsidiaries[1].Type)
}

func TestQueryResolver_Organization_WithParentForSubsidiary(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	parentOrganizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "parent org")
	subsidiaryOrganizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "sub org")
	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parentOrganizationId, subsidiaryOrganizationId1, "shop")

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "SUBSIDIARY_OF"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_parent_for_subsidiary"),
		client.Var("organizationId", subsidiaryOrganizationId1))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationStruct struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	organization := organizationStruct.Organization

	require.Equal(t, subsidiaryOrganizationId1, organization.ID)
	require.Equal(t, 1, len(organization.SubsidiaryOf))
	require.Equal(t, parentOrganizationId, organization.SubsidiaryOf[0].Organization.ID)
	require.Equal(t, "shop", *organization.SubsidiaryOf[0].Type)
}

func TestQueryResolver_Organization_WithSuggestedMerges(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	suggestedAt := utils.Now()
	primaryOrganizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "primary")
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org")
	neo4jt.LinkSuggestedMerge(ctx, driver, primaryOrganizationId, organizationId, "AI", suggestedAt, 0.55)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "SUGGESTED_MERGE"))

	rawResponse := callGraphQL(t, "organization/get_organization_with_suggested_merges", map[string]interface{}{"organizationId": organizationId})

	var organizationStruct struct {
		Organization model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	organization := organizationStruct.Organization

	require.Equal(t, organizationId, organization.ID)
	suggestedMerge := organization.SuggestedMergeTo
	require.Equal(t, 1, len(suggestedMerge))
	require.Equal(t, "AI", *suggestedMerge[0].SuggestedBy)
	require.Equal(t, 0.55, *suggestedMerge[0].Confidence)
	require.NotNil(t, *suggestedMerge[0].SuggestedAt)
	require.Equal(t, primaryOrganizationId, suggestedMerge[0].Organization.ID)
}

func TestQueryResolver_Organization_WithAccountDetails(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	nextRenewal := utils.Now().Add(time.Duration(10) * time.Hour)

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Name: "org",
		RenewalSummary: neo4jentity.RenewalSummary{
			ArrForecast:            utils.ToPtr[float64](500),
			MaxArrForecast:         utils.ToPtr[float64](1000),
			RenewalLikelihood:      "HIGH",
			RenewalLikelihoodOrder: utils.ToPtr(int64(1)),
			NextRenewalAt:          utils.TimePtr(nextRenewal),
		},
	})

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))

	rawResponse := callGraphQL(t, "organization/get_organization_with_account_details", map[string]interface{}{"organizationId": organizationId})

	var organizationStruct struct {
		Organization model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	organization := organizationStruct.Organization

	require.Equal(t, organizationId, organization.ID)
	require.Equal(t, "org", organization.Name)
	require.Equal(t, 500.0, *organization.AccountDetails.RenewalSummary.ArrForecast)
	require.Equal(t, 1000.0, *organization.AccountDetails.RenewalSummary.MaxArrForecast)
	require.Equal(t, model.OpportunityRenewalLikelihoodHighRenewal, *organization.AccountDetails.RenewalSummary.RenewalLikelihood)
	require.Equal(t, nextRenewal, *organization.AccountDetails.RenewalSummary.NextRenewalDate)
}

func TestQueryResolver_Organization_WithSocials(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name")

	socialId1 := neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{
		Url: "url1",
	})
	socialId2 := neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{
		Url: "url2",
	})
	neo4jt.LinkSocialWithEntity(ctx, driver, orgId, socialId1)
	neo4jt.LinkSocialWithEntity(ctx, driver, orgId, socialId2)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS"))

	rawResponse := callGraphQL(t, "organization/get_organization_with_socials",
		map[string]interface{}{"organizationId": orgId})

	var orgStruct struct {
		Organization model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &orgStruct)
	require.Nil(t, err)

	organization := orgStruct.Organization
	require.NotNil(t, organization)
	require.Equal(t, 2, len(organization.Socials))

	require.Equal(t, socialId1, organization.Socials[0].ID)
	require.Equal(t, "url1", organization.Socials[0].URL)
	require.NotNil(t, organization.Socials[0].CreatedAt)
	require.NotNil(t, organization.Socials[0].UpdatedAt)

	require.Equal(t, socialId2, organization.Socials[1].ID)
	require.Equal(t, "url2", organization.Socials[1].URL)
	require.NotNil(t, organization.Socials[1].CreatedAt)
	require.NotNil(t, organization.Socials[1].UpdatedAt)
}

func TestQueryResolver_Organization_WithOwner(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	userId := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name")
	neo4jtest.UserOwnsOrganization(ctx, driver, userId, organizationId)

	rawResponse := callGraphQL(t, "organization/get_organization_with_owner",
		map[string]interface{}{"organizationId": organizationId})

	var organizationStruct struct {
		Organization model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)

	organization := organizationStruct.Organization
	require.Equal(t, organizationId, organization.ID)
	require.Equal(t, userId, organization.Owner.ID)
	require.Equal(t, "first", organization.Owner.FirstName)
	require.Equal(t, "last", organization.Owner.LastName)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Organization", "Organization_" + tenantName})
}

func TestQueryResolver_Organization_WithExternalLinks(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name")

	neo4jt.CreateHubspotExternalSystem(ctx, driver, tenantName)
	syncDate1 := utils.Now()
	syncDate2 := syncDate1.Add(time.Hour * 1)
	neo4jt.LinkWithHubspotExternalSystem(ctx, driver, organizationId, "111", utils.StringPtr("www.external1.com"), nil, syncDate1)
	neo4jt.LinkWithHubspotExternalSystem(ctx, driver, organizationId, "222", utils.StringPtr("www.external2.com"), nil, syncDate2)

	rawResponse := callGraphQL(t, "organization/get_organization_with_external_links",
		map[string]interface{}{"organizationId": organizationId})

	var organizationStruct struct {
		Organization model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)

	organization := organizationStruct.Organization
	require.Equal(t, organizationId, organization.ID)
	require.Equal(t, 2, len(organization.ExternalLinks))
	require.Equal(t, "111", *organization.ExternalLinks[0].ExternalID)
	require.Equal(t, "222", *organization.ExternalLinks[1].ExternalID)
	require.Equal(t, "www.external1.com", *organization.ExternalLinks[0].ExternalURL)
	require.Equal(t, "www.external2.com", *organization.ExternalLinks[1].ExternalURL)
	require.Nil(t, organization.ExternalLinks[0].ExternalSource)
	require.Nil(t, organization.ExternalLinks[1].ExternalSource)
	require.Equal(t, syncDate1, *organization.ExternalLinks[0].SyncDate)
	require.Equal(t, syncDate2, *organization.ExternalLinks[1].SyncDate)
}

func TestQueryResolver_OrganizationDistinctOwners(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	userId1 := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)
	userId2 := neo4jtest.CreateUser(ctx, driver, tenantName, neo4jentity.UserEntity{
		FirstName: "first2",
		LastName:  "last2",
	})
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name 1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name 2")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name 3")
	neo4jtest.UserOwnsOrganization(ctx, driver, userId1, organizationId1)
	neo4jtest.UserOwnsOrganization(ctx, driver, userId2, organizationId2)
	neo4jtest.UserOwnsOrganization(ctx, driver, userId2, organizationId3)

	rawResponse := callGraphQL(t, "organization/get_organization_owners", map[string]interface{}{})

	var usersStruct struct {
		Organization_DistinctOwners []model.User
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &usersStruct)
	require.Nil(t, err)
	require.NotNil(t, usersStruct)

	users := usersStruct.Organization_DistinctOwners
	require.Equal(t, 2, len(users))
	require.Equal(t, userId1, users[0].ID)
	require.Equal(t, userId2, users[1].ID)

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))
}

func TestQueryResolver_Organization_WithContracts(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	now := utils.Now()
	yesterday := now.Add(time.Duration(-24) * time.Hour)
	daysAgo1 := now.Add(time.Duration(-24) * time.Hour)
	daysAgo2 := now.Add(time.Duration(-48) * time.Hour)
	daysAgo3 := now.Add(time.Duration(-72) * time.Hour)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "org name"})
	orgId2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "just another org"})
	contractId1 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		Name:                  "contract 1",
		CreatedAt:             now,
		UpdatedAt:             now,
		ServiceStartedAt:      &daysAgo3,
		SignedAt:              &daysAgo2,
		EndedAt:               &daysAgo1,
		LengthInMonths:        1,
		ContractStatus:        neo4jenum.ContractStatusDraft,
		ContractUrl:           "url1",
		Source:                neo4jentity.DataSourceOpenline,
		AppSource:             "test1",
		OrganizationLegalName: "legal name 1",
		Country:               "country 1",
		Locality:              "locality 1",
		Zip:                   "zip 1",
		InvoiceEmail:          "invoice email 1",
		InvoiceNote:           "invoice note 1",
	})
	contractId2 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		Name:                  "contract 2",
		CreatedAt:             yesterday,
		UpdatedAt:             yesterday,
		ServiceStartedAt:      &daysAgo1,
		SignedAt:              &daysAgo3,
		EndedAt:               &daysAgo2,
		LengthInMonths:        12,
		ContractStatus:        neo4jenum.ContractStatusLive,
		ContractUrl:           "url2",
		Source:                neo4jentity.DataSourceOpenline,
		AppSource:             "test2",
		OrganizationLegalName: "legal name 2",
		Country:               "country 2",
		Locality:              "locality 2",
		Zip:                   "zip 2",
		InvoiceEmail:          "invoice email 2",
		InvoiceNote:           "invoice note 2",
	})
	contractId3 := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId2, neo4jentity.ContractEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		"Organization":           2,
		"Contract":               3,
		"Contract_" + tenantName: 3,
	})
	neo4jtest.AssertRelationship(ctx, t, driver, orgId, "HAS_CONTRACT", contractId1)
	neo4jtest.AssertRelationship(ctx, t, driver, orgId, "HAS_CONTRACT", contractId2)
	neo4jtest.AssertRelationship(ctx, t, driver, orgId2, "HAS_CONTRACT", contractId3)

	rawResponse := callGraphQL(t, "organization/get_organization_with_contracts",
		map[string]interface{}{"organizationId": orgId})

	var orgStruct struct {
		Organization model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &orgStruct)
	require.Nil(t, err)

	organization := orgStruct.Organization
	require.NotNil(t, organization)
	require.Equal(t, 2, len(organization.Contracts))

	firstContract := organization.Contracts[0]
	require.Equal(t, contractId1, firstContract.ID)
	require.Equal(t, "contract 1", firstContract.Name)
	require.Equal(t, now, firstContract.CreatedAt)
	require.Equal(t, now, firstContract.UpdatedAt)
	require.Equal(t, utils.ToDate(daysAgo3), *firstContract.ServiceStartedAt)
	require.Equal(t, utils.ToDate(daysAgo2), *firstContract.SignedAt)
	require.Equal(t, utils.ToDate(daysAgo1), *firstContract.EndedAt)
	require.Equal(t, model.ContractRenewalCycleMonthlyRenewal, firstContract.RenewalCycle)
	require.Equal(t, model.ContractStatusDraft, firstContract.Status)
	require.Equal(t, "url1", *firstContract.ContractURL)
	require.Equal(t, model.DataSourceOpenline, firstContract.Source)
	require.Equal(t, "test1", firstContract.AppSource)
	require.Equal(t, "legal name 1", *firstContract.OrganizationLegalName)
	require.Equal(t, "country 1", *firstContract.Country)
	require.Equal(t, "locality 1", *firstContract.Locality)
	require.Equal(t, "zip 1", *firstContract.Zip)
	require.Equal(t, "invoice email 1", *firstContract.InvoiceEmail)
	require.Equal(t, "invoice note 1", *firstContract.InvoiceNote)

	secondContract := organization.Contracts[1]
	require.Equal(t, contractId2, secondContract.ID)
	require.Equal(t, "contract 2", secondContract.Name)
	require.Equal(t, yesterday, secondContract.CreatedAt)
	require.Equal(t, yesterday, secondContract.UpdatedAt)
	require.Equal(t, utils.ToDate(daysAgo1), *secondContract.ServiceStartedAt)
	require.Equal(t, utils.ToDate(daysAgo3), *secondContract.SignedAt)
	require.Equal(t, utils.ToDate(daysAgo2), *secondContract.EndedAt)
	require.Equal(t, model.ContractRenewalCycleAnnualRenewal, secondContract.RenewalCycle)
	require.Equal(t, model.ContractStatusLive, secondContract.Status)
	require.Equal(t, "url2", *secondContract.ContractURL)
	require.Equal(t, model.DataSourceOpenline, secondContract.Source)
	require.Equal(t, "test2", secondContract.AppSource)
	require.Equal(t, "legal name 2", *secondContract.OrganizationLegalName)
	require.Equal(t, "country 2", *secondContract.Country)
	require.Equal(t, "locality 2", *secondContract.Locality)
	require.Equal(t, "zip 2", *secondContract.Zip)
	require.Equal(t, "invoice email 2", *secondContract.InvoiceEmail)
	require.Equal(t, "invoice note 2", *secondContract.InvoiceNote)
}

//func TestMutationResolver_OrganizationMerge_Properties(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//
//	parentOrgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "main organization"})
//	mergedOrgId1 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "to merge 1"})
//	mergedOrgId2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "to merge 2"})
//
//	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
//
//	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
//		RefreshArr: func(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
//			return &organizationpb.OrganizationIdGrpcResponse{
//				Id: parentOrgId,
//			}, nil
//		},
//		RefreshRenewalSummary: func(ctx context.Context, proto *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
//			return &organizationpb.OrganizationIdGrpcResponse{
//				Id: parentOrgId,
//			}, nil
//		},
//	}
//	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)
//
//	rawResponse, err := c.RawPost(getQuery("organization/merge_organizations"),
//		client.Var("parentOrganizationId", parentOrgId),
//		client.Var("mergedOrganizationId1", mergedOrgId1),
//		client.Var("mergedOrganizationId2", mergedOrgId2),
//	)
//	assertRawResponseSuccess(t, rawResponse, err)
//
//	var organizationStruct struct {
//		Organization_Merge model.Organization
//	}
//	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
//	organization := organizationStruct.Organization_Merge
//	require.NotNil(t, organization)
//
//	require.Equal(t, parentOrgId, organization.ID)
//	require.Equal(t, "main organization", organization.Name)
//
//	// Check only 1 organization remains after merge
//	// other 2 converted into MergedOrganization
//	// Other notes not impacted
//	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
//	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "MergedOrganization"))
//}

//func TestMutationResolver_OrganizationMerge_CheckSubsidiariesMerge(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//
//	parentOrgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "main organization")
//	mergedOrgId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "to merge 1")
//	mergedOrgId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "to merge 2")
//
//	subsidiaryOrgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "")
//	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, mergedOrgId1, subsidiaryOrgId, "shop")
//
//	parentForSubsidiaryOrgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "")
//	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parentForSubsidiaryOrgId, mergedOrgId2, "factory")
//
//	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
//	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "SUBSIDIARY_OF"))
//
//	rawResponse, err := c.RawPost(getQuery("organization/merge_organizations"),
//		client.Var("parentOrganizationId", parentOrgId),
//		client.Var("mergedOrganizationId1", mergedOrgId1),
//		client.Var("mergedOrganizationId2", mergedOrgId2),
//	)
//	assertRawResponseSuccess(t, rawResponse, err)
//
//	var organizationStruct struct {
//		Organization_Merge model.Organization
//	}
//	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
//	organization := organizationStruct.Organization_Merge
//	require.NotNil(t, organization)
//
//	require.Equal(t, parentOrgId, organization.ID)
//	require.Equal(t, 1, len(organization.Subsidiaries))
//	require.Equal(t, subsidiaryOrgId, organization.Subsidiaries[0].Organization.ID)
//	require.Equal(t, "shop", *organization.Subsidiaries[0].Type)
//	require.Equal(t, 1, len(organization.SubsidiaryOf))
//	require.Equal(t, parentForSubsidiaryOrgId, organization.SubsidiaryOf[0].Organization.ID)
//
//	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
//	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "MergedOrganization"))
//
//	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "SUBSIDIARY_OF"))
//}

//func TestMutationResolver_OrganizationMerge_MergeBetweenParentAndSubsidiaryOrg(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jtest.CreateTenant(ctx, driver, tenantName)
//
//	parentOrgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "main"})
//	mergedOrgId1 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "to merge 1"})
//	mergedOrgId2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "to merge 2"})
//
//	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parentOrgId, mergedOrgId1, "A")
//	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, mergedOrgId2, parentOrgId, "B")
//
//	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
//	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "SUBSIDIARY_OF"))
//
//	rawResponse, err := c.RawPost(getQuery("organization/merge_organizations"),
//		client.Var("parentOrganizationId", parentOrgId),
//		client.Var("mergedOrganizationId1", mergedOrgId1),
//		client.Var("mergedOrganizationId2", mergedOrgId2),
//	)
//	assertRawResponseSuccess(t, rawResponse, err)
//
//	var organizationStruct struct {
//		Organization_Merge model.Organization
//	}
//	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
//	organization := organizationStruct.Organization_Merge
//	require.NotNil(t, organization)
//
//	require.Equal(t, parentOrgId, organization.ID)
//	require.Equal(t, 0, len(organization.Subsidiaries))
//	require.Equal(t, 0, len(organization.SubsidiaryOf))
//
//	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
//	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "MergedOrganization"))
//
//	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "SUBSIDIARY_OF"))
//}

func TestMutationResolver_OrganizationAddSubsidiary(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	parentOrgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "parent")
	subOrgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "sub")
	subsidiaryType := "shop"

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})

	calledAddParent := false

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		AddParent: func(context context.Context, org *organizationpb.AddParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, subsidiaryType, org.Type)
			require.Equal(t, subOrgId, org.OrganizationId)
			require.Equal(t, parentOrgId, org.ParentOrganizationId)
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, constants.AppSourceCustomerOsApi, org.AppSource)
			require.Equal(t, testUserId, org.LoggedInUserId)
			calledAddParent = true
			neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parentOrgId, subOrgId, subsidiaryType)
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: subOrgId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	rawResponse, err := c.RawPost(getQuery("organization/add_subsidiary"),
		client.Var("organizationId", parentOrgId),
		client.Var("subsidiaryId", subOrgId),
		client.Var("type", subsidiaryType),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationStruct struct {
		Organization_AddSubsidiary model.Organization
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	organization := organizationStruct.Organization_AddSubsidiary
	require.NotNil(t, organization)
	require.Equal(t, parentOrgId, organization.ID)
	require.True(t, calledAddParent)

	// below changes were done directly in DB, checking for returned data consistency
	require.Equal(t, 1, len(organization.Subsidiaries))
	require.Equal(t, subOrgId, organization.Subsidiaries[0].Organization.ID)
	require.Equal(t, "shop", *organization.Subsidiaries[0].Type)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	neo4jtest.AssertRelationship(ctx, t, driver, subOrgId, "SUBSIDIARY_OF", parentOrgId)
}

func TestMutationResolver_OrganizationRemoveSubsidiary(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	parentOrgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "main")
	subOrgId := neo4jt.CreateOrganization(ctx, driver, tenantName, "sub")
	neo4jt.LinkOrganizationAsSubsidiary(ctx, driver, parentOrgId, subOrgId, "shop")

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})

	calledRemoveParent := false

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		RemoveParent: func(context context.Context, org *organizationpb.RemoveParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, subOrgId, org.OrganizationId)
			require.Equal(t, parentOrgId, org.ParentOrganizationId)
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, constants.AppSourceCustomerOsApi, org.AppSource)
			require.Equal(t, testUserId, org.LoggedInUserId)
			calledRemoveParent = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: subOrgId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	rawResponse, err := c.RawPost(getQuery("organization/remove_subsidiary"),
		client.Var("organizationId", parentOrgId),
		client.Var("subsidiaryId", subOrgId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationStruct struct {
		Organization_RemoveSubsidiary model.Organization
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	organization := organizationStruct.Organization_RemoveSubsidiary
	require.NotNil(t, organization)
	require.Equal(t, parentOrgId, organization.ID)
	require.True(t, calledRemoveParent)
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
}

func TestMutationResolver_OrganizationSetOwner_NewOwner(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	userId := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name")

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		UpdateOrganizationOwner: func(context context.Context, org *organizationpb.UpdateOrganizationOwnerGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, organizationId, org.OrganizationId)
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, constants.AppSourceCustomerOsApi, org.AppSource)
			require.Equal(t, userId, org.OwnerUserId)
			neo4jtest.UserOwnsOrganization(ctx, driver, userId, organizationId)
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: organizationId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	rawResponse := callGraphQL(t, "organization/set_owner",
		map[string]interface{}{"organizationId": organizationId, "userId": userId})

	var organizationStruct struct {
		Organization_SetOwner model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)

	organization := organizationStruct.Organization_SetOwner
	require.Equal(t, organizationId, organization.ID)
	require.Equal(t, userId, organization.Owner.ID)
	test.AssertRecentTime(t, organization.UpdatedAt)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Organization", "Organization_" + tenantName})
}

func TestMutationResolver_OrganizationSetOwner_ReplaceOwner(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	previousOwnerId := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)
	newOwnerId := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name")
	neo4jtest.UserOwnsOrganization(ctx, driver, previousOwnerId, organizationId)

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		UpdateOrganizationOwner: func(context context.Context, org *organizationpb.UpdateOrganizationOwnerGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, organizationId, org.OrganizationId)
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, constants.AppSourceCustomerOsApi, org.AppSource)
			require.Equal(t, newOwnerId, org.OwnerUserId)
			neo4jtest.UserOwnsOrganization(ctx, driver, newOwnerId, organizationId)
			neo4jt.DeleteUserOwnsOrganization(ctx, driver, previousOwnerId, organizationId)
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: organizationId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	rawResponse := callGraphQL(t, "organization/set_owner",
		map[string]interface{}{"organizationId": organizationId, "userId": newOwnerId})

	var organizationStruct struct {
		Organization_SetOwner model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)

	organization := organizationStruct.Organization_SetOwner
	require.Equal(t, organizationId, organization.ID)
	require.Equal(t, newOwnerId, organization.Owner.ID)
	test.AssertRecentTime(t, organization.UpdatedAt)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Organization", "Organization_" + tenantName})
}

func TestMutationResolver_OrganizationUnsetOwner(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	ownerId := neo4jtest.CreateDefaultUser(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org name")
	neo4jtest.UserOwnsOrganization(ctx, driver, ownerId, organizationId)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))

	rawResponse := callGraphQL(t, "organization/unset_owner",
		map[string]interface{}{"organizationId": organizationId})

	var organizationStruct struct {
		Organization_UnsetOwner model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)

	organization := organizationStruct.Organization_UnsetOwner
	require.Equal(t, organizationId, organization.ID)
	require.Nil(t, organization.Owner)
	test.AssertRecentTime(t, organization.UpdatedAt)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 0, neo4jtest.GetCountOfRelationships(ctx, driver, "OWNS"))
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "User", "User_" + tenantName, "Organization", "Organization_" + tenantName})
}

func TestMutationResolver_OrganizationUpdateOnboardingStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})

	calledEventsPlatform := false

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		UpdateOnboardingStatus: func(context context.Context, org *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, tenantName, org.Tenant)
			require.Equal(t, organizationId, org.OrganizationId)
			require.Equal(t, testUserId, org.LoggedInUserId)
			require.Equal(t, constants.AppSourceCustomerOsApi, org.AppSource)
			require.Equal(t, organizationpb.OnboardingStatus_ONBOARDING_STATUS_DONE, org.OnboardingStatus)
			require.Equal(t, "Set to done", org.Comments)
			calledEventsPlatform = true
			return &organizationpb.OrganizationIdGrpcResponse{
				Id: organizationId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	rawResponse := callGraphQL(t, "organization/update_onboarding_status",
		map[string]interface{}{"organizationId": organizationId})

	var organizationStruct struct {
		Organization_UpdateOnboardingStatus model.Organization
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)

	organization := organizationStruct.Organization_UpdateOnboardingStatus
	require.Equal(t, organizationId, organization.ID)
	require.True(t, calledEventsPlatform)
}

func TestMutationResolver_OrganizationUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{})

	calledUpdateOrganization := false

	organizationServiceCallbacks := events_platform.MockOrganizationServiceCallbacks{
		UpsertOrganization: func(context context.Context, request *organizationpb.UpsertOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
			require.Equal(t, organizationId, request.Id)
			require.Equal(t, tenantName, request.Tenant)
			require.Equal(t, "slackChannelId", request.SlackChannelId)
			require.Equal(t, true, request.IsPublic)
			require.Equal(t, "PROSPECT", request.Relationship)
			require.Equal(t, "TARGET", request.Stage)
			require.ElementsMatch(t, []organizationpb.OrganizationMaskField{
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_SLACK_CHANNEL_ID,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_IS_PUBLIC,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_RELATIONSHIP,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_STAGE,
			}, request.FieldsMask)
			calledUpdateOrganization = true

			return &organizationpb.OrganizationIdGrpcResponse{
				Id: organizationId,
			}, nil
		},
	}
	events_platform.SetOrganizationCallbacks(&organizationServiceCallbacks)

	rawResponse := callGraphQL(t, "organization/update_organization",
		map[string]interface{}{"input": map[string]interface{}{
			"id":             organizationId,
			"slackChannelId": "slackChannelId",
			"isCustomer":     true,
			"public":         true,
			"relationship":   "PROSPECT",
			"stage":          "TARGET",
		},
		})

	var organizationStruct struct {
		Organization_Update model.Organization
	}

	require.Equal(t, true, calledUpdateOrganization)

	err := decode.Decode(rawResponse.Data.(map[string]any), &organizationStruct)
	require.Nil(t, err)
	require.NotNil(t, organizationStruct)

	organization := organizationStruct.Organization_Update
	require.Equal(t, organizationId, organization.Metadata.ID)
}
