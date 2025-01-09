import type { Operation } from '@store/types';
import type { rdiffResult } from 'recursive-diff';

import get from 'lodash/get';
import { P, match } from 'ts-pattern';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';

import {
  type Tag,
  EntityType,
  OnboardingStatus,
  type OrganizationUpdateInput,
} from '@graphql/types';

import type { Organization } from '../Organization.dto';

import AddTagDocument from './addTag.graphql';
import AddSocialDocument from './addSocial.graphql';
import RemoveTagDocument from './removeTag.graphql';
import AddDomainDocument from './addDomain.graphql';
import CheckDomainDocument from './checkDomain.graphql';
import UpdateSocialDocument from './updateSocial.graphql';
import RemoveSocialDocument from './removeSocial.graphql';
import CheckWebsiteDocument from './checkWebsite.graphql';
import RemoveDomainDocument from './removeDomain.graphql';
import RemoveDomainsDocument from './removeDomains.graphql';
import AddSubsidiaryDocument from './addSubsidiary.graphql';
import GetOrganizationsDocument from './getOrganizations.graphql';
import SaveOrganizationDocument from './saveOrganization.graphql';
import RemoveSubsidiaryDocument from './removeSubsidiary.graphql';
import HideOrganizationsDocument from './hideOrganizations.graphql';
import ImportOrganizationDocument from './importOrganization.graphql';
import MergeOrganizationsDocument from './mergeOrganizations.graphql';
import UpdateOrganizationDocument from './updateOrganization.graphql';
import SearchOrganizationsDocument from './searchOrganizations.graphql';
import GetOrganizationsByIdsDocument from './getOrganizationsByIds.graphql';
import UpdateOnboardingStatusDocument from './updateOnboardingStatus.graphql';
import SearchGlobalOrganizationsDocument from './searchGlobalOrganizations.graphql';
import GetArchivedOrganizationsAfterDocument from './getArchivedOrganizations.graphql';
import UpdateAllOpportunityRenewalsDocument from './updateAllOpportunityRenewals.graphql';
import {
  AddSocialMutation,
  AddSocialMutationVariables,
} from './addSocial.generated';
import {
  CheckDomainQuery,
  CheckDomainQueryVariables,
} from './checkDomain.generated';
import {
  AddDomainMutation,
  AddDomainMutationVariables,
} from './addDomain.generated';
import {
  CheckWebsiteQuery,
  CheckWebsiteQueryVariables,
} from './checkWebsite.generated';
import {
  UpdateSocialMutation,
  UpdateSocialMutationVariables,
} from './updateSocial.generated';
import {
  RemoveSocialMutation,
  RemoveSocialMutationVariables,
} from './removeSocial.generated';
import {
  RemoveDomainMutation,
  RemoveDomainMutationVariables,
} from './removeDomain.generated';
import {
  RemoveDomainsMutation,
  RemoveDomainsMutationVariables,
} from './removeDomains.generated';
import {
  GetOrganizationsQuery,
  GetOrganizationsQueryVariables,
} from './getOrganizations.generated';
import {
  SaveOrganizationMutation,
  SaveOrganizationMutationVariables,
} from './saveOrganization.generated';
import {
  AddTagsToOrganizationMutation,
  AddTagsToOrganizationMutationVariables,
} from './addTag.generated';
import {
  SearchOrganizationsQuery,
  SearchOrganizationsQueryVariables,
} from './searchOrganizations.generated';
import {
  HideOrganizationsMutation,
  HideOrganizationsMutationVariables,
} from './hideOrganizations.generated';
import {
  MergeOrganizationsMutation,
  MergeOrganizationsMutationVariables,
} from './mergeOrganizations.generated';
import {
  UpdateOrganizationMutation,
  UpdateOrganizationMutationVariables,
} from './updateOrganization.generated';
import {
  ImportOrganizationMutation,
  ImportOrganizationMutationVariables,
} from './importOrganization.generated';
import {
  GetOrganizationsByIdsQuery,
  GetOrganizationsByIdsQueryVariables,
} from './getOrganizationsByIds.generated';
import {
  RemoveTagFromOrganizationMutation,
  RemoveTagFromOrganizationMutationVariables,
} from './removeTag.generated';
import {
  UpdateOnboardingStatusMutation,
  UpdateOnboardingStatusMutationVariables,
} from './updateOnboardingStatus.generated';
import {
  AddSubsidiaryToOrganizationMutation,
  AddSubsidiaryToOrganizationMutationVariables,
} from './addSubsidiary.generated';
import {
  SearchGlobalOrganizationsQuery,
  SearchGlobalOrganizationsQueryVariables,
} from './searchGlobalOrganizations.generated';
import {
  GetArchivedOrganizationsAfterQuery,
  GetArchivedOrganizationsAfterQueryVariables,
} from './getArchivedOrganizations.generated';
import {
  RemoveSubsidiaryToOrganizationMutation,
  RemoveSubsidiaryToOrganizationMutationVariables,
} from './removeSubsidiary.generated';
import {
  BulkUpdateOpportunityRenewalMutation,
  BulkUpdateOpportunityRenewalMutationVariables,
} from './updateAllOpportunityRenewals.generated';

export class OrganizationsService {
  private static instance: OrganizationsService | null = null;
  private transport = Transport.getInstance();

  constructor() {}

  public static getInstance(): OrganizationsService {
    if (!OrganizationsService.instance) {
      OrganizationsService.instance = new OrganizationsService();
    }

    return OrganizationsService.instance;
  }

  async searchOrganizations(payload: SearchOrganizationsQueryVariables) {
    return this.transport.graphql.request<
      SearchOrganizationsQuery,
      SearchOrganizationsQueryVariables
    >(SearchOrganizationsDocument, payload);
  }

  async importOrganization(payload: ImportOrganizationMutationVariables) {
    return this.transport.graphql.request<
      ImportOrganizationMutation,
      ImportOrganizationMutationVariables
    >(ImportOrganizationDocument, payload);
  }

  async searchGlobalOrganizations(
    payload: SearchGlobalOrganizationsQueryVariables,
  ) {
    return this.transport.graphql.request<
      SearchGlobalOrganizationsQuery,
      SearchGlobalOrganizationsQueryVariables
    >(SearchGlobalOrganizationsDocument, payload);
  }

  async getOrganizations(payload: GetOrganizationsQueryVariables) {
    return this.transport.graphql.request<
      GetOrganizationsQuery,
      GetOrganizationsQueryVariables
    >(GetOrganizationsDocument, payload);
  }

  async getOrganizationsByIds(payload: GetOrganizationsByIdsQueryVariables) {
    return this.transport.graphql.request<
      GetOrganizationsByIdsQuery,
      GetOrganizationsByIdsQueryVariables
    >(GetOrganizationsByIdsDocument, payload);
  }

  async getOrganization(id: string) {
    const { ui_organizations } = await this.getOrganizationsByIds({
      ids: [id],
    });

    return ui_organizations[0];
  }

  async checkWebsite(payload: CheckWebsiteQueryVariables) {
    return this.transport.graphql.request<
      CheckWebsiteQuery,
      CheckWebsiteQueryVariables
    >(CheckWebsiteDocument, payload);
  }

  async getArchivedOrganizationsAfter(
    payload: GetArchivedOrganizationsAfterQueryVariables,
  ) {
    return this.transport.graphql.request<
      GetArchivedOrganizationsAfterQuery,
      GetArchivedOrganizationsAfterQueryVariables
    >(GetArchivedOrganizationsAfterDocument, payload);
  }

  async saveOrganization(payload: SaveOrganizationMutationVariables) {
    return this.transport.graphql.request<
      SaveOrganizationMutation,
      SaveOrganizationMutationVariables
    >(SaveOrganizationDocument, payload);
  }

  async removeDomain(payload: RemoveDomainMutationVariables) {
    return this.transport.graphql.request<
      RemoveDomainMutation,
      RemoveDomainMutationVariables
    >(RemoveDomainDocument, payload);
  }

  async addDomain(payload: AddDomainMutationVariables) {
    return this.transport.graphql.request<
      AddDomainMutation,
      AddDomainMutationVariables
    >(AddDomainDocument, payload);
  }

  async removeDomains(payload: RemoveDomainsMutationVariables) {
    return this.transport.graphql.request<
      RemoveDomainsMutation,
      RemoveDomainsMutationVariables
    >(RemoveDomainsDocument, payload);
  }

  async checkDomain(payload: CheckDomainQueryVariables) {
    return this.transport.graphql.request<
      CheckDomainQuery,
      CheckDomainQueryVariables
    >(CheckDomainDocument, payload);
  }

  async hideOrganizations(payload: HideOrganizationsMutationVariables) {
    return this.transport.graphql.request<
      HideOrganizationsMutation,
      HideOrganizationsMutationVariables
    >(HideOrganizationsDocument, payload);
  }

  async mergeOrganizations(payload: MergeOrganizationsMutationVariables) {
    return this.transport.graphql.request<
      MergeOrganizationsMutation,
      MergeOrganizationsMutationVariables
    >(MergeOrganizationsDocument, payload);
  }

  /**
   * @deprecated
   * use saveOrganization instead
   * */
  async updateOrganization(payload: UpdateOrganizationMutationVariables) {
    return this.transport.graphql.request<
      UpdateOrganizationMutation,
      UpdateOrganizationMutationVariables
    >(UpdateOrganizationDocument, payload);
  }

  async addSocial(payload: AddSocialMutationVariables) {
    return this.transport.graphql.request<
      AddSocialMutation,
      AddSocialMutationVariables
    >(AddSocialDocument, payload);
  }

  async removeSocial(payload: RemoveSocialMutationVariables) {
    return this.transport.graphql.request<
      RemoveSocialMutation,
      RemoveSocialMutationVariables
    >(RemoveSocialDocument, payload);
  }

  async addTag(payload: AddTagsToOrganizationMutationVariables) {
    return this.transport.graphql.request<
      AddTagsToOrganizationMutation,
      AddTagsToOrganizationMutationVariables
    >(AddTagDocument, payload);
  }

  async removeTag(payload: RemoveTagFromOrganizationMutationVariables) {
    return this.transport.graphql.request<
      RemoveTagFromOrganizationMutation,
      RemoveTagFromOrganizationMutationVariables
    >(RemoveTagDocument, payload);
  }

  async updateSocial(payload: UpdateSocialMutationVariables) {
    return this.transport.graphql.request<
      UpdateSocialMutation,
      UpdateSocialMutationVariables
    >(UpdateSocialDocument, payload);
  }

  async updateAllOpportunityRenewals(
    payload: BulkUpdateOpportunityRenewalMutationVariables,
  ) {
    return this.transport.graphql.request<
      BulkUpdateOpportunityRenewalMutation,
      BulkUpdateOpportunityRenewalMutationVariables
    >(UpdateAllOpportunityRenewalsDocument, payload);
  }

  async addSubsidiary(payload: AddSubsidiaryToOrganizationMutationVariables) {
    return this.transport.graphql.request<
      AddSubsidiaryToOrganizationMutation,
      AddSubsidiaryToOrganizationMutationVariables
    >(AddSubsidiaryDocument, payload);
  }

  async removeSubsidiary(
    payload: RemoveSubsidiaryToOrganizationMutationVariables,
  ) {
    return this.transport.graphql.request<
      RemoveSubsidiaryToOrganizationMutation,
      RemoveSubsidiaryToOrganizationMutationVariables
    >(RemoveSubsidiaryDocument, payload);
  }

  async updateOnboardingStatus(
    payload: UpdateOnboardingStatusMutationVariables,
  ) {
    return this.transport.graphql.request<
      UpdateOnboardingStatusMutation,
      UpdateOnboardingStatusMutationVariables
    >(UpdateOnboardingStatusDocument, payload);
  }

  public async mutateOperation(operation: Operation, store: Organization) {
    const diff = operation.diff?.[0];
    const type = diff?.op;
    const path = diff?.path;
    const value = diff?.val;
    const organizationId = operation?.entityId;
    const oldValue = (diff as rdiffResult & { oldVal: unknown })?.oldVal;

    if (!operation.diff.length) {
      return;
    }

    if (!organizationId) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }

    return match(path)
      .with(['owner', ...P.array()], () => {
        this.saveOrganization({
          input: {
            id: organizationId,
            ownerId: store?.owner?.id || '',
          },
        });
      })
      .with(['contracts', ...P.array()], () => {})
      .with(['contacts', ...P.array()], () => {})
      .with([P.string.startsWith('renewalSummary'), ...P.array()], async () => {
        const amount = store?.value.renewalSummaryArrForecast ?? 0;
        const potentialAmount = store?.value.renewalSummaryMaxArrForecast ?? 0;
        const rate =
          amount === 0 || potentialAmount === 0
            ? 0
            : (amount / potentialAmount) * 100;

        await this.updateAllOpportunityRenewals({
          input: {
            organizationId,
            renewalAdjustedRate: rate,
            renewalLikelihood: store.value.renewalSummaryRenewalLikelihood,
          },
        });
      })
      .with([P.string.startsWith('onboarding'), ...P.array()], async () => {
        await this.updateOnboardingStatus({
          input: {
            organizationId,
            status:
              store?.value.onboardingStatus ?? OnboardingStatus.NotApplicable,
            comments: store?.value.onboardingComments ?? '',
          },
        });
      })
      .with(['socialMedia', ...P.array()], () => {
        match(type)
          .with('add', async () => {
            await this.addSocial({
              organizationId,
              input: {
                url: value.url,
              },
            });
          })
          .with('update', async () => {
            const index = path[1] as number;

            const foundSocial = get(store.value, `socialMedia[${index}]`, null);

            if (!foundSocial) return;

            await this.updateSocial({
              input: { id: foundSocial.id, url: foundSocial.url },
            });
          })
          .with(
            'delete',
            async () => await this.removeSocial({ socialId: oldValue?.id }),
          );
      })
      .with(['subsidiaries', ...P.array()], async () => {
        if (type === 'delete') {
          const subsidiaryId = oldValue?.organization?.metadata?.id;

          await this.removeSubsidiary({ organizationId, subsidiaryId });

          return;
        }

        const subsidiaryId = match(typeof value)
          .with('string', () => value)
          .otherwise(
            () =>
              value[0]?.organization?.metadata?.id ||
              value?.organization?.metadata?.id,
          );

        if (typeof value === 'string' && type === 'update') {
          this.removeSubsidiary({
            organizationId: value,
            subsidiaryId: oldValue,
          });

          return;
        }

        await this.addSubsidiary({
          input: { organizationId, subsidiaryId, removeExisting: false },
        });
      })
      .with([P.union('parentId', 'parentName'), ...P.array()], async () => {})
      .with(['tags', ...P.array()], () => {
        match(type)
          .with('add', async () => {
            await this.addTag({
              input: {
                organizationId,
                tag: {
                  name: value.name,
                  entityType: EntityType.Organization,
                },
              },
            });
          })
          .with('delete', async () => {
            await this.removeTag({
              input: { organizationId, tag: { id: oldValue.metadata.id } },
            });
          })
          .with('update', async () => {
            match(operation.diff)
              .with(
                [
                  { op: 'update', path: ['tags', P.number, 'name'] },
                  {
                    op: 'update',
                    path: ['tags', P.number, 'metadata', 'id'],
                  },
                  ...P.array(),
                  {
                    op: 'delete',
                    path: ['tags', P.number],
                  },
                ],
                async () => {
                  const oldValue = (
                    operation.diff[1] as rdiffResult & {
                      oldVal: unknown;
                    }
                  )?.oldVal;

                  await this.removeTag({
                    input: {
                      organizationId,
                      tag: {
                        id: oldValue,
                      },
                    },
                  });
                },
              )
              .otherwise(async () => {
                if (!oldValue) {
                  (value as Array<Tag>)?.forEach(async (tag) => {
                    await this.addTag({
                      input: {
                        organizationId,
                        tag: {
                          name: tag?.name,
                          entityType: EntityType.Organization,
                        },
                      },
                    });
                  });
                }

                if (oldValue) {
                  await this.removeTag({
                    input: {
                      organizationId,
                      tag: { id: oldValue.metadata.id },
                    },
                  });
                }
              });
          });
      })
      .with(['updatedAt'], () => undefined)
      .with(['domainsDetails', ...P.array()], async () => {
        if (type === 'update' || type === 'delete') {
          return;

          return await this.removeDomain({
            organizationId,
            domain: oldValue,
          });
        }

        return await this.addDomain({
          organizationId,
          domain: value?.domain,
        });
      })
      .otherwise(async () => {
        const payload = makePayload<OrganizationUpdateInput>(operation);

        return await this.saveOrganization({
          input: { ...payload, id: organizationId },
        });
      });
  }
}
