package resolver

type FindEmailRequest struct {
	Domain    string `json:"domain" `
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type FindEmailResponse struct {
	FirsName string  `json:"firstName"`
	LastName string  `json:"lastName"`
	Email    string  `json:"email"`
	Score    float64 `json:"score"`
}
