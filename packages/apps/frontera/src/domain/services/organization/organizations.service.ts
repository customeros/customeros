import { runInAction } from 'mobx';
import { RootStore } from '@store/root';
import { TagStore } from '@store/Tags/Tag.store';
import {
  Organization,
  OrganizationDatum,
} from '@store/Organizations/Organization.dto';
import { OrganizationRepository } from '@infra/repositories/organization/organization.repository.ts';

import { unwrap } from '@utils/unwrap';
import {
  Social,
  Domain,
  EntityType,
  FlagWrongFields,
  SortingDirection,
  OrganizationStage,
  ComparisonOperator,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@shared/types/__generated__/graphql.types';

export class OrganizationService {
  private root = RootStore.getInstance();
  private orgRepo = OrganizationRepository.getInstance();

  constructor() {}

  public async searchTenant(searchTerm: string) {
    try {
      const { ui_organizations_search } =
        await this.orgRepo.searchOrganizations({
          limit: 30,
          sort: {
            by: 'ORGANIZATIONS_NAME',
            caseSensitive: false,
            direction: SortingDirection.Asc,
          },
          where: {
            OR: [
              {
                filter: {
                  property: 'ORGANIZATIONS_NAME',
                  value: searchTerm,
                  caseSensitive: false,
                  includeEmpty: false,
                  operation: ComparisonOperator.Contains,
                },
              },
              {
                filter: {
                  property: 'ORGANIZATIONS_WEBSITE',
                  value: searchTerm,
                  caseSensitive: false,
                  includeEmpty: false,
                  operation: ComparisonOperator.Contains,
                },
              },
            ],
          },
        });

      const results = ui_organizations_search?.ids ?? [];

      if (results.length === 0) {
        return;
      }

      await this.root.organizations.retrieve(results);
    } catch (_err) {
      throw new Error('Failed to search for tenant');
    }
  }

  public async merge(primaryId: string, secondaryIds: string[]) {
    try {
      const { organization_Merge } = await this.orgRepo.mergeOrganizations({
        primaryOrganizationId: primaryId,
        mergedOrganizationIds: secondaryIds,
      });

      runInAction(() => {
        if (organization_Merge.id) {
          this.root.organizations.mergeOrganizations(primaryId, secondaryIds);

          this.root.ui.toastSuccess(`Merged companies`, `merge-${primaryId}`);
        }
      });
    } catch (err) {
      throw new Error('Failed to merge companies');
    }
  }

  public async flagWrongField(id: string, field: FlagWrongFields) {
    const organization = this.root.organizations.getById(id);

    try {
      const { flagWrongField } = await this.orgRepo.flagWrongField({
        input: {
          entityId: id,
          entityType: EntityType.Organization,
          field,
        },
      });

      runInAction(() => {
        if (flagWrongField?.result) {
          organization?.flagIncorrectIndustry();
          this.root.ui.toastSuccess(
            `Noted, we're looking into it`,
            `flag-field-${field}`,
          );
        }

        if (!flagWrongField?.result) {
          throw new Error('Failed to flag wrong field');
        }
      });
    } catch (err) {
      throw new Error('Failed to flag wrong field');
    }
  }

  public async addTag(organization: Organization, tag: TagStore) {
    organization.addTag(tag.id);

    const [res, err] = await unwrap(
      this.orgRepo.addTag({
        input: { organizationId: organization.id, tag: { name: tag.tagName } },
      }),
    );

    if (err) {
      console.error(err);
      organization.deleteTag(tag.id);

      this.root.ui.toastError('Failed to add tag to company', 'tag-add-failed');

      return [null, err];
    }

    return [res, err];
  }

  public async removeTag(organization: Organization, tag: TagStore) {
    organization.deleteTag(tag.id);

    const [res, err] = await unwrap(
      this.orgRepo.removeTag({
        input: { organizationId: organization.id, tag: { name: tag.tagName } },
      }),
    );

    if (err) {
      console.error(err);
      organization.addTag(tag.id);

      this.root.ui.toastError(
        'Failed to remove tag from company',
        'tag-remove-failed',
      );

      return [null, err];
    }

    return [res, err];
  }

  public async removeSocialMediaItem(
    organization: Organization,
    social: Social,
  ) {
    organization.deleteSocialMedia(social.id);

    const [res, err] = await unwrap(
      this.orgRepo.removeSocial({
        socialId: social.id,
      }),
    );

    if (err) {
      console.error(err);
      organization.revertSocialMedia(social);

      this.root.ui.toastError(
        "We couldn't remove this link",
        `${social.id}-remove-social-media`,
      );

      return [null, err];
    }

    return [res, err];
  }

  public async addDomain(organization: Organization, domain: Domain) {
    organization.addDomain(domain);

    const [res, err] = await unwrap(
      this.orgRepo.addDomain({
        organizationId: organization.id,
        domain: domain.domain,
      }),
    );

    if (err) {
      console.error(err);
      organization.deleteDomain(domain.domain);

      this.root.ui.toastError(
        "We couldn't add this domain",
        `${domain.domain}-add-domain`,
      );

      return [null, err];
    }

    return [res, err];
  }

  public async removeDomain(organization: Organization, domain: string) {
    const domainDetailsToRestore = organization.value.domainsDetails.find(
      (d) => d.domain !== domain,
    );

    organization.deleteDomain(domain);

    const [res, err] = await unwrap(
      this.orgRepo.removeDomain({
        organizationId: organization.id,
        domain: domain,
      }),
    );

    if (err) {
      console.error(err);

      if (domainDetailsToRestore) {
        organization.addDomain(domainDetailsToRestore);
      }

      this.root.ui.toastError(
        "We couldn't remove this domain",
        `${domain}-remove-domain`,
      );

      return [null, err];
    }

    return [res, err];
  }

  public async editOwner(orgId: string, ownerId: string) {
    const prevOwner = this.root.organizations.getById(orgId)?.value.owner
      ? this.root.organizations.getById(orgId)?.value.owner
      : null;

    this.root.organizations.getById(orgId)?.setOwner(ownerId);

    const [res, err] = await unwrap(
      this.orgRepo.saveOrganization({
        input: {
          id: orgId,
          ownerId,
        },
      }),
    );

    if (err) {
      console.error(err);
      this.root.organizations.getById(orgId)?.setOwner(prevOwner?.id ?? '');

      this.root.ui.toastError(
        "We couldn't change the owner",
        `${ownerId}-edit-owner`,
      );

      return [null, err];
    }

    if (res) {
      this.root.ui.toastSuccess(
        'Owner changed successfully',
        `${ownerId}-edit-owner`,
      );
    }

    return [res, err];
  }

  public async setRelationshipAndStage(
    payload: Partial<OrganizationDatum>,
    orgId: string,
  ) {
    const organization = this.root.organizations.getById(orgId);
    const prevRelationship = organization?.value.relationship;

    organization?.setRelationship(
      payload.relationship ?? OrganizationRelationship.Prospect,
    );
    organization?.setStage(
      payload.stage ??
        this.root.organizations.getById(orgId)?.value.stage ??
        OrganizationStage.Target,
    );

    const [res, err] = await unwrap(
      this.orgRepo.saveOrganization({
        input: {
          id: orgId,
          relationship:
            payload.relationship ?? OrganizationRelationship.Prospect,
          stage: this.root.organizations.getById(orgId)?.value.stage,
        },
      }),
    );

    if (err) {
      console.error(err);
      this.root.organizations
        .getById(orgId)
        ?.setRelationship(
          prevRelationship ?? OrganizationRelationship.Prospect,
        );
    }

    return [res, err];
  }

  public async updateOrganizationHealth(
    orgId: string,
    renewalLikelihood: OpportunityRenewalLikelihood,
  ) {
    const prevRenewalLikelihood =
      this.root.organizations.getById(orgId)?.value
        .renewalSummaryRenewalLikelihood;

    this.root.organizations
      .getById(orgId)
      ?.setRenewalAdjustedRate(renewalLikelihood);

    const amount =
      this.root.organizations.getById(orgId)?.value
        ?.renewalSummaryArrForecast ?? 0;
    const potentialAmount =
      this.root.organizations.getById(orgId)?.value
        ?.renewalSummaryMaxArrForecast ?? 0;
    const rate =
      amount === 0 || potentialAmount === 0
        ? 0
        : (amount / potentialAmount) * 100;

    const [res, err] = await unwrap(
      this.orgRepo.updateAllOpportunityRenewals({
        input: {
          organizationId: orgId,
          renewalAdjustedRate: rate,
          renewalLikelihood:
            this.root.organizations.getById(orgId)?.value
              .renewalSummaryRenewalLikelihood,
        },
      }),
    );

    if (err) {
      console.error(err);

      this.root.organizations
        .getById(orgId)
        ?.setRenewalAdjustedRate(
          prevRenewalLikelihood ?? OpportunityRenewalLikelihood.MediumRenewal,
        );
      this.root.ui.toastError(
        'Failed to update opportunity renewals',
        'update-opportunity-renewals',
      );

      return [null, err];
    }

    return [res, err];
  }

  public async updateNotes(organization: Organization, notes: string) {
    const prevNotes = organization.value.notes ?? '';

    organization.setNotes(notes);

    const [res, err] = await unwrap(
      this.orgRepo.saveOrganization({
        input: {
          id: organization.id,
          notes,
        },
      }),
    );

    if (err) {
      console.error(err);
      organization.setNotes(prevNotes);

      this.root.ui.toastError('Failed to update notes', 'update-notes');

      return [null, err];
    }

    return [res, err];
  }
}
