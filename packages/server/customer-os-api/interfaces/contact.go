package cosapi_interfaces

import (
	"context"

	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

type ContactService interface {
	Create(ctx context.Context, contact *ContactCreateData) (string, error)
	GetById(ctx context.Context, id string) (*neo4jentity.ContactEntity, error)
	GetFirstContactByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4jentity.ContactEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error)
	PermanentDelete(ctx context.Context, id string) (bool, error)
	RestoreFromArchive(ctx context.Context, contactId string) (bool, error)
	GetContactsForJobRoles(ctx context.Context, jobRoleIds []string) (*neo4jentity.ContactEntities, error)
	GetContactsForOrganization(ctx context.Context, organizationId string, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error)
	Merge(ctx context.Context, primaryContactId, mergedContactId string) error
	GetContactsForEmails(ctx context.Context, emailIds []string) (*neo4jentity.ContactEntities, error)
	GetContactsForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*neo4jentity.ContactEntities, error)
	RemoveLocation(ctx context.Context, contactId string, locationId string) error
	CustomerContactCreate(ctx context.Context, entity *CustomerContactCreateData) (*model.CustomerContact, error)
	GetContactCountByOrganizations(ctx context.Context, ids []string) (map[string]int64, error)
}

type ContactCreateData struct {
	ContactEntity     *neo4jentity.ContactEntity
	EmailEntity       *neo4jentity.EmailEntity
	PhoneNumberEntity *neo4jentity.PhoneNumberEntity
	ExternalReference *neo4jentity.ExternalSystemEntity
	Source            neo4jentity.DataSource
	SocialUrl         string
	AppSource         string
}

type CustomerContactCreateData struct {
	ContactEntity *neo4jentity.ContactEntity
	EmailEntity   *neo4jentity.EmailEntity
}
