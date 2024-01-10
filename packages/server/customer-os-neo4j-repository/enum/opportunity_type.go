package enum

type OpportunityInternalType string

const (
	OpportunityInternalTypeNBO       OpportunityInternalType = "NBO"
	OpportunityInternalTypeUpsell    OpportunityInternalType = "UPSELL"
	OpportunityInternalTypeCrossSell OpportunityInternalType = "CROSS_SELL"
	OpportunityInternalTypeRenewal   OpportunityInternalType = "RENEWAL"
)

var AllOpportunityInternalType = []OpportunityInternalType{
	OpportunityInternalTypeNBO,
	OpportunityInternalTypeUpsell,
	OpportunityInternalTypeCrossSell,
	OpportunityInternalTypeRenewal,
}

func (e OpportunityInternalType) IsValid() bool {
	switch e {
	case OpportunityInternalTypeNBO, OpportunityInternalTypeUpsell, OpportunityInternalTypeCrossSell, OpportunityInternalTypeRenewal:
		return true
	}
	return false
}

func (e OpportunityInternalType) String() string {
	return string(e)
}
