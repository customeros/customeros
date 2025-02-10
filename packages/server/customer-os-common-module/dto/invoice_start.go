package dto

type InvoiceStart struct {
	AgentExecutionId string `json:"agentExecutionId"`
	ContractID       string `json:"contractId"`
	DryRun           bool   `json:"dryRun"`
	Preview          bool   `json:"preview"`
}
