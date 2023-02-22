package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// PhoneNumberMergeToContact is the resolver for the phoneNumberMergeToContact field.
func (r *mutationResolver) PhoneNumberMergeToContact(ctx context.Context, contactID string, input model.PhoneNumberInput) (*model.PhoneNumber, error) {
	panic(fmt.Errorf("not implemented: PhoneNumberMergeToContact - phoneNumberMergeToContact"))
}

// PhoneNumberUpdateInContact is the resolver for the phoneNumberUpdateInContact field.
func (r *mutationResolver) PhoneNumberUpdateInContact(ctx context.Context, contactID string, input model.PhoneNumberUpdateInput) (*model.PhoneNumber, error) {
	panic(fmt.Errorf("not implemented: PhoneNumberUpdateInContact - phoneNumberUpdateInContact"))
}

// PhoneNumberRemoveFromContactByE164 is the resolver for the phoneNumberRemoveFromContactByE164 field.
func (r *mutationResolver) PhoneNumberRemoveFromContactByE164(ctx context.Context, contactID string, e164 string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: PhoneNumberRemoveFromContactByE164 - phoneNumberRemoveFromContactByE164"))
}

// PhoneNumberRemoveFromContactByID is the resolver for the phoneNumberRemoveFromContactById field.
func (r *mutationResolver) PhoneNumberRemoveFromContactByID(ctx context.Context, contactID string, id string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: PhoneNumberRemoveFromContactByID - phoneNumberRemoveFromContactById"))
}
