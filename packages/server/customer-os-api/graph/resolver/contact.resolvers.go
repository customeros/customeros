package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// Companies is the resolver for the companies field.
func (r *contactResolver) Companies(ctx context.Context, obj *model.Contact) ([]*model.Company, error) {
	companyPositionEntities, err := r.ServiceContainer.CompanyPositionService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToCompanyPositiones(companyPositionEntities), err
}

// Groups is the resolver for the groups field.
func (r *contactResolver) Groups(ctx context.Context, obj *model.Contact) ([]*model.ContactGroup, error) {
	contactGroupEntities, err := r.ServiceContainer.ContactGroupService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToContactGroups(contactGroupEntities), err
}

// PhoneNumbers is the resolver for the phoneNumbers field.
func (r *contactResolver) PhoneNumbers(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error) {
	phoneNumberEntities, err := r.ServiceContainer.PhoneNumberService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToPhoneNumbers(phoneNumberEntities), err
}

// Emails is the resolver for the emails field.
func (r *contactResolver) Emails(ctx context.Context, obj *model.Contact) ([]*model.Email, error) {
	emailEntities, err := r.ServiceContainer.EmailService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToEmails(emailEntities), err
}

// CustomFields is the resolver for the customFields field.
func (r *contactResolver) CustomFields(ctx context.Context, obj *model.Contact) ([]*model.CustomField, error) {
	var customFields []*model.CustomField
	textCustomFieldEntities, err := r.ServiceContainer.TextCustomFieldService.FindAllForContact(ctx, obj)
	for _, v := range mapper.MapEntitiesToTextCustomFields(textCustomFieldEntities) {
		customFields = append(customFields, v)
	}
	return customFields, err
}

// FieldSets is the resolver for the fieldSets field.
func (r *contactResolver) FieldSets(ctx context.Context, obj *model.Contact) ([]*model.FieldSet, error) {
	fieldSetEntities, err := r.ServiceContainer.FieldSetService.FindAllForContact(ctx, obj)
	return mapper.MapEntitiesToFieldSets(fieldSetEntities), err
}

// Definition is the resolver for the definition field.
func (r *contactResolver) Definition(ctx context.Context, obj *model.Contact) (*model.EntityDefinition, error) {
	entity, err := r.ServiceContainer.EntityDefinitionService.FindLinkedWithContact(ctx, obj.ID)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to get contact definition for contact ", obj.ID)
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	return mapper.MapEntityToEntityDefinition(entity), err
}

// Contact returns generated.ContactResolver implementation.
func (r *Resolver) Contact() generated.ContactResolver { return &contactResolver{r} }

type contactResolver struct{ *Resolver }
