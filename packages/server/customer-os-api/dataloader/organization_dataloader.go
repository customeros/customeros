package dataloader

import (
	"context"
	"errors"
	"github.com/graph-gophers/dataloader"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

func (i *Loaders) GetOrganizationsForEmail(ctx context.Context, emailId string) (*entity.OrganizationEntities, error) {
	thunk := i.OrganizationsForEmail.Load(ctx, dataloader.StringKey(emailId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetOrganizationsForPhoneNumber(ctx context.Context, phoneNumberId string) (*entity.OrganizationEntities, error) {
	thunk := i.OrganizationsForPhoneNumber.Load(ctx, dataloader.StringKey(phoneNumberId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetSubsidiariesForOrganization(ctx context.Context, organizationId string) (*entity.OrganizationEntities, error) {
	thunk := i.SubsidiariesForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetSubsidiariesOfForOrganization(ctx context.Context, organizationId string) (*entity.OrganizationEntities, error) {
	thunk := i.SubsidiariesOfForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrganizationEntities)
	return &resultObj, nil
}

func (i *Loaders) GetOrganizationForJobRole(ctx context.Context, jobRoleId string) (*entity.OrganizationEntity, error) {
	thunk := i.OrganizationForJobRole.Load(ctx, dataloader.StringKey(jobRoleId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*entity.OrganizationEntity), nil
}

func (i *Loaders) GetSuggestedMergeToForOrganization(ctx context.Context, organizationId string) (*entity.OrganizationEntities, error) {
	thunk := i.SuggestedMergeToForOrganization.Load(ctx, dataloader.StringKey(organizationId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(entity.OrganizationEntities)
	return &resultObj, nil
}

func (b *organizationBatcher) getOrganizationsForEmails(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForEmails", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
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

	organizationEntitiesByEmailId := make(map[string]entity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByEmailId[val.DataloaderKey]; ok {
			organizationEntitiesByEmailId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByEmailId[val.DataloaderKey] = entity.OrganizationEntities{val}
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
		results[ix] = &dataloader.Result{Data: entity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getOrganizationsForPhoneNumbers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForPhoneNumbers", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
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

	organizationEntitiesByPhoneNumberId := make(map[string]entity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByPhoneNumberId[val.DataloaderKey]; ok {
			organizationEntitiesByPhoneNumberId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByPhoneNumberId[val.DataloaderKey] = entity.OrganizationEntities{val}
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
		results[ix] = &dataloader.Result{Data: entity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getSubsidiariesForOrganization(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getSubsidiariesForOrganization", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
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

	organizationEntitiesByOrgId := make(map[string]entity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByOrgId[val.DataloaderKey]; ok {
			organizationEntitiesByOrgId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByOrgId[val.DataloaderKey] = entity.OrganizationEntities{val}
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
		results[ix] = &dataloader.Result{Data: entity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getSubsidiariesOfForOrganization(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getSubsidiariesOfForOrganization", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
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

	organizationEntitiesByOrgId := make(map[string]entity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByOrgId[val.DataloaderKey]; ok {
			organizationEntitiesByOrgId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByOrgId[val.DataloaderKey] = entity.OrganizationEntities{val}
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
		results[ix] = &dataloader.Result{Data: entity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getOrganizationsForJobRoles(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getOrganizationsForJobRoles", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
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

	organizationEntityByJobRoleId := make(map[string]entity.OrganizationEntity)
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

	if err = assertEntitiesPtrType(results, reflect.TypeOf(entity.OrganizationEntity{}), true); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}

func (b *organizationBatcher) getSuggestedMergeToForOrganization(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationDataLoader.getSuggestedMergeToForOrganization", opentracing.ChildOf(tracing.ExtractSpanCtx(ctx)))
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

	organizationEntitiesByOrgId := make(map[string]entity.OrganizationEntities)
	for _, val := range *organizationEntitiesPtr {
		if list, ok := organizationEntitiesByOrgId[val.DataloaderKey]; ok {
			organizationEntitiesByOrgId[val.DataloaderKey] = append(list, val)
		} else {
			organizationEntitiesByOrgId[val.DataloaderKey] = entity.OrganizationEntities{val}
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
		results[ix] = &dataloader.Result{Data: entity.OrganizationEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(entity.OrganizationEntities{})); err != nil {
		tracing.TraceErr(span, err)
		return []*dataloader.Result{{nil, err}}
	}

	span.LogFields(log.Int("results_length", len(results)))

	return results
}
