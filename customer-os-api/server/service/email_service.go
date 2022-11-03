package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type EmailService interface {
	FindAllForContact(ctx context.Context, obj *model.Contact) (*entity.EmailEntities, error)
	MergeEmailToContact(ctx context.Context, id string, toEntity *entity.EmailEntity) (*entity.EmailEntity, error)
	Delete(ctx context.Context, contactId string, email string) (bool, error)
	DeleteById(ctx context.Context, contactId string, emailId string) (bool, error)
}

type emailService struct {
	driver *neo4j.Driver
}

func NewEmailService(driver *neo4j.Driver) EmailService {
	return &emailService{
		driver: driver,
	}
}

func (s *emailService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.EmailEntities, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c:Contact {id:$id})-[r:EMAILED_AT]->(e:Email) 
				RETURN e, r`,
			map[string]interface{}{
				"id":     contact.ID,
				"tenant": common.GetContext(ctx).Tenant})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
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

func (s *emailService) MergeEmailToContact(ctx context.Context, contactId string, entity *entity.EmailEntity) (*entity.EmailEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		if entity.Primary == true {
			err := setOtherContactEmailsNonPrimaryInTx(ctx, contactId, entity.Email, tx)
			if err != nil {
				return nil, err
			}
		}
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (c)-[r:EMAILED_AT]->(e:Email {email: $email})
            ON CREATE SET e.label=$label, r.primary=$primary
            ON MATCH SET e.label=$label, r.primary=$primary
			RETURN e, r`,
			map[string]interface{}{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"email":     entity.Email,
				"label":     entity.Label,
				"primary":   entity.Primary,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record, nil
	})
	if err != nil {
		return nil, err
	}

	var emailEntity = s.mapDbNodeToEmailEntity(queryResult.(*db.Record).Values[0].(dbtype.Node))
	s.addDbRelationshipToEmailEntity(queryResult.(*db.Record).Values[1].(dbtype.Relationship), emailEntity)
	return emailEntity, nil
}

func (s *emailService) Delete(ctx context.Context, contactId string, email string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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

func addEmailToContactInTx(ctx context.Context, contactId string, input entity.EmailEntity, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})
			MERGE (e:Email {
				id: randomUUID(),
				email: $email,
				label: $label
			})<-[:EMAILED_AT {primary:$primary}]-(c)
			RETURN e`,
		map[string]interface{}{
			"contactId": contactId,
			"primary":   input.Primary,
			"email":     input.Email,
			"label":     input.Label,
		})
	return err
}

func setOtherContactEmailsNonPrimaryInTx(ctx context.Context, contactId string, email string, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[r:EMAILED_AT]->(e:Email)
			WHERE e.email <> $email
            SET r.primary=false`,
		map[string]interface{}{
			"tenant":    common.GetContext(ctx).Tenant,
			"contactId": contactId,
			"email":     email,
		})
	return err
}

func (s *emailService) mapDbNodeToEmailEntity(node dbtype.Node) *entity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.EmailEntity{
		Id:    utils.GetStringPropOrEmpty(props, "id"),
		Email: utils.GetStringPropOrEmpty(props, "email"),
		Label: utils.GetStringPropOrEmpty(props, "label"),
	}
	return &result
}

func (s *emailService) addDbRelationshipToEmailEntity(relationship dbtype.Relationship, emailEntity *entity.EmailEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	emailEntity.Primary = utils.GetBoolPropOrFalse(props, "primary")
}
