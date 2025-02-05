package postgres_entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type CapabilitiesConfig struct {
	Capabilities []Capability `json:"capabilities"`
}

func (c *CapabilitiesConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for CapabilitiesConfig, got %T", value)
	}

	return json.Unmarshal(bytes, c)
}

func (c CapabilitiesConfig) Value() (driver.Value, error) {
	if c.Capabilities == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

type Capability struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Type        enum.AgentCapability `json:"type"`
	Error       string               `json:"error"`
	Config      json.RawMessage      `json:"config"`
	Active      bool                 `json:"active"`
	Description string               `json:"description"`
}

func (c *Capability) SetConfig(config interface{}) error {
	if config == nil || config == "" {
		c.Config = nil
		return nil
	}

	switch v := config.(type) {
	case string:
		if v == "" {
			c.Config = nil
			return nil
		}
		// Parse the JSON string into a map to modify values
		var configMap map[string]interface{}
		if err := json.Unmarshal([]byte(v), &configMap); err != nil {
			return err
		}

		replaceNullWithEmptyString(configMap)

		// Marshal back to JSON
		data, err := json.Marshal(configMap)
		if err != nil {
			return err
		}
		c.Config = json.RawMessage(data)

	default:
		data, err := json.Marshal(config)
		if err != nil {
			return err
		}
		if string(data) == "{}" || string(data) == `""` {
			c.Config = nil
			return nil
		}
		c.Config = json.RawMessage(data)
	}
	return nil
}

func replaceNullWithEmptyString(m map[string]interface{}) {
	for k, v := range m {
		switch val := v.(type) {
		case nil:
			m[k] = ""
		case map[string]interface{}:
			replaceNullWithEmptyString(val)
		}
	}
}

func (c *Capability) GetConfig(configPtr interface{}) error {
	if c.Config == nil || string(c.Config) == "" {
		return nil
	}
	return json.Unmarshal(c.Config, configPtr)
}

func (c *Capability) GetConfigString() string {
	str := string(c.Config)
	if str == "null" {
		str = ""
	}
	return str
}

type Agent struct {
	ID                 string             `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Type               enum.AgentType     `gorm:"column:type;type:varchar(50);not null;" json:"type"`
	Tenant             string             `gorm:"column:tenant;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"tenant" binding:"required"`
	Name               string             `gorm:"column:name;type:varchar(255);not null" json:"name" binding:"required"`
	CapabilitiesConfig CapabilitiesConfig `gorm:"column:capabilities_config;type:jsonb" json:"capabilities"`
	Configured         bool               `gorm:"column:configured;type:boolean;default:false" json:"capabilitiesConfigured"`
	Goal               string             `gorm:"column:goal;type:text" json:"goal"`
	Status             string             `gorm:"column:status;type:varchar(32)" json:"status"`
	IsActive           bool               `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
	FlowID             string             `gorm:"column:flow_id;type:varchar(255)" json:"flowId"`
	VisibleInUI        bool               `gorm:"column:visible_in_ui;type:boolean;default:true" json:"visibleInUI"`
	CreatedAt          time.Time          `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt          *time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	ErrorMessage       *string            `gorm:"column:error_message;type:varchar(255)" json:"errorMessage"`
	Color              string             `gorm:"column:color;type:varchar(255)" json:"color"`
	Icon               string             `gorm:"column:icon;type:varchar(255)" json:"icon"`
	RegistryID         string             `gorm:"column:registry_id;type:varchar(32)" json:"registryId"`
}

func (Agent) TableName() string {
	return "agents"
}

func (r *Agent) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("agent", 16)
	return nil
}

func (r *Agent) GetCapabilitiesConfigAsString() string {
	bytes, err := json.Marshal(r.CapabilitiesConfig)
	if err != nil {
		return ""
	}
	return string(bytes)
}
