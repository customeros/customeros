package flows

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
)

type CreateFlowNodeRequest struct {
	Type      string  `json:"type"`
	Event     *string `json:"event"`
	EventData *any    `json:"eventData"`
}

type FlowNodeRecord struct {
	ID        string     `json:"id"`
	FlowID    string     `json:"flowId"`
	Type      string     `json:"type"`
	Event     *string    `json:"event,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	EventData *any       `json:"eventData,omitempty"`
}

type FlowNodeResponse struct {
	enum.BaseResponse
	Node FlowNodeRecord `json:"node"`
}

type FlowNodesResponse struct {
	enum.BaseResponse
	Nodes []FlowNodeRecord `json:"nodes"`
}
