package enums

// NatsEventType represents a NATS event
type NatsEventType string

const (
	// tenant events
	EventTenantCreated NatsEventType = "tenant.created"

	// Organization events
	EventOrganizationCreated       NatsEventType = "organization.created"
	EventOrganizationUpdated       NatsEventType = "organization.updated"
	EventOrganizationDomainAdded   NatsEventType = "organization.domain.added"
	EventOrganizationDomainRemoved NatsEventType = "organization.domain.removed"
	EventRequestEnrichOrganization NatsEventType = "request.organization.enrich"

	// Contact events
	EventContactCreated       NatsEventType = "contact.created"
	EventContactUpdated       NatsEventType = "contact.updated"
	EventRequestEnrichContact NatsEventType = "request.contact.enrich"

	// Web-related events
	EventWebsiteCrawled    NatsEventType = "web.site.crawled"
	EventWebpageScraped    NatsEventType = "web.page.scraped"
	EventWebpageClassified NatsEventType = "web.page.classified"
	EventWebpageProfiled   NatsEventType = "web.page.profiled"

	// AI and analysis events
	EventRequestWebpageClassification NatsEventType = "request.webpage_classification"
	EventRequestWebpageIntent         NatsEventType = "request.webpage_intent"

	// ICP Profile events
	EventIdealCustomerProfileCreated NatsEventType = "icp.created"
	EventRequestIdealCustomerProfile NatsEventType = "request.icp"

	// WebTracker
	EventWebtrackerCreated           NatsEventType = "webtracker.created"
	EventWebtrackerUpdated           NatsEventType = "webtracker.updated"
	EventWebtrackerArchived          NatsEventType = "webtracker.archived"
	EventWebtrackerSessionCreated    NatsEventType = "webtracker.session.created"
	EventWebtrackerSessionClosed     NatsEventType = "webtracker.session.closed"
	EventWebtrackerSessionAnalyzed   NatsEventType = "webtracker.session.analyzed"
	EventWebtrackerVisitorIdentified NatsEventType = "webtracker.visitor.identified"
	EventWebtrackerPageView          NatsEventType = "webtracker.event.page_viewed"
	EventWebtrackerPageExit          NatsEventType = "webtracker.event.page_exited"
	EventWebtrackerClick             NatsEventType = "webtracker.event.clicked"

	// WebTracker - proxy configuration
	EventProxyWebtrackerCnameConfigured    NatsEventType = "proxy.webtracker.cname.configured"
	EventProxyWebtrackerCnameNotConfigured NatsEventType = "proxy.webtracker.cname.not_configured"
	EventProxyWebtrackerActivated          NatsEventType = "proxy.webtracker.activated"
	EventProxyWebtrackerDeactivated        NatsEventType = "proxy.webtracker.deactivated"

	// Lead events
	EventLeadIdentified               NatsEventType = "lead.identified"
	EventLeadStageUpdate              NatsEventType = "lead.stage.updated"
	EventLeadInitialTargetListCreated NatsEventType = "lead.initial_target_list.created"

	// IP Address verification events
	EventAskIPData   NatsEventType = "request.verify_ipaddress.ipdata"
	EventAskSnitcher NatsEventType = "request.identify_ipaddress.snitcher"

	// DLQ
	EventDLQContact      NatsEventType = "dlq.contact"
	EventDLQICP          NatsEventType = "dlq.icp"
	EventDLQLead         NatsEventType = "dlq.lead"
	EventDLQOrganization NatsEventType = "dlq.organization"
	EventDLQProxy        NatsEventType = "dlq.proxy"
	EventDLQRequest      NatsEventType = "dlq.request"
	EventDLQTenant       NatsEventType = "dlq.tenant"
	EventDLQWeb          NatsEventType = "dlq.web"
	EventDLQWebtracker   NatsEventType = "dlq.webtracker"
)

// String returns the string representation of the NatsEventType
func (e NatsEventType) String() string {
	return string(e)
}
