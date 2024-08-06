package model

type RancherEmailResponseDTO struct {
	Input       string `json:"input"`
	IsReachable string `json:"is_reachable"`
	Misc        struct {
		IsDisposable   bool        `json:"is_disposable"`
		IsRoleAccount  bool        `json:"is_role_account"`
		GravatarUrl    interface{} `json:"gravatar_url"`
		Haveibeenpwned interface{} `json:"haveibeenpwned"`
	} `json:"misc"`
	Mx struct {
		AcceptsMail bool     `json:"accepts_mail"`
		Records     []string `json:"records"`
	} `json:"mx"`
	Smtp struct {
		CanConnectSmtp bool `json:"can_connect_smtp"`
		HasFullInbox   bool `json:"has_full_inbox"`
		IsCatchAll     bool `json:"is_catch_all"`
		IsDeliverable  bool `json:"is_deliverable"`
		IsDisabled     bool `json:"is_disabled"`
	} `json:"smtp"`
	Syntax struct {
		Address         string      `json:"address"`
		Domain          string      `json:"domain"`
		IsValidSyntax   bool        `json:"is_valid_syntax"`
		Username        string      `json:"username"`
		NormalizedEmail string      `json:"normalized_email"`
		Suggestion      interface{} `json:"suggestion"`
	} `json:"syntax"`
}
