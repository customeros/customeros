package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// ContactType is the resolver for the contactType field.
func (r *contactResolver) ContactType(ctx context.Context, obj *model.Contact) (*model.ContactType, error) {
	entity, err := r.ServiceContainer.ContactTypeService.FindContactTypeForContact(ctx, obj.ID)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to get contact type for contact %s", obj.ID)
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	return mapper.MapEntityToContactType(entity), nil
}

// CompanyPositions is the resolver for the companyPositions field.
func (r *contactResolver) CompanyPositions(ctx context.Context, obj *model.Contact) ([]*model.CompanyPosition, error) {
	companyPositionEntities, err := r.ServiceContainer.CompanyService.GetCompanyPositionsForContact(ctx, obj.ID)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to get company positions %s", obj.ID)
		return nil, err
	}
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
	customFieldEntities, err := r.ServiceContainer.CustomFieldService.FindAllForContact(ctx, obj)
	for _, v := range mapper.MapEntitiesToCustomFields(customFieldEntities) {
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
		graphql.AddErrorf(ctx, "Failed to get contact definition for contact %s", obj.ID)
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	return mapper.MapEntityToEntityDefinition(entity), err
}

// Contact is the resolver for the contact field.
func (r *queryResolver) Contact(ctx context.Context, id string) (*model.Contact, error) {
	contactEntity, err := r.ServiceContainer.ContactService.FindContactById(ctx, id)
	if err != nil || contactEntity == nil {
		graphql.AddErrorf(ctx, "Contact with id %s not found", id)
		return nil, err
	}
	return mapper.MapEntityToContact(contactEntity), nil
}

// Contacts is the resolver for the contacts field.
func (r *queryResolver) Contacts(ctx context.Context, paginationFilter *model.PaginationFilter) (*model.ContactsPage, error) {
	if paginationFilter == nil {
		paginationFilter = &model.PaginationFilter{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.ContactService.FindAll(ctx, paginationFilter.Page, paginationFilter.Limit)
	return &model.ContactsPage{
		Content:       mapper.MapEntitiesToContacts(paginatedResult.Rows.(*entity.ContactEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// ContactByEmail is the resolver for the contactByEmail field.
func (r *queryResolver) ContactByEmail(ctx context.Context, email string) (*model.Contact, error) {
	contactEntity, err := r.ServiceContainer.ContactService.FindContactByEmail(ctx, email)
	if err != nil || contactEntity == nil {
		graphql.AddErrorf(ctx, "Contact with email %s not identified", email)
		return nil, err
	}
	return mapper.MapEntityToContact(contactEntity), nil
}

// ContactByPhone is the resolver for the contactByPhone field.
func (r *queryResolver) ContactByPhone(ctx context.Context, e164 string) (*model.Contact, error) {
	contactEntity, err := r.ServiceContainer.ContactService.FindContactByPhoneNumber(ctx, e164)
	if err != nil || contactEntity == nil {
		graphql.AddErrorf(ctx, "Contact with phone number %s not identified", e164)
		return nil, err
	}
	return mapper.MapEntityToContact(contactEntity), nil
}

// Contact returns generated.ContactResolver implementation.
func (r *Resolver) Contact() generated.ContactResolver { return &contactResolver{r} }

type contactResolver struct{ *Resolver }
