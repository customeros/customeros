package generate

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"io/ioutil"
	"net/http"
)

func AddDemoTenantRoutes(rg *gin.RouterGroup, config *config.Config, cosClient service.CustomerOsClient) {
	rg.GET("/demo-tenant", func(context *gin.Context) {

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

		appSource := "user-admin-api"

		tenant := context.GetHeader("TENANT_NAME")
		if tenant == "" {
			context.JSON(http.StatusBadRequest, gin.H{
				"result": fmt.Sprintf("tenant is required"),
			})
			return
		}

		username := context.GetHeader("MASTER_USERNAME")
		if username == "" {
			context.JSON(http.StatusBadRequest, gin.H{
				"result": fmt.Sprintf("username is required"),
			})
			return
		}

		multipartFileHeader, err := context.FormFile("file")
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"result": fmt.Sprintf("file is required"),
			})
			return
		}

		multipartFile, err := multipartFileHeader.Open()
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"result": fmt.Sprintf("file is required"),
			})
			return
		}

		bytes, err := ioutil.ReadAll(multipartFile)
		if err != nil {
			panic(err)
		}

		// Parse the JSON file into the User struct
		var sourceData SourceData
		if err := json.Unmarshal(bytes, &sourceData); err != nil {
			panic(err)
		}

		var userIds = make([]EmailAddressWithId, len(sourceData.Users))
		var contactIds = make([]EmailAddressWithId, len(sourceData.Contacts))

		//users creation
		for _, user := range sourceData.Users {
			userId, err := cosClient.CreateUser(&cosModel.UserInput{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email: cosModel.EmailInput{
					Email: user.Email,
				},
				AppSource: &appSource,
			}, tenant, []service.Role{service.ROLE_USER})
			if err != nil {
				context.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			userIds = append(userIds, EmailAddressWithId{
				Email: user.Email,
				Id:    userId,
			})
		}

		//contacts creation
		for _, contact := range sourceData.Contacts {
			contactId, err := cosClient.CreateContact(tenant, username, contact.FirstName, contact.LastName, contact.Email)
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
			organizationId, err := cosClient.CreateOrganization(tenant, username, organization.Name, organization.Domain)
			if err != nil {
				context.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
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

				err = cosClient.AddOrganizationToContact(tenant, username, contactId, organizationId)
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

				sessionId, err := cosClient.CreateInteractionSession(tenant, username, sessionOpts...)
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

				interactionEventId, err := cosClient.CreateInteractionEvent(tenant, username, eventOpts...)
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
				meetingId, err := cosClient.CreateMeeting(tenant, username, input)
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

				interactionEventId, err := cosClient.CreateInteractionEvent(tenant, username, eventOpts...)
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
		context.JSON(200, gin.H{
			"tenant": "tenant initiated",
		})
	})
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
