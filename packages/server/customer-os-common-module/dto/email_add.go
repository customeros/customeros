package dto

type AddEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func NewAddEmailEvent(email string, primary bool) AddEmail {
	output := AddEmail{
		Email:   email,
		Primary: primary,
	}
	return output
}
