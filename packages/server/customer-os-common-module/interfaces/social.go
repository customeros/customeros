package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type SocialService interface {
	SetContactService(contact ContactService)
	IsInitialized() bool

	GetById(ctx context.Context, socialId string) (*neo4j_entity.SocialEntity, error)
	AddSocialToEntity(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkWith common_srv.LinkWith, socialEntity neo4j_entity.SocialEntity) (string, error)
	RemoveSocialFromEntity(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkWith common_srv.LinkWith, socialId string) error
	Update(ctx context.Context, entity neo4j_entity.SocialEntity) (*neo4j_entity.SocialEntity, error)
	PermanentlyDelete(ctx context.Context, tenant, socialId string) error
	GetAllForEntities(ctx context.Context, tenant string, linkedEntityType model.EntityType, linkedEntityIds []string) (*neo4j_entity.SocialEntities, error)
	GetAllLinkedinForEntities(ctx context.Context, tenant string, linkedEntityType model.EntityType, linkedEntityIds []string) (*neo4j_entity.SocialEntities, error)
}
