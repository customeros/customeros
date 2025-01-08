package data_fields

type WebsiteVisitEvent struct {
	ID                  string      `json:"id"`
	Tenant              string      `json:"tenant"`
	IPAddress           *string     `json:"ipAddress"`
	OrgnanizationDomain string      `json:"organizationDomain"`
	OrganizationID      *string     `json:"organizationId"`
	VisitorId           string      `json:"visitorId"`
	VisitorEmail        *string     `json:"visitorEmail"`
	PageVisited         string      `json:"pageVisited"`
	Referrer            string      `json:"referrer"`
	Params              []URLParams `json:"params"`
}

type URLParams struct {
	Param string `json:"param"`
	Value string `json:"value"`
}

func (f WebsiteVisitEvent) Type() string {
	return "WebsiteVisitEvent"
}
