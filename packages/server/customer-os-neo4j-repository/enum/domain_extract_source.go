package enum

type DomainEnrichSource string

const (
	Brandfetch DomainEnrichSource = "BRANDFETCH"
)

func (e DomainEnrichSource) String() string {
	return string(e)
}

func DecodeDomainEnrichSource(str string) DomainEnrichSource {
	switch str {
	case Brandfetch.String():
		return Brandfetch
	default:
		return ""
	}
}

func (e DomainEnrichSource) IsValid() bool {
	switch e {
	case Brandfetch:
		return true
	}
	return false
}
