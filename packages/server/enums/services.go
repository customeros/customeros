package enums

type Services string

const (
	ServicesICP          Services = "core-crm-icp"
	ServicesOrganization Services = "organization"
)

func (s Services) String() string {
	return string(s)
}
