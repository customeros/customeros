package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	emailDbNode, err := h.repositories.EmailRepository.GetEmailForUser(ctx, tenant, userId)

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

	// target user
	userDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, userId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var user *entity.UserEntity
	if userDbNode != nil {
		user = graph_db.MapDbNodeToUserEntity(*userDbNode)
	}

	// actor user
	actorDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, actorUserId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var actor *entity.UserEntity
	if userDbNode != nil {
		actor = graph_db.MapDbNodeToUserEntity(*actorDbNode)
	}

	// Organization
	orgDbNode, err := h.repositories.OrganizationRepository.GetOrganization(ctx, tenant, orgId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.OrganizationRepository.GetOrganization")
	}
	var org *entity.OrganizationEntity
	if orgDbNode != nil {
		org = graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)
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

	// call notification service
	err = h.notificationProvider.SendNotification(ctx, &NotifiableUser{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        email.Email,
		SubscriberID: userId,
	}, payload, workflowId)

	return err
}

func (h *OrganizationEventHandler) notificationProviderSendInAppNotification(ctx context.Context, span opentracing.Span, workflowId, userId, actorUserId, orgId, tenant string) error {
	///////////////////////////////////       Get User, Actor, Org Content       ///////////////////////////////////
	// target user email
	emailDbNode, err := h.repositories.EmailRepository.GetEmailForUser(ctx, tenant, userId)

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
	var user *entity.UserEntity
	if userDbNode != nil {
		user = graph_db.MapDbNodeToUserEntity(*userDbNode)
	}

	// actor user
	actorDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, actorUserId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.UserRepository.GetUser")
	}
	var actor *entity.UserEntity
	if userDbNode != nil {
		actor = graph_db.MapDbNodeToUserEntity(*actorDbNode)
	}

	// Organization
	orgDbNode, err := h.repositories.OrganizationRepository.GetOrganization(ctx, tenant, orgId)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "h.repositories.OrganizationRepository.GetOrganization")
	}
	var org *entity.OrganizationEntity
	if orgDbNode != nil {
		org = graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)
	}
	/////////////////////////////////// Notification Provider Payload And Call ///////////////////////////////////

	payload := map[string]interface{}{
		"notificationText": fmt.Sprintf("%s %s made you the owner of %s", actor.FirstName, actor.LastName, org.Name),
		"orgId":            orgId,
		"isArchived":       org.Hide,
	}

	// call notification service
	err = h.notificationProvider.SendNotification(ctx, &NotifiableUser{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        email.Email,
		SubscriberID: userId,
	}, payload, workflowId)

	return err
}

func (h *OrganizationEventHandler) parseOrgOwnerUpdateEmail(actor, target *entity.UserEntity, orgId, orgName string) (string, error) {
	rawMjml, _ := os.ReadFile("./email_templates/ownership.single.mjml")
	mjmlf := strings.Replace(string(rawMjml[:]), "{{userFirstName}}", target.FirstName, -1)
	mjmlf = strings.Replace(mjmlf, "{{actorFirstName}}", actor.FirstName, -1)
	mjmlf = strings.Replace(mjmlf, "{{actorLastName}}", actor.LastName, -1)
	mjmlf = strings.Replace(mjmlf, "{{orgName}}", orgName, -1)
	mjmlf = strings.Replace(mjmlf, "{{orgLink}}", fmt.Sprintf("%s/organization/%s", h.cfg.Subscriptions.NotificationsSubscription.RedirectUrl, orgId), -1)
	mjmlMap := map[string]string{
		"mjml": mjmlf,
	}

	mjmlSecret := h.cfg.Services.MJML.SecretKey
	mjmlAppId := h.cfg.Services.MJML.ApplicationId

	if mjmlSecret == "" || mjmlAppId == "" {
		html, err := mjml.ToHTML(context.Background(), mjmlf) // mjml.WithMinify(true)
		var mjmlError mjml.Error
		if errors.As(err, &mjmlError) {
			return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s", mjmlError.Message)
		}
		return html, err
	}
	mjmlJSON, err := json.Marshal(mjmlMap)
	if err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s", err.Error())
	}
	req, err := http.NewRequest("POST", "https://api.mjml.io/v1/render", bytes.NewReader(mjmlJSON))
	if err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(mjmlAppId, mjmlSecret)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s", err.Error())
	}
	defer response.Body.Close()
	var result struct {
		HTML        string   `json:"html"`
		Errors      []string `json:"errors"`
		MJML        string   `json:"mjml"`
		MJMLVersion string   `json:"mjml_version"`
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s", response.Status)
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("(OrganizationEventHandler.parseOrgOwnerUpdateEmail) error: %s", err.Error())
	}
	return result.HTML, err
}
