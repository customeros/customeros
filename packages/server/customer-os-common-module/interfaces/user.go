package interfaces

import (
	"context"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
