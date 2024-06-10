package resolver

import neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

func OrganizationStageAndRelationshipCompatible(stageStr, relationshipStr string) bool {
	stage := neo4jenum.OrganizationStage(stageStr)
	relationship := neo4jenum.OrganizationRelationship(relationshipStr)

	if stage == "" || relationship == "" {
		return true
	}

	if relationship == neo4jenum.NotAFit && stage != neo4jenum.Unqualified {
		return false
	} else if relationship == neo4jenum.FormerCustomer && stage != neo4jenum.Target {
		return false
	} else if relationship == neo4jenum.Prospect && stage != neo4jenum.Lead && stage != neo4jenum.Target &&
		stage != neo4jenum.Engaged && stage != neo4jenum.ReadyToBuy && stage != neo4jenum.Trial {
		return false
	} else if relationship == neo4jenum.Customer && stage != neo4jenum.Onboarding && stage != neo4jenum.InitialValue &&
		stage != neo4jenum.RecurringValue && stage != neo4jenum.MaxValue && stage != neo4jenum.PendingChurn {
		return false
	}
	return true
}
