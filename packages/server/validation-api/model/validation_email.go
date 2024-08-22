package model

type ValidateEmailRequest struct {
	Email string `json:"email"`
}

type ValidateEmailResponse struct {
	Status          string                       `json:"status"`
	Message         string                       `json:"message,omitempty"`
	InternalMessage string                       `json:"internalMessage,omitempty"`
	Data            *ValidateEmailMailsherpaData `json:"data,omitempty"`
}

type ValidateEmailMailsherpaData struct {
	Email  string `json:"email"`
	Syntax struct {
		IsValid    bool   `json:"isValid"`
		User       string `json:"user"`
		Domain     string `json:"domain"`
		CleanEmail string `json:"cleanEmail"`
	} `json:"syntax"`
	DomainData struct {
		IsFirewalled   bool   `json:"isFirewalled"`
		Provider       string `json:"provider"`
		Firewall       string `json:"firewall"`
		IsCatchAll     bool   `json:"isCatchAll"`
		CanConnectSMTP bool   `json:"canConnectSMTP"`
	} `json:"domainData"`
	EmailData struct {
		IsDeliverable   bool   `json:"isDeliverable"`
		IsMailboxFull   bool   `json:"isMailboxFull"`
		IsRoleAccount   bool   `json:"isRoleAccount"`
		IsFreeAccount   bool   `json:"isFreeAccount"`
		SmtpSuccess     bool   `json:"smtpSuccess"`
		ResponseCode    string `json:"responseCode"`
		ErrorCode       string `json:"errorCode"`
		Description     string `json:"description"`
		RetryValidation bool   `json:"retryValidation"`
		SmtpResponse    string `json:"smtpResponse"`
	} `json:"emailData"`
}
