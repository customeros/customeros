package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type DomainService interface {
	GetPrimaryDomainForOrganizationWebsite(ctx context.Context, websiteUrl string) (string, string)
	IsKnownCompanyHostingUrl(ctx context.Context, website string) bool
	GetAllDomainsForOrganizations(ctx context.Context, organizationIds []string) (*entity.DomainEntities, error)
	UpdateDomainPrimaryDetails(ctx context.Context, domain string) error
	MergeDomain(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, domain string) error
	GetDomain(ctx context.Context, domain string) (*entity.DomainEntity, error)
	IsAcceptedDomainForOrganization(ctx context.Context, domain string) bool
	CheckDomainWithMailsherpa(ctx context.Context, domain string) (bool, bool, string)
}
