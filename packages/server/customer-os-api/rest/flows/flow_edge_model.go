package flows

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
)

type CreateFlowEdgeRequest struct {
	FromNodeID string  `json:"fromNodeId"`
	ToNodeID   string  `json:"toNodeId"`
	Condition  *string `json:"condition"`
	Data       *any    `json:"data"`
}

type FlowEdgeRecord struct {
	ID         string    `json:"id"`
	FlowID     string    `json:"flowId"`
	FromNodeID string    `json:"fromNodeId"`
	ToNodeID   string    `json:"toNodeId"`
	Condition  *string   `json:"condition,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	Data       *any      `json:"data,omitempty"`
}

type FlowEdgeResponse struct {
	enum.BaseResponse
	Edge FlowEdgeRecord `json:"edge"`
}

type FlowEdgesResponse struct {
	enum.BaseResponse
	Edges []FlowEdgeRecord `json:"edges"`
}
