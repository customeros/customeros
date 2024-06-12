package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"reflect"
)

func (i *Loaders) GetOrganizationsForEmail(ctx context.Context, emailId string) (*neo4jentity.OrganizationEntities, error) {
	thunk := i.OrganizationsForEmail.Load(ctx, dataloader.StringKey(emailId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetOrganizationsForPhoneNumber(ctx context.Context, phoneNumberId string) (*neo4jentity.OrganizationEntities, error) {
	thunk := i.OrganizationsForPhoneNumber.Load(ctx, dataloader.StringKey(phoneNumberId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetSubsidiariesForOrganization(ctx context.Context, organizationId string) (*neo4jentity.OrganizationEntities, error) {
	thunk := i.SubsidiariesForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetSubsidiariesOfForOrganization(ctx context.Context, organizationId string) (*neo4jentity.OrganizationEntities, error) {
	thunk := i.SubsidiariesOfForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetOrganizationForJobRole(ctx context.Context, jobRoleId string) (*neo4jentity.OrganizationEntity, error) {
	thunk := i.OrganizationForJobRole.Load(ctx, dataloader.StringKey(jobRoleId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.OrganizationEntity), nil
}

func (i *Loaders) GetSuggestedMergeToForOrganization(ctx context.Context, organizationId string) (*neo4jentity.OrganizationEntities, error) {
	thunk := i.SuggestedMergeToForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetOrganization(ctx context.Context, organizationId string) (*neo4jentity.OrganizationEntity, error) {
	thunk := i.Organization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.OrganizationEntity), nil
}

func (i *Loaders) GetOrganizationForInvoice(ctx context.Context, invoiceId string) (*neo4jentity.OrganizationEntity, error) {
	thunk := i.OrganizationForInvoice.Load(ctx, dataloader.StringKey(invoiceId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.OrganizationEntity), nil
}

func (i *Loaders) GetOrganizationForSlackChannel(ctx context.Context, slackChannelId string) (*neo4jentity.OrganizationEntity, error) {
	thunk := i.OrganizationForInvoice.Load(ctx, dataloader.StringKey(slackChannelId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.OrganizationEntity), nil
}

func (i *Loaders) GetOrganizationForOpportunityOptional(ctx context.Context, opportunityId string) (*neo4jentity.OrganizationEntity, error) {
	thunk := i.OrganizationForOpportunity.Load(ctx, dataloader.StringKey(opportunityId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.OrganizationEntity), nil
}

func (b *organizationBatcher) getOrganizationsForEmails(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForEmails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	organizationEntitiesPtr, err := b.organizationService.GetOrganizationsForEmails(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get organizations for emails")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntitiesByEmailId := make(map[string]neo4jentity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByEmailId[val.DataloaderKey]; ok {
			organizationEntitiesByEmailId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByEmailId[val.DataloaderKey] = neo4jentity.OrganizationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for emailId, record := range organizationEntitiesByEmailId {
		if ix, ok := keyOrder[emailId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, emailId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getOrganizationsForPhoneNumbers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	organizationEntitiesPtr, err := b.organizationService.GetOrganizationsForPhoneNumbers(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get organizations for phone numbers")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntitiesByPhoneNumberId := make(map[string]neo4jentity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByPhoneNumberId[val.DataloaderKey]; ok {
			organizationEntitiesByPhoneNumberId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByPhoneNumberId[val.DataloaderKey] = neo4jentity.OrganizationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for phoneNumberId, record := range organizationEntitiesByPhoneNumberId {
		if ix, ok := keyOrder[phoneNumberId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, phoneNumberId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getSubsidiariesForOrganization(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getSubsidiariesForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	organizationEntitiesPtr, err := b.organizationService.GetSubsidiariesForOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get subsidiaries for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntitiesByOrgId := make(map[string]neo4jentity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByOrgId[val.DataloaderKey]; ok {
			organizationEntitiesByOrgId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByOrgId[val.DataloaderKey] = neo4jentity.OrganizationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for orgId, record := range organizationEntitiesByOrgId {
		if ix, ok := keyOrder[orgId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, orgId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getSubsidiariesOfForOrganization(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getSubsidiariesOfForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	organizationEntitiesPtr, err := b.organizationService.GetSubsidiariesOfForOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get subsidiaries of for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntitiesByOrgId := make(map[string]neo4jentity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByOrgId[val.DataloaderKey]; ok {
			organizationEntitiesByOrgId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByOrgId[val.DataloaderKey] = neo4jentity.OrganizationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for orgId, record := range organizationEntitiesByOrgId {
		if ix, ok := keyOrder[orgId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, orgId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getOrganizationsForJobRoles(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForJobRoles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	organizationEntities, err := b.organizationService.GetOrganizationsForJobRoles(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get organizations for job roles")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntityByJobRoleId := make(map[string]neo4jentity.OrganizationEntity)
	for _, val := range *organizationEntities {
		organizationEntityByJobRoleId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for jobRoleId, _ := range organizationEntityByJobRoleId {
		if ix, ok := keyOrder[jobRoleId]; ok {
			val := organizationEntityByJobRoleId[jobRoleId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, jobRoleId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.OrganizationEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getSuggestedMergeToForOrganization(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getSuggestedMergeToForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	organizationEntitiesPtr, err := b.organizationService.GetSuggestedMergeToForOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if ctx.Err() == context.DeadlineExceeded {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get suggested merges for organizations")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntitiesByOrgId := make(map[string]neo4jentity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByOrgId[val.DataloaderKey]; ok {
			organizationEntitiesByOrgId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByOrgId[val.DataloaderKey] = neo4jentity.OrganizationEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for orgId, record := range organizationEntitiesByOrgId {
		if ix, ok := keyOrder[orgId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, orgId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getOrganizationsForInvoices(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForInvoices")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	organizationEntities, err := b.organizationService.GetOrganizationsForInvoices(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get organizations for invoices")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntityByInvoiceId := make(map[string]neo4jentity.OrganizationEntity)
	for _, val := range *organizationEntities {
		organizationEntityByInvoiceId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for jobRoleId, _ := range organizationEntityByInvoiceId {
		if ix, ok := keyOrder[jobRoleId]; ok {
			val := organizationEntityByInvoiceId[jobRoleId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, jobRoleId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.OrganizationEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("result.length", len(results)))

	return results
}

func (b *organizationBatcher) getOrganizationsForSlackChannels(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForSlackChannels")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	organizationEntities, err := b.organizationService.GetOrganizationsForSlackChannels(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.New("deadline exceeded to get organizations for invoices")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntityBySlackChannelId := make(map[string]neo4jentity.OrganizationEntity)
	for _, val := range *organizationEntities {
		organizationEntityBySlackChannelId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for jobRoleId, _ := range organizationEntityBySlackChannelId {
		if ix, ok := keyOrder[jobRoleId]; ok {
			val := organizationEntityBySlackChannelId[jobRoleId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, jobRoleId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.OrganizationEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("result.length", len(results)))

	return results
}

func (b *organizationBatcher) getOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	organizationEntities, err := b.organizationService.GetOrganizations(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	organizationEntityById := make(map[string]neo4jentity.OrganizationEntity)
	for _, val := range *organizationEntities {
		organizationEntityById[val.ID] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for id, _ := range organizationEntityById {
		if ix, ok := keyOrder[id]; ok {
			val := organizationEntityById[id]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, id)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.OrganizationEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Object("output - results_length", len(results)))

	return results
}

func (b *organizationBatcher) getOrganizationsForOpportunities(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("keys", keys), log.Int("keys_length", len(keys)))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	organizationEntities, err := b.organizationService.GetOrganizationsForOpportunities(ctx, ids)
	if err != nil {
		tracing.TraceErr(span, err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByOpportunityId := make(map[string]entity.UserEntity)
	for _, val := range *userEntities {
		userEntityByOpportunityId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for opportunityID, _ := range userEntityByOpportunityId {
		if ix, ok := keyOrder[opportunityID]; ok {
			val := userEntityByOpportunityId[opportunityID]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, opportunityID)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(entity.UserEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Object("output - results_length", len(results)))

	return results
}
