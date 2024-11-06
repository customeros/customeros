import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { User, Filter, UserPage, Pagination } from '@graphql/types';

import mock from './mock.json';
import { UserStore } from './User.store';

export class UsersStore implements GroupStore<User> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, UserStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<User>();
  subscribe = makeAutoSyncableGroup.subscribe;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncableGroup(this, {
      channelName: 'Users',
      ItemStore: UserStore,
      getItemId: (user) => user.id,
    });
    makeAutoObservable(this);
  }

  get usersWithMailboxes() {
    return this.toComputedArray((users) =>
      users.filter((user) => user.value.mailboxes.length > 0),
    );
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.load(mock.data.users.content as unknown as User[]);
      this.isBootstrapped = true;
      this.totalElements = mock.data.users.totalElements;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;

      const { users } = await this.transport.graphql.request<
        USERS_QUERY_RESPONSE,
        USERS_QUERY_PAYLOAD
      >(USERS_QUERY, {
        pagination: {
          limit: 1000,
          page: 0,
        },
      });

      this.load(users.content);

      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = users.totalElements;
      });
    } catch (error) {
      runInAction(() => {
        this.error = (error as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray(compute: (arr: UserStore[]) => UserStore[]) {
    return compute(this.toArray());
  }
}

type USERS_QUERY_PAYLOAD = {
  where?: Filter;
  pagination: Pagination;
};
type USERS_QUERY_RESPONSE = {
  users: UserPage;
};
const USERS_QUERY = gql`
  query getUsers($pagination: Pagination!, $where: Filter) {
    users(pagination: $pagination, where: $where) {
      content {
        id
        firstName
        lastName
        name
        profilePhotoUrl
        mailboxes
        bot
        timezone
        hasLinkedInToken
        emails {
          email
        }
      }
      totalElements
    }
  }
`;
