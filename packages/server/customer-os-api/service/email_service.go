package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type EmailService interface {
	FindAllFor(ctx context.Context, entityType repository.EntityType, entityId string) (*entity.EmailEntities, error)
	MergeEmailTo(ctx context.Context, entityType repository.EntityType, entityId string, entity *entity.EmailEntity) (*entity.EmailEntity, error)
	UpdateEmailFor(ctx context.Context, entityType repository.EntityType, entityId string, entity *entity.EmailEntity) (*entity.EmailEntity, error)
	Delete(ctx context.Context, contactId string, email string) (bool, error)
	DeleteById(ctx context.Context, contactId string, emailId string) (bool, error)
}

type emailService struct {
	repositories *repository.Repositories
}

func NewEmailService(repository *repository.Repositories) EmailService {
	return &emailService{
		repositories: repository,
	}
}

func (s *emailService) getDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *emailService) FindAllFor(ctx context.Context, entityType repository.EntityType, entityId string) (*entity.EmailEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	queryResult, err := s.repositories.EmailRepository.FindAllFor(session, common.GetContext(ctx).Tenant, entityType, entityId)
	if err != nil {
		return nil, err
	}

	emailEntities := entity.EmailEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		emailEntity := s.mapDbNodeToEmailEntity(dbRecord.Values[0].(dbtype.Node))
		s.addDbRelationshipToEmailEntity(dbRecord.Values[1].(dbtype.Relationship), emailEntity)
		emailEntities = append(emailEntities, *emailEntity)
	}

	return &emailEntities, nil
}

func (s *emailService) MergeEmailTo(ctx context.Context, entityType repository.EntityType, entityId string, entity *entity.EmailEntity) (*entity.EmailEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	var err error
	var emailNode *dbtype.Node
	var emailRelationship *dbtype.Relationship

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		if entity.Primary == true {
			err := s.repositories.EmailRepository.SetOtherEmailsNonPrimaryInTx(tx, common.GetContext(ctx).Tenant, entityType, entityId, entity.Email)
			if err != nil {
				return nil, err
			}
		}
		emailNode, emailRelationship, err = s.repositories.EmailRepository.MergeEmailToInTx(tx, common.GetContext(ctx).Tenant, entityType, entityId, *entity)
		return nil, err
	})
	if err != nil {
		return nil, err
	}

	var emailEntity = s.mapDbNodeToEmailEntity(*emailNode)
	s.addDbRelationshipToEmailEntity(*emailRelationship, emailEntity)
	return emailEntity, nil
}

func (s *emailService) UpdateEmailFor(ctx context.Context, entityType repository.EntityType, entityId string, entity *entity.EmailEntity) (*entity.EmailEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	var err error
	var emailNode *dbtype.Node
	var emailRelationship *dbtype.Relationship

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		emailNode, emailRelationship, err = s.repositories.EmailRepository.UpdateEmailByInTx(tx, common.GetContext(ctx).Tenant, entityType, entityId, *entity)
		if err != nil {
			return nil, err
		}
		if entity.Primary == true {
			err := s.repositories.EmailRepository.SetOtherEmailsNonPrimaryInTx(tx, common.GetContext(ctx).Tenant, entityType, entityId, entity.Email)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	})
	if err != nil {
		return nil, err
	}

	var emailEntity = s.mapDbNodeToEmailEntity(*emailNode)
	s.addDbRelationshipToEmailEntity(*emailRelationship, emailEntity)
	return emailEntity, nil
}

func (s *emailService) Delete(ctx context.Context, contactId string, email string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact {id:$id})-[:EMAILED_AT]->(p:Email {email:$email})
            DETACH DELETE p
			`,
			map[string]interface{}{
				"id":     contactId,
				"email":  email,
				"tenant": common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *emailService) DeleteById(ctx context.Context, contactId string, emailId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact {id:$contactId})-[:EMAILED_AT]->(p:Email {id:$emailId})
            DETACH DELETE p
			`,
			map[string]interface{}{
				"contactId": contactId,
				"emailId":   emailId,
				"tenant":    common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *emailService) mapDbNodeToEmailEntity(node dbtype.Node) *entity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.EmailEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Email:         utils.GetStringPropOrEmpty(props, "email"),
		Label:         utils.GetStringPropOrEmpty(props, "label"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:     utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrNow(props, "updatedAt"),
	}
	return &result
}

func (s *emailService) addDbRelationshipToEmailEntity(relationship dbtype.Relationship, emailEntity *entity.EmailEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	emailEntity.Primary = utils.GetBoolPropOrFalse(props, "primary")
}
