package enum

type Events string

const (
	EventWebtrackerCreated  Events = "webtracker.created"
	EventWebtrackerUpdated  Events = "webtracker.updated"
	EventWebtrackerArchived Events = "webtracker.archived"

	EventProxyWebtrackerCnameConfigured    Events = "proxy.webtracker.cname.configured"
	EventProxyWebtrackerCnameNotConfigured Events = "proxy.webtracker.cname.not_configured"
	EventProxyWebtrackerActivated          Events = "proxy.webtracker.activated"
	EventProxyWebtrackerDeactivated        Events = "proxy.webtracker.deactivated"

	EventWebtrackerSessionCreated  Events = "webtracker.session.created"
	EventWebtrackerSessionClosed   Events = "webtracker.session.closed"
	EventWebtrackerSessionAnalyzed Events = "webtracker.session.analyzed"

	EventWebtrackerVisitorIdentified Events = "webtracker.visitor.identified"

	EventWebtrackerPageView Events = "webtracker.event.page_viewed"
	EventWebtrackerPageExit Events = "webtracker.event.page_exited"
	EventWebtrackerClick    Events = "webtracker.event.clicked"

	EventLeadCreated Events = "lead.created"

	EventAskIPData   Events = "ipaddress.verify.ipdata"
	EventAskSnitcher Events = "ipaddress.identify.snitcher"
)

func (e Events) String() string {
	return string(e)
}
