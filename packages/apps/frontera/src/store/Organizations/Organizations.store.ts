import type { RootStore } from '@store/root';
import type { Transport } from '@store/transport';

import set from 'lodash/set';
import { Store } from '@store/_store';
import { TagDatum } from '@store/Tags/Tag.store';
import { action, computed, observable, runInAction } from 'mobx';

import {
  relationshipStageMap,
  stageRelationshipMap,
  validRelationshipsForStage,
} from '@utils/orgStageAndRelationshipStatusMap';
import {
  SortingDirection,
  OrganizationStage,
  ComparisonOperator,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import type { SaveOrganizationMutationVariables } from './__service__/saveOrganization.generated';

import { TeamViews } from './__views__/Team.view';
import { CustomView } from './__views__/Custom.view';
import { TargetsView } from './__views__/Targets.view';
import { CustomersView } from './__views__/Customers.view';
import { AllOrganizationsView } from './__views__/AllOrganizations.view';
import { Organization, type OrganizationDatum } from './Organization.dto';
import { OrganizationsService } from './__service__/Organizations.service';

export class OrganizationsStore extends Store<OrganizationDatum, Organization> {
  private chunkSize = 50;
  private service: OrganizationsService;
  @observable accessor searchedIds: string[] = [];
  @observable accessor chunk = 0;
  @observable accessor availableCounts: Map<string, number> = new Map();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Organizations',
      getId: (data) => data?.id,
      factory: Organization,
    });

    this.service = OrganizationsService.getInstance(this.transport);

    // new CustomersView(this);
    // new TargetsView(this);
    // new AllOrganizationsView(this);
    new CustomView(this);
    // new TeamViews(this);
  }

  @computed
  get canLoadNext() {
    return this.searchedIds.length > this.chunkSize * (this.chunk + 1);
  }

  @computed
  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  @action
  async getRecentChanges() {
    try {
      if (this.isBootstrapping) {
        return;
      }

      runInAction(() => {
        this.isLoading = true;
      });

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

      const { ui_organizations_search } =
        await this.service.searchOrganizations({
          limit: this.chunkSize,
          where,
        });

      await this.retrieve(ui_organizations_search.ids);

      if (this.isHydrated) {
        await this.drop(idsToDrop);
      } else {
        await this.hydrate({
          idsToDrop,
        });
      }
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

  @action
  async getAllData() {
    runInAction(() => {
      this.isBootstrapping = true;
    });

    try {
      const urlId = (() => {
        // get organization id from url if possible
        // necessary to bootstrap the targeted organization on a profile view
        const parts = window?.location?.pathname?.split('/');

        if (parts.length !== 3 && parts[parts.length - 1].length !== 36) {
          return null;
        }

        return parts[parts.length - 1];
      })();

      const { ui_organizations_search } =
        await this.service.searchOrganizations({
          limit: this.chunkSize,
          sort: {
            by: 'ORGANIZATIONS_LAST_TOUCHPOINT',
            caseSensitive: false,
            direction: SortingDirection.Desc,
          },
        });

      if (urlId) {
        const existingIndex = ui_organizations_search.ids.indexOf(urlId);

        if (existingIndex === -1) {
          ui_organizations_search.ids.unshift(urlId);
        }
      }

      await this.retrieve(ui_organizations_search.ids);

      // const totalElements = ui_organizations_search.totalElements;

      runInAction(() => {
        this.size = this.value.size;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
        this.isBootstrapped = true;
        this.isBootstrapping = false;
      });
    }
  }

  @action
  async search(viewDefPrest: string) {
    const viewDef = this.root.tableViewDefs.getById(viewDefPrest);

    if (!viewDef) {
      console.error(`viewDef with preset=${viewDefPrest} not found`);

      return;
    }

    try {
      runInAction(() => {
        if (this.chunk > 0) {
          // reset chunk if new search is performed
          this.chunk = 0;
        }
        this.isLoading = true;
      });

      const payload = viewDef.toSearchPayload();

      const { ui_organizations_search: searchResult } =
        await this.service.searchOrganizations({ ...payload });

      if (this.chunk === 0) {
        const ids = (searchResult?.ids ?? []).slice(
          this.chunkSize * this.chunk,
          this.chunkSize * this.chunk + this.chunkSize,
        );

        // retrieve first chunk of data after new search is performed
        await this.retrieve(ids);
      }

      runInAction(() => {
        this.isLoading = false;
        this.availableCounts.set(viewDefPrest, searchResult?.totalElements);
        this.searchedIds = searchResult?.ids ?? [];
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    }
  }

  @action
  async retrieve(ids: string[]) {
    try {
      const { ui_organizations } = await this.service.getOrganizationsByIds({
        ids,
      });

      runInAction(() => {
        ui_organizations.forEach((raw) => {
          if (this.value.has(raw.id)) {
            Object.assign(raw, this.value.get(raw.id)?.value);
          } else {
            const record = new Organization(this, raw);

            this.value.set(record.id, record);
          }
        });

        this.size = this.value.size;
        this.version++;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    }
  }

  @action
  async bootstrap() {
    if (this.isLoading) return;

    try {
      // const canHydrate = await this.checkIfCanHydrate();
      //
      // if (canHydrate) {
      //   this.getRecentChanges();
      // } else {
      //   this.getAllData();
      // }
      this.getAllData();
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    }
  }

  @action
  public async loadNext() {
    if (!this.canLoadNext) return;

    runInAction(() => {
      this.chunk++;
    });

    const ids = this.searchedIds.slice(
      this.chunkSize * this.chunk,
      this.chunkSize * this.chunk + this.chunkSize,
    );

    try {
      runInAction(() => {
        this.isLoading = true;
      });

      await this.retrieve(ids);
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

  @action
  public async invalidate(id: string) {
    try {
      const { ui_organizations } = await this.service.getOrganizationsByIds({
        ids: [id],
      });

      if (!ui_organizations) return;

      const data = ui_organizations[0];

      if (!data) return;

      runInAction(() => {
        const record = this.value.get(id);

        if (record) {
          Object.assign(record.value, data);
        } else {
          const record = new Organization(this, data);

          this.value.set(record.id, record);
        }
      });
    } catch (e) {
      console.error('Failed invalidating organization with ID: ' + id);
    }
  }

  @action
  public async create(
    payload: SaveOrganizationMutationVariables['input'],
    opts?: { onSucces?: (serverId: string) => void },
  ) {
    let tempId = '';

    try {
      const record = new Organization(this, Organization.default(payload));

      this.value.set(record.id, record);
      tempId = record.id;

      const { organization_Save } = await this.service.saveOrganization({
        input: payload,
      });

      runInAction(() => {
        record.id = organization_Save.metadata.id;

        this.value.set(record.id, record);
        this.value.delete(tempId);

        this.version++;

        tempId = record.id;

        this.sync({
          action: 'APPEND',
          ids: [record.id],
        });
        opts?.onSucces?.(record.id);

        this.root.ui.toastSuccess(
          'Organization created successfully!',
          record.id,
        );
      });
    } catch (error) {
      runInAction(() => {
        this.value.delete(tempId);
        this.root.ui.toastError(
          'Failed to create organization.',
          'create-org-faillure',
        );
      });
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

        this.root.ui.toastSuccess(
          `Successfully archived ${ids.length} ${
            ids.length > 1 ? 'organizations' : 'organization'
          }`,
          crypto.randomUUID(),
        );
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
        this.root.ui.toastError(
          `Failed archiving ${
            ids.length > 1 ? 'organizations' : 'organization'
          }`,
          crypto.randomUUID(),
        );
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
        this.sync({ action: 'INVALIDATE', ids: [primaryId] });

        this.root.ui.toastSuccess(
          `Successfully merged ${mergeIds.length} ${
            mergeIds.length > 1 ? 'organizations' : 'organization'
          }`,
          primaryId,
        );
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
        this.root.ui.toastSuccess(
          `Failed merging ${mergeIds.length} ${
            mergeIds.length > 1 ? 'organizations' : 'organization'
          }`,
          primaryId,
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  updateTags = (ids: string[], tags: TagDatum[]) => {
    const tagIdsToUpdate = new Set(tags.map((tag) => tag.metadata.id));

    const shouldRemoveTags = ids.every((id) => {
      const organization = this.value.get(id);

      if (!organization) return false;

      const organizationTagIds = new Set(
        (organization.value.tags ?? []).map((tag) => tag.metadata.id),
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
          (t) => !tagIdsToUpdate.has(t.metadata.id),
        );
      } else {
        const existingIds = new Set(
          organization.value.tags?.map((t) => t.metadata.id) ?? [],
        );
        const newTags = tags.filter((t) => !existingIds.has(t.metadata.id));

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
}
