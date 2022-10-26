package service

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/model"
)

type ContactService interface {
	FindAll() ([]model.ContactDB, error)
	FindAllByName(name string) ([]model.ContactDB, error)
	Create(contact model.ContactDB) (*model.ContactDB, error)
}

type neo4jContactService struct {
}

func NewContactService() ContactService {
	return &neo4jContactService{}
}

func (cs *neo4jContactService) Create(aContact model.ContactDB) (*model.ContactDB, error) {
	contact := model.ContactDB{}
	driver, err := neo4j.NewDriver("neo4j://neo4j-customer-os.openline-development.svc.cluster.local:7687",
		neo4j.BasicAuth("neo4j", "StrongLocalPa$$", ""))
	if err != nil {
		return &contact, err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			CREATE (c:ContactDB {
				  id: randomUuid(),
				  firstName: $firstName,
				  lastName: $lastName,
				  label: $label,
				  contactType: $contactType
			})
			RETURN c { .id, .firstName, .lastName, .label, .contactType } as c`,
			map[string]interface{}{
				"firstName":   aContact.LastName,
				"lastName":    aContact.LastName,
				"label":       aContact.Label,
				"contactType": aContact.ContactType,
			})

		// Extract safe properties from the aContact node (`c`) in the first row
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		contact, _ := record.Get("c")
		return contact, nil
	})
	if err != nil {
		return &contact, err
	}

	contactData := result.(map[string]interface{})

	mapstructure.Decode(contactData, &contact)

	if err != nil {
		return &contact, err
	}
	return &contact, err
}

func (cs *neo4jContactService) FindAll() ([]model.ContactDB, error) {
	//TODO implement me
	panic("implement me")
}

func (n neo4jContactService) FindAllByName(name string) ([]model.ContactDB, error) {
	//TODO implement me
	panic("implement me")
}
