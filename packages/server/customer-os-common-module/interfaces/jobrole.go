package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type JobRoleService interface {
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, jobRoleId, contactId, organizationId *string, dataFields data_fields.JobRoleFields) (string, error)
	GetAllForContact(ctx context.Context, contactId string) (*neo4j_entity.JobRoleEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*neo4j_entity.JobRoleEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*neo4j_entity.JobRoleEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4j_entity.JobRoleEntities, error)
	DeleteJobRole(ctx context.Context, contactId, roleId string) (bool, error)
	GetAllForUsers(ctx context.Context, userIds []string) (*neo4j_entity.JobRoleEntities, error)
	GetJobRolesByIds(ctx context.Context, ids []string) (*neo4j_entity.JobRoleEntities, error)
	IdentifyJobRole(ctx context.Context, contactId, organizationId string) (*neo4j_entity.JobRoleEntity, error)
	GetById(ctx context.Context, jobRoleId string) (*neo4j_entity.JobRoleEntity, error)
}
