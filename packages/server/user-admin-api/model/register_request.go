package model

type RegisterRequest struct {
	Properties struct {
		FirstName string  `json:"firstname"`
		LastName  string  `json:"lastname"`
		Email     string  `json:"email"`
		Workspace *string `json:"workspace"`
		Provider  string  `json:"provider"`
	} `json:"properties"`
}
