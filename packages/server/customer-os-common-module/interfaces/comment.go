package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type CommentService interface {
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, commentFields data_fields.CommentFields) (string, error)
}
