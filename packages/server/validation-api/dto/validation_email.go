package dto

type ValidationEmailRequest struct {
	Email string `json:"email"`
}
type ValidationEmailResponse struct {
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

func MapValidationEmailResponse(reacherResponse *RancherEmailResponseDTO, error *string) ValidationEmailResponse {
	return ValidationEmailResponse{
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
