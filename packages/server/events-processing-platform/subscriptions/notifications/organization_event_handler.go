package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"

	"github.com/Boostport/mjml-go"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OrganizationEventHandler struct {
	repositories         *repository.Repositories
	log                  logger.Logger
	notificationProvider NotificationProvider
	cfg                  *config.Config
}

func NewOrganizationEventHandler(log logger.Logger, repositories *repository.Repositories, cfg *config.Config) *OrganizationEventHandler {
	return &OrganizationEventHandler{
		repositories:         repositories,
		log:                  log,
		notificationProvider: NewNotificationProvider(log, cfg.Services.Novu.ApiKey),
		cfg:                  cfg,
	}
}

func (h *OrganizationEventHandler) OnOrganizationUpdateOwner(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Notifications.OrganizationEventHandler.OnOrganizationUpdateOwner")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationOwnerUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	if eventData.ActorUserId == eventData.OwnerUserId {
		// do not send notification if the actor is the same as the target
		h.log.Info("actor is the same as the target, skipping notification")
		return nil
	}

	err := h.notificationProviderSendEmail(
		ctx,
		span,
		WorkflowIdOrgOwnerUpdateEmail,
		eventData.OwnerUserId,
		eventData.ActorUserId,
		eventData.OrganizationId,
		eventData.Tenant,
	)

	if err != nil {
		tracing.TraceErr(span, err)
	}

	err = h.notificationProviderSendInAppNotification(
		ctx,
		span,
		WorkflowIdOrgOwnerUpdateAppNotification,
		eventData.OwnerUserId,
		eventData.ActorUserId,
		eventData.OrganizationId,
		eventData.Tenant,
	)

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (h *OrganizationEventHandler) notificationProviderSendEmail(ctx context.Context, span opentracing.Span, workflowId, userId, actorUserId, orgId, tenant string) error {
	///////////////////////////////////       Get User, Actor, Org Content       ///////////////////////////////////
	// target user email
	emailDbNode, err := h.repositories.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.EmailRepository.GetEmailForUser")
	}

	var email *entity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	email = graph_db.MapDbNodeToEmailEntity(*emailDbNode)

	// actor user email
	actorEmailDbNode, err := h.repositories.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.EmailRepository.GetEmailForUser")
	}

	var actorEmail *entity.EmailEntity
	if actorEmailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("actor email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendEmail")
	}
	actorEmail = graph_db.MapDbNodeToEmailEntity(*actorEmailDbNode)

	// target user
	userDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var user *neo4jentity.UserEntity
	if userDbNode != nil {
		user = neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}

	// actor user
	actorDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, actorUserId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var actor *neo4jentity.UserEntity
	if userDbNode != nil {
		actor = neo4jmapper.MapDbNodeToUserEntity(actorDbNode)
	}

	// Organization
	orgDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, orgId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.OrganizationRepository.GetOrganization")
	}
	var org *neo4jentity.OrganizationEntity
	if orgDbNode != nil {
		org = neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	}
	///////////////////////////////////       Get Email Content       ///////////////////////////////////
	html, err := h.parseOrgOwnerUpdateEmail(actor, user, orgId, org.Name)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "notifications.parseOrgOwnerUpdateEmail")
	}
	/////////////////////////////////// Notification Provider Payload And Call ///////////////////////////////////

	payload := map[string]interface{}{
		"html":    html,
		"subject": fmt.Sprintf("%s %s added you as an owner", actor.FirstName, actor.LastName),
		"email":   email.Email,
		"orgName": org.Name,
	}

	overrides := map[string]interface{}{
		"email": map[string]string{
			"replyTo": actorEmail.Email,
		},
	}

	// call notification service
	err = h.notificationProvider.SendNotification(ctx, &NotifiableUser{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        email.Email,
		SubscriberID: userId,
	}, payload, overrides, workflowId)

	return err
}

func (h *OrganizationEventHandler) notificationProviderSendInAppNotification(ctx context.Context, span opentracing.Span, workflowId, userId, actorUserId, orgId, tenant string) error {
	///////////////////////////////////       Get User, Actor, Org Content       ///////////////////////////////////
	// target user email
	emailDbNode, err := h.repositories.Neo4jRepositories.EmailReadRepository.GetEmailForUser(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.EmailRepository.GetEmailForUser")
	}

	var email *entity.EmailEntity
	if emailDbNode == nil {
		tracing.TraceErr(span, err)
		err = errors.New("email db node not found")
		return errors.Wrap(err, "h.notificationProviderSendInAppNotification")
	}
	email = graph_db.MapDbNodeToEmailEntity(*emailDbNode)

	// target user
	userDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var user *neo4jentity.UserEntity
	if userDbNode != nil {
		user = neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	}

	// actor user
	actorDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, actorUserId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var actor *neo4jentity.UserEntity
	if userDbNode != nil {
		actor = neo4jmapper.MapDbNodeToUserEntity(actorDbNode)
	}

	// Organization
	orgDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, orgId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.OrganizationRepository.GetOrganization")
	}
	var org *neo4jentity.OrganizationEntity
	if orgDbNode != nil {
		org = neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
	}
	/////////////////////////////////// Notification Provider Payload And Call ///////////////////////////////////

	payload := map[string]interface{}{
		"notificationText": fmt.Sprintf("%s %s made you the owner of %s", actor.FirstName, actor.LastName, org.Name),
		"orgId":            orgId,
		"isArchived":       org.Hide,
	}

	overrides := map[string]interface{}{}

	// call notification service
	err = h.notificationProvider.SendNotification(ctx, &NotifiableUser{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        email.Email,
		SubscriberID: userId,
	}, payload, overrides, workflowId)

	return err
}

func (h *OrganizationEventHandler) parseOrgOwnerUpdateEmail(actor, target *neo4jentity.UserEntity, orgId, organizationName string) (string, error) {
	orgName := organizationName
	if organizationName == "" {
		orgName = "Unnamed Organization"
	}
	if _, err := os.Stat(h.cfg.Subscriptions.NotificationsSubscription.EmailTemplatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s", err.Error())
	}
	emailPath := fmt.Sprintf("%s/ownership.single.mjml", h.cfg.Subscriptions.NotificationsSubscription.EmailTemplatePath)
	if _, err := os.Stat(emailPath); err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s", err.Error())
	}

	rawMjml, _ := os.ReadFile(emailPath)
	mjmlf := strings.Replace(string(rawMjml[:]), "{{userFirstName}}", target.FirstName, -1)
	mjmlf = strings.Replace(mjmlf, "{{actorFirstName}}", actor.FirstName, -1)
	mjmlf = strings.Replace(mjmlf, "{{actorLastName}}", actor.LastName, -1)
	mjmlf = strings.Replace(mjmlf, "{{orgName}}", orgName, -1)
	mjmlf = strings.Replace(mjmlf, "{{orgLink}}", fmt.Sprintf("%s/organization/%s", h.cfg.Subscriptions.NotificationsSubscription.RedirectUrl, orgId), -1)

	html, err := mjml.ToHTML(context.Background(), mjmlf) // mjml.WithMinify(true)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		mjmlSecret := h.cfg.Services.MJML.SecretKey
		mjmlAppId := h.cfg.Services.MJML.ApplicationId
		html, err = mjmlToHtmlApi(mjmlf, mjmlAppId, mjmlSecret)
		if err != nil {
			return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s:%s", mjmlError.Message, err.Error())
		}
	}
	return html, err
}

func mjmlToHtmlApi(mjml, mjmlAppId, mjmlSecret string) (string, error) {
	mjmlMap := map[string]string{
		"mjml": mjml,
	}
	mjmlJSON, err := json.Marshal(mjmlMap)
	if err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.mjmlToHtmlApi) error: %s", err.Error())
	}
	req, err := http.NewRequest("POST", "https://api.mjml.io/v1/render", bytes.NewReader(mjmlJSON))
	if err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.mjmlToHtmlApi) error: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(mjmlAppId, mjmlSecret)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.mjmlToHtmlApi) error: %s", err.Error())
	}
	defer response.Body.Close()
	var result struct {
		HTML        string   `json:"html"`
		Errors      []string `json:"errors"`
		MJML        string   `json:"mjml"`
		MJMLVersion string   `json:"mjml_version"`
	}

	if response.StatusCode != http.StatusOK {
		var badResponse struct {
			Message   string `json:"message"`
			RequestID string `json:"request_id"`
			StartedAt string `json:"started_at"`
		}
		err = json.NewDecoder(response.Body).Decode(&badResponse)
		if err != nil {
			return "", fmt.Errorf("(OrganizationEventHandler.mjmlToHtmlApi) error: %s", err.Error())
		}
		return "", fmt.Errorf("(OrganizationEventHandler.mjmlToHtmlApi) error: %s: %s", response.Status, badResponse.Message)
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.mjmlToHtmlApi) error: %s", err.Error())
	}
	return result.HTML, err
}
