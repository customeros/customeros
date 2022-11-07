package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// CompanyPositions is the resolver for the companyPositions field.
func (r *contactResolver) CompanyPositions(ctx context.Context, obj *model.Contact) ([]*model.CompanyPosition, error) {
	companyPositionEntities, err := r.ServiceContainer.CompanyPositionService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToCompanyPositiones(companyPositionEntities), err
}

// Groups is the resolver for the groups field.
func (r *contactResolver) Groups(ctx context.Context, obj *model.Contact) ([]*model.ContactGroup, error) {
	contactGroupEntities, err := r.ServiceContainer.ContactGroupService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToContactGroups(contactGroupEntities), err
}

// TextCustomFields is the resolver for the textCustomFields field.
func (r *contactResolver) TextCustomFields(ctx context.Context, obj *model.Contact) ([]*model.TextCustomField, error) {
	textCustomFieldEntities, err := r.ServiceContainer.TextCustomFieldService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToTextCustomFields(textCustomFieldEntities), err
}

// PhoneNumbers is the resolver for the phoneNumbers field.
func (r *contactResolver) PhoneNumbers(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumberInfo, error) {
	phoneNumberEntities, err := r.ServiceContainer.PhoneNumberService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToPhoneNumbers(phoneNumberEntities), err
}

// Emails is the resolver for the emails field.
func (r *contactResolver) Emails(ctx context.Context, obj *model.Contact) ([]*model.EmailInfo, error) {
	emailEntities, err := r.ServiceContainer.EmailService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToEmails(emailEntities), err
}

// FieldsSets is the resolver for the fieldsSets field.
func (r *contactResolver) FieldsSets(ctx context.Context, obj *model.Contact) ([]*model.FieldsSet, error) {
	panic(fmt.Errorf("not implemented: FieldsSets - fieldsSets"))
}

// Contact returns generated.ContactResolver implementation.
func (r *Resolver) Contact() generated.ContactResolver { return &contactResolver{r} }

type contactResolver struct{ *Resolver }
