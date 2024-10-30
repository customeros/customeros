import type { Operation } from '@store/types';
import type { Transport } from '@store/transport';
import type { rdiffResult } from 'recursive-diff';

import get from 'lodash/get';
import { P, match } from 'ts-pattern';
import { makePayload } from '@store/util';

import {
  type Tag,
  OnboardingStatus,
  type OrganizationUpdateInput,
} from '@graphql/types';

import type { OrganizationStore } from '../Organization.store';

import AddTagDocument from './addTag.graphql';
import AddSocialDocument from './addSocial.graphql';
import RemoveTagDocument from './removeTag.graphql';
import RemoveOwnerDocument from './removeOwner.graphql';
import UpdateOwnerDocument from './updateOwner.graphql';
import UpdateSocialDocument from './updateSocial.graphql';
import RemoveSocialDocument from './removeSocial.graphql';
import AddSubsidiaryDocument from './addSubsidiary.graphql';
import GetOrganizationDocument from './getOrganization.graphql';
import GetOrganizationsDocument from './getOrganizations.graphql';
import SaveOrganizationDocument from './saveOrganization.graphql';
import RemoveSubsidiaryDocument from './removeSubsidiary.graphql';
import HideOrganizationsDocument from './hideOrganizations.graphql';
import MergeOrganizationsDocument from './mergeOrganizations.graphql';
import UpdateOrganizationDocument from './updateOrganization.graphql';
import UpdateOnboardingStatusDocument from './updateOnboardingStatus.graphql';
import GetArchivedOrganizationsAfterDocument from './getArchivedOrganizations.graphql';
import UpdateAllOpportunityRenewalsDocument from './updateAllOpportunityRenewals.graphql';
import {
  AddSocialMutation,
  AddSocialMutationVariables,
} from './addSocial.generated';
import {
  OrganizationQuery,
  OrganizationQueryVariables,
} from './getOrganization.generated';
import {
  UpdateSocialMutation,
  UpdateSocialMutationVariables,
} from './updateSocial.generated';
import {
  RemoveSocialMutation,
  RemoveSocialMutationVariables,
} from './removeSocial.generated';
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
  HideOrganizationsMutation,
  HideOrganizationsMutationVariables,
} from './hideOrganizations.generated';
import {
  SetOrganizationOwnerMutation,
  SetOrganizationOwnerMutationVariables,
} from './updateOwner.generated';
import {
  MergeOrganizationsMutation,
  MergeOrganizationsMutationVariables,
} from './mergeOrganizations.generated';
import {
  UpdateOrganizationMutation,
  UpdateOrganizationMutationVariables,
} from './updateOrganization.generated';
import {
  RemoveOrganizationOwnerMutation,
  RemoveOrganizationOwnerMutationVariables,
} from './removeOwner.generated';
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
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  public static getInstance(transport: Transport): OrganizationsService {
    if (!OrganizationsService.instance) {
      OrganizationsService.instance = new OrganizationsService(transport);
    }

    return OrganizationsService.instance;
  }

  async getOrganization(id: string) {
    return this.transport.graphql.request<
      OrganizationQuery,
      OrganizationQueryVariables
    >(GetOrganizationDocument, { id });
  }

  async getOrganizations(payload: GetOrganizationsQueryVariables) {
    return this.transport.graphql.request<
      GetOrganizationsQuery,
      GetOrganizationsQueryVariables
    >(GetOrganizationsDocument, payload);
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

  async updateOwner(payload: SetOrganizationOwnerMutationVariables) {
    return this.transport.graphql.request<
      SetOrganizationOwnerMutation,
      SetOrganizationOwnerMutationVariables
    >(UpdateOwnerDocument, payload);
  }

  async removeOwner(payload: RemoveOrganizationOwnerMutationVariables) {
    return this.transport.graphql.request<
      RemoveOrganizationOwnerMutation,
      RemoveOrganizationOwnerMutationVariables
    >(RemoveOwnerDocument, payload);
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

  public async mutateOperation(operation: Operation, store: OrganizationStore) {
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
            ownerId: store.value.owner?.id || '',
          },
        });
      })
      .with(['contracts', ...P.array()], () => {})
      .with(['contacts', ...P.array()], () => {})
      .with(['accountDetails', 'renewalSummary', ...P.array()], async () => {
        const amount =
          store.value.accountDetails?.renewalSummary?.arrForecast ?? 0;
        const potentialAmount =
          store.value.accountDetails?.renewalSummary?.maxArrForecast ?? 0;
        const rate =
          amount === 0 || potentialAmount === 0
            ? 0
            : (amount / potentialAmount) * 100;

        await this.updateAllOpportunityRenewals({
          input: {
            organizationId,
            renewalAdjustedRate: rate,
            renewalLikelihood:
              store.value.accountDetails?.renewalSummary?.renewalLikelihood,
          },
        });
      })
      .with(['accountDetails', 'onboarding', ...P.array()], async () => {
        await this.updateOnboardingStatus({
          input: {
            organizationId,
            status:
              store.value?.accountDetails?.onboarding?.status ??
              OnboardingStatus.NotApplicable,
            comments: store.value?.accountDetails?.onboarding?.comments ?? '',
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

            const foundSocial = get(store, `value.socialMedia[${index}]`, null);

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
      .with(['parentCompanies', ...P.array()], async () => {})
      .with(['tags', ...P.array()], () => {
        match(type)
          .with('add', async () => {
            await this.addTag({
              input: {
                organizationId,
                tag: { id: value.id, name: value.name },
              },
            });
          })
          .with('delete', async () => {
            await this.removeTag({
              input: { organizationId, tag: { id: oldValue.id } },
            });
          })
          .with('update', async () => {
            if (!oldValue) {
              (value as Array<Tag>)?.forEach(async (tag) => {
                await this.addTag({
                  input: {
                    organizationId,
                    tag: { id: tag?.id, name: tag?.name },
                  },
                });
              });
            }

            if (oldValue) {
              await this.removeTag({
                input: { organizationId, tag: { id: oldValue.id } },
              });
            }
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
