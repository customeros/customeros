package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/DusanKasan/parsemail"
	mimemail "github.com/emersion/go-message/mail"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"google.golang.org/api/gmail/v1"
	"io"
	"log"
	"net/mail"
	"strings"
	"time"
)

type mailService struct {
	services *Services
	config   *c.Config
}

type MailService interface {
	SaveMail(email *parsemail.Email, tenant *string, user *string) (*model.InteractionEventCreateResponse, error)
	SendMail(request *model.MailReplyRequest, username *string) (*parsemail.Email, error)
}

func (s *mailService) SaveMail(email *parsemail.Email, tenant *string, user *string) (*model.InteractionEventCreateResponse, error) {
	refSize := len(email.References)
	threadId := ""
	if refSize > 0 {
		threadId = utils.EnsureEmailRfcId(email.References[0])
	} else {
		threadId = utils.EnsureEmailRfcId(email.MessageID)
	}

	sessionId, err := s.services.CustomerOsService.GetInteractionSession(&threadId, tenant, user)

	if err != nil {
		log.Printf("failed retriving interaction session: error=%s", err.Error())
	}

	channelValue := "EMAIL"
	appSource := "COMMS_API"
	sessionStatus := "ACTIVE"
	sessionType := "THREAD"
	if sessionId == nil {
		sessionOpts := []SessionOption{
			WithSessionIdentifier(&threadId),
			WithSessionChannel(&channelValue),
			WithSessionName(&email.Subject),
			WithSessionAppSource(&appSource),
			WithSessionStatus(&sessionStatus),
			WithSessionTenant(tenant),
			WithSessionUsername(user),
			WithSessionType(&sessionType),
		}

		sessionId, err = s.services.CustomerOsService.CreateInteractionSession(sessionOpts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create interaction session: %v", err)
		}
		log.Printf("interaction session created: %s", *sessionId)
	}

	participantTypeTO, participantTypeCC, participantTypeBCC := "TO", "CC", "BCC"
	participantsTO := toParticipantInputArr(email.To, &participantTypeTO)
	participantsCC := toParticipantInputArr(email.Cc, &participantTypeCC)
	participantsBCC := toParticipantInputArr(email.Bcc, &participantTypeBCC)
	sentTo := append(append(participantsTO, participantsCC...), participantsBCC...)
	sentBy := toParticipantInputArr(email.From, nil)

	emailChannelData, err := buildEmailChannelData(email, err)
	if err != nil {
		return nil, err
	}
	eventOpts := []EventOption{
		WithTenant(tenant),
		WithUsername(user),
		WithSessionId(sessionId),
		WithEventIdentifier(utils.EnsureEmailRfcId(email.MessageID)),
		WithExternalId(utils.EnsureEmailRfcId(email.MessageID)),
		WithExternalSystemId("gmail"),
		WithChannel(&channelValue),
		WithChannelData(emailChannelData),
		WithContent(utils.FirstNotEmpty(email.HTMLBody, email.TextBody)),
		WithContentType(&email.ContentType),
		WithSentBy(sentBy),
		WithSentTo(sentTo),
		WithAppSource(&appSource),
	}
	response, err := s.services.CustomerOsService.CreateInteractionEvent(eventOpts...)

	if err != nil {
		return nil, fmt.Errorf("failed to create interaction event: %v", err)
	}

	return response, nil
}

func buildEmailChannelData(email *parsemail.Email, err error) (*string, error) {
	emailContent := model.EmailChannelData{
		Subject:   email.Subject,
		InReplyTo: utils.EnsureEmailRfcIds(email.InReplyTo),
		Reference: utils.EnsureEmailRfcIds(email.References),
	}
	jsonContent, err := json.Marshal(emailContent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email content: %v", err)
	}
	jsonContentString := string(jsonContent)

	return &jsonContentString, nil
}

func (s *mailService) SendMail(request *model.MailReplyRequest, username *string) (*parsemail.Email, error) {
	retMail := parsemail.Email{}
	retMail.HTMLBody = request.Content
	subject := request.Subject
	var h mimemail.Header

	tenant, err := s.services.CustomerOsService.GetTenant(username)
	if err != nil {
		log.Printf("unable to retrieve tenant for %s", *username)
		return nil, fmt.Errorf("unable to retrieve tenant for %s", *username)
	}

	gSrv, err := s.services.AuthServices.GoogleService.GetGmailService(*username, tenant.Tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}

	if gSrv == nil {
		return nil, fmt.Errorf("unable to build a gmail service with service account or auth token: %v", err)
	}

	user, err := s.services.CustomerOsService.GetUserByEmail(username)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve user for %s", *username)
	}

	sentByName := ""
	if user.UserByEmail.Name != nil && *user.UserByEmail.Name != "" {
		sentByName = *user.UserByEmail.Name
	} else if user.UserByEmail.FirstName != nil && user.UserByEmail.LastName != nil {
		if *user.UserByEmail.FirstName != "" && *user.UserByEmail.LastName != "" {
			sentByName = *user.UserByEmail.FirstName + " " + *user.UserByEmail.LastName
		} else if *user.UserByEmail.LastName != "" {
			sentByName = *user.UserByEmail.LastName
		} else if *user.UserByEmail.FirstName != "" {
			sentByName = *user.UserByEmail.FirstName
		}
	} else {
		sentByName = *username
	}

	fromAddress := []*mimemail.Address{{sentByName, *username}}
	retMail.From = fromAddress
	var toAddress []*mimemail.Address
	var ccAddress []*mimemail.Address
	var bccAddress []*mimemail.Address
	for _, to := range request.To {
		toAddress = append(toAddress, &mimemail.Address{Address: to})
		retMail.To = toAddress
	}
	if request.Cc != nil {
		for _, cc := range request.Cc {
			ccAddress = append(ccAddress, &mimemail.Address{Address: cc})
			retMail.Cc = toAddress
		}
	}
	if request.Bcc != nil {
		for _, bcc := range request.Bcc {
			bccAddress = append(bccAddress, &mimemail.Address{Address: bcc})
			retMail.Bcc = toAddress
		}
	}

	var b bytes.Buffer

	h.SetDate(time.Now())
	h.SetAddressList("From", fromAddress)
	h.SetAddressList("To", toAddress)
	h.SetAddressList("Cc", ccAddress)
	h.SetAddressList("Bcc", bccAddress)

	if subject != nil {
		h.SetSubject(*subject)
	}

	if request.ReplyTo != nil {
		event, err := s.services.CustomerOsService.GetInteractionEvent(request.ReplyTo, username)
		if err != nil {
			return nil, err
		}
		emailChannelData := model.EmailChannelData{}
		err = json.Unmarshal([]byte(event.InteractionEvent.ChannelData), &emailChannelData)
		if err != nil {
			log.Printf("unable to parse email channel data for %s", *request.ReplyTo)
			return nil, fmt.Errorf("unable to parse email channel data for %s", *request.ReplyTo)
		}
		subject = &emailChannelData.Subject

		retMail.References = append(emailChannelData.Reference, event.InteractionEvent.EventIdentifier)
		retMail.InReplyTo = []string{event.InteractionEvent.EventIdentifier}

		h.Set("References", strings.Join(retMail.References, " "))
		h.Set("In-Reply-To", event.InteractionEvent.EventIdentifier)
	}

	// Create a new mail writer
	mw, err := mimemail.CreateWriter(&b, h)
	if err != nil {
		log.Fatal(err)
	}

	// Create a text part
	tw, err := mw.CreateInline()
	if err != nil {
		log.Fatal(err)
	}
	var th mimemail.InlineHeader
	th.Set("Content-Type", "text/html")
	w, err := tw.CreatePart(th)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, request.Content)
	w.Close()
	tw.Close()

	mw.Close()

	raw := base64.StdEncoding.EncodeToString(b.Bytes())
	msgToSend := &gmail.Message{
		Raw: raw,
	}
	result, err := gSrv.Users.Messages.Send("me", msgToSend).Do()
	if err != nil {
		log.Printf("Unable to send email: %v", err)
		return nil, err
	}

	generatedMessage, err := gSrv.Users.Messages.Get("me", result.Id).Do()
	if err != nil {
		log.Printf("Unable to get email: %v", err)
		return nil, err
	}
	for _, header := range generatedMessage.Payload.Headers {
		log.Printf("Comparing %s to %s", header.Name, "Message-ID")
		if strings.EqualFold(header.Name, "Message-ID") {
			retMail.MessageID = header.Value
			retMail.References = append(retMail.References, header.Value)
			break
		}
	}

	if subject != nil {
		retMail.Subject = *subject
	}

	log.Printf("Email successfully sent id %v", retMail.MessageID)
	return &retMail, nil
}

func toParticipantInputArr(from []*mail.Address, participantType *string) []cosModel.InteractionEventParticipantInput {
	var to []cosModel.InteractionEventParticipantInput
	for _, a := range from {
		participantInput := cosModel.InteractionEventParticipantInput{
			Email: &a.Address,
			Type:  participantType,
		}
		to = append(to, participantInput)
	}
	return to
}

func NewMailService(config *c.Config, services *Services) MailService {
	return &mailService{
		config:   config,
		services: services,
	}
}
