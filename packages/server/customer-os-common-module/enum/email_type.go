package enum

type EmailType string

const (
	EmailPersonal EmailType = "personal"
	EmailBusiness EmailType = "business"
)

func (t EmailType) String() string {
	return string(t)
}
