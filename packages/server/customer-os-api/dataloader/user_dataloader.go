package dataloader

import (
	"context"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
)

func (i *Loaders) GetUsersForEmail(ctx context.Context, emailID string) (*neo4jentity.UserEntities, error) {
	thunk := i.UsersForEmail.Load(ctx, dataloader.StringKey(emailID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.UserEntities)
	return &resultObj, nil
}

func (i *Loaders) GetUsersConnectedForContact(ctx context.Context, contactId string) (*neo4jentity.UserEntities, error) {
	thunk := i.UsersConnectedForContact.Load(ctx, dataloader.StringKey(contactId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.UserEntities)
	return &resultObj, nil
}

func (i *Loaders) GetUsersAssigneesForTask(ctx context.Context, taskId string) (*neo4jentity.UserEntities, error) {
	thunk := i.UsersForTask.Load(ctx, dataloader.StringKey(taskId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.UserEntities)
	return &resultObj, nil
}

func (i *Loaders) GetUsersForPhoneNumber(ctx context.Context, phoneNumberID string) (*neo4jentity.UserEntities, error) {
	thunk := i.UsersForPhoneNumber.Load(ctx, dataloader.StringKey(phoneNumberID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	resultObj := result.(neo4jentity.UserEntities)
	return &resultObj, nil
}

func (i *Loaders) GetUserOwnerForOrganization(ctx context.Context, organizationID string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserOwnerForOrganization.Load(ctx, dataloader.StringKey(organizationID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUserOwnerForOpportunity(ctx context.Context, opportunityId string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserOwnerForOpportunity.Load(ctx, dataloader.StringKey(opportunityId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUserCreatorForOpportunity(ctx context.Context, opportunityId string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserCreatorForOpportunity.Load(ctx, dataloader.StringKey(opportunityId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUserCreatorForServiceLineItem(ctx context.Context, serviceLineItemId string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserCreatorForServiceLineItem.Load(ctx, dataloader.StringKey(serviceLineItemId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUserCreatorForContract(ctx context.Context, contractId string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserCreatorForContract.Load(ctx, dataloader.StringKey(contractId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUser(ctx context.Context, userId string) (*neo4jentity.UserEntity, error) {
	thunk := i.User.Load(ctx, dataloader.StringKey(userId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUserAuthorForLogEntry(ctx context.Context, logEntryId string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserAuthorForLogEntry.Load(ctx, dataloader.StringKey(logEntryId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUserForFlowSenders(ctx context.Context, flowSenderId string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserForFlowSender.Load(ctx, dataloader.StringKey(flowSenderId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUserAuthorForComment(ctx context.Context, logEntryId string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserAuthorForComment.Load(ctx, dataloader.StringKey(logEntryId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (i *Loaders) GetUserCreatorForTask(ctx context.Context, taskId string) (*neo4jentity.UserEntity, error) {
	thunk := i.UserCreatorForTask.Load(ctx, dataloader.StringKey(taskId))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*neo4jentity.UserEntity), nil
}

func (b *userBatcher) getUsersConnectedForContact(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUsersConnectedForContact")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	userEntitiesPtr, err := b.userCommonService.GetUsersConnectedForContacts(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntitiesByEmailId := make(map[string]neo4jentity.UserEntities)
	for _, val := range *userEntitiesPtr {
		if list, ok := userEntitiesByEmailId[val.DataloaderKey]; ok {
			userEntitiesByEmailId[val.DataloaderKey] = append(list, val)
		} else {
			userEntitiesByEmailId[val.DataloaderKey] = neo4jentity.UserEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for emailId, record := range userEntitiesByEmailId {
		if ix, ok := keyOrder[emailId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, emailId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.UserEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.UserEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUsersForEmails(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUsersForEmails")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	userEntitiesPtr, err := b.userCommonService.GetUsersByEmailIds(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntitiesByEmailId := make(map[string]neo4jentity.UserEntities)
	for _, val := range *userEntitiesPtr {
		if list, ok := userEntitiesByEmailId[val.DataloaderKey]; ok {
			userEntitiesByEmailId[val.DataloaderKey] = append(list, val)
		} else {
			userEntitiesByEmailId[val.DataloaderKey] = neo4jentity.UserEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for emailId, record := range userEntitiesByEmailId {
		if ix, ok := keyOrder[emailId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, emailId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.UserEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.UserEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUsersForTasks(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUsersForTasks")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	userEntitiesPtr, err := b.userCommonService.GetUserAssigneesForTasks(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntitiesByTaskId := make(map[string]neo4jentity.UserEntities)
	for _, val := range *userEntitiesPtr {
		if list, ok := userEntitiesByTaskId[val.DataloaderKey]; ok {
			userEntitiesByTaskId[val.DataloaderKey] = append(list, val)
		} else {
			userEntitiesByTaskId[val.DataloaderKey] = neo4jentity.UserEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for taskId, record := range userEntitiesByTaskId {
		if ix, ok := keyOrder[taskId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, taskId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.UserEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.UserEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUsersForPhoneNumbers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUsersForPhoneNumbers")
	defer spans.Finish()
	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	userEntitiesPtr, err := b.userCommonService.GetUsersForPhoneNumbers(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntitiesByPhoneNumberId := make(map[string]neo4jentity.UserEntities)
	for _, val := range *userEntitiesPtr {
		if list, ok := userEntitiesByPhoneNumberId[val.DataloaderKey]; ok {
			userEntitiesByPhoneNumberId[val.DataloaderKey] = append(list, val)
		} else {
			userEntitiesByPhoneNumberId[val.DataloaderKey] = neo4jentity.UserEntities{val}
		}
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for phoneNumberId, record := range userEntitiesByPhoneNumberId {
		if ix, ok := keyOrder[phoneNumberId]; ok {
			results[ix] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, phoneNumberId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: neo4jentity.UserEntities{}, Error: nil}
	}

	if err = assertEntitiesType(results, reflect.TypeOf(neo4jentity.UserEntities{})); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserOwnersForOrganizations(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserOwnersForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	userEntities, err := b.userCommonService.GetUserOwnersForOrganizations(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByOrganizationId := make(map[string]neo4jentity.UserEntity)
	for _, val := range *userEntities {
		userEntityByOrganizationId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for organizationID, _ := range userEntityByOrganizationId {
		if ix, ok := keyOrder[organizationID]; ok {
			val := userEntityByOrganizationId[organizationID]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, organizationID)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserOwnersForOpportunities(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserOwnersForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	userEntities, err := b.userCommonService.GetUserOwnersForOpportunities(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByOpportunityId := make(map[string]neo4jentity.UserEntity)
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

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserCreatorsForOpportunities(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserOwnersForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	userEntities, err := b.userCommonService.GetUserCreatorsForOpportunities(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByOpportunityId := make(map[string]neo4jentity.UserEntity)
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

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserCreatorsForServiceLineItems(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserCreatorsForServiceLineItems")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	userEntities, err := b.userCommonService.GetUserCreatorsForServiceLineItems(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByServiceLineItemId := make(map[string]neo4jentity.UserEntity)
	for _, val := range *userEntities {
		userEntityByServiceLineItemId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for serviceLineItemID, _ := range userEntityByServiceLineItemId {
		if ix, ok := keyOrder[serviceLineItemID]; ok {
			val := userEntityByServiceLineItemId[serviceLineItemID]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, serviceLineItemID)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserCreatorsForContracts(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserCreatorsForContracts")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	userEntities, err := b.userCommonService.GetUserCreatorsForContracts(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByContractId := make(map[string]neo4jentity.UserEntity)
	for _, val := range *userEntities {
		userEntityByContractId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for contractId, _ := range userEntityByContractId {
		if ix, ok := keyOrder[contractId]; ok {
			val := userEntityByContractId[contractId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, contractId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserCreatorsForTasks(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserCreatorsForTasks")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	userEntities, err := b.userCommonService.GetUserCreatorsForTasks(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByTaskId := make(map[string]neo4jentity.UserEntity)
	for _, val := range *userEntities {
		userEntityByTaskId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for taskId, _ := range userEntityByTaskId {
		if ix, ok := keyOrder[taskId]; ok {
			val := userEntityByTaskId[taskId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, taskId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUsers")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	userEntities, err := b.userCommonService.GetUsers(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityById := make(map[string]neo4jentity.UserEntity)
	for _, val := range *userEntities {
		userEntityById[val.Id] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for id, _ := range userEntityById {
		if ix, ok := keyOrder[id]; ok {
			val := userEntityById[id]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, id)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserAuthorsForLogEntries(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserAuthorsForLogEntries")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	userEntities, err := b.userCommonService.GetUserAuthorsForLogEntries(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByLogEntryId := make(map[string]neo4jentity.UserEntity)
	for _, val := range *userEntities {
		userEntityByLogEntryId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for logEntryId, _ := range userEntityByLogEntryId {
		if ix, ok := keyOrder[logEntryId]; ok {
			val := userEntityByLogEntryId[logEntryId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, logEntryId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserAuthorsForComments(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserAuthorsForLogEntries")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	userEntities, err := b.userCommonService.GetUserAuthorsForComments(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByLogEntryId := make(map[string]neo4jentity.UserEntity)
	for _, val := range *userEntities {
		userEntityByLogEntryId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for logEntryId, _ := range userEntityByLogEntryId {
		if ix, ok := keyOrder[logEntryId]; ok {
			val := userEntityByLogEntryId[logEntryId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, logEntryId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}

func (b *userBatcher) getUserForFlowSenders(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserDataLoader.getUserForFlowSenders")
	defer spans.Finish()

	spans.LogObjectAsJson("keys", keys)
	spans.LogKV("keys_length", len(keys))

	ids, keyOrder := sortKeys(keys)

	ctx, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	userEntities, err := b.userCommonService.GetUserForFlowSenders(ctx, ids)
	if err != nil {
		spans.TraceError(err)
		// check if context deadline exceeded error occurred
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return []*dataloader.Result{{Data: nil, Error: errors.Wrap(err, "context deadline exceeded")}}
		}
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	userEntityByLogEntryId := make(map[string]neo4jentity.UserEntity)
	for _, val := range *userEntities {
		userEntityByLogEntryId[val.DataloaderKey] = val
	}

	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for logEntryId, _ := range userEntityByLogEntryId {
		if ix, ok := keyOrder[logEntryId]; ok {
			val := userEntityByLogEntryId[logEntryId]
			results[ix] = &dataloader.Result{Data: &val, Error: nil}
			delete(keyOrder, logEntryId)
		}
	}
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: nil}
	}

	if err = assertEntitiesPtrType(results, reflect.TypeOf(neo4jentity.UserEntity{}), true); err != nil {
		spans.TraceError(err)
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	spans.LogKV("result.length", len(results))

	return results
}
