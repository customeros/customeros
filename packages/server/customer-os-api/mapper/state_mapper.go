package mapper

import (
	neo4jEntity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToState(state *neo4jEntity.StateEntity) *model.State {
	return &model.State{
		ID:   state.Id,
		Name: state.Name,
		Code: state.Code,
	}
}
