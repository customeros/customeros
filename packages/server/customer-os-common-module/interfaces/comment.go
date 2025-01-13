package comment

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type CommentService interface {
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, commentFields data_fields.CommentFields) (string, error)
}
