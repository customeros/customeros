package enum

type QualifiedBy string

const (
	QualifiedBySystem QualifiedBy = "SYSTEM"
	QualifiedByUser   QualifiedBy = "USER"
)

func (e QualifiedBy) String() string {
	return string(e)
}

func DecodeQualifiedBy(str string) QualifiedBy {
	switch str {
	case QualifiedBySystem.String():
		return QualifiedBySystem
	case QualifiedByUser.String():
		return QualifiedByUser
	default:
		return ""
	}
}
