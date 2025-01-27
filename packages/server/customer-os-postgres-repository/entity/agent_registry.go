package postgres_entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type CapabilitiesConfig struct {
	Capabilities []Capability `json:"capabilities"`
}

type Capability struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Type        enum.AgentCapabilityType `json:"type"`
	Error       string                   `json:"error"`
	Values      string                   `json:"values"`
	Active      bool                     `json:"active"`
	Description string                   `json:"description"`
}

// Scan implements the sql.Scanner interface so GORM can read from the DB.
func (c *CapabilitiesConfig) Scan(value interface{}) error {
	if value == nil {
		// Handle NULL
		*c = CapabilitiesConfig{}
		return nil
	}

	// Postgres jsonb columns come in as []byte
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan type %T into CapabilitiesConfig", value)
	}

	return json.Unmarshal(bytes, &c)
}

// Value implements the driver.Valuer interface so GORM can write to the DB.
func (c CapabilitiesConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type AgentRegistry struct {
	ID                 string             `gorm:"primaryKey;type:varchar(20)" json:"id"`
	Type               enum.AgentType     `gorm:"column:type;type:varchar(255)" json:"type" binding:"required"`
	Name               string             `gorm:"column:name;type:varchar(255);not null;index" json:"name" binding:"required"`
	CapabilitiesConfig CapabilitiesConfig `gorm:"column:capabilities_config;type:jsonb" json:"capabilities"`
	Goal               string             `gorm:"column:goal;type:text" json:"goal"`
	IsActive           bool               `gorm:"column:is_active;type:boolean;default:true" json:"isActive"`
	Icon               string             `gorm:"column:icon;type:text" json:"icon"`
}

func (AgentRegistry) TableName() string {
	return "agent_registry"
}

func (r *AgentRegistry) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("ar", 16)
	return nil
}
