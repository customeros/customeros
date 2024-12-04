package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UserService interface {
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, userFields data_fields.UserFields) (string, error)

	GetById(ctx context.Context, userId string) (*neo4jentity.UserEntity, error)
	GetAllUsersForTenant(ctx context.Context, tenant string) ([]*neo4jentity.UserEntity, error)
	FindUserByEmail(parentCtx context.Context, email string) (*neo4jentity.UserEntity, error)
	IsOwner(ctx context.Context, id string) (bool, error)
	GetContactOwner(ctx context.Context, contactId string) (*neo4jentity.UserEntity, error)
	GetNoteCreator(ctx context.Context, noteId string) (*neo4jentity.UserEntity, error)
	GetUsersConnectedForContacts(ctx context.Context, contactIds []string) (*neo4jentity.UserEntities, error)
	GetUsersForEmails(ctx context.Context, emailIds []string) (*neo4jentity.UserEntities, error)
	GetUsersForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*neo4jentity.UserEntities, error)
	GetUserOwnersForOrganizations(ctx context.Context, organizationIDs []string) (*neo4jentity.UserEntities, error)
	GetUserOwnersForOpportunities(ctx context.Context, opportunityIds []string) (*neo4jentity.UserEntities, error)
	GetUserCreatorsForOpportunities(ctx context.Context, opportunityIds []string) (*neo4jentity.UserEntities, error)
	GetUserCreatorsForServiceLineItems(ctx context.Context, serviceLineItemIds []string) (*neo4jentity.UserEntities, error)
	GetUsersWithMailboxes(ctx context.Context) (*neo4jentity.UserEntities, error)
	GetUserCreatorsForContracts(ctx context.Context, contractIds []string) (*neo4jentity.UserEntities, error)
	GetUserAuthorsForLogEntries(ctx context.Context, logEntryIDs []string) (*neo4jentity.UserEntities, error)
	GetUserAuthorsForComments(ctx context.Context, commentIds []string) (*neo4jentity.UserEntities, error)
	GetUserForFlowSenders(ctx context.Context, flowSenderIds []string) (*neo4jentity.UserEntities, error)
	GetUsers(ctx context.Context, userIds []string) (*neo4jentity.UserEntities, error)
	GetAllOwnersForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.UserEntities, error)
	GetDistinctOrganizationOwners(ctx context.Context) (*neo4jentity.UserEntities, error)
	GetReminderOwner(ctx context.Context, reminderId string) (*neo4jentity.UserEntity, error)
	GetContractOwner(ctx context.Context, contractId string) (*neo4jentity.UserEntity, error)
}

type userService struct {
	services *Services
}

func NewUserService(service *Services) UserService {
	return &userService{
		services: service,
	}
}

func (s *userService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, userFields data_fields.UserFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "userFields", userFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	userId := ""

	if id == nil || *id == "" {
		createFlow = true
		span.LogKV("flow", "create")
	} else {
		span.LogKV("flow", "update")
	}

	if createFlow {
		// generate id
		userId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelUser)
		if err != nil {
			tracing.TraceErr(span, err)
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
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, userId, model.NodeLabelUser)
		if err != nil || !exists {
			err = errors.New("user not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, userId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		var innerErr error
		if createFlow {
			innerErr = s.services.Neo4jRepositories.UserWriteRepository.CreateUserInTx(ctx, txWithPostCommit.Tx, tenant, userId, userFields)
		} else {
			innerErr = s.services.Neo4jRepositories.UserWriteRepository.UpdateUserInTx(ctx, txWithPostCommit.Tx, tenant, userId, userFields)
		}
		if innerErr != nil {
			return nil, innerErr
		}

		if userFields.ExternalSystemAvailable() {
			innerErr = s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, userId, model.NodeLabelUser, *userFields.ExternalSystem)
			if err != nil {
				return nil, innerErr
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if createFlow {
				err = s.services.RabbitMQService.PublishEvent(ctx, userId, model.USER, dto.CreateUser{userFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateUser"))
				}
			} else {
				err = s.services.RabbitMQService.PublishEvent(ctx, userId, model.USER, dto.UpdateUser{userFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateUser"))
				}
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createFlow {
		span.LogFields(log.Bool("response.userCreated", true))
	} else {
		span.LogFields(log.Bool("response.userUpdated", true))
	}
	return userId, nil
}

func (s *userService) GetById(parentCtx context.Context, userId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.services.Neo4jRepositories.UserReadRepository.GetUserById(ctx, common.GetContext(ctx).Tenant, userId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(node), nil
}

func (s *userService) GetAllUsersForTenant(ctx context.Context, tenant string) ([]*neo4jentity.UserEntity, error) {
	nodes, err := s.services.Neo4jRepositories.UserReadRepository.GetAllForTenant(ctx, tenant)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.FindFirstUserWithRolesByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	userDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, email)
	if err != nil {
		return nil, err
	}

	if userDbNode == nil {
		return nil, nil
	}

	return mapper.MapDbNodeToUserEntity(userDbNode), nil
}

func (s *userService) IsOwner(parentCtx context.Context, userId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.IsOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	isOwner, err := s.services.Neo4jRepositories.UserReadRepository.IsOwner(ctx, common.GetContext(ctx).Tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	return isOwner, nil
}

func (s *userService) GetContactOwner(parentCtx context.Context, contactId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetContactOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	ownerDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetOwnerForContact(ctx, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(ownerDbNode), nil
}

func (s *userService) GetNoteCreator(parentCtx context.Context, noteId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetNoteCreator")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	userDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetCreatorForNote(ctx, common.GetContext(ctx).Tenant, noteId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(userDbNode), nil
}

func (s *userService) GetUsersConnectedForContacts(ctx context.Context, contactIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUsersConnectedForContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetUsersConnectedForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
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

func (s *userService) GetUsersForEmails(parentCtx context.Context, emailIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUsersForEmails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), emailIds)
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

func (s *userService) GetUsersForPhoneNumbers(parentCtx context.Context, phoneNumberIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUsersForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), phoneNumberIds)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserOwnersForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationIDs", organizationIDs))

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllOwnersForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserOwnersForOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("opportunityIds", opportunityIds))

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllOwnersForOpportunities(ctx, common.GetTenantFromContext(ctx), opportunityIds)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserCreatorsForOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("opportunityIds", opportunityIds))

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllCreatorsForOpportunities(ctx, common.GetTenantFromContext(ctx), opportunityIds)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserCreatorsForOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("serviceLineItemIds", serviceLineItemIds))

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllCreatorsForServiceLineItems(ctx, common.GetTenantFromContext(ctx), serviceLineItemIds)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUsersWithMailboxes")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	mailboxes, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
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
		userNode, err := s.services.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, userEmail)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserCreatorsForContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractIds", contractIds))

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllCreatorsForContracts(ctx, common.GetTenantFromContext(ctx), contractIds)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserAuthorsForLogEntries")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("logEntryIDs", logEntryIDs))

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllAuthorsForLogEntries(ctx, common.GetTenantFromContext(ctx), logEntryIDs)
	if err != nil {
		tracing.TraceErr(span, err)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUserAuthorsForComments")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("commentIds", commentIds))

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllAuthorsForComments(ctx, common.GetTenantFromContext(ctx), commentIds)
	if err != nil {
		tracing.TraceErr(span, err)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUserForFlowSenders")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("flowSenderIds", flowSenderIds))

	users, err := s.services.Neo4jRepositories.UserReadRepository.GetAllSendersForFlowSenders(ctx, common.GetTenantFromContext(ctx), flowSenderIds)
	if err != nil {
		tracing.TraceErr(span, err)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUsers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("userIds", userIds))

	userDbNodes, err := s.services.Neo4jRepositories.UserReadRepository.GetUsers(ctx, common.GetTenantFromContext(ctx), userIds)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetAllOwnersForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	dbNodes, err := s.services.Neo4jRepositories.UserReadRepository.GetAllOwnersForOrganizations(ctx, tenant, organizationIds)
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
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetDistinctOrganizationOwners")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNodes, err := s.services.Neo4jRepositories.UserReadRepository.GetDistinctOrganizationOwners(ctx, common.GetTenantFromContext(ctx))
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetReminderOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	ownerDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetOwnerForReminder(ctx, common.GetContext(ctx).Tenant, reminderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(ownerDbNode), nil
}

func (s *userService) GetContractOwner(parentCtx context.Context, contractId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetContractOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	ownerDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetOwnerForContract(ctx, common.GetContext(ctx).Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(ownerDbNode), nil
}
