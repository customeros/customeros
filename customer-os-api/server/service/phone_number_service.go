package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
)

type PhoneNumberService interface {
}

type phoneNumberService struct {
	driver *neo4j.Driver
}

func NewPhoneNumberService(driver *neo4j.Driver) PhoneNumberService {
	return &phoneNumberService{
		driver: driver,
	}
}

func addPhoneNumberToContactInTx(ctx context.Context, contactId string, input entity.PhoneNumberEntity, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})
			CREATE (p:PhoneNumber {
				  number: $number,
				  label: $label
			})<-[:CALLED_AT]-(c)
			RETURN p`,
		map[string]interface{}{
			"contactId": contactId,
			"number":    input.Number,
			"label":     input.Label,
		})

	return err
}
