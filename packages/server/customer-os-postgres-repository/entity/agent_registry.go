package postgres_entity

import (
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AgentRegistry struct {
	ID               string          `gorm:"primaryKey;type:varchar(21)" json:"id"`
	Type             enum.AgentType  `gorm:"column:type;type:varchar(255);not null;index" json:"type" binding:"required"`
	AgentName        string          `gorm:"column:agent_name;type:varchar(255)" json:"agentName"`
	Description      string          `gorm:"column:description;type:text" json:"description"`
	Scope            enum.AgentScope `gorm:"column:agent_scope;type:varchar(255)" json:"agentScope"`
	Filename         string          `gorm:"column:filename;type:varchar(255)" json:"filename"`
	CompletionEvents pq.StringArray  `gorm:"column:completion_events;type:varchar[]" json:"completionEvents"`
	ListenerEvents   pq.StringArray  `gorm:"column:listener_events;type:varchar[]" json:"listenerEvents"`
	Capabilities     pq.StringArray  `gorm:"column:capabilities;type:varchar[]" json:"capabilities"`
	Version          string          `gorm:"column:version;type:varchar(21)" json:"version"`
	Goal             string          `gorm:"column:goal;type:varchar(255)" json:"goal"`
	IsActive         bool            `gorm:"column:is_active;type:boolean;default:true" json:"isActive"`
	Icon             string          `gorm:"column:icon;type:text" json:"icon"`
}

func (AgentRegistry) TableName() string {
	return "agent_registry"
}

func (r *AgentRegistry) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("ar", 16)
	err := r.ValidateCapabilities()
	if err != nil {
		return err
	}
	err = r.ValidateEvents()
	if err != nil {
		return err
	}
	return nil
}

func (r *AgentRegistry) ValidateCapabilities() error {
	if len(r.Capabilities) == 0 {
		err := errors.New("No capabilities configured for agent")
		return err
	}

	for _, capabilityType := range r.Capabilities {
		_, err := enum.GetAgentCapability(capabilityType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AgentRegistry) ValidateEvents() error {
	if len(r.CompletionEvents) == 0 {
		return errors.New("No completion events configured for agent")
	}
	for _, event := range r.CompletionEvents {
		_, err := enum.GetAgentListener(event)
		if err != nil {
			return err
		}
	}

	if len(r.ListenerEvents) == 0 {
		return nil
	}
	for _, event := range r.ListenerEvents {
		_, err := enum.GetAgentListener(event)
		if err != nil {
			return err
		}
	}
	return nil
}

type AgentPlay struct {
	AgentRegistryID string         `gorm:"column:agent_registry_id;type:varchar(21);primaryKey" json:"agentRegistryId"`
	AgentType       enum.AgentType `gorm:"column:agent_type;type:varchar(255);not null;index" json:"agentType" binding:"required"`
	TriggerEvent    string         `gorm:"column:trigger_event;type:varchar(255);primaryKey" json:"triggerEvent"`
	Capabilities    pq.StringArray `gorm:"column:capabilities;type:varchar[]" json:"capabilities"`
}

func (AgentPlay) TableName() string {
	return "agent_plays"
}

func (r *AgentPlay) BeforeCreate(tx *gorm.DB) error {
	err := r.ValidatePlay()
	if err != nil {
		return err
	}
	return r.ValidateTrigger()
}

func (r *AgentPlay) ValidatePlay() error {
	if len(r.Capabilities) == 0 {
		return nil
	}

	for _, capabilityType := range r.Capabilities {
		_, err := enum.GetAgentCapability(capabilityType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AgentPlay) ValidateTrigger() error {
	if r.TriggerEvent == "" {
		return errors.New("play trigger event is empty")
	}

	_, err := enum.GetAgentListener(r.TriggerEvent)
	return err
}
