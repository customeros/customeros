package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type EmailService interface {
	SetContactService(contact ContactService)
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	Merge(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant string, emailFields EmailFields, linkWith *common_srv.LinkWith) (*string, error)
	ReplaceEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, previousEmail string, emailFields EmailFields, linkWith common_srv.LinkWith) (*string, error)
	UnlinkEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, email, appSource string, linkWith common_srv.LinkWith) error
	DeleteOrphanEmail(ctx context.Context, emailId string) error
	GetAllEmailsForEntityIds(ctx context.Context, tenant string, entityType model.EntityType, entityIds []string) (*entity.EmailEntities, error)
	SetPrimary(ctx context.Context, email string, forEntity common_srv.LinkWith) error
	GetPrimaryEmailForEntityId(ctx context.Context, entityType model.EntityType, entityId string) (*entity.EmailEntity, error)
	GetPrimaryEmailsForEntityIds(ctx context.Context, entityType model.EntityType, entityIds []string) (*entity.EmailEntities, error)
	UpdateEmailValidationDetails(ctx context.Context, emailId string, validationFields data_fields.EmailValidationFields) error
	RequestEmailValidation(ctx context.Context, emailId string) error
}

type EmailFields struct {
	Email     string            `json:"email"`
	Source    entity.DataSource `json:"source"`
	AppSource string            `json:"appSource"`
	Primary   bool              `json:"primary"`
}
