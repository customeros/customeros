package model

type OrganizationsResponse struct {
	Organizations struct {
		Content []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"content"`
	} `json:"organizations"`
}
