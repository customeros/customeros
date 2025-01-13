package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type JobRoleService interface {
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, jobRoleId, contactId, organizationId *string, dataFields data_fields.JobRoleFields) (string, error)
	GetAllForContact(ctx context.Context, contactId string) (*entity.JobRoleEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*entity.JobRoleEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*entity.JobRoleEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*entity.JobRoleEntities, error)
	DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error)
	GetAllForUsers(ctx context.Context, userIds []string) (*entity.JobRoleEntities, error)
	GetJobRolesByIds(ctx context.Context, ids []string) (*entity.JobRoleEntities, error)
	IdentifyJobRole(ctx context.Context, contactId, organizationId string) (*entity.JobRoleEntity, error)
	GetById(ctx context.Context, jobRoleId string) (*entity.JobRoleEntity, error)
}
