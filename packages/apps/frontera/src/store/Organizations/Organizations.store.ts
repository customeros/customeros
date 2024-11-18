import type { RootStore } from '@store/root';
import type { Transport } from '@store/transport';

import { Store } from '@store/_store';
import { action, computed, runInAction } from 'mobx';

import { SortingDirection, ComparisonOperator } from '@graphql/types';

import { AllOrganizationsView } from './__views__/AllOrganizations.view';
import { Organization, type OrganizationDatum } from './Organization.dto';
import { OrganizationsService } from './__service__/Organizations.service';

export class OrganizationsStore extends Store<Organization> {
  private service: OrganizationsService;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Organizations',
      getId: (data) => data?.metadata?.id,
      factory: Organization,
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
        (dashboardView_Organizations?.content as OrganizationDatum[]) ?? [];

      runInAction(() => {
        this.size = this.value.size;
        data.forEach((raw) => {
          if (!raw) return;

          const record = new Organization(this, raw);

          this.value.set(record.id, record);
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

      const data = dashboardView_Organizations?.content ?? [];
      const totalElements = dashboardView_Organizations?.totalElements;

      runInAction(() => {
        data.forEach((raw) => {
          if (!raw) return;

          const record = new Organization(this, raw);

          this.value.set(record.id, record);
        });

        this.size = this.value.size;

        if (this.totalElements !== totalElements) {
          this.totalElements = totalElements;
        }
      });
      // await this.bootstrapRest();
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

        const data = dashboardView_Organizations?.content ?? [];

        page++;

        runInAction(() => {
          data.forEach((raw) => {
            if (!raw) return;

            const record = new Organization(this, raw);

            this.value.set(record.id, record);
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
      const { organization: raw } = await this.service.getOrganization(id);

      if (!raw) return;

      runInAction(() => {
        const record = this.value.get(id);

        if (record) {
          Object.assign(record, raw);
        }
      });
    } catch (e) {
      console.error('Failed invalidating organization with ID: ' + id);
    }
  }
}
