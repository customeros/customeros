package interfaces

import (
	"context"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type LocationService interface {
	GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.LocationEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*neo4jentity.LocationEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*neo4jentity.LocationEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.LocationEntities, error)
	ExtractAndEnrichLocation(ctx context.Context, tenant, address string) (*data_fields.LocationFields, error)
	Create(ctx context.Context, commit *utils.TxWithPostCommit, locationFields data_fields.LocationFields, linkWith *common_srv.LinkWith) (string, error)
}
