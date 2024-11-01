package dto

type RegisterEmail struct {
	Email  string `json:"email"`
	Source string `json:"source"`
}

func NewRegisterEmailEvent(email, source string) RegisterEmail {
	output := RegisterEmail{
		Email:  email,
		Source: source,
	}
	return output
}
