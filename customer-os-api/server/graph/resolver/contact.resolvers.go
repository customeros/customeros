package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// Groups is the resolver for the groups field.
func (r *contactResolver) Groups(ctx context.Context, obj *model.Contact) ([]*model.ContactGroup, error) {
	contactGroupEntities, err := r.ServiceContainer.ContactGroupService.FindAllForContact(obj)
	return mapper.MapEntitiesToContactGroups(contactGroupEntities), err
}

// Contact returns generated.ContactResolver implementation.
func (r *Resolver) Contact() generated.ContactResolver { return &contactResolver{r} }

type contactResolver struct{ *Resolver }
