package model

type PageView struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Order       int    `json:"order"`
	EngagedTime int    `json:"engagedTime"`
	SessionId   string
}
