import type { RootStore } from '@store/root';
import type { Transport } from '@store/transport';

import { Store } from '@store/active-record';
import { action, computed, runInAction } from 'mobx';

import {
  SortingDirection,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

import { OrganizationDTO, type Organization } from './Organization';
import { AllOrganizationsView } from './__views__/AllOrganizations.view';
import { OrganizationsService } from './__service__/Organizations.service';

export class OrganizationsStore extends Store<Organization> {
  private service: OrganizationsService;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Organizations',
      getId: (data) => data!.metadata.id,
      factory: OrganizationDTO,
    });

    this.service = OrganizationsService.getInstance(this.transport);

    new AllOrganizationsView(this);
  }

  @computed
  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  @action
  async getRecentChanges() {
    try {
      if (this.root.demoMode || this.isBootstrapping) {
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
        });
      }

      const data =
        (dashboardView_Organizations?.content as Organization[]) ?? [];

      runInAction(() => {
        this.size = this.value.size;
        data.forEach((raw) => {
          if (!raw) return;

          const organization = OrganizationDTO.of(this.root, raw);

          this.value.set(organization.id, organization);
        });
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

  @action
  async getAllData() {
    runInAction(() => {
      this.isBootstrapping = true;
    });

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

      runInAction(() => {
        data.forEach((raw) => {
          if (!raw) return;

          const organization = OrganizationDTO.of(this.root, raw);

          this.value.set(organization.id, organization);
        });

        this.size = this.value.size;

        if (this.totalElements !== totalElements) {
          this.totalElements = totalElements;
        }
      });
      await this.bootstrapRest();
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      this.isLoading = false;
    }
  }

  @action
  async bootstrap() {
    if (this.root.demoMode) {
      // this.load(
      //   mock.data.dashboardView_Organizations
      //     .content as unknown as Organization[],
      //   { getId: (data) => data.metadata.id },
      // );
      // this.totalElements = mock.data.dashboardView_Organizations.totalElements;

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

  @action
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

        const data =
          (dashboardView_Organizations?.content as Organization[]) ?? [];

        page++;

        runInAction(() => {
          data.forEach((raw) => {
            if (!raw) return;

            const organization = OrganizationDTO.of(this.root, raw);

            this.value.set(organization.id, organization);
          });

          this.size = this.value.size;
        });
      } catch (e) {
        runInAction(() => {
          this.error = (e as Error)?.message;
        });
        break;
      }
    }

    runInAction(() => {
      this.isBootstrapped = this.totalElements === this.value.size;
      this.isBootstrapping = false;
    });
  }

  @action
  public async invalidate(id: string) {
    try {
      const { organization } = await this.service.getOrganization(id);

      if (!organization) return;

      runInAction(() => {
        this.value.set(id, OrganizationDTO.of(this.root, organization));

        if (this.active.has(id)) {
          const active = this.active.get(id);

          Object.assign(active as Organization, organization);
        }
      });
    } catch (e) {
      console.error('Failed invalidating organization with ID: ' + id);
    }
  }
}
