package agent_capability

type NoConfig struct{}

type CapabilityOutput struct {
	ExecutionValidated bool `json:"executionValidated"`
	Completed          bool `json:"completed"`
}
