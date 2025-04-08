package user

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type userService struct {
	neo4j     *neo4j_repository.Repositories
	postgres  *postgres_repository.Repositories
	events    *events.EventsService
	mailstack interfaces.MailstackService
}

func NewUserService(neo4j *neo4j_repository.Repositories, postgres *postgres_repository.Repositories, events *events.EventsService) interfaces.UserService {
	return &userService{
		neo4j:    neo4j,
		postgres: postgres,
		events:   events,
	}
}

func (s *userService) SetMailstack(mailstack interfaces.MailstackService) {
	s.mailstack = mailstack
}

func (s *userService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, userFields data_fields.UserFields) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.Save")
	defer spans.Finish()

	spans.LogObjectAsJson("userFields", userFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	userId := ""

	if id == nil || *id == "" {
		createFlow = true
		spans.LogKV("flow", "create")
	} else {
		spans.LogKV("flow", "update")
	}

	if createFlow {
		// generate id
		userId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelUser)
		if err != nil {
			spans.TraceError(err)
			return "", err
		}

		// prepare missing fields
		if userFields.CreatedAt == nil {
			userFields.CreatedAt = utils.NowPtr()
		} else {
			userFields.CreatedAt = utils.TimePtr(utils.NowIfZero(*userFields.CreatedAt))
		}
		if utils.IfNotNilString(userFields.Source) == "" {
			userFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(userFields.AppSource) == "" {
			userFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
	} else {
		userId = *id
		// validate user exists
		exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, userId, model.NodeLabelUser)
		if err != nil || !exists {
			err = errors.New("user not found")
			spans.TraceError(err)
			return "", err
		}
	}
	spans.TagEntity(userId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		var innerErr error
		if createFlow {
			innerErr = s.neo4j.UserWriteRepository.CreateUserInTx(ctx, txWithPostCommit.Tx, tenant, userId, userFields)
		} else {
			innerErr = s.neo4j.UserWriteRepository.UpdateUserInTx(ctx, txWithPostCommit.Tx, tenant, userId, userFields)
		}
		if innerErr != nil {
			return nil, innerErr
		}

		if userFields.ExternalSystemAvailable() {
			innerErr = s.neo4j.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, userId, model.NodeLabelUser, *userFields.ExternalSystem)
			if err != nil {
				return nil, innerErr
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if createFlow {
				err = s.events.Publisher.PublishFanoutEvent(ctx, userId, model.USER, dto.CreateUser{UserFields: userFields})
				if err != nil {
					spans.TraceError(errors.Wrap(err, "unable to publish message CreateUser"))
				}
			} else {
				err = s.events.Publisher.PublishFanoutEvent(ctx, userId, model.USER, dto.UpdateUser{UserFields: userFields})
				if err != nil {
					spans.TraceError(errors.Wrap(err, "unable to publish message UpdateUser"))
				}
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	if createFlow {
		spans.LogKV("response.userCreated", true)
	} else {
		spans.LogKV("response.userUpdated", true)
	}
	return userId, nil
}

func (s *userService) GetById(parentCtx context.Context, userId string) (*neo4jentity.UserEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetById")
	defer spans.Finish()

	node, err := s.neo4j.UserReadRepository.GetUserById(ctx, common.GetContext(ctx).Tenant, userId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(node), nil
}

func (s *userService) GetAllUsersForTenant(ctx context.Context, tenant string) ([]*neo4jentity.UserEntity, error) {
	nodes, err := s.neo4j.UserReadRepository.GetAllForTenant(ctx, tenant)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsersForTenant: %w", err)
	}

	users := make([]*neo4jentity.UserEntity, len(nodes))

	for i, node := range nodes {
		users[i] = mapper.MapDbNodeToUserEntity(node)
	}

	return users, nil
}

func (s *userService) FindUserByEmail(parentCtx context.Context, email string) (*neo4jentity.UserEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.FindFirstUserWithRolesByEmail")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	userDbNode, err := s.neo4j.UserReadRepository.GetFirstUserByEmail(ctx, tenant, email)
	if err != nil {
		return nil, err
	}

	if userDbNode == nil {
		return nil, nil
	}

	return mapper.MapDbNodeToUserEntity(userDbNode), nil
}

func (s *userService) GetContactOwner(parentCtx context.Context, contactId string) (*neo4jentity.UserEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetContactOwner")
	defer spans.Finish()

	ownerDbNode, err := s.neo4j.UserReadRepository.GetOwnerForContact(ctx, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(ownerDbNode), nil
}

func (s *userService) GetNoteCreator(parentCtx context.Context, noteId string) (*neo4jentity.UserEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetNoteCreator")
	defer spans.Finish()

	userDbNode, err := s.neo4j.UserReadRepository.GetCreatorForNote(ctx, common.GetContext(ctx).Tenant, noteId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(userDbNode), nil
}

func (s *userService) GetUsersConnectedForContacts(ctx context.Context, contactIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetUsersConnectedForContacts")
	defer spans.Finish()

	users, err := s.neo4j.UserReadRepository.GetUsersConnectedForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUsersByEmailIds(parentCtx context.Context, emailIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUsersByEmailIds")
	defer spans.Finish()

	users, err := s.neo4j.UserReadRepository.GetUsersByEmailIds(ctx, common.GetTenantFromContext(ctx), emailIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	spans.LogKV("result.count", len(userEntities))
	return &userEntities, nil
}

func (s *userService) GetUsersByEmailAddresses(parentCtx context.Context, emailAddresses []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUsersByEmailAddresses")
	defer spans.Finish()

	spans.LogObjectAsJson("emailAddresses", emailAddresses)

	users, err := s.neo4j.UserReadRepository.GetUsersByEmailAddresses(ctx, common.GetTenantFromContext(ctx), emailAddresses)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	spans.LogKV("result.count", len(userEntities))
	return &userEntities, nil
}

func (s *userService) GetUsersForPhoneNumbers(parentCtx context.Context, phoneNumberIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUsersForPhoneNumbers")
	defer spans.Finish()

	users, err := s.neo4j.UserReadRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), phoneNumberIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserOwnersForOrganizations(parentCtx context.Context, organizationIDs []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUserOwnersForOrganizations")
	defer spans.Finish()

	spans.LogFields(log.Object("organizationIDs", organizationIDs))

	users, err := s.neo4j.UserReadRepository.GetAllOwnersForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserOwnersForOpportunities(parentCtx context.Context, opportunityIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUserOwnersForOpportunities")
	defer spans.Finish()

	spans.LogFields(log.Object("opportunityIds", opportunityIds))

	users, err := s.neo4j.UserReadRepository.GetAllOwnersForOpportunities(ctx, common.GetTenantFromContext(ctx), opportunityIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserCreatorsForOpportunities(parentCtx context.Context, opportunityIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUserCreatorsForOpportunities")
	defer spans.Finish()

	spans.LogFields(log.Object("opportunityIds", opportunityIds))

	users, err := s.neo4j.UserReadRepository.GetAllCreatorsForOpportunities(ctx, common.GetTenantFromContext(ctx), opportunityIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserCreatorsForServiceLineItems(parentCtx context.Context, serviceLineItemIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUserCreatorsForOpportunities")
	defer spans.Finish()

	spans.LogFields(log.Object("serviceLineItemIds", serviceLineItemIds))

	users, err := s.neo4j.UserReadRepository.GetAllCreatorsForServiceLineItems(ctx, common.GetTenantFromContext(ctx), serviceLineItemIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserCreatorsForTasks(ctx context.Context, taskIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetUserCreatorsForTasks")
	defer spans.Finish()

	spans.LogObjectAsJson("taskIds", taskIds)

	users, err := s.neo4j.UserReadRepository.GetAllCreatorsForTasks(ctx, common.GetTenantFromContext(ctx), taskIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserAssigneesForTasks(ctx context.Context, taskIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetUserCreatorsForTasks")
	defer spans.Finish()

	spans.LogObjectAsJson("taskIds", taskIds)

	users, err := s.neo4j.UserReadRepository.GetAllAssigneesForTasks(ctx, common.GetTenantFromContext(ctx), taskIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUsersWithMailboxes(ctx context.Context) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetUsersWithMailboxes")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	_, _, mailboxes, err := s.mailstack.GetMailboxes(ctx, tenant, "", "")
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Create a map to store unique usernames
	userEmailMap := make(map[string]struct{})
	for _, mailbox := range mailboxes {
		// Add the username to the map (map ensures uniqueness)
		userEmailMap[mailbox.Username] = struct{}{}
	}

	entities := make(neo4jentity.UserEntities, 0, len(userEmailMap))
	for userEmail := range userEmailMap {
		userNode, err := s.neo4j.UserReadRepository.GetFirstUserByEmail(ctx, tenant, userEmail)
		if err != nil {
			return nil, err
		}

		if userNode != nil {
			userEntity := mapper.MapDbNodeToUserEntity(userNode)
			entities = append(entities, *userEntity)
		}
	}

	return &entities, nil
}

func (s *userService) GetUserCreatorsForContracts(parentCtx context.Context, contractIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUserCreatorsForContracts")
	defer spans.Finish()

	spans.LogFields(log.Object("contractIds", contractIds))

	users, err := s.neo4j.UserReadRepository.GetAllCreatorsForContracts(ctx, common.GetTenantFromContext(ctx), contractIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserAuthorsForLogEntries(parentCtx context.Context, logEntryIDs []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUserAuthorsForLogEntries")
	defer spans.Finish()

	spans.LogFields(log.Object("logEntryIDs", logEntryIDs))

	users, err := s.neo4j.UserReadRepository.GetAllAuthorsForLogEntries(ctx, common.GetTenantFromContext(ctx), logEntryIDs)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserAuthorsForComments(ctx context.Context, commentIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetUserAuthorsForComments")
	defer spans.Finish()

	spans.LogFields(log.Object("commentIds", commentIds))

	users, err := s.neo4j.UserReadRepository.GetAllAuthorsForComments(ctx, common.GetTenantFromContext(ctx), commentIds)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserForFlowSenders(ctx context.Context, flowSenderIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetUserForFlowSenders")
	defer spans.Finish()

	spans.LogFields(log.Object("flowSenderIds", flowSenderIds))

	users, err := s.neo4j.UserReadRepository.GetAllSendersForFlowSenders(ctx, common.GetTenantFromContext(ctx), flowSenderIds)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUsers(parentCtx context.Context, userIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetUsers")
	defer spans.Finish()

	spans.LogFields(log.Object("userIds", userIds))

	userDbNodes, err := s.neo4j.UserReadRepository.GetUsers(ctx, common.GetTenantFromContext(ctx), userIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(userDbNodes))
	for _, dbNode := range userDbNodes {
		userEntity := mapper.MapDbNodeToUserEntity(dbNode)
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetAllOwnersForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetAllOwnersForOrganizations")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	dbNodes, err := s.neo4j.UserReadRepository.GetAllOwnersForOrganizations(ctx, tenant, organizationIds)
	if err != nil {
		return nil, err
	}

	userEntities := make(neo4jentity.UserEntities, 0, len(dbNodes))
	for _, v := range dbNodes {
		userEntity := mapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetDistinctOrganizationOwners(parentCtx context.Context) (*neo4jentity.UserEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetDistinctOrganizationOwners")
	defer spans.Finish()

	dbNodes, err := s.neo4j.UserReadRepository.GetDistinctOrganizationOwners(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return nil, err
	}

	userEntities := make(neo4jentity.UserEntities, 0, len(dbNodes))
	for _, dbNode := range dbNodes {
		entity := mapper.MapDbNodeToUserEntity(dbNode)
		userEntities = append(userEntities, *entity)
	}
	return &userEntities, nil
}

func (s *userService) GetReminderOwner(ctx context.Context, reminderId string) (*neo4jentity.UserEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetReminderOwner")
	defer spans.Finish()

	ownerDbNode, err := s.neo4j.UserReadRepository.GetOwnerForReminder(ctx, common.GetContext(ctx).Tenant, reminderId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(ownerDbNode), nil
}

func (s *userService) GetContractOwner(parentCtx context.Context, contractId string) (*neo4jentity.UserEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.GetContractOwner")
	defer spans.Finish()

	ownerDbNode, err := s.neo4j.UserReadRepository.GetOwnerForContract(ctx, common.GetContext(ctx).Tenant, contractId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(ownerDbNode), nil
}
