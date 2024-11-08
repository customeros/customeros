import set from 'lodash/set';
import merge from 'lodash/merge';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { SyncableGroup } from '@store/syncable-group';
import {
  action,
  computed,
  override,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import {
  relationshipStageMap,
  stageRelationshipMap,
  validRelationshipsForStage,
} from '@utils/orgStageAndRelationshipStatusMap.ts';
import {
  Tag,
  Organization,
  SortingDirection,
  OrganizationInput,
  OrganizationStage,
  ComparisonOperator,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import mock from './mock.json';
import { OrganizationStore } from './Organization.store';
import { OrganizationsService } from './__service__/Organizations.service';

export class OrganizationsStore extends SyncableGroup<
  Organization,
  OrganizationStore
> {
  totalElements = 0;
  private service: OrganizationsService;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, OrganizationStore);
    this.service = OrganizationsService.getInstance(this.transport);

    makeObservable(this, {
      maxLtv: computed,
      hide: action.bound,
      merge: action.bound,
      create: action.bound,
      channelName: override,
      isFullyLoaded: computed,
      updateStage: action.bound,
      totalElements: observable,
      getRecentChanges: override,
    });
  }

  get channelName() {
    return 'Organizations';
  }

  get persisterKey() {
    return 'Organizations';
  }

  get maxLtv() {
    return Math.max(
      ...this.toArray().map(
        (org) => Math.round(org.value.accountDetails?.ltv ?? 0) + 1,
      ),
    );
  }

  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  async bootstrapStream() {
    try {
      await this.transport.stream<Organization>('/organizations', {
        onData: (data) => {
          runInAction(() => {
            this.load([data], { getId: (data) => data.metadata.id });
          });
        },
      });

      runInAction(() => {
        this.isBootstrapped = true;
      });
    } catch (e) {
      runInAction(() => {
        console.error(e);
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async getRecentChanges() {
    try {
      if (this.root.demoMode || this.isBootstrapping) {
        return;
      }

      this.isLoading = true;

      const lastActiveAtUTC = this.root.windowManager
        .getLastActiveAtUTC()
        .toISOString();

      const where = {
        AND: [
          {
            filter: {
              property: 'UPDATED_AT',
              value: lastActiveAtUTC,
              operation: ComparisonOperator.Gte,
            },
          },
        ],
      };

      const { organizations_HiddenAfter: idsToDrop } =
        await this.service.getArchivedOrganizationsAfter({
          date: lastActiveAtUTC,
        });

      const { dashboardView_Organizations } =
        await this.service.getOrganizations({
          pagination: { limit: 1000, page: 0 },
          sort: {
            by: 'LAST_TOUCHPOINT',
            caseSensitive: false,
            direction: SortingDirection.Desc,
          },
          where,
        });

      if (this.isHydrated) {
        await this.drop(idsToDrop);
      } else {
        await this.hydrate({
          idsToDrop,
          getId: (data) => data.metadata.id,
        });
      }

      const data =
        (dashboardView_Organizations?.content as Organization[]) ?? [];
      const totalElements = dashboardView_Organizations?.totalElements;

      this.load(data, {
        getId: (item) => item.metadata.id,
      });
      runInAction(() => {
        this.totalElements = totalElements;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async getAllData() {
    this.isBootstrapping = true;

    try {
      const { dashboardView_Organizations } =
        await this.service.getOrganizations({
          pagination: { limit: 1000, page: 0 },
          sort: {
            by: 'LAST_TOUCHPOINT',
            caseSensitive: false,
            direction: SortingDirection.Desc,
          },
        });

      const data =
        (dashboardView_Organizations?.content as Organization[]) ?? [];
      const totalElements = dashboardView_Organizations?.totalElements;

      this.load(data, {
        getId: (item) => item.metadata.id,
      });
      runInAction(() => {
        this.totalElements = totalElements;
      });
      await this.bootstrapRest();
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.load(
        mock.data.dashboardView_Organizations
          .content as unknown as Organization[],
        { getId: (data) => data.metadata.id },
      );
      this.totalElements = mock.data.dashboardView_Organizations.totalElements;

      return;
    }

    if (this.isLoading) return;

    try {
      const canHydrate = await this.checkIfCanHydrate();

      if (canHydrate) {
        this.getRecentChanges();
      } else {
        this.getAllData();
      }
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    }
  }

  async bootstrapRest() {
    let page = 1;

    while (this.totalElements > this.value.size) {
      try {
        const { dashboardView_Organizations } =
          await this.service.getOrganizations({
            pagination: { limit: 1000, page },
            sort: {
              by: 'LAST_TOUCHPOINT',
              caseSensitive: false,
              direction: SortingDirection.Desc,
            },
          });

        runInAction(() => {
          page++;
          this.load(dashboardView_Organizations?.content as Organization[], {
            getId: (data) => data.metadata.id,
          });
        });
      } catch (e) {
        runInAction(() => {
          this.error = (e as Error)?.message;
        });
        break;
      }
    }

    this.isBootstrapped = this.totalElements === this.value.size;
    this.isBootstrapping = false;
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray<T extends OrganizationStore>(
    compute: (arr: OrganizationStore[]) => T[],
  ) {
    const arr = this.toArray();

    return compute(arr);
  }

  async create(
    payload?: OrganizationInput,
    options?: { onSucces?: (serverId: string) => void },
  ) {
    const newOrganization = new OrganizationStore(
      this.root,
      this.transport,
      merge(OrganizationStore.getDefaultValue(), payload),
    );
    const tempId = newOrganization.id;
    let serverId = '';

    this.value.set(tempId, newOrganization);
    this.isLoading = true;

    try {
      const { organization_Save } = await this.service.saveOrganization({
        input: {
          website: payload?.website ?? '',
          name: payload?.name ?? 'Unnamed',
          relationship: newOrganization.value.relationship,
          stage: newOrganization.value.stage,
        },
      });

      runInAction(() => {
        serverId = organization_Save.metadata.id;

        newOrganization.setId(serverId);

        this.value.set(serverId, newOrganization);
        this.value.delete(tempId);

        this.sync({
          action: 'APPEND',
          ids: [serverId],
        });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      this.isLoading = false;

      if (serverId) {
        // Invalidate the cache after 1 second to allow the server to process the data
        // invalidating immediately would cause the server to return the organization data without
        // lastTouchpoint properties populated
        setTimeout(() => {
          this.value.get(serverId)?.invalidate();
          options?.onSucces?.(serverId);
        }, 1000);
      }
    }
  }

  async hide(ids: string[]) {
    ids.forEach((id) => {
      this.value.delete(id);
    });

    try {
      this.isLoading = true;
      await this.service.hideOrganizations({ ids });

      runInAction(() => {
        this.sync({ action: 'DELETE', ids });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async merge(
    primaryId: string,
    mergeIds: string[],
    callback?: (id: string) => void,
  ) {
    mergeIds.forEach((id) => {
      this.value.delete(id);
    });
    callback?.(primaryId);

    try {
      this.isLoading = true;
      await this.service.mergeOrganizations({
        primaryOrganizationId: primaryId,
        mergedOrganizationIds: mergeIds,
      });

      runInAction(() => {
        this.sync({ action: 'DELETE', ids: mergeIds });
        this.sync({ action: 'INVALIDATE', ids: mergeIds });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  updateTags = (ids: string[], tags: Tag[]) => {
    const tagIdsToUpdate = new Set(tags.map((tag) => tag.id));

    const shouldRemoveTags = ids.every((id) => {
      const organization = this.value.get(id);

      if (!organization) return false;

      const organizationTagIds = new Set(
        (organization.value.tags ?? []).map((tag) => tag.id),
      );

      return Array.from(tagIdsToUpdate).every((tagId) =>
        organizationTagIds.has(tagId),
      );
    });

    ids.forEach((id) => {
      const organization = this.value.get(id);

      if (!organization) return;

      if (shouldRemoveTags) {
        organization.value.tags = organization.value.tags?.filter(
          (t) => !tagIdsToUpdate.has(t.id),
        );
      } else {
        const existingIds = new Set(
          organization.value.tags?.map((t) => t.id) ?? [],
        );
        const newTags = tags.filter((t) => !existingIds.has(t.id));

        if (!Array.isArray(organization.value.tags)) {
          organization.value.tags = [];
        }

        organization.value.tags = [
          ...(organization.value.tags ?? []),
          ...newTags,
        ];

        organization.commit();
      }
    });
  };

  removeTags = (ids: string[]) => {
    ids.forEach((id) => {
      const organization = this.value.get(id);

      if (!organization) return;

      const count = organization.value.tags?.length ?? 0;

      for (let i = 0; i < count; i++) {
        organization.value.tags?.pop();
        organization.commit();
      }
    });
  };

  updateStage = (ids: string[], stage: OrganizationStage, mutate = true) => {
    let invalidCustomerStageCount = 0;

    ids.forEach((id) => {
      const organization = this.value.get(id);

      if (!organization) return;

      const currentRelationship = organization.value.relationship;
      const newDefaultRelationship = stageRelationshipMap[stage];
      const validRelationships = validRelationshipsForStage[stage];

      if (
        currentRelationship &&
        validRelationships?.includes(currentRelationship)
      ) {
        organization.value.stage = stage;
      } else if (currentRelationship === OrganizationRelationship.Customer) {
        invalidCustomerStageCount++;

        // Do not update if current relationship is Customer and new stage is not valid
      } else {
        organization.value.stage = stage;
        organization.value.relationship =
          newDefaultRelationship || organization.value.relationship;
      }

      organization.commit({ syncOnly: !mutate });
    });

    if (invalidCustomerStageCount) {
      this.root.ui.toastError(
        `${invalidCustomerStageCount} customer${
          invalidCustomerStageCount > 1 ? 's' : ''
        } remain unchanged`,
        'stage-update-failed-due-to-relationship-mismatch',
      );
    }
  };

  updateRelationship = (
    ids: string[],
    relationship: OrganizationRelationship,
    mutate = true,
  ) => {
    let invalidCustomerStageCount = 0;

    ids.forEach((id) => {
      const organization = this.value.get(id);

      if (!organization) return;

      if (
        organization.value.relationship === OrganizationRelationship.Customer &&
        ![
          OrganizationRelationship.FormerCustomer,
          OrganizationRelationship.NotAFit,
        ].includes(relationship)
      ) {
        invalidCustomerStageCount++;

        return; // Do not update if current is customer and new is not formet customer or not a fit
      }

      organization.value.relationship = relationship;
      organization.value.stage =
        relationshipStageMap[organization.value.relationship];

      organization.commit({ syncOnly: !mutate });
    });

    if (invalidCustomerStageCount) {
      this.root.ui.toastError(
        `${invalidCustomerStageCount} customer${
          invalidCustomerStageCount > 1 ? 's' : ''
        } remain unchanged`,
        'stage-update-failed-due-to-relationship-mismatch',
      );
    }
  };

  updateHealth = (
    ids: string[],
    health: OpportunityRenewalLikelihood,
    mutate = true,
  ) => {
    ids.forEach((id) => {
      const organization = this.value.get(id);

      if (!organization) return;

      set(
        organization.value,
        'accountDetails.renewalSummary.renewalLikelihood',
        health,
      );

      organization.commit({ syncOnly: !mutate });
    });
  };

  async getById(id: string) {
    try {
      this.isLoading = true;

      const { organization } = await this.service.getOrganization(id);

      if (!organization) return;

      this.load([organization as Organization], {
        getId: (d) => d.metadata.id,
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}
