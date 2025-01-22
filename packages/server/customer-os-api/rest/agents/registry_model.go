package agents

type RegisterMasterAgentRequest struct {
	Type                string            `json:"type"`
	Name                string            `json:"name"`
	Goal                string            `json:"goal"`
	Icon                string            `json:"icon"`
	DefaultCapabilities []AgentCapability `json:"defaultCapabilities"`
}

type RegisterMasterAgentResponse struct {
	ID string `json:"id"`
}

type GetAgentRegistryResponse struct {
	Agents []AgentRegistryRecord `json:"agents"`
}

type AgentRegistryRecord struct {
	ID                  string            `json:"id"`
	Type                string            `json:"type"`
	Name                string            `json:"name"`
	DefaultCapabilities []AgentCapability `json:"defaultCapabilities"`
	Goal                string            `json:"goal"`
}

type AgentCapability struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}
