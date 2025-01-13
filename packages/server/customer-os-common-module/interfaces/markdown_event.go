package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type MarkdownEventService interface {
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, input data_fields.MarkdownEventFields) (string, error)
}
