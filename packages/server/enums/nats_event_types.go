package enums

// NatsEventType represents a NATS event
type NatsEventType string

const (
	// Core events
	EventTenantCreated NatsEventType = "core.tenant.created"

	// Organization events
	EventCreateOrganization        NatsEventType = "organization.create"
	EventUpdateOrganization        NatsEventType = "organization.update"
	EventRequestEnrichOrganization NatsEventType = "organization.request_enrich"
	EventAddDomain                 NatsEventType = "organization.add_domain"
	EventRemoveDomain              NatsEventType = "organization.remove_domain"

	// Contact events
	EventCreateContact        NatsEventType = "contact.create"
	EventUpdateContact        NatsEventType = "contact.update"
	EventRequestEnrichContact NatsEventType = "contact.request_enrich"

	// Web-related events
	EventWebsiteCrawled    NatsEventType = "website.crawled"
	EventWebpageScraped    NatsEventType = "webpage.scraped"
	EventWebpageClassified NatsEventType = "webpage.classified"
	EventWebpageProfiled   NatsEventType = "webpage.profiled"

	// AI and analysis events
	EventAIRequestGeneric             NatsEventType = "ai.request.generic"
	EventRequestWebpageClassification NatsEventType = "ai.request.webpage_classification"
	EventRequestWebpageIntent         NatsEventType = "ai.request.webpage_intent"

	// ICP Profile events
	EventRequestICPProfile NatsEventType = "request.icp_profile"
	EventICPProfileCreated NatsEventType = "icp_profile.created"

	// WebTracker - main events
	EventWebtrackerCreated  NatsEventType = "webtracker.created"
	EventWebtrackerUpdated  NatsEventType = "webtracker.updated"
	EventWebtrackerArchived NatsEventType = "webtracker.archived"

	// WebTracker - session events
	EventWebtrackerSessionCreated  NatsEventType = "webtracker.session.created"
	EventWebtrackerSessionClosed   NatsEventType = "webtracker.session.closed"
	EventWebtrackerSessionAnalyzed NatsEventType = "webtracker.session.analyzed"

	// WebTracker - visitor events
	EventWebtrackerVisitorIdentified NatsEventType = "webtracker.visitor.identified"

	// WebTracker - page events
	EventWebtrackerPageView NatsEventType = "webtracker.event.page_viewed"
	EventWebtrackerPageExit NatsEventType = "webtracker.event.page_exited"
	EventWebtrackerClick    NatsEventType = "webtracker.event.clicked"

	// WebTracker - proxy configuration
	EventProxyWebtrackerCnameConfigured    NatsEventType = "proxy.webtracker.cname.configured"
	EventProxyWebtrackerCnameNotConfigured NatsEventType = "proxy.webtracker.cname.not_configured"
	EventProxyWebtrackerActivated          NatsEventType = "proxy.webtracker.activated"
	EventProxyWebtrackerDeactivated        NatsEventType = "proxy.webtracker.deactivated"

	// Lead events
	EventLeadIdentified               NatsEventType = "lead.identified"
	EventLeadStageUpdate              NatsEventType = "lead.update.stage"
	EventLeadInitialTargetListCreated NatsEventType = "lead.initial_target_list.created"
	EventLeadError                    NatsEventType = "lead.error"

	// IP Address verification events
	EventAskIPData   NatsEventType = "request.verify_ipaddress.ipdata"
	EventAskSnitcher NatsEventType = "request.identify_ipaddress.snitcher"
)

// String returns the string representation of the NatsEventType
func (e NatsEventType) String() string {
	return string(e)
}
