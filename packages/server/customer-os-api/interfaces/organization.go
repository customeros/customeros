package cosapi_interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type OrganizationService interface {
	CountOrganizations(ctx context.Context, tenant string) (int64, error)
	ExistsById(ctx context.Context, organizationId string) (bool, error)
	GetByCustomerOsId(ctx context.Context, customerOsId string) (*neo4jentity.OrganizationEntity, error)
	GetByReferenceId(ctx context.Context, referenceId string) (*neo4jentity.OrganizationEntity, error)
	GetSubsidiariesForOrganizations(ctx context.Context, parentOrganizationIds []string) (*neo4jentity.OrganizationEntities, error)
	AddSubsidiary(ctx context.Context, parentOrganizationId, subsidiaryOrganizationId, subsidiaryType string, removeExisting bool) error
	RemoveSubsidiary(ctx context.Context, parentOrganizationId, subsidiaryOrganizationId string) error
	GetSubsidiariesOfForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error)
	UpdateLastTouchpoint(ctx context.Context, organizationId string)
	GetOrganizationsForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForSlackChannels(ctx context.Context, slackChannelIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForJobRoles(ctx context.Context, jobRoleIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForOpportunities(ctx context.Context, opportunityIds []string) (*neo4jentity.OrganizationEntities, error)
	GetOrganizationsForEmails(ctx context.Context, emailIds []string) (*neo4jentity.OrganizationEntities, error)
	UpdateLastTouchpointByContactId(ctx context.Context, contactId string)
	GetMinMaxRenewalForecastArr(ctx context.Context) (float64, float64, error)
	GetOrganizationsForContact(ctx context.Context, contactId string, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error)
	RemoveOwner(ctx context.Context, organizationId string) (*neo4jentity.OrganizationEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error)
	Merge(ctx context.Context, primaryOrganizationId, mergedOrganizationId string) error
	GetOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error)
	GetSuggestedMergeToForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OrganizationEntities, error)
}
