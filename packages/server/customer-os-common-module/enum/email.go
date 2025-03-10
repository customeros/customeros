package enum

type EmailType string

const (
	EmailPersonal EmailType = "personal"
	EmailBusiness EmailType = "business"
)

func (t EmailType) String() string {
	return string(t)
}

type EmailProvider string

const (
	EmailGoogleWorkspace EmailProvider = "google_workspace"
	EmailOutlook         EmailProvider = "outlook"
	EmailMailstack       EmailProvider = "mailstack"
	EmailGeneric         EmailProvider = "generic"
)

func (t EmailProvider) String() string {
	return string(t)
}

type EmailClassification string

const (
	EmailAutoResponder      EmailClassification = "auto_responder"
	EmailBounceNotification EmailClassification = "bounce_notification"
	EmailBulk               EmailClassification = "bulk_email"
	EmailInternal           EmailClassification = "internal"
	EmailOK                 EmailClassification = "ok"
	EmailSensitive          EmailClassification = "sensitive"
	EmailSpam               EmailClassification = "spam"
	EmailWarmer             EmailClassification = "email_warmer"
)

func (t EmailClassification) String() string {
	return string(t)
}
