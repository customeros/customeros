package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapExternalSystemEntitiesToExternalSystemInstances(entities *neo4jentity.ExternalSystemEntities) []*model.ExternalSystemInstance {
	var instances []*model.ExternalSystemInstance
	for _, entity := range *entities {
		instances = append(instances, MapExternalSystemEntityToExternalSystemInstance(&entity))
	}
	return instances

}

func MapExternalSystemEntityToExternalSystemInstance(e *neo4jentity.ExternalSystemEntity) *model.ExternalSystemInstance {
	if e == nil {
		return nil
	}
	externalSystemInstance := model.ExternalSystemInstance{
		Type: MapExternalSystemTypeToModel(e.ExternalSystemId),
	}
	if externalSystemInstance.Type == model.ExternalSystemTypeStripe {
		externalSystemInstance.StripeDetails = &model.ExternalSystemStripeDetails{
			PaymentMethodTypes: e.Stripe.PaymentMethodTypes,
		}
	}
	return &externalSystemInstance
}
