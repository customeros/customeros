import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import type { Organization } from '@graphql/types';

import { OrganizationStore } from './Organization.store';

export class NewBusinessTableStore implements GroupStore<Organization> {
  value: Map<string, OrganizationStore> = new Map();
  isLoading = false;
  channel?: Channel;
  version: number = 0;
  page: number = 1;
  totalElements: number = 1;
  totalPages: number = 1;
  history: GroupOperation[] = [];
  isBootstrapped = false;
  error: string | null = null;
  subscribe = makeAutoSyncableGroup.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<Organization>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncableGroup(this, {
      channelName: 'Organizations',
      ItemStore: OrganizationStore,
      getItemId: (item) => item.metadata.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.isBootstrapped) return;

    try {
      this.isLoading = true;
      const res =
        await this.transport.graphql.request<NEW_BUSINESS_VIEW_QUERY_RESULT>(
          NEW_BUSINESS_VIEW_QUERY_DOCUMENT,
          {
            pagination: { page: this.page, limit: 80 },
          },
        );

      this.load(res?.dashboardView_Organizations.content);
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = res?.dashboardView_Organizations.totalElements;
        this.totalPages = res?.dashboardView_Organizations.totalPages;
      });
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  async loadMore() {
    this.page = this.page + 1;

    try {
      this.isLoading = true;
      const res =
        await this.transport.graphql.request<NEW_BUSINESS_VIEW_QUERY_RESULT>(
          NEW_BUSINESS_VIEW_QUERY_DOCUMENT,
          {
            pagination: { page: this.page, limit: 80 },
          },
        );
      this.load(res?.dashboardView_Organizations.content.flat());
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }
  async save() {
    // TODO: Implement save
    // this could call one or several mutations to save the data
    // operations should be group based and not per item
    // e.g. bulk update, bulk delete, create item, etc.
  }

  getById(id: string) {
    return this.value.get(id);
  }

  toArray(): OrganizationStore[] {
    return Array.from(this.value).flatMap(
      ([, organizationStore]) => organizationStore,
    );
  }
}

type NEW_BUSINESS_VIEW_QUERY_RESULT = {
  dashboardView_Organizations: {
    totalPages: number;
    totalElements: number;
    content: Organization[];
  };
};
export const NEW_BUSINESS_VIEW_QUERY_DOCUMENT = gql`
  query getOrganizationsKanban($pagination: Pagination!, $sort: SortBy) {
    dashboardView_Organizations(
      pagination: $pagination
      where: {
        AND: [
          { filter: { property: "RELATIONSHIP", value: "PROSPECT" } }
          {
            filter: {
              property: "STAGE"
              value: ["TARGET", "INTERESTED", "ENGAGED", "CLOSED_WON"]
              operation: IN
            }
          }
        ]
      }
      sort: $sort
    ) {
      content {
        name
        metadata {
          id
          created
          lastUpdated
        }
        stage
        owner {
          id
          firstName
          lastName
          name
          profilePhotoUrl
        }
      }
      totalElements
      totalAvailable
      totalPages
    }
  }
`;
