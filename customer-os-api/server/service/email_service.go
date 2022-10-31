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
              			(c:Contact {id:$id})-[:EMAILED_AT]->(e:Email) 
				RETURN e`,
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
		emailEntities = append(emailEntities, *emailEntity)
	}

	return &emailEntities, nil
}

func addEmailToContactInTx(ctx context.Context, contactId string, input entity.EmailEntity, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})
			CREATE (p:Email {
				  email: $email,
				  label: $label
			})<-[:EMAILED_AT]-(c)
			RETURN p`,
		map[string]interface{}{
			"contactId": contactId,
			"email":     input.Email,
			"label":     input.Label,
		})

	return err
}

func (s *emailService) mapDbNodeToEmailEntity(dbContactGroupNode dbtype.Node) *entity.EmailEntity {
	props := utils.GetPropsFromNode(dbContactGroupNode)
	result := entity.EmailEntity{
		Email: utils.GetStringPropOrEmpty(props, "email"),
		Label: utils.GetStringPropOrEmpty(props, "label"),
	}
	return &result
}
