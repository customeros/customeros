package entity

type CountryEntity struct {
	Name      string
	CodeA2    string
	CodeA3    string
	PhoneCode string

	DataloaderKey string
}

type CountryEntities []CountryEntity

func (CountryEntity) Labels(string) []string {
	return []string{
		NodeLabel_Country,
	}
}
