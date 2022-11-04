package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/customer-os-api/service"
)

// CreateTenantUser is the resolver for the createTenantUser field.
func (r *mutationResolver) CreateTenantUser(ctx context.Context, input model.TenantUserInput) (*model.TenantUser, error) {
	createdTenantEntity, err := r.ServiceContainer.TenantUserService.Create(ctx, mapper.MapTenantUserInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create user %s %s", input.FirstName, input.LastName)
		return nil, err
	}
	return mapper.MapEntityToTenantUser(createdTenantEntity), nil
}

// CreateContact is the resolver for the createContact field.
func (r *mutationResolver) CreateContact(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
	contactNodeCreated, err := r.ServiceContainer.ContactService.Create(ctx, &service.ContactCreateData{
		ContactEntity:     mapper.MapContactInputToEntity(input),
		TextCustomFields:  mapper.MapTextCustomFieldInputsToEntities(input.TextCustomFields),
		PhoneNumberEntity: mapper.MapPhoneNumberInputToEntity(input.PhoneNumber),
		EmailEntity:       mapper.MapEmailInputToEntity(input.Email),
		CompanyPosition:   mapper.MapCompanyPositionInputToEntity(input.CompanyPosition),
	})
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact %s %s", input.FirstName, input.LastName)
		return nil, err
	}

	return mapper.MapEntityToContact(contactNodeCreated), nil
}

// UpdateContact is the resolver for the updateContact field.
func (r *mutationResolver) UpdateContact(ctx context.Context, input model.ContactUpdateInput) (*model.Contact, error) {
	updatedContact, err := r.ServiceContainer.ContactService.Update(ctx, mapper.MapContactUpdateInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to update contact %s", input.ID)
		return nil, err
	}
	return mapper.MapEntityToContact(updatedContact), nil
}

// HardDeleteContact is the resolver for the hardDeleteContact field.
func (r *mutationResolver) HardDeleteContact(ctx context.Context, contactID string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.ContactService.HardDelete(ctx, contactID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not hard delete contact %s", contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// SoftDeleteContact is the resolver for the softDeleteContact field.
func (r *mutationResolver) SoftDeleteContact(ctx context.Context, contactID string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.ContactService.SoftDelete(ctx, contactID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not soft delete contact %s", contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// MergeTextCustomFieldToContact is the resolver for the mergeTextCustomFieldToContact field.
func (r *mutationResolver) MergeTextCustomFieldToContact(ctx context.Context, contactID string, input model.TextCustomFieldInput) (*model.TextCustomField, error) {
	result, err := r.ServiceContainer.TextCustomFieldService.MergeTextCustomFieldToContact(ctx, contactID, mapper.MapTextCustomFieldInputToEntity(&input))
	if err != nil {
		graphql.AddErrorf(ctx, "Could not add custom field %s to contact %s", input.Name, contactID)
		return nil, err
	}
	return mapper.MapEntityToTextCustomField(result), nil
}

// UpdateTextCustomFieldInContact is the resolver for the updateTextCustomFieldInContact field.
func (r *mutationResolver) UpdateTextCustomFieldInContact(ctx context.Context, contactID string, input model.TextCustomFieldUpdateInput) (*model.TextCustomField, error) {
	result, err := r.ServiceContainer.TextCustomFieldService.UpdateTextCustomFieldInContact(ctx, contactID, mapper.MapTextCustomFieldUpdateInputToEntity(&input))
	if err != nil {
		graphql.AddErrorf(ctx, "Could not update text custom field %s in contact %s", input.ID, contactID)
		return nil, err
	}
	return mapper.MapEntityToTextCustomField(result), nil
}

// RemoveTextCustomFieldFromContact is the resolver for the removeTextCustomFieldFromContact field.
func (r *mutationResolver) RemoveTextCustomFieldFromContact(ctx context.Context, contactID string, fieldName string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.TextCustomFieldService.Delete(ctx, contactID, fieldName)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove text field %s from contact %s", fieldName, contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// RemoveTextCustomFieldFromContactByID is the resolver for the removeTextCustomFieldFromContactById field.
func (r *mutationResolver) RemoveTextCustomFieldFromContactByID(ctx context.Context, contactID string, id string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.TextCustomFieldService.DeleteById(ctx, contactID, id)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove text field %s from contact %s", id, contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// MergePhoneNumberToContact is the resolver for the mergePhoneNumberToContact field.
func (r *mutationResolver) MergePhoneNumberToContact(ctx context.Context, contactID string, input model.PhoneNumberInput) (*model.PhoneNumberInfo, error) {
	result, err := r.ServiceContainer.PhoneNumberService.MergePhoneNumberToContact(ctx, contactID, mapper.MapPhoneNumberInputToEntity(&input))
	if err != nil {
		graphql.AddErrorf(ctx, "Could not add phone number %s to contact %s", input.Number, contactID)
		return nil, err
	}
	return mapper.MapEntityToPhoneNumber(result), nil
}

// UpdatePhoneNumberInContact is the resolver for the updatePhoneNumberInContact field.
func (r *mutationResolver) UpdatePhoneNumberInContact(ctx context.Context, contactID string, input model.PhoneNumberUpdateInput) (*model.PhoneNumberInfo, error) {
	result, err := r.ServiceContainer.PhoneNumberService.UpdatePhoneNumberInContact(ctx, contactID, mapper.MapPhoneNumberUpdateInputToEntity(&input))
	if err != nil {
		graphql.AddErrorf(ctx, "Could not update email %s in contact %s", input.ID, contactID)
		return nil, err
	}
	return mapper.MapEntityToPhoneNumber(result), nil
}

// RemovePhoneNumberFromContact is the resolver for the removePhoneNumberFromContact field.
func (r *mutationResolver) RemovePhoneNumberFromContact(ctx context.Context, contactID string, phoneNumber string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.PhoneNumberService.Delete(ctx, contactID, phoneNumber)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove phone number %s from contact %s", phoneNumber, contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// RemovePhoneNumberFromContactByID is the resolver for the removePhoneNumberFromContactById field.
func (r *mutationResolver) RemovePhoneNumberFromContactByID(ctx context.Context, contactID string, id string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.PhoneNumberService.DeleteById(ctx, contactID, id)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove phone number %s from contact %s", id, contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// MergeEmailToContact is the resolver for the mergeEmailToContact field.
func (r *mutationResolver) MergeEmailToContact(ctx context.Context, contactID string, input model.EmailInput) (*model.EmailInfo, error) {
	result, err := r.ServiceContainer.EmailService.MergeEmailToContact(ctx, contactID, mapper.MapEmailInputToEntity(&input))
	if err != nil {
		graphql.AddErrorf(ctx, "Could not add email %s to contact %s", input.Email, contactID)
		return nil, err
	}
	return mapper.MapEntityToEmail(result), nil
}

// UpdateEmailInContact is the resolver for the updateEmailInContact field.
func (r *mutationResolver) UpdateEmailInContact(ctx context.Context, contactID string, input model.EmailUpdateInput) (*model.EmailInfo, error) {
	result, err := r.ServiceContainer.EmailService.UpdateEmailInContact(ctx, contactID, mapper.MapEmailUpdateInputToEntity(&input))
	if err != nil {
		graphql.AddErrorf(ctx, "Could not update email %s in contact %s", input.ID, contactID)
		return nil, err
	}
	return mapper.MapEntityToEmail(result), nil
}

// RemoveEmailFromContact is the resolver for the removeEmailFromContact field.
func (r *mutationResolver) RemoveEmailFromContact(ctx context.Context, contactID string, email string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.EmailService.Delete(ctx, contactID, email)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove email %s from contact %s", email, contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// RemoveEmailFromContactByID is the resolver for the removeEmailFromContactById field.
func (r *mutationResolver) RemoveEmailFromContactByID(ctx context.Context, contactID string, id string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.EmailService.DeleteById(ctx, contactID, id)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove email %s from contact %s", id, contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// CreateContactGroup is the resolver for the createContactGroup field.
func (r *mutationResolver) CreateContactGroup(ctx context.Context, input model.ContactGroupInput) (*model.ContactGroup, error) {
	contactGroupEntityCreated, err := r.ServiceContainer.ContactGroupService.Create(ctx, &entity.ContactGroupEntity{
		Name: input.Name,
	})
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact group %s", input.Name)
		return nil, err
	}
	return mapper.MapEntityToContactGroup(contactGroupEntityCreated), nil
}

// UpdateContactGroup is the resolver for the updateContactGroup field.
func (r *mutationResolver) UpdateContactGroup(ctx context.Context, input model.ContactGroupUpdateInput) (*model.ContactGroup, error) {
	updatedContactGroup, err := r.ServiceContainer.ContactGroupService.Update(ctx, &entity.ContactGroupEntity{
		Id:   input.ID,
		Name: input.Name,
	})
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to update contact group %s", input.ID)
		return nil, err
	}
	return mapper.MapEntityToContactGroup(updatedContactGroup), nil
}

// DeleteContactGroupAndUnlinkAllContacts is the resolver for the deleteContactGroupAndUnlinkAllContacts field.
func (r *mutationResolver) DeleteContactGroupAndUnlinkAllContacts(ctx context.Context, id string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.ContactGroupService.Delete(ctx, id)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not delete contact group %s", id)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// AddContactToGroup is the resolver for the addContactToGroup field.
func (r *mutationResolver) AddContactToGroup(ctx context.Context, contactID string, groupID string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.ContactGroupService.AddContactToGroup(ctx, contactID, groupID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not add contact to group")
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// RemoveContactFromGroup is the resolver for the removeContactFromGroup field.
func (r *mutationResolver) RemoveContactFromGroup(ctx context.Context, contactID string, groupID string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.ContactGroupService.RemoveContactFromGroup(ctx, contactID, groupID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove contact from group")
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
