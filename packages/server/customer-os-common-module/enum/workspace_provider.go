package enum

type WorkspaceProvider string

const (
	WorkspaceProviderGoogle    WorkspaceProvider = "google"
	WorkspaceProviderAzure     WorkspaceProvider = "azure-ad"
	WorkspaceProviderMagicLink WorkspaceProvider = "magic-link"
)

func DecodeWorkspaceProvider(s string) WorkspaceProvider {
	switch WorkspaceProvider(s) {
	case WorkspaceProviderGoogle, WorkspaceProviderAzure, WorkspaceProviderMagicLink:
		return WorkspaceProvider(s)
	}
	return ""
}

func IsValidWorkspaceProvider(s string) bool {
	switch WorkspaceProvider(s) {
	case WorkspaceProviderGoogle, WorkspaceProviderAzure, WorkspaceProviderMagicLink:
		return true
	}
	return false
}

func (w WorkspaceProvider) String() string {
	return string(w)
}
