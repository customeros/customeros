package dto

type RemoveEmail struct {
	Email string `json:"email"`
}

func NewRemoveEmailEvent(email string) RemoveEmail {
	output := RemoveEmail{
		Email: email,
	}
	return output
}
