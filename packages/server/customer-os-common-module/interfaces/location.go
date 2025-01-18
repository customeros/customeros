package interfaces

import (
	"context"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type LocationService interface {
	SetContactService(contact ContactService)
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.LocationEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*neo4jentity.LocationEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*neo4jentity.LocationEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.LocationEntities, error)
	ExtractAndEnrichLocation(ctx context.Context, tenant, address string) (*data_fields.LocationFields, error)
	Create(ctx context.Context, commit *utils.TxWithPostCommit, locationFields data_fields.LocationFields, linkWith *common_srv.LinkWith) (string, error)
}
