package model

type WorkspaceInput struct {
	Name      string  `json:"name"`
	Provider  string  `json:"provider"`
	AppSource *string `json:"appSource"`
}

type EmailInput struct {
	Email     string  `json:"email"`
	Primary   bool    `json:"primary"`
	AppSource *string `json:"appSource"`
}

type PlayerInput struct {
	IdentityId string  `json:"identityId"`
	AuthId     string  `json:"authId"`
	Provider   string  `json:"provider"`
	AppSource  *string `json:"appSource"`
}

type UserInput struct {
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     EmailInput  `json:"email"`
	Person    PlayerInput `json:"player"`
}

type TenantInput struct {
	Name      string  `json:"name"`
	AppSource *string `json:"appSource"`
}
