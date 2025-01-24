import { runInAction } from 'mobx';
import { RootStore } from '@store/root';
import { TagStore } from '@store/Tags/Tag.store';
import { Organization } from '@store/Organizations/Organization.dto';
import { OrganizationsService } from '@store/Organizations/__service__/Organizations.service';
import { OrganizationRepository } from '@infra/repositories/organization/organization.repository.ts';

import { unwrap } from '@shared/util/unwrap';
import {
  Social,
  Domain,
  EntityType,
  FlagWrongFields,
  SortingDirection,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

export class OrganizationService {
  private root = RootStore.getInstance();
  private orgRepo = OrganizationsService.getInstance();
  private orgInfraRepo = OrganizationRepository.getInstance();

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

          this.root.ui.toastSuccess(
            `Merged organizations`,
            `merge-${primaryId}`,
          );
        }
      });
    } catch (err) {
      throw new Error('Failed to merge organizations');
    }
  }

  public async flagWrongField(id: string, field: FlagWrongFields) {
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

      this.root.ui.toastError(
        'Failed to add tag to organization',
        'tag-add-failed',
      );

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
        'Failed to remove tag from organization',
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
      this.orgInfraRepo.addDomain({
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
      this.orgInfraRepo.removeDomain({
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
}
