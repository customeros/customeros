package agent_capability

type CapabilityOutput struct {
	ExecutionValidated bool `json:"executionValidated"`
	Completed          bool `json:"completed"`
}
