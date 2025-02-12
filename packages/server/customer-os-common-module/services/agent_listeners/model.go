package agent_listeners

type ConfigSingleValue struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

type ConfigMultipleValues struct {
	Value []string `json:"value"`
	Error string   `json:"error"`
}
