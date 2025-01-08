package data_fields

type WebsiteVisitEvent struct {
	ID           string      `json:"id"`
	Tenant       string      `json:"tenant"`
	IPAddress    *string     `json:"ipAddress"`
	VisitorId    string      `json:"visitorId"`
	VisitorEmail *string     `json:"visitorEmail"`
	Website      string      `json:"website"`
	PageVisited  string      `json:"pageVisited"`
	Referrer     string      `json:"referrer"`
	Params       []URLParams `json:"params"`
}

type URLParams struct {
	Param string `json:"param"`
	Value string `json:"value"`
}

func (f WebsiteVisitEvent) Type() string {
	return "WebsiteVisitEvent"
}
