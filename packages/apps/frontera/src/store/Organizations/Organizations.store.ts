import type { RootStore } from '@store/root';
import type { Transport } from '@infra/transport';

import set from 'lodash/set';
import { Store } from '@store/_store';
import { action, computed, observable, runInAction } from 'mobx';
import { OrganizationRepository } from '@infra/repositories/organization';
import { SaveOrganizationMutationVariables } from '@infra/repositories/organization/mutations/saveOrganization.generated';

import {
  relationshipStageMap,
  stageRelationshipMap,
  validRelationshipsForStage,
} from '@utils/orgStageAndRelationshipStatusMap';
import {
  OrganizationStage,
  ComparisonOperator,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import { TeamViews } from './__views__/Team.view';
import { CustomView } from './__views__/Custom.view';
import { ProfileView } from './__views__/Profile.view';
import { TargetsView } from './__views__/Targets.view';
import { CustomersView } from './__views__/Customers.view';
import { AllOrganizationsView } from './__views__/AllOrganizations.view';
import { Organization, type OrganizationDatum } from './Organization.dto';

type SaveOrganizationPayload = SaveOrganizationMutationVariables['input'];

export class OrganizationsStore extends Store<OrganizationDatum, Organization> {
  chunkSize = 500;
  private repository = OrganizationRepository.getInstance();
  @observable accessor cursors: Map<string, number> = new Map();
  @observable accessor availableCounts: Map<string, number> = new Map();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Organizations',
      getId: (data) => data?.id,
      factory: Organization,
    });

    new ProfileView(this);
    new CustomersView(this);
    new TargetsView(this);
    new AllOrganizationsView(this);
    new CustomView(this);
    new TeamViews(this);
  }

  canLoadNext(preset: string) {
    const ids = this.searchResults.get(preset);
    const cursor = this.cursors.get(preset) ?? 0;

    if (!ids) return false;

    return ids.length > this.chunkSize * (cursor + 1);
  }

  @computed
  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  public getById(id: string): Organization | null {
    if (!this.value || typeof id !== 'string') return null;

    if (!this?.value.has(id)) {
      setTimeout(() => {
        this.retrieve([id]);
      }, 0);
    }

    return this.value.get(id) as Organization;
  }

  // temporary unused
  @action
  private async _getRecentChanges() {
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
        await this.repository.getArchivedOrganizationsAfter({
          date: lastActiveAtUTC,
        });

      const { ui_organizations_search } =
        await this.repository.searchOrganizations({
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
  async search(preset: string) {
    const viewDef = this.root.tableViewDefs.getById(preset);
    const cursor = (
      this.cursors.has(preset)
        ? this.cursors.get(preset)
        : this.cursors.set(preset, 0).get(preset)
    ) as number;

    if (!viewDef) {
      console.error(`viewDef with preset=${preset} not found`);

      return;
    }

    try {
      runInAction(() => {
        if (cursor > 0) {
          // reset chunk if new search is performed
          this.cursors.set(preset, 0);
        }
        this.isLoading = true;
      });

      const payload = viewDef.toSearchPayload();

      const { ui_organizations_search: searchResult } =
        await this.repository.searchOrganizations({ ...payload });

      if (cursor === 0) {
        const ids = (searchResult?.ids ?? []).slice(
          this.chunkSize * cursor,
          this.chunkSize * cursor + this.chunkSize,
        );

        // retrieve first chunk of data after new search is performed
        await this.retrieve(ids);
      }

      runInAction(() => {
        this.isLoading = false;
        this.availableCounts.set(preset, searchResult?.totalElements);
        this.totalElements = searchResult?.totalAvailable ?? 0;
        this.searchResults.set(preset, searchResult?.ids ?? []);
        this.version++;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    }
  }

  async preload(ids: string[]) {
    ids.forEach((id) => {
      if (this.value.has(id)) return;

      this.value.set(id, new Organization(this, Organization.default({ id })));
    });

    this.retrieve(ids);
  }

  @action
  async retrieve(ids: string[]) {
    if (ids.length === 0) return;
    const invalidIds = ids.filter((id) => id.length !== 36);

    if (invalidIds.length > 0) return;

    try {
      const { ui_organizations } = await this.repository.getOrganizationsByIds({
        ids,
      });

      runInAction(() => {
        ui_organizations.forEach((raw) => {
          if (raw.hide) return;

          const foundRecord = this.value.get(raw.id);

          if (foundRecord) {
            Object.assign(foundRecord.value, raw);
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
  public async loadNext(preset: string) {
    let cursor = this.cursors.get(preset) ?? 0;

    runInAction(() => {
      cursor++;
      this.cursors.set(preset, cursor);
    });

    const ids = this.searchResults.get(preset);

    const chunkedIds = (ids ?? []).slice(
      this.chunkSize * cursor,
      this.chunkSize * cursor + this.chunkSize,
    );

    try {
      runInAction(() => {
        this.isLoading = true;
      });

      await this.retrieve(chunkedIds);
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
    if (!id) return;

    try {
      const { ui_organizations } = await this.repository.getOrganizationsByIds({
        ids: [id],
      });

      if (!ui_organizations) return;

      const data = ui_organizations[0];

      if (!data) return;

      runInAction(() => {
        const record = this.value.get(id);

        if (record) {
          Object.assign(record.value, data);
          record.version++;
        } else {
          const record = new Organization(this, data);

          this.value.set(record.id, record);
        }
      });
    } catch (e) {
      console.error('Failed invalidating company with ID: ' + id);
    }
  }

  @action
  public async create(
    payload: SaveOrganizationPayload,
    opts?: { onSucces?: (serverId: string) => void },
  ) {
    let tempId = '';

    try {
      const record = new Organization(this, Organization.default(payload));

      this.value.set(record.id, record);
      tempId = record.id;

      const { organization_Save } = await this.repository.saveOrganization({
        input: payload,
      });

      runInAction(() => {
        record.id = organization_Save.metadata.id;

        this.value.set(record.id, record);
        this.value.delete(tempId);

        this.version++;
        this.totalElements++;

        tempId = record.id;

        this.sync({
          action: 'APPEND',
          ids: [record.id],
        });
        opts?.onSucces?.(record.id);

        this.root.ui.toastSuccess('Company added', record.id);
      });
    } catch (error) {
      runInAction(() => {
        this.value.delete(tempId);
        this.root.ui.toastError(
          'Failed to create company.',
          'create-org-faillure',
        );
      });
    } finally {
      this.refreshCurrentView();
    }
  }

  @action
  public async import(globalId: number, payload?: SaveOrganizationPayload) {
    let serverId = '';

    try {
      const { organization_SaveByGlobalOrganization } =
        await this.repository.importOrganization({
          globalOrganizationId: globalId,
          input: {
            relationship: payload?.relationship,
            stage: payload?.stage,
          },
        });

      serverId = organization_SaveByGlobalOrganization?.id;

      runInAction(() => {
        if (!serverId) return;

        const record = new Organization(
          this,
          Organization.default({ id: serverId, ...(payload ?? {}) }),
        );

        this.value.set(record.id, record);
        this.version++;

        this.sync({
          action: 'APPEND',
          ids: [record.id],
        });

        this.root.ui.toastSuccess('Company added', record.id);
      });
    } catch (error) {
      runInAction(() => {
        this.root.ui.toastError(
          'Failed to add company.',
          'create-org-faillure',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
        this.invalidate(serverId);
      });
      this.refreshCurrentView();
    }
  }

  async hide(ids: string[]) {
    ids.forEach((id) => {
      this.getById(id)?.contacts.forEach((c) => {
        runInAction(() => {
          const contact = this.root.contacts.getById(c.id);

          if (!contact) return;

          contact?.draft();
          contact?.removeOrganization();
          contact?.commit({ syncOnly: true });
        });
      });

      this.value.delete(id);
    });

    try {
      this.isLoading = true;
      await this.repository.hideOrganizations({ ids });

      runInAction(() => {
        this.sync({ action: 'DELETE', ids });

        this.root.ui.toastSuccess(
          `Archived ${ids.length} ${
            ids.length === 1 ? 'company' : 'companies'
          }`,
          crypto.randomUUID(),
        );
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
        this.root.ui.toastError(
          `Failed archiving ${ids.length === 1 ? 'company' : 'companies'}`,
          crypto.randomUUID(),
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
      this.refreshCurrentView();
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
      await this.repository.mergeOrganizations({
        primaryOrganizationId: primaryId,
        mergedOrganizationIds: mergeIds,
      });

      runInAction(() => {
        this.sync({ action: 'DELETE', ids: mergeIds });
        this.sync({ action: 'INVALIDATE', ids: [primaryId] });

        this.root.ui.toastSuccess(
          `Merged ${mergeIds.length} ${
            mergeIds.length === 1 ? 'company' : 'companies'
          }`,
          primaryId,
        );
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
        this.root.ui.toastSuccess(
          `Failed merging ${mergeIds.length} ${
            mergeIds.length === 1 ? 'company' : 'companies'
          }`,
          primaryId,
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
      this.refreshCurrentView();
    }
  }

  mergeOrganizations = (primaryId: string, secondaryIds: string[]) => {
    secondaryIds.forEach((id) => {
      this.root.organizations.value.delete(id);
    });

    this.root.organizations.sync({
      action: 'DELETE',
      ids: secondaryIds,
    });

    this.root.organizations.sync({
      action: 'INVALIDATE',
      ids: [primaryId],
    });

    this.root.ui.toastSuccess(`Merged companies`, `merge-${primaryId}`);
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

  private refreshCurrentView() {
    const currentPreset = new URLSearchParams(window.location.search).get(
      'preset',
    );

    if (currentPreset) {
      this.search(currentPreset);
    }
  }
}
