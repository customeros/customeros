package agent_capability

type NoOutput struct{}

type ConfigSingleValue struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

type ConfigSingleBoolValue struct {
	Value bool   `json:"value"`
	Error string `json:"error"`
}

type ConfigMultipleValues struct {
	Value []string `json:"value"`
	Error string   `json:"error"`
}
