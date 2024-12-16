package flows

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
)

type FlowRecord struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Trigger       string    `json:"triggerOn,omitempty"`
	TriggerNodeID string    `json:"triggerNodeId,omitempty"`
	VisibleUI     *bool     `json:"visible"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"UpadatedAt,omitempty"`
}

type FlowResponse struct {
	enum.BaseResponse
	Flow FlowRecord `json:"flow"`
}

type FlowsResponse struct {
	enum.BaseResponse
	Flows []FlowRecord `json:"flows"`
}
