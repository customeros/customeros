package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j/entity"
)

func MapEntityToState(state *entity.StateEntity) *model.State {
	return &model.State{
		ID:   state.Id,
		Name: state.Name,
		Code: state.Code,
	}
}
