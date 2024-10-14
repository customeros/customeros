import type { Operation } from '@store/types';
import type { Transport } from '@store/transport';
import type { rdiffResult } from 'recursive-diff';

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
import RemoveSubsidiaryDocument from './removeSubsidiary.graphql';
import UpdateOrganizationDocument from './updateOrganization.graphql';
import UpdateOnboardingStatusDocument from './updateOnboardingStatus.graphql';
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
  AddTagsToOrganizationMutation,
  AddTagsToOrganizationMutationVariables,
} from './addTag.generated';
import {
  SetOrganizationOwnerMutation,
  SetOrganizationOwnerMutationVariables,
} from './updateOwner.generated';
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

    if (!organizationId) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }

    return match(path)
      .with(['owner', ...P.array()], () => {
        if (type === 'update') {
          match(value)
            .with(null, async () => {
              await this.removeOwner({ organizationId });
            })
            .otherwise(async () => {
              await this.updateOwner({ organizationId, userId: value });
            });
        }
      })
      .with(['contracts', ...P.array()], () => {})
      .with(['accountDetails', 'renewalSummary', ...P.array()], async () => {
        const amount =
          store.value.accountDetails?.renewalSummary?.arrForecast ?? 0;
        const potentialAmount =
          store.value.accountDetails?.renewalSummary?.maxArrForecast ?? 0;
        const rate = (amount / potentialAmount) * 100;

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
        const index = path[1] as number;
        const { id, url } = store.value.socialMedia[index];

        match(type)
          .with(
            'add',
            async () =>
              await this.addSocial({
                organizationId,
                input: {
                  url,
                },
              }),
          )
          .with(
            'update',
            async () => await this.updateSocial({ input: { id, url } }),
          )
          .with(
            'delete',
            async () => await this.removeSocial({ socialId: id }),
          );
      })
      .with(['subsidiaries', ...P.array()], async () => {
        const subsidiaryId = match(typeof value)
          .with('string', () => value)
          .otherwise(
            () =>
              value[0]?.organization?.metadata?.id ||
              value?.organization?.metadata?.id,
          );

        await this.addSubsidiary({
          input: { organizationId, subsidiaryId, removeExisting: false },
        });
      })
      .with(['parentCompanies', ...P.array()], async () => {
        if (type === 'delete') {
          const subsidiaryId = oldValue?.organization?.metadata?.id;

          await this.removeSubsidiary({ organizationId, subsidiaryId });
        }
      })
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

        return await this.updateOrganization({
          input: { ...payload, id: organizationId, patch: true },
        });
      });
  }
}
