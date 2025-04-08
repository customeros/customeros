package contact

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"strings"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type LinkedInContact struct {
	Url        string
	Alias      string
	ExternalID string
}

func NewLinkedInContact(url string, alias string, externalID string) *LinkedInContact {
	return &LinkedInContact{
		Url:        url,
		Alias:      alias,
		ExternalID: externalID,
	}
}

func (l *LinkedInContact) Validate() error {
	if l.Url != "" {
		l.ConvertLinkedinSalesUrlToPublic()
		socialEntity := neo4jentity.SocialEntity{Url: l.Url}
		if !socialEntity.IsLinkedin() {
			return errors.New("not a valid linkedin url")
		}

		// Extract identifier from URL
		identifier := socialEntity.ExtractLinkedinPersonIdentifierFromUrl()

		// If no alias is provided, use the identifier from URL
		if l.Alias == "" {
			l.Alias = identifier
		}
	}

	return nil
}

func (l *LinkedInContact) ConvertLinkedinSalesUrlToPublic() {
	if l.Url == "" {
		return
	}

	if !strings.Contains(l.Url, "/sales/lead") {
		return
	}
	l.Url = strings.Replace(l.Url, "/sales/lead/", "/in/", 1)
}

func (s *contactService) CheckContactExistsWithLinkedIn(ctx context.Context, url, alias, externalId string) (bool, string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContactService.CheckContactExistsWithLinkedIn")
	defer spans.Finish()

	spans.LogKV("url", url, "alias", alias, "externalId", externalId)

	linkedIn := NewLinkedInContact(url, alias, externalId)
	err := linkedIn.Validate()
	if err != nil {
		spans.TraceError(err)
		return false, "", err
	}

	contacts, err := s.neo4j.ContactReadRepository.GetContactsByLinkedIn(
		ctx,
		common.GetTenantFromContext(ctx),
		linkedIn.Url,
		linkedIn.Alias,
		linkedIn.ExternalID,
	)
	if err != nil {
		spans.TraceError(err)
		return false, "", err
	}

	if len(contacts) == 0 {
		return false, "", nil
	}

	contactID := contacts[0].Props["id"].(string)
	return true, contactID, nil
}

func (s *contactService) CreateContactByLinkedIn(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkedInUrl string, options ...common_srv.ServiceOptions) (string, string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContactService.CreateContactByLinkedIn")
	defer spans.Finish()

	spans.LogKV("linkedInUrl", linkedInUrl)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return "", "", err
	}

	linkedin := NewLinkedInContact(linkedInUrl, "", "")
	err = linkedin.Validate()
	if err != nil {
		spans.TraceError(err)
		return "", "", err
	}

	existingID, err := s.handleExistingLinkedinContactCheck(ctx, linkedin, txWithPostCommit, options...)
	if err != nil || existingID != "" {
		return existingID, linkedin.Url, err
	}

	return s.createNewContact(ctx, txWithPostCommit, linkedin.Url, options...)
}

func (s *contactService) handleExistingLinkedinContactCheck(ctx context.Context, linkedin *LinkedInContact, txWithPostCommit *utils.TxWithPostCommit, options ...common_srv.ServiceOptions) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContactService.handleExistingLinkedinContact")
	defer spans.Finish()

	spans.LogObjectAsJson("linkedinContact", &linkedin)

	exists, existingID, err := s.CheckContactExistsWithLinkedIn(ctx, linkedin.Url, linkedin.Alias, "")
	if err != nil || !exists {
		return "", err
	}

	return s.updateExistingContact(ctx, existingID, txWithPostCommit, options...)
}

func (s *contactService) updateExistingContact(ctx context.Context, existingID string, txWithPostCommit *utils.TxWithPostCommit, options ...common_srv.ServiceOptions) (string, error) {
	contactEntity, err := s.GetContactById(ctx, existingID)
	if err != nil {
		return "", errors.Wrap(err, "unable to get contact by id")
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit,
		func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
			return nil, s.processExistingContact(ctx, txWithPostCommit, contactEntity, existingID, options...)
		})

	return existingID, err
}

func (s *contactService) processExistingContact(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactEntity *neo4jentity.ContactEntity, existingID string, options ...common_srv.ServiceOptions) error {
	tenant := common.GetTenantFromContext(ctx)

	if contactEntity.IsHidden() {
		return s.ShowContact(ctx, txWithPostCommit, existingID)
	}

	if err := s.neo4j.CommonWriteRepository.TouchEntity(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelContact, existingID); err != nil {
		return errors.Wrap(err, "error on updating contact updatedAt")
	}

	txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
		if common_srv.PublishCompletedEvents(options...) {
			s.events.Publisher.PublishNotification(ctx, tenant, existingID, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
		}
		return nil
	})

	return nil
}

func (s *contactService) createNewContact(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkedInUrl string, options ...common_srv.ServiceOptions) (string, string, error) {
	var createdContactID string
	var err error

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit,
		func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
			if createdContactID, err = s.Save(ctx, txWithPostCommit, nil, data_fields.ContactFields{}, false, options...); err != nil {
				return "", errors.Wrap(err, "failed to create contact")
			}

			return nil, s.addSocialLink(ctx, txWithPostCommit, createdContactID, linkedInUrl)
		})

	return createdContactID, linkedInUrl, err
}

func (s *contactService) addSocialLink(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contactID, linkedInUrl string) error {
	socialEntity := neo4jentity.SocialEntity{
		Url:       linkedInUrl,
		Source:    neo4jentity.DecodeDataSource(neo4jentity.DataSourceOpenline.String()),
		AppSource: common.GetAppSourceFromContext(ctx),
	}

	_, err := s.social.AddSocialToEntity(ctx, txWithPostCommit,
		common_srv.LinkWith{
			Id:   contactID,
			Type: model.CONTACT,
		},
		socialEntity)
	if err != nil {
		return errors.Wrap(err, "failed to merge social with contact")
	}

	return nil
}
