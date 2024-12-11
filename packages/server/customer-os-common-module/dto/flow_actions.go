package dto

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
)

type CreateMarkdownEvent struct {
	FlowExecutionID string
	MarkdownEvent   data_fields.MarkdownEventFields
}
