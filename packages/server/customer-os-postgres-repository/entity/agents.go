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

// Agent entity without the embedded capabilities
type Agent struct {
	ID           string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Type         enum.AgentType `gorm:"column:type;type:varchar(50);not null;" json:"type"`
	Tenant       string         `gorm:"column:tenant;type:varchar(255);not null" json:"tenant" binding:"required"`
	Name         string         `gorm:"column:name;type:varchar(255);not null" json:"name" binding:"required"`
	Configured   bool           `gorm:"column:configured;type:boolean;default:false" json:"capabilitiesConfigured"`
	Goal         string         `gorm:"column:goal;type:text" json:"goal"`
	Status       string         `gorm:"column:status;type:varchar(32)" json:"status"`
	IsActive     bool           `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
	FlowID       string         `gorm:"column:flow_id;type:varchar(255)" json:"flowId"`
	VisibleInUI  bool           `gorm:"column:visible_in_ui;type:boolean;default:true" json:"visibleInUI"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    *time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	ErrorMessage *string        `gorm:"column:error_message;type:varchar(255)" json:"errorMessage"`
	Color        string         `gorm:"column:color;type:varchar(255)" json:"color"`
	Icon         string         `gorm:"column:icon;type:varchar(255)" json:"icon"`
	RegistryID   string         `gorm:"column:registry_id;type:varchar(32)" json:"registryId"`
	Capabilities []Capability   `gorm:"foreignKey:AgentID" json:"capabilities"`
}

func (Agent) TableName() string {
	return "agents"
}

func (r *Agent) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("agent", 16)
	return nil
}

// Capability as a separate entity
type Capability struct {
	ID          string               `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Tenant      string               `gorm:"column:tenant;type:varchar(255)" json:"tenant"`
	Position    int                  `gorm:"column:position;type:integer" json:"order"`
	AgentID     string               `gorm:"column:agent_id;type:varchar(32);not null" json:"agentId"`
	Name        string               `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Type        enum.AgentCapability `gorm:"column:type;type:varchar(50);not null" json:"type"`
	Error       string               `gorm:"column:error;type:varchar(255)" json:"error"`
	Config      JSONConfig           `gorm:"column:config;type:jsonb" json:"config"`
	Active      bool                 `gorm:"column:active;type:boolean;default:true" json:"active"`
	Description string               `gorm:"column:description;type:text" json:"description"`
	CreatedAt   time.Time            `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   *time.Time           `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Capability) TableName() string {
	return "agent_capabilities"
}

func (c *Capability) BeforeCreate(tx *gorm.DB) error {
	c.ID = utils.GenerateNanoIdWithPrefix("cap", 16)
	return nil
}

// JSONConfig type for handling the config JSON field
type JSONConfig json.RawMessage

func (j *JSONConfig) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for JSONConfig, got %T", value)
	}

	*j = JSONConfig(bytes)
	return nil
}

func (j JSONConfig) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return []byte(j), nil
}

// Helper methods for Capability
func (c *Capability) SetConfig(config any) error {
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
		var configMap map[string]any
		if err := json.Unmarshal([]byte(v), &configMap); err != nil {
			return err
		}

		replaceNullWithEmptyString(configMap)

		data, err := json.Marshal(configMap)
		if err != nil {
			return err
		}
		c.Config = data

	default:
		data, err := json.Marshal(config)
		if err != nil {
			return err
		}
		if string(data) == "{}" || string(data) == `""` {
			c.Config = nil
			return nil
		}
		c.Config = JSONConfig(data)
	}
	return nil
}

func (c *Capability) GetConfig(configPtr any) error {
	if c.Config == nil {
		return nil
	}
	if _, ok := configPtr.(*NoConfig); ok {
		return nil
	}
	return json.Unmarshal([]byte(c.Config), configPtr)
}

func (c *Capability) GetConfigString() string {
	if c.Config == nil {
		return ""
	}
	str := string(c.Config)
	if str == "null" {
		str = ""
	}
	return str
}

type NoConfig struct{}

func replaceNullWithEmptyString(m map[string]any) {
	for k, v := range m {
		switch val := v.(type) {
		case nil:
			m[k] = ""
		case map[string]any:
			replaceNullWithEmptyString(val)
		}
	}
}
