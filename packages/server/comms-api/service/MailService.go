package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/DusanKasan/parsemail"
	mimemail "github.com/emersion/go-message/mail"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/util"
	oryClient "github.com/ory/client-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/mail"
	"strings"
	"time"
)

type MailService struct {
	customerOSService *CustomerOSService
	config            *c.Config
	oauthConfig       *oauth2.Config
}

type MailServiceInterface interface {
	SaveMail(email *parsemail.Email, tenant *string, user *string) (*model.InteractionEventCreateResponse, error)
	ReplyToMail(request *model.MailReplyRequest, interactionEventId *string, username *string, identityId *string) (*parsemail.Email, error)
}

func (s *MailService) SaveMail(email *parsemail.Email, tenant *string, user *string) (*model.InteractionEventCreateResponse, error) {
	refSize := len(email.References)
	threadId := ""
	if refSize > 0 {
		threadId = util.EnsureRfcId(email.References[0])
	} else {
		threadId = util.EnsureRfcId(email.MessageID)
	}

	cosService := s.customerOSService
	sessionId, err := cosService.GetInteractionSession(&threadId, tenant, user)

	if err != nil {
		log.Printf("failed retriving interaction session: error=%s", err.Error())
	}

	channelValue := "EMAIL"
	appSource := "COMMS_API"
	sessionStatus := "ACTIVE"
	if sessionId == nil {
		sessionId, err = cosService.CreateInteractionSession(
			cosService.WithSessionIdentifier(&threadId),
			cosService.WithSessionChannel(&channelValue),
			cosService.WithSessionName(&email.Subject),
			cosService.WithSessionAppSource(&appSource),
			cosService.WithSessionStatus(&sessionStatus),
			cosService.WithSessionTenant(tenant))
		if err != nil {
			return nil, fmt.Errorf("failed to create interaction session: %v", err)
		}
		log.Printf("interaction session created: %s", *sessionId)
	}

	participantTypeTO, participantTypeCC := "TO", "CC"
	participantsTO := toParticipantInputArr(email.To, &participantTypeTO)
	participantsCC := toParticipantInputArr(email.Cc, &participantTypeCC)
	sentTo := append(participantsTO, participantsCC...)
	sentBy := toParticipantInputArr(email.From, nil)

	emailChannelData, err := buildEmailChannelData(email, err)
	if err != nil {
		return nil, err
	}

	response, err := cosService.CreateInteractionEvent(
		cosService.WithTenant(tenant),
		cosService.WithSessionId(sessionId),
		cosService.WithEventIdentifier(email.MessageID),
		cosService.WithChannel(&channelValue),
		cosService.WithChannelData(emailChannelData),
		cosService.WithContent(util.FirstNotEmpty(email.HTMLBody, email.TextBody)),
		cosService.WithContentType(&email.ContentType),
		cosService.WithSentBy(sentBy),
		cosService.WithSentTo(sentTo),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create interaction event: %v", err)
	}

	return response, nil
}

func buildEmailChannelData(email *parsemail.Email, err error) (*string, error) {
	emailContent := model.EmailChannelData{
		Subject:   email.Subject,
		InReplyTo: util.EnsureRfcIds(email.InReplyTo),
		Reference: util.EnsureRfcIds(email.References),
	}
	jsonContent, err := json.Marshal(emailContent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email content: %v", err)
	}
	jsonContentString := string(jsonContent)

	return &jsonContentString, nil
}

func (s *MailService) ReplyToMail(request *model.MailReplyRequest, interactionEventId *string, username *string, identityId *string) (*parsemail.Email, error) {
	retMail := parsemail.Email{}
	event, err := s.customerOSService.GetInteractionEvent(interactionEventId, username)
	if err != nil {
		return nil, err
	}
	emailChannelData := model.EmailChannelData{}
	err = json.Unmarshal([]byte(event.InteractionEvent.ChannelData), &emailChannelData)
	if err != nil {
		log.Printf("unable to parse email channel data for %s", *interactionEventId)
		return nil, fmt.Errorf("unable to parse email channel data for %s", *interactionEventId)
	}

	tok, err := s.getMailAuthToken(identityId)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token: %v", err)
	}
	log.Printf("Got Auth Token of %v", tok)
	client := s.oauthConfig.Client(context.Background(), tok)

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))

	fromAddress := []*mimemail.Address{{"", *username}}
	toAddress := []*mimemail.Address{}
	for _, to := range request.Destination {
		toAddress = append(toAddress, &mimemail.Address{"", to})
	}

	var b bytes.Buffer
	user := "me"

	// Create our mail header
	var h mimemail.Header
	h.SetDate(time.Now())
	h.SetAddressList("From", fromAddress)
	h.SetAddressList("To", toAddress)
	h.SetSubject(emailChannelData.Subject)

	if len(emailChannelData.Reference) > 0 {
		h.Set("References", strings.Join(emailChannelData.Reference, " "))
		copy(retMail.References, emailChannelData.Reference)

	}

	if len(emailChannelData.InReplyTo) > 0 {
		h.Set("In-Reply-To", strings.Join(emailChannelData.InReplyTo, " "))
		copy(retMail.InReplyTo, emailChannelData.InReplyTo)
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
	result, err := srv.Users.Messages.Send(user, msgToSend).Do()
	if err != nil {
		log.Printf("Unable to send email: %v", err)
		return nil, err
	}

	generatedMessage, err := srv.Users.Messages.Get("me", result.Id).Do()
	if err != nil {
		log.Printf("Unable to get email: %v", err)
		return nil, err
	}
	for _, header := range generatedMessage.Payload.Headers {
		log.Printf("Comparing %s to %s", header.Name, "Content-ID")
		if strings.EqualFold(header.Name, "Content-ID") {
			retMail.MessageID = header.Value
			retMail.References = append(retMail.References, header.Value)
			break
		}
	}
	retMail.Subject = event.InteractionEvent.Subject
	log.Printf("Email successfully sent id %v", retMail.MessageID)
	return &retMail, nil
}

func toParticipantInputArr(from []*mail.Address, participantType *string) []model.InteractionEventParticipantInput {
	var to []model.InteractionEventParticipantInput
	for _, a := range from {
		participantInput := model.InteractionEventParticipantInput{
			Email:           &a.Address,
			ParticipantType: participantType,
		}
		to = append(to, participantInput)
	}
	return to
}

func (s *MailService) getMailAuthToken(identityId *string) (*oauth2.Token, error) {
	configuration := oryClient.NewConfiguration()
	configuration.Servers = []oryClient.ServerConfiguration{
		{
			URL: s.config.GMail.OryServerUrl,
		},
	}
	ory := oryClient.NewAPIClient(configuration)
	ctx := context.Background()
	ctx = context.WithValue(ctx, oryClient.ContextAccessToken, s.config.GMail.OryApiKey)
	identity, _, err := ory.IdentityApi.GetIdentity(ctx, *identityId).IncludeCredential([]string{"oidc"}).Execute()
	if err != nil {
		log.Printf("Unable to get gmail auth token for %s, (%s)", *identityId, err.Error())
		return nil, err
	}
	credentials := identity.GetCredentials()["oidc"]
	log.Printf("Got credentials of %v", credentials)

	providers, ok := credentials.GetConfig()["providers"].([]interface{})
	log.Printf("Got providers of %T", providers[0])

	if !ok {
		log.Printf("unable to get provider list %s", *identityId)
		return nil, err
	}

	provider, ok := providers[0].(map[string]interface{})
	if !ok {
		log.Printf("unable to get provider list %s", *identityId)
		return nil, err
	}
	token, ok := provider["initial_access_token"].(string)

	if !ok {
		log.Printf("unable to get access token %s", *identityId)
		return nil, err
	}
	tok := &oauth2.Token{AccessToken: token, TokenType: "Bearer"}

	refreshToken, ok := provider["initial_refresh_token"].(string)

	if !ok {
		log.Printf("unable to get refresh token`` %s", *identityId)
	} else {
		log.Printf("Setting refresh token to %s", refreshToken)
		tok.RefreshToken = refreshToken
	}
	tok.Expiry = time.Now().Add(time.Hour * -1)
	return tok, nil
}

func NewMailService(config *c.Config, customerOSService *CustomerOSService) *MailService {
	mailService := new(MailService)
	mailService.config = config
	mailService.oauthConfig = &oauth2.Config{
		ClientID:     config.GMail.ClientId,
		ClientSecret: config.GMail.ClientSecret,
		RedirectURL:  strings.Split(config.GMail.RedirectUris, " ")[0],
		Scopes:       []string{gmail.GmailReadonlyScope, gmail.GmailComposeScope, "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
	}
	mailService.customerOSService = customerOSService
	return mailService
}
