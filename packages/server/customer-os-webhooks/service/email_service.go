package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
)

type EmailService interface {
	GetMessageId(data *model.PostmarkEmailWebhookData) (string, error)
	GetReferences(data *model.PostmarkEmailWebhookData) (string, error)
	GetInReplyTo(data *model.PostmarkEmailWebhookData) (string, error)
	IsExcluded(ctx context.Context, tenant string, data *model.PostmarkEmailWebhookData) bool
	ExtractParticipants(tenant string, data *model.PostmarkEmailWebhookData) []string
	ResolveMailbox(ctx context.Context, tenant string, participants []string) (string, error)
	StoreIfNotExists(ctx context.Context, tenant, username, messageId string, data *model.PostmarkEmailWebhookData) error
}

type emailService struct {
	services *Services
	logger   logger.Logger
}

func NewEmailService(services *Services) EmailService {
	return &emailService{
		services: services,
		logger:   services.Logger,
	}
}

func (s *emailService) GetMessageId(data *model.PostmarkEmailWebhookData) (string, error) {
	messageId := ""
	for _, header := range data.Headers {
		if strings.EqualFold(header.Name, "message-id") {
			messageId = header.Value
			break
		}
	}
	if messageId == "" {
		return "", errors.New("Message-ID not found in headers")
	}
	return messageId, nil
}

func (s *emailService) GetReferences(data *model.PostmarkEmailWebhookData) (string, error) {
	for _, header := range data.Headers {
		if header.Name == "References" {
			return header.Value, nil
		}
	}
	return "", nil
}

func (s *emailService) GetInReplyTo(data *model.PostmarkEmailWebhookData) (string, error) {
	for _, header := range data.Headers {
		if header.Name == "In-Reply-To" {
			return header.Value, nil
		}
	}
	return "", nil
}

func (s *emailService) IsExcluded(ctx context.Context, tenant string, data *model.PostmarkEmailWebhookData) bool {
	htmlData := strings.ReplaceAll(data.HtmlBody, "&amp;", "&")
	textData := strings.ReplaceAll(data.TextBody, "&amp;", "&")

	emailExclusion := s.services.CommonServices.Cache.GetEmailExclusion(tenant)
	for _, exclusion := range emailExclusion {
		if s.matchesExclusion(exclusion, data.Subject, htmlData, textData) {
			return true
		}
	}
	return false
}

func (s *emailService) matchesExclusion(exclusion dto.EmailExclusion, subject, htmlBody, textBody string) bool {
	if exclusion.ExcludeSubject != nil && strings.Contains(subject, *exclusion.ExcludeSubject) {
		return true
	}
	if exclusion.ExcludeBody != nil {
		return strings.Contains(htmlBody, *exclusion.ExcludeBody) || strings.Contains(textBody, *exclusion.ExcludeBody)
	}
	return false
}

func (s *emailService) ExtractParticipants(tenant string, data *model.PostmarkEmailWebhookData) []string {
	participants := make([]string, 0)
	tenantEmail := "bcc@" + strings.ToLower(tenant) + ".customeros.ai"

	// Add From participant
	participants = append(participants, data.FromFull.Email)

	// Add To participants
	for _, to := range data.ToFull {
		if to.Email != "" {
			participants = append(participants, to.Email)
		}
	}

	// Add CC participants
	if data.CcFull != nil {
		for _, cc := range data.CcFull {
			if cc.Email != "" {
				participants = append(participants, cc.Email)
			}
		}
	}

	// Add BCC participants
	if data.BccFull != nil {
		for _, bcc := range data.BccFull {
			if bcc.Email != "" && bcc.Email != tenantEmail {
				participants = append(participants, bcc.Email)
			}
		}
	}

	return participants
}

func (s *emailService) ResolveMailbox(ctx context.Context, tenant string, participants []string) (string, error) {
	for _, email := range participants {
		user, err := s.services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, email)
		if err != nil {
			s.logger.Errorf("Error getting user by email: %s", err.Error())
			continue
		}
		if user != nil {
			return email, nil
		}
	}
	return "", errors.New("no valid mailbox found")
}

func (s *emailService) StoreIfNotExists(ctx context.Context, tenant, username, messageId string, data *model.PostmarkEmailWebhookData) error {
	exists, err := s.services.CommonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(
		ctx, "mailstack", tenant, username, messageId,
	)
	if err != nil {
		return err
	}

	if !exists {
		emailRawData, err := s.mapToEmailRawData(tenant, data)
		if err != nil {
			return err
		}

		jsonContent, err := s.marshalEmailData(emailRawData)
		if err != nil {
			return err
		}

		return s.services.CommonServices.PostgresRepositories.RawEmailRepository.Store(
			ctx,
			"mailstack",
			tenant,
			username,
			emailRawData.ProviderMessageId,
			messageId,
			string(jsonContent),
			emailRawData.Sent,
			entity.REAL_TIME,
		)
	}

	return nil
}

func (s *emailService) mapToEmailRawData(tenant string, data *model.PostmarkEmailWebhookData) (entity.EmailRawData, error) {
	sentTime, err := utils.UnmarshalDateTime(data.Date)
	if err != nil {
		return entity.EmailRawData{}, err
	}

	headers := make(map[string]string)
	for _, header := range data.Headers {
		headers[header.Name] = header.Value
	}

	messageId, err := s.GetMessageId(data)
	if err != nil {
		return entity.EmailRawData{}, err
	}

	references, err := s.GetReferences(data)
	if err != nil {
		return entity.EmailRawData{}, err
	}

	inReplyTo, err := s.GetInReplyTo(data)
	if err != nil {
		return entity.EmailRawData{}, err
	}

	return entity.EmailRawData{
		ProviderMessageId: messageId,
		MessageId:         messageId,
		Sent:              *sentTime,
		Subject:           data.Subject,
		From:              "<" + data.FromFull.Email + ">",
		To:                s.formatEmailList(tenant, data.ToFull),
		Cc:                s.formatEmailList(tenant, data.CcFull),
		Bcc:               s.formatBccList(tenant, data.BccFull),
		Html:              data.HtmlBody,
		Text:              data.TextBody,
		ThreadId:          "",
		InReplyTo:         inReplyTo,
		Reference:         references,
		Headers:           headers,
	}, nil
}

func (s *emailService) formatEmailList(tenant string, emails []model.EmailAddress) string {
	formatted := make([]string, 0)
	for _, e := range emails {
		if e.Email != "" {
			formatted = append(formatted, "<"+e.Email+">")
		}
	}
	return strings.Join(formatted, ", ")
}

func (s *emailService) formatBccList(tenant string, emails []model.EmailAddress) string {
	formatted := make([]string, 0)
	tenantEmail := "bcc@" + strings.ToLower(tenant) + ".customeros.ai"
	for _, e := range emails {
		if e.Email != "" && e.Email != tenantEmail {
			formatted = append(formatted, e.Email)
		}
	}
	return strings.Join(formatted, ", ")
}

func (s *emailService) marshalEmailData(data entity.EmailRawData) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	return buffer.Bytes(), err
}
