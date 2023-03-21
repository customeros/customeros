package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestQueryResolver_Organizations_FilterByNameLike(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateOrganization(ctx, driver, tenantName, "A closed organization")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "OPENLINE")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "the openline")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "some other open organization")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "OpEnLiNe")

	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))

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

func TestQueryResolver_Organization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationInput := entity.OrganizationEntity{
		Name:        "Organization name",
		Description: "Organization description",
		Website:     "Organization_website.com",
		Industry:    "tech",
		IsPublic:    true,
	}
	organizationId1 := neo4jt.CreateFullOrganization(ctx, driver, tenantName, organizationInput)
	neo4jt.AddDomainToOrg(ctx, driver, organizationId1, "domain1.com")
	neo4jt.AddDomainToOrg(ctx, driver, organizationId1, "domain2.com")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "otherOrganization")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_by_id"),
		client.Var("organizationId", organizationId1),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organization struct {
		Organization model.Organization
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &organization)
	require.Nil(t, err)
	require.NotNil(t, organization)
	require.Equal(t, organizationId1, organization.Organization.ID)
	require.Equal(t, organizationInput.Name, organization.Organization.Name)
	require.Equal(t, organizationInput.Description, *organization.Organization.Description)
	require.Equal(t, []string{"domain1.com", "domain2.com"}, organization.Organization.Domains)
	require.Equal(t, organizationInput.Website, *organization.Organization.Website)
	require.Equal(t, organizationInput.IsPublic, *organization.Organization.IsPublic)
	require.Equal(t, organizationInput.Industry, *organization.Organization.Industry)
	require.NotNil(t, organization.Organization.CreatedAt)
}

func TestQueryResolver_Organizations_WithLocations(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "OPENLINE")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "some other organization")
	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:      "WORK",
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
		Country:   "testCountry",
		Region:    "testRegion",
		Locality:  "testLocality",
		Address:   "testAddress",
		Address2:  "testAddress2",
		Zip:       "testZip",
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Name:      "UNKNOWN",
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
	})
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId1, locationId2)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

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
	require.Equal(t, "WORK", locationWithAddressDtls.Name)
	require.NotNil(t, locationWithAddressDtls.CreatedAt)
	require.NotNil(t, locationWithAddressDtls.UpdatedAt)
	require.Equal(t, "test", *locationWithAddressDtls.AppSource)
	require.Equal(t, model.DataSourceOpenline, *locationWithAddressDtls.Source)
	require.Equal(t, "testCountry", *locationWithAddressDtls.Country)
	require.Equal(t, "testLocality", *locationWithAddressDtls.Locality)
	require.Equal(t, "testRegion", *locationWithAddressDtls.Region)
	require.Equal(t, "testAddress", *locationWithAddressDtls.Address)
	require.Equal(t, "testAddress2", *locationWithAddressDtls.Address2)
	require.Equal(t, "testZip", *locationWithAddressDtls.Zip)

	require.Equal(t, locationId2, locationWithoutAddressDtls.ID)
	require.Equal(t, "UNKNOWN", locationWithoutAddressDtls.Name)
	require.NotNil(t, locationWithoutAddressDtls.CreatedAt)
	require.NotNil(t, locationWithoutAddressDtls.UpdatedAt)
	require.Equal(t, "test", *locationWithoutAddressDtls.AppSource)
	require.Equal(t, model.DataSourceOpenline, *locationWithoutAddressDtls.Source)
	require.Equal(t, "", *locationWithoutAddressDtls.Country)
	require.Equal(t, "", *locationWithoutAddressDtls.Region)
	require.Equal(t, "", *locationWithoutAddressDtls.Locality)
	require.Equal(t, "", *locationWithoutAddressDtls.Address)
	require.Equal(t, "", *locationWithoutAddressDtls.Address2)
	require.Equal(t, "", *locationWithoutAddressDtls.Zip)
}

func TestQueryResolver_Organizations_WithTags(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "Org 1 with 2 tags")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "Org 2 with 1 tag")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "Org 3 with 0 tags")
	tag1 := neo4jt.CreateTag(ctx, driver, tenantName, "tag1")
	tag2 := neo4jt.CreateTag(ctx, driver, tenantName, "tag2")

	neo4jt.TagOrganization(ctx, driver, organizationId1, tag1)
	neo4jt.TagOrganization(ctx, driver, organizationId1, tag2)
	neo4jt.TagOrganization(ctx, driver, organizationId2, tag1)

	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 3, neo4jt.GetCountOfRelationships(ctx, driver, "TAGGED"))

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

func TestQueryResolver_Organization_WithNotes_ById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org")
	userId := neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
	noteId1 := neo4jt.CreateNoteForOrganization(ctx, driver, tenantName, organizationId, "note1", utils.Now())
	noteId2 := neo4jt.CreateNoteForOrganization(ctx, driver, tenantName, organizationId, "note2", utils.Now())
	neo4jt.NoteCreatedByUser(ctx, driver, noteId1, userId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "NOTED"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "CREATED"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_notes_by_id"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedOrganization struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedOrganization)
	require.Nil(t, err)
	require.Equal(t, organizationId, searchedOrganization.Organization.ID)

	notes := searchedOrganization.Organization.Notes.Content
	require.Equal(t, 2, len(notes))
	var noteWithUser, noteWithoutUser *model.Note
	if noteId1 == notes[0].ID {
		noteWithUser = notes[0]
		noteWithoutUser = notes[1]
	} else {
		noteWithUser = notes[1]
		noteWithoutUser = notes[0]
	}
	require.Equal(t, noteId1, noteWithUser.ID)
	require.Equal(t, "note1", noteWithUser.HTML)
	require.NotNil(t, noteWithUser.CreatedAt)
	require.NotNil(t, noteWithUser.CreatedBy)
	require.Equal(t, userId, noteWithUser.CreatedBy.ID)
	require.Equal(t, "first", noteWithUser.CreatedBy.FirstName)
	require.Equal(t, "last", noteWithUser.CreatedBy.LastName)

	require.Equal(t, noteId2, noteWithoutUser.ID)
	require.Equal(t, "note2", noteWithoutUser.HTML)
	require.NotNil(t, noteWithoutUser.CreatedAt)
	require.Nil(t, noteWithoutUser.CreatedBy)
}

func TestMutationResolver_OrganizationCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationTypeId := neo4jt.CreateOrganizationType(ctx, driver, tenantName, "COMPANY")

	// Ensure that the tenant and organization type nodes were created in the database.
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationType"))
	require.Equal(t, 2, neo4jt.GetTotalCountOfNodes(ctx, driver))

	// Call the "create_organization" mutation.
	rawResponse, err := c.RawPost(getQuery("organization/create_organization"),
		client.Var("organizationTypeId", organizationTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the "organization" struct.
	var organization struct {
		Organization_Create model.Organization
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &organization)
	require.Nil(t, err)
	require.NotNil(t, organization)

	// Assign the organization to a shorter variable for easier reference.
	createdOrganization := organization.Organization_Create

	// Ensure that the organization was created correctly.
	require.NotNil(t, createdOrganization.ID)
	require.NotNil(t, createdOrganization.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdOrganization.CreatedAt)
	require.NotNil(t, createdOrganization.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdOrganization.UpdatedAt)
	require.Equal(t, "organization name", createdOrganization.Name)
	require.Equal(t, "organization description", *createdOrganization.Description)
	require.Equal(t, []string{"domain1", "domain2"}, createdOrganization.Domains)
	require.Equal(t, "organization website", *createdOrganization.Website)
	require.Equal(t, "organization industry", *createdOrganization.Industry)
	require.Equal(t, true, *createdOrganization.IsPublic)
	require.Equal(t, organizationTypeId, createdOrganization.OrganizationType.ID)
	require.Equal(t, "COMPANY", createdOrganization.OrganizationType.Name)
	require.Equal(t, model.DataSourceOpenline, createdOrganization.Source)
	require.Equal(t, model.DataSourceOpenline, createdOrganization.SourceOfTruth)
	require.Equal(t, "test", createdOrganization.AppSource)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Domain"))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(ctx, t, driver, []string{"Domain", "Tenant", "OrganizationType", "Organization", "Organization_" + tenantName})
}

func TestMutationResolver_OrganizationUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "some organization")
	organizationTypeIdOrig := neo4jt.CreateOrganizationType(ctx, driver, tenantName, "ORIG")
	organizationTypeIdUpdate := neo4jt.CreateOrganizationType(ctx, driver, tenantName, "UPDATED")
	neo4jt.SetOrganizationTypeForOrganization(ctx, driver, organizationTypeIdOrig, organizationTypeIdOrig)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationType"))

	rawResponse, err := c.RawPost(getQuery("organization/update_organization"),
		client.Var("organizationId", organizationId),
		client.Var("organizationTypeId", organizationTypeIdUpdate))
	assertRawResponseSuccess(t, rawResponse, err)

	var organization struct {
		Organization_Update model.Organization
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &organization)
	require.Nil(t, err)
	require.NotNil(t, organization)
	updatedOrganization := organization.Organization_Update
	require.Equal(t, organizationId, updatedOrganization.ID)
	require.NotNil(t, updatedOrganization.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), updatedOrganization.UpdatedAt)
	require.Equal(t, "updated name", updatedOrganization.Name)
	require.Equal(t, "updated description", *updatedOrganization.Description)
	require.Equal(t, []string{"updated domain"}, updatedOrganization.Domains)
	require.Equal(t, "updated website", *updatedOrganization.Website)
	require.Equal(t, "updated industry", *updatedOrganization.Industry)
	require.Equal(t, true, *updatedOrganization.IsPublic)
	require.Equal(t, organizationTypeIdUpdate, updatedOrganization.OrganizationType.ID)
	require.Equal(t, "UPDATED", updatedOrganization.OrganizationType.Name)
	require.Equal(t, model.DataSourceOpenline, updatedOrganization.SourceOfTruth)

	// Check still single organization node exists after update, no new node created
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
}

func TestMutationResolver_OrganizationDelete(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "LLC LLC")
	locationId := neo4jt.CreateLocation(ctx, driver, tenantName, entity.LocationEntity{
		Source: "manual",
	})
	neo4jt.OrganizationAssociatedWithLocation(ctx, driver, organizationId, locationId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	rawResponse, err := c.RawPost(getQuery("organization/delete_organization"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		Organization_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.Organization_Delete.Result)

	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant"})
}

func TestQueryResolver_Organization_WithRoles_ById(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "some organization")
	role1 := neo4jt.ContactWorksForOrganization(ctx, driver, contactId1, organizationId, "CTO", false)
	role2 := neo4jt.ContactWorksForOrganization(ctx, driver, contactId2, organizationId, "CEO", true)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "JobRole"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "WORKS_AS"))

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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization2")
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId3 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId4 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId1, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId2, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId3, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId4, organizationId2)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(ctx, driver, "CONTACT_OF"))

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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId1, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId2, organizationId)

	now := time.Now().UTC()
	secInFuture10 := now.Add(time.Duration(10) * time.Second)
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo20 := now.Add(time.Duration(-20) * time.Second)
	secAgo30 := now.Add(time.Duration(-30) * time.Second)
	secAgo50 := now.Add(time.Duration(-50) * time.Second)
	secAgo60 := now.Add(time.Duration(-60) * time.Second)
	secAgo70 := now.Add(time.Duration(-70) * time.Second)
	secAgo1000 := now.Add(time.Duration(-1000) * time.Second)

	// prepare contact and org notes
	contactNoteId1 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId1, "contact note 1", secAgo10)
	contactNoteId2 := neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId2, "contact note 2", secAgo20)
	orgNoteId3 := neo4jt.CreateNoteForOrganization(ctx, driver, tenantName, organizationId, "org note 1", secAgo30)
	neo4jt.CreateNoteForOrganization(ctx, driver, tenantName, organizationId, "org note 2", secAgo1000)
	neo4jt.CreateNoteForOrganization(ctx, driver, tenantName, organizationId, "org note 3", secInFuture10)

	// prepare contact and org interaction events
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 1", "application/json", "EMAIL", secAgo50)
	interactionEventId2 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 2", "application/json", "EMAIL", secAgo60)
	interactionEventId3 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 3", "application/json", "EMAIL", secAgo70)
	emailIdContact := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId1, "email1", false, "WORK")
	emailIdOrg := neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId, "email2", false, "WORK")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId2, "+1234", false, "WORK")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailIdContact, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId2, phoneNumberId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, emailIdOrg, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, phoneNumberId, "")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 8, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_timeline_events_direct_and_via_contacts"),
		client.Var("organizationId", organizationId),
		client.Var("from", now),
		client.Var("size", 6))
	assertRawResponseSuccess(t, rawResponse, err)

	organization := rawResponse.Data.(map[string]interface{})["organization"]
	require.Equal(t, organizationId, organization.(map[string]interface{})["id"])

	timelineEvents := organization.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 6, len(timelineEvents))

	timelineEvent1 := timelineEvents[0].(map[string]interface{})
	require.Equal(t, "Note", timelineEvent1["__typename"].(string))
	require.Equal(t, contactNoteId1, timelineEvent1["id"].(string))
	require.NotNil(t, timelineEvent1["createdAt"].(string))
	require.Equal(t, "contact note 1", timelineEvent1["html"].(string))

	timelineEvent2 := timelineEvents[1].(map[string]interface{})
	require.Equal(t, "Note", timelineEvent2["__typename"].(string))
	require.Equal(t, contactNoteId2, timelineEvent2["id"].(string))
	require.NotNil(t, timelineEvent2["createdAt"].(string))
	require.Equal(t, "contact note 2", timelineEvent2["html"].(string))

	timelineEvent3 := timelineEvents[2].(map[string]interface{})
	require.Equal(t, "Note", timelineEvent3["__typename"].(string))
	require.Equal(t, orgNoteId3, timelineEvent3["id"].(string))
	require.NotNil(t, timelineEvent3["createdAt"].(string))
	require.Equal(t, "org note 1", timelineEvent3["html"].(string))

	timelineEvent4 := timelineEvents[3].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent4["__typename"].(string))
	require.Equal(t, interactionEventId1, timelineEvent4["id"].(string))
	require.NotNil(t, timelineEvent4["createdAt"].(string))
	require.Equal(t, "IE text 1", timelineEvent4["content"].(string))

	timelineEvent5 := timelineEvents[4].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent5["__typename"].(string))
	require.Equal(t, interactionEventId2, timelineEvent5["id"].(string))
	require.NotNil(t, timelineEvent5["createdAt"].(string))
	require.Equal(t, "IE text 2", timelineEvent5["content"].(string))

	timelineEvent6 := timelineEvents[5].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent6["__typename"].(string))
	require.Equal(t, interactionEventId3, timelineEvent6["id"].(string))
	require.NotNil(t, timelineEvent6["createdAt"].(string))
	require.Equal(t, "IE text 3", timelineEvent6["content"].(string))
}

func TestQueryResolver_Organization_WithTimelineEventsTotalCount(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId1, organizationId)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId2, organizationId)

	now := time.Now().UTC()

	// prepare contact amd org notes
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId1, "contact note 1", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId2, "contact note 2", now)
	neo4jt.CreateNoteForOrganization(ctx, driver, tenantName, organizationId, "org note 1", now)

	// prepare contact and org interaction events
	interactionEventId1 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 1", "application/json", "EMAIL", now)
	interactionEventId2 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 2", "application/json", "EMAIL", now)
	interactionEventId3 := neo4jt.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 3", "application/json", "EMAIL", now)
	emailIdContact := neo4jt.AddEmailTo(ctx, driver, entity.CONTACT, tenantName, contactId1, "email1", false, "WORK")
	emailIdOrg := neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId, "email2", false, "WORK")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId2, "+1234", false, "WORK")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailIdContact, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId2, phoneNumberId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId3, emailIdOrg, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId3, phoneNumberId, "")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 6, neo4jt.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_timeline_events_total_count"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	organization := rawResponse.Data.(map[string]interface{})["organization"]
	require.Equal(t, organizationId, organization.(map[string]interface{})["id"])
	require.Equal(t, float64(6), organization.(map[string]interface{})["timelineEventsTotalCount"].(float64))
}

func TestQueryResolver_Organization_WithEmails(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org")
	emailId1 := neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId, "email1", true, "MAIN")
	emailId2 := neo4jt.AddEmailTo(ctx, driver, entity.ORGANIZATION, tenantName, organizationId, "email2", false, "WORK")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

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
	require.Equal(t, model.EmailLabelMain, *emailA.Label)

	require.Equal(t, emailId2, emailB.ID)
	require.NotNil(t, emailB.CreatedAt)
	require.Equal(t, false, emailB.Primary)
	require.Equal(t, "email2", *emailB.RawEmail)
	require.Equal(t, "email2", *emailB.Email)
	require.Equal(t, model.EmailLabelWork, *emailB.Label)
}

func TestQueryResolver_Organization_WithPhoneNumbers(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	organizationId := neo4jt.CreateOrganization(ctx, driver, tenantName, "test org")
	phoneNumberId1 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId, "+1111", true, "MAIN")
	phoneNumberId2 := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, organizationId, "+2222", false, "WORK")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(ctx, driver, "HAS"))

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
