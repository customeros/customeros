package enum

type Events string

const (
	EventWebtrackerCreated  Events = "webtracker.created"
	EventWebtrackerUpdated  Events = "webtracker.updated"
	EventWebtrackerArchived Events = "webtracker.archived"

	EventProxyWebtrackerCnameConfigured    Events = "proxy.webtracker.cname.configured"
	EventProxyWebtrackerCnameNotConfigured Events = "proxy.webtracker.cname.not_configured"
	EventProxyWebtrackerActive             Events = "proxy.webtracker.active"
	EventProxyWebtrackerNotActive          Events = "proxy.webtracker.not_active"

	EventWebtrackerSessionNew    Events = "webtracker.session.new"
	EventWebtrackerSessionClosed Events = "webtracker.session.closed"

	EventWebtrackerVisitorIdentified Events = "webtracker.visitor.identified"

	EventWebtrackerPageView Events = "webtracker.event.page_view"
	EventWebtrackerPageExit Events = "webtracker.event.page_exit"
	EventWebtrackerClick    Events = "webtracker.event.click"

	EventLeadCreated Events = "lead.created"

	EventAskIPData   Events = "ipaddress.verify.ipdata"
	EventAskSnitcher Events = "ipaddress.identify.snitcher"
)

func (e Events) String() string {
	return string(e)
}
