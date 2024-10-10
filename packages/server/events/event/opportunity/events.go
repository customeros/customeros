package opportunity

const (
	OpportunityCreateV1              = "V1_OPPORTUNITY_CREATE"
	OpportunityUpdateV1              = "V1_OPPORTUNITY_UPDATE"
	OpportunityCreateRenewalV1       = "V1_OPPORTUNITY_CREATE_RENEWAL"
	OpportunityUpdateRenewalV1       = "V1_OPPORTUNITY_UPDATE_RENEWAL"
	OpportunityUpdateNextCycleDateV1 = "V1_OPPORTUNITY_UPDATE_NEXT_CYCLE_DATE"
	OpportunityCloseWinV1            = "V1_OPPORTUNITY_CLOSE_WIN"
	OpportunityCloseLooseV1          = "V1_OPPORTUNITY_CLOSE_LOOSE"
)

const (
	FieldMaskName              = "name"
	FieldMaskAmount            = "amount"
	FieldMaskMaxAmount         = "maxAmount"
	FieldMaskComments          = "comments"
	FieldMaskRenewalLikelihood = "renewalLikelihood"
	FieldMaskRenewalApproved   = "renewalApproved"
	FieldMaskRenewedAt         = "renewedAt"
	FieldMaskAdjustedRate      = "adjustedRate"
	FieldMaskExternalType      = "externalType"
	FieldMaskInternalType      = "internalType"
	FieldMaskExternalStage     = "externalStage"
	FieldMaskInternalStage     = "internalStage"
	FieldMaskEstimatedClosedAt = "estimatedClosedAt"
	FieldMaskOwnerUserId       = "ownerUserId"
	FieldMaskCurrency          = "currency"
	FieldMaskNextSteps         = "nextSteps"
	FieldMaskLikelihoodRate    = "likelihoodRate"
)
