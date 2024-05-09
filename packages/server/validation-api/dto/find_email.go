package dto

type FindEmailRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Domain    string `json:"domain"`
}

type FindEmailResponse struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Score     float64 `json:"score"`
}
