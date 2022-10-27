package service

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactService interface {
	Create(contact *entity.ContactNode) (*entity.ContactNode, error)
	FindAll() ([]entity.ContactNode, error)
	FindContactById(id string) (*entity.ContactNode, error)
}

type neo4jContactService struct {
	driver *neo4j.Driver
}

func NewContactService(driver *neo4j.Driver) ContactService {
	return &neo4jContactService{
		driver: driver,
	}
}

func (s *neo4jContactService) Create(newContact *entity.ContactNode) (*entity.ContactNode, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			CREATE (c:Contact {
				  id: randomUUID(),
				  firstName: $firstName,
				  lastName: $lastName,
				  label: $label,
				  contactType: $contactType
			})
			RETURN c { .id, .firstName, .lastName, .label, .contactType } as c`,
			map[string]interface{}{
				"firstName":   newContact.FirstName,
				"lastName":    newContact.LastName,
				"label":       newContact.Label,
				"contactType": newContact.ContactType,
			})

		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}
	contact := entity.ContactNode{}
	mapstructure.Decode(result.(map[string]interface{}), &contact)

	return &contact, nil
}

func (s *neo4jContactService) FindContactById(id string) (*entity.ContactNode, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (c:Contact) WHERE c.id=$id RETURN c`,
			map[string]interface{}{
				"id": id,
			})

		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	contact := entity.ContactNode{}
	err = mapstructure.Decode(utils.GetPropsFromNode(result.(dbtype.Node)), &contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (cs *neo4jContactService) FindAll() ([]entity.ContactNode, error) {

	// Open a new Session
	session := (*cs.driver).NewSession(neo4j.SessionConfig{})

	// Close the session once this function has completed
	defer (*cs.driver).Close()

	// Execute a query in a new Read Transaction
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get an array of IDs for the User's favorite movies

		// Retrieve a list of movies

		result, err := tx.Run("MATCH (c:Contact) RETURN c { .* } AS contact", map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		// Get a list of Movies from the Result
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			person, _ := record.Get("contact")
			results = append(results, person.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	return results.([]entity.ContactNode), nil
}
