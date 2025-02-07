package postgres_entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
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
	Goal         enum.AgentGoal `gorm:"column:goal;type:text" json:"goal"`
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
	Listeners    []Listener     `gorm:"foreignKey:AgentID" json:"listeners"`
}

func (Agent) TableName() string {
	return "agents"
}

func (r *Agent) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("agent", 16)
	return nil
}

func (a *Agent) UpdateCapabilities(updatedCaps []Capability) {
	// Build a lookup map of existing capabilities by ID.
	existingMap := make(map[string]*Capability, len(a.Capabilities))
	for i := range a.Capabilities {
		existingMap[a.Capabilities[i].ID] = &a.Capabilities[i]
	}

	// Iterate over the provided updated capabilities.
	for _, upd := range updatedCaps {
		// Skip if the updated capability has no ID.
		if upd.ID == "" {
			continue
		}
		// Check if a capability with this ID exists in the agent.
		if existingCap, ok := existingMap[upd.ID]; ok {
			// Update allowed fields.
			existingCap.Name = upd.Name
			existingCap.Config = upd.Config
			existingCap.Active = upd.Active
		}
	}
}

func (a *Agent) UpdateListeners(updatedListeners []Listener) {
	existingMap := make(map[string]*Listener, len(a.Listeners))
	for i := range a.Listeners {
		existingMap[a.Listeners[i].ID] = &a.Listeners[i]
	}

	for _, upd := range updatedListeners {
		if upd.ID == "" {
			continue
		}
		if existingListener, ok := existingMap[upd.ID]; ok {
			existingListener.Name = upd.Name
			existingListener.Config = upd.Config
			existingListener.Active = upd.Active
		}
	}
}

// Capability as a separate entity
type Capability struct {
	ID        string               `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Tenant    string               `gorm:"column:tenant;type:varchar(255)" json:"tenant"`
	Position  int                  `gorm:"column:position;type:integer" json:"order"`
	AgentID   string               `gorm:"column:agent_id;type:varchar(32);not null" json:"agentId"`
	Name      string               `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Type      enum.AgentCapability `gorm:"column:type;type:varchar(50);not null" json:"type"`
	Error     string               `gorm:"column:error;type:varchar(255)" json:"error"`
	Config    JSONConfig           `gorm:"column:config;type:jsonb" json:"config"`
	Active    bool                 `gorm:"column:active;type:boolean;default:true" json:"active"`
	CreatedAt time.Time            `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time           `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	configHandlerImpl
}

func (Capability) TableName() string {
	return "agent_capabilities"
}

func (c *Capability) BeforeCreate(tx *gorm.DB) error {
	c.ID = utils.GenerateNanoIdWithPrefix("cap", 16)
	return nil
}

type Listener struct {
	ID        string                  `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Tenant    string                  `gorm:"column:tenant;type:varchar(255)" json:"tenant"`
	Position  int                     `gorm:"column:position;type:integer" json:"order"`
	AgentID   string                  `gorm:"column:agent_id;type:varchar(32);not null" json:"agentId"`
	Name      string                  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Type      enum.AgentListenerEvent `gorm:"column:type;type:varchar(50);not null" json:"type"`
	Error     string                  `gorm:"column:error;type:varchar(255)" json:"error"`
	Config    JSONConfig              `gorm:"column:config;type:jsonb" json:"config"`
	Active    bool                    `gorm:"column:active;type:boolean;default:true" json:"active"`
	CreatedAt time.Time               `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time              `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	configHandlerImpl
}

func (Listener) TableName() string {
	return "agent_listeners"
}

func (l *Listener) BeforeCreate(tx *gorm.DB) error {
	l.ID = utils.GenerateNanoIdWithPrefix("lst", 16)
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

// ConfigHandler interface for common config operations
type ConfigHandler interface {
	SetConfig(config any) error
	GetConfig(configPtr any) error
	GetConfigString() string
}

// configHandlerImpl implements common config handling
type configHandlerImpl struct {
	Config JSONConfig `gorm:"column:config;type:jsonb" json:"config"`
}

func (ch *configHandlerImpl) SetConfig(config any) error {
	if config == nil || config == "" {
		ch.Config = nil
		return nil
	}

	switch v := config.(type) {
	case string:
		if v == "" {
			ch.Config = nil
			return nil
		}
		var configMap map[string]any
		if err := json.Unmarshal([]byte(v), &configMap); err != nil {
			return err
		}

		replaceNullWithEmptyStringOrArray(configMap)

		data, err := json.Marshal(configMap)
		if err != nil {
			return err
		}
		ch.Config = data

	default:
		data, err := json.Marshal(config)
		if err != nil {
			return err
		}
		if string(data) == "{}" || string(data) == `""` {
			ch.Config = nil
			return nil
		}
		ch.Config = data
	}
	return nil
}

func (ch *configHandlerImpl) GetConfig(configPtr any) error {
	if ch.Config == nil {
		// If the config type expects an array, initialize it as empty
		if reflect.TypeOf(configPtr).Elem().Kind() == reflect.Slice {
			reflect.ValueOf(configPtr).Elem().Set(reflect.MakeSlice(reflect.TypeOf(configPtr).Elem(), 0, 0))
		}
		return nil
	}
	if _, ok := configPtr.(*NoConfig); ok {
		return nil
	}
	return json.Unmarshal(ch.Config, configPtr)
}

func (ch *configHandlerImpl) GetConfigString() string {
	if ch.Config == nil {
		return ""
	}
	str := string(ch.Config)
	if str == "null" {
		str = ""
	}
	return str
}

type NoConfig struct{}

func replaceNullWithEmptyStringOrArray(m map[string]any) {
	for k, v := range m {
		switch val := v.(type) {
		case nil:
			m[k] = ""
		case []interface{}:
			// Handle array values
			for i, item := range val {
				if item == nil {
					val[i] = ""
				}
			}
			m[k] = val
		case map[string]any:
			replaceNullWithEmptyStringOrArray(val)
		}
	}
}
