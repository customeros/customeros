package enum

type QualificationStatus string

const (
	QualificationStatusPending      QualificationStatus = "PENDING"
	QualificationStatusQualifying   QualificationStatus = "QUALIFYING"
	QualificationStatusQualified    QualificationStatus = "QUALIFIED"
	QualificationStatusNotQualified QualificationStatus = "UNQUALIFIED"
)

func (e QualificationStatus) String() string {
	return string(e)
}

func (e QualificationStatus) Order() int64 {
	switch e {
	case QualificationStatusPending:
		return 0
	case QualificationStatusQualifying:
		return 1
	case QualificationStatusQualified:
		return 2
	case QualificationStatusNotQualified:
		return 3
	default:
		return 0
	}
}

func DecodeQualificationStatus(str string) QualificationStatus {
	switch str {
	case QualificationStatusPending.String():
		return QualificationStatusPending
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
