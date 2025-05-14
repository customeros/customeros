package enum

type QualificationStatus string

const (
	QualificationStatusQualifying   QualificationStatus = "QUALIFYING"
	QualificationStatusQualified    QualificationStatus = "QUALIFIED"
	QualificationStatusNotQualified QualificationStatus = "UNQUALIFIED"
)

func (e QualificationStatus) String() string {
	return string(e)
}

func DecodeQualificationStatus(str string) QualificationStatus {
	switch str {
	case QualificationStatusQualifying.String():
		return QualificationStatusQualifying
	case QualificationStatusQualified.String():
		return QualificationStatusQualified
	case QualificationStatusNotQualified.String():
		return QualificationStatusNotQualified
	default:
		return ""
	}
}
