package service

import (
	"errors"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	commonModuleService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository/entity"
	"golang.org/x/net/context"
	"time"
)

type CustomerOSService struct {
	driver               *neo4j.DriverWithContext
	postgresRepositories *repository.PostgresRepositories
	commonStoreService   *commonStoreService
	graphqlClient        *graphql.Client
	conf                 *config.Config
}

type EmailContent struct {
	MessageId string   `json:"messageId"`
	Html      string   `json:"html"`
	Subject   string   `json:"subject"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Cc        []string `json:"cc"`
	Bcc       []string `json:"bcc"`
	InReplyTo []string `json:"InReplyTo"`
	Reference []string `json:"Reference"`
}

type CustomerOSServiceInterface interface {
	ContactByIdExists(contactId string) (bool, error)
	ContactByPhoneExists(e164 string) (bool, error)

	GetUserByEmail(email string) (*User, error)
	GetContactById(id string) (*Contact, error)
	GetContactByEmail(email string) (*Contact, error)

	CreateContactWithEmail(tenant string, email string) (*Contact, error)
	CreateContactWithPhone(tenant string, phone string) (*Contact, error)

	ConversationByIdExists(tenant string, conversationId string) (bool, error)

	GetConversations(tenant string) ([]Conversation, error)
	GetConversationById(tenant string, conversationId string) (Conversation, error)

	GetWebChatConversationWithContactInitiator(tenant string, contactId string) (*Conversation, error)
	GetWebChatConversationWithUserInitiator(tenant string, userId string) (*Conversation, error)

	CreateConversation(tenant string, initiatorId string, initiatorFirstName string, initiatorLastName string, initiatorUsername string, initiatorType entity.SenderType, channel entity.EventType) (*Conversation, error)
	UpdateConversation(tenant string, conversationId string, participantId string, participantType entity.SenderType, lastSenderFirstName string, lastSenderLastName string, lastContentPreview string) (string, error)
}

type Conversation struct {
	Id        string
	StartedAt time.Time
	UpdatedAt time.Time
	Channel   string
	Status    string

	InitiatorFirstName  string
	InitiatorLastName   string
	InitiatorUsername   string
	InitiatorType       string
	LastSenderId        string
	LastSenderType      string
	LastSenderFirstName string
	LastSenderLastName  string
	LastContentPreview  string
}

type emailObject struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

type phoneNumberObject struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

type Contact struct {
	FirstName    string              `json:"firstName"`
	LastName     string              `json:"lastName"`
	Id           string              `json:"id"`
	Emails       []emailObject       `json:"emails"`
	phoneNumbers []phoneNumberObject `json:"phoneNumbers"`
}

type User struct {
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Id        string        `json:"id"`
	Emails    []emailObject `json:"emails"`
}

const contactFieldSelection = `firstName,lastName,id,
    				emails {
  					  email
                      primary
					  label
            		},
					phoneNumbers {
						  e164,
						  primary,
						  label
            		}`

const userFieldSelection = `firstName,lastName,id,
    				emails {
  					  email
                      primary
					  label
            		},
					phoneNumbers {
						  e164,
						  primary,
						  label
            		}`

func (s *CustomerOSService) addHeadersToGraphRequest(req *graphql.Request, ctx context.Context) error {
	req.Header.Add("X-Openline-API-KEY", s.conf.Service.CustomerOsAPIKey)
	user, err := commonModuleService.GetUsernameMetadataForGRPC(ctx)
	if err != nil {
		return err
	}
	req.Header.Add("X-Openline-USERNAME", *user)
	return nil
}

func (s *CustomerOSService) ContactByIdExists(ctx context.Context, contactId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.driver)
	defer session.Close(ctx)

	params := map[string]interface{}{
		"contactId": contactId,
	}
	query := "MATCH (n:Contact {id:$contactId}) RETURN count(*)"
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	//TODO check here if the error is correct
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *CustomerOSService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	graphqlRequest := graphql.NewRequest(`
  				query ($email: String!) {
  					user_ByEmail(email: $email){` + userFieldSelection + `}
  				}
    `)
	graphqlRequest.Var("email", email)

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx)

	if err != nil {
		return nil, err
	}

	var graphqlResponse struct {
		userByEmail User `json:"user_ByEmail"`
	}

	if err = s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return &graphqlResponse.userByEmail, nil
}

func (s *CustomerOSService) GetContactById(ctx context.Context, id string) (*Contact, error) {
	graphqlRequest := graphql.NewRequest(`
  				query (id: ID!) {
  					contact(id: $id){` + contactFieldSelection + `}
  				}
    `)
	graphqlRequest.Var("id", id)

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx)
	if err != nil {
		return nil, err
	}

	var graphqlResponse struct {
		contact Contact `json:"contact"`
	}

	if err = s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return &graphqlResponse.contact, nil
}

func (s *CustomerOSService) GetContactByEmail(ctx context.Context, email string) (*Contact, error) {
	graphqlRequest := graphql.NewRequest(`
  				query ($email: String!) {
  					contact_ByEmail(email: $email){` + contactFieldSelection + `}
  				}
    `)
	graphqlRequest.Var("email", email)

	err := s.addHeadersToGraphRequest(graphqlRequest, ctx)
	if err != nil {
		return nil, err
	}

	var graphqlResponse struct {
		contactByEmail Contact `json:"contact_ByEmail"`
	}

	if err = s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return &graphqlResponse.contactByEmail, nil
}

func (s *CustomerOSService) CreateContactWithEmail(ctx context.Context, tenant string, email string) (*Contact, error) {
	graphqlRequest := graphql.NewRequest(`
		mutation CreateContact ($email: String!) {
		  contact_Create(input: {
		  email:{email:  $email, label: WORK}}) {
			id
          }
		}
    `)

	graphqlRequest.Var("email", email)
	err := s.addHeadersToGraphRequest(graphqlRequest, ctx)

	if err != nil {
		return nil, err
	}

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	id := graphqlResponse["contact_Create"]["id"]

	return s.GetContactById(ctx, id)

}

func (s *CustomerOSService) ConversationByIdExists(ctx context.Context, tenant string, conversationId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, "MATCH (c:Conversation_"+tenant+"{id:$conversationId}) RETURN c",
			map[string]any{
				"conversationId": conversationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})

	if err != nil {
		return false, err
	}

	if len(dbRecords.([]*db.Record)) == 0 {
		return false, nil
	}

	return true, nil
}

func (s *CustomerOSService) GetConversations(ctx context.Context, tenant string, onlyContacts bool) ([]Conversation, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		//todo move order by as param
		cypher := ""
		if onlyContacts {
			//cypher = "match (t:Tenant{name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email {email:$user}), (u)-[:PARTICIPATES]->(o:Conversation)<-[:PARTICIPATES]-(c:Contact) return distinct o"
			cypher = "match (t:Tenant{name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User), (u)-[:PARTICIPATES]->(o:Conversation)<-[:PARTICIPATES]-(c:Contact) return distinct o order by o.updatedAt desc"
		} else {
			cypher = "MATCH (c:Conversation_" + tenant + ") RETURN c order by c.updatedAt desc"
		}

		if queryResult, err := tx.Run(ctx, cypher, map[string]any{
			"tenant": tenant,
		}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})

	if err != nil {
		return nil, err
	}

	var conversations []Conversation
	for _, v := range dbRecords.([]*neo4j.Record) {
		node := v.Values[0].(neo4j.Node)
		conversations = append(conversations, *mapNodeToConversation(&node))
	}

	return conversations, nil
}

func (s *CustomerOSService) GetConversationById(ctx context.Context, tenant string, conversationId string) (*Conversation, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.driver)
	defer session.Close(ctx)

	conversationNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		if queryResult, err := tx.Run(ctx, "MATCH (c:Conversation_"+tenant+"{id: $conversationId}) RETURN c", map[string]any{
			"conversationId": conversationId,
		}); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single(ctx)
			if err != nil {
				return nil, err
			}
			if record == nil {
				return nil, errors.New("conversation not found")
			}

			return utils.NodePtr(record.Values[0].(neo4j.Node)), nil
		}
	})
	if err != nil {
		return nil, err
	}

	return mapNodeToConversation(conversationNode.(*dbtype.Node)), nil
}

func (s *CustomerOSService) GetConversationParticipants(ctx context.Context, tenant string, conversationId string) ([]string, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, "MATCH (c:Conversation_"+tenant+"{id:$conversationId})<-[PARTICIPATES]-(p)-[HAS]->(e:Email) RETURN DISTINCT(e.email) AS email",
			map[string]interface{}{
				"conversationId": conversationId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return []string{}, err
	}
	emails := make([]string, 0)
	if len(records.([]*neo4j.Record)) > 0 {
		for _, record := range records.([]*neo4j.Record) {
			emails = append(emails, record.Values[0].(string))
		}
		return emails, nil
	} else {
		return []string{}, nil
	}
}

func (s *CustomerOSService) CreateConversation(ctx context.Context, tenant string, initiator Participant, initiatorUsername string, channel entity.EventType, threadId string) (*Conversation, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.driver)
	defer session.Close(ctx)

	contactIds := []string{}
	userIds := []string{}
	initiatorTypeStr := ""

	if initiator.Type == entity.CONTACT {
		contactIds = append(contactIds, initiator.Id)
		initiatorTypeStr = "CONTACT"
	} else if initiator.Type == entity.USER {
		userIds = append(userIds, initiator.Id)
		initiatorTypeStr = "USER"
	}

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:randomUUID()}) " +
			" ON CREATE SET o.startedAt=$startedAt, o.updatedAt=$updatedAt, o.threadId=$threadId, " +
			" o.messageCount=0, o.channel=$channel, o.status=$status, o.source=$source, o.sourceOfTruth=$sourceOfTruth, " +
			" o.initiatorFirstName=$initiatorFirstName, o.initiatorLastName=$initiatorLastName, o.initiatorUsername=$initiatorUsername, o.initiatorType=$initiatorType, " +
			" o.source=$source, o.sourceOfTruth=$sourceOfTruth, " +
			" o.lastSenderId=$lastSenderId, o.lastSenderType=$lastSenderType, o.lastSenderFirstName=$lastSenderFirstName, o.lastSenderLastName=$lastSenderLastName, o.lastContentPreview=$lastContentPreview, " +
			" o.appSource=$appSource, o:%s " +
			" %s %s " +
			" RETURN DISTINCT o"
		queryLinkWithContacts := ""
		if len(contactIds) > 0 {
			queryLinkWithContacts = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" MERGE (c)-[:PARTICIPATES]->(o) " +
				" WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" MERGE (c)-[:INITIATED]->(o) "
		}
		queryLinkWithUsers := ""
		if len(userIds) > 0 {
			queryLinkWithUsers = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" MERGE (u)-[:PARTICIPATES]->(o) " +
				" WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" MERGE (u)-[:INITIATED]->(o) "
		}
		utc := time.Now().UTC()
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Conversation_"+tenant, queryLinkWithContacts, queryLinkWithUsers),
			map[string]interface{}{
				"tenant":              tenant,
				"source":              "openline",
				"sourceOfTruth":       "openline",
				"appSource":           "manual",
				"status":              "ACTIVE",
				"initiatorFirstName":  "",
				"initiatorLastName":   "",
				"initiatorUsername":   initiatorUsername,
				"initiatorType":       initiatorTypeStr,
				"startedAt":           utc,
				"updatedAt":           utc,
				"threadId":            threadId,
				"channel":             channel,
				"contactIds":          contactIds,
				"userIds":             userIds,
				"lastSenderId":        "",
				"lastSenderType":      "",
				"lastSenderFirstName": "",
				"lastSenderLastName":  "",
				"lastContentPreview":  "",
			})

		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		dbNode := result.(*dbtype.Node)
		return mapNodeToConversation(dbNode), nil
	}
}

func (s *CustomerOSService) UpdateConversation(ctx context.Context, tenant string, conversationId string, lastSenderId string, lastSenderType string, contactIds []string, userIds []string, lastContentPreview string) (string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:$conversationId}) " +
			" ON MATCH SET " +
			" o.messageCount=o.messageCount+1, o.updatedAt=$updatedAt, " +
			" o.lastSenderId=$lastSenderId, o.lastSenderType=$lastSenderType, o.lastSenderFirstName=$lastSenderFirstName, o.lastSenderLastName=$lastSenderLastName, o.lastContentPreview=$lastContentPreview, " +
			" o:%s " +
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
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Conversation_"+tenant, queryLinkWithContacts, queryLinkWithUsers),
			map[string]interface{}{
				"tenant":              tenant,
				"conversationId":      conversationId,
				"updatedAt":           time.Now().UTC(),
				"contactIds":          contactIds,
				"userIds":             userIds,
				"lastSenderId":        lastSenderId,
				"lastSenderType":      lastSenderType,
				"lastSenderFirstName": "",
				"lastSenderLastName":  "",
				"lastContentPreview":  lastContentPreview,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return "", err
	} else {
		dbNode := result.(*dbtype.Node)
		return utils.GetPropsFromNode(*dbNode)["id"].(string), err
	}
}

func mapNodeToConversation(node *dbtype.Node) *Conversation {
	if node == nil {
		return nil
	}

	props := utils.GetPropsFromNode(*node)

	conversation := new(Conversation)
	conversation.Id = utils.GetPropsFromNode(*node)["id"].(string)
	conversation.Status = utils.GetPropsFromNode(*node)["status"].(string)
	conversation.Channel = utils.GetPropsFromNode(*node)["channel"].(string)
	conversation.StartedAt = utils.GetPropsFromNode(*node)["startedAt"].(time.Time)
	conversation.UpdatedAt = utils.GetPropsFromNode(*node)["updatedAt"].(time.Time)
	conversation.InitiatorFirstName = utils.GetPropsFromNode(*node)["initiatorFirstName"].(string)
	conversation.InitiatorLastName = utils.GetPropsFromNode(*node)["initiatorLastName"].(string)
	conversation.InitiatorUsername = utils.GetPropsFromNode(*node)["initiatorUsername"].(string)
	conversation.InitiatorType = utils.GetPropsFromNode(*node)["initiatorType"].(string)

	conversation.LastSenderId = utils.GetStringPropOrEmpty(props, "lastSenderId")
	conversation.LastSenderType = utils.GetStringPropOrEmpty(props, "lastSenderType")
	conversation.LastSenderFirstName = utils.GetStringPropOrEmpty(props, "lastSenderFirstName")
	conversation.LastSenderLastName = utils.GetStringPropOrEmpty(props, "lastSenderLastName")
	conversation.LastContentPreview = utils.GetStringPropOrEmpty(props, "lastContentPreview")

	return conversation
}

func (s *CustomerOSService) GetContactWithEmailOrCreate(ctx context.Context, tenant string, email string) (Contact, error) {
	contact, err := s.GetContactByEmail(ctx, email)
	if err != nil {
		contact, err = s.CreateContactWithEmail(ctx, tenant, email)
		if err != nil {
			return Contact{}, err
		}
		if contact == nil {
			return Contact{}, errors.New("contact not found and could not be created")
		}
		return *contact, nil
	} else {
		return *contact, nil
	}
}

func (s *CustomerOSService) GetActiveConversationOrCreate(
	ctx context.Context,
	tenant string,
	initiator Participant,
	initiatorUsername string,
	eventType entity.EventType,
	threadId string,
) (*Conversation, error) {
	var conversation *Conversation
	var err error
	if eventType == entity.WEB_CHAT {
		if initiator.Type == entity.CONTACT {
			conversation, err = s.GetWebChatConversationWithContactInitiator(ctx, tenant, initiator.Id)
		} else if initiator.Type == entity.USER {
			conversation, err = s.GetWebChatConversationWithUserInitiator(ctx, tenant, initiator.Id)
		}
	} else if eventType == entity.EMAIL {
		if err != nil {
			return nil, err
		}
		if initiator.Type == entity.CONTACT {
			conversation, err = s.GetEmailConversationWithContactInitiator(ctx, tenant, initiator.Id, threadId)
		} else if initiator.Type == entity.USER {
			conversation, err = s.GetEmailConversationWithUserInitiator(ctx, tenant, initiator.Id, threadId)
		}
	}

	if err != nil {
		return nil, err
	}

	if conversation == nil {
		conversation, err = s.CreateConversation(ctx, tenant, initiator, initiatorUsername, eventType, threadId)
	}
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func (s *CustomerOSService) GetEmailConversationWithContactInitiator(ctx context.Context, tenant string, contactId string, threadId string) (*Conversation, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.driver)
	defer session.Close(ctx)

	conversationNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `match(o:Conversation{status:"ACTIVE", channel: "EMAIL", threadId: $threadId})<-[:INITIATED]-(c:Contact{id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant{name:$tenant}) return o`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"threadId":  threadId,
			}); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single(ctx)
			if err != nil && err.Error() != "Result contains no more records" {
				return nil, err
			}
			if record != nil {
				return record.Values[0].(dbtype.Node), nil
			} else {
				return nil, nil
			}
		}
	})

	if err != nil {
		return nil, err
	}

	if conversationNode != nil {
		node := conversationNode.(dbtype.Node)
		return mapNodeToConversation(&node), nil
	} else {
		return nil, nil
	}
}

func (s *CustomerOSService) GetEmailConversationWithUserInitiator(ctx context.Context, tenant string, userId string, threadId string) (*Conversation, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.driver)
	defer session.Close(ctx)

	conversationNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `match(o:Conversation{status:"ACTIVE", channel: "EMAIL", threadId: $threadId})<-[:INITIATED]-(u:User{id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant{name:$tenant}) return o`,
			map[string]any{
				"tenant":   tenant,
				"userId":   userId,
				"threadId": threadId,
			}); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single(ctx)
			if err != nil && err.Error() != "Result contains no more records" {
				return nil, err
			}
			if record != nil {
				return record.Values[0].(dbtype.Node), nil
			} else {
				return nil, nil
			}
		}
	})

	if err != nil {
		return nil, err
	}

	if conversationNode != nil {
		node := conversationNode.(dbtype.Node)
		return mapNodeToConversation(&node), nil
	} else {
		return nil, nil
	}
}

func (s *CustomerOSService) GetWebChatConversationWithContactInitiator(ctx context.Context, tenant string, contactId string) (*Conversation, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.driver)
	defer session.Close(ctx)

	conversationNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `match(o:Conversation{status:"ACTIVE", channel: "WEB_CHAT"})<-[:INITIATED]-(c:Contact{id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant{name:$tenant}) return o`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			}); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single(ctx)
			if err != nil && err.Error() != "Result contains no more records" {
				return nil, err
			}
			if record != nil {
				return record.Values[0].(dbtype.Node), nil
			} else {
				return nil, nil
			}
		}
	})

	if err != nil {
		return nil, err
	}

	if conversationNode != nil {
		node := conversationNode.(dbtype.Node)
		return mapNodeToConversation(&node), nil
	} else {
		return nil, nil
	}
}

func (s *CustomerOSService) GetWebChatConversationWithUserInitiator(ctx context.Context, tenant string, userId string) (*Conversation, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.driver)
	defer session.Close(ctx)

	conversationNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `match(o:Conversation{status:"ACTIVE", channel: "WEB_CHAT"})<-[:INITIATED]-(u:User{id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant{name:$tenant}) return o`,
			map[string]any{
				"tenant": tenant,
				"userId": userId,
			}); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single(ctx)
			if err != nil && err.Error() != "Result contains no more records" {
				return nil, err
			}
			if record != nil {
				return record.Values[0].(dbtype.Node), nil
			} else {
				return nil, nil
			}
		}
	})

	if err != nil {
		return nil, err
	}

	if conversationNode != nil {
		node := conversationNode.(dbtype.Node)
		return mapNodeToConversation(&node), nil
	} else {
		return nil, nil
	}
}

func NewCustomerOSService(driver *neo4j.DriverWithContext, graphqlClient *graphql.Client, postgresRepositories *repository.PostgresRepositories, commonStoreService *commonStoreService, config *config.Config) *CustomerOSService {
	customerOsService := new(CustomerOSService)
	customerOsService.driver = driver
	customerOsService.postgresRepositories = postgresRepositories
	customerOsService.commonStoreService = commonStoreService
	customerOsService.graphqlClient = graphqlClient
	customerOsService.conf = config
	return customerOsService
}
