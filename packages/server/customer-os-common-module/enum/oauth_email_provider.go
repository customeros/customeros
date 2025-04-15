package enum

type OAuthEmailProvider string

const (
	ProviderGoogle  OAuthEmailProvider = "gmail" // TODO replace gmail with google_workspace, also update prod db
	ProviderOutlook OAuthEmailProvider = "outlook"
)
