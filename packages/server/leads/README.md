# Services

event_processor
monitoring
proxy_manager
session_manager
ssl_certificate
website_registation
lead_generator

# Flow

- logger listens on all events and logs them in data warehouse

- receive request to start tracking domain
  tracker.{domain}.create
- website registration service listens and provisions tracker
  tracker.{domain}.provisioned
- ssl manager serviece listens and provisions ssl
  tracker.{domain}.ssl.ok
- proxy manager service listens and provisions proxy
  tracker.{domain}.proxy.ok

- API is called to register event
- handler validates payload, throws event
  tracker.{domain}.event.page_exit
  tracker.{doamin}.event.page_view
  tracker.{domain}.event.click
  tracker.{domain}.event.identify

- session manager listens for events and attaches them to sessions.
- session closer is a cron that closes sessions
  tracker.{domain}.session.closed

- session analyzer listens for closed sessions and anlyzes them
  -- visitor identification
  -- unique page views
  -- intent signals
  tracker.{domain}.session.analyzed

- receive request to stop tracking domain
  tracker.{domain}.delete

- tracker error
  tracker.error.{domain}

## Events

webtracker.error.{tenant}

webtracker.{tenant}.tracker.create
webtracker.{tenant}.tracker.provisioned
webtracker.{tenant}.proxy.ssl_setup
webtracker.{tenant}.proxy.ready
webtracker.{tenant}.tracker.delete

webtracker.{tenant}.event.page_exit
webtracker.{tenant}.event.page_view
webtracker.{tenant}.event.click
webtracker.{tenant}.event.identify

webtracker.{tenant}.session.closed
webtracker.{tenant}.session.analyzed

webtracker.visitor.identify
