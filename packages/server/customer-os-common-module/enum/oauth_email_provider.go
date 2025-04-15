package enum

import (
	"fmt"
)

type OAuthEmailProvider string

const (
	ProviderGoogle  OAuthEmailProvider = "gmail" // TODO replace gmail with google_workspace, also update prod db
	ProviderOutlook OAuthEmailProvider = "outlook"
)

func (p OAuthEmailProvider) String() string {
	return string(p)
}

func GetOAuthEmailProvider(s string) (OAuthEmailProvider, error) {
	switch OAuthEmailProvider(s) {
	case ProviderGoogle, ProviderOutlook:
		return OAuthEmailProvider(s), nil
	default:
		return "", fmt.Errorf("invalid OAuthEmailProvider: %s", s)
	}
}
