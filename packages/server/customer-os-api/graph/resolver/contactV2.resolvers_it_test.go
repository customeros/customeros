package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_UIContactsSearch_FilterByName(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "", LastName: "", Name: ""})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "AAA", LastName: "BBB", Name: ""})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "ccc", LastName: "aaa", Name: ""})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "ccc", LastName: "ddd", Name: "--Aa--"})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "ccc", LastName: "ddd", Name: "eee"})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelContact))

	searchBy := postgresEntity.ColumnViewTypeContactsName
	searchTerm := "Aa"

	assertContactSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsEmpty, 5, 1)
	assertContactSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsNotEmpty, 5, 4)
	assertContactSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorContains, 5, 3)
	assertContactSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorNotContains, 5, 2)
}

func TestQueryResolver_UIContactsSearch_SortByName(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty", FirstName: "", LastName: "", Name: ""})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "B1-AB1", Name: "B1", FirstName: "AB1", LastName: ""})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "A1", Name: "", FirstName: "", LastName: "A1"})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "B1-AB2", Name: "B1", FirstName: "", LastName: "AB2"})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "A2", Name: "A2", FirstName: "", LastName: ""})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelContact))

	expectedAsc := []string{"A1", "A2", "B1-AB1", "B1-AB2", "empty"}
	expectedDesc := []string{"B1-AB2", "B1-AB1", "A2", "A1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsName, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsName, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByPrimaryEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contact1 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "A"})
	contact2 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "B"})
	contact3 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "C"})
	contact4 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "D"})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{FirstName: "E"})

	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact1, neo4jentity.EmailEntity{Email: "aaa111@gmail.com", Primary: true})
	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact1, neo4jentity.EmailEntity{Email: "bbb222@gmail.com", Primary: false})
	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact2, neo4jentity.EmailEntity{Email: "aaa333@gmail.com", Primary: false})
	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact2, neo4jentity.EmailEntity{Email: "ccc444@gmail.com", Primary: true})
	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact3, neo4jentity.EmailEntity{RawEmail: "aaa.aaa@gmail.com", Primary: true})
	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact4, neo4jentity.EmailEntity{RawEmail: "aaa111@gmail.com", Primary: false})

	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelContact))
	require.Equal(t, 6, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelEmail))

	searchBy := postgresEntity.ColumnViewTypeContactsPrimaryEmail
	searchTerm := "aa"

	assertContactSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsEmpty, 5, 2)
	assertContactSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorIsNotEmpty, 5, 3)
	assertContactSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorContains, 5, 2)
	assertContactSearch(t, searchBy, searchTerm, commonModel.ComparisonOperatorNotContains, 5, 1)
}

func TestQueryResolver_UIContactsSearch_SortByPrimaryEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contact1 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c1"})
	contact2 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c2"})
	contact3 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c3"})

	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact1, neo4jentity.EmailEntity{Email: "bbb@gmail.com", Primary: true})
	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact2, neo4jentity.EmailEntity{Email: "aaa@gmail.com", Primary: false})
	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact3, neo4jentity.EmailEntity{RawEmail: "aaa@gmail.com", Primary: true})

	expectedAsc := []string{contact3, contact1, contact2}
	expectedDesc := []string{contact1, contact3, contact2}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsPrimaryEmail, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsPrimaryEmail, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_SortByPrimaryEmail_WithNoPrimaryEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contact1 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{})
	contact2 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{})
	contact3 := neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{})

	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact1, neo4jentity.EmailEntity{Email: "aaa@gmail.com", Primary: true})
	neo4jtest.CreateEmailForEntity(ctx, driver, tenantName, contact2, neo4jentity.EmailEntity{RawEmail: "bbb@gmail.com", Primary: true})

	expectedAsc := []string{contact1, contact2, contact3}
	expectedDesc := []string{contact2, contact1, contact3}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsPrimaryEmail, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsPrimaryEmail, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByCountry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", CountryCodeA2: "US"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", CountryCodeA2: "CA"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l3"})
	neo4jtest.LinkNodes(ctx, driver, "3", "l3", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "4"})

	searchBy := postgresEntity.ColumnViewTypeContactsCountry

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 4, 2)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 4, 2)
	assertContactSearch(t, searchBy, []string{"X"}, commonModel.ComparisonOperatorIn, 4, 0)
	assertContactSearch(t, searchBy, []string{"US"}, commonModel.ComparisonOperatorIn, 4, 1)
	assertContactSearch(t, searchBy, []string{"CA"}, commonModel.ComparisonOperatorIn, 4, 1)
	assertContactSearch(t, searchBy, []string{"US", "CA"}, commonModel.ComparisonOperatorIn, 4, 2)

	assertContactSearch(t, searchBy, []string{"X"}, commonModel.ComparisonOperatorNotIn, 4, 4)
	assertContactSearch(t, searchBy, []string{"CA"}, commonModel.ComparisonOperatorNotIn, 4, 3)
}

func TestQueryResolver_UIContactsSearch_SortByCountry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", Country: "C1"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", Country: "C2"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsCountry, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsCountry, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByRegion(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", Region: "r1"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", Region: "r2"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})

	searchBy := postgresEntity.ColumnViewTypeContactsRegion

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorContains, 3, 2)
	assertContactSearch(t, searchBy, "r", commonModel.ComparisonOperatorContains, 3, 2)
	assertContactSearch(t, searchBy, "r1", commonModel.ComparisonOperatorContains, 3, 1)
	assertContactSearch(t, searchBy, "r2", commonModel.ComparisonOperatorContains, 3, 1)
	assertContactSearch(t, searchBy, "r3", commonModel.ComparisonOperatorContains, 3, 0)

	assertContactSearch(t, searchBy, "r1", commonModel.ComparisonOperatorNotContains, 3, 1)
	assertContactSearch(t, searchBy, "r2", commonModel.ComparisonOperatorNotContains, 3, 1)
	assertContactSearch(t, searchBy, "r3", commonModel.ComparisonOperatorNotContains, 3, 2)
}

func TestQueryResolver_UIContactsSearch_SortByRegion(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", Region: "region-aaa"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", Region: "region-BBB"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsRegion, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsRegion, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByCity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", Locality: "c1"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", Locality: "c2"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})

	searchBy := postgresEntity.ColumnViewTypeContactsCity

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorContains, 3, 2)
	assertContactSearch(t, searchBy, "c", commonModel.ComparisonOperatorContains, 3, 2)
	assertContactSearch(t, searchBy, "c1", commonModel.ComparisonOperatorContains, 3, 1)
	assertContactSearch(t, searchBy, "c2", commonModel.ComparisonOperatorContains, 3, 1)
	assertContactSearch(t, searchBy, "c3", commonModel.ComparisonOperatorContains, 3, 0)

	assertContactSearch(t, searchBy, "c1", commonModel.ComparisonOperatorNotContains, 3, 1)
	assertContactSearch(t, searchBy, "c2", commonModel.ComparisonOperatorNotContains, 3, 1)
	assertContactSearch(t, searchBy, "c3", commonModel.ComparisonOperatorNotContains, 3, 2)
}

func TestQueryResolver_UIContactsSearch_SortByCity(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l1", Locality: "L1"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{Id: "l2", Locality: "L2"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "ASSOCIATED_WITH")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsCity, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsCity, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_SortByCreatedDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1", CreatedAt: firstOfDecember})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2", CreatedAt: firstOfJanuary})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelContact))

	expectedAsc := []string{"1", "2"}
	expectedDesc := []string{"2", "1"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsCreatedAt, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsCreatedAt, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_SortByUpdatedDate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	firstOfDecember := utils.FirstTimeOfMonth(2023, 12)
	firstOfJanuary := utils.FirstTimeOfMonth(2024, 1)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1", UpdatedAt: firstOfDecember})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2", UpdatedAt: firstOfJanuary})

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelContact))

	expectedAsc := []string{"1", "2"}
	expectedDesc := []string{"2", "1"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsUpdatedAt, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsUpdatedAt, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByTags(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "contact1"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag1", Name: "A"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag2", Name: "B"})
	neo4jtest.LinkNodes(ctx, driver, "contact1", "tag1", "TAGGED")
	neo4jtest.LinkNodes(ctx, driver, "contact1", "tag2", "TAGGED")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "contact2"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag3", Name: "A"})
	neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{Id: "tag4", Name: "c"})
	neo4jtest.LinkNodes(ctx, driver, "contact2", "tag3", "TAGGED")
	neo4jtest.LinkNodes(ctx, driver, "contact2", "tag4", "TAGGED")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "contact3"})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelContact))
	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, commonModel.NodeLabelTag))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "TAGGED"))

	searchBy := postgresEntity.ColumnViewTypeContactsTags

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertContactSearch(t, searchBy, []string{"A", "B", "C"}, commonModel.ComparisonOperatorIn, 3, 2)
	assertContactSearch(t, searchBy, []string{"c"}, commonModel.ComparisonOperatorIn, 3, 1)
	assertContactSearch(t, searchBy, []string{"X", "Y", "Z"}, commonModel.ComparisonOperatorIn, 3, 0)

	assertContactSearch(t, searchBy, []string{"X"}, commonModel.ComparisonOperatorNotIn, 3, 3)
	assertContactSearch(t, searchBy, []string{"B", "Y"}, commonModel.ComparisonOperatorNotIn, 3, 2)
}

func TestQueryResolver_UIContactsSearch_FilterByLinkedIn(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "l1", Url: "linkedin.com/in/aaa", Alias: "ab"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "l2", Url: "linkedin.com/in/ababab", Alias: "bc"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "l3", Url: "ab", Alias: "aaa"})
	neo4jtest.LinkNodes(ctx, driver, "3", "l3", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "4"})

	searchBy := postgresEntity.ColumnViewTypeContactsLinkedin

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 4, 2)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 4, 2)
	assertContactSearch(t, searchBy, "aaa", commonModel.ComparisonOperatorContains, 4, 1)
	assertContactSearch(t, searchBy, "ab", commonModel.ComparisonOperatorContains, 4, 2)
	assertContactSearch(t, searchBy, "bc", commonModel.ComparisonOperatorContains, 4, 1)

	assertContactSearch(t, searchBy, "xx", commonModel.ComparisonOperatorNotContains, 4, 4)
	assertContactSearch(t, searchBy, "ab", commonModel.ComparisonOperatorNotContains, 4, 2)
	assertContactSearch(t, searchBy, "aaa", commonModel.ComparisonOperatorNotContains, 4, 3)
}

func TestQueryResolver_UIContactsSearch_SortByLinkedIn(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "l1", Url: "linkedin.com/in/aaa", Alias: "ab"})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "l2", Url: "linkedin.com/in/aaa", Alias: "bc"})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})

	expectedAsc := []string{"1", "2", "3"}
	expectedDesc := []string{"2", "1", "3"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsLinkedin, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsLinkedin, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	org1 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "Aaa", Hide: false})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "1", org1, neo4jentity.JobRoleEntity{Primary: true})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	org2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "Aaa", Hide: false})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "2", org2, neo4jentity.JobRoleEntity{Primary: false})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})
	org3 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "Aaa", Hide: true})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "3", org3, neo4jentity.JobRoleEntity{Primary: true})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "4"})
	org4 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "BB-AA", Hide: false})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "4", org4, neo4jentity.JobRoleEntity{Primary: true})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "5"})

	searchBy := postgresEntity.ColumnViewTypeContactsOrganization

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 3)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 2)
	assertContactSearch(t, searchBy, "aaa", commonModel.ComparisonOperatorContains, 5, 1)
	assertContactSearch(t, searchBy, "aa", commonModel.ComparisonOperatorContains, 5, 2)

	assertContactSearch(t, searchBy, "aaa", commonModel.ComparisonOperatorNotContains, 5, 4)
	assertContactSearch(t, searchBy, "aa", commonModel.ComparisonOperatorNotContains, 5, 3)
}

func TestQueryResolver_UIContactsSearch_SortByOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	org1 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "aaa", Hide: false})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "1", org1, neo4jentity.JobRoleEntity{Primary: true})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	org2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Name: "BB-AA", Hide: false})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "2", org2, neo4jentity.JobRoleEntity{Primary: true})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsOrganization, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsOrganization, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByJobTitle(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1"})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "1", "1", neo4jentity.JobRoleEntity{Primary: true, JobTitle: "Aaa"})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2"})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "2", "2", neo4jentity.JobRoleEntity{Primary: false, JobTitle: "Aaa"})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "3"})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "3", "3", neo4jentity.JobRoleEntity{Primary: true, JobTitle: "BB"})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "4"})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "5"})

	searchBy := postgresEntity.ColumnViewTypeContactsJobTitle

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 3)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 2)
	assertContactSearch(t, searchBy, "aaa", commonModel.ComparisonOperatorContains, 5, 1)
	assertContactSearch(t, searchBy, "bb", commonModel.ComparisonOperatorContains, 5, 1)

	assertContactSearch(t, searchBy, "aa", commonModel.ComparisonOperatorNotContains, 5, 4)
}

func TestQueryResolver_UIContactsSearch_SortByJobTitle(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "1"})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "1", "1", neo4jentity.JobRoleEntity{Primary: true, JobTitle: "Aaa"})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{ID: "2"})
	neo4jtest.LinkContactWithOrganization(ctx, driver, tenantName, "2", "2", neo4jentity.JobRoleEntity{Primary: true, JobTitle: "BBB"})

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsJobTitle, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsJobTitle, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByPhoneNumber(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{Id: "1", RawPhoneNumber: "1"})
	neo4jtest.LinkNodes(ctx, driver, "1", "1", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{Id: "2", RawPhoneNumber: "2"})
	neo4jtest.LinkNodes(ctx, driver, "2", "2", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})
	neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{Id: "3", RawPhoneNumber: "11"})
	neo4jtest.LinkNodes(ctx, driver, "3", "3", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "4"})
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "5"})

	searchBy := postgresEntity.ColumnViewTypeContactsPhoneNumbers

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 5, 2)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 5, 3)
	assertContactSearch(t, searchBy, "1", commonModel.ComparisonOperatorContains, 5, 2)
	assertContactSearch(t, searchBy, "11", commonModel.ComparisonOperatorContains, 5, 1)
	assertContactSearch(t, searchBy, "2", commonModel.ComparisonOperatorContains, 5, 1)

	assertContactSearch(t, searchBy, "1", commonModel.ComparisonOperatorNotContains, 5, 3)
	assertContactSearch(t, searchBy, "11", commonModel.ComparisonOperatorNotContains, 5, 4)
}

func TestQueryResolver_UIContactsSearch_SortByPhoneNumber(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{Id: "1", RawPhoneNumber: "1"})
	neo4jtest.LinkNodes(ctx, driver, "1", "1", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreatePhoneNumber(ctx, driver, tenantName, neo4jentity.PhoneNumberEntity{Id: "2", RawPhoneNumber: "2"})
	neo4jtest.LinkNodes(ctx, driver, "2", "2", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsPhoneNumbers, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsPhoneNumbers, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByFlow(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f1", Name: "A"})
	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f2", Name: "AA"})
	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f3", Name: "B"})

	//contact 1 in flow 1 and flow 2
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c1"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp11"})
	neo4jtest.LinkNodes(ctx, driver, "f1", "fp11", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp11", "c1", "HAS")

	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp12"})
	neo4jtest.LinkNodes(ctx, driver, "f2", "fp12", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp12", "c1", "HAS")

	//contact 2 in flow 1
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c2"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp2"})
	neo4jtest.LinkNodes(ctx, driver, "f1", "fp2", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp2", "c2", "HAS")

	//contact 3 in flow 3
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c3"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp3"})
	neo4jtest.LinkNodes(ctx, driver, "f3", "fp3", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp3", "c3", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "4"})

	searchBy := postgresEntity.ColumnViewTypeContactsFlows

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 4, 1)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 4, 3)
	assertContactSearch(t, searchBy, []string{"f1"}, commonModel.ComparisonOperatorIn, 4, 2)
	assertContactSearch(t, searchBy, []string{"f2"}, commonModel.ComparisonOperatorIn, 4, 1)
	assertContactSearch(t, searchBy, []string{"f3"}, commonModel.ComparisonOperatorIn, 4, 1)
	assertContactSearch(t, searchBy, []string{"f4"}, commonModel.ComparisonOperatorIn, 4, 0)
	assertContactSearch(t, searchBy, []string{"f1", "f2"}, commonModel.ComparisonOperatorIn, 4, 2)
	assertContactSearch(t, searchBy, []string{"f1", "f3"}, commonModel.ComparisonOperatorIn, 4, 3)
	assertContactSearch(t, searchBy, []string{"f1", "f4"}, commonModel.ComparisonOperatorIn, 4, 2)

	assertContactSearch(t, searchBy, []string{"f1"}, commonModel.ComparisonOperatorNotIn, 4, 2)
	assertContactSearch(t, searchBy, []string{"f2"}, commonModel.ComparisonOperatorNotIn, 4, 3)
	assertContactSearch(t, searchBy, []string{"f3"}, commonModel.ComparisonOperatorNotIn, 4, 3)
	assertContactSearch(t, searchBy, []string{"f4"}, commonModel.ComparisonOperatorNotIn, 4, 4)
}

func TestQueryResolver_UIContactsSearch_SortByFlow(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f1", Name: "A"})
	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f2", Name: "B"})

	//contact 1 in flow 1
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c1"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp1"})
	neo4jtest.LinkNodes(ctx, driver, "f1", "fp1", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp1", "c1", "HAS")

	//contact 2 in flow 1
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c2"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp2"})
	neo4jtest.LinkNodes(ctx, driver, "f2", "fp2", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp2", "c2", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"c1", "c2", "empty"}
	expectedDesc := []string{"c2", "c1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsFlows, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsFlows, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_FilterByFlowParticipantStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f1", Name: "A"})
	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f2", Name: "AA"})
	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f3", Name: "B"})

	//contact 1 in flow 1 and flow 2
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c1"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp11", Status: neo4jentity.FlowParticipantStatusCompleted})
	neo4jtest.LinkNodes(ctx, driver, "f1", "fp11", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp11", "c1", "HAS")

	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp12", Status: neo4jentity.FlowParticipantStatusCompleted})
	neo4jtest.LinkNodes(ctx, driver, "f2", "fp12", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp12", "c1", "HAS")

	//contact 2 in flow 1
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c2"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp2", Status: neo4jentity.FlowParticipantStatusReady})
	neo4jtest.LinkNodes(ctx, driver, "f1", "fp2", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp2", "c2", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "3"})

	searchBy := postgresEntity.ColumnViewTypeContactsFlowStatus

	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsEmpty, 3, 1)
	assertContactSearch(t, searchBy, "", commonModel.ComparisonOperatorIsNotEmpty, 3, 2)
	assertContactSearch(t, searchBy, []string{"COMPLETED"}, commonModel.ComparisonOperatorIn, 3, 1)
	assertContactSearch(t, searchBy, []string{"READY"}, commonModel.ComparisonOperatorIn, 3, 1)
	assertContactSearch(t, searchBy, []string{"ON_HOLD"}, commonModel.ComparisonOperatorIn, 3, 0)
	assertContactSearch(t, searchBy, []string{"COMPLETED", "READY"}, commonModel.ComparisonOperatorIn, 3, 2)
	assertContactSearch(t, searchBy, []string{"COMPLETED", "ON_HOLD"}, commonModel.ComparisonOperatorIn, 3, 1)

	assertContactSearch(t, searchBy, []string{"COMPLETED"}, commonModel.ComparisonOperatorNotIn, 3, 2)
	assertContactSearch(t, searchBy, []string{"READY"}, commonModel.ComparisonOperatorNotIn, 3, 2)
	assertContactSearch(t, searchBy, []string{"ON_HOLD"}, commonModel.ComparisonOperatorNotIn, 3, 3)
}

func TestQueryResolver_UIContactsSearch_SortByFlowParticipantStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f1", Name: "A"})
	neo4jtest.CreateFlow(ctx, driver, tenantName, neo4jentity.FlowEntity{Id: "f2", Name: "B"})

	//contact 1 in flow 1
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c1"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp1", Status: neo4jentity.FlowParticipantStatusCompleted})
	neo4jtest.LinkNodes(ctx, driver, "f1", "fp1", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp1", "c1", "HAS")

	//contact 2 in flow 1
	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "c2"})
	neo4jtest.CreateFlowParticipant(ctx, driver, tenantName, neo4jentity.FlowParticipantEntity{Id: "fp2", Status: neo4jentity.FlowParticipantStatusReady})
	neo4jtest.LinkNodes(ctx, driver, "f2", "fp2", "HAS")
	neo4jtest.LinkNodes(ctx, driver, "fp2", "c2", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"c1", "c2", "empty"}
	expectedDesc := []string{"c2", "c1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsFlowStatus, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsFlowStatus, commonModel.SortingDirectionDesc, expectedDesc)
}

func TestQueryResolver_UIContactsSearch_SortByLinkedInFollowerCount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "1"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "l1", Url: "linkedin.com/in/aaa", FollowersCount: 10})
	neo4jtest.LinkNodes(ctx, driver, "1", "l1", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "2"})
	neo4jtest.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{Id: "l2", Url: "linkedin.com/in/aaa", FollowersCount: 20})
	neo4jtest.LinkNodes(ctx, driver, "2", "l2", "HAS")

	neo4jtest.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{Id: "empty"})

	expectedAsc := []string{"1", "2", "empty"}
	expectedDesc := []string{"2", "1", "empty"}

	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsLinkedinFollowerCount, commonModel.SortingDirectionAsc, expectedAsc)
	verifyContactSortOrder(t, postgresEntity.ColumnViewTypeContactsLinkedinFollowerCount, commonModel.SortingDirectionDesc, expectedDesc)
}

func assertContactSearch(t *testing.T, filterName postgresEntity.ColumnViewType, searchValue any, operator commonModel.ComparisonOperator, totalAvailable int64, totalElements int64) {
	rawResponse, err := c.RawPost(getQuery("contact/ui_contacts_search"),
		client.Var("limit", 10),
		client.Var("filterName", filterName),
		client.Var("filterValue", searchValue),
		client.Var("filterOperation", operator),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var contacts struct {
		Ui_Contacts_Search model.ContactSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contacts)
	require.Nil(t, err)
	require.NotNil(t, contacts)

	searchResult := contacts.Ui_Contacts_Search
	require.Equal(t, totalAvailable, searchResult.TotalAvailable)
	require.Equal(t, totalElements, searchResult.TotalElements)
	require.Equal(t, int(totalElements), len(searchResult.Ids))
}

func verifyContactSortOrder(t *testing.T, sortBy postgresEntity.ColumnViewType, direction commonModel.SortingDirection, expectedOrder []string) {
	sortedResult := assertContactSort(t, sortBy, direction)
	assert.Equal(t, len(expectedOrder), len(sortedResult), "Mismatch in result length")
	for i, expected := range expectedOrder {
		assert.Equal(t, expected, sortedResult[i], "Mismatch at index %d", i)
	}
}

func assertContactSort(t *testing.T, sortBy postgresEntity.ColumnViewType, sortDirection commonModel.SortingDirection) []string {
	rawResponse, err := c.RawPost(getQuery("contact/ui_contacts_sort"),
		client.Var("limit", 10),
		client.Var("sortByField", string(sortBy)),
		client.Var("sortByDirection", sortDirection),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var contacts struct {
		Ui_Contacts_Search model.ContactSearchResult
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contacts)
	require.Nil(t, err)
	require.NotNil(t, contacts)

	return contacts.Ui_Contacts_Search.Ids
}
