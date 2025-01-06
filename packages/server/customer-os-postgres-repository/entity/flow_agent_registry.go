package entity

type FlowAgentRegistry struct {
	ID           uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Agent        string       `gorm:"column:agent;type:varchar(255);not null;index" json:"agent" binding:"required"`
	FriendlyName string       `gorm:"column:friendly_name;type:varchar(255);not null;index" json:"FriendlyName" binding:"required"`
	Description  string       `gorm:"column:description;type:varchar(255);not null" json:"description" binding:"required"`
	Objective    string       `gorm:"column:objective;type:text" json:"objective"`
	Schema       *AgentSchema `gorm:"column:schema;type:jsonb" json:"schema"`
	Status       string       `gorm:"column:status;type:varchar(50)" json:"status"`
}

type AgentSchema struct {
	Input  map[string]AgentSchemaProperty `json:"input"`  // Input requirements
	Output map[string]AgentSchemaProperty `json:"output"` // Output format
}

// SchemaProperty defines a single property in the schema
type AgentSchemaProperty struct {
	Name        string                         `json:"name"`
	Type        string                         `json:"type"`                  // string, number, boolean, array, object
	Required    bool                           `json:"required"`              // Whether this field is required
	Description string                         `json:"description,omitempty"` // Description of what this property is for
	Format      string                         `json:"format,omitempty"`      // Any specific format (e.g., "email", "date-time")
	Enum        []interface{}                  `json:"enum,omitempty"`        // Possible values if this is an enum
	Properties  map[string]AgentSchemaProperty `json:"properties,omitempty"`  // Nested properties if type is object
	Items       *AgentSchemaProperty           `json:"items,omitempty"`       // Schema for array items if type is array
}

func (FlowAgentRegistry) TableName() string {
	return "flow_agent_registry"
}
