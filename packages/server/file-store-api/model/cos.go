package model

import (
	graph_model "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
)

type AttachmentCreateResponse struct {
	graph_model.Attachment `json:"attachment_Create"`
}

type AttachmentResponse struct {
	graph_model.Attachment `json:"attachment"`
}
