package dto

import (
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
)

type UpdateAgent struct {
	data_fields.AgentFields
	Capabilities []postgresentity.Capability `json:"capabilities,omitempty"`
	Listeners    []postgresentity.Listener   `json:"listeners,omitempty"`
}
