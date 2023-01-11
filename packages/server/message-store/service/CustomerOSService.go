package service

import (
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository/entity"
	"regexp"
	"strings"
	"time"
)

type customerOSService struct {
	driver               *neo4j.Driver
	postgresRepositories *repository.PostgresRepositories
}

type CustomerOSService interface {
	ContactByIdExists(contactId string) (bool, error)
	ContactByPhoneExists(e164 string) (bool, error)

	GetUserByEmail(email string) (*UserInfo, error)
	GetContactById(contactId string) (*ContactInfo, error)
	GetContactByEmail(email string) (*ContactInfo, error)

	CreateContactWithEmail(email string) (string, error)
	CreateContactWithPhone(phone string) (string, error)

	GetWebChatConversationIdWithContactInitiator(contactId string) (string, error)
	CreateConversation(tenant string, initiatorId string, initiatorType entity.SenderType, channel string) (string, error)
}

type UserInfo struct {
	Id string
}

type ContactInfo struct {
	Id string
}

func parseEmail(email string) (string, string) {
	re := regexp.MustCompile("^\"{0,1}([^\"]*)\"{0,1}[ ]*<(.*)>$")
	matches := re.FindStringSubmatch(strings.Trim(email, " "))
	if matches != nil {
		return strings.Trim(matches[1], " "), matches[2]
	}
	return "", email
}

func (s *customerOSService) ContactByIdExists(contactId string) (bool, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	params := map[string]interface{}{
		"contactId": contactId,
	}
	query := "MATCH (n:Contact {id:$contactId}) RETURN count(*)"
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(query, params)
		return nil, err
	})

	//TODO check here if the error is correct
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *customerOSService) ContactByPhoneExists(e164 string) (bool, error) {

	//graphqlRequest := graphql.NewRequest(`
	//			query ($e164: String!) {
	//				contact_ByPhone(e164: $e164){firstName,lastName,id}
	//			}
	//`)
	//
	//graphqlRequest.Var("e164", e164)
	//graphqlRequest.Header.Add("X-Openline-API-KEY", s.config.Service.CustomerOsAPIKey)
	//var graphqlResponse map[string]map[string]string
	//if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
	//	return false, err
	//}
	return true, nil
}

func (s *customerOSService) GetUserByEmail(email string) (*UserInfo, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`MATCH (u:User{email:$email}) return u`,
			map[string]any{
				"email": email,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	}
	if len(dbRecords.([]*db.Record)) == 0 {
		return nil, nil
	}
	dbNode := (dbRecords.([]*db.Record))[0].Values[0].(dbtype.Node)
	props := utils.GetPropsFromNode(dbNode)

	userInfo := UserInfo{Id: utils.GetStringPropOrEmpty(props, "id")}

	return &userInfo, nil
}

func (s *customerOSService) GetContactById(contactId string) (*ContactInfo, error) {

	//graphqlRequest := graphql.NewRequest(`
	//			query ($id: ID!) {
	//				contact(id: $id){
	//					firstName,
	//					lastName,
	//					id,
	//					phoneNumbers {
	//					   e164
	//					 }, emails {
	//					   email
	//					 }
	//  				}
	//			}
	//`)
	//
	//graphqlRequest.Var("id", id)
	//graphqlRequest.Header.Add("X-Openline-API-KEY", s.config.Service.CustomerOsAPIKey)
	//var graphqlResponse contactResponse
	//if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
	//	log.Printf("Grapql got error %s", err.Error())
	//	return nil, err
	//}
	//contactInfo := &ContactInfo{firstName: graphqlResponse.Contact.FirstName,
	//	lastName: graphqlResponse.Contact.LastName,
	//	id:       graphqlResponse.Contact.ID}
	//if len(graphqlResponse.Contact.Emails) > 0 {
	//	contactInfo.email = &graphqlResponse.Contact.Emails[0].Email
	//}
	//if len(graphqlResponse.Contact.PhoneNumbers) > 0 {
	//	contactInfo.phone = &graphqlResponse.Contact.PhoneNumbers[0].E164
	//}
	return nil, nil
}

func (s *customerOSService) GetContactByEmail(email string) (*ContactInfo, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact)-[:EMAILED_AT]->(p:Email {email:$email})
            RETURN c`,
			map[string]any{
				"tenant": "openline", //TODO discuss with customerOS team
				"email":  email,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	if err != nil {
		return nil, err
	}
	if len(dbRecords.([]*db.Record)) == 0 {
		return nil, errors.New("no contact found")
	}
	dbNode := (dbRecords.([]*db.Record))[0].Values[0].(dbtype.Node)
	props := utils.GetPropsFromNode(dbNode)

	contactInfo := ContactInfo{Id: utils.GetStringPropOrEmpty(props, "id")}

	return &contactInfo, nil
}

func (s *customerOSService) CreateContactWithEmail(tenant string, email string) (string, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		//create the contact
		contactQuery := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (c:Contact {id:randomUUID()})-[:CONTACT_BELONGS_TO_TENANT]->(t) ON CREATE SET" +
			" c.createdAt=$createdAt, " +
			" c:%s " +
			" RETURN c"

		contactQueryResult, err := tx.Run(fmt.Sprintf(contactQuery, "Contact_"+tenant),
			map[string]interface{}{
				"tenant":    tenant,
				"createdAt": time.Now().UTC(),
			})

		node, err := utils.ExtractSingleRecordFirstValueAsNode(contactQueryResult, err)
		if err != nil {
			return nil, err
		}

		//create the email
		emailQuery := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
			" MERGE (c)-[r:EMAILED_AT]->(e:Email {email: $email}) " +
			" ON CREATE SET e.label=$label, r.primary=$primary, e.id=randomUUID(), e:%s " +
			" ON MATCH SET e.label=$label, r.primary=$primary " +
			" RETURN e, r"

		contactId := fmt.Sprintf("%d", node.Id)
		_, err = tx.Run(fmt.Sprintf(emailQuery, "Email_"+tenant),
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"email":     email,
				"label":     "WORK",
				"primary":   true,
			})

		if err != nil {
			return nil, err
		}

		return contactId, nil
	})
	if err != nil {
		return "", err
	} else {
		return result.(string), nil
	}
}

// TODO
func (s *customerOSService) GetWebChatConversationIdWithContactInitiator(contactId string) (string, error) {
	//graphqlRequest := graphql.NewRequest(`
	//	//mutation CreateContact ($firstName: String!, $lastName: String!, $e164: String!) {
	//	//  contact_Create(input: {
	//    //    firstName: $firstName,
	//	//	lastName: $lastName,
	//	//    phoneNumber:{e164:  $e164, label: WORK}
	//	//  }) {
	//	//	  id
	//	//  }
	//	//}
	//`)
	//
	//graphqlRequest.Header.Add("X-Openline-API-KEY", s.config.Service.CustomerOsAPIKey)
	//var graphqlResponse map[string]map[string]string
	//if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
	//	return "", err
	//}
	//return graphqlResponse["contact_Create"]["id"], nil
	return "", nil
}

func (s *customerOSService) CreateContactWithPhone(phone string) (string, error) {
	//graphqlRequest := graphql.NewRequest(`
	//	mutation CreateContact ($firstName: String!, $lastName: String!, $e164: String!) {
	//	  contact_Create(input: {
	//        firstName: $firstName,
	//		lastName: $lastName,
	//	    phoneNumber:{e164:  $e164, label: WORK}
	//	  }) {
	//		  id
	//	  }
	//	}
	//`)
	//
	//graphqlRequest.Var("firstName", firstName)
	//graphqlRequest.Var("lastName", lastName)
	//graphqlRequest.Var("e164", phone)
	//graphqlRequest.Header.Add("X-Openline-API-KEY", s.config.Service.CustomerOsAPIKey)
	//var graphqlResponse map[string]map[string]string
	//if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
	//	return "", err
	//}
	//return graphqlResponse["contact_Create"]["id"], nil
	return "", nil
}

func (s *customerOSService) CreateConversation(tenant string, initiatorId string, initiatorType entity.SenderType, channel entity.EventType) (string, error) {
	session := (*s.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	contactIds := []string{}
	userIds := []string{}

	if initiatorType == entity.CONTACT {
		contactIds = append(contactIds, initiatorId)
	} else if initiatorType == entity.USER {
		userIds = append(userIds, initiatorId)
	}

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:randomUUID()}) " +
			" ON CREATE SET o.startedAt=$startedAt, o.messageCount=0, o.channel=$channel, o.status=$status, o:%s " +
			" %s %s " +
			" RETURN DISTINCT o"
		queryLinkWithContacts := ""
		if len(contactIds) > 0 {
			queryLinkWithContacts = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" MERGE (c)-[:PARTICIPATES]->(o) "
		}
		queryLinkWithUsers := ""
		if len(userIds) > 0 {
			queryLinkWithUsers = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" MERGE (u)-[:PARTICIPATES]->(o) "
		}
		queryResult, err := tx.Run(fmt.Sprintf(query, "Conversation_"+tenant, queryLinkWithContacts, queryLinkWithUsers),
			map[string]interface{}{
				"tenant":     tenant,
				"status":     "ACTIVE",
				"startedAt":  time.Now().UTC(),
				"channel":    channel,
				"contactIds": contactIds,
				"userIds":    userIds,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return "", err
	} else {
		dbNode := result.(*dbtype.Node)
		return fmt.Sprintf("%d", dbNode.Id), err
	}
}

func NewCustomerOSService(driver *neo4j.Driver, postgresRepositories *repository.PostgresRepositories) *customerOSService {
	customerOsService := new(customerOSService)
	customerOsService.driver = driver
	customerOsService.postgresRepositories = postgresRepositories
	return customerOsService
}
