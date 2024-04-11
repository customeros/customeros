package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToState(state *neo4jEntity.StateEntity) *model.State {
	return &model.State{
		ID:   state.Id,
		Name: state.Name,
		Code: state.Code,
	}
}
