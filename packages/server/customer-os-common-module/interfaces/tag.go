package interfaces

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
)

type TagService interface {
	Save(ctx context.Context, tx *neo4j.ManagedTransaction, inputTag *neo4jentity.TagEntity) (*neo4jentity.TagEntity, error)
	AddTagToEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, tagId, tagName string) (string, error)
	RemoveTagFromEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, tagId string) error
	Update(ctx context.Context, tagId string, name, colorCode *string) error
	UnlinkAndDelete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) (*neo4jentity.TagEntities, error)
	GetById(ctx context.Context, tagId string) (*neo4jentity.TagEntity, error)
	GetTagByEntityTypeAndName(ctx context.Context, entityType model.EntityType, name string) (*neo4jentity.TagEntity, error)
	GetTagsByEntityType(ctx context.Context, entityType model.EntityType) (*neo4jentity.TagEntities, error)
	GetTagsForContacts(ctx context.Context, contactIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.TagEntities, error)
	GetTagsForLogEntries(ctx context.Context, logEntryIds []string) (*neo4jentity.TagEntities, error)
}
