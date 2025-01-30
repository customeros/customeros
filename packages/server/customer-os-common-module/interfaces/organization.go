package interfaces

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"time"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type OrganizationService interface {
	SetSocialService(social SocialService)
	IsInitialized() bool

	GetById(ctx context.Context, tenant, organizationId string) (*neo4j_entity.OrganizationEntity, error)
	GetPrimaryDomainByOrgID(ctx context.Context, organizationId string) (string, error)
	GetOrganizationByDomain(ctx context.Context, domain string, includePrimaryDomainCheck bool) (*neo4j_entity.OrganizationEntity, error)

	CreateFromGlobalOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, globalOrgId uint64, dataFields data_fields.OrganizationFields) (string, error)
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, dataFields data_fields.OrganizationFields) (string, error)
	LinkWithDomain(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId, domain string) (bool, error)
	UnlinkDomain(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId, domain string) error

	Hide(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string) error
	Show(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string) error

	AddParentOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, parentOrganizationId, subOrganizationId, relationType string) error
	RemoveParentOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, parentOrganizationId, subOrganizationId string) error

	UpdateOnboardingStatus(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string, dataFields data_fields.OrganizationOnboardingStatusFields) error
	UpdateDerivedData(ctx context.Context, organizationId string) error
	UpdateRenewalSummary(ctx context.Context, organizationId string) error

	GetHiddenOrganizationIds(ctx context.Context, hiddenAfter time.Time) ([]string, error)
	GetMergedOrganizationIds(ctx context.Context, mergedAfter time.Time) ([]string, error)
	GetOrganizationsByStage(ctx context.Context, stage enum.OrganizationStage) (*neo4j_entity.OrganizationEntities, error)
	RequestRefreshLastTouchpoint(ctx context.Context, organizationId string) error
	RefreshLastTouchpoint(ctx context.Context, organizationId string) error
	CheckOrganizationExistsWithEmail(ctx context.Context, email string) (bool, string, error)
	CheckOrganizationExistsWithLinkedIn(ctx context.Context, url, alias, externalId string) (bool, string, error)
	GetPrimaryOrganizationsWithJobRoleForContacts(ctx context.Context, contactIds []string) (*neo4j_entity.OrganizationWithJobRoleEntities, error)
	ValidateOrganizationExists(ctx context.Context, tx *neo4j.ManagedTransaction, organizationId string) error
}
