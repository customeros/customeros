import { RootStore } from '@store/root';
import { Store } from '@store/active-record';
import { Transport } from '@store/transport';
import { action, override, computed, observable, runInAction } from 'mobx';

import {
  SortingDirection,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

import { OrganizationsService } from './__service__/Organizations.service';
import { OrganizationQuery } from './__service__/getOrganization.generated';

type Organization = OrganizationQuery['organization'];

export class OrganizationsStore extends Store<Organization> {
  @observable private accessor totalElements = 0;
  private service: OrganizationsService;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, { name: 'Organizations' });

    this.service = OrganizationsService.getInstance(this.transport);
  }

  @computed
  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  @override
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
        });
      }

      const data =
        (dashboardView_Organizations?.content as Organization[]) ?? [];
      const totalElements = dashboardView_Organizations?.totalElements;

      runInAction(() => {
        data.forEach((item) => {
          if (!item) return;

          this.value.set(item?.metadata?.id, item);
        });
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

  @action
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

      runInAction(() => {
        data.forEach((item) => {
          if (!item) return;

          this.value.set(item?.metadata?.id, item);
        });
        this.totalElements = totalElements;
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
        const totalElements = dashboardView_Organizations?.totalElements;

        runInAction(() => {
          page++;

          data.forEach((item) => {
            if (!item) return;

            this.value.set(item?.metadata?.id, item);
          });

          this.totalElements = totalElements;
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
}
