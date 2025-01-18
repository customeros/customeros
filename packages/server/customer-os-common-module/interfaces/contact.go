package interfaces

import (
	"context"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type ContactService interface {
	SetEmailService(email EmailService)
	SetOrganizationService(org OrganizationService)
	SetJobRoleService(jobrole JobRoleService)
	SetSocialService(social SocialService)
	SetFlowService(flow FlowService)
	IsInitialized() bool

	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, contactFields data_fields.ContactFields, updateOnlyIfEmpty bool, options ...common_srv.ServiceOptions) (string, error)
	CreateContactByLinkedIn(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkedInUrl string, options ...common_srv.ServiceOptions) (string, error)
	CreateContactWithOrganizationByEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, email string) (string, error)
	CreateContactByEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, email string, options ...common_srv.ServiceOptions) (string, error)
	HideContact(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId string) error
	ShowContact(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId string) error
	LinkContactWithOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId, organizationId, jobTitle, description, source string, primary bool, startedAt, endedAt *time.Time) error
	CheckContactExistsWithLinkedIn(ctx context.Context, url, alias, externalId string) (bool, string, error)
	CheckContactExistsWithEmail(ctx context.Context, email string) (bool, string, error)
	GetContactById(ctx context.Context, contactId string) (*entity.ContactEntity, error)
	GetContactsByIds(ctx context.Context, contactIds []string) ([]*entity.ContactEntity, error)
	SetPrimaryJobRole(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactId string, primaryOrganizationId *string) error
	GetFirstContactByEmail(ctx context.Context, email string) (*entity.ContactEntity, error)
}
