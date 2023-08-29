package mapper

import (
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
)

func MapRenewalLikelihoodToModels(likelihood organization_grpc_service.Likelihood) models.RenewalLikelihoodProbability {
	switch likelihood {
	case organization_grpc_service.Likelihood_HIGH:
		return models.RenewalLikelihoodHIGH
	case organization_grpc_service.Likelihood_MEDIUM:
		return models.RenewalLikelihoodMEDIUM
	case organization_grpc_service.Likelihood_LOW:
		return models.RenewalLikelihoodLOW
	case organization_grpc_service.Likelihood_ZERO:
		return models.RenewalLikelihoodZERO
	}

	return ""
}

func MapRenewalLikelihoodToGraphDb(likelihood models.RenewalLikelihoodProbability) entity.RenewalLikelihoodProbability {
	switch likelihood {
	case models.RenewalLikelihoodHIGH:
		return entity.RenewalLikelihoodHigh
	case models.RenewalLikelihoodMEDIUM:
		return entity.RenewalLikelihoodMedium
	case models.RenewalLikelihoodLOW:
		return entity.RenewalLikelihoodLow
	case models.RenewalLikelihoodZERO:
		return entity.RenewalLikelihoodZero
	}

	return ""
}

func MapRenewalLikelihoodFromGraphDb(likelihood entity.RenewalLikelihoodProbability) models.RenewalLikelihoodProbability {
	switch likelihood {
	case entity.RenewalLikelihoodHigh:
		return models.RenewalLikelihoodHIGH
	case entity.RenewalLikelihoodMedium:
		return models.RenewalLikelihoodMEDIUM
	case entity.RenewalLikelihoodLow:
		return models.RenewalLikelihoodLOW
	case entity.RenewalLikelihoodZero:
		return models.RenewalLikelihoodZERO
	}

	return ""
}

func MapFrequencyToString(frequency *organization_grpc_service.Frequency) string {
	if frequency == nil {
		return ""
	}
	switch *frequency {
	case organization_grpc_service.Frequency_WEEKLY:
		return "WEEKLY"
	case organization_grpc_service.Frequency_BIWEEKLY:
		return "BIWEEKLY"
	case organization_grpc_service.Frequency_MONTHLY:
		return "MONTHLY"
	case organization_grpc_service.Frequency_QUARTERLY:
		return "QUARTERLY"
	case organization_grpc_service.Frequency_BIANNUALLY:
		return "BIANNUALLY"
	case organization_grpc_service.Frequency_ANNUALLY:
		return "ANNUALLY"
	}

	return ""
}
