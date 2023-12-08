package generate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"net/http"
)

func AddDemoTenantRoutes(rg *gin.RouterGroup, config *config.Config, services *service.Services) {
	appSource := "user-admin-api"

	rg.GET("/demo-tenant-users", func(context *gin.Context) {
		apiKey := context.GetHeader("X-Openline-Api-Key")
		if apiKey != config.Service.ApiKey {
			context.JSON(http.StatusUnauthorized, gin.H{
				"result": fmt.Sprintf("invalid api key"),
			})
			return
		}

		sourceData, err := validateRequestAndGetFileBytes(context)
		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		tenant := context.GetHeader("TENANT_NAME")

		//users creation
		for _, user := range sourceData.Users {
			_, err := services.CustomerOsClient.CreateUser(&cosModel.UserInput{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email: cosModel.EmailInput{
					Email: user.Email,
				},
				AppSource:       &appSource,
				ProfilePhotoURL: user.ProfilePhotoURL,
			}, tenant, []cosModel.Role{cosModel.RoleUser, cosModel.RoleOwner})
			if err != nil {
				context.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		context.JSON(200, gin.H{
			"OK": "users initiated",
		})
	})

	rg.GET("/demo-tenant-data", func(context *gin.Context) {

		apiKey := context.GetHeader("X-Openline-Api-Key")
		if apiKey != config.Service.ApiKey {
			context.JSON(http.StatusUnauthorized, gin.H{
				"result": fmt.Sprintf("invalid api key"),
			})
			return
		}

		//match (n:User_LightBlok)--(e:Email) where e.email <> "customerosdemo@gmail.com" detach delete n;
		//match (n:Email_LightBlok) where n.email <> "customerosdemo@gmail.com" detach delete n;
		//match (n:Contact_LightBlok) detach delete n;
		//match (n:JobRole_LightBlok) detach delete n;
		//match (n:Organization_LightBlok) detach delete n;
		//match (n:InteractionSession_LightBlok) detach delete n;
		//match (n:InteractionEvent_LightBlok) detach delete n;
		//match (n:Note_LightBlok) detach delete n;
		//match (n:Action_LightBlok) detach delete n;
		//match (n:Meeting_LightBlok) detach delete n;
		//match (n:Issue_LightBlok) detach delete n;
		//match (n:LogEntry_LightBlok) detach delete n;

		sourceData, err := validateRequestAndGetFileBytes(context)
		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		tenant := context.GetHeader("TENANT_NAME")
		username := context.GetHeader("MASTER_USERNAME")

		var userIds = make([]EmailAddressWithId, len(sourceData.Users))
		var contactIds = make([]EmailAddressWithId, len(sourceData.Contacts))

		//read users
		for _, user := range sourceData.Users {
			userResponse, err := services.CustomerOsClient.GetUserByEmail(tenant, user.Email)
			if err != nil {
				context.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			userIds = append(userIds, EmailAddressWithId{
				Email: user.Email,
				Id:    userResponse.ID,
			})
		}

		//contacts creation
		for _, contact := range sourceData.Contacts {
			contactId, err := services.CustomerOsClient.CreateContact(tenant, username, contact.FirstName, contact.LastName, contact.Email, contact.ProfilePhotoURL)
			if err != nil {
				context.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			contactIds = append(contactIds, EmailAddressWithId{
				Email: contact.Email,
				Id:    contactId,
			})
		}

		//create orgs
		for _, organization := range sourceData.Organizations {

			var organizationId string
			if organization.Id != "" {
				organizationId = organization.Id
			} else {
				organizationId, err = services.CustomerOsClient.CreateOrganization(tenant, username, organization.Name, organization.Domain)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}
			}

			//create Contracts with Service Lines in org
			for _, contract := range organization.Contracts {
				contractInput := cosModel.ContractInput{
					OrganizationId:   organizationId,
					Name:             contract.Name,
					RenewalCycle:     contract.RenewalCycle,
					RenewalPeriods:   contract.RenewalPeriods,
					ContractUrl:      contract.ContractUrl,
					ServiceStartedAt: contract.ServiceStartedAt,
					SignedAt:         contract.SignedAt,
				}
				contractId, err := services.CustomerOsClient.CreateContract(tenant, username, contractInput)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}

				if contractId == "" {
					context.JSON(500, gin.H{
						"error": "contractId is nil",
					})
					return
				}

				for _, serviceLine := range contract.ServiceLines {

					serviceLineInput := func() interface{} {
						if serviceLine.EndedAt.IsZero() {
							return cosModel.ServiceLineInput{
								ContractId: contractId,
								Name:       serviceLine.Name,
								Billed:     serviceLine.Billed,
								Price:      serviceLine.Price,
								Quantity:   serviceLine.Quantity,
								StartedAt:  serviceLine.StartedAt,
							}
						}
						return cosModel.ServiceLineEndedInput{
							ContractId: contractId,
							Name:       serviceLine.Name,
							Billed:     serviceLine.Billed,
							Price:      serviceLine.Price,
							Quantity:   serviceLine.Quantity,
							StartedAt:  serviceLine.StartedAt,
							EndedAt:    serviceLine.EndedAt,
						}
					}()
					serviceLineId, err := services.CustomerOsClient.CreateServiceLine(tenant, username, serviceLineInput)
					if err != nil {
						context.JSON(500, gin.H{
							"error": err.Error(),
						})
						return
					}

					if serviceLineId == "" {
						context.JSON(500, gin.H{
							"error": "serviceLineId is nil",
						})
						return
					}
				}
			}

			//create people in org
			for _, people := range organization.People {
				var contactId string
				for _, contact := range contactIds {
					if contact.Email == people.Email {
						contactId = contact.Id
						break
					}
				}

				if contactId == "" {
					context.JSON(500, gin.H{
						"error": fmt.Sprintf("contact not found for email %s", people.Email),
					})
					return
				}

				err = services.CustomerOsClient.AddContactToOrganization(tenant, username, contactId, organizationId, people.JobRole, people.Description)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}
			}

			//create emails
			for _, email := range organization.Emails {
				sig, err := uuid.NewUUID()
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}
				sigs := sig.String()

				channelValue := "EMAIL"
				appSource := appSource
				sessionStatus := "ACTIVE"
				sessionType := "THREAD"
				sessionOpts := []service.InteractionSessionBuilderOption{
					service.WithSessionIdentifier(&sigs),
					service.WithSessionChannel(&channelValue),
					service.WithSessionName(&email.Subject),
					service.WithSessionAppSource(&appSource),
					service.WithSessionStatus(&sessionStatus),
					service.WithSessionType(&sessionType),
				}

				sessionId, err := services.CustomerOsClient.CreateInteractionSession(tenant, username, sessionOpts...)
				if sessionId == nil {
					context.JSON(500, gin.H{
						"error": "sessionId is nil",
					})
					return
				}

				participantTypeTO, participantTypeCC, participantTypeBCC := "TO", "CC", "BCC"
				participantsTO := toParticipantInputArr(email.To, &participantTypeTO)
				participantsCC := toParticipantInputArr(email.Cc, &participantTypeCC)
				participantsBCC := toParticipantInputArr(email.Bcc, &participantTypeBCC)
				sentTo := append(append(participantsTO, participantsCC...), participantsBCC...)
				sentBy := toParticipantInputArr([]string{email.From}, nil)

				emailChannelData, err := buildEmailChannelData(email.Subject, err)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}

				iig, err := uuid.NewUUID()
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}
				iigs := iig.String()
				eventOpts := []service.InteractionEventBuilderOption{
					service.WithCreatedAt(&email.Date),
					service.WithSessionId(sessionId),
					service.WithEventIdentifier(iigs),
					service.WithChannel(&channelValue),
					service.WithChannelData(emailChannelData),
					service.WithContent(&email.Body),
					service.WithContentType(&email.ContentType),
					service.WithSentBy(sentBy),
					service.WithSentTo(sentTo),
					service.WithAppSource(&appSource),
				}

				interactionEventId, err := services.CustomerOsClient.CreateInteractionEvent(tenant, username, eventOpts...)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}

				if interactionEventId == nil {
					context.JSON(500, gin.H{
						"error": "interactionEventId is nil",
					})
					return
				}
			}

			//create meetings
			for _, meeting := range organization.Meetings {
				var createdBy []*cosModel.MeetingParticipantInput
				createdBy = append(createdBy, getMeetingParticipantInput(meeting.CreatedBy, userIds, contactIds))

				var attendedBy []*cosModel.MeetingParticipantInput
				for _, attendee := range meeting.Attendees {
					attendedBy = append(attendedBy, getMeetingParticipantInput(attendee, userIds, contactIds))
				}

				contentType := "text/plain"
				noteInput := cosModel.NoteInput{Content: &meeting.Agenda, ContentType: &contentType, AppSource: &appSource}
				input := cosModel.MeetingInput{
					Name:       &meeting.Subject,
					CreatedAt:  &meeting.StartedAt,
					CreatedBy:  createdBy,
					AttendedBy: attendedBy,
					StartedAt:  &meeting.StartedAt,
					EndedAt:    &meeting.EndedAt,
					Note:       &noteInput,
					AppSource:  &appSource,
				}
				meetingId, err := services.CustomerOsClient.CreateMeeting(tenant, username, input)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}

				if meetingId == "" {
					context.JSON(500, gin.H{
						"error": "meetingId is nil",
					})
					return
				}

				eventType := "meeting"
				eventOpts := []service.InteractionEventBuilderOption{
					service.WithSentBy([]cosModel.InteractionEventParticipantInput{*getInteractionEventParticipantInput(meeting.CreatedBy, userIds, contactIds)}),
					service.WithSentTo(getInteractionEventParticipantInputList(meeting.Attendees, userIds, contactIds)),
					service.WithMeetingId(&meetingId),
					service.WithCreatedAt(&meeting.StartedAt),
					service.WithEventType(&eventType),
					service.WithAppSource(&appSource),
				}

				interactionEventId, err := services.CustomerOsClient.CreateInteractionEvent(tenant, username, eventOpts...)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}

				if interactionEventId == nil {
					context.JSON(500, gin.H{
						"error": "interactionEventId is nil",
					})
					return
				}
			}

			//log entries
			for _, logEntry := range organization.LogEntries {

				interactionEventId, err := services.CustomerOsClient.CreateLogEntry(tenant, username, organizationId, logEntry.CreatedBy, logEntry.Content, logEntry.ContentType, logEntry.Date)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}

				if interactionEventId == nil {
					context.JSON(500, gin.H{
						"error": "interactionEventId is nil",
					})
					return
				}
			}

			//issues
			for index, issue := range organization.Issues {
				issueGrpcRequest := issuepb.UpsertIssueGrpcRequest{
					Tenant:      tenant,
					Subject:     issue.Subject,
					Status:      issue.Status,
					Priority:    issue.Priority,
					Description: issue.Description,
					CreatedAt:   timestamppb.New(issue.CreatedAt),
					UpdatedAt:   timestamppb.New(issue.CreatedAt),
					SourceFields: &commonpb.SourceFields{
						Source:    "zendesk_support",
						AppSource: appSource,
					},
					ExternalSystemFields: &commonpb.ExternalSystemFields{
						ExternalSystemId: "zendesk_support",
						ExternalId:       "random-thing-" + fmt.Sprintf("%d", index),
						ExternalUrl:      "https://random-thing.zendesk.com/agent/tickets/" + fmt.Sprintf("%d", index),
						SyncDate:         timestamppb.New(issue.CreatedAt),
					},
				}

				issueGrpcRequest.ReportedByOrganizationId = &organizationId

				for _, userWithId := range userIds {
					if userWithId.Email == issue.CreatedBy {
						issueGrpcRequest.SubmittedByUserId = &userWithId.Id
						break
					}
				}

				_, err = services.GrpcClients.IssueClient.UpsertIssue(context, &issueGrpcRequest)
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}
			}

			//slack
			for _, slackThread := range organization.Slack {

				sig, err := uuid.NewUUID()
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}
				sigs := sig.String()

				channelValue := "CHAT"
				appSource := appSource
				sessionStatus := "ACTIVE"
				sessionType := "THREAD"
				sessionName := slackThread[0].Message
				sessionOpts := []service.InteractionSessionBuilderOption{
					service.WithSessionIdentifier(&sigs),
					service.WithSessionChannel(&channelValue),
					service.WithSessionName(&sessionName),
					service.WithSessionAppSource(&appSource),
					service.WithSessionStatus(&sessionStatus),
					service.WithSessionType(&sessionType),
				}

				sessionId, err := services.CustomerOsClient.CreateInteractionSession(tenant, username, sessionOpts...)
				if sessionId == nil {
					context.JSON(500, gin.H{
						"error": "sessionId is nil",
					})
					return
				}

				for _, slackMessage := range slackThread {

					sentBy := toParticipantInputArr([]string{slackMessage.CreatedBy}, nil)

					iig, err := uuid.NewUUID()
					if err != nil {
						context.JSON(500, gin.H{
							"error": err.Error(),
						})
						return
					}
					iigs := iig.String()
					eventType := "MESSAGE"
					contentType := "text/plain"
					eventOpts := []service.InteractionEventBuilderOption{
						service.WithCreatedAt(&slackMessage.CreatedAt),
						service.WithSessionId(sessionId),
						service.WithEventIdentifier(iigs),
						service.WithExternalId(iigs),
						service.WithExternalSystemId("slack"),
						service.WithChannel(&channelValue),
						service.WithEventType(&eventType),
						service.WithContent(&slackMessage.Message),
						service.WithContentType(&contentType),
						service.WithSentBy(sentBy),
						service.WithAppSource(&appSource),
					}

					interactionEventId, err := services.CustomerOsClient.CreateInteractionEvent(tenant, username, eventOpts...)
					if err != nil {
						context.JSON(500, gin.H{
							"error": err.Error(),
						})
						return
					}

					if interactionEventId == nil {
						context.JSON(500, gin.H{
							"error": "interactionEventId is nil",
						})
						return
					}

				}

			}

			//intercom
			for _, intercomThread := range organization.Intercom {

				sig, err := uuid.NewUUID()
				if err != nil {
					context.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}
				sigs := sig.String()

				channelValue := "CHAT"
				appSource := appSource
				sessionStatus := "ACTIVE"
				sessionType := "THREAD"
				sessionName := intercomThread[0].Message
				sessionOpts := []service.InteractionSessionBuilderOption{
					service.WithSessionIdentifier(&sigs),
					service.WithSessionChannel(&channelValue),
					service.WithSessionName(&sessionName),
					service.WithSessionAppSource(&appSource),
					service.WithSessionStatus(&sessionStatus),
					service.WithSessionType(&sessionType),
				}

				sessionId, err := services.CustomerOsClient.CreateInteractionSession(tenant, username, sessionOpts...)
				if sessionId == nil {
					context.JSON(500, gin.H{
						"error": "sessionId is nil",
					})
					return
				}

				for _, intercomMessage := range intercomThread {

					sentById := ""
					for _, contactWithId := range contactIds {
						if contactWithId.Email == intercomMessage.CreatedBy {
							sentById = contactWithId.Id
							break
						}
					}
					sentBy := toContactParticipantInputArr([]string{sentById})

					iig, err := uuid.NewUUID()
					if err != nil {
						context.JSON(500, gin.H{
							"error": err.Error(),
						})
						return
					}
					iigs := iig.String()
					eventType := "MESSAGE"
					contentType := "text/html"
					eventOpts := []service.InteractionEventBuilderOption{
						service.WithCreatedAt(&intercomMessage.CreatedAt),
						service.WithSessionId(sessionId),
						service.WithEventIdentifier(iigs),
						service.WithExternalId(iigs),
						service.WithExternalSystemId("intercom"),
						service.WithChannel(&channelValue),
						service.WithEventType(&eventType),
						service.WithContent(&intercomMessage.Message),
						service.WithContentType(&contentType),
						service.WithSentBy(sentBy),
						service.WithAppSource(&appSource),
					}

					interactionEventId, err := services.CustomerOsClient.CreateInteractionEvent(tenant, username, eventOpts...)
					if err != nil {
						context.JSON(500, gin.H{
							"error": err.Error(),
						})
						return
					}

					if interactionEventId == nil {
						context.JSON(500, gin.H{
							"error": "interactionEventId is nil",
						})
						return
					}

				}

			}

		}
		context.JSON(200, gin.H{
			"tenant": "tenant initiated",
		})
	})
}

func validateRequestAndGetFileBytes(context *gin.Context) (*SourceData, error) {
	tenant := context.GetHeader("TENANT_NAME")
	if tenant == "" {
		return nil, errors.New("tenant is required")
	}

	username := context.GetHeader("MASTER_USERNAME")
	if username == "" {
		return nil, errors.New("username is required")
	}

	multipartFileHeader, err := context.FormFile("file")
	if err != nil {
		return nil, err
	}

	multipartFile, err := multipartFileHeader.Open()
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(multipartFile)
	if err != nil {
		return nil, err
	}

	var sourceData SourceData
	if err := json.Unmarshal(bytes, &sourceData); err != nil {
		return nil, err
	}

	return &sourceData, nil
}

func toParticipantInputArr(from []string, participantType *string) []cosModel.InteractionEventParticipantInput {
	var to []cosModel.InteractionEventParticipantInput
	for _, a := range from {
		participantInput := cosModel.InteractionEventParticipantInput{
			Email: &a,
			Type:  participantType,
		}
		to = append(to, participantInput)
	}
	return to
}

func toContactParticipantInputArr(from []string) []cosModel.InteractionEventParticipantInput {
	var to []cosModel.InteractionEventParticipantInput
	for _, a := range from {
		participantInput := cosModel.InteractionEventParticipantInput{
			ContactID: &a,
		}
		to = append(to, participantInput)
	}
	return to
}

func buildEmailChannelData(subject string, err error) (*string, error) {
	emailContent := cosModel.EmailChannelData{
		Subject: subject,
		//InReplyTo: utils.EnsureEmailRfcIds(email.InReplyTo),
		//Reference: utils.EnsureEmailRfcIds(email.References),
	}
	jsonContent, err := json.Marshal(emailContent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email content: %v", err)
	}
	jsonContentString := string(jsonContent)

	return &jsonContentString, nil
}

func getMeetingParticipantInput(emailAddress string, userIds, contactIds []EmailAddressWithId) *cosModel.MeetingParticipantInput {
	for _, userWithId := range userIds {
		if userWithId.Email == emailAddress {
			return &cosModel.MeetingParticipantInput{UserID: &userWithId.Id}
		}
	}

	for _, contactWithId := range contactIds {
		if contactWithId.Email == emailAddress {
			return &cosModel.MeetingParticipantInput{ContactID: &contactWithId.Id}
		}
	}

	return nil
}

func getInteractionEventParticipantInput(emailAddress string, userIds, contactIds []EmailAddressWithId) *cosModel.InteractionEventParticipantInput {
	for _, userWithId := range userIds {
		if userWithId.Email == emailAddress {
			return &cosModel.InteractionEventParticipantInput{UserID: &userWithId.Id}
		}
	}

	for _, contactWithId := range contactIds {
		if contactWithId.Email == emailAddress {
			return &cosModel.InteractionEventParticipantInput{ContactID: &contactWithId.Id}
		}
	}

	return nil
}

func getInteractionEventParticipantInputList(emailAddresses []string, userIds, contactIds []EmailAddressWithId) []cosModel.InteractionEventParticipantInput {
	var interactionEventParticipantInputList []cosModel.InteractionEventParticipantInput
	for _, emailAddress := range emailAddresses {
		interactionEventParticipantInputList = append(interactionEventParticipantInputList, *getInteractionEventParticipantInput(emailAddress, userIds, contactIds))
	}
	return interactionEventParticipantInputList
}

type EmailAddressWithId struct {
	Email string
	Id    string
}
