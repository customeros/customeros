package enum

type EnrichSource string

const (
	Brandfetch EnrichSource = "BRANDFETCH"
	ScrapIn    EnrichSource = "SCRAPIN"
)

func (e EnrichSource) String() string {
	return string(e)
}

func DecodeDomainEnrichSource(str string) EnrichSource {
	switch str {
	case Brandfetch.String():
		return Brandfetch
	case ScrapIn.String():
		return ScrapIn
	default:
		return ""
	}
}
