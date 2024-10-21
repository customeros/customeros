package dto

type AddEmail struct {
	Email string `json:"email"`
}

func NewAddEmailEvent(email string) AddEmail {
	output := AddEmail{
		Email: email,
	}
	return output
}
