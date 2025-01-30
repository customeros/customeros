package data_fields

import postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

type AgentFields struct {
	Active             *bool                              `json:"active,omitempty"`
	VisibleInUI        *bool                              `json:"visibleInUi,omitempty"`
	Color              *string                            `json:"color,omitempty"`
	Icon               *string                            `json:"icon,omitempty"`
	Name               *string                            `json:"name,omitempty"`
	Goal               *string                            `json:"goal,omitempty"`
	FlowID             *string                            `json:"flowId,omitempty"`
	CapabilitiesConfig *postgresentity.CapabilitiesConfig `json:"capabilities,omitempty"`
}
