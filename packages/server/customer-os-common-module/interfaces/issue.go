package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type IssueService interface {
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, issueFields data_fields.IssueFields) (string, error)
	AddUserAssignee(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error
	RemoveUserAssignee(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error
	AddUserFollower(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error
	RemoveUserFollower(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, issueId, userId string) error
}
