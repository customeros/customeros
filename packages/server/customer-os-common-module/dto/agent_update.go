package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type UpdateAgent struct {
	data_fields.AgentFields
	Capabilities *postgresentity.CapabilitiesConfig `json:"capabilities,omitempty"`
}
