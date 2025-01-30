package dto

type CreateAgent struct {
	Active       bool   `json:"active"`
	VisibleInUI  bool   `json:"visibleInUi"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	Capabilities string `json:"capabilities"`
	Goal         string `json:"goal"`
	Status       string `json:"status"`
	FlowID       string `json:"flowId"`
	RegistryID   string `json:"registryId"`
}
