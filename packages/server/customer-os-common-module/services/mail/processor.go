package mail

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (s *mailService) GetOrganizationIdForEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant, email string, source string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MailService.GetOrganizationIdForEmail")
	defer spans.Finish()
	spans.LogKV("email", email)

	if email == "" {
		return "", nil
	}
	emailSyntax := mailsherpa.ValidateEmailSyntax(email)
	if !emailSyntax.IsValid {
		return "", nil
	}

	// if email and contact exists, return the org
	emailId, err := s.neo4j.EmailReadRepository.GetEmailIdIfExists(ctx, nil, tenant, email)
	if err != nil {
		err = errors.Wrap(err, "failed to get email id for email")
		spans.TraceError(err)
		return "", err
	}
	if emailId != "" {
		contactByEmail, err := s.contact.GetFirstContactByEmail(ctx, email)
		if err != nil {
			err := errors.Wrap(err, "failed to get contact by email")
			spans.TraceError(err)
			return "", err
		}

		if contactByEmail != nil {
			orgsForContact, err := s.org.GetPrimaryOrganizationsWithJobRoleForContacts(ctx, []string{contactByEmail.Id})
			if err != nil {
				err := errors.Wrap(err, "failed to get primary organizations with job role for contacts")
				spans.TraceError(err)
				return "", err
			}

			if orgsForContact != nil && len(*orgsForContact) > 0 {
				return (*orgsForContact)[0].Organization.ID, nil
			}
		}
	}

	output, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		var organizationNode *neo4j.Node
		var organizationId string

		// if it's a personal email or system generated email or role account eamil, we create just the email node in tenant
		domain := utils.ExtractDomainFromEmail(email)
		if domain == "" {
			err = errors.New("unable to extract domain from email: " + email)
			return "", err
		}

		tenantWorkspaces, err := s.workspace.GetWorkspaceDomainsForTenant(ctx)
		if err != nil {
			err = errors.Wrap(err, "failed to get workspace domains for tenant")
			return "", err
		}

		if utils.Contains(tenantWorkspaces, domain) {
			return "", nil
		}

		if utils.Contains(s.cache.GetPersonalEmailProviders(), domain) ||
			emailSyntax.IsSystemGenerated ||
			emailSyntax.IsRoleAccount {
			emailIdPtr, err := s.email.Merge(ctx, txWithPostCommit, tenant, interfaces.EmailFields{
				Email:     email,
				Source:    neo4jentity.DecodeDataSource(source),
				AppSource: constants.AppSourceSyncEmail,
			}, nil)
			if err != nil {
				return "", err
			}
			emailId = utils.IfNotNilString(emailIdPtr)
			return "", nil
		}

		organizationNode, err = s.neo4j.OrganizationReadRepository.GetOrganizationByDomain(ctx, txWithPostCommit.Tx, tenant, domain)
		if err != nil {
			return "", fmt.Errorf("unable to retrieve organization for tenant: %v", err)
		}

		if organizationNode == nil {
			stage := neo4jenum.Lead
			leadSource := ""

			if source == neo4jentity.DataSourceGmail.String() {
				leadSource = "Gmail"
			} else if source == neo4jentity.DataSourceOutlook.String() {
				leadSource = "Outlook"
			} else if source == neo4jentity.DataSourceMailstack.String() {
				leadSource = "Mailstack"
				stage = neo4jenum.Target
			} else {
				leadSource = "Email"
			}

			organizationFields := data_fields.OrganizationFields{
				LeadSource:    utils.StringPtr(leadSource),
				Stage:         utils.ToPtr(stage),
				Relationship:  utils.ToPtr(neo4jenum.OrganizationRelationshipProspect),
				PrimaryDomain: &domain,
				Source:        utils.StringPtr(source),
			}

			organizationId, err = s.org.Save(ctx, txWithPostCommit, nil, organizationFields)
			if err != nil {
				return "", fmt.Errorf("failed to create organization: %w", err)
			}
		} else {
			organizationId = utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id")
		}

		contactId, err := s.contact.CreateContactByEmail(ctx, txWithPostCommit, email)
		if err != nil {
			return "", fmt.Errorf("unable to create contact: %w", err)
		}

		err = s.contact.LinkContactWithOrganization(ctx, txWithPostCommit, contactId, organizationId, "", "", source, false, nil, nil)

		emailId, err = s.neo4j.EmailReadRepository.GetEmailIdIfExists(ctx, txWithPostCommit.Tx, tenant, email)
		if err != nil {
			return "", fmt.Errorf("unable to retrieve email id for tenant: %w", err)
		}

		return organizationId, nil
	})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to get email id for email"))
		return "", err
	}

	return output.(string), nil
}
