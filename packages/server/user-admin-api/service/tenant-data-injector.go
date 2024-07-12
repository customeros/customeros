package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	issuepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/issue"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type SourceData struct {
	Users []struct {
		FirstName       string  `json:"firstName"`
		LastName        string  `json:"lastName"`
		Email           string  `json:"email"`
		ProfilePhotoURL *string `json:"profilePhotoUrl,omitempty"`
	} `json:"users"`
	Contacts []struct {
		FirstName       string  `json:"firstName"`
		LastName        string  `json:"lastName"`
		Email           string  `json:"email"`
		ProfilePhotoURL *string `json:"profilePhotoUrl,omitempty"`
		Timezone        *string `json:"timezone,omitempty"`
		PhoneNumber     string  `json:"phoneNumber,omitempty"`
		Social          *string `json:"social,omitempty"`
		Description     string  `json:"description,omitempty"`
		Note            *string `json:"note,omitempty"`
	} `json:"contacts"`
	TenantBillingProfiles []struct {
		LegalName                     string `json:"legalName"`
		Email                         string `json:"email"`
		AddressLine1                  string `json:"addressLine1"`
		Locality                      string `json:"locality"`
		Country                       string `json:"country"`
		Zip                           string `json:"zip"`
		DomesticPaymentsBankInfo      string `json:"domesticPaymentsBankInfo"`
		InternationalPaymentsBankInfo string `json:"internationalPaymentsBankInfo"`
		VatNumber                     string `json:"vatNumber"`
		SendInvoicesFrom              string `json:"sendInvoicesFrom"`
		CanPayWithCard                bool   `json:"canPayWithCard"`
		CanPayWithDirectDebitSEPA     bool   `json:"canPayWithDirectDebitSEPA"`
		CanPayWithDirectDebitACH      bool   `json:"canPayWithDirectDebitACH"`
		CanPayWithDirectDebitBacs     bool   `json:"canPayWithDirectDebitBacs"`
		CanPayWithPigeon              bool   `json:"canPayWithPigeon"`
		CanPayWithBankTransfer        bool   `json:"canPayWithBankTransfer"`
		Check                         bool   `json:"check"`
	} `json:"tenantBillingProfiles"`
	Organizations []struct {
		Id                    string  `json:"id"`
		Name                  string  `json:"name"`
		ValueProposition      string  `json:"valueProposition"`
		Website               string  `json:"website"`
		Logo                  string  `json:"logo"`
		Domain                string  `json:"domain"`
		Notes                 string  `json:"notes"`
		Industry              string  `json:"industry"`
		LastFundingRound      string  `json:"lastFundingRound"`
		TargetAudience        string  `json:"targetAudience"`
		Market                string  `json:"market"`
		Employees             int64   `json:"employees"`
		Relationship          string  `json:"relationship"`
		LastFundingAmount     string  `json:"lastFundingAmount"`
		OrganizationSocial    *string `json:"organizationSocial,omitempty"`
		OnboardingStatusInput []struct {
			Status   string `json:"status"`
			Comments string `json:"comments"`
		} `json:"onboardingStatusInput"`
		Contracts []struct {
			ContractName            string     `json:"contractName"`
			CommittedPeriodInMonths int64      `json:"committedPeriodInMonths"`
			ContractUrl             string     `json:"contractUrl"`
			ServiceStarted          time.Time  `json:"serviceStarted"`
			ContractSigned          time.Time  `json:"contractSigned"`
			InvoicingStartDate      *time.Time `json:"invoicingStartDate"`
			BillingCycle            string     `json:"billingCycle"`
			Currency                string     `json:"currency"`
			AddressLine1            string     `json:"addressLine1"`
			AddressLine2            string     `json:"addressLine2"`
			Zip                     string     `json:"zip"`
			Locality                string     `json:"locality"`
			Country                 string     `json:"country"`
			OrganizationLegalName   string     `json:"organizationLegalName"`
			InvoiceEmail            string     `json:"invoiceEmail"`
			InvoiceNote             string     `json:"invoiceNote"`
			Approved                bool       `json:"approved"`
			ServiceLines            []struct {
				Description    string     `json:"description"`
				BillingCycle   string     `json:"billingCycle"`
				Price          int        `json:"price"`
				Quantity       int        `json:"quantity"`
				ServiceStarted *time.Time `json:"serviceStarted"`
				ServiceEnded   *time.Time `json:"serviceEnded,omitempty"`
			} `json:"serviceLines"`
		} `json:"contracts,omitempty"`
		People []struct {
			Email       string `json:"email"`
			JobRole     string `json:"jobRole"`
			Description string `json:"description"`
		} `json:"people"`
		Emails []struct {
			From        string    `json:"from"`
			To          []string  `json:"to"`
			Cc          []string  `json:"cc"`
			Bcc         []string  `json:"bcc"`
			Subject     string    `json:"subject"`
			Body        string    `json:"body"`
			ContentType string    `json:"contentType"`
			Date        time.Time `json:"date"`
		} `json:"emails"`
		Meetings []struct {
			CreatedBy string    `json:"createdBy"`
			Attendees []string  `json:"attendees"`
			Subject   string    `json:"subject"`
			Agenda    string    `json:"agenda"`
			StartedAt time.Time `json:"startedAt"`
			EndedAt   time.Time `json:"endedAt"`
		} `json:"meetings"`
		LogEntries []struct {
			CreatedBy   string    `json:"createdBy"`
			Content     string    `json:"content"`
			ContentType string    `json:"contentType"`
			Date        time.Time `json:"date"`
		} `json:"logEntries"`
		Issues []struct {
			CreatedBy   string    `json:"createdBy"`
			CreatedAt   time.Time `json:"createdAt"`
			Subject     string    `json:"subject"`
			Status      string    `json:"status"`
			Priority    string    `json:"priority"`
			Description string    `json:"description"`
		} `json:"issues"`
		Slack [][]struct {
			CreatedBy string    `json:"createdBy"`
			CreatedAt time.Time `json:"createdAt"`
			Message   string    `json:"message"`
		} `json:"slack"`
		Intercom [][]struct {
			CreatedBy string    `json:"createdBy"`
			CreatedAt time.Time `json:"createdAt"`
			Message   string    `json:"message"`
		} `json:"intercom"`
	} `json:"organizations"`
	MasterPlans []struct {
		Name       string `json:"name"`
		Milestones []struct {
			Name          string   `json:"name"`
			Order         int64    `json:"order"`
			DurationHours int64    `json:"durationHours"`
			Optional      bool     `json:"optional"`
			Items         []string `json:"items"`
		} `json:"milestones"`
	} `json:"masterPlans"`
}

type TenantDataInjector interface {
	InjectTenantData(context context.Context, tenant, username string, sourceData *SourceData) error
	CleanupTenantData(context context.Context, tenant, username, reqTenant, reqConfirmTenant string) error
}

type tenantDataInjector struct {
	services *Services
}

func NewTenantDataInjector(services *Services) TenantDataInjector {
	return &tenantDataInjector{
		services: services,
	}
}

func (t *tenantDataInjector) InjectTenantData(ctx context.Context, tenant, username string, sourceData *SourceData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantDataInjector.InjectTenantData")
	defer span.Finish()

	appSource := "user-admin-api"

	var userIds = make([]EmailAddressWithId, len(sourceData.Users))
	var contactIds = make([]EmailAddressWithId, len(sourceData.Contacts))

	//users creation
	for _, user := range sourceData.Users {
		userResponse, err := t.services.CustomerOsClient.GetUserByEmail(tenant, user.Email)
		if err != nil {
			return err
		}
		if userResponse == nil {
			userResponse, err := t.services.CustomerOsClient.CreateUser(&cosModel.UserInput{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email: cosModel.EmailInput{
					Email: user.Email,
				},
				AppSource:       &appSource,
				ProfilePhotoURL: user.ProfilePhotoURL,
			}, tenant, []cosModel.Role{cosModel.RoleUser, cosModel.RoleOwner})
			if err != nil {
				return err
			}

			userIds = append(userIds, EmailAddressWithId{
				Email: user.Email,
				Id:    userResponse.ID,
			})
		} else {
			userIds = append(userIds, EmailAddressWithId{
				Email: user.Email,
				Id:    userResponse.ID,
			})
		}
	}

	//contacts creation
	for _, contact := range sourceData.Contacts {
		contactInput := model.ContactInput{
			FirstName: &contact.FirstName,
			LastName:  &contact.LastName,
			Email: &model.EmailInput{
				Email: contact.Email,
			},
			ProfilePhotoURL: contact.ProfilePhotoURL,
			Timezone:        contact.Timezone,
			PhoneNumber: &model.PhoneNumberInput{
				PhoneNumber: contact.PhoneNumber,
			},
			Description: &contact.Description,
		}

		contactId, err := t.services.CustomerOSApiClient.CreateContact(tenant, username, contactInput)
		if err != nil {
			return err
		}

		socialInput := cosModel.SocialInput{
			Url: *contact.Social,
		}
		_, err = t.services.CustomerOsClient.AddSocialContact(tenant, username, contactId, socialInput)
		if err != nil {
			return err
		}

		noteInput := cosModel.NoteInput{
			Content: contact.Note,
		}
		_, err = t.services.CustomerOsClient.CreateNoteForContact(tenant, username, contactId, noteInput)
		if err != nil {
			return err
		}

		contactIds = append(contactIds, EmailAddressWithId{
			Email: contactInput.Email.Email,
			Id:    contactId,
		})
	}

	//create tenant billingProfile
	for _, tenantBillingProfile := range sourceData.TenantBillingProfiles {
		tenantBillingProfileInput := cosModel.TenantBillingProfileInput{
			LegalName:                     tenantBillingProfile.LegalName,
			Email:                         tenantBillingProfile.Email,
			AddressLine1:                  tenantBillingProfile.AddressLine1,
			Locality:                      tenantBillingProfile.Locality,
			Country:                       tenantBillingProfile.Country,
			Zip:                           tenantBillingProfile.Zip,
			DomesticPaymentsBankInfo:      tenantBillingProfile.DomesticPaymentsBankInfo,
			InternationalPaymentsBankInfo: tenantBillingProfile.InternationalPaymentsBankInfo,
			VatNumber:                     tenantBillingProfile.VatNumber,
			SendInvoicesFrom:              tenantBillingProfile.SendInvoicesFrom,
			CanPayWithCard:                tenantBillingProfile.CanPayWithCard,
			CanPayWithDirectDebitSEPA:     tenantBillingProfile.CanPayWithDirectDebitSEPA,
			CanPayWithDirectDebitACH:      tenantBillingProfile.CanPayWithDirectDebitACH,
			CanPayWithDirectDebitBacs:     tenantBillingProfile.CanPayWithDirectDebitBacs,
			CanPayWithPigeon:              tenantBillingProfile.CanPayWithPigeon,
			CanPayWithBankTransfer:        tenantBillingProfile.CanPayWithBankTransfer,
		}
		tenantBillingProfileId, err := t.services.CustomerOsClient.CreateTenantBillingProfile(tenant, username, tenantBillingProfileInput)
		if err != nil {
			return err
		}
		if tenantBillingProfileId == "" {
			return errors.New("tenantBillingProfileId is nil")
		}

	}
	//create orgs
	for _, organization := range sourceData.Organizations {

		var organizationId string
		if organization.Id != "" {
			organizationId = organization.Id
		} else {
			var err error
			b := model.MarketB2b
			var organizationInput = model.OrganizationInput{
				Name:      &organization.Name,
				Website:   &organization.Website,
				Logo:      &organization.Logo,
				Notes:     &organization.Notes,
				Industry:  &organization.Industry,
				Market:    &b,
				Employees: &organization.Employees,
				Relationship: func() *model.OrganizationRelationship {
					rel := &organization.Relationship
					return (*model.OrganizationRelationship)(rel)
				}(),
				Domains: []string{
					organization.Domain,
				},
			}
			organizationId, err = t.services.CustomerOSApiClient.CreateOrganization(tenant, username, organizationInput)
			if err != nil {
				return err
			}

			var organizationUpdateInput = cosModel.OrganizationUpdateInput{
				Id:                organizationId,
				LastFundingAmount: &organization.LastFundingAmount,
				LastFundingRound:  &organization.LastFundingRound,
				TargetAudience:    &organization.TargetAudience,
				ValueProposition:  &organization.ValueProposition,
			}
			organizationId, err = t.services.CustomerOsClient.UpdateOrganization(tenant, username, organizationUpdateInput)
			if err != nil {
				return err
			}

			socialInput := cosModel.SocialInput{
				Url: *organization.OrganizationSocial,
			}
			_, err = t.services.CustomerOsClient.AddSocialOrganization(tenant, username, organizationId, socialInput)
			if err != nil {
				return err
			}

			for _, onboardingStatusInput := range organization.OnboardingStatusInput {
				if onboardingStatusInput.Status != "" {
					organizationOnboardingStatus := cosModel.OrganizationUpdateOnboardingStatus{
						OrganizationId: organizationId,
						Status:         onboardingStatusInput.Status,
						Comments:       onboardingStatusInput.Comments,
					}
					_, err := t.services.CustomerOsClient.UpdateOrganizationOnboardingStatus(tenant, username, organizationOnboardingStatus)
					if err != nil {
						return err
					}
				}
			}
		}

		//create Contracts with Service Lines in org
		for _, contract := range organization.Contracts {
			contractInput := cosModel.ContractInput{
				OrganizationId:          organizationId,
				ContractName:            contract.ContractName,
				CommittedPeriodInMonths: contract.CommittedPeriodInMonths,
				ContractUrl:             contract.ContractUrl,
				ServiceStarted:          contract.ServiceStarted,
				ContractSigned:          contract.ContractSigned,
				Approved:                contract.Approved,
			}
			contractId, err := t.services.CustomerOsClient.CreateContract(tenant, username, contractInput)
			if err != nil {
				return err
			}
			if contractId == "" {
				return errors.New("contractId is nil")
			}

			repository.WaitForNodeCreatedInNeo4jWithConfig(ctx, span, t.services.CommonServices.Neo4jRepositories, contractId, model2.NodeLabelContract, 10*time.Second)

			contractUpdateInput := cosModel.ContractUpdateInput{
				ContractId:            contractId,
				Patch:                 true,
				InvoicingStartDate:    contract.InvoicingStartDate,
				BillingCycle:          contract.BillingCycle,
				Currency:              contract.Currency,
				AddressLine1:          contract.AddressLine1,
				AddressLine2:          contract.AddressLine2,
				Zip:                   contract.Zip,
				Locality:              contract.Locality,
				Country:               contract.Country,
				OrganizationLegalName: contract.OrganizationLegalName,
				InvoiceEmail:          contract.InvoiceEmail,
				InvoiceNote:           contract.InvoiceNote,
			}
			contractId, err = t.services.CustomerOsClient.UpdateContract(tenant, username, contractUpdateInput)
			if err != nil {
				return err
			}
			if contractId == "" {
				return errors.New("contractId is nil")
			}

			for _, serviceLine := range contract.ServiceLines {

				serviceLineInput := func() interface{} {
					if serviceLine.ServiceEnded == nil {
						return cosModel.ServiceLineInput{
							ContractId:     contractId,
							Description:    serviceLine.Description,
							BillingCycle:   serviceLine.BillingCycle,
							Price:          serviceLine.Price,
							Quantity:       serviceLine.Quantity,
							ServiceStarted: serviceLine.ServiceStarted,
						}
					}
					return cosModel.ServiceLineEndedInput{
						ContractId:     contractId,
						Description:    serviceLine.Description,
						BillingCycle:   serviceLine.BillingCycle,
						Price:          serviceLine.Price,
						Quantity:       serviceLine.Quantity,
						ServiceStarted: serviceLine.ServiceStarted,
						ServiceEnded:   serviceLine.ServiceEnded,
					}
				}()
				serviceLineId, err := t.services.CustomerOsClient.CreateServiceLine(tenant, username, serviceLineInput)
				if err != nil {
					return err
				}

				if serviceLineId == "" {
					return errors.New("serviceLineId is nil")
				}

				repository.WaitForNodeCreatedInNeo4jWithConfig(ctx, span, t.services.CommonServices.Neo4jRepositories, serviceLineId, model2.NodeLabelServiceLineItem, 10*time.Second)
			}

			invoiceId, err := t.services.CustomerOsClient.DryRunNextInvoiceForContractInput(tenant, username, contractId)
			if err != nil {
				return err
			}
			if invoiceId == "" {
				return errors.New("invoiceId is nil")
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
				return errors.New("contactId is nil")
			}

			err := t.services.CustomerOsClient.AddContactToOrganization(tenant, username, contactId, organizationId, people.JobRole, people.Description)
			if err != nil {
				return err
			}
		}

		//create emails
		for _, email := range organization.Emails {
			sig, _ := uuid.NewUUID()
			sigs := sig.String()

			channelValue := "EMAIL"
			appSource := appSource
			sessionStatus := "ACTIVE"
			sessionType := "THREAD"
			sessionOpts := []InteractionSessionBuilderOption{
				WithSessionIdentifier(&sigs),
				WithSessionChannel(&channelValue),
				WithSessionName(&email.Subject),
				WithSessionAppSource(&appSource),
				WithSessionStatus(&sessionStatus),
				WithSessionType(&sessionType),
			}

			sessionId, err := t.services.CustomerOsClient.CreateInteractionSession(tenant, username, sessionOpts...)
			if sessionId == nil {
				return errors.New("sessionId is nil")
			}

			participantTypeTO, participantTypeCC, participantTypeBCC := "TO", "CC", "BCC"
			participantsTO := toParticipantInputArr(email.To, &participantTypeTO)
			participantsCC := toParticipantInputArr(email.Cc, &participantTypeCC)
			participantsBCC := toParticipantInputArr(email.Bcc, &participantTypeBCC)
			sentTo := append(append(participantsTO, participantsCC...), participantsBCC...)
			sentBy := toParticipantInputArr([]string{email.From}, nil)

			emailChannelData, err := dto.BuildEmailChannelData("", "", email.Subject, nil, nil)
			if err != nil {
				return err
			}

			iig, err := uuid.NewUUID()
			if err != nil {
				return err
			}
			iigs := iig.String()
			eventOpts := []InteractionEventBuilderOption{
				WithCreatedAt(&email.Date),
				WithSessionId(sessionId),
				WithEventIdentifier(iigs),
				WithChannel(&channelValue),
				WithChannelData(emailChannelData),
				WithContent(&email.Body),
				WithContentType(&email.ContentType),
				WithSentBy(sentBy),
				WithSentTo(sentTo),
				WithAppSource(&appSource),
			}

			interactionEventId, err := t.services.CustomerOsClient.CreateInteractionEvent(tenant, username, eventOpts...)
			if err != nil {
				return err
			}

			if interactionEventId == nil {
				return errors.New("interactionEventId is nil")
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
			meetingId, err := t.services.CustomerOsClient.CreateMeeting(tenant, username, input)
			if err != nil {
				return err
			}

			if meetingId == "" {
				return errors.New("meetingId is nil")
			}

			eventType := "meeting"
			eventOpts := []InteractionEventBuilderOption{
				WithSentBy([]cosModel.InteractionEventParticipantInput{*getInteractionEventParticipantInput(meeting.CreatedBy, userIds, contactIds)}),
				WithSentTo(getInteractionEventParticipantInputList(meeting.Attendees, userIds, contactIds)),
				WithMeetingId(&meetingId),
				WithCreatedAt(&meeting.StartedAt),
				WithEventType(&eventType),
				WithAppSource(&appSource),
			}

			interactionEventId, err := t.services.CustomerOsClient.CreateInteractionEvent(tenant, username, eventOpts...)
			if err != nil {
				return err
			}

			if interactionEventId == nil {
				return errors.New("interactionEventId is nil")
			}
		}

		//log entries
		for _, logEntry := range organization.LogEntries {

			interactionEventId, err := t.services.CustomerOsClient.CreateLogEntry(tenant, username, organizationId, logEntry.CreatedBy, logEntry.Content, logEntry.ContentType, logEntry.Date)
			if err != nil {
				return err
			}

			if interactionEventId == nil {
				return errors.New("interactionEventId is nil")
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

			_, err := t.services.GrpcClients.IssueClient.UpsertIssue(ctx, &issueGrpcRequest)
			if err != nil {
				return err
			}
		}

		//slack
		for _, slackThread := range organization.Slack {

			sig, err := uuid.NewUUID()
			if err != nil {
				return err
			}
			sigs := sig.String()

			channelValue := "CHAT"
			appSource := appSource
			sessionStatus := "ACTIVE"
			sessionType := "THREAD"
			sessionName := slackThread[0].Message
			sessionOpts := []InteractionSessionBuilderOption{
				WithSessionIdentifier(&sigs),
				WithSessionChannel(&channelValue),
				WithSessionName(&sessionName),
				WithSessionAppSource(&appSource),
				WithSessionStatus(&sessionStatus),
				WithSessionType(&sessionType),
			}

			sessionId, err := t.services.CustomerOsClient.CreateInteractionSession(tenant, username, sessionOpts...)
			if sessionId == nil {
				return errors.New("sessionId is nil")
			}

			for _, slackMessage := range slackThread {

				sentBy := toParticipantInputArr([]string{slackMessage.CreatedBy}, nil)

				iig, err := uuid.NewUUID()
				if err != nil {
					return err
				}
				iigs := iig.String()
				eventType := "MESSAGE"
				contentType := "text/plain"
				eventOpts := []InteractionEventBuilderOption{
					WithCreatedAt(&slackMessage.CreatedAt),
					WithSessionId(sessionId),
					WithEventIdentifier(iigs),
					WithExternalId(iigs),
					WithExternalSystemId("slack"),
					WithChannel(&channelValue),
					WithEventType(&eventType),
					WithContent(&slackMessage.Message),
					WithContentType(&contentType),
					WithSentBy(sentBy),
					WithAppSource(&appSource),
				}

				interactionEventId, err := t.services.CustomerOsClient.CreateInteractionEvent(tenant, username, eventOpts...)
				if err != nil {
					return err
				}

				if interactionEventId == nil {
					return errors.New("interactionEventId is nil")
				}

			}

		}

		//intercom
		for _, intercomThread := range organization.Intercom {

			sig, err := uuid.NewUUID()
			if err != nil {
				return err
			}
			sigs := sig.String()

			channelValue := "CHAT"
			appSource := appSource
			sessionStatus := "ACTIVE"
			sessionType := "THREAD"
			sessionName := intercomThread[0].Message
			sessionOpts := []InteractionSessionBuilderOption{
				WithSessionIdentifier(&sigs),
				WithSessionChannel(&channelValue),
				WithSessionName(&sessionName),
				WithSessionAppSource(&appSource),
				WithSessionStatus(&sessionStatus),
				WithSessionType(&sessionType),
			}

			sessionId, err := t.services.CustomerOsClient.CreateInteractionSession(tenant, username, sessionOpts...)
			if sessionId == nil {
				return errors.New("sessionId is nil")
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
					return err
				}
				iigs := iig.String()
				eventType := "MESSAGE"
				contentType := "text/html"
				eventOpts := []InteractionEventBuilderOption{
					WithCreatedAt(&intercomMessage.CreatedAt),
					WithSessionId(sessionId),
					WithEventIdentifier(iigs),
					WithExternalId(iigs),
					WithExternalSystemId("intercom"),
					WithChannel(&channelValue),
					WithEventType(&eventType),
					WithContent(&intercomMessage.Message),
					WithContentType(&contentType),
					WithSentBy(sentBy),
					WithAppSource(&appSource),
				}
				interactionEventId, err := t.services.CustomerOsClient.CreateInteractionEvent(tenant, username, eventOpts...)
				if err != nil {
					return err
				}

				if interactionEventId == nil {
					return errors.New("interactionEventId is nil")
				}

			}

		}
	}

	for _, masterPlan := range sourceData.MasterPlans {
		masterPlanId, err := t.services.CustomerOsClient.CreateMasterPlan(tenant, username, masterPlan.Name)
		if err != nil {
			return err
		}
		if masterPlanId == "" {
			return errors.New("masterPlanId is nil")
		}
		for _, milestone := range masterPlan.Milestones {
			masterPlanMilestoneInput := cosModel.MasterPlanMilestoneInput{
				MasterPlanId:  masterPlanId,
				Name:          milestone.Name,
				Order:         milestone.Order,
				DurationHours: milestone.DurationHours,
				Optional:      milestone.Optional,
				Items:         milestone.Items,
			}
			masterPlanMilestoneId, err := t.services.CustomerOsClient.CreateMasterPlanMilestone(tenant, username, masterPlanMilestoneInput)
			if err != nil {
				return err
			}
			if masterPlanMilestoneId == "" {
				return errors.New("masterPlanMilestoneId is nil")
			}
		}
	}

	return nil
}

func (t *tenantDataInjector) CleanupTenantData(ctx context.Context, tenant, username, reqTenant, reqConfirmTenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantDataInjector.CleanupTenantData")
	defer span.Finish()

	return t.services.CustomerOsClient.HardDeleteTenant(ctx, tenant, username, reqTenant, reqConfirmTenant)
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
