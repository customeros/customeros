package dto

type AddDomain struct {
	Domain string `json:"domain"`
}

func NewAddDomainEvent(domain string) AddDomain {
	output := AddDomain{
		Domain: domain,
	}
	return output
}
