package interfaces

import "context"

type VerifyService interface {
	Threats(ctx context.Context, ipAddress string) (*IpThreats, error)
	IdentifyCompanyDomain(ctx context.Context, ipAddress string) (*string, error)
}

type IpThreats struct {
	IsThreat      bool
	IsAnonymous   bool
	IsBogon       bool
	IsDatacenter  bool
	IsICloudRelay bool
	IsKnownAbuser bool
	IsProxy       bool
	IsTor         bool
	IsVpn         bool
}
