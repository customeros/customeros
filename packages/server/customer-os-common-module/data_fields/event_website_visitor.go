package data_fields

type WebsiteVisitEventType string

const (
	WebsiteVisitPageExit WebsiteVisitEventType = "page_exit"
	WebsiteVisitPageView WebsiteVisitEventType = "page_view"
	WebsiteVisitClice    WebsiteVisitEventType = "click"
)

type WebsiteVisitEvent struct {
	ID              string                `json:"id"`
	Tenant          string                `json:"tenant"`
	IPAddress       *string               `json:"ipAddress"`
	EventType       WebsiteVisitEventType `json:"eventType"`
	VisitorId       string                `json:"visitorId"`
	VisitorEmail    *string               `json:"visitorEmail"`
	Website         string                `json:"website"`
	PageVisited     string                `json:"pageVisited"`
	Referrer        string                `json:"referrer"`
	MinutesDuration string                `json:"minutesDuration"`
	Params          []URLParams           `json:"params"`
}

type URLParams struct {
	Param string `json:"param"`
	Value string `json:"value"`
}

func (f WebsiteVisitEvent) Type() string {
	return "WebsiteVisitEvent"
}
