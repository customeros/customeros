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

type ValidationEmailReacherResponse struct {
	Error           *string `json:"error"`
	Email           string  `json:"email"`
	AcceptsMail     bool    `json:"acceptsMail"`
	CanConnectSmtp  bool    `json:"canConnectSmtp"`
	HasFullInbox    bool    `json:"hasFullInbox"`
	IsCatchAll      bool    `json:"isCatchAll"`
	IsDeliverable   bool    `json:"isDeliverable"`
	IsDisabled      bool    `json:"isDisabled"`
	IsReachable     string  `json:"isReachable"`
	Address         string  `json:"address"`
	Domain          string  `json:"domain"`
	IsValidSyntax   bool    `json:"isValidSyntax"`
	Username        string  `json:"username"`
	NormalizedEmail string  `json:"normalizedEmail"`
	IsDisposable    bool    `json:"isDisposable"`
	IsRoleAccount   bool    `json:"isRoleAccount"`
}

func MapValidationEmailResponse(reacherResponse *RancherEmailResponseDTO, error *string) ValidationEmailReacherResponse {
	return ValidationEmailReacherResponse{
		Error:           error,
		Email:           reacherResponse.Input,
		IsReachable:     reacherResponse.IsReachable,
		AcceptsMail:     reacherResponse.Mx.AcceptsMail,
		CanConnectSmtp:  reacherResponse.Smtp.CanConnectSmtp,
		HasFullInbox:    reacherResponse.Smtp.HasFullInbox,
		IsCatchAll:      reacherResponse.Smtp.IsCatchAll,
		IsDeliverable:   reacherResponse.Smtp.IsDeliverable,
		IsDisabled:      reacherResponse.Smtp.IsDisabled,
		Address:         reacherResponse.Syntax.Address,
		Domain:          reacherResponse.Syntax.Domain,
		IsValidSyntax:   reacherResponse.Syntax.IsValidSyntax,
		Username:        reacherResponse.Syntax.Username,
		NormalizedEmail: reacherResponse.Syntax.NormalizedEmail,
		IsDisposable:    reacherResponse.Misc.IsDisposable,
		IsRoleAccount:   reacherResponse.Misc.IsRoleAccount,
	}
}
