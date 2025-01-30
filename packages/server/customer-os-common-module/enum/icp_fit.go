package enum

type IcpFit string

const (
	IcpIsFit  IcpFit = "ICP_FIT"
	IcpNotFit IcpFit = "ICP_NOT_FIT"
	IcpNotSet IcpFit = "ICP_NOT_SET"
)

var AllIcpFits = []IcpFit{
	IcpIsFit,
	IcpNotFit,
	IcpNotSet,
}

func DecodeIcpFit(s string) IcpFit {
	if IsValidIcpFit(s) {
		return IcpFit(s)
	}
	return IcpNotSet
}

func IsValidIcpFit(s string) bool {
	for _, icp := range AllIcpFits {
		if icp == IcpFit(s) {
			return true
		}
	}
	return false
}

func (i IcpFit) String() string {
	return string(i)
}
