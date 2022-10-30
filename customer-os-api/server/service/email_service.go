package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
)

type EmailService interface {
}

type emailService struct {
	driver *neo4j.Driver
}

func NewEmailService(driver *neo4j.Driver) EmailService {
	return &emailService{
		driver: driver,
	}
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
