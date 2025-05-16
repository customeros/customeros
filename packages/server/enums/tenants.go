package enums

type Tenant string

const (
	TenantCustomerOS Tenant = "customerosai"
)

func (s Tenant) String() string {
	return string(s)
}
