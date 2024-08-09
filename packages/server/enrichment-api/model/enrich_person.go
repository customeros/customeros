package model

type EnrichPersonRequest struct {
	Email       string `json:"ip"`
	LinkedInUrl string `json:"linkedinUrl"`
}

type EnrichPersonResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Data    *EnrichPersonData `json:"data,omitempty"`
}

type EnrichPersonData struct {
}
