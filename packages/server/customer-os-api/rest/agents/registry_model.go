package agents

type RegisterMasterAgentRequest struct {
	Type                string            `json:"type,omitempty"`
	Name                string            `json:"name,omitempty"`
	Goal                string            `json:"goal,omitempty"`
	Icon                string            `json:"icon,omitempty"`
	DefaultCapabilities []AgentCapability `json:"defaultCapabilities,omitempty"`
}

type RegisterMasterAgentResponse struct {
	ID string `json:"id,omitempty"`
}

type GetAgentRegistryResponse struct {
	Agents []AgentRegistryRecord `json:"agents,omitempty"`
}

type AgentRegistryRecord struct {
	ID                  string            `json:"id,omitempty"`
	Type                string            `json:"type,omitempty"`
	Name                string            `json:"name,omitempty"`
	DefaultCapabilities []AgentCapability `json:"defaultCapabilities,omitempty"`
	Goal                string            `json:"goal,omitempty"`
}

type AgentCapability struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}
