package agent_listeners

type ConfigSingleValue struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

type ConfigMultipleValues struct {
	Value []string `json:"value"`
	Error string   `json:"error"`
}

type ConfigMultipleValuesWithObject struct {
	Value []interface{} `json:"value"`
	Error string        `json:"error"`
}

type ConfigSingleIntValue struct {
	Value int64  `json:"value"`
	Error string `json:"error"`
}
