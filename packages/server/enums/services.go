package enums

type Services string

const (
	ServicesICP Services = "core_crm.icp"
)

func (s Services) String() string {
	return string(s)
}
